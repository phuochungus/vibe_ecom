import { createBrowserRouter, Navigate, Outlet } from 'react-router-dom'
import { useAuth } from '@/lib/auth'
import { AdminLayout } from '@/components/layout/admin-layout'
import LoginPage from '@/pages/login'
import DashboardPage from '@/pages/dashboard'
import ProductListPage from '@/pages/products/product-list'
import ProductFormPage from '@/pages/products/product-form'
import OrderListPage from '@/pages/orders/order-list'
import OrderDetailPage from '@/pages/orders/order-detail'
import RevenuePage from '@/pages/revenue/revenue-page'

function ProtectedRoute() {
  const { isAuthenticated } = useAuth()
  if (!isAuthenticated) return <Navigate to="/login" replace />
  return (
    <AdminLayout>
      <Outlet />
    </AdminLayout>
  )
}

function GuestRoute() {
  const { isAuthenticated } = useAuth()
  if (isAuthenticated) return <Navigate to="/" replace />
  return <Outlet />
}

export const router = createBrowserRouter([
  {
    element: <GuestRoute />,
    children: [{ path: '/login', element: <LoginPage /> }],
  },
  {
    element: <ProtectedRoute />,
    children: [
      { path: '/', element: <DashboardPage /> },
      { path: '/products', element: <ProductListPage /> },
      { path: '/products/new', element: <ProductFormPage /> },
      { path: '/products/:id/edit', element: <ProductFormPage /> },
      { path: '/orders', element: <OrderListPage /> },
      { path: '/orders/:id', element: <OrderDetailPage /> },
      { path: '/revenue', element: <RevenuePage /> },
    ],
  },
  { path: '*', element: <Navigate to="/" replace /> },
])
