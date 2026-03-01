import { Layout, Row, Col, Space } from 'antd'
import { MailOutlined, PhoneOutlined, EnvironmentOutlined } from '@ant-design/icons'

const { Footer: AntFooter } = Layout

export default function Footer() {
  return (
    <AntFooter style={{ background: '#f5f5f5', marginTop: 48 }}>
      <Row gutter={[32, 32]}>
        <Col xs={24} md={8}>
          <h3 style={{ fontSize: 16, fontWeight: 600, marginBottom: 16 }}>Golf Store</h3>
          <p style={{ color: '#666' }}>
            Cửa hàng chuyên cung cấp thiết bị golf chất lượng cao
          </p>
        </Col>
        <Col xs={24} md={8}>
          <h3 style={{ fontSize: 16, fontWeight: 600, marginBottom: 16 }}>Liên hệ</h3>
          <Space direction="vertical">
            <Space>
              <PhoneOutlined />
              <span>1900 xxxx</span>
            </Space>
            <Space>
              <MailOutlined />
              <span>support@golfstore.vn</span>
            </Space>
            <Space>
              <EnvironmentOutlined />
              <span>TP. Hồ Chí Minh</span>
            </Space>
          </Space>
        </Col>
        <Col xs={24} md={8}>
          <h3 style={{ fontSize: 16, fontWeight: 600, marginBottom: 16 }}>Chính sách</h3>
          <Space direction="vertical">
            <a href="#" style={{ color: '#666' }}>Chính sách đổi trả</a>
            <a href="#" style={{ color: '#666' }}>Chính sách bảo hành</a>
            <a href="#" style={{ color: '#666' }}>Điều khoản sử dụng</a>
          </Space>
        </Col>
      </Row>
      <div style={{ textAlign: 'center', marginTop: 32, paddingTop: 24, borderTop: '1px solid #e8e8e8' }}>
        © {new Date().getFullYear()} Golf Store. All rights reserved.
      </div>
    </AntFooter>
  )
}
