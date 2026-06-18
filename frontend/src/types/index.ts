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
  id: number
  name: string
  price: number
  image: string
  capacity: number
  description: string
  features: string[]
}
