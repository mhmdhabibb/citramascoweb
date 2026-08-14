<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import UnifiedDashboard from '@/views/admin/ReceptionDashboardView.vue'
import InventoryDashboard from '@/views/admin/InventoryDashboardView.vue'
import FinanceDashboard from '@/views/admin/FinanceView.vue'

const authStore = useAuthStore()

const isInventory = computed(() => authStore.role === 'inventory')
const isFinance = computed(() => authStore.role === 'finance')
</script>

<template>
  <InventoryDashboard v-if="isInventory" />
  <FinanceDashboard v-else-if="isFinance" />
  <UnifiedDashboard v-else />
</template>