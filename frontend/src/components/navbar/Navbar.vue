<script setup lang="ts">
import { ref } from 'vue'

const menuOpen = ref(false)
</script>

<template>
  <nav class="navbar">
    <div class="navbar-inner">
      <!-- Logo (Left) -->
      <div class="navbar-logo">
        <router-link to="/" class="logo-link">
          <img src="@/assets/logo.svg" alt="CitraMas Co Living" class="logo-img" />
        </router-link>
      </div>

      <!-- Center Menu Links (Desktop) -->
      <div class="nav-center">
        <router-link to="/" class="nav-link">Home</router-link>
        <router-link to="/rooms" class="nav-link">Rooms & Suites</router-link>
        <router-link to="/offers" class="nav-link">Offers</router-link>
        <router-link to="/contact" class="nav-link">Contact</router-link>
      </div>

      <!-- Right Section (Desktop) -->
      <div class="nav-right">
        <router-link to="/booking" class="book-now-btn">Book Now</router-link>
      </div>

      <!-- Mobile Hamburger -->
      <button @click="menuOpen = !menuOpen" class="mobile-menu-btn">
        <span class="sr-only">Open main menu</span>
        <svg class="h-6 w-6" stroke="currentColor" fill="none" viewBox="0 0 24 24">
          <path :class="{'hidden': menuOpen, 'inline-flex': !menuOpen }" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          <path :class="{'hidden': !menuOpen, 'inline-flex': menuOpen }" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Mobile menu -->
    <div :class="{'block': menuOpen, 'hidden': !menuOpen}" class="mobile-menu">
      <div class="mobile-menu-inner">
        <router-link to="/" class="mobile-link" @click="menuOpen = false">Home</router-link>
        <router-link to="/rooms" class="mobile-link" @click="menuOpen = false">Rooms & Suites</router-link>
        <router-link to="/offers" class="mobile-link" @click="menuOpen = false">Offers</router-link>
        <router-link to="/contact" class="mobile-link" @click="menuOpen = false">Contact</router-link>
        <router-link to="/booking" class="mobile-book-btn" @click="menuOpen = false">Book Now</router-link>
      </div>
    </div>
  </nav>
</template>

<style scoped>
.navbar {
  background-color: #F8F5F1;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  position: sticky;
  top: 0;
  z-index: 50;
}

.navbar-inner {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 72px;
  position: relative;
}

.navbar-logo {
  flex-shrink: 0;
  z-index: 1;
}

.logo-link {
  display: flex;
  align-items: center;
}

.logo-img {
  height: 36px;
  width: auto;
}

/* Center nav links - absolutely centered */
.nav-center {
  display: none;
  align-items: center;
  gap: 2rem;
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
}

@media (min-width: 768px) {
  .nav-center {
    display: flex;
  }
}

.nav-link {
  font-size: 16px;
  font-weight: 500;
  color: #3D3329;
  text-decoration: none;
  letter-spacing: 0.01em;
  padding: 0.5rem 0;
  position: relative;
  transition: color 0.3s ease;
  white-space: nowrap;
}

.nav-link::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 0;
  height: 2px;
  background-color: #1C1612;
  transition: width 0.3s ease;
}

.nav-link:hover {
  color: #1C1612;
}

.nav-link:hover::after {
  width: 100%;
}

.nav-link.router-link-exact-active {
  color: #1C1612;
}

.nav-link.router-link-exact-active::after {
  width: 100%;
}

/* Right section */
.nav-right {
  display: none;
  align-items: center;
  z-index: 1;
}

@media (min-width: 768px) {
  .nav-right {
    display: flex;
  }
}

.book-now-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.6rem 1.5rem;
  border-radius: 100px;
  font-size: 16px;
  font-weight: 500;
  color: #fff;
  background-color: #1C1612;
  text-decoration: none;
  letter-spacing: 0.01em;
  transition: background-color 0.3s ease, transform 0.2s ease;
  white-space: nowrap;
}

.book-now-btn:hover {
  background-color: #3D3329;
  transform: translateY(-1px);
}

/* Mobile hamburger */
.mobile-menu-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem;
  border-radius: 6px;
  border: none;
  background: transparent;
  color: #3D3329;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.mobile-menu-btn:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

@media (min-width: 768px) {
  .mobile-menu-btn {
    display: none;
  }
}

/* Mobile menu */
.mobile-menu {
  background: #F8F5F1;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
}

@media (min-width: 768px) {
  .mobile-menu {
    display: none !important;
  }
}

.mobile-menu-inner {
  padding: 1rem 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.mobile-link {
  display: block;
  padding: 0.75rem 1rem;
  font-size: 16px;
  font-weight: 500;
  color: #3D3329;
  text-decoration: none;
  letter-spacing: 0.01em;
  border-radius: 6px;
  transition: background-color 0.2s ease, color 0.2s ease;
}

.mobile-link:hover {
  background-color: rgba(0, 0, 0, 0.04);
  color: #1C1612;
}

.mobile-link.router-link-exact-active {
  color: #1C1612;
  background-color: rgba(28, 22, 18, 0.06);
}

.mobile-book-btn {
  display: block;
  margin-top: 0.5rem;
  padding: 0.75rem 1rem;
  font-size: 16px;
  font-weight: 500;
  color: #fff;
  background-color: #1C1612;
  text-decoration: none;
  text-align: center;
  letter-spacing: 0.01em;
  border-radius: 100px;
  transition: background-color 0.3s ease;
}

.mobile-book-btn:hover {
  background-color: #3D3329;
}
</style>
