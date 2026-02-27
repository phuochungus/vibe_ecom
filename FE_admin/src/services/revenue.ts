import { api } from '@/lib/api'
import type { RevenueSummary, Order, PaginatedResponse } from '@/types'

export const revenueApi = {
  summary: (params: { from: string; to: string; group_by?: 'day' | 'month' | 'year' }) =>
    api
      .get<{ data: RevenueSummary }>('/admin/revenue/summary', { params })
      .then((r) => r.data.data),

  orders: (params: {
    from: string
    to: string
    page?: number
    page_size?: number
  }) =>
    api
      .get<{ data: PaginatedResponse<Order> }>('/admin/revenue/orders', { params })
      .then((r) => r.data.data),
}
