import { useNavigate, useSearchParams } from 'react-router-dom'
import { Result, Button } from 'antd'
import { CloseCircleOutlined } from '@ant-design/icons'

export default function PaymentCancelPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()

  const orderId = searchParams.get('orderId')

  return (
    <div style={{ padding: 48 }}>
      <Result
        icon={<CloseCircleOutlined style={{ color: '#ff4d4f' }} />}
        status="error"
        title="Thanh toán bị hủy"
        subTitle={`Bạn đã hủy thanh toán. Đơn hàng vẫn được lưu và chờ thanh toán.`}
        extra={[
          orderId && (
            <Button type="primary" key="retry" onClick={() => navigate(`/orders/${orderId}`)}>
              Thử lại thanh toán
            </Button>
          ),
          <Button key="orders" onClick={() => navigate('/orders')}>
            Xem đơn hàng
          </Button>,
          <Button key="continue" onClick={() => navigate('/products')}>
            Tiếp tục mua sắm
          </Button>,
        ].filter(Boolean)}
      />
    </div>
  )
}
