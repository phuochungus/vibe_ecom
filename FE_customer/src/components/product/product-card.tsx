import { Card, Button, Badge, Space, Tag } from 'antd'
import { ShoppingCartOutlined, EyeOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import type { Product } from '@/types'
import { formatCurrency } from '@/lib/utils'
import { useCart } from '@/lib/cart'
import { message } from 'antd'

const { Meta } = Card

interface ProductCardProps {
  product: Product
}

export default function ProductCard({ product }: ProductCardProps) {
  const navigate = useNavigate()
  const { addItem } = useCart()

  const handleAddToCart = (e: React.MouseEvent) => {
    e.stopPropagation()
    try {
      addItem(product, 1)
      message.success('Đã thêm vào giỏ hàng')
    } catch (error) {
      message.error(error instanceof Error ? error.message : 'Không thể thêm vào giỏ hàng')
    }
  }

  const handleViewDetail = () => {
    navigate(`/products/${product.id}`)
  }

  const isOutOfStock = product.stock === 0
  const isDiscontinued = product.status === 'DISCONTINUED'
  const isUnavailable = isOutOfStock || isDiscontinued

  return (
    <Badge.Ribbon
      text={isDiscontinued ? 'Ngừng bán' : isOutOfStock ? 'Hết hàng' : undefined}
      color={isDiscontinued ? 'red' : 'volcano'}
      style={{ display: isUnavailable ? 'block' : 'none' }}
    >
      <Card
        hoverable
        className="product-card"
        cover={
          <div style={{ height: 240, overflow: 'hidden', background: '#f5f5f5' }}>
            <img
              alt={product.name}
              src={product.image_url || 'https://via.placeholder.com/300x240?text=No+Image'}
              style={{ width: '100%', height: '100%', objectFit: 'cover' }}
              loading="lazy"
            />
          </div>
        }
        actions={[
          <Button
            key="view"
            type="text"
            icon={<EyeOutlined />}
            onClick={handleViewDetail}
            data-testid={`product-view-${product.id}`}
          >
            Xem chi tiết
          </Button>,
          <Button
            key="cart"
            type="primary"
            icon={<ShoppingCartOutlined />}
            onClick={handleAddToCart}
            disabled={isUnavailable}
            data-testid={`product-add-${product.id}`}
          >
            Thêm vào giỏ
          </Button>,
        ]}
        onClick={handleViewDetail}
      >
        <Meta
          title={
            <div style={{ fontSize: 16, height: 48, overflow: 'hidden', textOverflow: 'ellipsis' }}>
              {product.name}
            </div>
          }
          description={
            <Space direction="vertical" style={{ width: '100%' }}>
              <div className="price" style={{ fontSize: 20, fontWeight: 600 }}>
                {formatCurrency(product.price)}
              </div>
              {product.stock > 0 && product.stock < 10 && (
                <Tag color="warning">Chỉ còn {product.stock} sản phẩm</Tag>
              )}
            </Space>
          }
        />
      </Card>
    </Badge.Ribbon>
  )
}
