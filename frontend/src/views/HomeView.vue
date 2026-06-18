<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { roomService, type Room } from '@/services/roomService'
import { useCategoryStore } from '@/stores/categoryStore'
import { useTypeStore } from '@/stores/typeStore'
import RoomCard from '@/components/room/RoomCard.vue'

const featuredRooms = ref<Room[]>([])
const loading = ref(true)

const categoryStore = useCategoryStore()
const typeStore = useTypeStore()

onMounted(async () => {
  try {
    const [rooms] = await Promise.all([
      roomService.getRooms(),
      categoryStore.fetchCategories(),
      typeStore.fetchTypes(),
    ])
    featuredRooms.value = rooms.slice(0, 3)
  } catch (error) {
    console.error('Error fetching data:', error)
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

    <!-- About Us -->
    <section id="about" class="about-section">
      <div class="about-container">
        <span class="about-badge">About Us</span>
        <h3 class="about-text">
          <strong>Hotel Citramas Co-Living</strong>, established in 2026, is a modern accommodation located in the heart of Batam. Strategically situated near major shopping centers, business areas, and popular city attractions, we offer convenience and comfort for both business and leisure travelers.
        </h3>
        <p class="about-text-sub">
          With our co-living concept, guests can enjoy a welcoming environment designed for relaxation, productivity, and a true sense of home.
        </p>
      </div>
    </section>

    <!-- Categories Section (from Backend API) -->
    <section class="categories-section">
      <div class="categories-container">
        <div class="categories-header">
          <span class="categories-badge">Categories</span>
          <h2 class="categories-title">Room Categories</h2>
          <p class="categories-subtitle">Browse our available room categories to find your perfect match.</p>
        </div>

        <div v-if="categoryStore.loading" class="categories-loading">
          <div class="spinner"></div>
        </div>

        <div v-else-if="categoryStore.error" class="categories-error">
          <p>{{ categoryStore.error }}</p>
        </div>

        <div v-else-if="categoryStore.categories.length > 0" class="categories-grid">
          <div
            v-for="category in categoryStore.categories"
            :key="category.id"
            class="category-card"
          >
            <div class="category-icon">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
                <polyline points="9 22 9 12 15 12 15 22" />
              </svg>
            </div>
            <h3 class="category-name">{{ category.name }}</h3>
            <span class="category-slug">{{ category.slug }}</span>
          </div>
        </div>

        <div v-else class="categories-empty">
          <p>No categories available yet.</p>
        </div>
      </div>
    </section>

    <!-- Room Types Section (from Backend API) -->
    <section v-if="typeStore.types.length > 0" class="types-section">
      <div class="types-container">
        <div class="types-header">
          <span class="types-badge">Room Types</span>
          <h2 class="types-title">Available Room Types</h2>
          <p class="types-subtitle">Choose from our range of room types tailored to your needs.</p>
        </div>

        <div v-if="typeStore.loading" class="categories-loading">
          <div class="spinner"></div>
        </div>

        <div v-else class="types-grid">
          <div
            v-for="roomType in typeStore.types"
            :key="roomType.id"
            class="type-card"
          >
            <div class="type-icon">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3" y="3" width="7" height="7" />
                <rect x="14" y="3" width="7" height="7" />
                <rect x="14" y="14" width="7" height="7" />
                <rect x="3" y="14" width="7" height="7" />
              </svg>
            </div>
            <span class="type-name">{{ roomType.name }}</span>
          </div>
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

/* About Us Section */
.about-section {
  background-color: #F8F5F1;
  padding: 5rem 2rem;
}

.about-container {
  max-width: 800px;
  margin: 0 auto;
  text-align: center;
}

.about-badge {
  display: inline-block;
  padding: 0.4rem 1.25rem;
  border-radius: 100px;
  font-size: 14px;
  font-weight: 500;
  color: #3D3329;
  background-color: rgba(28, 22, 18, 0.06);
  border: 1px solid rgba(28, 22, 18, 0.1);
  margin-bottom: 1.75rem;
  letter-spacing: 0.02em;
}

.about-text {
  font-size: clamp(20px, 3vw, 28px);
  font-weight: 400;
  line-height: 1.5;
  color: #1C1612;
  margin: 0 0 1.25rem;
}

.about-text strong {
  font-weight: 700;
}

.about-text-sub {
  font-size: clamp(16px, 2.5vw, 20px);
  line-height: 1.6;
  color: #6B5E52;
  margin: 0;
}

/* Categories Section */
.categories-section {
  padding: 5rem 2rem;
  background-color: #FFFFFF;
}

.categories-container {
  max-width: 1280px;
  margin: 0 auto;
}

.categories-header {
  text-align: center;
  margin-bottom: 3rem;
}

.categories-badge,
.types-badge {
  display: inline-block;
  padding: 0.4rem 1.25rem;
  border-radius: 100px;
  font-size: 14px;
  font-weight: 500;
  color: #3D3329;
  background-color: rgba(28, 22, 18, 0.06);
  border: 1px solid rgba(28, 22, 18, 0.1);
  margin-bottom: 1rem;
  letter-spacing: 0.02em;
}

.categories-title,
.types-title {
  font-size: clamp(24px, 4vw, 36px);
  font-weight: 700;
  color: #1C1612;
  margin: 0 0 0.75rem;
  letter-spacing: -0.01em;
}

.categories-subtitle,
.types-subtitle {
  font-size: 16px;
  color: #6B5E52;
  margin: 0;
}

.categories-loading {
  display: flex;
  justify-content: center;
  padding: 3rem 0;
}

.spinner {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 3px solid rgba(28, 22, 18, 0.1);
  border-top-color: #3D3329;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.categories-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 1.5rem;
}

.category-card {
  background: #FAFAF8;
  border: 1px solid rgba(28, 22, 18, 0.08);
  border-radius: 16px;
  padding: 2rem 1.5rem;
  text-align: center;
  transition: transform 0.3s ease, box-shadow 0.3s ease, border-color 0.3s ease;
  cursor: default;
}

.category-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(28, 22, 18, 0.08);
  border-color: rgba(28, 22, 18, 0.15);
}

.category-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  border-radius: 14px;
  background: linear-gradient(135deg, #F8F5F1, #EDE8E1);
  color: #3D3329;
  margin-bottom: 1rem;
}

.category-name {
  font-size: 18px;
  font-weight: 600;
  color: #1C1612;
  margin: 0 0 0.25rem;
}

.category-slug {
  font-size: 13px;
  color: #9B8E82;
  letter-spacing: 0.02em;
}

.categories-error {
  text-align: center;
  padding: 2rem;
  color: #B44A3F;
  font-size: 14px;
}

.categories-empty {
  text-align: center;
  padding: 3rem 0;
  color: #9B8E82;
  font-size: 15px;
}

/* Types Section */
.types-section {
  padding: 4rem 2rem 5rem;
  background-color: #F8F5F1;
}

.types-container {
  max-width: 1280px;
  margin: 0 auto;
}

.types-header {
  text-align: center;
  margin-bottom: 2.5rem;
}

.types-grid {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 1rem;
}

.type-card {
  display: inline-flex;
  align-items: center;
  gap: 0.625rem;
  padding: 0.75rem 1.5rem;
  background: #FFFFFF;
  border: 1px solid rgba(28, 22, 18, 0.08);
  border-radius: 100px;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  cursor: default;
}

.type-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(28, 22, 18, 0.06);
}

.type-icon {
  color: #6B5E52;
  display: flex;
}

.type-name {
  font-size: 15px;
  font-weight: 500;
  color: #1C1612;
}
</style>
