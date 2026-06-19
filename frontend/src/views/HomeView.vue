<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { roomService, type Room } from '@/services/roomService'
import { useScrollReveal } from '@/composables/useScrollReveal'
import RoomCard from '@/components/room/RoomCard.vue'

const featuredRooms = ref<Room[]>([])
const loading = ref(true)

const { elementRef: heroRef, isVisible: heroVisible } = useScrollReveal(0.1)
const { elementRef: aboutRef, isVisible: aboutVisible } = useScrollReveal(0.2)
const { elementRef: valuesRef, isVisible: valuesVisible } = useScrollReveal(0.2)
const { elementRef: featuredRef, isVisible: featuredVisible } = useScrollReveal(0.1)
const { elementRef: ctaRef, isVisible: ctaVisible } = useScrollReveal(0.2)

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
    <section ref="heroRef" class="hero-section" :class="{ 'reveal-visible': heroVisible }">
      <div class="hero-overlay"></div>
      <div class="hero-container">
        <h1 class="text-[50px] md:text-[70px] lg:text-[89px] leading-[1.1] font-serif italic text-white mb-6 reveal-item reveal-delay-0">
          Experience luxury<br />
          living at Citramas<br />
          Co Living.
        </h1>
        <p class="hero-subtitle reveal-item reveal-delay-1">
          Your premium co-living destination offering beautifully designed rooms,<br class="hidden-mobile" />
          modern amenities, and a vibrant community in the heart of the city.
        </p>
        <div class="hero-buttons reveal-item reveal-delay-2">
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
    <section id="about" ref="aboutRef" class="w-full bg-[#FAF6F1] py-24 md:py-32">
      <div class="max-w-[1400px] mx-auto px-4 sm:px-6 lg:px-12" :class="{ 'reveal-visible': aboutVisible }">
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-24 items-center">
          
          <!-- Left Column: Heading -->
          <div class="lg:col-span-5 reveal-item reveal-delay-0">
            <h2 class="text-4xl md:text-5xl lg:text-6xl font-serif italic text-primary-dark mb-6">
              A House Built<br />
              For Slow Living
            </h2>
          </div>
          
          <!-- Right Column: Text -->
          <div class="lg:col-span-7 reveal-item reveal-delay-1">
            <p class="text-[16px] md:text-[18px] text-[#6B5E52] leading-relaxed mb-6 font-light">
              CitraMas began as a private residence on the tranquil coastline. 
              When we opened our doors to guests, we kept its quiet — its hand-plastered walls, its sun-drenched terraces, its old citrus garden.
            </p>
            <p class="text-[16px] md:text-[18px] text-[#6B5E52] leading-relaxed font-light">
              Today it is a small resort of curated suites and pavilions, run 
              by a small team who care deeply about the craft of hospitality.
            </p>
          </div>
          
        </div>
      </div>
    </section>

    <!-- Our Values -->
    <section ref="valuesRef" class="w-full bg-[#F2EFEA] py-24 md:py-32">
      <div class="max-w-[1200px] mx-auto px-4 sm:px-6 lg:px-8" :class="{ 'reveal-visible': valuesVisible }">
        
        <!-- Header -->
        <div class="text-center mb-16 reveal-item reveal-delay-0">
          <div class="flex items-center justify-center gap-4 mb-4">
            <span class="text-[#B59E75] text-[10px]">◆</span>
            <span class="text-[#C2A379] text-[11px] font-bold tracking-[0.2em] uppercase">What Guides Us</span>
            <span class="text-[#B59E75] text-[10px]">◆</span>
          </div>
          <h2 class="text-4xl md:text-5xl lg:text-6xl font-serif italic text-primary-dark mb-6">
            Our Values
          </h2>
        </div>

        <!-- Values Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-12 text-center">
          
          <!-- Value 1 -->
          <div class="flex flex-col items-center reveal-item reveal-delay-1">
            <div class="w-20 h-20 mb-6 flex items-center justify-center text-[#1C1612]">
              <svg viewBox="0 0 100 100" fill="none" stroke="currentColor" stroke-width="2" class="w-full h-full">
                 <path d="M50 5C53.5 15 65 20 75 15C73 25 80 32 90 32C82 39 85 50 95 50C85 50 82 61 90 68C80 68 73 75 75 85C65 80 53.5 85 50 95C46.5 85 35 80 25 85C27 75 20 68 10 68C18 61 15 50 5 50C15 50 18 39 10 32C20 32 27 25 25 15C35 20 46.5 15 50 5Z" stroke-width="1.5" stroke-linejoin="round" />
                 <circle cx="50" cy="40" r="8" fill="currentColor" />
                 <path d="M50 52C42 52 38 60 38 68H62C62 60 58 52 50 52Z" fill="currentColor" />
              </svg>
            </div>
            <h3 class="text-[22px] font-serif text-[#1C1612] mb-3" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">Thoughtful</h3>
            <p class="text-[14px] text-[#6B5E52] leading-[1.8] font-light max-w-[220px] mx-auto">
              Every stay is designed with care and attention to detail.
            </p>
          </div>

          <!-- Value 2 -->
          <div class="flex flex-col items-center reveal-item reveal-delay-2">
            <div class="w-20 h-20 mb-6 flex items-center justify-center text-[#1C1612]">
              <svg viewBox="0 0 100 100" fill="none" stroke="currentColor" stroke-width="2" class="w-full h-full">
                 <path d="M50 5C53.5 15 65 20 75 15C73 25 80 32 90 32C82 39 85 50 95 50C85 50 82 61 90 68C80 68 73 75 75 85C65 80 53.5 85 50 95C46.5 85 35 80 25 85C27 75 20 68 10 68C18 61 15 50 5 50C15 50 18 39 10 32C20 32 27 25 25 15C35 20 46.5 15 50 5Z" stroke-width="1.5" stroke-linejoin="round" />
                 <circle cx="50" cy="50" r="12" fill="currentColor" />
                 <path d="M50 25V18M50 82V75M25 50H18M82 50H75M32 32L27 27M68 68L73 73M32 68L27 73M68 32L73 27" stroke-width="3" stroke-linecap="round" />
              </svg>
            </div>
            <h3 class="text-[22px] font-serif text-[#1C1612] mb-3" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">Warm Hospitality</h3>
            <p class="text-[14px] text-[#6B5E52] leading-[1.8] font-light max-w-[220px] mx-auto">
              Service that feels personal, genuine, and welcoming.
            </p>
          </div>

          <!-- Value 3 -->
          <div class="flex flex-col items-center reveal-item reveal-delay-3">
            <div class="w-20 h-20 mb-6 flex items-center justify-center text-[#1C1612]">
              <svg viewBox="0 0 100 100" fill="none" stroke="currentColor" stroke-width="2" class="w-full h-full">
                 <path d="M50 5C53.5 15 65 20 75 15C73 25 80 32 90 32C82 39 85 50 95 50C85 50 82 61 90 68C80 68 73 75 75 85C65 80 53.5 85 50 95C46.5 85 35 80 25 85C27 75 20 68 10 68C18 61 15 50 5 50C15 50 18 39 10 32C20 32 27 25 25 15C35 20 46.5 15 50 5Z" stroke-width="1.5" stroke-linejoin="round" />
                 <path d="M50 65C50 65 35 52 35 42C35 34 42 30 48 35C50 37 50 37 52 35C58 30 65 34 65 42C65 52 50 65 50 65Z" fill="currentColor" />
                 <path d="M42 45L50 52L58 45" stroke="white" stroke-width="2" stroke-linecap="round" />
              </svg>
            </div>
            <h3 class="text-[22px] font-serif text-[#1C1612] mb-3" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">Authentic</h3>
            <p class="text-[14px] text-[#6B5E52] leading-[1.8] font-light max-w-[220px] mx-auto">
              Local culture, flavors, and experiences woven into every stay.
            </p>
          </div>

          <!-- Value 4 -->
          <div class="flex flex-col items-center reveal-item reveal-delay-4">
            <div class="w-20 h-20 mb-6 flex items-center justify-center text-[#1C1612]">
              <svg viewBox="0 0 100 100" fill="none" stroke="currentColor" stroke-width="2" class="w-full h-full">
                 <path d="M50 5C53.5 15 65 20 75 15C73 25 80 32 90 32C82 39 85 50 95 50C85 50 82 61 90 68C80 68 73 75 75 85C65 80 53.5 85 50 95C46.5 85 35 80 25 85C27 75 20 68 10 68C18 61 15 50 5 50C15 50 18 39 10 32C20 32 27 25 25 15C35 20 46.5 15 50 5Z" stroke-width="1.5" stroke-linejoin="round" />
                 <path d="M40 30H60L52 48H48L40 66H60" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" />
                 <path d="M40 30L60 30M40 66L60 66" stroke-width="4" stroke-linecap="round" />
                 <circle cx="50" cy="48" r="3" fill="currentColor" />
              </svg>
            </div>
            <h3 class="text-[22px] font-serif text-[#1C1612] mb-3" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">Memorable</h3>
            <p class="text-[14px] text-[#6B5E52] leading-[1.8] font-light max-w-[220px] mx-auto">
              Creating moments guests will want to return to again and again.
            </p>
          </div>

        </div>
      </div>
    </section>

    <!-- Featured Rooms -->
    <div class="bg-[#FAF7F2] py-24">
      <section ref="featuredRef" class="max-w-[1200px] mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center max-w-2xl mx-auto mb-16" :class="{ 'reveal-visible': featuredVisible }">
          <div class="reveal-item reveal-delay-0">
            <div class="flex items-center justify-center gap-2 mb-4">
              <span class="text-[#B59E75] text-[10px]">◆</span>
              <span class="text-[#B59E75] text-xs font-semibold tracking-[0.2em] uppercase">Curated</span>
              <span class="text-[#B59E75] text-[10px]">◆</span>
            </div>
            <h2 class="text-4xl md:text-5xl lg:text-6xl font-serif italic text-primary-dark mb-6">
              Exclusive Rooms
            </h2>
          </div>
        </div>

        <div v-if="loading" class="flex justify-center items-center py-12">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
        </div>

        <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-8 md:gap-6 lg:gap-8" :class="{ 'reveal-visible': featuredVisible }">
          <div
            v-for="(room, index) in featuredRooms"
            :key="room.id"
            class="reveal-item"
            :class="`reveal-delay-${index + 2}`"
          >
            <RoomCard :room="room" />
          </div>
        </div>
        
        <div class="mt-16 text-center reveal-item reveal-delay-4" :class="{ 'reveal-visible': featuredVisible }">
          <router-link to="/rooms" class="inline-flex items-center justify-center px-8 py-3.5 border border-[#1A1A1A] text-[14px] font-medium tracking-wide text-[#1A1A1A] hover:bg-[#1A1A1A] hover:text-white transition-colors duration-300">
            Explore All Rooms
          </router-link>
        </div>
      </section>
    </div>

    <!-- CTA Section -->
    <section ref="ctaRef" class="relative py-24 md:py-32 flex items-center justify-center min-h-[400px] md:min-h-[500px]">
      <!-- Background Image -->
      <div class="absolute inset-0 z-0">
        <img src="@/assets/room2.jpg" alt="Resort landscape" class="w-full h-full object-cover" />
        <div class="absolute inset-0 bg-black/40"></div>
      </div>
      
      <!-- Content -->
      <div class="relative z-10 text-center px-4 max-w-4xl mx-auto" :class="{ 'reveal-visible': ctaVisible }">
        <div class="reveal-item reveal-delay-0">
          <div class="flex items-center justify-center gap-2 mb-6">
            <span class="text-white/80 text-[10px]">◆</span>
            <span class="text-white text-xs font-bold tracking-[0.2em] uppercase">The Invitation</span>
            <span class="text-white/80 text-[10px]">◆</span>
          </div>
          
          <h2 class="text-4xl md:text-5xl lg:text-7xl font-serif italic text-white mb-10 leading-tight tracking-wide">
            Reserve Your<br/>Stay Today
          </h2>
          
          <router-link to="/booking" class="inline-flex items-center justify-center px-10 py-3.5 border border-white/80 text-[13px] font-bold tracking-[0.1em] text-white uppercase hover:bg-white hover:text-[#1A1A1A] transition-colors duration-400 rounded-[100px]">
            Reserve Now
          </router-link>
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
  background-image: url('@/assets/hero-bg.jpg');
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
