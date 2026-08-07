// ============================================
// API Response wrapper (matches backend gin.H format)
// ============================================
export interface ApiResponse<T> {
  success: boolean
  message: string
  data?: T
  error?: string
}

// ============================================
// Category (from backend: internal/modules/category/entity.go)
// ============================================
export interface Category {
  id: string
  name: string
  slug: string
  is_deleted: boolean
  created_at: string
  updated_at: string
  deleted_at: string | null
}

// ============================================
// RoomType (from backend: internal/modules/types/entity.go)
// ============================================
export interface RoomType {
  id: string
  name: string
  slug: string
  created_at: string
  updated_at: string
}

// ============================================
// Room (mock data — TODO: replace when backend Room module exists)
// ============================================
export interface Room {
  id: string
  code: string
  name: string
  slug: string
  image: string
  category_id: string
  category?: Category
  price: number
  capacity: number
  size: number
  type_id: string
  type?: RoomType
  description: string
  status?: string
  availability_status?: string
  features?: string[]
}

export interface Reservation {
  id: string
  code: string
  full_name: string
  email: string
  room_id: string
  status: string
  number_of_guest: number
  checkin_date: string
  checkout_date: string
  room: {
    name: string
  }
  price: number
  total_night: number
  total_price: number
  deposit?: number
  is_offer: boolean
  offer_code: string
}

export interface Offer {
  id: string
  title: string
  image: string
  status: string
  price?: number | null
  discount?: number | null
  code: string
  discounteed?: number | null
  valid_start?: string | null
  valid_end?: string | null
  max_quota?: number | null
  description?: string | null
  created_at?: string
  updated_at?: string
}

export type UserRole =
  | 'admin'
  | 'manager'
  | 'user'
  | 'reception'
  | 'finance'
  | 'inventory'

export interface User {
  id: string
  first_name: string
  last_name: string
  username: string
  phone: string
  email: string
  address: string
  role: UserRole
  created_at: string
  updated_at: string
}

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

export type InventoryTransactionType = 'stock-in' | 'usage' | 'adjustment'

export interface InventoryItem {
  id: string
  name: string
  category: string
  unit: string
  current_stock: number
  reorder_level: number
  status: string
  created_at: string
  updated_at: string
}

export interface InventoryTransaction {
  id: string
  item_id: string
  item: InventoryItem
  type: InventoryTransactionType
  quantity: number
  balance_after: number
  reference: string
  note: string
  created_by: string
  created_at: string
}

export interface InventoryReportItem {
  item: {
    id: string
    name: string
    category: string
    unit: string
    reorder_level: number
  }
  opening_stock: number
  stock_in: number
  usage: number
  closing_stock: number
  is_low_stock: boolean
}

export interface InventoryReport {
  from: string
  to: string
  summary: InventoryReportItem[]
  low_stock: InventoryReportItem[]
}

export interface AppNotification {
  id: string
  user_id: string
  type: string
  title: string
  message: string
  reference_id: string
  is_read: boolean
  created_at: string
}
