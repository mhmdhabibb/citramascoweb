import api from '../api'
import type { ApiResponse } from '@/types'

// 1. Definisikan sub-interface detail sesuai JSON Golang
export interface BookingBreakdown {
  booked: number
  canceled: number
  pending_confirmation: number
  no_show: number
  total_bookings: number
  booked_percentage: number
  canceled_percentage: number
  pending_percentage: number
  no_show_percentage: number
}

export interface RevenueMonthly {
  month: string
  income: number
}

export interface DashboardData {
  occupancy_rate: number
  available_rooms: number
  total_rooms: number
  total_revenue: number
  booking_breakdown: BookingBreakdown
  revenue_chart: RevenueMonthly[]
}

export const dashboardService = {
  // 2. Gunakan tipe generik ApiResponse<DashboardData> jika interface ApiResponse kamu mendukung generik
  async getSummary(range: 'today' | 'weekly' | 'monthly'): Promise<ApiResponse<DashboardData>> {
    const response = await api.get(`/dashboard/summary`, {
      params: { range }, // Menggunakan objek params lebih rapi dibanding string literal template
    })

    return response.data
  },
}
