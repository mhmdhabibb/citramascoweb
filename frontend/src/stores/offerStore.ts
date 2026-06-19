import { ref } from 'vue'
import { defineStore } from 'pinia'
import { offerService } from '@/services/offerService'
import type { Offer } from '@/types'

export const useOfferStore = defineStore('offer', () => {
  const offers = ref<Offer[]>([])
  const currentOffer = ref<Offer | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchOffers() {
    loading.value = true
    error.value = null
    try {
      offers.value = await offerService.getAll()
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Failed to fetch active offers'
      console.error('Failed to fetch active offers:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchAdminOffers() {
    loading.value = true
    error.value = null
    try {
      offers.value = await offerService.getAllAdmin()
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Failed to fetch offers list'
      console.error('Failed to fetch offers list:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchOfferById(id: string) {
    loading.value = true
    error.value = null
    try {
      currentOffer.value = await offerService.getById(id)
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Offer not found'
      currentOffer.value = null
    } finally {
      loading.value = false
    }
  }

  async function store(payload: FormData) {
    loading.value = true
    error.value = null
    try {
      await offerService.create(payload)
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Failed to create offer'
      console.error('Failed to create offer:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function update(id: string, payload: FormData) {
    loading.value = true
    error.value = null
    try {
      await offerService.update(id, payload)
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Failed to update offer'
      console.error('Failed to update offer:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function destroy(id: string) {
    loading.value = true
    error.value = null
    try {
      await offerService.delete(id)
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Failed to delete offer'
      console.error('Failed to delete offer:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    offers,
    currentOffer,
    loading,
    error,
    fetchOffers,
    fetchAdminOffers,
    fetchOfferById,
    store,
    update,
    destroy,
  }
})
