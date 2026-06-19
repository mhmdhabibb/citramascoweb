   <script setup lang="ts">
import { ref, onMounted } from 'vue'
import { roomService, type Room } from '@/services/roomService'
import { useCategoryStore } from '@/stores/categoryStore'
import { useTypeStore } from '@/stores/typeStore'
import { useScrollReveal } from '@/composables/useScrollReveal'
import RoomCard from '@/components/room/RoomCard.vue'

const rooms = ref<Room[]>([])
const loading = ref(true)

const categoryStore = useCategoryStore()
const typeStore = useTypeStore()

const { elementRef: headerRef, isVisible: headerVisible } = useScrollReveal(0.2)
const { elementRef: gridRef, isVisible: gridVisible } = useScrollReveal(0.05)

onMounted(async () => {
  try {
    const [fetchedRooms] = await Promise.all([
      roomService.getRooms(),
      categoryStore.fetchCategories(),
      typeStore.fetchTypes(),
    ])
    rooms.value = fetchedRooms
  } catch (error) {
    console.error('Error fetching rooms:', error)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="rooms-page bg-[#FAF7F2] min-h-screen pt-24 pb-24">
    <div class="max-w-[1200px] mx-auto px-4 sm:px-6 lg:px-8">
      
      <!-- Header -->
      <div ref="headerRef" class="mb-16 text-center max-w-2xl mx-auto" :class="{ 'reveal-visible': headerVisible }">
        <div class="reveal-item reveal-delay-0">
          <div class="flex items-center justify-center gap-2 mb-4">
            <span class="text-[#B59E75] text-[10px]">◆</span>
            <span class="text-[#B59E75] text-xs font-semibold tracking-[0.2em] uppercase">Our Collection</span>
            <span class="text-[#B59E75] text-[10px]">◆</span>
          </div>
          <h1 class="text-4xl md:text-5xl lg:text-6xl text-[#1A1A1A] mb-6" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif; font-style: italic;">
            Rooms & Suites
          </h1>
          <p class="text-[#6B5E52] text-[15px] leading-relaxed max-w-lg mx-auto">
            Find the perfect room for your trip, business stay, or weekend getaway.
          </p>
        </div>
      </div>

      <!-- Categories Filter (from Backend API) -->
      <div v-if="categoryStore.categories.length > 0" class="categories-filter justify-center mb-6">
        <div class="filter-label hidden sm:flex">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3" />
          </svg>
          <span>Categories</span>
        </div>
        <div class="filter-chips justify-center">
          <span
            v-for="category in categoryStore.categories"
            :key="category.id"
            class="filter-chip"
          >
            {{ category.name }}
          </span>
        </div>
      </div>

      <!-- Room Types (from Backend API) -->
      <div v-if="typeStore.types.length > 0" class="types-bar justify-center mb-12">
        <div class="filter-label hidden sm:flex">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="3" width="7" height="7" />
            <rect x="14" y="3" width="7" height="7" />
            <rect x="14" y="14" width="7" height="7" />
            <rect x="3" y="14" width="7" height="7" />
          </svg>
          <span>Room Types</span>
        </div>
        <div class="filter-chips justify-center">
          <span
            v-for="roomType in typeStore.types"
            :key="roomType.id"
            class="filter-chip filter-chip--type"
          >
            {{ roomType.name }}
          </span>
        </div>
      </div>

      <!-- Room Cards -->
      <div v-if="loading" class="flex justify-center items-center py-20">
        <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-primary"></div>
      </div>

      <div v-else ref="gridRef" class="grid grid-cols-1 md:grid-cols-3 gap-8 md:gap-6 lg:gap-8" :class="{ 'reveal-visible': gridVisible }">
        <div
          v-for="(room, index) in rooms"
          :key="room.id"
          class="reveal-item"
          :class="`reveal-delay-${Math.min(index, 5)}`"
        >
          <RoomCard :room="room" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
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

/* Filters */
.categories-filter,
.types-bar {
  margin-bottom: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.filter-label {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  font-size: 13px;
  font-weight: 600;
  color: #6B5E52;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  flex-shrink: 0;
}

.filter-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.filter-chip {
  display: inline-flex;
  align-items: center;
  padding: 0.4rem 1rem;
  font-size: 13px;
  font-weight: 500;
  color: #3D3329;
  background: #F8F5F1;
  border: 1px solid rgba(28, 22, 18, 0.1);
  border-radius: 100px;
  cursor: pointer;
  transition: all 0.35s cubic-bezier(0.16, 1, 0.3, 1);
}

.filter-chip:hover {
  background: #EDE8E1;
  border-color: rgba(28, 22, 18, 0.2);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(28, 22, 18, 0.06);
}

.filter-chip--type {
  background: #FFFFFF;
  border-color: rgba(28, 22, 18, 0.08);
}

.filter-chip--type:hover {
  background: #F8F5F1;
}
</style>
