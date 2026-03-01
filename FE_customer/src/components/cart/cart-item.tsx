import { Card, Row, Col, InputNumber, Button, Typography } from 'antd'
import { DeleteOutlined } from '@ant-design/icons'
import type { CartItem as CartItemType } from '@/types'
import { formatCurrency } from '@/lib/utils'
import { useCart } from '@/lib/cart'
import { message } from 'antd'

const { Text } = Typography

interface CartItemProps {
  item: CartItemType
}

export default function CartItem({ item }: CartItemProps) {
  const { updateQuantity, removeItem } = useCart()

  const handleQuantityChange = (value: number | null) => {
    if (!value) return
    try {
      updateQuantity(item.product_id, value)
    } catch (error) {
      message.error(error instanceof Error ? error.message : 'Không thể cập nhật số lượng')
    }
  }

  const handleRemove = () => {
    removeItem(item.product_id)
    message.success('Đã xóa sản phẩm khỏi giỏ hàng')
  }

  const lineTotal = parseFloat(item.product.price_cents) * item.quantity

  return (
    <Card size="small" style={{ marginBottom: 12 }}>
      <Row gutter={16} align="middle">
        {/* Image */}
        <Col xs={24} sm={6}>
          <div style={{ width: '100%', paddingBottom: '100%', position: 'relative', background: '#f5f5f5', borderRadius: 4 }}>
            <img
              src={item.product.image_url || 'https://via.placeholder.com/100?text=No+Image'}
              alt={item.product.name}
              style={{
                position: 'absolute',
                top: 0,
                left: 0,
                width: '100%',
                height: '100%',
                objectFit: 'cover',
                borderRadius: 4,
              }}
            />
          </div>
        </Col>

        {/* Info */}
        <Col xs={24} sm={18}>
          <Row gutter={[16, 8]} align="middle">
            <Col xs={24} md={12}>
              <div style={{ fontWeight: 500, marginBottom: 4 }}>{item.product.name}</div>
              <Text type="secondary" style={{ fontSize: 12 }}>SKU: {item.product.sku}</Text>
            </Col>

            <Col xs={8} md={4}>
              <div className="price">{formatCurrency(item.product.price_cents)}</div>
            </Col>

            <Col xs={8} md={4}>
              <InputNumber
                min={1}
                max={item.product.stock}
                value={item.quantity}
                onChange={handleQuantityChange}
                style={{ width: '100%' }}
              />
            </Col>

            <Col xs={6} md={3}>
              <div className="price" style={{ fontWeight: 600 }}>
                {formatCurrency(lineTotal)}
              </div>
            </Col>

            <Col xs={2} md={1}>
              <Button
                type="text"
                danger
                icon={<DeleteOutlined />}
                onClick={handleRemove}
              />
            </Col>
          </Row>
        </Col>
      </Row>
    </Card>
  )
}
