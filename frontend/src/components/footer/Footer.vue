<script setup lang="ts">
const currentYear = new Date().getFullYear()
</script>

<template>
  <footer class="bg-card text-card-foreground border-t border-border mt-auto">
    <div class="max-w-7xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
      <div class="flex flex-col md:flex-row justify-between items-center gap-4">
        <div>
          <span class="text-lg font-bold text-primary">CitraMas Co Living</span>
          <p class="text-sm text-muted-foreground mt-1">Comfortable stay and exceptional reservation services.</p>
        </div>
        <div class="flex space-x-6 text-sm text-muted-foreground">
          <router-link to="/" class="hover:text-primary">Home</router-link>
          <router-link to="/rooms" class="hover:text-primary">Rooms & Suites</router-link>
          <router-link to="/offers" class="hover:text-primary">Offers</router-link>
          <router-link to="/contact" class="hover:text-primary">Contact</router-link>
        </div>
      </div>
      <div class="border-t border-border mt-8 pt-8 flex justify-between items-center text-xs text-muted-foreground">
        <p>&copy; {{ currentYear }} CitraMas Co Living. All rights reserved.</p>
      </div>
    </div>
  </footer>
</template>
