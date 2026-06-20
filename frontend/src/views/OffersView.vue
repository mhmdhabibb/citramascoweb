<script setup lang="ts">
import { useScrollReveal } from '@/composables/useScrollReveal'

const { elementRef: headerRef, isVisible: headerVisible } = useScrollReveal(0.2)
const { elementRef: cardsRef, isVisible: cardsVisible } = useScrollReveal(0.1)

const offers = [
  { 
    id: 1,
    category: 'RENEWAL',
    title: 'Spa Rituals',
    description: 'Restorative spa rituals inspired by calm, warmth, and deep relaxation.',
    image: 'https://images.unsplash.com/photo-1540555700478-4be289fbecef?q=80&w=800&auto=format&fit=crop'
  },
  {
    id: 2,
    category: 'ATMOSPHERE',
    title: 'Coastal Dining',
    description: 'Seasonal coastal dining crafted with refined flavors and regional ingredients.',
    image: 'https://images.unsplash.com/photo-1544148103-0773bf10d330?q=80&w=800&auto=format&fit=crop'
  },
  {
    id: 3,
    category: 'DISCOVERY',
    title: 'Curated Escapes',
    description: 'Private coastal escapes designed for immersive journeys and peaceful exploration.',
    image: 'https://images.unsplash.com/photo-1506501139174-099022df5260?q=80&w=800&auto=format&fit=crop'
  }
]
</script>

<template>
  <div class="offers-page min-h-screen pt-24 pb-20">
    <div class="max-w-[1200px] mx-auto px-4 sm:px-6 lg:px-8">
      
      <!-- Header Section -->
      <div 
        ref="headerRef" 
        class="text-center max-w-2xl mx-auto mb-16"
        :class="{ 'reveal-visible': headerVisible }"
      >
        <div class="reveal-item reveal-delay-0">
          <div class="flex items-center justify-center gap-2 mb-4">
            <span class="text-accent text-[10px]">◆</span>
            <span class="text-accent text-xs font-semibold tracking-[0.2em] uppercase">Our Offers</span>
            <span class="text-accent text-[10px]">◆</span>
          </div>
          
          <h1 class="text-4xl md:text-5xl lg:text-6xl font-serif italic text-primary-dark mb-6">
            Comfort & Experiences
          </h1>
          
          <p class="text-gray-600 text-[15px] leading-relaxed max-w-lg mx-auto">
            Curated experiences shaped by wellness, seasonal<br class="hidden md:block" />
            dining, and moments designed for quiet luxury living.
          </p>
        </div>
      </div>

      <!-- Cards Grid -->
      <div 
        ref="cardsRef"
        class="grid grid-cols-1 md:grid-cols-3 gap-8 md:gap-6 lg:gap-8"
        :class="{ 'reveal-visible': cardsVisible }"
      >
        <div 
          v-for="(offer, index) in offers" 
          :key="offer.id"
          class="offer-card group reveal-item border border-[#EAE1D8] bg-[#F5F1E8] rounded-[20px] overflow-hidden transition-all duration-300 hover:shadow-lg hover:shadow-[#EAE1D8]/50"
          :class="`reveal-delay-${index + 1}`"
        >
          <!-- Image Container -->
          <div class="overflow-hidden aspect-[4/5] relative">
            <img 
              :src="offer.image" 
              :alt="offer.title"
              class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105"
            />
            <div class="absolute inset-0 bg-black/10 opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
          </div>
          
          <!-- Text Content -->
          <div class="p-6">
            <div class="text-accent text-[11px] font-bold tracking-[0.15em] uppercase mb-2.5">
              {{ offer.category }}
            </div>
            <h3 class="text-[26px] font-serif text-primary-dark mb-3">
              {{ offer.title }}
            </h3>
            <p class="text-[#6B5E52] text-[14.5px] leading-relaxed">
              {{ offer.description }}
            </p>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped>
.offers-page {
  background-color: #FAF7F2;
}

.text-accent {
  color: #B59E75;
}

.text-primary-dark {
  color: #1A1A1A;
}

/* Serif font style for headings to match the image */
.font-serif {
  font-family: 'Playfair Display', ui-serif, Georgia, Cambria, "Times New Roman", Times, serif;
}

/* Scroll Reveal Animations */
.reveal-item {
  opacity: 0;
  transform: translateY(30px);
  transition: opacity 0.8s cubic-bezier(0.16, 1, 0.3, 1),
              transform 0.8s cubic-bezier(0.16, 1, 0.3, 1);
}

.reveal-visible .reveal-item {
  opacity: 1;
  transform: translateY(0);
}

.reveal-delay-0 { transition-delay: 0s; }
.reveal-delay-1 { transition-delay: 0.15s; }
.reveal-delay-2 { transition-delay: 0.3s; }
.reveal-delay-3 { transition-delay: 0.45s; }
</style>
