import axios from 'axios'

const api = axios.create({
  // baseURL: 'https://citramas-production.up.railway.app/api/',
  baseURL: 'http://localhost:4000/api/',
  timeout: 10000,
})

// Request interceptor — attach auth token if available
api.interceptors.request.use(
  (config) => {
    // Prevent absolute path resolution from overriding the /api/ baseURL path
    if (config.url && config.url.startsWith('/')) {
      config.url = config.url.substring(1)
    }

    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)

// Response interceptor — normalize errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      const { status, data } = error.response
      console.error(`[API Error] ${status}:`, data?.message || data?.error || 'Unknown error')

      if (status === 401) {
        localStorage.removeItem('token')
        // Could redirect to login in the future
      }
    } else if (error.request) {
      console.error('[API Error] No response received:', error.message)
    } else {
      console.error('[API Error] Request setup failed:', error.message)
    }

    return Promise.reject(error)
  },
)

export default api
