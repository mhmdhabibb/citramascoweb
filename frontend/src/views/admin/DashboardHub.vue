<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import AdminDashboard from '@/views/admin/DashboardView.vue'
import ReceptionDashboard from '@/views/admin/ReceptionDashboardView.vue'

const authStore = useAuthStore()

const isReception = computed(() => authStore.role === 'reception')
</script>

<template>
  <ReceptionDashboard v-if="isReception" />
  <AdminDashboard v-else />
</template>