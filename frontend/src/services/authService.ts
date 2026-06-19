import api from './api'
import type { ApiResponse } from '@/types'

export interface LoginResponse {
  access_token: string
}

export const authService = {
  /**
   * POST /api/auth/sign-in — Log in user
   */
  login: async (username: string, password: string): Promise<LoginResponse> => {
    const response = await api.post<ApiResponse<LoginResponse>>('/auth/sign-in', {
      username,
      password,
    })
    if (!response.data.success || !response.data.data) {
      throw new Error(response.data.message || 'Login failed')
    }
    return response.data.data
  },

  /**
   * GET /api/auth/me — Retrieve currently logged-in user profile
   */
  getProfile: async (): Promise<any> => {
    const response = await api.get<ApiResponse<any>>('/auth/me')
    if (!response.data.success || !response.data.data) {
      throw new Error(response.data.message || 'Failed to retrieve profile')
    }
    return response.data.data
  },

  /**
   * POST /api/auth/sign-out — Log out user
   */
  logout: async (): Promise<void> => {
    try {
      await api.post('/auth/sign-out')
    } finally {
      localStorage.removeItem('token')
    }
  },
}
