import { defineStore } from 'pinia'
import { ref } from 'vue'
import { roomService } from '@/services/roomService'
import type { Room } from '@/types'

export const useRoomStore = defineStore('room', () => {
  const rooms = ref<Room[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 1. Mengambil Semua Data Kamar
  async function fetchRooms() {
    loading.value = true
    error.value = null
    try {
      rooms.value = await roomService.getRooms()
    } catch (err: any) {
      loading.value = false
      error.value = err.response?.data?.message || err.message || 'Failed to fetch rooms'
    } finally {
      loading.value = false
    }
  }

  // 2. Membuat Data Kamar Baru
  async function storeRoom(payload: FormData) {
    loading.value = true
    error.value = null
    try {
      await roomService.create(payload)
      await fetchRooms() // Auto-refresh data setelah berhasil menambahkan
    } catch (err: any) {
      loading.value = false
      error.value = err.response?.data?.message || err.message || 'Failed to create room'
      throw err
    } finally {
      loading.value = false
    }
  }

  // 3. AKSI BARU: Mengubah Status Operasional Fisik Kamar (Active, Maintenance, Dirty)
  async function updateRoomStatus(id: string, status: string) {
    loading.value = true
    error.value = null
    try {
      // Tembak endpoint PATCH /room/status/:id via roomService
      await roomService.updateStatus(id, status)

      // Sinkronisasi state lokal secara instan tanpa perlu reload page penuh
      const targetRoom = rooms.value.find((r) => r.id === id)
      if (targetRoom) {
        targetRoom.status = status
      }
    } catch (err: any) {
      error.value = err.response?.data?.message || err.message || 'Failed to update room status'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    rooms,
    loading,
    error,
    fetchRooms,
    storeRoom,
    updateRoomStatus, // Ekspor fungsi agar bisa dibaca komponen Vue
  }
})
