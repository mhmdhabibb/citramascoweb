<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { roomService, type Room } from '@/services/roomService'

const route = useRoute()
const room = ref<Room | null>(null)
const loading = ref(true)
const imageLoaded = ref(false)

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
  <div class="max-w-[1400px] mx-auto py-12 px-4 sm:px-6 lg:px-8">
    <div v-if="loading" class="flex justify-center items-center py-20">
      <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-primary"></div>
    </div>

    <div v-else-if="room" class="grid grid-cols-1 lg:grid-cols-2 gap-12">
      <!-- Image with fade-in -->
      <div class="detail-image-wrapper" :class="{ 'detail-image-loaded': imageLoaded }">
        <img
          :src="room.image"
          :alt="room.name"
          class="detail-image"
          @load="imageLoaded = true"
        />
      </div>

      <!-- Info with staggered entrance -->
      <div class="detail-info">
        <h1 class="text-3xl font-extrabold tracking-tight detail-animate detail-delay-0">{{ room.name }}</h1>
        <p class="text-2xl font-bold text-primary mt-2 detail-animate detail-delay-1">
          Rp {{ room.price.toLocaleString('id-ID') }}
          <span class="text-sm font-normal text-muted-foreground">/ night</span>
        </p>
        
        <div class="border-t border-border mt-6 pt-6 detail-animate detail-delay-2">
          <h2 class="text-lg font-semibold mb-2">Description</h2>
          <p class="text-muted-foreground">{{ room.description }}</p>
        </div>

        <div class="border-t border-border mt-6 pt-6 detail-animate detail-delay-3">
          <h2 class="text-lg font-semibold mb-2">Features & Amenities</h2>
          <ul class="grid grid-cols-2 gap-2">
            <li
              v-for="(feature, index) in room.features"
              :key="feature"
              class="feature-item"
              :style="{ transitionDelay: `${0.4 + index * 0.06}s` }"
            >
              <svg class="h-4 w-4 text-green-500 mr-2 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              {{ feature }}
            </li>
          </ul>
        </div>

        <div class="mt-8 flex gap-4 detail-animate detail-delay-4">
          <router-link
            :to="{ name: 'booking', query: { roomId: room.id } }"
            class="detail-btn-primary"
          >
            Book This Room
          </router-link>
          <router-link to="/rooms" class="detail-btn-secondary">
            Back to Rooms
          </router-link>
        </div>
      </div>
    </div>

    <div v-else class="text-center py-20 not-found-animate">
      <h2 class="text-2xl font-bold">Room Not Found</h2>
      <p class="text-muted-foreground mt-2">The room you are looking for does not exist.</p>
      <router-link to="/rooms" class="mt-4 inline-block bg-primary text-primary-foreground px-4 py-2 rounded-md">
        Back to Rooms
      </router-link>
    </div>
  </div>
</template>

<style scoped>
/* Image entrance */
.detail-image-wrapper {
  overflow: hidden;
  border-radius: 12px;
  box-shadow: 0 4px 24px rgba(28, 22, 18, 0.08);
  height: 100%;
}

.detail-image {
  width: 100%;
  height: 100%;
  min-height: 24rem;
  object-fit: cover;
  opacity: 0;
  transform: scale(1.04);
  transition: opacity 0.8s ease, transform 1s cubic-bezier(0.16, 1, 0.3, 1);
}

.detail-image-loaded .detail-image {
  opacity: 1;
  transform: scale(1);
}

/* Info staggered entrance */
.detail-animate {
  opacity: 0;
  transform: translateY(20px);
  animation: detailSlideUp 0.7s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.detail-delay-0 { animation-delay: 0.1s; }
.detail-delay-1 { animation-delay: 0.2s; }
.detail-delay-2 { animation-delay: 0.3s; }
.detail-delay-3 { animation-delay: 0.4s; }
.detail-delay-4 { animation-delay: 0.55s; }

@keyframes detailSlideUp {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Feature items stagger */
.feature-item {
  display: flex;
  align-items: center;
  font-size: 14px;
  color: var(--muted-foreground);
  opacity: 0;
  transform: translateX(-8px);
  animation: featureSlideIn 0.5s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  animation-delay: inherit;
}

@keyframes featureSlideIn {
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

/* Buttons */
.detail-btn-primary {
  flex-grow: 1;
  text-align: center;
  border-radius: 8px;
  background: var(--primary);
  color: var(--primary-foreground);
  padding: 0.75rem 1.5rem;
  font-size: 16px;
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(28, 22, 18, 0.1);
  transition: all 0.35s cubic-bezier(0.16, 1, 0.3, 1);
  text-decoration: none;
}

.detail-btn-primary:hover {
  opacity: 0.9;
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(28, 22, 18, 0.15);
}

.detail-btn-secondary {
  text-align: center;
  border-radius: 8px;
  border: 1px solid var(--border);
  color: var(--foreground);
  padding: 0.75rem 1.5rem;
  font-size: 16px;
  font-weight: 600;
  transition: all 0.35s cubic-bezier(0.16, 1, 0.3, 1);
  text-decoration: none;
}

.detail-btn-secondary:hover {
  background: var(--muted);
  transform: translateY(-2px);
}

@media (min-width: 640px) {
  .detail-btn-primary {
    flex-grow: 0;
  }
}

/* Not found animation */
.not-found-animate {
  animation: detailSlideUp 0.6s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}
</style>
