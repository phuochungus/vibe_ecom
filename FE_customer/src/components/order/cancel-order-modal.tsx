import { Modal, Form, Input, message } from 'antd'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { ordersApi } from '@/services/orders'
import { getErrorMessage } from '@/lib/api'

interface CancelOrderModalProps {
  orderId: string
  open: boolean
  onClose: () => void
}

export default function CancelOrderModal({ orderId, open, onClose }: CancelOrderModalProps) {
  const [form] = Form.useForm()
  const queryClient = useQueryClient()

  const cancelMutation = useMutation({
    mutationFn: async (reason: string) => {
      return await ordersApi.cancel(orderId, reason)
    },
    onSuccess: () => {
      message.success('Đã hủy đơn hàng')
      queryClient.invalidateQueries({ queryKey: ['orders'] })
      queryClient.invalidateQueries({ queryKey: ['order', orderId] })
      onClose()
      form.resetFields()
    },
    onError: (error) => {
      message.error(getErrorMessage(error))
    },
  })

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      await cancelMutation.mutateAsync(values.reason)
    } catch (error) {
      // Form validation error
    }
  }

  const handleCancel = () => {
    onClose()
    form.resetFields()
  }

  return (
    <Modal
      title="Hủy đơn hàng"
      open={open}
      onOk={handleSubmit}
      onCancel={handleCancel}
      confirmLoading={cancelMutation.isPending}
      okText="Xác nhận hủy"
      cancelText="Đóng"
      okButtonProps={{ danger: true }}
    >
      <Form form={form} layout="vertical">
        <Form.Item
          name="reason"
          label="Lý do hủy đơn hàng"
          rules={[
            { required: true, message: 'Vui lòng nhập lý do hủy' },
            { min: 10, message: 'Lý do phải có ít nhất 10 ký tự' },
          ]}
        >
          <Input.TextArea
            rows={4}
            placeholder="Vui lòng cho chúng tôi biết lý do bạn muốn hủy đơn hàng này..."
          />
        </Form.Item>
      </Form>
    </Modal>
  )
}
