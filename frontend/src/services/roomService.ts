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

  create: async (formData: FormData): Promise<void> => {
    await api.post('/room', formData)
  },

  update: async (id: string, formData: FormData): Promise<void> => {
    await api.patch(`/room/${id}`, formData)
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/room/${id}`)
  }
}