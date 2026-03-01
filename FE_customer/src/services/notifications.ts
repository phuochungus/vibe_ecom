import { api } from '@/lib/api'
import type { Notification, PaginatedResponse } from '@/types'

export interface NotificationListParams {
  page?: number
  page_size?: number
  is_read?: boolean
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

  getUnreadCount: async () => {
    const { data } = await api.get<{ success: boolean; data: { count: number } }>('/notifications/unread-count')
    return data.data.count
  },
}
