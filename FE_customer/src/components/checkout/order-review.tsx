import { Card, Descriptions, List, Space, Divider, Typography, Radio } from 'antd'
import { DollarOutlined, WalletOutlined } from '@ant-design/icons'
import type { CartItem, PaymentMethod } from '@/types'
import { formatCurrency } from '@/lib/utils'
import type { AddressFormValues } from './address-form'

const { Text } = Typography

interface OrderReviewProps {
  items: CartItem[]
  address: AddressFormValues
  subtotal: number
  shippingFee: number
  total: number
  paymentMethod: PaymentMethod
  onPaymentMethodChange: (method: PaymentMethod) => void
}

export default function OrderReview({ items, address, subtotal, shippingFee, total, paymentMethod, onPaymentMethodChange }: OrderReviewProps) {
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

      {/* Payment Method */}
      <Card title="Hình thức thanh toán" size="small">
        <Radio.Group
          value={paymentMethod}
          onChange={(e) => onPaymentMethodChange(e.target.value)}
          style={{ width: '100%' }}
        >
          <Space direction="vertical" style={{ width: '100%' }}>
            <Radio value="COD" style={{ padding: '8px 0' }}>
              <Space>
                <DollarOutlined style={{ fontSize: 20, color: '#52c41a' }} />
                <div>
                  <Text strong>Thanh toán khi nhận hàng (COD)</Text>
                  <br />
                  <Text type="secondary" style={{ fontSize: 12 }}>
                    Thanh toán bằng tiền mặt khi nhận hàng
                  </Text>
                </div>
              </Space>
            </Radio>
            <Radio value="PAYOS" style={{ padding: '8px 0' }}>
              <Space>
                <WalletOutlined style={{ fontSize: 20, color: '#a50064' }} />
                <div>
                  <Text strong>Thanh toán online PAYOS</Text>
                  <br />
                  <Text type="secondary" style={{ fontSize: 12 }}>
                    Thanh toán tự động qua cổng thanh toán online
                  </Text>
                </div>
              </Space>
            </Radio>
          </Space>
        </Radio.Group>
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
