<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { roomService, type Room } from '@/services/roomService'
import { useScrollReveal } from '@/composables/useScrollReveal'

const route = useRoute()
const room = ref<Room | null>(null)
const loading = ref(true)

const { elementRef: heroRef, isVisible: heroVisible } = useScrollReveal(0.1)
const { elementRef: contentRef, isVisible: contentVisible } = useScrollReveal(0.1)
const { elementRef: cardRef, isVisible: cardVisible } = useScrollReveal(0.1)

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
      <section class="relative w-full h-[50vh] md:h-[65vh] overflow-hidden">
        <img :src="room.image" :alt="room.name" class="w-full h-full object-cover animate-image-reveal" />
        <div class="absolute inset-0 bg-black/20"></div>
      </section>

      <!-- Room Header -->
      <section ref="heroRef" class="max-w-[1000px] mx-auto text-center px-4 -mt-16 md:-mt-24 relative z-10" :class="{ 'reveal-visible': heroVisible }">
        <div class="bg-[#FAF7F2] pt-12 pb-8 px-6 md:px-12 rounded-t-3xl reveal-item reveal-delay-0">
          <div class="flex items-center justify-center gap-2 mb-4">
            <span class="text-[#8C7A6B] text-[10px]">◆</span>
            <span class="text-[#8C7A6B] text-xs font-bold tracking-[0.2em] uppercase">Room Details</span>
            <span class="text-[#8C7A6B] text-[10px]">◆</span>
          </div>
          <h1 class="text-4xl md:text-5xl lg:text-6xl font-serif italic text-[#1C1612] mb-4" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">
            {{ room.name }}
          </h1>
        </div>
      </section>

      <!-- Main Content -->
      <section class="max-w-[1200px] mx-auto px-4 sm:px-6 lg:px-8 mt-4 md:mt-8">
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-20">
          
          <!-- Left Column: Details -->
          <div ref="contentRef" class="lg:col-span-8" :class="{ 'reveal-visible': contentVisible }">
            <div class="reveal-item reveal-delay-0 mb-16">
              <h2 class="text-2xl font-serif text-[#1C1612] mb-6 border-b border-[#EAE1D8] pb-4" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">
                The Experience
              </h2>
              <p class="text-gray-600 leading-relaxed text-[15px] md:text-[16px]">
                {{ room.description }}
              </p>
            </div>

            <div class="reveal-item reveal-delay-1">
              <h2 class="text-2xl font-serif text-[#1C1612] mb-6 border-b border-[#EAE1D8] pb-4" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">
                Amenities & Features
              </h2>
              <ul class="grid grid-cols-1 md:grid-cols-2 gap-y-5 gap-x-8">
                <li v-for="feature in room.features" :key="feature" class="flex items-center text-gray-600 text-[15px]">
                  <span class="text-[#8C7A6B] text-[10px] mr-3">◆</span>
                  {{ feature }}
                </li>
              </ul>
            </div>
          </div>

          <!-- Right Column: Sticky Booking Card -->
          <div ref="cardRef" class="lg:col-span-4" :class="{ 'reveal-visible': cardVisible }">
            <div class="sticky top-32 bg-[#F5F1E8] border border-[#EAE1D8] rounded-[24px] p-8 shadow-sm reveal-item reveal-delay-2">
              <div class="text-center mb-8 border-b border-[#EAE1D8] pb-8">
                <p class="text-gray-500 text-[11px] font-bold tracking-[0.15em] uppercase mb-3">Starting from</p>
                <div class="flex justify-center items-baseline gap-1.5">
                  <span class="text-3xl lg:text-4xl font-serif text-[#1C1612]" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">
                    Rp {{ room.price.toLocaleString('id-ID') }}
                  </span>
                </div>
                <p class="text-gray-500 text-[13px] mt-2">/ night (excluding taxes)</p>
              </div>
              
              <div class="flex flex-col gap-4">
                <router-link 
                  :to="{ name: 'booking', query: { roomId: room.id } }" 
                  class="w-full inline-flex items-center justify-center px-8 py-4 bg-[#1C1612] text-[12px] font-bold tracking-[0.1em] text-white uppercase hover:bg-[#3D3329] transition-colors duration-400 rounded-full shadow-md"
                >
                  Reserve This Room
                </router-link>
                <router-link 
                  to="/rooms" 
                  class="w-full inline-flex items-center justify-center px-8 py-4 border border-[#1C1612] text-[12px] font-bold tracking-[0.1em] text-[#1C1612] uppercase hover:bg-black/5 transition-colors duration-400 rounded-full"
                >
                  Explore Other Rooms
                </router-link>
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
/* Scroll Reveal Animations */
.reveal-item {
  opacity: 0;
  transform: translateY(24px);
  transition: opacity 0.8s cubic-bezier(0.16, 1, 0.3, 1),
              transform 0.8s cubic-bezier(0.16, 1, 0.3, 1);
}

.reveal-visible .reveal-item {
  opacity: 1;
  transform: translateY(0);
}

.reveal-delay-0 { transition-delay: 0s; }
.reveal-delay-1 { transition-delay: 0.1s; }
.reveal-delay-2 { transition-delay: 0.2s; }

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
</style>
