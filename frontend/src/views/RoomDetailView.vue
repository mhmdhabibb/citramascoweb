<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { roomService, type Room } from '@/services/roomService'

const route = useRoute()
const room = ref<Room | null>(null)
const loading = ref(true)

const policies = [
  {
    title: 'Check-In & Check-Out',
    items: [
      'Check-In: From 3:00 PM',
      'Check-Out: By 11:00 AM',
      'Early check-in and late check-out may be available upon request (fees may apply).'
    ]
  },
  {
    title: 'Cancellation Policy',
    items: [
      'Flexible Rate: Free cancellation up to 72 hours before arrival.',
      'Non-Refundable Rate: No refunds or changes allowed after booking.',
      'No-shows will be charged the full stay.',
      'We recommend travel insurance for peace of mind.'
    ]
  },
  {
    title: 'Pet Policy',
    items: [
      'Small pets (under 20 lbs) are welcome in select rooms.',
      'A one-time $75 cleaning fee applies per pet.',
      'Limit: 2 pets per room.',
      'Pets must be leashed in public areas and cannot be left unattended.'
    ]
  },
  {
    title: 'Smoking Policy',
    items: [
      'All rooms and suites are 100% smoke-free.',
      'A $250 cleaning fee will apply if smoking is detected in the room.',
      'Designated outdoor smoking areas are available near the garden terrace.'
    ]
  }
]

const openPolicies = ref<number[]>([0, 1, 2, 3]) // All open by default
const togglePolicy = (index: number) => {
  if (openPolicies.value.includes(index)) {
    openPolicies.value = openPolicies.value.filter(i => i !== index)
  } else {
    openPolicies.value.push(index)
  }
}

const fetchRoom = async () => {
  loading.value = true
  const roomId = String(route.params.id)
  try {
    const data = await roomService.getRoomById(roomId)
    if (data) {
      room.value = data
    }
  } catch (error) {
    console.error('Error fetching room details:', error)
  } finally {
    loading.value = false
  }
}

onMounted(fetchRoom)
watch(() => route.params.id, fetchRoom)
</script>

