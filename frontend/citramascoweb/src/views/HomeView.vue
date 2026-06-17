<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { roomService, type Room } from '@/services/roomService'
import RoomCard from '@/components/room/RoomCard.vue'

const featuredRooms = ref<Room[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const rooms = await roomService.getRooms()
    featuredRooms.value = rooms.slice(0, 3)
  } catch (error) {
    console.error('Error fetching rooms:', error)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <!-- Hero Section -->
    <section class="hero-section">
      <div class="hero-overlay"></div>
      <div class="hero-container">
        <h1 class="hero-heading">
          Experience luxury<br />
          living at Citramas<br />
          Co Living.
        </h1>
        <p class="hero-subtitle">
          Your premium co-living destination offering beautifully designed rooms,<br class="hidden-mobile" />
          modern amenities, and a vibrant community in the heart of the city.
        </p>
        <div class="hero-buttons">
          <router-link to="/rooms" class="hero-btn-primary">
            View Rooms &nbsp;&rarr;
          </router-link>
          <router-link to="/booking" class="hero-btn-secondary">
            Book Now
          </router-link>
        </div>
      </div>
    </section>

    <!-- Featured Rooms -->
    <section class="max-w-7xl mx-auto py-16 px-4 sm:px-6 lg:px-8">
      <div class="text-center max-w-xl mx-auto mb-12">
        <h2 class="text-3xl font-extrabold tracking-tight">Our Featured Rooms</h2>
        <p class="mt-4 text-muted-foreground">Select from our range of beautifully designed rooms equipped with premium amenities.</p>
      </div>

      <div v-if="loading" class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <RoomCard v-for="room in featuredRooms" :key="room.id" :room="room" />
      </div>
    </section>
  </div>
</template>

<style scoped>
.hero-section {
  position: relative;
  background-image: url('@/assets/hero-bg.png');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  padding: 8rem 2rem 6rem;
  min-height: 85vh;
  display: flex;
  align-items: flex-end;
}

.hero-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    to top,
    rgba(28, 22, 18, 0.85) 0%,
    rgba(28, 22, 18, 0.5) 50%,
    rgba(28, 22, 18, 0.2) 100%
  );
}

.hero-container {
  position: relative;
  z-index: 1;
  max-width: 1280px;
  margin: 0 auto;
  width: 100%;
}

.hero-heading {
  font-size: clamp(40px, 6vw, 64px);
  font-weight: 700;
  line-height: 1.1;
  letter-spacing: -0.02em;
  color: #FFFFFF;
  margin: 0 0 1.5rem;
}

.hero-subtitle {
  font-size: 18px;
  line-height: 1.6;
  color: rgba(255, 255, 255, 0.75);
  margin: 0 0 2.5rem;
  max-width: 600px;
}

.hidden-mobile {
  display: none;
}

@media (min-width: 768px) {
  .hidden-mobile {
    display: inline;
  }
}

.hero-buttons {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.hero-btn-primary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.75rem 1.75rem;
  border-radius: 100px;
  font-size: 16px;
  font-weight: 500;
  color: #1C1612;
  background-color: #F8F5F1;
  text-decoration: none;
  transition: background-color 0.3s ease, transform 0.2s ease;
}

.hero-btn-primary:hover {
  background-color: #FFFFFF;
  transform: translateY(-1px);
}

.hero-btn-secondary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.75rem 1.75rem;
  border-radius: 100px;
  font-size: 16px;
  font-weight: 500;
  color: #FFFFFF;
  background-color: rgba(255, 255, 255, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.3);
  text-decoration: none;
  backdrop-filter: blur(4px);
  transition: background-color 0.3s ease, transform 0.2s ease;
}

.hero-btn-secondary:hover {
  background-color: rgba(255, 255, 255, 0.25);
  transform: translateY(-1px);
}
</style>
