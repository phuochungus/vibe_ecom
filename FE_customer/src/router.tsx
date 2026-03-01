import { createBrowserRouter, Navigate } from 'react-router-dom'
import { useAuth } from './lib/auth'
import CustomerLayout from './components/layout/customer-layout'
import LoginPage from './pages/login'
import HomePage from './pages/home'
import ProductListPage from './pages/products/product-list'
import ProductDetailPage from './pages/products/product-detail'
import CartPage from './pages/cart/cart-page'
import CheckoutPage from './pages/checkout/checkout-page'
import PaymentReturnPage from './pages/checkout/payment-return'
import PaymentCancelPage from './pages/checkout/payment-cancel'
import OrderListPage from './pages/orders/order-list'
import OrderDetailPage from './pages/orders/order-detail'
import NotificationListPage from './pages/notifications/notification-list'
import ProfilePage from './pages/profile'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth()
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />
}

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />,
  },
  {
    path: '/',
    element: <CustomerLayout />,
    children: [
      {
        index: true,
        element: <HomePage />,
      },
      {
        path: 'products',
        element: <ProductListPage />,
      },
      {
        path: 'products/:id',
        element: <ProductDetailPage />,
      },
      {
        path: 'cart',
        element: <CartPage />,
      },
      {
        path: 'checkout',
        element: (
          <ProtectedRoute>
            <CheckoutPage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'checkout/return',
        element: (
          <ProtectedRoute>
            <PaymentReturnPage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'checkout/cancel',
        element: (
          <ProtectedRoute>
            <PaymentCancelPage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'orders',
        element: (
          <ProtectedRoute>
            <OrderListPage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'orders/:id',
        element: (
          <ProtectedRoute>
            <OrderDetailPage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'notifications',
        element: (
          <ProtectedRoute>
            <NotificationListPage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'profile',
        element: (
          <ProtectedRoute>
            <ProfilePage />
          </ProtectedRoute>
        ),
      },
    ],
  },
])
