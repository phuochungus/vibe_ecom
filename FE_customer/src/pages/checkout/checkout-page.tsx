import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Steps, Button, Space, Typography, Form, message, Spin } from 'antd'
import { useMutation } from '@tanstack/react-query'
import { useCart } from '@/lib/cart'
import { ordersApi } from '@/services/orders'
import { paymentsApi } from '@/services/payments'
import { getErrorMessage } from '@/lib/api'
import type { PaymentMethod } from '@/types'
import AddressForm, { type AddressFormValues } from '@/components/checkout/address-form'
import OrderReview from '@/components/checkout/order-review'

const { Title } = Typography

export default function CheckoutPage() {
  const navigate = useNavigate()
  const { items, totalAmount, clearCart } = useCart()
  const [currentStep, setCurrentStep] = useState(0)
  const [addressForm] = Form.useForm<AddressFormValues>()
  const [addressData, setAddressData] = useState<AddressFormValues | null>(null)
  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>('COD')

  const shippingFee = 30000
  const total = totalAmount + shippingFee

  const createOrderMutation = useMutation({
    mutationFn: async (data: AddressFormValues) => {
      // Generate idempotency key (BR-PAY-04)
      const idempotencyKey = crypto.randomUUID()

      const order = await ordersApi.create(
        {
          items: items.map((item) => ({
            product_id: item.product_id,
            quantity: item.quantity,
          })),
          payment_method: paymentMethod,
          shipping_recipient_name: data.recipient_name,
          shipping_phone: data.phone,
          shipping_line1: data.address_line1,
          shipping_line2: data.address_line2,
          shipping_ward: data.ward,
          shipping_district: data.district,
          shipping_city: data.city,
          shipping_province: data.province,
          shipping_country_code: 'VN',
          customer_note: data.note,
        },
        idempotencyKey
      )

      return order
    },
  })

  const initiatePaymentMutation = useMutation({
    mutationFn: async (orderId: string) => {
      const idempotencyKey = crypto.randomUUID()
      return await paymentsApi.initiate(orderId, idempotencyKey)
    },
  })

  const handleNext = async () => {
    if (currentStep === 0) {
      // Validate address form
      try {
        const values = await addressForm.validateFields()
        setAddressData(values)
        setCurrentStep(1)
      } catch (error) {
        message.error('Vui lòng điền đầy đủ thông tin giao hàng')
      }
    } else if (currentStep === 1) {
      // Create order and handle payment based on selected method
      if (!addressData) return

      try {
        message.loading({ content: 'Đang tạo đơn hàng...', key: 'checkout' })

        const order = await createOrderMutation.mutateAsync(addressData)

        // Clear cart after order created
        clearCart()

        if (paymentMethod === 'COD') {
          message.success({ content: 'Đặt hàng thành công!', key: 'checkout' })
          navigate(`/checkout/return?orderId=${order.id}&status=COD`)
        } else {
          // MOMO: initiate payment and redirect
          message.loading({ content: 'Đang chuyển đến trang thanh toán...', key: 'checkout' })

          const payment = await initiatePaymentMutation.mutateAsync(order.id)

          // Redirect to payment gateway
          window.location.href = payment.checkout_url
        }
      } catch (error) {
        message.error({ content: getErrorMessage(error), key: 'checkout' })
      }
    }
  }

  const handleBack = () => {
    setCurrentStep(currentStep - 1)
  }

  if (items.length === 0) {
    return (
      <div style={{ padding: 24, textAlign: 'center' }}>
        <Title level={3}>Giỏ hàng trống</Title>
        <Button type="primary" onClick={() => navigate('/products')}>
          Tiếp tục mua sắm
        </Button>
      </div>
    )
  }

  const isLoading = createOrderMutation.isPending || initiatePaymentMutation.isPending

  return (
    <div style={{ padding: 24, maxWidth: 800, margin: '0 auto' }}>
      <Title level={2}>Thanh toán</Title>

      <Steps
        current={currentStep}
        items={[
          { title: 'Địa chỉ giao hàng' },
          { title: 'Xác nhận đơn hàng' },
        ]}
        style={{ marginBottom: 32 }}
      />

      <Spin spinning={isLoading}>
        <div style={{ minHeight: 400 }}>
          {currentStep === 0 && <AddressForm form={addressForm as any} />}

          {currentStep === 1 && addressData && (
            <OrderReview
              items={items}
              address={addressData}
              subtotal={totalAmount}
              shippingFee={shippingFee}
              total={total}
              paymentMethod={paymentMethod}
              onPaymentMethodChange={setPaymentMethod}
            />
          )}
        </div>

        <div style={{ marginTop: 24, display: 'flex', justifyContent: 'space-between' }}>
          <Button size="large" onClick={() => navigate('/cart')} disabled={isLoading}>
            Quay lại giỏ hàng
          </Button>

          <Space>
            {currentStep > 0 && (
              <Button size="large" onClick={handleBack} disabled={isLoading}>
                Quay lại
              </Button>
            )}
            <Button type="primary" size="large" onClick={handleNext} loading={isLoading}>
              {currentStep === 0 ? 'Tiếp tục' : paymentMethod === 'COD' ? 'Đặt hàng' : 'Thanh toán'}
            </Button>
          </Space>
        </div>
      </Spin>
    </div>
  )
}

