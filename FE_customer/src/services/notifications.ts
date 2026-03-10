import { api } from '@/lib/api'
import type { Notification, PaginatedResponse } from '@/types'

export interface NotificationListParams {
  page?: number
  page_size?: number
  status?: 'PENDING' | 'SENT' | 'FAILED'
}

export const notificationsApi = {
  list: async (params?: NotificationListParams) => {
    const { data } = await api.get<{ success: boolean; data: PaginatedResponse<Notification> }>('/notifications', { params })
    return data.data
  },

  markAsRead: async (id: string) => {
    const { data } = await api.patch<{ success: boolean; data: Notification }>(`/notifications/${id}/read`)
    return data.data
  },

  markAllAsRead: async () => {
    const { data } = await api.patch<{ success: boolean; data: { updated_count: number } }>('/notifications/read-all')
    return data.data.updated_count
  },
}
