import { useNavigate } from 'react-router-dom'
import { Typography, Button, Row, Col, Card } from 'antd'
import { ShoppingOutlined, TrophyOutlined, SafetyOutlined, RocketOutlined } from '@ant-design/icons'

const { Title, Paragraph } = Typography

export default function HomePage() {
  const navigate = useNavigate()

  return (
    <div>
      {/* Hero Section */}
      <div
        style={{
          background: 'linear-gradient(135deg, #EE4D2D 0%, #ff6b45 100%)',
          padding: '80px 24px',
          textAlign: 'center',
          color: '#fff',
        }}
      >
        <Title level={1} style={{ color: '#fff', fontSize: 48, marginBottom: 16 }}>
          Golf Store
        </Title>
        <Paragraph style={{ fontSize: 20, color: '#fff', marginBottom: 32, maxWidth: 600, margin: '0 auto 32px' }}>
          Chuyên cung cấp thiết bị golf chất lượng cao với giá tốt nhất thị trường
        </Paragraph>
        <Button
          type="primary"
          size="large"
          icon={<ShoppingOutlined />}
          onClick={() => navigate('/products')}
          style={{ background: '#fff', color: '#EE4D2D', border: 'none', height: 48, fontSize: 16, padding: '0 32px' }}
        >
          Mua sắm ngay
        </Button>
      </div>

      {/* Features */}
      <div style={{ padding: '64px 24px', maxWidth: 1200, margin: '0 auto' }}>
        <Title level={2} style={{ textAlign: 'center', marginBottom: 48 }}>
          Tại sao chọn chúng tôi?
        </Title>

        <Row gutter={[32, 32]}>
          <Col xs={24} md={8}>
            <Card hoverable style={{ textAlign: 'center', height: '100%' }}>
              <TrophyOutlined style={{ fontSize: 48, color: '#EE4D2D', marginBottom: 16 }} />
              <Title level={4}>Sản phẩm chính hãng</Title>
              <Paragraph type="secondary">
                100% sản phẩm chính hãng từ các thương hiệu golf hàng đầu thế giới
              </Paragraph>
            </Card>
          </Col>

          <Col xs={24} md={8}>
            <Card hoverable style={{ textAlign: 'center', height: '100%' }}>
              <SafetyOutlined style={{ fontSize: 48, color: '#EE4D2D', marginBottom: 16 }} />
              <Title level={4}>Bảo hành uy tín</Title>
              <Paragraph type="secondary">
                Chính sách bảo hành rõ ràng, đổi trả dễ dàng trong 30 ngày
              </Paragraph>
            </Card>
          </Col>

          <Col xs={24} md={8}>
            <Card hoverable style={{ textAlign: 'center', height: '100%' }}>
              <RocketOutlined style={{ fontSize: 48, color: '#EE4D2D', marginBottom: 16 }} />
              <Title level={4}>Giao hàng nhanh</Title>
              <Paragraph type="secondary">
                Giao hàng toàn quốc, nhanh chóng và an toàn
              </Paragraph>
            </Card>
          </Col>
        </Row>
      </div>

      {/* CTA Section */}
      <div
        style={{
          background: '#f5f5f5',
          padding: '64px 24px',
          textAlign: 'center',
        }}
      >
        <Title level={2} style={{ marginBottom: 24 }}>
          Sẵn sàng nâng tầm trải nghiệm golf của bạn?
        </Title>
        <Button
          type="primary"
          size="large"
          onClick={() => navigate('/products')}
          style={{ height: 48, fontSize: 16, padding: '0 32px' }}
        >
          Khám phá sản phẩm
        </Button>
      </div>
    </div>
  )
}
