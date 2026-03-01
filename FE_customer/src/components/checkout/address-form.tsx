import { Form, Input, Row, Col } from 'antd'

export interface AddressFormValues {
  recipient_name: string
  phone: string
  address_line1: string
  address_line2?: string
  ward?: string
  district?: string
  city: string
  province?: string
  note?: string
}

interface AddressFormProps {
  form: ReturnType<typeof Form.useForm>[0]
}

export default function AddressForm({ form }: AddressFormProps) {
  return (
    <Form form={form} layout="vertical" size="large">
      <Row gutter={16}>
        <Col xs={24} md={12}>
          <Form.Item
            name="recipient_name"
            label="Họ và tên người nhận"
            rules={[{ required: true, message: 'Vui lòng nhập họ tên' }]}
          >
            <Input placeholder="Nguyễn Văn A" />
          </Form.Item>
        </Col>
        <Col xs={24} md={12}>
          <Form.Item
            name="phone"
            label="Số điện thoại"
            rules={[
              { required: true, message: 'Vui lòng nhập số điện thoại' },
              { pattern: /^[0-9]{10}$/, message: 'Số điện thoại không hợp lệ' },
            ]}
          >
            <Input placeholder="0901234567" />
          </Form.Item>
        </Col>
      </Row>

      <Form.Item
        name="address_line1"
        label="Địa chỉ"
        rules={[{ required: true, message: 'Vui lòng nhập địa chỉ' }]}
      >
        <Input placeholder="Số nhà, tên đường" />
      </Form.Item>

      <Form.Item name="address_line2" label="Địa chỉ chi tiết (tùy chọn)">
        <Input placeholder="Tòa nhà, số phòng..." />
      </Form.Item>

      <Row gutter={16}>
        <Col xs={24} md={8}>
          <Form.Item name="ward" label="Phường/Xã">
            <Input placeholder="Phường 1" />
          </Form.Item>
        </Col>
        <Col xs={24} md={8}>
          <Form.Item name="district" label="Quận/Huyện">
            <Input placeholder="Quận 1" />
          </Form.Item>
        </Col>
        <Col xs={24} md={8}>
          <Form.Item
            name="city"
            label="Thành phố"
            rules={[{ required: true, message: 'Vui lòng nhập thành phố' }]}
          >
            <Input placeholder="TP. Hồ Chí Minh" />
          </Form.Item>
        </Col>
      </Row>

      <Form.Item name="province" label="Tỉnh/Thành phố">
        <Input placeholder="TP. Hồ Chí Minh" />
      </Form.Item>

      <Form.Item name="note" label="Ghi chú đơn hàng (tùy chọn)">
        <Input.TextArea rows={3} placeholder="Ghi chú về đơn hàng, ví dụ: giao hàng giờ hành chính" />
      </Form.Item>
    </Form>
  )
}
