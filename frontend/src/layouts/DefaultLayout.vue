<script setup lang="ts">
import Navbar from '@/components/navbar/Navbar.vue'
import Footer from '@/components/footer/Footer.vue'
</script>

<template>
  <div class="min-h-screen flex flex-col bg-background text-foreground font-sans">
    <Navbar />
    <main class="flex-grow">
      <RouterView v-slot="{ Component }">
        <Transition
          name="page"
          mode="out-in"
        >
          <component :is="Component" />
        </Transition>
      </RouterView>
    </main>
    <Footer />
  </div>
</template>

<style>
/* Smooth page transitions */
.page-enter-active {
  transition: opacity 0.4s ease, transform 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

.page-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}

.page-enter-from {
  opacity: 0;
  transform: translateY(12px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
