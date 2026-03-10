import { api } from '@/lib/api'
import type { Order, OrderStatus, PaginatedResponse } from '@/types'

export const ordersApi = {
  adminList: (params?: {
    status?: string
    from?: string
    to?: string
    page?: number
    page_size?: number
  }) =>
    api.get<{ data: PaginatedResponse<Order> }>('/admin/orders', { params }).then((r) => r.data.data),

  adminGet: (id: string) =>
    api.get<{ data: Order }>(`/admin/orders/${id}`).then((r) => r.data.data),

  updateStatus: (id: string, to_status: OrderStatus, reason?: string) =>
    api
      .patch<{ data: Order }>(`/admin/orders/${id}/status`, { to_status, reason })
      .then((r) => r.data.data),

  getTracking: (id: string) =>
    api.get<{ data: { order_id: string; current_status: OrderStatus; timeline: Order['status_history'] } }>(
      `/admin/orders/${id}/tracking`
    ).then((r) => r.data.data),
}
