import { ref } from 'vue'
import { defineStore } from 'pinia'
import { typeService } from '@/services/typeService'
import type { RoomType } from '@/types'

export const useTypeStore = defineStore('type', () => {
  const types = ref<RoomType[]>([])
  const currentType = ref<RoomType | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchTypes() {
    loading.value = true
    error.value = null
    try {
      types.value = await typeService.getAll()
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Failed to fetch types'
      console.error('Failed to fetch types:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchTypeById(id: string) {
    loading.value = true
    error.value = null
    try {
      currentType.value = await typeService.getById(id)
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Type not found'
      currentType.value = null
    } finally {
      loading.value = false
    }
  }

  return {
    types,
    currentType,
    loading,
    error,
    fetchTypes,
    fetchTypeById,
  }
})
