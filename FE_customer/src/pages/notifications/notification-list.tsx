import { Typography, List, Card, Empty, Spin, Badge, Space, Button } from 'antd'
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

  const markAllAsReadMutation = useMutation({
    mutationFn: () => notificationsApi.markAllAsRead(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] })
    },
  })

  const handleNotificationClick = (id: string, isRead: boolean) => {
    if (!isRead) {
      markAsReadMutation.mutate(id)
    }
  }

  const unreadCount = data?.items.filter((notification) => !notification.read).length ?? 0

  if (isLoading) {
    return (
      <div style={{ textAlign: 'center', padding: 48 }}>
        <Spin size="large" />
      </div>
    )
  }

  return (
    <div style={{ padding: 24, maxWidth: 800, margin: '0 auto' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Title level={2} style={{ marginBottom: 0 }} data-testid="notifications-title">Thông báo</Title>
        <Button data-testid="notifications-mark-all" onClick={() => markAllAsReadMutation.mutate()} disabled={unreadCount === 0} loading={markAllAsReadMutation.isPending}>
          Đánh dấu tất cả đã đọc
        </Button>
      </div>

      {!data || data.items.length === 0 ? (
        <Empty description="Chưa có thông báo nào" />
      ) : (
        <List
          dataSource={data.items}
          renderItem={(notification) => (
            <Card
              key={notification.id}
              data-testid={`notification-item-${notification.id}`}
              style={{
                marginBottom: 12,
                background: notification.read ? '#fff' : '#f0f9ff',
                cursor: 'pointer',
              }}
              onClick={() => handleNotificationClick(notification.id, notification.read)}
            >
              <Space direction="vertical" style={{ width: '100%' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Text strong>{notification.title}</Text>
                  {!notification.read && <Badge status="processing" />}
                </div>
                <Text>{notification.content}</Text>
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
