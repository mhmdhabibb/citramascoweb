<script setup lang="ts">
interface Room {
  id: number
  name: string
  price: number
  image: string
  capacity: number
  description: string
}

defineProps<{
  room: Room
}>()
</script>

<template>
  <div class="bg-card text-card-foreground rounded-lg border border-border overflow-hidden shadow-sm hover:shadow-md transition-shadow">
    <img :src="room.image" :alt="room.name" class="w-full h-48 object-cover" />
    <div class="p-6">
      <h3 class="text-xl font-bold">{{ room.name }}</h3>
      <p class="text-sm text-muted-foreground mt-1 mb-4 line-clamp-2">{{ room.description }}</p>
      
      <div class="flex items-center justify-between">
        <div>
          <span class="text-xs text-muted-foreground">Capacity: {{ room.capacity }} guests</span>
          <p class="text-lg font-bold text-primary">Rp {{ room.price.toLocaleString('id-ID') }} <span class="text-xs font-normal text-muted-foreground">/ night</span></p>
        </div>
        <router-link :to="`/rooms/${room.id}`" class="inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors bg-primary text-primary-foreground hover:bg-primary/90 h-9 px-4 py-2">
          Details
        </router-link>
      </div>
    </div>
  </div>
</template>
