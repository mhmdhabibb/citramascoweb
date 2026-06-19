import api from './api'
import type { ApiResponse, RoomType } from '@/types'

export const typeService = {
  /**
   * GET /api/types — Fetch all room types
   */
  getAll: async (): Promise<RoomType[]> => {
    const response = await api.get<ApiResponse<RoomType[]>>('/types')
    return response.data.data ?? []
  },

  /**
   * GET /api/type/:id — Fetch room type by ID
   */
  getById: async (id: string): Promise<RoomType> => {
    const response = await api.get<ApiResponse<RoomType>>(`/type/${id}`)
    if (!response.data.data) {
      throw new Error(response.data.message || 'Type not found')
    }
    return response.data.data
  },
  create: async (data: { name: string }): Promise<void> => {
    await api.post('/type', data)
  },
  update: async (id: string, data: { name: string }): Promise<void> => {
    await api.patch(`/type/${id}`, data)
  },
  delete: async (id: string): Promise<void> => {
    await api.delete(`/type/${id}`)
  }
}
