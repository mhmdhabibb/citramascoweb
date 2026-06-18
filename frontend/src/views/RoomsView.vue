<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { roomService, type Room } from '@/services/roomService'
import { useCategoryStore } from '@/stores/categoryStore'
import { useTypeStore } from '@/stores/typeStore'
import RoomCard from '@/components/room/RoomCard.vue'

const rooms = ref<Room[]>([])
const loading = ref(true)

const categoryStore = useCategoryStore()
const typeStore = useTypeStore()

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
  <div class="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
    <div class="mb-10 text-center md:text-left">
      <h1 class="text-3xl font-extrabold tracking-tight">Our Rooms & Suites</h1>
      <p class="mt-2 text-muted-foreground">Find the perfect room for your trip, business stay, or weekend getaway.</p>
    </div>

    <!-- Categories Filter (from Backend API) -->
    <div v-if="categoryStore.categories.length > 0" class="categories-filter">
      <div class="filter-label">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3" />
        </svg>
        <span>Categories</span>
      </div>
      <div class="filter-chips">
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
    <div v-if="typeStore.types.length > 0" class="types-bar">
      <div class="filter-label">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="3" y="3" width="7" height="7" />
          <rect x="14" y="3" width="7" height="7" />
          <rect x="14" y="14" width="7" height="7" />
          <rect x="3" y="14" width="7" height="7" />
        </svg>
        <span>Room Types</span>
      </div>
      <div class="filter-chips">
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

    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-8">
      <RoomCard v-for="room in rooms" :key="room.id" :room="room" />
    </div>
  </div>
</template>

<style scoped>
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
  transition: all 0.2s ease;
}

.filter-chip:hover {
  background: #EDE8E1;
  border-color: rgba(28, 22, 18, 0.2);
  transform: translateY(-1px);
}

.filter-chip--type {
  background: #FFFFFF;
  border-color: rgba(28, 22, 18, 0.08);
}

.filter-chip--type:hover {
  background: #F8F5F1;
}
</style>
