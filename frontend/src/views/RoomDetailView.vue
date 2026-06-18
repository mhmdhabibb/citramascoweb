<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { roomService, type Room } from '@/services/roomService'

const route = useRoute()
const room = ref<Room | null>(null)
const loading = ref(true)

onMounted(async () => {
  const roomId = Number(route.params.id)
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
})
</script>

<template>
  <div class="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
    <div v-if="loading" class="flex justify-center items-center py-20">
      <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-primary"></div>
    </div>

    <div v-else-if="room" class="grid grid-cols-1 lg:grid-cols-2 gap-12">
      <div>
        <img :src="room.image" :alt="room.name" class="w-full h-96 object-cover rounded-lg shadow-sm" />
      </div>
      <div>
        <h1 class="text-3xl font-extrabold tracking-tight">{{ room.name }}</h1>
        <p class="text-2xl font-bold text-primary mt-2">Rp {{ room.price.toLocaleString('id-ID') }} <span class="text-sm font-normal text-muted-foreground">/ night</span></p>
        
        <div class="border-t border-border mt-6 pt-6">
          <h2 class="text-lg font-semibold mb-2">Description</h2>
          <p class="text-muted-foreground">{{ room.description }}</p>
        </div>

        <div class="border-t border-border mt-6 pt-6">
          <h2 class="text-lg font-semibold mb-2">Features & Amenities</h2>
          <ul class="grid grid-cols-2 gap-2">
            <li v-for="feature in room.features" :key="feature" class="flex items-center text-sm text-muted-foreground">
              <svg class="h-4 w-4 text-green-500 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              {{ feature }}
            </li>
          </ul>
        </div>

        <div class="mt-8 flex gap-4">
          <router-link :to="{ name: 'booking', query: { roomId: room.id } }" class="flex-grow sm:flex-none text-center rounded-md bg-primary px-6 py-3 text-base font-semibold text-primary-foreground shadow-sm hover:bg-primary/90">
            Book This Room
          </router-link>
          <router-link to="/rooms" class="text-center rounded-md border border-border px-6 py-3 text-base font-semibold text-foreground hover:bg-muted">
            Back to Rooms
          </router-link>
        </div>
      </div>
    </div>

    <div v-else class="text-center py-20">
      <h2 class="text-2xl font-bold">Room Not Found</h2>
      <p class="text-muted-foreground mt-2">The room you are looking for does not exist.</p>
      <router-link to="/rooms" class="mt-4 inline-block bg-primary text-primary-foreground px-4 py-2 rounded-md">
        Back to Rooms
      </router-link>
    </div>
  </div>
</template>
