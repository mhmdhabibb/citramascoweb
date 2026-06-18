<script setup lang="ts">
import { useScrollReveal } from '@/composables/useScrollReveal'

const currentYear = new Date().getFullYear()
const { elementRef: footerRef, isVisible: footerVisible } = useScrollReveal(0.1)
</script>

<template>
  <footer ref="footerRef" class="footer" :class="{ 'footer-visible': footerVisible }">
    <div class="footer-container">
      <!-- Top section: Brand + Link Columns -->
      <div class="footer-top">
        <!-- Brand -->
        <div class="footer-brand footer-item footer-delay-0">
          <span class="footer-logo">
            <span class="footer-logo-accent">CitraMas</span> Co Living
          </span>
          <p class="footer-tagline">Comfortable stay and exceptional reservation services.</p>
        </div>

        <!-- Link Columns -->
        <div class="footer-links footer-item footer-delay-1">
          <div class="footer-col">
            <h4 class="footer-col-title">Main Page</h4>
            <ul class="footer-col-list">
              <li><router-link to="/" class="footer-link">Home</router-link></li>
              <li><router-link to="/rooms" class="footer-link">Rooms</router-link></li>
              <li><router-link to="/contact" class="footer-link">Contact</router-link></li>
            </ul>
          </div>

          <div class="footer-col">
            <h4 class="footer-col-title">Quick Links</h4>
            <ul class="footer-col-list">
              <li><router-link to="/offers" class="footer-link">Offers</router-link></li>
              <li><router-link to="/booking" class="footer-link">Booking</router-link></li>
              <li><router-link to="/rooms" class="footer-link">Suites</router-link></li>
            </ul>
          </div>

          <div class="footer-col">
            <h4 class="footer-col-title">Legal</h4>
            <ul class="footer-col-list">
              <li><a href="#" class="footer-link">Privacy Policy</a></li>
              <li><a href="#" class="footer-link">Terms &amp; Condition</a></li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Bottom bar: Copyright + Social Icons -->
      <div class="footer-bottom footer-item footer-delay-2">
        <p class="footer-copyright">&copy; {{ currentYear }} CitraMas Co Living. All rights reserved.</p>

        <div class="footer-socials">
          <!-- Instagram -->
          <a href="#" class="footer-social-link" aria-label="Instagram">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="2" y="2" width="20" height="20" rx="5" ry="5"/>
              <path d="M16 11.37A4 4 0 1 1 12.63 8 4 4 0 0 1 16 11.37z"/>
              <line x1="17.5" y1="6.5" x2="17.51" y2="6.5"/>
            </svg>
          </a>
          <!-- Facebook -->
          <a href="#" class="footer-social-link" aria-label="Facebook">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M18 2h-3a5 5 0 0 0-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 0 1 1-1h3z"/>
            </svg>
          </a>
          <!-- X (Twitter) -->
          <a href="#" class="footer-social-link" aria-label="X">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
              <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
            </svg>
          </a>
          <!-- LinkedIn -->
          <a href="#" class="footer-social-link" aria-label="LinkedIn">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z"/>
              <rect x="2" y="9" width="4" height="12"/>
              <circle cx="4" cy="4" r="2"/>
            </svg>
          </a>
        </div>
      </div>
    </div>
  </footer>
</template>

<style scoped>
.footer {
  background: var(--background);
  color: var(--foreground);
  margin-top: auto;
}

.footer-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 3rem 2rem 2rem;
}

/* Reveal animations */
.footer-item {
  opacity: 0;
  transform: translateY(12px);
  transition: opacity 0.6s cubic-bezier(0.16, 1, 0.3, 1),
              transform 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

.footer-visible .footer-item {
  opacity: 1;
  transform: translateY(0);
}

.footer-delay-0 { transition-delay: 0s; }
.footer-delay-1 { transition-delay: 0.1s; }
.footer-delay-2 { transition-delay: 0.2s; }

/* Top section */
.footer-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 3rem;
}

/* Brand */
.footer-brand {
  flex-shrink: 0;
  max-width: 280px;
}

.footer-logo {
  font-size: 18px;
  font-weight: 700;
  color: var(--foreground);
}

.footer-logo-accent {
  color: var(--primary);
}

.footer-tagline {
  font-size: 14px;
  color: var(--muted-foreground);
  margin-top: 0.5rem;
  line-height: 1.5;
}

/* Link columns */
.footer-links {
  display: flex;
  gap: 4rem;
}

.footer-col-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--foreground);
  margin: 0 0 1rem;
  letter-spacing: 0.01em;
}

.footer-col-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}

.footer-link {
  font-size: 14px;
  color: var(--muted-foreground);
  text-decoration: none;
  transition: color 0.3s ease;
  position: relative;
}

.footer-link:hover {
  color: var(--primary);
}

/* Bottom bar */
.footer-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px solid var(--border);
  margin-top: 2.5rem;
  padding-top: 1.5rem;
}

.footer-copyright {
  font-size: 13px;
  color: var(--muted-foreground);
  margin: 0;
}

/* Social icons */
.footer-socials {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.footer-social-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  color: var(--muted-foreground);
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.footer-social-link:hover {
  color: var(--foreground);
  background: rgba(28, 22, 18, 0.05);
  transform: translateY(-2px);
}

/* Responsive */
@media (max-width: 768px) {
  .footer-top {
    flex-direction: column;
    gap: 2rem;
  }

  .footer-brand {
    max-width: 100%;
  }

  .footer-links {
    gap: 2.5rem;
    flex-wrap: wrap;
  }

  .footer-bottom {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
}
</style>
