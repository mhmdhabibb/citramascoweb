<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import type { Room } from '@/types'

const props = defineProps<{
  room: Room
}>()

const router = useRouter()
const goToDetail = () => {
  router.push({ name: 'room-detail', params: { id: props.room.id } })
}
</script>

<template>
  <div class="room-card group cursor-pointer" @click="goToDetail">
    <div class="room-card-image-wrapper">
      <img :src="room.image" :alt="room.name" class="room-card-image" />
    </div>
    
    <div class="pt-4 pb-2">
      <div class="flex justify-between items-start mb-3">
        <h3 class="text-[22px] font-serif text-[#1C1612] leading-tight">{{ room.name }}</h3>
        <span class="text-[#4A4A4A] text-[14px] font-medium mt-1">Rp {{ room.price.toLocaleString('id-ID') }} / night</span>
      </div>
      
      <div class="flex items-center gap-6 text-[#6B5E52] text-[13px]">
        <div class="flex items-center gap-2">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
            <circle cx="12" cy="7" r="4"></circle>
          </svg>
          <span>{{ room.capacity }} adults</span>
        </div>
        <div class="flex items-center gap-2">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="15 3 21 3 21 9"></polyline>
            <polyline points="9 21 3 21 3 15"></polyline>
            <line x1="21" y1="3" x2="14" y2="10"></line>
            <line x1="3" y1="21" x2="10" y2="14"></line>
          </svg>
          <span>{{ room.size }} Sq.m</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.room-card {
  background: transparent;
}

.font-serif {
  font-family: 'Playfair Display', ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
}

.room-card-image-wrapper {
  position: relative;
  overflow: hidden;
  border-radius: 0; /* Match mockup straight edges */
  aspect-ratio: 4/3; /* Adjust aspect ratio to fit the wide cards */
}

.room-card-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.7s cubic-bezier(0.16, 1, 0.3, 1);
}

.room-card:hover .room-card-image {
  transform: scale(1.03);
}
</style>
