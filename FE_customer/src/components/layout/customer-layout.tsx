import { Outlet } from 'react-router-dom'
import { Layout } from 'antd'
import Header from './header'
import Footer from './footer'

const { Content } = Layout

export default function CustomerLayout() {
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header />
      <Content style={{ marginTop: 64 }}>
        <Outlet />
      </Content>
      <Footer />
    </Layout>
  )
}
