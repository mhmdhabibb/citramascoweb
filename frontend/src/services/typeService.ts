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
  create: async (data: { name: string }): Promise<string> => {
    const response = await api.post<ApiResponse<any>>('/type', data)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Validation failed!')
    }
    return response.data.message
  },
  update: async (id: string, data: { name: string }): Promise<string> => {
    const response = await api.patch<ApiResponse<any>>(`/type/${id}`, data)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Validation failed!')
    }
    return response.data.message
  },
  delete: async (id: string): Promise<string> => {
    const response = await api.delete<ApiResponse<any>>(`/type/${id}`)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Deletion failed!')
    }
    return response.data.message
  }
}
