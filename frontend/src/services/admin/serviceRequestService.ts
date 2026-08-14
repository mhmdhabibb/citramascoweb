import api from '../api'
import type { ApiResponse, ServiceRequest, ServiceRequestCategory, ServiceRequestPriority } from '@/types'

export interface CreateServiceRequestPayload {
  room_id: string
  reservation_id?: string
  guest_name: string
  guest_phone?: string
  category: ServiceRequestCategory
  title: string
  description: string
  priority?: ServiceRequestPriority
}

export interface AssignTaskPayload {
  assigned_to_user_id?: string
  notes?: string
  priority?: string
}

export interface CompleteTaskPayload {
  damage_charge?: number
}

export const serviceRequestService = {
  /**
   * POST /api/service-requests (Public / In-Room QR / Guest endpoint)
   */
  create: async (payload: CreateServiceRequestPayload): Promise<ServiceRequest> => {
    const response = await api.post<ApiResponse<ServiceRequest>>('/service-requests', payload)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Gagal mengirim laporan layanan kamar')
    }
    return response.data.data!
  },

  /**
   * GET /api/service-requests
   */
  getAll: async (status: string = 'all'): Promise<ServiceRequest[]> => {
    const response = await api.get<ApiResponse<ServiceRequest[]>>(`/service-requests?status=${status}`)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Gagal memuat daftar permintaan layanan')
    }
    return response.data.data || []
  },

  /**
   * GET /api/service-requests/:id
   */
  getById: async (id: string): Promise<ServiceRequest> => {
    const response = await api.get<ApiResponse<ServiceRequest>>(`/service-requests/${id}`)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Laporan tidak ditemukan')
    }
    return response.data.data!
  },

  /**
   * PATCH /api/service-requests/:id/assign (Reception assigns to Housekeeping)
   */
  assignToHousekeeping: async (id: string, payload: AssignTaskPayload): Promise<void> => {
    const response = await api.patch<ApiResponse<null>>(`/service-requests/${id}/assign`, payload)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Gagal menugaskan ke Housekeeping')
    }
  },

  /**
   * PATCH /api/service-requests/:id/complete (Housekeeping completes task)
   */
  complete: async (id: string, payload: CompleteTaskPayload = { damage_charge: 0 }): Promise<void> => {
    const response = await api.patch<ApiResponse<null>>(`/service-requests/${id}/complete`, payload)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Gagal menyelesaikan tugas')
    }
  },

  /**
   * PATCH /api/service-requests/:id/cancel
   */
  cancel: async (id: string): Promise<void> => {
    const response = await api.patch<ApiResponse<null>>(`/service-requests/${id}/cancel`)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Gagal membatalkan laporan')
    }
  },
}
