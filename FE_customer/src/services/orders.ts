import { api } from '@/lib/api'
import type { Order, OrderCreatePayload, OrderTracking, PaginatedResponse } from '@/types'

export interface OrderListParams {
  page?: number
  page_size?: number
  status?: string
  from?: string
  to?: string
}

export const ordersApi = {
  create: async (payload: OrderCreatePayload, idempotencyKey: string) => {
    const { data } = await api.post<{ success: boolean; data: Order }>(
      '/orders',
      payload,
      {
        headers: {
          'Idempotency-Key': idempotencyKey,
        },
      }
    )
    return data.data
  },

  list: async (params?: OrderListParams) => {
    const { data } = await api.get<{ success: boolean; data: PaginatedResponse<Order> }>('/orders', { params })
    return data.data
  },

  getById: async (id: string) => {
    const { data } = await api.get<{ success: boolean; data: Order }>(`/orders/${id}`)
    return data.data
  },

  getTracking: async (id: string) => {
    const { data } = await api.get<{ success: boolean; data: OrderTracking }>(`/orders/${id}/tracking`)
    return data.data
  },

  cancel: async (id: string, reason: string) => {
    const { data } = await api.post<{ success: boolean; data: Order }>(`/orders/${id}/cancel`, { reason })
    return data.data
  },
}
