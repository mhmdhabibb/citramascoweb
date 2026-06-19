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
  <div class="room-card group cursor-pointer">
    <div class="room-card-image-wrapper">
      <img :src="room.image" :alt="room.name" class="room-card-image" />
      <div class="room-card-image-overlay"></div>
    </div>
    
    <div class="pt-5 px-1">
      <div class="text-accent text-[11px] font-bold tracking-[0.15em] uppercase mb-2">{{ category }}</div>
      <h3 class="text-[26px] font-serif text-primary-dark mb-3 leading-tight">{{ room.name }}</h3>
      <p class="text-[#6B5E52] text-[14.5px] leading-relaxed mb-6">
        {{ room.description }}
      </p>
      
      <div class="flex items-center justify-between">
        <span class="text-[#4A4A4A] text-[14.5px]">from Rp {{ room.price.toLocaleString('id-ID') }}</span>
        <span class="text-[#1C1612] text-[14.5px] border-b border-[#1C1612] pb-0.5 font-medium transition-colors hover:text-accent hover:border-accent">Inquire</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.room-card {
  background: transparent;
}

.text-accent {
  color: #B59E75;
}

.text-primary-dark {
  color: #1A1A1A;
}

.font-serif {
  font-family: 'Playfair Display', ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
}

.room-card-image-wrapper {
  position: relative;
  overflow: hidden;
  border-radius: 16px;
  aspect-ratio: 4/5;
  margin-bottom: 1.25rem;
}

.room-card-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.7s cubic-bezier(0.16, 1, 0.3, 1);
}

.room-card:hover .room-card-image {
  transform: scale(1.05);
}

.room-card-image-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.1);
  opacity: 0;
  transition: opacity 0.5s ease;
}

.room-card:hover .room-card-image-overlay {
  opacity: 1;
}
</style>
