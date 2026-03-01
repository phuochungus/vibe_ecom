import { Card, Descriptions, Button, Divider } from 'antd'
import { useNavigate } from 'react-router-dom'
import { formatCurrency } from '@/lib/utils'
import { useCart } from '@/lib/cart'
import { useAuth } from '@/lib/auth'

export default function CartSummary() {
  const navigate = useNavigate()
  const { totalAmount, itemCount } = useCart()
  const { isAuthenticated } = useAuth()

  const shippingFee = 30000 // Fixed shipping fee for now
  const total = totalAmount + shippingFee

  const handleCheckout = () => {
    if (!isAuthenticated) {
      navigate('/login?redirect=/checkout')
      return
    }
    navigate('/checkout')
  }

  return (
    <Card title="Tổng đơn hàng">
      <Descriptions column={1} size="small">
        <Descriptions.Item label="Tạm tính">
          {formatCurrency(totalAmount)}
        </Descriptions.Item>
        <Descriptions.Item label="Phí vận chuyển">
          {formatCurrency(shippingFee)}
        </Descriptions.Item>
      </Descriptions>

      <Divider style={{ margin: '12px 0' }} />

      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <span style={{ fontSize: 16, fontWeight: 600 }}>Tổng cộng:</span>
        <span className="price" style={{ fontSize: 20, fontWeight: 700 }}>
          {formatCurrency(total)}
        </span>
      </div>

      <Button
        type="primary"
        size="large"
        block
        onClick={handleCheckout}
        disabled={itemCount === 0}
      >
        Tiến hành thanh toán ({itemCount} sản phẩm)
      </Button>
    </Card>
  )
}
