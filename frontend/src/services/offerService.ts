import api from './api'
import type { ApiResponse, Offer } from '@/types'

export const offerService = {
  /**
   * GET /api/offers - Fetch active public offers
   */
  getAll: async (): Promise<Offer[]> => {
    const response = await api.get<ApiResponse<Offer[]>>('/offers')
    return response.data.data ?? []
  },

  /**
   * GET /api/offer - Fetch all offers for admin/manager
   */
  getAllAdmin: async (): Promise<Offer[]> => {
    const response = await api.get<ApiResponse<Offer[]>>('/offer')
    return response.data.data ?? []
  },

  /**
   * GET /api/offer/:id - Fetch offer by ID
   */
  getById: async (id: string): Promise<Offer> => {
    const response = await api.get<ApiResponse<Offer>>(`/offer/${id}`)
    if (!response.data.data) {
      throw new Error(response.data.message || 'Offer not found')
    }
    return response.data.data
  },

  /**
   * POST /api/offer/ - Create a new offer (multipart form-data)
   */
  create: async (formData: FormData): Promise<string> => {
    const response = await api.post<ApiResponse<any>>('/offer/', formData)
    return response.data.message
  },

  /**
   * PATCH /api/offer/:id - Update an offer (multipart form-data)
   */
  update: async (id: string, formData: FormData): Promise<string> => {
    const response = await api.patch<ApiResponse<any>>(`/offer/${id}`, formData)
    return response.data.message
  },

  /**
   * DELETE /api/offer/:id - Delete an offer
   */
  delete: async (id: string): Promise<string> => {
    const response = await api.delete<ApiResponse<any>>(`/offer/${id}`)
    return response.data.message
  }
}
