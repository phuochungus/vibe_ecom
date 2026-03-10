// --- Enums ---
export type UserRole = 'USER' | 'ADMIN'
export type UserStatus = 'ACTIVE' | 'LOCKED' | 'DISABLED'

export type ProductStatus = 'ACTIVE' | 'INACTIVE' | 'DISCONTINUED'

export type OrderStatus =
  | 'NEW'
  | 'PENDING_PAYMENT'
  | 'PAID'
  | 'PROCESSING'
  | 'SHIPPING'
  | 'COMPLETED'
  | 'CANCELLED'

export type PaymentStatus = 'UNPAID' | 'PAID' | 'FAILED' | 'REFUNDED'
export type PaymentTxnState = 'PENDING' | 'SUCCESS' | 'FAILED' | 'CANCELLED'
export type NotificationStatus = 'PENDING' | 'SENT' | 'FAILED'

// --- User ---
export interface User {
  id: string
  email: string
  phone: string
  full_name: string
  role: UserRole
  status: UserStatus
}

// --- Product ---
export interface Product {
  id: string
  sku: string
  name: string
  description?: string
  price: string
  stock: number
  status: ProductStatus
  image_url?: string
  created_at: string
  updated_at: string
}

export interface ProductCreatePayload {
  sku: string
  name: string
  description?: string
  price: string
  stock: number
  status: ProductStatus
  image_url?: string
}

export type ProductUpdatePayload = Partial<ProductCreatePayload>

// --- Order ---
export interface OrderItem {
  id: string
  product_id: string
  product_sku_snapshot?: string
  product_name_snapshot?: string
  sku?: string
  name?: string
  unit_price: string
  quantity: number
  line_total: string
}

export interface OrderStatusHistory {
  id: string
  from_status: OrderStatus | null
  to_status: OrderStatus
  changed_by_type: 'SYSTEM' | 'USER' | 'ADMIN' | 'PAYMENT_GATEWAY' | 'CARRIER'
  change_reason?: string
  occurred_at: string
  status?: OrderStatus
  source_type?: 'SYSTEM' | 'USER' | 'ADMIN' | 'PAYMENT_GATEWAY' | 'CARRIER'
  description?: string
}

export interface Order {
  id: string
  order_code: string
  user_id: string
  user?: { full_name: string; email: string; phone: string }
  order_status: OrderStatus
  payment_status: PaymentStatus
  currency_code: string
  subtotal_amount: string
  discount_amount: string
  shipping_fee: string
  total_amount: string
  payment_due_at?: string
  customer_note?: string
  cancel_reason?: string
  shipping_recipient_name: string
  shipping_phone: string
  shipping_line1: string
  shipping_line2?: string
  shipping_ward?: string
  shipping_district?: string
  shipping_city: string
  shipping_province?: string
  shipping_country_code: string
  placed_at: string
  created_at: string
  items?: OrderItem[]
  status_history?: OrderStatusHistory[]
}

// --- Payment ---
export interface PaymentTransaction {
  id: string
  order_id: string
  txn_type: 'PAYMENT' | 'REFUND'
  provider_name: string
  provider_txn_code: string
  amount: string
  currency_code: string
  status: PaymentTxnState
  provider_message?: string
  processed_at?: string
  created_at: string
}

// --- Revenue ---
export interface RevenueSummary {
  from: string
  to: string
  gross_revenue: string
  refund_amount: string
  net_revenue: string
  completed_orders: number
  payment_success_rate: string
}

export interface RevenueDataPoint {
  date: string
  gross_revenue: string
  net_revenue: string
  completed_orders: number
}

// --- API Response Envelope ---
export interface ApiResponse<T> {
  success: boolean
  data: T
  meta?: Record<string, unknown>
  request_id?: string
  timestamp: string
}

export interface ApiError {
  success: false
  error: {
    code: string
    message: string
    details?: unknown[]
  }
  request_id?: string
  timestamp: string
}

export interface PaginationMeta {
  page: number
  page_size: number
  total: number
  total_pages: number
}

export interface PaginatedResponse<T> {
  items: T[]
  pagination: PaginationMeta
}

// --- Auth ---
export interface LoginResponse {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
  user: User
}
