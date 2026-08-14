import { ref } from 'vue'
import { defineStore } from 'pinia'
import { channelService } from '@/services/admin/channelService'
import type { Channel } from '@/types'

export const useChannelStore = defineStore('channel', () => {
  const channels = ref<Channel[]>([])
  const currentChannel = ref<Channel | null>(null)
  const total = ref(0)
  const page = ref(1)
  const limit = ref(10)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchChannels() {
    loading.value = true
    error.value = null
    try {
      channels.value = await channelService.getAll()
      total.value = channels.value.length
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Failed to fetch channels'
      console.error('Failed to fetch channels:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchPaginated(targetPage = 1) {
    loading.value = true
    error.value = null
    try {
      const res = await channelService.get(targetPage)
      channels.value = res.channels
      total.value = res.total
      page.value = targetPage
      limit.value = res.limit
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Failed to fetch channels'
      console.error('Failed to fetch paginated channels:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchChannelById(id: string) {
    loading.value = true
    error.value = null
    try {
      currentChannel.value = await channelService.getById(id)
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Channel not found'
      currentChannel.value = null
    } finally {
      loading.value = false
    }
  }

  async function store(payload: { name: string }): Promise<string> {
    loading.value = true
    error.value = null
    try {
      return await channelService.create(payload)
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Failed to create channel'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function update(id: string, payload: { name: string }): Promise<string> {
    loading.value = true
    error.value = null
    try {
      return await channelService.update(id, payload)
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Failed to update channel'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function destroy(id: string): Promise<string> {
    loading.value = true
    error.value = null
    try {
      return await channelService.delete(id)
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Failed to delete channel'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    channels,
    currentChannel,
    total,
    page,
    limit,
    loading,
    error,
    fetchChannels,
    fetchPaginated,
    fetchChannelById,
    store,
    update,
    destroy,
  }
})
