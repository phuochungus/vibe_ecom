import { Typography, Empty, Row, Col, Space, Button } from 'antd'
import { ShoppingOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { useCart } from '@/lib/cart'
import CartItem from '@/components/cart/cart-item'
import CartSummary from '@/components/cart/cart-summary'

const { Title } = Typography

export default function CartPage() {
  const navigate = useNavigate()
  const { items } = useCart()

  if (items.length === 0) {
    return (
      <div style={{ padding: 24, textAlign: 'center', minHeight: 400 }}>
        <Empty
          description="Giỏ hàng trống"
          image={Empty.PRESENTED_IMAGE_SIMPLE}
        >
          <Button type="primary" icon={<ShoppingOutlined />} onClick={() => navigate('/products')}>
            Tiếp tục mua sắm
          </Button>
        </Empty>
      </div>
    )
  }

  return (
    <div style={{ padding: 24, maxWidth: 1200, margin: '0 auto' }}>
      <Title level={2} data-testid="cart-title">Giỏ hàng của bạn</Title>

      <Row gutter={[24, 24]}>
        {/* Cart Items */}
        <Col xs={24} lg={16}>
          <Space direction="vertical" style={{ width: '100%' }}>
            {items.map((item) => (
              <CartItem key={item.product_id} item={item} />
            ))}
          </Space>
        </Col>

        {/* Summary */}
        <Col xs={24} lg={8}>
          <div style={{ position: 'sticky', top: 88 }}>
            <CartSummary />
          </div>
        </Col>
      </Row>
    </div>
  )
}
