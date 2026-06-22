import api from '../api'
import type { ApiResponse, Reservation } from '@/types'

export const reservationService = {
  /**
   * GET /api/reservations
   */
  getAll: async (): Promise<Reservation[]> => {
    const response = await api.get<ApiResponse<Reservation[]>>('/reservations')
    if (!response.data.data) {
      throw new Error(response.data.message || 'Failed to fetch reservations')
    }
    return response.data.data
  },

  /**
   * PATCH /api/reservation/approve/:id
   */
  approve: async (id: string): Promise<string> => {
    const response = await api.patch<ApiResponse<any>>(`/reservation/approve/${id}`)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to approve reservation')
    }
    return response.data.message
  },

  /**
   * PATCH /api/reservation/reject/:id
   */
  reject: async (id: string): Promise<string> => {
    const response = await api.patch<ApiResponse<any>>(`/reservation/reject/${id}`)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to reject reservation')
    }
    return response.data.message
  },

  /**
   * PATCH /api/reservation/cancel/:id
   */
  cancel: async (id: string): Promise<string> => {
    const response = await api.patch<ApiResponse<any>>(`/reservation/cancel/${id}`)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to cancel reservation')
    }
    return response.data.message
  },

  /**
   * PATCH /api/reservation/check-in/:id
   */
  checkIn: async (id: string): Promise<string> => {
    const response = await api.patch<ApiResponse<any>>(`/reservation/check-in/${id}`)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to check in guest')
    }
    return response.data.message
  },

  /**
   * PATCH /api/reservation/check-out/:id
   */
  checkOut: async (id: string): Promise<string> => {
    const response = await api.patch<ApiResponse<any>>(`/reservation/check-out/${id}`)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to check out guest')
    }
    return response.data.message
  },

  /**
   * POST /api/reservation
   */
  create: async (data: {
    room_id: string
    full_name: string
    email: string
    check_in_date: string
    check_out_date: string
    number_of_guest: number
    is_offer?: boolean
    offer_code?: string
  }): Promise<string> => {
    const response = await api.post<ApiResponse<any>>('/reservation', data)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to submit reservation')
    }
    return response.data.message
  },
}
