<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { roomService, type Room } from '@/services/roomService'
import { reservationService } from '@/services/admin/reservationService'
import { useOfferStore } from '@/stores/offerStore'

const route = useRoute()
const rooms = ref<Room[]>([])
const offers = ref<any[]>([])
const offerStore = useOfferStore()

const selectedRoomId = ref<string>('')
const checkInDate = ref('')
const checkOutDate = ref('')
const guests = ref(1)
const fullName = ref('')
const email = ref('')
const promoCode = ref('')

const loading = ref(true)
const submitting = ref(false)
const errorMsg = ref('')
const successMsg = ref('')

const isOffer = ref(false)
const offerCode = ref('')
const isSubmitting = ref(false)
const submitError = ref<string | null>(null)

onMounted(async () => {
  try {
    const [fetchedRooms] = await Promise.all([
      roomService.getRooms(),
      offerStore.fetchOffers()
    ])
    rooms.value = fetchedRooms
    offers.value = offerStore.offers.filter((o: any) => o.code)

    if (route.query.roomId) {
      selectedRoomId.value = String(route.query.roomId)
    } else if (rooms.value.length > 0 && rooms.value[0]) {
      selectedRoomId.value = String(rooms!.value[0].id)
    }
  } catch (error) {
    console.error('Error loading data:', error)
  } finally {
    loading.value = false
  }
})

const selectedRoom = computed(() => {
  return rooms.value.find(r => String(r.id) === selectedRoomId.value)
})

const selectedOffer = computed(() => {
  if (!promoCode.value) return null
  return offers.value.find(o => o.code === promoCode.value) || null
})

const normalPrice = computed(() => selectedRoom.value ? selectedRoom.value.price : 0)

const discountAmount = computed(() => {
  if (!selectedOffer.value || !selectedRoom.value) return 0
  
  if (selectedOffer.value.discount && selectedOffer.value.discount > 0) {
    return normalPrice.value * (selectedOffer.value.discount / 100)
  } else if (selectedOffer.value.price && selectedOffer.value.price > 0) {
    return selectedOffer.value.price
  }
  return 0
})

