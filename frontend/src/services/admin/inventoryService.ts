import api from '../api'
import type {
  ApiResponse,
  InventoryItem,
  InventoryReport,
  InventoryTransaction,
} from '@/types'

export const inventoryService = {
  /**
   * GET /api/inventory/items
   */
  getItems: async (): Promise<InventoryItem[]> => {
    const response = await api.get<ApiResponse<InventoryItem[]>>('/inventory/items')
    if (!response.data.data) {
      throw new Error(response.data.message || 'Failed to fetch inventory items')
    }
    return response.data.data
  },

  /**
   * POST /api/inventory/items
   */
  createItem: async (data: {
    name: string
    category: string
    unit: string
    current_stock?: number
    reorder_level?: number
  }): Promise<InventoryItem> => {
    const response = await api.post<ApiResponse<InventoryItem>>('/inventory/items', data)
    if (!response.data.success || !response.data.data) {
      throw new Error(response.data.message || 'Failed to create item')
    }
    return response.data.data
  },

  /**
   * DELETE /api/inventory/items/:id
   */
  deleteItem: async (id: string): Promise<string> => {
    const response = await api.delete<ApiResponse<{ message: string }>>(`/inventory/items/${id}`)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to delete item')
    }
    return response.data.message
  },

  /**
   * POST /api/inventory/items/:id/stock-in
   */
  stockIn: async (
    id: string,
    data: { quantity: number; reference?: string; note?: string },
  ): Promise<string> => {
    const response = await api.post<ApiResponse<{ message: string }>>(`/inventory/items/${id}/stock-in`, data)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to add stock')
    }
    return response.data.message
  },

  /**
   * POST /api/inventory/items/:id/usage
   */
  usage: async (
    id: string,
    data: { quantity: number; reference?: string; note?: string },
  ): Promise<string> => {
    const response = await api.post<ApiResponse<{ message: string }>>(`/inventory/items/${id}/usage`, data)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to record usage')
    }
    return response.data.message
  },

  /**
   * POST /api/inventory/stock-take
   */
  stockTake: async (data: {
    item_id: string
    actual_stock: number
    note?: string
  }): Promise<string> => {
    const response = await api.post<ApiResponse<{ message: string }>>('/inventory/stock-take', data)
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to record stock take')
    }
    return response.data.message
  },

  /**
   * GET /api/inventory/transactions
   */
  getTransactions: async (params?: {
    from?: string
    to?: string
  }): Promise<InventoryTransaction[]> => {
    const response = await api.get<ApiResponse<InventoryTransaction[]>>('/inventory/transactions', {
      params,
    })
    if (!response.data.data) {
      throw new Error(response.data.message || 'Failed to fetch transactions')
    }
    return response.data.data
  },

  /**
   * GET /api/inventory/report
   */
  getReport: async (params?: { from?: string; to?: string }): Promise<InventoryReport> => {
    const response = await api.get<ApiResponse<InventoryReport>>('/inventory/report', { params })
    if (!response.data.data) {
      throw new Error(response.data.message || 'Failed to fetch report')
    }
    return response.data.data
  },
}