<script setup lang="ts">
import { useScrollReveal } from '@/composables/useScrollReveal'

const { elementRef: ctaRef, isVisible: ctaVisible } = useScrollReveal(0.2)
</script>

<template>
  <section ref="ctaRef" class="cta-section">
    <div class="cta-container" :class="{ 'cta-visible': ctaVisible }">
      <!-- Decorative gradient orb -->
      <div class="cta-orb"></div>

      <h2 class="cta-heading cta-item cta-delay-0">
        Your Perfect Stay Awaits at<br />
        CitraMas Co Living
      </h2>
      <p class="cta-description cta-item cta-delay-1">
        Experience modern co-living in the heart of Batam. Beautifully designed rooms,
        premium amenities, and a welcoming community — all at an affordable price.
      </p>
      <div class="cta-buttons cta-item cta-delay-2">
        <router-link to="/booking" class="cta-btn-primary">
          Book Now
        </router-link>
        <router-link to="/contact" class="cta-btn-secondary">
          Contact Us
        </router-link>
      </div>
    </div>
  </section>
</template>

<style scoped>
.cta-section {
  padding: 0 1.5rem;
  margin-top: var(--spacing-section);
  margin-bottom: var(--spacing-section);
}

.cta-container {
  position: relative;
  max-width: 1400px;
  margin: 0 auto;
  background: #1C1612;
  border-radius: 24px;
  padding: 5rem 2rem;
  text-align: center;
  overflow: hidden;
}

.cta-orb {
  position: absolute;
  bottom: -60px;
  left: 50%;
  transform: translateX(-50%);
  width: 500px;
  height: 280px;
  background: radial-gradient(
    ellipse at center,
    rgba(200, 120, 80, 0.35) 0%,
    rgba(180, 80, 100, 0.2) 40%,
    transparent 70%
  );
  border-radius: 50%;
  pointer-events: none;
  filter: blur(40px);
}

/* Reveal animations */
.cta-item {
  opacity: 0;
  transform: translateY(24px);
  transition: opacity 0.7s cubic-bezier(0.16, 1, 0.3, 1),
              transform 0.7s cubic-bezier(0.16, 1, 0.3, 1);
}

.cta-visible .cta-item {
  opacity: 1;
  transform: translateY(0);
}

.cta-delay-0 { transition-delay: 0s; }
.cta-delay-1 { transition-delay: 0.12s; }
.cta-delay-2 { transition-delay: 0.24s; }

.cta-heading {
  position: relative;
  z-index: 1;
  font-size: clamp(28px, 4vw, 44px);
  font-weight: 700;
  line-height: 1.2;
  letter-spacing: -0.02em;
  color: #FFFFFF;
  margin: 0 0 1.25rem;
}

.cta-description {
  position: relative;
  z-index: 1;
  font-size: clamp(15px, 2vw, 18px);
  line-height: 1.7;
  color: rgba(255, 255, 255, 0.65);
  max-width: 580px;
  margin: 0 auto 2.5rem;
}

.cta-buttons {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.cta-btn-primary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.8rem 2rem;
  border-radius: 100px;
  font-size: 15px;
  font-weight: 600;
  color: #1C1612;
  background: #F8F5F1;
  text-decoration: none;
  border: 1px solid transparent;
  transition: all 0.35s cubic-bezier(0.16, 1, 0.3, 1);
}

.cta-btn-primary:hover {
  background: #FFFFFF;
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(248, 245, 241, 0.25);
}

.cta-btn-secondary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.8rem 2rem;
  border-radius: 100px;
  font-size: 15px;
  font-weight: 500;
  color: #FFFFFF;
  background: transparent;
  text-decoration: none;
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.35s cubic-bezier(0.16, 1, 0.3, 1);
}

.cta-btn-secondary:hover {
  background: rgba(255, 255, 255, 0.1);
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(255, 255, 255, 0.08);
}

@media (max-width: 640px) {
  .cta-container {
    padding: 3.5rem 1.5rem;
    border-radius: 20px;
  }

  .cta-buttons {
    flex-direction: column;
  }

  .cta-btn-primary,
  .cta-btn-secondary {
    width: 100%;
  }
}
</style>
