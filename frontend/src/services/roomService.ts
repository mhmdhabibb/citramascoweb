import api from './api'
import type { ApiResponse, Room } from '@/types'

export type { Room }

export const roomService = {
  getRooms: async (): Promise<Room[]> => {
    try {
      const response = await api.get<ApiResponse<Room[]>>('/rooms')
      return response.data.data ?? []
    } catch (error) {
      console.error('Error fetching rooms:', error)
      return []
    }
  },

  getAll: async (): Promise<Room[]> => {
    try {
      const response = await api.get<ApiResponse<Room[]>>('/rooms')
      return response.data.data ?? []
    } catch (error) {
      console.error('Error fetching rooms:', error)
      return []
    }
  },

  getRoomById: async (id: string): Promise<Room | undefined> => {
    try {
      const response = await api.get<ApiResponse<Room>>(`/room/${id}`)
      return response.data.data
    } catch (error) {
      console.error('Error fetching room details:', error)
      return undefined
    }
  },

  create: async (formData: FormData): Promise<string> => {
    const response = await api.post<ApiResponse<any>>('/room', formData)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Validation failed!')
    }
    return response.data.message
  },

  update: async (id: string, formData: FormData): Promise<string> => {
    const response = await api.patch<ApiResponse<any>>(`/room/${id}`, formData)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Validation failed!')
    }
    return response.data.message
  },

  delete: async (id: string): Promise<string> => {
    const response = await api.delete<ApiResponse<any>>(`/room/${id}`)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Deletion failed!')
    }
    return response.data.message
  },

  // Tambahkan method ini di dalam objek roomService
  updateStatus: async (id: string, status: string): Promise<string> => {
    const response = await api.patch<ApiResponse<any>>(`/room/status/${id}`, { status })
    if (!response.data.success) {
      throw new Error(response.data.message || 'Gagal memperbarui status operasional!')
    }
    return response.data.message
  },
}
