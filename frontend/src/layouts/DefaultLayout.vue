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
          name="fade"
          mode="out-in"
        >
          <component :is="Component" />
        </Transition>
      </RouterView>
    </main>
    <Footer />
  </div>
</template>
