import { Link, useNavigate } from 'react-router-dom'
import { Layout, Input, Badge, Dropdown, Avatar, Button, Space } from 'antd'
import {
  ShoppingCartOutlined,
  BellOutlined,
  UserOutlined,
  LogoutOutlined,
  UnorderedListOutlined,
} from '@ant-design/icons'
import type { MenuProps } from 'antd'
import { useAuth } from '@/lib/auth'
import { useCart } from '@/lib/cart'
import { useState } from 'react'

const { Header: AntHeader } = Layout
const { Search } = Input

export default function Header() {
  const navigate = useNavigate()
  const { user, isAuthenticated, logout } = useAuth()
  const { itemCount } = useCart()
  const [searchQuery, setSearchQuery] = useState('')

  const handleSearch = (value: string) => {
    if (value.trim()) {
      navigate(`/products?search=${encodeURIComponent(value.trim())}`)
    }
  }

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  const userMenuItems: MenuProps['items'] = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: 'Tài khoản',
      onClick: () => navigate('/profile'),
    },
    {
      key: 'orders',
      icon: <UnorderedListOutlined />,
      label: 'Đơn hàng',
      onClick: () => navigate('/orders'),
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: 'Đăng xuất',
      onClick: handleLogout,
    },
  ]

  return (
    <AntHeader
      style={{
        position: 'fixed',
        top: 0,
        zIndex: 1000,
        width: '100%',
        height: 72,
        background: 'rgba(255,255,255,0.96)',
        backdropFilter: 'blur(10px)',
        boxShadow: '0 1px 0 rgba(15, 23, 42, 0.08)',
        padding: 0,
      }}
    >
      <div
        style={{
          width: '100%',
          maxWidth: 1440,
          margin: '0 auto',
          padding: '0 24px',
          height: '100%',
          display: 'flex',
          alignItems: 'center',
          gap: 24,
        }}
      >
        <Link
          to="/"
          style={{
            display: 'flex',
            alignItems: 'center',
            textDecoration: 'none',
            whiteSpace: 'nowrap',
            flexShrink: 0,
          }}
        >
          <div style={{ fontSize: 24, fontWeight: 700, color: '#EE4D2D' }}>Golf Store</div>
        </Link>

        <div style={{ flex: 1, minWidth: 220, maxWidth: 760 }}>
          <Search
            placeholder="Tìm kiếm sản phẩm..."
            allowClear
            enterButton
            size="large"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            onSearch={handleSearch}
          />
        </div>

        <div
          style={{
            marginLeft: 'auto',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'flex-end',
            gap: 20,
            flexShrink: 0,
          }}
        >
          <Badge count={itemCount} showZero={false} offset={[-2, 4]}>
            <Link to="/cart" data-testid="header-cart-link">
              <ShoppingCartOutlined style={{ fontSize: 24, color: '#333' }} />
            </Link>
          </Badge>

          {isAuthenticated ? (
            <>
              <Badge dot offset={[-2, 4]}>
                <Link to="/notifications" data-testid="header-notifications-link">
                  <BellOutlined style={{ fontSize: 24, color: '#333' }} />
                </Link>
              </Badge>

              <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
                <Space style={{ cursor: 'pointer', whiteSpace: 'nowrap' }}>
                  <Avatar icon={<UserOutlined />} />
                  <span>{user?.full_name}</span>
                </Space>
              </Dropdown>
            </>
          ) : (
            <Button type="primary" onClick={() => navigate('/login')} data-testid="header-login-button">
              Đăng nhập
            </Button>
          )}
        </div>
      </div>
    </AntHeader>
  )
}
