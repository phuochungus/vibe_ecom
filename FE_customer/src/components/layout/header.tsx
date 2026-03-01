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
        background: '#fff',
        boxShadow: '0 1px 4px rgba(0,0,0,.08)',
        padding: '0 24px',
        display: 'flex',
        alignItems: 'center',
        gap: 24,
      }}
    >
      {/* Logo */}
      <Link to="/" style={{ display: 'flex', alignItems: 'center', textDecoration: 'none' }}>
        <div style={{ fontSize: 24, fontWeight: 700, color: '#EE4D2D' }}>
          Golf Store
        </div>
      </Link>

      {/* Search */}
      <Search
        placeholder="Tìm kiếm sản phẩm..."
        allowClear
        enterButton
        size="large"
        style={{ flex: 1, maxWidth: 600 }}
        value={searchQuery}
        onChange={(e) => setSearchQuery(e.target.value)}
        onSearch={handleSearch}
      />

      {/* Actions */}
      <Space size="large">
        {/* Cart */}
        <Badge count={itemCount} showZero={false} offset={[-2, 4]}>
          <Link to="/cart">
            <ShoppingCartOutlined style={{ fontSize: 24, color: '#333' }} />
          </Link>
        </Badge>

        {isAuthenticated ? (
          <>
            {/* Notifications */}
            <Badge dot offset={[-2, 4]}>
              <Link to="/notifications">
                <BellOutlined style={{ fontSize: 24, color: '#333' }} />
              </Link>
            </Badge>

            {/* User Menu */}
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <Space style={{ cursor: 'pointer' }}>
                <Avatar icon={<UserOutlined />} />
                <span>{user?.full_name}</span>
              </Space>
            </Dropdown>
          </>
        ) : (
          <Button type="primary" onClick={() => navigate('/login')}>
            Đăng nhập
          </Button>
        )}
      </Space>
    </AntHeader>
  )
}
