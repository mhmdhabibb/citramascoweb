<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { roomService } from '@/services/roomService'
import { reservationService } from '@/services/admin/reservationService'
import { useToastStore } from '@/stores/toastStore'
import type { Room } from '@/types'

const router = useRouter()
const toastStore = useToastStore()

const rooms = ref<Room[]>([])
const loadingRooms = ref(true)
const submitting = ref(false)

const form = ref({
  full_name: '',
  email: '',
  room_id: '',
  check_in_date: '',
  check_out_date: '',
  number_of_guest: 1,
  deposit: 0,
  offer_code: '',
})

const today = new Date()
const todayStr = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(
  today.getDate(),
).padStart(2, '0')}`

const selectedRoom = computed(() => rooms.value.find((r) => r.id === form.value.room_id))

const maxDeposit = computed(() => {
  if (!selectedRoom.value || !form.value.check_in_date || !form.value.check_out_date) return 0
  const nights = Math.max(
    1,
    Math.round(
      (new Date(form.value.check_out_date).getTime() -
        new Date(form.value.check_in_date).getTime()) /
        86400000,
    ),
  )
  return selectedRoom.value.price * nights
})

const submit = async () => {
  if (submitting.value) return

  if (!form.value.full_name || !form.value.email || !form.value.room_id) {
    toastStore.error('Guest name, email and room are required.')
    return
  }
  if (!form.value.check_in_date || !form.value.check_out_date) {
    toastStore.error('Check-in and check-out dates are required.')
    return
  }
  if (new Date(form.value.check_out_date) <= new Date(form.value.check_in_date)) {
    toastStore.error('Check-out date must be after check-in date.')
    return
  }

  submitting.value = true
  try {
    const depositNum = Number(form.value.deposit || 0)
    await reservationService.create({
      room_id: form.value.room_id,
      full_name: form.value.full_name,
      email: form.value.email,
      check_in_date: form.value.check_in_date,
      check_out_date: form.value.check_out_date,
      number_of_guest: Number(form.value.number_of_guest) || 1,
      deposit: depositNum,
      is_offer: form.value.offer_code ? true : false,
      offer_code: form.value.offer_code || undefined,
    })

    toastStore.success(
      depositNum > 0
        ? 'Reservation created. Deposit sent to the finance team.'
        : 'Reservation created successfully!',
    )
    router.push('/admin/reservations')
  } catch (error) {
    toastStore.error((error as Error).message || 'Failed to create reservation')
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  try {
    loadingRooms.value = true
    rooms.value = await roomService.getAll()
  } finally {
    loadingRooms.value = false
  }
})
</script>

<template>
  <div class="form-view">
    <div class="form-header">
      <div>
        <h1>New Reservation</h1>
        <p class="subtitle">Record a new booking at the front desk, including any deposit.</p>
      </div>
      <button class="btn-ghost" @click="router.push('/admin/reservations')">Back to Reservations</button>
    </div>

    <div class="form-card">
      <div class="section-title">Guest Information</div>
      <div class="form-grid">
        <div class="field">
          <label>Guest Name*</label>
          <input v-model="form.full_name" type="text" placeholder="Guest full name" />
        </div>
        <div class="field">
          <label>Email*</label>
          <input v-model="form.email" type="email" placeholder="guest@email.com" />
        </div>
      </div>

      <div class="section-title">Stay Details</div>
      <div class="form-grid">
        <div class="field">
          <label>Room*</label>
          <select v-model="form.room_id">
            <option value="" disabled>Select room</option>
            <option v-for="room in rooms" :key="room.id" :value="room.id">
              {{ room.name }} — Rp{{ (room.price || 0).toLocaleString('id-ID') }}/night
            </option>
          </select>
        </div>
        <div class="field">
          <label>Number of Guests</label>
          <input v-model.number="form.number_of_guest" type="number" min="1" />
        </div>
        <div class="field">
          <label>Check-in Date</label>
          <input v-model="form.check_in_date" type="date" :min="todayStr" />
        </div>
        <div class="field">
          <label>Check-out Date</label>
          <input v-model="form.check_out_date" type="date" :min="form.check_in_date || todayStr" />
        </div>
      </div>

      <div class="section-title">Deposit & Offer</div>
      <div class="form-grid">
        <div class="field">
          <label>Deposit (IDR)</label>
          <input v-model="form.deposit" type="number" min="0" placeholder="0" />
          <p v-if="maxDeposit > 0" class="field-hint">
            Max deposit is {{ maxDeposit.toLocaleString('id-ID') }} (full stay). A deposit is
            automatically sent to the finance team.
          </p>
        </div>
        <div class="field">
          <label>Offer Code (optional)</label>
          <input v-model="form.offer_code" type="text" placeholder="e.g. DISKON10" />
        </div>
      </div>

      <div class="form-actions">
        <button class="btn-secondary" type="button" @click="router.push('/admin/reservations')">
          Cancel
        </button>
        <button class="btn-primary" type="button" :disabled="submitting" @click="submit">
          {{ submitting ? 'Creating...' : 'Create Reservation' }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.form-view {
  display: flex;
  flex-direction: column;
  gap: 20px;
  background-color: #f8fafc;
  font-family: 'Plus Jakarta Sans', sans-serif;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header h1 {
  font-size: 1.5rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}
.subtitle {
  color: #64748b;
  font-size: 0.9rem;
  margin: 4px 0 0;
}
.btn-ghost {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #475569;
  padding: 9px 16px;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
}
.btn-ghost:hover {
  background: #f8fafc;
  color: #0f172a;
}

.form-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  padding: 24px;
  max-width: 860px;
}
.section-title {
  font-size: 0.8rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: #94a3b8;
  margin: 22px 0 10px;
}
.section-title:first-of-type {
  margin-top: 0;
}
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}
.field label {
  display: block;
  font-size: 0.8rem;
  font-weight: 600;
  color: #475569;
  margin-bottom: 6px;
}
.field input,
.field select {
  width: 100%;
  box-sizing: border-box;
  padding: 10px 14px;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  background: #f8fafc;
  color: #1e293b;
  font-size: 0.875rem;
  outline: none;
  transition: all 0.2s;
}
.field input:focus,
.field select:focus {
  border-color: #e4793b;
  background: #ffffff;
  box-shadow: 0 0 0 3px rgba(228, 121, 59, 0.08);
}
.field-hint {
  font-size: 0.75rem;
  color: #94a3b8;
  margin: 6px 0 0;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}
.btn-primary {
  background: #e4793b;
  color: white;
  border: none;
  padding: 11px 22px;
  border-radius: 10px;
  font-weight: 700;
  font-size: 0.9rem;
  cursor: pointer;
}
.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.btn-secondary {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #475569;
  padding: 11px 22px;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
}
</style>