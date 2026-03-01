import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Typography,
  Card,
  List,
  Tag,
  Space,
  Button,
  Select,
  Empty,
  Spin,
  Pagination,
  DatePicker,
} from 'antd'
import { EyeOutlined } from '@ant-design/icons'
import { useQuery } from '@tanstack/react-query'
import { ordersApi } from '@/services/orders'
import type { OrderListParams } from '@/services/orders'
import { formatCurrency, formatOrderStatus, getOrderStatusColor, formatDateTime } from '@/lib/utils'

const { Title } = Typography
const { RangePicker } = DatePicker

export default function OrderListPage() {
  const navigate = useNavigate()
  const [filters, setFilters] = useState<OrderListParams>({
    page: 1,
    page_size: 10,
  })

  const { data, isLoading } = useQuery({
    queryKey: ['orders', filters],
    queryFn: () => ordersApi.list(filters),
  })

  const handleStatusChange = (status: string) => {
    setFilters((prev) => ({ ...prev, status: status || undefined, page: 1 }))
  }

  const handleDateChange = (dates: any) => {
    if (dates && dates[0] && dates[1]) {
      setFilters((prev) => ({
        ...prev,
        from: dates[0].format('YYYY-MM-DD'),
        to: dates[1].format('YYYY-MM-DD'),
        page: 1,
      }))
    } else {
      setFilters((prev) => {
        const { from, to, ...rest } = prev
        return rest
      })
    }
  }

  const handlePageChange = (page: number) => {
    setFilters((prev) => ({ ...prev, page }))
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  return (
    <div style={{ padding: 24, maxWidth: 1200, margin: '0 auto' }}>
      <Title level={2}>Đơn hàng của tôi</Title>

      {/* Filters */}
      <Card size="small" style={{ marginBottom: 16 }}>
        <Space wrap>
          <Select
            placeholder="Trạng thái"
            allowClear
            style={{ width: 200 }}
            onChange={handleStatusChange}
            options={[
              { label: 'Tất cả', value: '' },
              { label: 'Chờ thanh toán', value: 'PENDING_PAYMENT' },
              { label: 'Đã thanh toán', value: 'PAID' },
              { label: 'Đang xử lý', value: 'PROCESSING' },
              { label: 'Đang giao', value: 'SHIPPING' },
              { label: 'Hoàn tất', value: 'COMPLETED' },
              { label: 'Đã hủy', value: 'CANCELLED' },
            ]}
          />
          <RangePicker
            placeholder={['Từ ngày', 'Đến ngày']}
            onChange={handleDateChange}
            format="DD/MM/YYYY"
          />
        </Space>
      </Card>

      {isLoading ? (
        <div style={{ textAlign: 'center', padding: 48 }}>
          <Spin size="large" />
        </div>
      ) : !data || data.items.length === 0 ? (
        <Empty description="Chưa có đơn hàng nào" />
      ) : (
        <>
          <List
            dataSource={data.items}
            renderItem={(order) => (
              <Card key={order.id} style={{ marginBottom: 16 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 12 }}>
                  <Space>
                    <span style={{ fontWeight: 600 }}>#{order.order_code}</span>
                    <Tag color={getOrderStatusColor(order.order_status)}>
                      {formatOrderStatus(order.order_status)}
                    </Tag>
                  </Space>
                  <span style={{ color: '#999', fontSize: 12 }}>
                    {formatDateTime(order.placed_at)}
                  </span>
                </div>

                <div style={{ marginBottom: 12 }}>
                  <div style={{ color: '#666', fontSize: 12, marginBottom: 4 }}>
                    Người nhận: {order.shipping_recipient_name} - {order.shipping_phone}
                  </div>
                  <div style={{ color: '#666', fontSize: 12 }}>
                    Tổng tiền: <span className="price" style={{ fontWeight: 600 }}>{formatCurrency(order.total_amount)}</span>
                  </div>
                </div>

                <div style={{ textAlign: 'right' }}>
                  <Button
                    type="primary"
                    icon={<EyeOutlined />}
                    onClick={() => navigate(`/orders/${order.id}`)}
                  >
                    Xem chi tiết
                  </Button>
                </div>
              </Card>
            )}
          />

          {data.pagination.total_pages > 1 && (
            <div style={{ textAlign: 'center', marginTop: 24 }}>
              <Pagination
                current={data.pagination.page}
                pageSize={data.pagination.page_size}
                total={data.pagination.total}
                onChange={handlePageChange}
                showSizeChanger={false}
                showTotal={(total) => `Tổng ${total} đơn hàng`}
              />
            </div>
          )}
        </>
      )}
    </div>
  )
}
