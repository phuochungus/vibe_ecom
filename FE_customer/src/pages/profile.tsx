import { Avatar, Card, Col, Divider, List, Row, Space, Statistic, Tag, Typography } from 'antd'
import {
  CheckCircleOutlined,
  MailOutlined,
  PhoneOutlined,
  SafetyCertificateOutlined,
  ShoppingOutlined,
  UserOutlined,
} from '@ant-design/icons'
import { useAuth } from '@/lib/auth'

const { Paragraph, Text, Title } = Typography

const quickActions = [
  'Theo dõi đơn hàng và trạng thái vận chuyển',
  'Nhận thông báo cho thay đổi đơn hàng và thanh toán',
  'Quản lý thông tin liên hệ cho các lần mua tiếp theo',
]

export default function ProfilePage() {
  const { user } = useAuth()

  if (!user) return null

  return (
    <div
      style={{
        padding: '32px 24px 56px',
        maxWidth: 1320,
        margin: '0 auto',
      }}
    >
      <Card
        bordered={false}
        style={{
          marginBottom: 24,
          borderRadius: 24,
          overflow: 'hidden',
          background: 'linear-gradient(135deg, #fff7f2 0%, #ffffff 52%, #eef6ff 100%)',
          boxShadow: '0 18px 48px rgba(15, 23, 42, 0.08)',
        }}
        bodyStyle={{ padding: 0 }}
      >
        <Row gutter={0}>
          <Col xs={24} xl={15}>
            <div style={{ padding: 32 }}>
              <Space align="start" size={20}>
                <Avatar
                  size={88}
                  icon={<UserOutlined />}
                  style={{
                    background: 'linear-gradient(135deg, #ee4d2d 0%, #ff8f5a 100%)',
                    boxShadow: '0 12px 28px rgba(238, 77, 45, 0.25)',
                  }}
                />
                <div>
                  <Tag color="blue" style={{ borderRadius: 999, paddingInline: 12, marginBottom: 12 }}>
                    Hồ sơ khách hàng
                  </Tag>
                  <Title level={2} style={{ margin: 0 }}>
                    {user.full_name}
                  </Title>
                  <Paragraph style={{ maxWidth: 620, marginTop: 12, marginBottom: 0, color: '#52606d', fontSize: 16 }}>
                    Đây là trung tâm quản lý tài khoản của bạn. Trên màn hình laptop, bố cục này mở rộng để thông tin
                    chính, trạng thái tài khoản và các tác vụ nhanh đều nằm trong cùng một vùng quan sát.
                  </Paragraph>
                </div>
              </Space>
            </div>
          </Col>
          <Col xs={24} xl={9}>
            <div
              style={{
                height: '100%',
                minHeight: 220,
                padding: 32,
                background: 'rgba(255,255,255,0.78)',
                borderLeft: '1px solid rgba(15, 23, 42, 0.06)',
              }}
            >
              <Row gutter={[16, 16]}>
                <Col xs={12}>
                  <Card bordered={false} style={{ borderRadius: 20, background: '#fff' }}>
                    <Statistic title="Vai trò" value={user.role} />
                  </Card>
                </Col>
                <Col xs={12}>
                  <Card bordered={false} style={{ borderRadius: 20, background: '#fff' }}>
                    <Statistic title="Trạng thái" value={user.status} valueStyle={{ color: '#16a34a' }} />
                  </Card>
                </Col>
                <Col span={24}>
                  <Card bordered={false} style={{ borderRadius: 20, background: '#fff' }}>
                    <Space align="start" size={12}>
                      <CheckCircleOutlined style={{ color: '#16a34a', fontSize: 20, marginTop: 2 }} />
                      <div>
                        <Text strong>Tài khoản sẵn sàng cho Sprint 1</Text>
                        <Paragraph style={{ margin: '6px 0 0', color: '#5b6572' }}>
                          Bạn có thể đăng nhập, theo dõi đơn hàng, nhận thông báo và tiếp tục checkout từ cùng một khu vực.
                        </Paragraph>
                      </div>
                    </Space>
                  </Card>
                </Col>
              </Row>
            </div>
          </Col>
        </Row>
      </Card>

      <Row gutter={[24, 24]} align="stretch">
        <Col xs={24} lg={14}>
          <Card
            title="Thông tin tài khoản"
            bordered={false}
            style={{ height: '100%', borderRadius: 24, boxShadow: '0 16px 40px rgba(15, 23, 42, 0.06)' }}
          >
            <Row gutter={[16, 16]}>
              <Col xs={24} md={12}>
                <Card bordered={false} style={{ background: '#f7f9fc', borderRadius: 18 }}>
                  <Space align="start">
                    <MailOutlined style={{ color: '#ee4d2d', fontSize: 18, marginTop: 2 }} />
                    <div>
                      <Text type="secondary">Email</Text>
                      <div style={{ fontWeight: 600, marginTop: 4 }}>{user.email}</div>
                    </div>
                  </Space>
                </Card>
              </Col>
              <Col xs={24} md={12}>
                <Card bordered={false} style={{ background: '#f7f9fc', borderRadius: 18 }}>
                  <Space align="start">
                    <PhoneOutlined style={{ color: '#ee4d2d', fontSize: 18, marginTop: 2 }} />
                    <div>
                      <Text type="secondary">Số điện thoại</Text>
                      <div style={{ fontWeight: 600, marginTop: 4 }}>{user.phone}</div>
                    </div>
                  </Space>
                </Card>
              </Col>
              <Col xs={24} md={12}>
                <Card bordered={false} style={{ background: '#f7f9fc', borderRadius: 18 }}>
                  <Space align="start">
                    <UserOutlined style={{ color: '#ee4d2d', fontSize: 18, marginTop: 2 }} />
                    <div>
                      <Text type="secondary">Họ và tên</Text>
                      <div style={{ fontWeight: 600, marginTop: 4 }}>{user.full_name}</div>
                    </div>
                  </Space>
                </Card>
              </Col>
              <Col xs={24} md={12}>
                <Card bordered={false} style={{ background: '#f7f9fc', borderRadius: 18 }}>
                  <Space align="start">
                    <SafetyCertificateOutlined style={{ color: '#ee4d2d', fontSize: 18, marginTop: 2 }} />
                    <div>
                      <Text type="secondary">Truy cập</Text>
                      <div style={{ fontWeight: 600, marginTop: 4 }}>{user.role === 'USER' ? 'Khách mua hàng' : user.role}</div>
                    </div>
                  </Space>
                </Card>
              </Col>
            </Row>
          </Card>
        </Col>

        <Col xs={24} lg={10}>
          <Card
            title="Tác vụ nhanh"
            bordered={false}
            style={{ borderRadius: 24, boxShadow: '0 16px 40px rgba(15, 23, 42, 0.06)' }}
          >
            <List
              dataSource={quickActions}
              split={false}
              renderItem={(item) => (
                <List.Item style={{ paddingInline: 0 }}>
                  <Card bordered={false} style={{ width: '100%', background: '#fff7f2', borderRadius: 18 }}>
                    <Space align="start">
                      <ShoppingOutlined style={{ color: '#ee4d2d', fontSize: 18, marginTop: 2 }} />
                      <Text>{item}</Text>
                    </Space>
                  </Card>
                </List.Item>
              )}
            />
            <Divider />
            <Paragraph style={{ marginBottom: 0, color: '#5b6572' }}>
              Màn hình này đã được mở rộng để sử dụng tốt hơn không gian desktop: khối thông tin chính nằm bên trái,
              thẻ trạng thái và tóm tắt nằm bên phải, tránh cảm giác trống trải trên laptop.
            </Paragraph>
          </Card>
        </Col>
      </Row>
    </div>
  )
}
