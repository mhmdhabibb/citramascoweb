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
  create: async (formData: FormData): Promise<void> => {
    await api.post('/offer/', formData)
  },

  /**
   * PATCH /api/offer/:id - Update an offer (multipart form-data)
   */
  update: async (id: string, formData: FormData): Promise<void> => {
    await api.patch(`/offer/${id}`, formData)
  },

  /**
   * DELETE /api/offer/:id - Delete an offer
   */
  delete: async (id: string): Promise<void> => {
    await api.delete(`/offer/${id}`)
  }
}
