import { dashboardService, type DashboardData } from '@/services/admin/dashboardService'
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useDashboardStore = defineStore('dashboard', () => {
  // States
  const dashboardData = ref<DashboardData | null | undefined>(null)
  const currentRange = ref<'today' | 'weekly' | 'monthly'>('weekly') // default ke weekly sesuai UI
  const isLoading = ref<boolean>(false)
  const errorMsg = ref<string | null>(null)

  // Getters (Computed properties)
  const occupancyRate = computed(() => dashboardData.value?.occupancy_rate ?? 0)
  const totalRooms = computed(() => dashboardData.value?.total_rooms ?? 0)
  const availableRooms = computed(() => dashboardData.value?.available_rooms ?? 0)

  const activeStats = computed(() => {
    // Jika data belum dimuat dari backend, return nilai fallback awal
    if (!dashboardData.value) {
      return { occupancy: 0, available: 0, revenue: 0 }
    }

    // Jika user memilih filter utama (weekly), gunakan langsung data asli backend
    if (currentRange.value === 'weekly') {
      return {
        occupancy: dashboardData.value.occupancy_rate,
        available: dashboardData.value.available_rooms,
        revenue: dashboardData.value.total_revenue,
      }
    }

    // Kamu bisa memberikan pengondisian data statis atau menghitungnya di sini
    // Contoh jika menggunakan data default dari state range saat ini:
    return {
      occupancy: dashboardData.value.occupancy_rate,
      available: dashboardData.value.available_rooms,
      revenue: dashboardData.value.total_revenue,
    }
  })

  // Format Revenue menjadi Rupiah otomatis di level store
  const formattedRevenue = computed(() => {
    const revenue = dashboardData.value?.total_revenue ?? 0
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
    }).format(activeStats.value.revenue)
  })

  const bookingBreakdown = computed(() => dashboardData.value?.booking_breakdown)
  const revenueChart = computed(() => dashboardData.value?.revenue_chart ?? [])

  // Actions
  async function fetchDashboard(range?: 'today' | 'weekly' | 'monthly') {
    if (range) {
      currentRange.value = range
    }
    isLoading.value = true
    errorMsg.value = null

    try {
      const response = await dashboardService.getSummary(currentRange.value)
      dashboardData.value = response.data
    } catch (error: any) {
      console.error('Gagal memuat data dashboard:', error)
      errorMsg.value = error.response?.data?.error || 'Terjadi kesalahan sistem'
    } finally {
      isLoading.value = false
    }
  }

  return {
    dashboardData,
    currentRange,
    isLoading,
    errorMsg,
    occupancyRate,
    totalRooms,
    availableRooms,
    formattedRevenue,
    bookingBreakdown,
    revenueChart,
    fetchDashboard,
  }
})
