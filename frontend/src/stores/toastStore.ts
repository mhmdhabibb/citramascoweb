import { ref } from 'vue'
import { defineStore } from 'pinia'

export interface Toast {
  id: string
  message: string
  type: 'success' | 'error' | 'info' | 'warning'
  duration?: number
}

export const useToastStore = defineStore('toast', () => {
  const toasts = ref<Toast[]>([])

  const add = (message: string, type: 'success' | 'error' | 'info' | 'warning' = 'success', duration = 3000) => {
    const id = Math.random().toString(36).substring(2, 9)
    const newToast: Toast = { id, message, type, duration }
    toasts.value.push(newToast)

    setTimeout(() => {
      remove(id)
    }, duration)
  }

  const remove = (id: string) => {
    toasts.value = toasts.value.filter((t) => t.id !== id)
  }

  const success = (message: string, duration = 3000) => add(message, 'success', duration)
  const error = (message: string, duration = 3000) => add(message, 'error', duration)
  const info = (message: string, duration = 3000) => add(message, 'info', duration)
  const warning = (message: string, duration = 3000) => add(message, 'warning', duration)

  return {
    toasts,
    add,
    remove,
    success,
    error,
    info,
    warning,
  }
})
