import api from '../api'
import type { ApiResponse, User, UserRole } from '@/types'

export const userService = {
  /**
   * GET /api/customers
   * List of customers (role = user)
   */
  getCustomers: async (): Promise<User[]> => {
    const response = await api.get<ApiResponse<User[]>>('/customers')
    if (!response.data.data) {
      throw new Error(response.data.message || 'Failed to fetch customers')
    }
    return response.data.data
  },

  /**
   * GET /api/users/?role={role}
   * List of users by role
   */
  getByRole: async (role: UserRole): Promise<User[]> => {
    const response = await api.get<ApiResponse<User[]>>('/users/', {
      params: { role },
    })
    if (!response.data.data) {
      throw new Error(response.data.message || 'Failed to fetch users')
    }
    return response.data.data
  },

  /**
   * POST /api/users
   * Create a new user with an assigned role
   */
  create: async (data: {
    first_name: string
    last_name: string
    username: string
    password: string
    phone: string
    email: string
    address?: string
    role: UserRole
  }): Promise<User> => {
    const response = await api.post<ApiResponse<User>>('/users/', data)
    if (!response.data.success || !response.data.data) {
      throw new Error(response.data.message || 'Failed to create user')
    }
    return response.data.data
  },

  /**
   * PATCH /api/users/:id/role
   * Assign/update a user's role
   */
  updateRole: async (id: string, role: UserRole): Promise<string> => {
    const response = await api.patch<ApiResponse<{ message: string }>>(`/users/${id}/role`, {
      role,
    })
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to update user role')
    }
    return response.data.message
  },

  /**
   * DELETE /api/users/:id
   */
  delete: async (id: string): Promise<string> => {
    const response = await api.delete<ApiResponse<{ message: string }>>(`/users/${id}`)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to delete user')
    }
    return response.data.message
  },
}
