import { Typography, List, Card, Empty, Spin, Badge, Space } from 'antd'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { notificationsApi } from '@/services/notifications'
import { formatDateTime } from '@/lib/utils'

const { Title, Text } = Typography

export default function NotificationListPage() {
  const queryClient = useQueryClient()

  const { data, isLoading } = useQuery({
    queryKey: ['notifications'],
    queryFn: () => notificationsApi.list({ page_size: 50 }),
  })

  const markAsReadMutation = useMutation({
    mutationFn: (id: string) => notificationsApi.markAsRead(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] })
    },
  })

  const handleNotificationClick = (id: string, isRead: boolean) => {
    if (!isRead) {
      markAsReadMutation.mutate(id)
    }
  }

  if (isLoading) {
    return (
      <div style={{ textAlign: 'center', padding: 48 }}>
        <Spin size="large" />
      </div>
    )
  }

  return (
    <div style={{ padding: 24, maxWidth: 800, margin: '0 auto' }}>
      <Title level={2}>Thông báo</Title>

      {!data || data.items.length === 0 ? (
        <Empty description="Chưa có thông báo nào" />
      ) : (
        <List
          dataSource={data.items}
          renderItem={(notification) => (
            <Card
              key={notification.id}
              style={{
                marginBottom: 12,
                background: notification.is_read ? '#fff' : '#f0f9ff',
                cursor: 'pointer',
              }}
              onClick={() => handleNotificationClick(notification.id, notification.is_read)}
            >
              <Space direction="vertical" style={{ width: '100%' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Text strong>{notification.title}</Text>
                  {!notification.is_read && <Badge status="processing" />}
                </div>
                <Text>{notification.message}</Text>
                <Text type="secondary" style={{ fontSize: 12 }}>
                  {formatDateTime(notification.created_at)}
                </Text>
              </Space>
            </Card>
          )}
        />
      )}
    </div>
  )
}