<template>
  <div class="min-h-screen pb-24 bg-[#FAF7F2]">
    <div v-if="loading" class="flex justify-center items-center py-32">
      <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-[#1C1612]"></div>
    </div>

    <div v-else-if="room">
      <!-- Full-Width Hero Section -->
      <section class="relative w-full h-[60vh] md:h-[75vh] overflow-hidden flex items-center justify-center">
        <img :src="room.image" :alt="room.name" class="absolute inset-0 w-full h-full object-cover animate-image-reveal" />
        <div class="absolute inset-0 bg-black/40"></div>
        
        <!-- Room Header Text inside Hero -->
        <div class="relative z-10 text-center px-4 max-w-3xl mx-auto mt-16">
          <h1 class="text-4xl md:text-6xl lg:text-7xl font-serif text-white mb-6 animate-slide-up" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">
            {{ room.name }}
          </h1>
        </div>
      </section>

      <!-- Main Content -->
      <section class="max-w-[1200px] mx-auto px-4 sm:px-6 lg:px-8 mt-16 md:mt-24 mb-16 md:mb-24">
        <!-- About Section -->
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-16 mb-24 animate-slide-up-delayed">
          <div class="lg:col-span-3 md:pt-2">
            <h2 class="text-4xl md:text-5xl font-serif text-[#1F2927]" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">
              About
            </h2>
          </div>
          <div class="lg:col-span-9">
            <!-- Meta Info Row -->
            <div class="flex flex-wrap items-center gap-6 md:gap-10 text-[#4A5553] text-[14px] font-medium mb-6">
              <div class="flex items-center gap-2.5">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="15 3 21 3 21 9"></polyline>
                  <polyline points="9 21 3 21 3 15"></polyline>
                  <line x1="21" y1="3" x2="14" y2="10"></line>
                  <line x1="3" y1="21" x2="10" y2="14"></line>
                </svg>
                <span>{{ room.size }} Sq.m</span>
              </div>
              <div class="flex items-center gap-2.5" v-if="room.type?.name">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M2 4v16"></path>
                  <path d="M2 8h18a2 2 0 0 1 2 2v10"></path>
                  <path d="M2 17h20"></path>
                  <path d="M6 8v9"></path>
                </svg>
                <span>{{ room.type.name }}</span>
              </div>
              <div class="flex items-center gap-2.5">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                  <circle cx="12" cy="7" r="4"></circle>
                </svg>
                <span>{{ room.capacity }} Guests</span>
              </div>
            </div>

            <hr class="border-[#E0DCD4] mb-8">

            <!-- Description -->
            <p class="text-[#4A5553] leading-relaxed text-[15px] md:text-[15.5px] whitespace-pre-line mb-8 font-medium">
              {{ room.description }}
            </p>

            <hr class="border-[#E0DCD4] mb-8">

            <!-- Price & CTA -->
            <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-6">
              <div>
                <span class="block text-[#A3998C] text-[10px] uppercase tracking-[0.15em] font-bold mb-1">Starting at</span>
                <div class="flex items-baseline gap-1.5">
                  <span class="text-3xl font-serif text-[#1F2927]" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">
                    Rp {{ room.price.toLocaleString('id-ID') }}
                  </span>
                  <span class="text-[#6C7573] text-[13px] font-medium">/ night</span>
                </div>
              </div>
              <router-link 
                :to="{ name: 'booking', query: { roomId: room.id } }" 
                class="px-8 py-3.5 border border-[#1F2927] text-[12px] font-bold tracking-[0.1em] text-[#1F2927] uppercase hover:bg-[#1F2927] hover:text-white transition-colors duration-400 inline-block text-center"
              >
                Book a Stay
              </router-link>
            </div>
          </div>
        </div>

        <!-- Diamond Divider -->
        <div class="relative flex items-center justify-center mb-16 animate-slide-up-delayed" style="animation-delay: 0.8s">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-[#E0DCD4]"></div>
          </div>
          <div class="relative bg-[#FAF7F2] px-4">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#E0DCD4" stroke-width="1.5" transform="rotate(45)">
              <rect x="3" y="3" width="18" height="18"></rect>
            </svg>
          </div>
        </div>

        <!-- Policies Section -->
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-16 mb-24 animate-slide-up-delayed" style="animation-delay: 1s">
          <div class="lg:col-span-3 md:pt-2">
            <h2 class="text-4xl md:text-5xl font-serif text-[#1F2927]" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif; line-height: 1.1;">
              Policies &<br>Check-in Info
            </h2>
          </div>
          <div class="lg:col-span-9">
            <div class="flex flex-col">
              <div v-for="(policy, index) in policies" :key="index" class="border-t border-[#EAE1D8] last:border-b">
                <!-- Accordion Header -->
                <button 
                  @click="togglePolicy(index)" 
                  class="w-full py-6 flex justify-between items-center text-left focus:outline-none cursor-pointer"
                >
                  <span class="text-[#4A5553] font-semibold text-[15px]">{{ policy.title }}</span>
                  <span class="text-[#A3998C] text-xl font-light ml-4">
                    {{ openPolicies.includes(index) ? '×' : '+' }}
                  </span>
                </button>
                
                <!-- Accordion Content -->
                <div 
                  class="grid transition-all duration-500 ease-in-out"
                  :class="openPolicies.includes(index) ? 'grid-rows-[1fr] opacity-100' : 'grid-rows-[0fr] opacity-0'"
                >
                  <div class="overflow-hidden">
                    <div class="pb-6 pr-4 pt-2">
                      <div class="border-l-2 border-[#EAE1D8] pl-5 py-1">
                        <p v-for="(item, iIndex) in policy.items" :key="iIndex" class="text-[#6C7573] text-[13px] font-medium mb-4 last:mb-0">
                          {{ item }}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

      </section>
    </div>

    <!-- 404 Room Not Found -->
    <div v-else class="flex flex-col items-center justify-center py-32 text-center">
      <h2 class="text-3xl font-serif text-[#1C1612] mb-4">Room Not Found</h2>
      <p class="text-gray-600 mb-8">We couldn't find the room you were looking for.</p>
      <router-link to="/rooms" class="inline-flex items-center justify-center px-8 py-3.5 border border-[#1A1A1A] text-[12px] font-bold tracking-[0.1em] text-[#1A1A1A] uppercase hover:bg-[#1A1A1A] hover:text-white transition-colors duration-400 rounded-full">
        Return to Rooms
      </router-link>
    </div>
  </div>
</template>

<style scoped>
/* Initial Image Reveal Animation */
.animate-image-reveal {
  animation: zoomOutReveal 1.5s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes zoomOutReveal {
  0% {
    transform: scale(1.1);
    filter: blur(4px);
  }
  100% {
    transform: scale(1);
    filter: blur(0);
  }
}

.animate-slide-up {
  opacity: 0;
  transform: translateY(30px);
  animation: slideUpFade 1s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  animation-delay: 0.2s;
}

.animate-slide-up-delayed {
  opacity: 0;
  transform: translateY(30px);
  animation: slideUpFade 1s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  animation-delay: 0.4s;
}

@keyframes slideUpFade {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
