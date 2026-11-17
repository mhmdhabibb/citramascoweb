<template>
  <div class="p-4">
    <h1 class="text-2xl font-bold mb-4">Daftar Reservasi</h1>
    <table class="min-w-full border">
      <thead>
        <tr class="bg-gray-100">
          <th class="px-2 py-1 border">ID</th>
          <th class="px-2 py-1 border">Tamu</th>
          <th class="px-2 py-1 border">Kamar</th>
          <th class="px-2 py-1 border">Check‑in</th>
          <th class="px-2 py-1 border">Check‑out</th>
          <th class="px-2 py-1 border">Status</th>
          <th class="px-2 py-1 border">Aksi</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="res in reservations" :key="res.id">
          <td class="px-2 py-1 border">{{ res.id }}</td>
          <td class="px-2 py-1 border">{{ res.guest_name }}</td>
          <td class="px-2 py-1 border">{{ res.room?.name || '-' }}</td>
          <td class="px-2 py-1 border">{{ formatDate(res.checkin_date) }}</td>
          <td class="px-2 py-1 border">{{ formatDate(res.checkout_date) }}</td>
          <td class="px-2 py-1 border capitalize">{{ res.status }}</td>
          <td class="px-2 py-1 border space-x-1">
            <button v-if="res.status === 'pending'" @click="approve(res.id)" class="btn btn-sm btn-success">Approve</button>
            <button v-if="res.status === 'pending'" @click="reject(res.id)" class="btn btn-sm btn-danger">Reject</button>
            <button v-if="res.status === 'approved'" @click="checkIn(res.id, res.room?.id)" class="btn btn-sm btn-primary">Check‑in</button>
            <button v-if="res.status === 'checked-in'" @click="checkOut(res.id, res.room?.id)" class="btn btn-sm btn-warning">Check‑out</button>
            <button v-if="['pending','approved','checked-in'].includes(res.status)" @click="cancel(res.id)" class="btn btn-sm btn-outline">Cancel</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { reservationService } from '@/services/admin/reservationService'
import { useToast } from '@/stores/toast'

interface Reservation {
  id: string
  guest_name: string
  room?: { id: string; name: string }
  checkin_date: string
  checkout_date: string
  status: string
}

const reservations = ref<Reservation[]>([])
const toast = useToast()

const load = async () => {
  try {
    const data = await reservationService.getAll()
    reservations.value = data
  } catch (e) {
    toast.error((e as Error).message)
  }
}

onMounted(load)

const approve = async (id: string) => {
  try { await reservationService.approve(id); await load() } catch (e) { toast.error((e as Error).message) }
}
const reject = async (id: string) => {
  try { await reservationService.reject(id); await load() } catch (e) { toast.error((e as Error).message) }
}
const cancel = async (id: string) => {
  try { await reservationService.cancel(id); await load() } catch (e) { toast.error((e as Error).message) }
}
const checkIn = async (id: string, roomId: string | undefined) => {
  if (!roomId) return toast.error('Room ID missing')
  try { await reservationService.checkIn(id); await load() } catch (e) { toast.error((e as Error).message) }
}
const checkOut = async (id: string, roomId: string | undefined) => {
  if (!roomId) return toast.error('Room ID missing')
  try { await reservationService.checkOut(id); await load() } catch (e) { toast.error((e as Error).message) }
}

const formatDate = (dateStr: string) => new Date(dateStr).toLocaleDateString()
</script>

<style scoped>
.btn { @apply px-2 py-1 rounded text-white; }
.btn-success { @apply bg-green-600; }
.btn-danger { @apply bg-red-600; }
.btn-primary { @apply bg-blue-600; }
.btn-warning { @apply bg-orange-600; }
.btn-outline { @apply bg-gray-300 text-gray-800; }
</style>