const finalPrice = computed(() => {
  const final = normalPrice.value - discountAmount.value
  return final > 0 ? final : 0
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
  submitting.value = true
  errorMsg.value = ''
  successMsg.value = ''

  const payload = {
    room_id: selectedRoomId.value,
    full_name: fullName.value.trim(),
    email: email.value.trim(),
    check_in_date: formatDateToApi(checkInDate.value),
    check_out_date: formatDateToApi(checkOutDate.value),
    number_of_guest: Number(guests.value),
    is_offer: !!promoCode.value,
    offer_code: promoCode.value ? promoCode.value.trim().toUpperCase() : ''
  }

  try {
    const msg = await reservationService.create(payload)
    successMsg.value = msg || 'Reservation created successfully!'
  } catch (error: any) {
    console.error('Failed to submit booking:', error)
    errorMsg.value = error.response?.data?.message || error.message || 'Failed to submit reservation. Please try again.'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-[#FAF7F2] py-24">
    <div v-if="loading" class="flex justify-center items-center py-32">
      <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-[#1C1612]"></div>
    </div>

    <div v-else class="max-w-[1200px] mx-auto px-4 sm:px-6 lg:px-8">
      
      <!-- Header -->
      <div class="text-center mb-16 animate-fade-in">
        <div class="flex items-center justify-center gap-2 mb-4">
          <span class="text-[#8C7A6B] text-[10px]">◆</span>
          <span class="text-[#8C7A6B] text-xs font-bold tracking-[0.2em] uppercase">Reservation</span>
          <span class="text-[#8C7A6B] text-[10px]">◆</span>
        </div>
        <h1 class="text-4xl md:text-5xl font-serif italic text-[#1C1612] mb-4" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">
          Complete Your Stay
        </h1>
      </div>

      <!-- Success State -->
      <div v-if="successMsg" class="max-w-2xl mx-auto bg-white border border-[#EAE1D8] p-12 text-center rounded-2xl shadow-sm animate-fade-in">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-[#F5F1E8] text-[#8C7A6B] mb-6">
          <svg class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h2 class="text-3xl font-serif text-[#1C1612] mb-4">Reservation Confirmed</h2>
        <p class="text-gray-600 mb-8">{{ successMsg }}</p>
        <button @click="successMsg = ''" class="inline-flex items-center justify-center px-8 py-3.5 border border-[#1A1A1A] text-[12px] font-bold tracking-[0.1em] text-[#1A1A1A] uppercase hover:bg-[#1A1A1A] hover:text-white transition-colors duration-400 rounded-full">
          Make Another Booking
        </button>
      </div>

      <!-- 2-Column Booking Layout -->
      <div v-else class="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-16">
        
        <!-- Left Col: Form -->
        <div class="lg:col-span-7 xl:col-span-8 animate-slide-up">
          <form @submit.prevent="handleSubmit" class="space-y-8 bg-white p-8 md:p-12 rounded-[24px] border border-[#EAE1D8] shadow-sm">
            
            <div v-if="errorMsg" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm mb-6">
              {{ errorMsg }}
            </div>

            <!-- Guest Details -->
            <div>
              <h3 class="text-lg font-serif text-[#1C1612] mb-6 border-b border-[#EAE1D8] pb-3" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">Guest Details</h3>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
                <div>
                  <label class="block text-[11px] font-bold tracking-[0.1em] text-gray-500 uppercase mb-2">Full Name</label>
                  <input type="text" v-model="fullName" required class="w-full bg-[#FAF7F2] border border-[#EAE1D8] text-[#1C1612] rounded-lg px-4 py-3 text-sm focus:outline-none focus:border-[#8C7A6B] transition-colors" placeholder="John Doe">
                </div>
                <div>
                  <label class="block text-[11px] font-bold tracking-[0.1em] text-gray-500 uppercase mb-2">Email Address</label>
                  <input type="email" v-model="email" required class="w-full bg-[#FAF7F2] border border-[#EAE1D8] text-[#1C1612] rounded-lg px-4 py-3 text-sm focus:outline-none focus:border-[#8C7A6B] transition-colors" placeholder="john@example.com">
                </div>
              </div>
            </div>

            <!-- Stay Details -->
            <div>
              <h3 class="text-lg font-serif text-[#1C1612] mb-6 border-b border-[#EAE1D8] pb-3 mt-8" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">Stay Details</h3>
              <div class="space-y-6">
                <div>
                  <label class="block text-[11px] font-bold tracking-[0.1em] text-gray-500 uppercase mb-2">Select Room</label>
                  <select v-model="selectedRoomId" class="w-full bg-[#FAF7F2] border border-[#EAE1D8] text-[#1C1612] rounded-lg px-4 py-3 text-sm focus:outline-none focus:border-[#8C7A6B] transition-colors appearance-none">
                    <option v-for="room in rooms" :key="room.id" :value="String(room.id)">
                      {{ room.name }}
                    </option>
                  </select>
                </div>
                
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
                  <div>
                    <label class="block text-[11px] font-bold tracking-[0.1em] text-gray-500 uppercase mb-2">Check-In Date</label>
                    <input type="date" v-model="checkInDate" required class="w-full bg-[#FAF7F2] border border-[#EAE1D8] text-[#1C1612] rounded-lg px-4 py-3 text-sm focus:outline-none focus:border-[#8C7A6B] transition-colors">
                  </div>
                  <div>
                    <label class="block text-[11px] font-bold tracking-[0.1em] text-gray-500 uppercase mb-2">Check-Out Date</label>
                    <input type="date" v-model="checkOutDate" required class="w-full bg-[#FAF7F2] border border-[#EAE1D8] text-[#1C1612] rounded-lg px-4 py-3 text-sm focus:outline-none focus:border-[#8C7A6B] transition-colors">
                  </div>
                </div>

                <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
                  <div>
                    <label class="block text-[11px] font-bold tracking-[0.1em] text-gray-500 uppercase mb-2">Number of Guests</label>
                    <input type="number" v-model="guests" min="1" max="10" required class="w-full bg-[#FAF7F2] border border-[#EAE1D8] text-[#1C1612] rounded-lg px-4 py-3 text-sm focus:outline-none focus:border-[#8C7A6B] transition-colors">
                  </div>
                  <div>
                    <label class="block text-[11px] font-bold tracking-[0.1em] text-gray-500 uppercase mb-2">Promo Code (Optional)</label>
                    <select v-model="promoCode" class="w-full bg-[#FAF7F2] border border-[#EAE1D8] text-[#1C1612] rounded-lg px-4 py-3 text-sm focus:outline-none focus:border-[#8C7A6B] transition-colors appearance-none">
                      <option value="">-- No Promo Code --</option>
                      <option v-for="offer in offers" :key="offer.id" :value="offer.code">
                        {{ offer.code }} - {{ offer.title }}
                      </option>
                    </select>
                  </div>
                </div>
              </div>
            </div>

            <!-- Submit -->
            <div class="pt-6">
              <button type="submit" :disabled="submitting" class="w-full inline-flex items-center justify-center px-8 py-4 bg-[#1C1612] text-[12px] font-bold tracking-[0.1em] text-white uppercase hover:bg-[#3D3329] transition-colors duration-400 rounded-full disabled:opacity-70 disabled:cursor-not-allowed">
                <svg v-if="submitting" class="animate-spin -ml-1 mr-3 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                {{ submitting ? 'Processing...' : 'Confirm Reservation' }}
              </button>
            </div>
          </form>
        </div>

        <!-- Right Col: Summary -->
        <div class="lg:col-span-5 xl:col-span-4 animate-slide-up" style="animation-delay: 0.1s;">
          <div class="sticky top-32">
            <h3 class="text-lg font-serif text-[#1C1612] mb-6 border-b border-[#EAE1D8] pb-3" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">Booking Summary</h3>
            
            <div v-if="selectedRoom" class="bg-white rounded-[24px] border border-[#EAE1D8] overflow-hidden shadow-sm">
              <div class="h-48 overflow-hidden">
                <img :src="selectedRoom.image" :alt="selectedRoom.name" class="w-full h-full object-cover transition-transform duration-700 hover:scale-105" />
              </div>
              <div class="p-8">
                <h4 class="text-2xl font-serif text-[#1C1612] mb-2" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">{{ selectedRoom.name }}</h4>
                <p class="text-gray-500 text-sm mb-6">{{ selectedRoom.capacity }} Guests Maximum</p>
                
                <div class="flex justify-between items-end border-t border-[#EAE1D8] pt-6">
                  <span class="text-gray-500 text-[11px] font-bold tracking-[0.1em] uppercase">Rate</span>
                  <div class="text-right">
                    <span class="block text-2xl font-serif text-[#1C1612]" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">
                      Rp {{ selectedRoom.price.toLocaleString('id-ID') }}
                    </span>
                    <span class="text-gray-500 text-xs">/ night</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Price Summary Card -->
            <div v-if="selectedRoom" class="mt-6 bg-white rounded-[24px] border border-[#EAE1D8] p-8 shadow-sm">
              <h4 class="text-lg font-serif text-[#1C1612] mb-6" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">Price Summary</h4>
              
              <div class="space-y-4 text-[15px]">
                <div class="flex justify-between items-start text-gray-600 gap-4">
                  <span class="leading-tight">Normal Rate / night</span>
                  <span class="whitespace-nowrap">Rp {{ normalPrice.toLocaleString('id-ID') }}</span>
                </div>
                
                <div v-if="discountAmount > 0" class="flex justify-between items-start text-[#10b981] font-medium gap-4">
                  <span class="leading-tight break-all">Discount <span class="text-xs opacity-80 font-normal block mt-0.5">({{ promoCode }})</span></span>
                  <span class="whitespace-nowrap">- Rp {{ discountAmount.toLocaleString('id-ID') }}</span>
                </div>
                
                <div class="pt-5 mt-5 border-t border-[#EAE1D8] flex justify-between items-end gap-4">
                  <span class="text-[11px] font-bold tracking-[0.1em] uppercase text-gray-500 leading-tight">Total Rate<br>per night</span>
                  <div class="text-right whitespace-nowrap">
                    <span class="block text-3xl font-serif text-[#1C1612]" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">
                      Rp {{ finalPrice.toLocaleString('id-ID') }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
            
            <div v-else class="bg-white rounded-[24px] border border-[#EAE1D8] p-8 text-center text-gray-500 shadow-sm">
              Please select a room to see the summary.
            </div>
          </div>
        </div>
      </div>
  </div>
  </div>
</template>

<style scoped>
.animate-fade-in {
  opacity: 0;
  animation: fadeIn 0.8s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.animate-slide-up {
  opacity: 0;
  transform: translateY(20px);
  animation: slideUp 0.8s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes fadeIn {
  to { opacity: 1; }
}

@keyframes slideUp {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

select {
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
  background-position: right 0.5rem center;
  background-repeat: no-repeat;
  background-size: 1.5em 1.5em;
}

.code-input {
  text-transform: uppercase;
  font-family: monospace;
  font-weight: 700;
  letter-spacing: 0.05em;
}
</style>
