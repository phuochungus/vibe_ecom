import { Card, Descriptions, List, Space, Divider, Typography } from 'antd'
import type { CartItem } from '@/types'
import { formatCurrency } from '@/lib/utils'
import type { AddressFormValues } from './address-form'

const { Text } = Typography

interface OrderReviewProps {
  items: CartItem[]
  address: AddressFormValues
  subtotal: number
  shippingFee: number
  total: number
}

export default function OrderReview({ items, address, subtotal, shippingFee, total }: OrderReviewProps) {
  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      {/* Shipping Address */}
      <Card title="Địa chỉ giao hàng" size="small">
        <Descriptions column={1} size="small">
          <Descriptions.Item label="Người nhận">{address.recipient_name}</Descriptions.Item>
          <Descriptions.Item label="Số điện thoại">{address.phone}</Descriptions.Item>
          <Descriptions.Item label="Địa chỉ">
            {[
              address.address_line1,
              address.address_line2,
              address.ward,
              address.district,
              address.city,
              address.province,
            ]
              .filter(Boolean)
              .join(', ')}
          </Descriptions.Item>
          {address.note && <Descriptions.Item label="Ghi chú">{address.note}</Descriptions.Item>}
        </Descriptions>
      </Card>

      {/* Order Items */}
      <Card title="Sản phẩm" size="small">
        <List
          dataSource={items}
          renderItem={(item) => (
            <List.Item>
              <List.Item.Meta
                avatar={
                  <img
                    src={item.product.image_url || 'https://via.placeholder.com/60?text=No+Image'}
                    alt={item.product.name}
                    style={{ width: 60, height: 60, objectFit: 'cover', borderRadius: 4 }}
                  />
                }
                title={item.product.name}
                description={`Số lượng: ${item.quantity}`}
              />
              <div className="price">{formatCurrency(parseFloat(item.product.price) * item.quantity)}</div>
            </List.Item>
          )}
        />
      </Card>

      {/* Order Summary */}
      <Card title="Tổng đơn hàng" size="small">
        <Space direction="vertical" style={{ width: '100%' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <Text>Tạm tính:</Text>
            <Text>{formatCurrency(subtotal)}</Text>
          </div>
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <Text>Phí vận chuyển:</Text>
            <Text>{formatCurrency(shippingFee)}</Text>
          </div>
          <Divider style={{ margin: '8px 0' }} />
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <Text style={{ fontSize: 16, fontWeight: 600 }}>Tổng cộng:</Text>
            <Text className="price" style={{ fontSize: 20, fontWeight: 700 }}>
              {formatCurrency(total)}
            </Text>
          </div>
        </Space>
      </Card>
    </Space>
  )
}
