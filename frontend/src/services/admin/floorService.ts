import api from '../api'
import type { ApiResponse, Floor } from '@/types'

export const floorService = {
  getAll: async (): Promise<Floor[]> => {
    const response = await api.get<ApiResponse<Floor[]>>('/floors')
    return response.data.data ?? []
  },

  getById: async (id: string): Promise<Floor> => {
    const response = await api.get<ApiResponse<Floor>>(`/floor/${id}`)
    if (!response.data.data) {
      throw new Error(response.data.message || 'Floor not found')
    }
    return response.data.data
  },

  create: async (data: { name: string; floor_number: number; description?: string }): Promise<string> => {
    const response = await api.post<ApiResponse<any>>('/floor', data)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to create floor')
    }
    return response.data.message
  },

  update: async (id: string, data: { name?: string; floor_number?: number; description?: string }): Promise<string> => {
    const response = await api.put<ApiResponse<any>>(`/floor/${id}`, data)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to update floor')
    }
    return response.data.message
  },

  delete: async (id: string): Promise<string> => {
    const response = await api.delete<ApiResponse<any>>(`/floor/${id}`)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to delete floor')
    }
    return response.data.message
  },
}
