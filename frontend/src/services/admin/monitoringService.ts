import api from '../api'
import type { ApiResponse, MonitoringReport } from '@/types'

export const monitoringService = {
  /**
   * GET /api/monitoring/services
   * Returns comprehensive service status, per-category metrics, and database connection health
   */
  getServicesStatus: async (): Promise<MonitoringReport> => {
    const response = await api.get<ApiResponse<MonitoringReport>>('/monitoring/services')
    if (!response.data.data) {
      throw new Error(response.data.message || 'Failed to fetch service monitoring report')
    }
    return response.data.data
  },

  /**
   * GET /api/monitoring/health
   * Lightweight quick health check
   */
  getHealth: async (): Promise<{ status: string; database: string; uptime: string; error_rate: number }> => {
    const response = await api.get('/monitoring/health')
    return response.data
  },

  /**
   * GET /metrics?format=json
   * Returns raw metrics in JSON structure
   */
  getMetricsJSON: async (): Promise<MonitoringReport> => {
    const response = await api.get<ApiResponse<MonitoringReport>>('/metrics', {
      params: { format: 'json' },
    })
    return response.data.data!
  },

  /**
   * GET /metrics
   * Returns raw Prometheus text format string
   */
  getMetricsPrometheus: async (): Promise<string> => {
    const response = await api.get<string>('/metrics', {
      headers: { Accept: 'text/plain' },
      responseType: 'text',
    })
    return response.data
  },
}

