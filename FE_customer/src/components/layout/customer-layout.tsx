import { Outlet } from 'react-router-dom'
import { Layout } from 'antd'
import Header from './header'
import Footer from './footer'

const { Content } = Layout

export default function CustomerLayout() {
  return (
    <Layout style={{ minHeight: '100vh', background: '#f4f6f8' }}>
      <Header />
      <Content style={{ marginTop: 72, background: '#f4f6f8' }}>
        <Outlet />
      </Content>
      <Footer />
    </Layout>
  )
}
