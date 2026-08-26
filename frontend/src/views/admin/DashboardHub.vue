<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import UnifiedDashboard from '@/views/admin/ReceptionDashboardView.vue'
import InventoryDashboard from '@/views/admin/InventoryDashboardView.vue'
import FinanceDashboard from '@/views/admin/FinanceView.vue'
import HousekeepingDashboard from '@/views/admin/HousekeepingDashboardView.vue'
import StaffKpiDashboard from '@/views/admin/StaffKpiDashboardView.vue'
import { LayoutDashboard, BarChart3 } from 'lucide-vue-next'

const authStore = useAuthStore()

const isHousekeeping = computed(() => authStore.role === 'housekeeping')
const isInventory = computed(() => authStore.role === 'inventory')
const isFinance = computed(() => authStore.role === 'finance')
const isAdminOrManager = computed(() => authStore.role === 'admin' || authStore.role === 'manager')

const activeViewTab = ref<'operational' | 'kpi'>('operational')
</script>

<template>
  <HousekeepingDashboard v-if="isHousekeeping" />
  <InventoryDashboard v-else-if="isInventory" />
  <FinanceDashboard v-else-if="isFinance" />
  <div v-else class="admin-dashboard-container">
    <!-- Executive Tab Switcher for Admin & Manager -->
    <div v-if="isAdminOrManager" class="view-mode-tabs">
      <button 
        class="mode-tab-btn" 
        :class="{ active: activeViewTab === 'operational' }"
        @click="activeViewTab = 'operational'"
      >
        <LayoutDashboard class="tab-icon" />
        <span>Operations & Reservations</span>
      </button>
      <button 
        class="mode-tab-btn" 
        :class="{ active: activeViewTab === 'kpi' }"
        @click="activeViewTab = 'kpi'"
      >
        <BarChart3 class="tab-icon" />
        <span>Staff KPI & SLA Monitoring</span>
      </button>
    </div>

    <!-- Active View Display -->
    <StaffKpiDashboard v-if="isAdminOrManager && activeViewTab === 'kpi'" />
    <UnifiedDashboard v-else />
  </div>
</template>

<style scoped>
.admin-dashboard-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.view-mode-tabs {
  display: inline-flex;
  gap: 8px;
  background: #ffffff;
  border: 1px solid rgba(250, 235, 198, 0.9);
  padding: 6px;
  border-radius: 14px;
  width: fit-content;
  box-shadow: 0 2px 10px rgba(180, 140, 60, 0.05);
}

.mode-tab-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 18px;
  border-radius: 10px;
  border: none;
  background: transparent;
  color: #8c7a62;
  font-size: 0.85rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
}

.mode-tab-btn:hover {
  background: rgba(250, 235, 198, 0.4);
  color: #4a3a20;
}

.mode-tab-btn.active {
  background: #1a1612;
  color: #faebc6;
  box-shadow: 0 4px 12px rgba(26, 22, 18, 0.2);
}

.tab-icon {
  width: 16px;
  height: 16px;
}
</style>