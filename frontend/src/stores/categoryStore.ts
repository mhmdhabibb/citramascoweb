import { ref } from 'vue'
import { defineStore } from 'pinia'
import { categoryService } from '@/services/categoryService'
import type { Category } from '@/types'

export const useCategoryStore = defineStore('category', () => {
  const categories = ref<Category[]>([])
  const currentCategory = ref<Category | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchCategories() {
    loading.value = true
    error.value = null
    try {
      categories.value = await categoryService.getAll()
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Failed to fetch categories'
      console.error('Failed to fetch categories:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchCategoryBySlug(slug: string) {
    loading.value = true
    error.value = null
    try {
      currentCategory.value = await categoryService.getBySlug(slug)
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Category not found'
      currentCategory.value = null
    } finally {
      loading.value = false
    }
  }

  async function fetchCategoryById(id: string) {
    loading.value = true
    error.value = null
    try {
      currentCategory.value = await categoryService.getById(id)
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Category not found'
      currentCategory.value = null
    } finally {
      loading.value = false
    }
  }

  async function store(payload: { name: string }): Promise<string> {
    loading.value = true
    error.value = null
    try {
      return await categoryService.create(payload)
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Failed to create category'
      throw err
    } finally {
      loading.value = false
    }
  }


  async function update(id: string, payload: { name: string }): Promise<string> {
    loading.value = true
    error.value = null
    try {
      return await categoryService.update(id, payload)
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Failed to update category'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function destroy(id: string): Promise<string> {
    loading.value = true
    error.value = null
    try {
      return await categoryService.delete(id)
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Failed to delete category'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    categories,
    currentCategory,
    loading,
    error,
    store,
    fetchCategories,
    update,
    destroy,
    fetchCategoryBySlug,
    fetchCategoryById,
  }
})
