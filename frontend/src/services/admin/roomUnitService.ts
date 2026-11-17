import api from '../api'
import type { ApiResponse, RoomUnit } from '@/types'

export const roomUnitService = {
  getAll: async (params?: { room_id?: string; floor_id?: string }): Promise<RoomUnit[]> => {
    const response = await api.get<ApiResponse<RoomUnit[]>>('/room-units', { params })
    return response.data.data ?? []
  },

  getById: async (id: string): Promise<RoomUnit> => {
    const response = await api.get<ApiResponse<RoomUnit>>(`/room-unit/${id}`)
    if (!response.data.data) {
      throw new Error(response.data.message || 'Room unit not found')
    }
    return response.data.data
  },

  getNextAvailable: async (roomId: string): Promise<RoomUnit | null> => {
    try {
      const response = await api.get<ApiResponse<RoomUnit>>(`/room-unit/next-available/${roomId}`)
      return response.data.data ?? null
    } catch {
      return null
    }
  },

  create: async (data: { room_number: string; room_id: string; floor_id: string; status?: string }): Promise<string> => {
    const response = await api.post<ApiResponse<any>>('/room-unit', data)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to create room unit')
    }
    return response.data.message
  },

  update: async (id: string, data: { room_number?: string; room_id?: string; floor_id?: string; status?: string }): Promise<string> => {
    const response = await api.put<ApiResponse<any>>(`/room-unit/${id}`, data)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to update room unit')
    }
    return response.data.message
  },

  // Update status and push notification to housekeeping
  setStatusAndNotify: async (id: string, status: string): Promise<string> => {
    const msg = await roomUnitService.updateStatus(id, status)
    // Notify housekeeping role
    await api.post<ApiResponse<null>>('/notifications/', {
      title: `Room ${id} status changed`,
      body: `Status set to ${status}. Please attend.`,
      target_roles: ['housekeeping'],
      room_unit_id: id,
    })
    return msg
  },

  delete: async (id: string): Promise<string> => {
    const response = await api.delete<ApiResponse<any>>(`/room-unit/${id}`)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to delete room unit')
    }
    return response.data.message
  },
}
