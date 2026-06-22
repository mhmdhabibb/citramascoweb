<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { roomService, type Room } from '@/services/roomService'
import { useScrollReveal } from '@/composables/useScrollReveal'
import RoomCard from '@/components/room/RoomCard.vue'

const rooms = ref<Room[]>([])
const loading = ref(true)

const { elementRef: sectionRef, isVisible: sectionVisible } = useScrollReveal(0.05)

onMounted(async () => {
  try {
    rooms.value = await roomService.getRooms()
  } catch (error) {
    console.error('Error fetching rooms:', error)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="bg-[#FAF7F2] min-h-screen py-24">
    <section ref="sectionRef" class="max-w-[1200px] mx-auto px-4 sm:px-6 lg:px-8">
      <div class="text-center max-w-2xl mx-auto mb-16" :class="{ 'reveal-visible': sectionVisible }">
        <div class="reveal-item reveal-delay-0">
          <div class="flex items-center justify-center gap-2 mb-4">
            <span class="text-[#B59E75] text-[10px]">◆</span>
            <span class="text-[#B59E75] text-xs font-semibold tracking-[0.2em] uppercase">Curated</span>
            <span class="text-[#B59E75] text-[10px]">◆</span>
          </div>
          <h2 class="text-4xl md:text-5xl lg:text-6xl font-serif italic text-[#1C1612] mb-6">
            Exclusive Rooms
          </h2>
        </div>
      </div>

      <div v-if="loading" class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-8 md:gap-6 lg:gap-8" :class="{ 'reveal-visible': sectionVisible }">
        <div
          v-for="(room, index) in rooms"
          :key="room.id"
          class="reveal-item"
          :class="`reveal-delay-${Math.min(index + 2, 5)}`"
        >
          <RoomCard :room="room" />
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.font-serif {
  font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;
}

/* Scroll reveal */
.reveal-item {
  opacity: 0;
  transform: translateY(24px);
  transition: opacity 0.7s cubic-bezier(0.16, 1, 0.3, 1),
              transform 0.7s cubic-bezier(0.16, 1, 0.3, 1);
}

.reveal-visible .reveal-item {
  opacity: 1;
  transform: translateY(0);
}

.reveal-delay-0 { transition-delay: 0s; }
.reveal-delay-1 { transition-delay: 0.1s; }
.reveal-delay-2 { transition-delay: 0.2s; }
.reveal-delay-3 { transition-delay: 0.3s; }
.reveal-delay-4 { transition-delay: 0.4s; }
.reveal-delay-5 { transition-delay: 0.5s; }
</style>
