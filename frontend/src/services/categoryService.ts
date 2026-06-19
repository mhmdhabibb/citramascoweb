import api from './api'
import type { ApiResponse, Category } from '@/types'

export const categoryService = {
  /**
   * GET /api/categories — Fetch all categories
   */
  getAll: async (): Promise<Category[]> => {
    const response = await api.get<ApiResponse<Category[]>>('/categories')
    return response.data.data ?? []
  },

  /**
   * GET /api/category/:slug — Fetch category by slug
   */
  getBySlug: async (slug: string): Promise<Category> => {
    const response = await api.get<ApiResponse<Category>>(`/category/${slug}`)
    if (!response.data.data) {
      throw new Error(response.data.message || 'Category not found')
    }
    return response.data.data
  },

  /**
   * GET /api/category/detail/:id — Fetch category by ID
   */
  getById: async (id: string): Promise<Category> => {
    const response = await api.get<ApiResponse<Category>>(`/category/detail/${id}`)
    if (!response.data.data) {
      throw new Error(response.data.message || 'Category not found')
    }
    return response.data.data
  },

  /**
   * POST /api/category/ — Create a new category
   */
  create: async (data: { name: string }): Promise<string> => {
    const response = await api.post<ApiResponse<any>>('/category/', data)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Validation failed!')
    }
    return response.data.message
  },

  /**
   * PATCH /api/category/:id — Update a category
   */
  update: async (id: string, data: { name: string }): Promise<string> => {
    const response = await api.patch<ApiResponse<any>>(`/category/${id}`, data)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Validation failed!')
    }
    return response.data.message
  },

  /**
   * DELETE /api/category/:id — Delete a category
   */
  delete: async (id: string): Promise<string> => {
    const response = await api.delete<ApiResponse<any>>(`/category/${id}`)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Deletion failed!')
    }
    return response.data.message
  },
}
