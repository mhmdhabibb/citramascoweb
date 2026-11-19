import type { ApiResponse, Reservation, ReservationItem, ReservationLog } from '@/types'
import api from '../api'

export const reservationService = {
  /**
   * GET /api/reservations
   */
  getAll: async (): Promise<Reservation[]> => {
    const response = await api.get<ApiResponse<Reservation[]>>('/reservations')
    return response.data.data ?? []
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
    number_of_adult: number
    number_of_children: number
    is_offer?: boolean
    offer_code?: string
    deposit?: number
    channel_id?: string
    is_early_checkin?: boolean
    estimated_arrival_time?: string
  }): Promise<string> => {
    const response = await api.post<ApiResponse<any>>('/reservation', data)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to submit reservation')
    }
    return response.data.message
  },

  /**
   * GET /api/reservation/:id/logs
   */
  getLogs: async (reservationId: string): Promise<ReservationLog[]> => {
    const response = await api.get<ApiResponse<ReservationLog[]>>(`/reservation/${reservationId}/logs`)
    return response.data.data ?? []
  },

  /**
   * POST /api/reservation/:id/logs
   */
  createLog: async (
    reservationId: string,
    data: { action: string; is_early?: boolean; notes?: string; handled_by_id?: string },
  ): Promise<ReservationLog> => {
    const response = await api.post<ApiResponse<ReservationLog>>(`/reservation/${reservationId}/logs`, data)
    return response.data.data!
  },

  /**
   * GET /api/reservation/:id/items
   */
  getItems: async (reservationId: string): Promise<ReservationItem[]> => {
    const response = await api.get<ApiResponse<ReservationItem[]>>(`/reservation/${reservationId}/items`)
    return response.data.data ?? []
  },

  /**
   * POST /api/reservation/:id/items
   */
  createItem: async (
    reservationId: string,
    data: { room_unit_id: string; stay_date: string; rate_applied: number },
  ): Promise<ReservationItem> => {
    const response = await api.post<ApiResponse<ReservationItem>>(`/reservation/${reservationId}/items`, data)
    return response.data.data!
  },
}
