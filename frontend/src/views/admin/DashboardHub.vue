<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import AdminDashboard from '@/views/admin/DashboardView.vue'
import ReceptionDashboard from '@/views/admin/ReceptionDashboardView.vue'
import InventoryDashboard from '@/views/admin/InventoryDashboardView.vue'

const authStore = useAuthStore()

const isReception = computed(() => authStore.role === 'reception')
const isInventory = computed(() => authStore.role === 'inventory')
</script>

<template>
  <ReceptionDashboard v-if="isReception" />
  <InventoryDashboard v-else-if="isInventory" />
  <AdminDashboard v-else />
</template>