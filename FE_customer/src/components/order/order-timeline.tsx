import { Timeline, Typography } from 'antd'
import { CheckCircleOutlined, ClockCircleOutlined } from '@ant-design/icons'
import type { OrderTrackingEvent } from '@/types'
import { formatOrderStatus } from '@/lib/utils'
import { formatDateTime } from '@/lib/utils'

const { Text } = Typography

interface OrderTimelineProps {
  history: OrderTrackingEvent[]
}

export default function OrderTimeline({ history }: OrderTimelineProps) {
  const sortedHistory = [...history].sort(
    (a, b) => new Date(b.occurred_at).getTime() - new Date(a.occurred_at).getTime()
  )

  return (
    <Timeline
      items={sortedHistory.map((event, index) => ({
        dot: index === 0 ? <CheckCircleOutlined style={{ fontSize: 16 }} /> : <ClockCircleOutlined />,
        color: index === 0 ? 'green' : 'gray',
        children: (
          <div>
            <div style={{ fontWeight: 500 }}>
              {formatOrderStatus(event.status)}
            </div>
            <Text type="secondary" style={{ fontSize: 12 }}>
              {formatDateTime(event.occurred_at)}
            </Text>
            {event.description && (
              <div style={{ marginTop: 4 }}>
                <Text type="secondary" style={{ fontSize: 12 }}>
                  {event.description}
                </Text>
              </div>
            )}
            <div style={{ marginTop: 4 }}>
              <Text type="secondary" style={{ fontSize: 11 }}>
                Bởi: {event.source_type}
              </Text>
            </div>
          </div>
        ),
      }))}
    />
  )
}
