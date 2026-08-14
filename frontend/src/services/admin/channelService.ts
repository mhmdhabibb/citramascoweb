import api from '../api'
import type { ApiResponse, Channel } from '@/types'

export interface ChannelPaginatedResponse {
  data: Channel[]
  meta: {
    page: number
    total: number
    limit: number
  }
}

export const channelService = {
  /**
   * GET /api/channels — Fetch all channels without pagination
   */
  getAll: async (): Promise<Channel[]> => {
    const response = await api.get<ApiResponse<Channel[]>>('/channels')
    return response.data.data ?? []
  },

  /**
   * GET /api/channel?page=:page — Fetch channels with pagination
   */
  get: async (page = 1): Promise<{ channels: Channel[]; total: number; limit: number }> => {
    const response = await api.get<ApiResponse<Channel[]> & { meta?: { page: number; total: number; limit: number } }>(`/channel?page=${page}`)
    return {
      channels: response.data.data ?? [],
      total: response.data.meta?.total ?? (response.data.data?.length ?? 0),
      limit: response.data.meta?.limit ?? 10,
    }
  },

  /**
   * GET /api/channel/:id — Fetch channel by ID
   */
  getById: async (id: string): Promise<Channel> => {
    const response = await api.get<ApiResponse<Channel>>(`/channel/${id}`)
    if (!response.data.data) {
      throw new Error(response.data.message || 'Channel not found')
    }
    return response.data.data
  },

  /**
   * POST /api/channel — Create a new channel
   */
  create: async (data: { name: string }): Promise<string> => {
    const response = await api.post<ApiResponse<any>>('/channel', data)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to create channel')
    }
    return response.data.message
  },

  /**
   * PUT /api/channel/:id — Update a channel
   */
  update: async (id: string, data: { name: string }): Promise<string> => {
    const response = await api.put<ApiResponse<any>>(`/channel/${id}`, data)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to update channel')
    }
    return response.data.message
  },

  /**
   * DELETE /api/channel/:id — Delete a channel
   */
  delete: async (id: string): Promise<string> => {
    const response = await api.delete<ApiResponse<any>>(`/channel/${id}`)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to delete channel')
    }
    return response.data.message
  },
}
