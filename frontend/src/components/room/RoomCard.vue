<script setup lang="ts">
import { computed } from 'vue'
import type { Room } from '@/types'

const props = defineProps<{
  room: Room
}>()

// Extract category from room name (e.g. "Standard Room" -> "STANDARD")
const category = computed(() => {
  return props.room.name.split(' ')[0].toUpperCase()
})
</script>

<template>
  <div class="room-card">
    <div class="room-card-image-wrapper">
      <img :src="room.image" :alt="room.name" class="room-card-image" />
      <div class="room-card-image-overlay"></div>
    </div>
    <div class="room-info">
      <div class="text-[11px] font-bold tracking-[0.1em] text-gray-500 uppercase mb-2">{{ category }}</div>
      <h3 class="text-xl font-bold text-[#1C1612] leading-tight mb-3">{{ room.name }}</h3>
      <p class="text-[15px] text-gray-500">From Rp {{ room.price.toLocaleString('id-ID') }} <span class="text-[13px]">/ night</span></p>
    </div>
  </div>
</template>

<style scoped>
.room-card {
  background: transparent;
  overflow: hidden;
  position: relative;
  cursor: pointer;
}

.room-info {
  background-color: #F8F9FA;
  padding: 1.5rem;
  transition: background-color 0.3s ease;
}

.room-card:hover .room-info {
  background-color: #F1F3F5;
}

.room-card-image-wrapper {
  position: relative;
  overflow: hidden;
}

.room-card-image {
  width: 100%;
  aspect-ratio: 16/10;
  height: auto;
  object-fit: cover;
  transition: transform 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

.room-card:hover .room-card-image {
  transform: scale(1.05);
}

.room-card-image-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(28, 22, 18, 0.15) 0%, transparent 50%);
  opacity: 0;
  transition: opacity 0.4s ease;
}

.room-card:hover .room-card-image-overlay {
  opacity: 1;
}

</style>
