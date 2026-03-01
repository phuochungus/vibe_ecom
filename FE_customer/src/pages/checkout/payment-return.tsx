import { useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { Result, Button, Spin } from 'antd'
import { CheckCircleOutlined } from '@ant-design/icons'

export default function PaymentReturnPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()

  const orderId = searchParams.get('orderId')
  const status = searchParams.get('status')

  useEffect(() => {
    // In production, you might want to verify payment status with backend here
    if (!orderId || (status !== 'PAID' && status !== 'COD')) {
      // Invalid payment result
      setTimeout(() => navigate('/orders'), 3000)
    }
  }, [orderId, status, navigate])

  if (!orderId) {
    return (
      <div style={{ padding: 48, textAlign: 'center' }}>
        <Spin size="large" />
      </div>
    )
  }

  if (status === 'COD') {
    return (
      <div style={{ padding: 48 }}>
        <Result
          icon={<CheckCircleOutlined style={{ color: '#52c41a' }} />}
          status="success"
          title="Đặt hàng thành công!"
          subTitle="Đơn hàng của bạn đã được tạo. Bạn sẽ thanh toán khi nhận hàng (COD)."
          extra={[
            <Button type="primary" key="orders" onClick={() => navigate('/orders')}>
              Xem đơn hàng
            </Button>,
            <Button key="continue" onClick={() => navigate('/products')}>
              Tiếp tục mua sắm
            </Button>,
          ]}
        />
      </div>
    )
  }

  if (status === 'PAID') {
    return (
      <div style={{ padding: 48 }}>
        <Result
          icon={<CheckCircleOutlined style={{ color: '#52c41a' }} />}
          status="success"
          title="Thanh toán thành công!"
          subTitle={`Đơn hàng của bạn đã được thanh toán. Chúng tôi sẽ xử lý và giao hàng sớm nhất có thể.`}
          extra={[
            <Button type="primary" key="orders" onClick={() => navigate('/orders')}>
              Xem đơn hàng
            </Button>,
            <Button key="continue" onClick={() => navigate('/products')}>
              Tiếp tục mua sắm
            </Button>,
          ]}
        />
      </div>
    )
  }

  return (
    <div style={{ padding: 48 }}>
      <Result
        status="warning"
        title="Thanh toán chưa hoàn tất"
        subTitle="Vui lòng kiểm tra trạng thái đơn hàng hoặc thử lại"
        extra={[
          <Button type="primary" key="orders" onClick={() => navigate('/orders')}>
            Xem đơn hàng
          </Button>,
        ]}
      />
    </div>
  )
}
