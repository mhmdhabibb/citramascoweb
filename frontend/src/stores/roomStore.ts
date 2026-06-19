import { roomService } from "@/services/roomService";
import type { Room } from "@/types";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useRoomStore = defineStore('room', () => {
    const rooms = ref<Room[]>([])
    const loading = ref(false)
    const error = ref<string | null>(null)

    async function fetchRooms() {
        loading.value = true
        error.value = null
        try {
            rooms.value = await roomService.getRooms()
        } catch (err: any) {
            loading.value = false
            error.value = err.response?.data?.message || err.message || 'Failed to fetch rooms '
        } finally {
            loading.value = false
        }

    }

    async function storeRoom(payload: any) {
        loading.value = true
        error.value = null
        try {
            await roomService.create(payload)
        } catch (err: any) {
            loading.value = false
            error.value = err.response?.data?.message || err.message || 'Failed to create room'
        } finally {
            loading.value = false
        }
    }

    return {
        rooms,
        loading,
        error,
        fetchRooms,
        storeRoom
    }
})