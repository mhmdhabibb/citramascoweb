import api from '../api'
import type { ApiResponse, AppNotification } from '@/types'

export const notificationService = {
  /**
   * GET /api/notifications
   */
  getAll: async (): Promise<AppNotification[]> => {
    const response = await api.get<ApiResponse<AppNotification[]>>('/notifications/')
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to fetch notifications')
    }
    return response.data.data || []
  },

  /**
   * GET /api/notifications/unread-count
   */
  getUnreadCount: async (): Promise<number> => {
    const response = await api.get<ApiResponse<number>>('/notifications/unread-count')
    if (!response.data.success) {
      return 0
    }
    return response.data.data || 0
  },

  /**
   * PATCH /api/notifications/:id/read
   */
  markRead: async (id: string): Promise<void> => {
    await api.patch<ApiResponse<null>>(`/notifications/${id}/read`)
  },
}