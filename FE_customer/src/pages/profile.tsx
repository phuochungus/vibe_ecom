import { Typography, Card, Descriptions, Tag } from 'antd'
import { useAuth } from '@/lib/auth'

const { Title } = Typography

export default function ProfilePage() {
  const { user } = useAuth()

  if (!user) return null

  return (
    <div style={{ padding: 24, maxWidth: 800, margin: '0 auto' }}>
      <Title level={2}>Tài khoản của tôi</Title>

      <Card>
        <Descriptions column={1} bordered>
          <Descriptions.Item label="Họ và tên">{user.full_name}</Descriptions.Item>
          <Descriptions.Item label="Email">{user.email}</Descriptions.Item>
          <Descriptions.Item label="Số điện thoại">{user.phone}</Descriptions.Item>
          <Descriptions.Item label="Vai trò">
            <Tag color="blue">{user.role}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="Trạng thái">
            <Tag color={user.status === 'ACTIVE' ? 'success' : 'error'}>
              {user.status}
            </Tag>
          </Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  )
}
