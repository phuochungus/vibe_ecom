import { useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  Typography,
  Card,
  Descriptions,
  List,
  Tag,
  Space,
  Button,
  Spin,
  Breadcrumb,
  Row,
  Col,
  Alert,
} from 'antd'
import { HomeOutlined } from '@ant-design/icons'
import { useQuery } from '@tanstack/react-query'
import { ordersApi } from '@/services/orders'
import { formatCurrency, formatOrderStatus, getOrderStatusColor, formatDateTime, canCancelOrder } from '@/lib/utils'
import OrderTimeline from '@/components/order/order-timeline'
import CancelOrderModal from '@/components/order/cancel-order-modal'

const { Title, Text } = Typography

export default function OrderDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [cancelModalOpen, setCancelModalOpen] = useState(false)

  const { data: order, isLoading } = useQuery({
    queryKey: ['order', id],
    queryFn: () => ordersApi.getById(id!),
    enabled: !!id,
  })

  const { data: tracking } = useQuery({
    queryKey: ['order-tracking', id],
    queryFn: () => ordersApi.getTracking(id!),
    enabled: !!id,
  })

  if (isLoading) {
    return (
      <div style={{ textAlign: 'center', padding: 48 }}>
        <Spin size="large" />
      </div>
    )
  }

  if (!order) {
    return (
      <div style={{ textAlign: 'center', padding: 48 }}>
        <Text>Không tìm thấy đơn hàng</Text>
      </div>
    )
  }

  const showCancelButton = canCancelOrder(order.order_status)

  return (
    <div style={{ padding: 24, maxWidth: 1200, margin: '0 auto' }}>
      {/* Breadcrumb */}
      <Breadcrumb
        items={[
          { title: <HomeOutlined />, href: '/' },
          { title: 'Đơn hàng', href: '/orders' },
          { title: order.order_code },
        ]}
        style={{ marginBottom: 24 }}
      />

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Title level={2}>Đơn hàng #{order.order_code}</Title>
        <Tag color={getOrderStatusColor(order.order_status)} style={{ fontSize: 14, padding: '4px 12px' }}>
          {formatOrderStatus(order.order_status)}
        </Tag>
      </div>

      {/* Payment Due Alert */}
      {order.order_status === 'PENDING_PAYMENT' && order.payment_due_at && (
        <Alert
          message="Đơn hàng chờ thanh toán"
          description={`Vui lòng thanh toán trước ${formatDateTime(order.payment_due_at)} để giữ đơn hàng.`}
          type="warning"
          showIcon
          style={{ marginBottom: 24 }}
        />
      )}

      <Row gutter={[24, 24]}>
        {/* Order Info */}
        <Col xs={24} lg={14}>
          <Space direction="vertical" size="large" style={{ width: '100%' }}>
            {/* Order Details */}
            <Card title="Thông tin đơn hàng" size="small">
              <Descriptions column={1} size="small">
                <Descriptions.Item label="Mã đơn hàng">{order.order_code}</Descriptions.Item>
                <Descriptions.Item label="Ngày đặt">{formatDateTime(order.placed_at)}</Descriptions.Item>
                <Descriptions.Item label="Trạng thái thanh toán">
                  <Tag color={order.payment_status === 'PAID' ? 'success' : 'warning'}>
                    {order.payment_status}
                  </Tag>
                </Descriptions.Item>
                {order.customer_note && (
                  <Descriptions.Item label="Ghi chú">{order.customer_note}</Descriptions.Item>
                )}
                {order.cancel_reason && (
                  <Descriptions.Item label="Lý do hủy">{order.cancel_reason}</Descriptions.Item>
                )}
              </Descriptions>
            </Card>

            {/* Shipping Address */}
            <Card title="Địa chỉ giao hàng" size="small">
              <Descriptions column={1} size="small">
                <Descriptions.Item label="Người nhận">{order.shipping_recipient_name}</Descriptions.Item>
                <Descriptions.Item label="Số điện thoại">{order.shipping_phone}</Descriptions.Item>
                <Descriptions.Item label="Địa chỉ">
                  {[
                    order.shipping_line1,
                    order.shipping_line2,
                    order.shipping_ward,
                    order.shipping_district,
                    order.shipping_city,
                    order.shipping_province,
                  ]
                    .filter(Boolean)
                    .join(', ')}
                </Descriptions.Item>
              </Descriptions>
            </Card>

            {/* Order Items */}
            <Card title="Sản phẩm" size="small">
              <List
                dataSource={order.items || []}
                renderItem={(item) => (
                  <List.Item>
                    <List.Item.Meta
                      title={item.product_name_snapshot}
                      description={`SKU: ${item.product_sku_snapshot} | SL: ${item.quantity}`}
                    />
                    <div className="price">{formatCurrency(item.line_total)}</div>
                  </List.Item>
                )}
              />

              <div style={{ marginTop: 16, paddingTop: 16, borderTop: '1px solid #f0f0f0' }}>
                <Space direction="vertical" style={{ width: '100%' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                    <Text>Tạm tính:</Text>
                    <Text>{formatCurrency(order.subtotal_amount)}</Text>
                  </div>
                  <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                    <Text>Phí vận chuyển:</Text>
                    <Text>{formatCurrency(order.shipping_fee)}</Text>
                  </div>
                  {parseFloat(order.discount_amount) > 0 && (
                    <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                      <Text>Giảm giá:</Text>
                      <Text type="danger">-{formatCurrency(order.discount_amount)}</Text>
                    </div>
                  )}
                  <div style={{ display: 'flex', justifyContent: 'space-between', paddingTop: 8, borderTop: '1px solid #f0f0f0' }}>
                    <Text strong>Tổng cộng:</Text>
                    <Text className="price" style={{ fontSize: 18, fontWeight: 700 }}>
                      {formatCurrency(order.total_amount)}
                    </Text>
                  </div>
                </Space>
              </div>
            </Card>

            {/* Actions */}
            {showCancelButton && (
              <div>
                <Button danger onClick={() => setCancelModalOpen(true)}>
                  Hủy đơn hàng
                </Button>
              </div>
            )}
          </Space>
        </Col>

        {/* Timeline */}
        <Col xs={24} lg={10}>
          <Card title="Lịch sử đơn hàng" size="small">
            {tracking && tracking.length > 0 ? (
              <OrderTimeline history={tracking} />
            ) : (
              <Text type="secondary">Chưa có lịch sử</Text>
            )}
          </Card>
        </Col>
      </Row>

      {/* Cancel Modal */}
      <CancelOrderModal
        orderId={order.id}
        open={cancelModalOpen}
        onClose={() => setCancelModalOpen(false)}
      />
    </div>
  )
}
