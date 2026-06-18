<script setup lang="ts">
import { useScrollReveal } from '@/composables/useScrollReveal'

const { elementRef: infoRef, isVisible: infoVisible } = useScrollReveal(0.2)
const { elementRef: mapRef, isVisible: mapVisible } = useScrollReveal(0.15)
</script>

<template>
  <div class="max-w-[1400px] mx-auto py-12 px-4 sm:px-6 lg:px-8">
    <div class="mb-10 text-center md:text-left contact-header">
      <h1 class="text-3xl font-extrabold tracking-tight">Contact Us</h1>
      <p class="mt-2 text-muted-foreground">Have any questions? We would love to hear from you.</p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-12">
      <!-- Contact Info -->
      <div ref="infoRef" class="space-y-8" :class="{ 'reveal-visible': infoVisible }">
        <div class="reveal-item reveal-delay-0">
          <h2 class="text-xl font-bold">Get in Touch</h2>
          <p class="text-muted-foreground mt-2">Feel free to reach out to us using the form or through our direct contact information below.</p>
        </div>

        <div class="space-y-4">
          <div class="contact-info-item reveal-item reveal-delay-1">
            <div class="contact-info-icon">
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </div>
            <div class="ml-3">
              <p class="font-medium text-foreground">Address</p>
              <p class="text-muted-foreground text-sm">Citramas Commercial Complex, Block B No. 14 - 15, Batam, Riau Islands Province, Indonesia</p>
            </div>
          </div>

          <div class="contact-info-item reveal-item reveal-delay-2">
            <div class="contact-info-icon">
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.94.725l.548 2.2a1 1 0 01-.321.988l-1.305.98a10.582 10.582 0 004.872 4.872l.98-1.305a1 1 0 01.988-.321l2.2.548a1 1 0 01.725.94V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
              </svg>
            </div>
            <div class="ml-3">
              <p class="font-medium text-foreground">Phone</p>
              <p class="text-muted-foreground text-sm">+62 21 1234 5678</p>
            </div>
          </div>

          <div class="contact-info-item reveal-item reveal-delay-3">
            <div class="contact-info-icon">
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
            </div>
            <div class="ml-3">
              <p class="font-medium text-foreground">Email</p>
              <p class="text-muted-foreground text-sm">info@citramascohotel.com</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Our Location Map -->
      <div ref="mapRef" :class="{ 'reveal-visible': mapVisible }">
        <div class="map-card reveal-item reveal-delay-0">
          <h2 class="text-xl font-bold mb-4">Our Location</h2>
          <div class="map-wrapper">
            <iframe
              src="https://www.google.com/maps?q=Komplek+Pertokoan+Citra+Mas,+Penuin,+Lubuk+Baja,+Batam,+Kepulauan+Riau&output=embed"
              width="100%"
              height="400"
              style="border: 0"
              :allowfullscreen="true"
              loading="lazy"
              referrerpolicy="no-referrer-when-downgrade"
              title="Ruko Citramas Blok B No 14-15, Batam, Kepulauan Riau"
            ></iframe>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Header entrance */
.contact-header {
  opacity: 0;
  transform: translateY(20px);
  animation: slideUp 0.6s cubic-bezier(0.16, 1, 0.3, 1) 0.1s forwards;
}

@keyframes slideUp {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Scroll reveal */
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
.reveal-delay-1 { transition-delay: 0.12s; }
.reveal-delay-2 { transition-delay: 0.24s; }
.reveal-delay-3 { transition-delay: 0.36s; }

/* Contact info items */
.contact-info-item {
  display: flex;
  align-items: flex-start;
  padding: 1rem;
  border-radius: 12px;
  transition: background-color 0.3s ease, transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.contact-info-item:hover {
  background-color: rgba(28, 22, 18, 0.03);
  transform: translateX(4px);
}

.contact-info-icon {
  flex-shrink: 0;
  color: var(--primary);
  transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.contact-info-item:hover .contact-info-icon {
  transform: scale(1.15);
}

/* Map card */
.map-card {
  padding: 0;
}

/* Google Maps */
.map-wrapper {
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid var(--border);
  box-shadow: 0 4px 16px rgba(28, 22, 18, 0.06);
  transition: box-shadow 0.3s ease;
}

.map-wrapper:hover {
  box-shadow: 0 8px 24px rgba(28, 22, 18, 0.12);
}

.map-wrapper iframe {
  display: block;
}
</style>
