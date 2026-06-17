<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { roomService, type Room } from '@/services/roomService'

const route = useRoute()
const rooms = ref<Room[]>([])
const selectedRoomId = ref<string>('')
const checkIn = ref('')
const checkOut = ref('')
const guests = ref(1)
const name = ref('')
const email = ref('')
const submitted = ref(false)

onMounted(async () => {
  try {
    rooms.value = await roomService.getRooms()
    if (route.query.roomId) {
      selectedRoomId.value = String(route.query.roomId)
    } else if (rooms.value.length > 0 && rooms.value[0]) {
      selectedRoomId.value = String(rooms.value[0].id)
    }
  } catch (error) {
    console.error('Error loading rooms:', error)
  }
})

const handleSubmit = () => {
  submitted.value = true
}
</script>

<template>
  <div class="max-w-3xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
    <div class="bg-card text-card-foreground border border-border rounded-lg p-8 shadow-sm">
      <h1 class="text-3xl font-extrabold tracking-tight text-center mb-8">Room Reservation</h1>
      
      <div v-if="submitted" class="text-center py-8">
        <div class="inline-flex items-center justify-center h-16 w-16 rounded-full bg-green-100 text-green-600 mb-4">
          <svg class="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h2 class="text-2xl font-bold">Booking Submitted!</h2>
        <p class="text-muted-foreground mt-2">Thank you, {{ name }}. We have received your booking and sent a confirmation to {{ email }}.</p>
        <button @click="submitted = false" class="mt-6 inline-flex items-center justify-center rounded-md bg-primary text-primary-foreground hover:bg-primary/90 px-4 py-2 text-sm font-medium">
          Make Another Booking
        </button>
      </div>

      <form v-else @submit.prevent="handleSubmit" class="space-y-6">
        <div>
          <label for="room" class="block text-sm font-medium text-foreground">Select Room</label>
          <select id="room" v-model="selectedRoomId" class="mt-1 block w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary">
            <option v-for="room in rooms" :key="room.id" :value="String(room.id)">
              {{ room.name }} - Rp {{ room.price.toLocaleString('id-ID') }} / night
            </option>
          </select>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <div>
            <label for="check-in" class="block text-sm font-medium text-foreground">Check-In Date</label>
            <input type="date" id="check-in" v-model="checkIn" required class="mt-1 block w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>
          <div>
            <label for="check-out" class="block text-sm font-medium text-foreground">Check-Out Date</label>
            <input type="date" id="check-out" v-model="checkOut" required class="mt-1 block w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>
        </div>

        <div>
          <label for="guests" class="block text-sm font-medium text-foreground">Number of Guests</label>
          <input type="number" id="guests" v-model="guests" min="1" max="10" required class="mt-1 block w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary" />
        </div>

        <div class="border-t border-border pt-6 space-y-6">
          <div>
            <label for="name" class="block text-sm font-medium text-foreground">Full Name</label>
            <input type="text" id="name" v-model="name" required placeholder="John Doe" class="mt-1 block w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>
          <div>
            <label for="email" class="block text-sm font-medium text-foreground">Email Address</label>
            <input type="email" id="email" v-model="email" required placeholder="john@example.com" class="mt-1 block w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>
        </div>

        <button type="submit" class="w-full inline-flex items-center justify-center rounded-md bg-primary text-primary-foreground hover:bg-primary/90 py-3 text-base font-semibold shadow-sm">
          Confirm Reservation
        </button>
      </form>
    </div>
  </div>
</template>
