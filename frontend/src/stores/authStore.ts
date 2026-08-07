import { authService } from '@/services/authService'
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const MOCK_DEV_TOKEN = 'mock-developer-token-citramas'

interface ProfileUser {
  id?: string
  first_name: string
  last_name?: string
  username?: string
  email?: string
  role: string
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<ProfileUser | null>(null)
  const loaded = ref(false)
  const loadingProfile = ref(false)

  const role = computed(() => user.value?.role || '')
  const isAuthenticated = computed(() => !!user.value)

  async function fetchProfile(stale = false): Promise<ProfileUser | null> {
    // if stale is false and we already have a cached profile, return it
    if (!stale && loaded.value && user.value) return user.value
    if (loadingProfile.value) return user.value

    const token = localStorage.getItem('token')
    if (!token) {
      user.value = null
      loadingProfile.value = false
      loaded.value = true
      return null
    }

    loadingProfile.value = true

    if (token === MOCK_DEV_TOKEN) {
      user.value = {
        first_name: 'Developer',
        last_name: 'Mock',
        username: 'admin',
        email: 'admin@citramas.com',
        role: 'admin',
      }
      loadingProfile.value = false
      loaded.value = true
      return user.value
    }

    try {
      const data = await authService.getProfile()
      user.value = data as ProfileUser
    } catch {
      user.value = null
    }
    loadingProfile.value = false
    loaded.value = true
    return user.value
  }

  function reset() {
    user.value = null
    loaded.value = false
  }

  return { user, loaded, loadingProfile, role, isAuthenticated, fetchProfile, reset }
})