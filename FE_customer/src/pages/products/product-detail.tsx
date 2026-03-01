import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import {
  Row,
  Col,
  Spin,
  Typography,
  Space,
  Button,
  InputNumber,
  Tag,
  Divider,
  message,
  Breadcrumb,
} from 'antd'
import { ShoppingCartOutlined, HomeOutlined } from '@ant-design/icons'
import { useQuery } from '@tanstack/react-query'
import { productsApi } from '@/services/products'
import { formatCurrency, formatProductStatus, getProductStatusColor } from '@/lib/utils'
import { useCart } from '@/lib/cart'

const { Title, Paragraph, Text } = Typography

export default function ProductDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { addItem } = useCart()
  const [quantity, setQuantity] = useState(1)

  const { data: product, isLoading } = useQuery({
    queryKey: ['product', id],
    queryFn: () => productsApi.getById(id!),
    enabled: !!id,
  })

  const handleAddToCart = () => {
    if (!product) return
    try {
      addItem(product, quantity)
      message.success(`Đã thêm ${quantity} sản phẩm vào giỏ hàng`)
    } catch (error) {
      message.error(error instanceof Error ? error.message : 'Không thể thêm vào giỏ hàng')
    }
  }

  const handleBuyNow = () => {
    if (!product) return
    try {
      addItem(product, quantity)
      navigate('/checkout')
    } catch (error) {
      message.error(error instanceof Error ? error.message : 'Không thể mua ngay')
    }
  }

  if (isLoading) {
    return (
      <div style={{ textAlign: 'center', padding: 48 }}>
        <Spin size="large" />
      </div>
    )
  }

  if (!product) {
    return (
      <div style={{ textAlign: 'center', padding: 48 }}>
        <Text>Không tìm thấy sản phẩm</Text>
      </div>
    )
  }

  const isOutOfStock = product.stock === 0
  const isDiscontinued = product.status === 'DISCONTINUED'
  const isUnavailable = isOutOfStock || isDiscontinued

  return (
    <div style={{ padding: 24, maxWidth: 1200, margin: '0 auto' }}>
      {/* Breadcrumb */}
      <Breadcrumb
        items={[
          { title: <HomeOutlined />, href: '/' },
          { title: 'Sản phẩm', href: '/products' },
          { title: product.name },
        ]}
        style={{ marginBottom: 24 }}
      />

      <Row gutter={[48, 24]}>
        {/* Image */}
        <Col xs={24} md={10}>
          <div
            style={{
              width: '100%',
              paddingBottom: '100%',
              position: 'relative',
              background: '#f5f5f5',
              borderRadius: 8,
              overflow: 'hidden',
            }}
          >
            <img
              src={product.image_url || 'https://via.placeholder.com/500?text=No+Image'}
              alt={product.name}
              style={{
                position: 'absolute',
                top: 0,
                left: 0,
                width: '100%',
                height: '100%',
                objectFit: 'cover',
              }}
            />
          </div>
        </Col>

        {/* Details */}
        <Col xs={24} md={14}>
          <Space direction="vertical" size="large" style={{ width: '100%' }}>
            {/* Title & Status */}
            <div>
              <Title level={2} style={{ marginBottom: 8 }}>
                {product.name}
              </Title>
              <Space>
                <Text type="secondary">SKU: {product.sku}</Text>
                <Tag color={getProductStatusColor(product.status)}>
                  {formatProductStatus(product.status)}
                </Tag>
              </Space>
            </div>

            <Divider style={{ margin: 0 }} />

            {/* Price */}
            <div>
              <div className="price" style={{ fontSize: 32, fontWeight: 700 }}>
                {formatCurrency(product.price_cents)}
              </div>
            </div>

            {/* Stock */}
            <div>
              <Space>
                <Text>Tình trạng:</Text>
                {isOutOfStock ? (
                  <Tag color="error">Hết hàng</Tag>
                ) : product.stock < 10 ? (
                  <Tag color="warning">Chỉ còn {product.stock} sản phẩm</Tag>
                ) : (
                  <Tag color="success">Còn hàng</Tag>
                )}
              </Space>
            </div>

            {/* Quantity */}
            {!isUnavailable && (
              <div>
                <Space>
                  <Text>Số lượng:</Text>
                  <InputNumber
                    min={1}
                    max={product.stock}
                    value={quantity}
                    onChange={(val) => setQuantity(val || 1)}
                    style={{ width: 120 }}
                  />
                </Space>
              </div>
            )}

            {/* Actions */}
            <Space size="large">
              <Button
                type="primary"
                size="large"
                icon={<ShoppingCartOutlined />}
                onClick={handleAddToCart}
                disabled={isUnavailable}
              >
                Thêm vào giỏ hàng
              </Button>
              <Button
                size="large"
                onClick={handleBuyNow}
                disabled={isUnavailable}
                style={{ background: '#EE4D2D', color: '#fff', borderColor: '#EE4D2D' }}
              >
                Mua ngay
              </Button>
            </Space>

            <Divider />

            {/* Description */}
            {product.description && (
              <div>
                <Title level={4}>Mô tả sản phẩm</Title>
                <Paragraph style={{ whiteSpace: 'pre-wrap' }}>
                  {product.description}
                </Paragraph>
              </div>
            )}
          </Space>
        </Col>
      </Row>
    </div>
  )
}
