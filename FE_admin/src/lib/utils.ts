import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"
import type { OrderStatus, ProductStatus } from "@/types"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

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

export function getOrderStatusVariant(
  status: OrderStatus
): 'default' | 'secondary' | 'destructive' | 'outline' {
  if (status === 'COMPLETED') return 'default'
  if (status === 'CANCELLED') return 'destructive'
  if (status === 'PAID' || status === 'PROCESSING' || status === 'SHIPPING') return 'secondary'
  return 'outline'
}

export function formatProductStatus(status: ProductStatus): string {
  const map: Record<ProductStatus, string> = {
    ACTIVE: 'Đang bán',
    INACTIVE: 'Tạm dừng',
    DISCONTINUED: 'Ngừng bán',
  }
  return map[status] ?? status
}

export function getProductStatusVariant(
  status: ProductStatus
): 'default' | 'secondary' | 'destructive' | 'outline' {
  if (status === 'ACTIVE') return 'default'
  if (status === 'INACTIVE') return 'secondary'
  return 'destructive'
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
