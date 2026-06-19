<script setup lang="ts">
import { useScrollReveal } from '@/composables/useScrollReveal'

const { elementRef: footerRef, isVisible: footerVisible } = useScrollReveal(0.1)
</script>

<template>
  <footer ref="footerRef" class="bg-[#151515] text-white/70 py-20 px-4 sm:px-6 lg:px-8 mt-auto" :class="{ 'reveal-visible': footerVisible }">
    <div class="max-w-[1400px] mx-auto grid grid-cols-1 md:grid-cols-3 gap-16 md:gap-8 reveal-item reveal-delay-0">
      
      <!-- Left Column: Navigation -->
      <div>
        <h4 class="text-[26px] font-serif text-white mb-6" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">Navigation</h4>
        <ul class="space-y-3.5">
          <li><router-link to="/" class="hover:text-white transition-colors text-[14.5px]">Home</router-link></li>
          <li><router-link to="/#about" class="hover:text-white transition-colors text-[14.5px]">About Us</router-link></li>
          <li><router-link to="/rooms" class="hover:text-white transition-colors text-[14.5px]">Rooms</router-link></li>
          <li><router-link to="/contact" class="hover:text-white transition-colors text-[14.5px]">Contact</router-link></li>
          <li><router-link to="/offers" class="hover:text-white transition-colors text-[14.5px]">Offers</router-link></li>
        </ul>
      </div>

      <!-- Center Column: Logo, Subscribe, Socials -->
      <div class="flex flex-col items-center text-center">
        <!-- Logo -->
        <div class="mb-16">
          <svg class="w-8 h-8 text-white mx-auto mb-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z"/>
          </svg>
          <span class="text-xl font-medium tracking-wide text-white">CitraMas</span>
        </div>

        <!-- Subscribe -->
        <div class="w-full max-w-[320px] mb-12">
          <div class="flex items-center justify-between border-b border-white/20 pb-2 group focus-within:border-white/50 transition-colors">
            <input 
              type="email" 
              placeholder="Subscribe to our journal" 
              class="bg-transparent border-none outline-none text-[13px] text-white placeholder-white/50 w-full" 
            />
            <button class="text-white/50 hover:text-white transition-colors ml-2" aria-label="Subscribe">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="5" y1="12" x2="19" y2="12"></line>
                <polyline points="12 5 19 12 12 19"></polyline>
              </svg>
            </button>
          </div>
        </div>

        <!-- Social Icons -->
        <div class="flex items-center gap-5">
          <a href="#" class="text-white/70 hover:text-white transition-colors" aria-label="Facebook">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
              <path d="M18 2h-3a5 5 0 0 0-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 0 1 1-1h3z"/>
            </svg>
          </a>
          <a href="#" class="text-white/70 hover:text-white transition-colors" aria-label="X (Twitter)">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
              <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
            </svg>
          </a>
          <a href="#" class="text-white/70 hover:text-white transition-colors" aria-label="Instagram">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="2" y="2" width="20" height="20" rx="5" ry="5"/>
              <path d="M16 11.37A4 4 0 1 1 12.63 8 4 4 0 0 1 16 11.37z"/>
              <line x1="17.5" y1="6.5" x2="17.51" y2="6.5"/>
            </svg>
          </a>
          <a href="#" class="text-white/70 hover:text-white transition-colors" aria-label="LinkedIn">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
              <path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z"/>
              <rect x="2" y="9" width="4" height="12"/>
              <circle cx="4" cy="4" r="2"/>
            </svg>
          </a>
        </div>
      </div>

      <!-- Right Column: Contact -->
      <div class="md:text-right">
        <h4 class="text-[26px] font-serif text-white mb-6" style="font-family: 'Playfair Display', ui-serif, Georgia, Cambria, 'Times New Roman', Times, serif;">Contact</h4>
        <ul class="space-y-3.5 text-[14.5px]">
          <li><a href="info@citramascohotel.com" class="hover:text-white transition-colors">info@citramascohotel.com</a></li>
          <li><a href="tel:+628111222333" class="hover:text-white transition-colors">+62 811-1222-333</a></li>
          <li>Citramas Commercial Complex, Block B No. 14 - 15</li>
          <li>Batam, Riau Islands Province, Indonesia</li>
        </ul>
      </div>
      
    </div>
  </footer>
</template>

<style scoped>
/* Scroll Reveal Animations */
.reveal-item {
  opacity: 0;
  transform: translateY(24px);
  transition: opacity 0.8s cubic-bezier(0.16, 1, 0.3, 1),
              transform 0.8s cubic-bezier(0.16, 1, 0.3, 1);
}

.reveal-visible .reveal-item {
  opacity: 1;
  transform: translateY(0);
}

.reveal-delay-0 { transition-delay: 0s; }
</style>
