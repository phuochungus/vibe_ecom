import type { OrderStatus, ProductStatus } from "@/types"

export function formatCurrency(value: string | number): string {
  const num = typeof value === 'string' ? parseFloat(value) : value
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND',
    maximumFractionDigits: 0,
  }).format(num)
}

export function formatOrderStatus(status: OrderStatus): string {
  const map: Record<OrderStatus, string> = {
    NEW: 'Mới tạo',
    PENDING_PAYMENT: 'Chờ thanh toán',
    PAID: 'Đã thanh toán',
    PROCESSING: 'Đang xử lý',
    SHIPPING: 'Đang giao',
    COMPLETED: 'Hoàn tất',
    CANCELLED: 'Đã hủy',
  }
  return map[status] ?? status
}

export function getOrderStatusColor(status: OrderStatus): string {
  const map: Record<OrderStatus, string> = {
    NEW: 'default',
    PENDING_PAYMENT: 'warning',
    PAID: 'processing',
    PROCESSING: 'processing',
    SHIPPING: 'processing',
    COMPLETED: 'success',
    CANCELLED: 'error',
  }
  return map[status] ?? 'default'
}

export function formatProductStatus(status: ProductStatus): string {
  const map: Record<ProductStatus, string> = {
    ACTIVE: 'Đang bán',
    INACTIVE: 'Tạm dừng',
    DISCONTINUED: 'Ngừng bán',
  }
  return map[status] ?? status
}

export function getProductStatusColor(status: ProductStatus): string {
  const map: Record<ProductStatus, string> = {
    ACTIVE: 'success',
    INACTIVE: 'warning',
    DISCONTINUED: 'error',
  }
  return map[status] ?? 'default'
}

// Next valid order statuses per state machine (BR-ORD-03)
export const ORDER_TRANSITIONS: Record<OrderStatus, OrderStatus[]> = {
  NEW: ['PENDING_PAYMENT', 'CANCELLED'],
  PENDING_PAYMENT: ['PAID', 'CANCELLED'],
  PAID: ['PROCESSING'],
  PROCESSING: ['SHIPPING'],
  SHIPPING: ['COMPLETED'],
  COMPLETED: [],
  CANCELLED: [],
}

// Can cancel order (before SHIPPING)
export function canCancelOrder(status: OrderStatus): boolean {
  return ['NEW', 'PENDING_PAYMENT', 'PAID', 'PROCESSING'].includes(status)
}

// Format date time
export function formatDateTime(dateStr: string): string {
  const date = new Date(dateStr)
  return new Intl.DateTimeFormat('vi-VN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}

// Calculate payment time remaining (30 minutes)
export function getPaymentTimeRemaining(paymentDueAt?: string): {
  isExpired: boolean
  minutes: number
  seconds: number
} {
  if (!paymentDueAt) {
    return { isExpired: true, minutes: 0, seconds: 0 }
  }

  const now = new Date().getTime()
  const dueTime = new Date(paymentDueAt).getTime()
  const diff = dueTime - now

  if (diff <= 0) {
    return { isExpired: true, minutes: 0, seconds: 0 }
  }

  const minutes = Math.floor(diff / 60000)
  const seconds = Math.floor((diff % 60000) / 1000)

  return { isExpired: false, minutes, seconds }
}
