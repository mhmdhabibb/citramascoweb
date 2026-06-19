<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { roomService, type Room } from '@/services/roomService'
import { reservationService } from '@/services/admin/reservationService'

const route = useRoute()
const rooms = ref<Room[]>([])
const selectedRoomId = ref<string>('')
const checkIn = ref('')
const checkOut = ref('')
const guests = ref(1)
const name = ref('')
const email = ref('')
const submitted = ref(false)

const isOffer = ref(false)
const offerCode = ref('')
const isSubmitting = ref(false)
const submitError = ref<string | null>(null)

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

const formatDateToApi = (dateStr: string) => {
  if (!dateStr) return ''
  const parts = dateStr.split('-')
  if (parts.length === 3) {
    return `${parts[2]}-${parts[1]}-${parts[0]}`
  }
  return dateStr
}

const handleSubmit = async () => {
  try {
    isSubmitting.value = true
    submitError.value = null

    const payload = {
      room_id: selectedRoomId.value,
      full_name: name.value.trim(),
      email: email.value.trim(),
      check_in_date: formatDateToApi(checkIn.value),
      check_out_date: formatDateToApi(checkOut.value),
      number_of_guest: Number(guests.value),
      is_offer: isOffer.value,
      offer_code: isOffer.value ? offerCode.value.trim().toUpperCase() : ''
    }

    await reservationService.create(payload)
    submitted.value = true
  } catch (err: any) {
    console.error('Failed to submit booking:', err)
    submitError.value = err.response?.data?.message || err.message || 'Failed to submit reservation. Please try again.'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="max-w-3xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
    <div class="booking-card">
      <h1 class="text-3xl font-extrabold tracking-tight text-center mb-8 booking-animate booking-delay-0">Room Reservation</h1>
      
      <!-- Success State -->
      <div v-if="submitted" class="text-center py-8">
        <div class="success-icon">
          <svg class="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h2 class="text-2xl font-bold booking-animate booking-delay-1">Booking Submitted!</h2>
        <p class="text-muted-foreground mt-2 booking-animate booking-delay-2">Thank you, {{ name }}. We have received your booking and sent a confirmation to {{ email }}.</p>
        <button @click="submitted = false" class="booking-btn-reset booking-animate booking-delay-3">
          Make Another Booking
        </button>
      </div>

      <!-- Form -->
      <form v-else @submit.prevent="handleSubmit" class="space-y-6">
        <div class="booking-animate booking-delay-1">
          <label for="room" class="block text-sm font-medium text-foreground">Select Room</label>
          <select id="room" v-model="selectedRoomId" class="booking-input mt-1">
            <option v-for="room in rooms" :key="room.id" :value="String(room.id)">
              {{ room.name }} - Rp {{ room.price.toLocaleString('id-ID') }} / night
            </option>
          </select>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6 booking-animate booking-delay-2">
          <div>
            <label for="check-in" class="block text-sm font-medium text-foreground">Check-In Date</label>
            <input type="date" id="check-in" v-model="checkIn" required class="booking-input mt-1" />
          </div>
          <div>
            <label for="check-out" class="block text-sm font-medium text-foreground">Check-Out Date</label>
            <input type="date" id="check-out" v-model="checkOut" required class="booking-input mt-1" />
          </div>
        </div>

        <div class="booking-animate booking-delay-3">
          <label for="guests" class="block text-sm font-medium text-foreground">Number of Guests</label>
          <input type="number" id="guests" v-model="guests" min="1" max="10" required class="booking-input mt-1" />
        </div>

        <div class="border-t border-border pt-6 space-y-6 booking-animate booking-delay-4">
          <div>
            <label for="name" class="block text-sm font-medium text-foreground">Full Name</label>
            <input type="text" id="name" v-model="name" required placeholder="John Doe" class="booking-input mt-1" />
          </div>
          <div>
            <label for="email" class="block text-sm font-medium text-foreground">Email Address</label>
            <input type="email" id="email" v-model="email" required placeholder="john@example.com" class="booking-input mt-1" />
          </div>
          
          <!-- Promo Code Section -->
          <div class="flex items-center gap-2 mt-4 pt-2">
            <input 
              type="checkbox" 
              id="is-offer" 
              v-model="isOffer" 
              class="w-4 h-4 rounded border-gray-300 text-[#e15b2b] focus:ring-[#e15b2b] accent-[#e15b2b]"
            />
            <label for="is-offer" class="text-sm font-medium text-foreground cursor-pointer">I have a promo code</label>
          </div>

          <div v-if="isOffer" class="booking-animate mt-4">
            <label for="offer-code" class="block text-sm font-medium text-foreground">Promo Code</label>
            <input 
              type="text" 
              id="offer-code" 
              v-model="offerCode" 
              placeholder="e.g. CITRA5050" 
              class="booking-input mt-1 code-input"
              required
            />
          </div>
        </div>

        <!-- Error Notification -->
        <div v-if="submitError" class="p-4 bg-red-50 border border-red-200 text-red-700 text-sm rounded-lg">
          {{ submitError }}
        </div>

        <button 
          type="submit" 
          :disabled="isSubmitting"
          class="booking-submit booking-animate booking-delay-5 disabled:opacity-50"
        >
          <span v-if="isSubmitting">Submitting...</span>
          <span v-else>Confirm Reservation</span>
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.booking-card {
  background: var(--card);
  color: var(--card-foreground);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 2.5rem;
  box-shadow: 0 4px 24px rgba(28, 22, 18, 0.06);
  animation: cardEnter 0.6s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  opacity: 0;
  transform: translateY(16px);
}

@keyframes cardEnter {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Staggered form fields */
.booking-animate {
  opacity: 0;
  transform: translateY(16px);
  animation: bookingSlideUp 0.6s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.booking-delay-0 { animation-delay: 0.15s; }
.booking-delay-1 { animation-delay: 0.25s; }
.booking-delay-2 { animation-delay: 0.35s; }
.booking-delay-3 { animation-delay: 0.45s; }
.booking-delay-4 { animation-delay: 0.55s; }
.booking-delay-5 { animation-delay: 0.65s; }

@keyframes bookingSlideUp {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Input styling with focus animation */
.booking-input {
  display: block;
  width: 100%;
  border-radius: 8px;
  border: 1px solid var(--input);
  background: var(--background);
  padding: 0.5rem 0.75rem;
  font-size: 14px;
  transition: border-color 0.3s ease, box-shadow 0.3s ease, transform 0.2s ease;
  outline: none;
}

.booking-input:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(28, 22, 18, 0.08);
  transform: translateY(-1px);
}

/* Success icon bounce */
.success-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 4rem;
  width: 4rem;
  border-radius: 50%;
  background: #dcfce7;
  color: #16a34a;
  margin-bottom: 1rem;
  animation: successBounce 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
  opacity: 0;
  transform: scale(0);
}

@keyframes successBounce {
  to {
    opacity: 1;
    transform: scale(1);
  }
}

/* Submit button */
.booking-submit {
  width: 100%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  background: var(--primary);
  color: var(--primary-foreground);
  padding: 0.75rem;
  font-size: 16px;
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(28, 22, 18, 0.1);
  border: none;
  cursor: pointer;
  transition: all 0.35s cubic-bezier(0.16, 1, 0.3, 1);
}

.booking-submit:hover {
  opacity: 0.9;
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(28, 22, 18, 0.15);
}

.booking-submit:active {
  transform: translateY(0);
}

/* Reset button */
.booking-btn-reset {
  margin-top: 1.5rem;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  background: var(--primary);
  color: var(--primary-foreground);
  padding: 0.5rem 1rem;
  font-size: 14px;
  font-weight: 500;
  border: none;
  cursor: pointer;
  transition: all 0.35s cubic-bezier(0.16, 1, 0.3, 1);
}

.booking-btn-reset:hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

.code-input {
  text-transform: uppercase;
  font-family: monospace;
  font-weight: 700;
  letter-spacing: 0.05em;
}
</style>
