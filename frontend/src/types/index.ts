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