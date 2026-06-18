<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { roomService, type Room } from '@/services/roomService'
import { useScrollReveal } from '@/composables/useScrollReveal'
import RoomCard from '@/components/room/RoomCard.vue'

const featuredRooms = ref<Room[]>([])
const loading = ref(true)

const { elementRef: aboutRef, isVisible: aboutVisible } = useScrollReveal(0.2)
const { elementRef: featuredRef, isVisible: featuredVisible } = useScrollReveal(0.1)

onMounted(async () => {
  try {
    const rooms = await roomService.getRooms()
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
        <h1 class="hero-heading animate-hero-text">
          Experience luxury<br />
          living at Citramas<br />
          Co Living.
        </h1>
        <p class="hero-subtitle animate-hero-text animate-hero-text--delay-1">
          Your premium co-living destination offering beautifully designed rooms,<br class="hidden-mobile" />
          modern amenities, and a vibrant community in the heart of the city.
        </p>
        <div class="hero-buttons animate-hero-text animate-hero-text--delay-2">
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
    <section id="about" ref="aboutRef" class="max-w-[1400px] mx-auto px-4 sm:px-6 lg:px-8 py-16">
      <div class="grid grid-cols-1 md:grid-cols-2" :class="{ 'reveal-visible': aboutVisible }">
        <!-- Image Column -->
        <div class="reveal-item reveal-delay-0 overflow-hidden">
          <img src="@/assets/hero-bg.png" alt="Hotel Details" class="w-full h-full object-cover min-h-[400px] md:min-h-[500px]" />
        </div>
        <!-- Text Column -->
        <div class="bg-[#F8F9FA] flex flex-col justify-center p-8 md:p-16 lg:p-24 reveal-item reveal-delay-1">
          <div class="text-[11px] font-bold tracking-[0.15em] text-gray-500 uppercase mb-3">WHO WE ARE</div>
          <h2 class="text-4xl md:text-[42px] font-bold tracking-tight text-[#1C1612] mb-6 leading-[1.15]">The details that matter.</h2>
          <p class="text-[15px] text-[#4A4A4A] leading-relaxed mb-6">
            The Hotel opened in 2026 with a simple idea: that a great hotel should feel like somewhere, not anywhere. Independently owned and operated, we have never chased scale over quality.
          </p>
          <p class="text-[15px] text-[#4A4A4A] leading-relaxed">
            Every decision we make, from the mattresses to the menu, comes back to the same question: does this make the stay better? If the answer is no, we do not do it.
          </p>
        </div>
      </div>
    </section>

    <!-- Featured Rooms -->
    <section ref="featuredRef" class="max-w-[1400px] mx-auto py-16 px-4 sm:px-6 lg:px-8">
      <div class="flex flex-col sm:flex-row sm:items-end justify-between gap-6 mb-12" :class="{ 'reveal-visible': featuredVisible }">
        <div class="reveal-item reveal-delay-0">
          <div class="text-[11px] font-bold tracking-[0.15em] text-gray-500 uppercase mb-3">ROOMS & SUITES</div>
          <h2 class="text-4xl md:text-5xl font-bold tracking-tight text-[#1C1612]">Find your room.</h2>
        </div>
        <div class="reveal-item reveal-delay-1">
          <router-link to="/rooms" class="inline-block border border-gray-300 text-[13px] font-medium px-5 py-2.5 hover:bg-gray-50 transition-colors text-[#1C1612]">
            View All Rooms
          </router-link>
        </div>
      </div>

      <div v-if="loading" class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-8" :class="{ 'reveal-visible': featuredVisible }">
        <div
          v-for="(room, index) in featuredRooms"
          :key="room.id"
          class="reveal-item"
          :class="`reveal-delay-${index + 2}`"
        >
          <RoomCard :room="room" />
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
/* ============================================
   Hero Section
   ============================================ */
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
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}


/* Hero entrance animations */
.animate-hero-text {
  opacity: 0;
  transform: translateY(30px);
  animation: heroSlideUp 0.8s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.animate-hero-text--delay-1 {
  animation-delay: 0.2s;
}

.animate-hero-text--delay-2 {
  animation-delay: 0.4s;
}

@keyframes heroSlideUp {
  to {
    opacity: 1;
    transform: translateY(0);
  }
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
  transition: all 0.35s cubic-bezier(0.16, 1, 0.3, 1);
}

.hero-btn-primary:hover {
  background-color: #FFFFFF;
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(248, 245, 241, 0.3);
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
  transition: all 0.35s cubic-bezier(0.16, 1, 0.3, 1);
}

.hero-btn-secondary:hover {
  background-color: rgba(255, 255, 255, 0.25);
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(255, 255, 255, 0.1);
}

/* ============================================
   Scroll Reveal System
   ============================================ */
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



/* ============================================
   Categories Section
   ============================================ */
.categories-section {
  padding: 5rem 2rem;
  background-color: #FFFFFF;
}

.categories-container {
  max-width: 1400px;
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
  transition: transform 0.4s cubic-bezier(0.16, 1, 0.3, 1),
              box-shadow 0.4s cubic-bezier(0.16, 1, 0.3, 1),
              border-color 0.3s ease;
  cursor: default;
}

.category-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 16px 40px rgba(28, 22, 18, 0.08);
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
  transition: transform 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

.category-card:hover .category-icon {
  transform: scale(1.1) rotate(-3deg);
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

/* ============================================
   Types Section
   ============================================ */
.types-section {
  padding: 4rem 2rem 5rem;
  background-color: #F8F5F1;
}

.types-container {
  max-width: 1400px;
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
  transition: all 0.35s cubic-bezier(0.16, 1, 0.3, 1);
  cursor: default;
}

.type-card:hover {
  transform: translateY(-3px) scale(1.03);
  box-shadow: 0 8px 20px rgba(28, 22, 18, 0.06);
  border-color: rgba(28, 22, 18, 0.15);
}

.type-icon {
  color: #6B5E52;
  display: flex;
  transition: color 0.3s ease;
}

.type-card:hover .type-icon {
  color: #3D3329;
}

.type-name {
  font-size: 15px;
  font-weight: 500;
  color: #1C1612;
}
</style>
