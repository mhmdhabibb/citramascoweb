<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { inventoryService } from '@/services/admin/inventoryService'
import { useToastStore } from '@/stores/toastStore'
import type { InventoryItem, InventoryTransaction } from '@/types'
import { AlertTriangle, Package, History, TrendingDown } from 'lucide-vue-next'

const toastStore = useToastStore()
const loading = ref(false)

const items = ref<InventoryItem[]>([])
const transactions = ref<InventoryTransaction[]>([])

const totalItems = computed(() => items.value.length)
const lowStockCount = computed(() => items.value.filter((i) => i.reorder_level > 0 && i.current_stock <= i.reorder_level).length)
const recentUsages = computed(() => transactions.value.filter(t => t.type === 'usage').slice(0, 5))
const recentStockIns = computed(() => transactions.value.filter(t => t.type === 'stock-in').slice(0, 5))

const fetchData = async () => {
  try {
    loading.value = true
    items.value = await inventoryService.getItems()
    transactions.value = await inventoryService.getTransactions()
  } catch (error) {
    toastStore.error((error as Error).message || 'Failed to load inventory data')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric', month: 'short', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  })
}
</script>

<template>
  <div class="inventory-dashboard">
    <div class="dashboard-header">
      <div>
        <h1 class="page-title">Inventory Dashboard</h1>
        <p class="page-subtitle">Overview of your items and stock status.</p>
      </div>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Loading dashboard...</p>
    </div>

    <div v-else>
      <div class="stats-grid">
        <!-- Total Items -->
        <div class="stat-card">
          <div class="stat-icon-wrapper blue">
            <Package class="stat-icon" />
          </div>
          <div class="stat-content">
            <h3 class="stat-label">Total Inventory Items</h3>
            <p class="stat-value">{{ totalItems }}</p>
          </div>
        </div>

        <!-- Low Stock Alerts -->
        <div class="stat-card" :class="{ 'warning-bg': lowStockCount > 0 }">
          <div class="stat-icon-wrapper orange">
            <AlertTriangle class="stat-icon" />
          </div>
          <div class="stat-content">
            <h3 class="stat-label">Low Stock Alerts</h3>
            <p class="stat-value">{{ lowStockCount }}</p>
          </div>
        </div>
      </div>

      <div class="tables-grid">
        <!-- Recent Usages -->
        <div class="card table-card">
          <div class="card-header">
            <TrendingDown class="card-icon" />
            <h2 class="card-title">Recent Usages</h2>
          </div>
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Date</th>
                  <th>Item</th>
                  <th>Quantity</th>
                  <th>Note</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="t in recentUsages" :key="t.id">
                  <td class="date-col">{{ formatDate(t.created_at) }}</td>
                  <td>{{ t.item?.name || 'Unknown Item' }}</td>
                  <td class="qty-col negative">-{{ t.quantity }}</td>
                  <td>{{ t.note || '-' }}</td>
                </tr>
                <tr v-if="recentUsages.length === 0">
                  <td colspan="4" class="empty-state">No recent usage.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Recent Stock In -->
        <div class="card table-card">
          <div class="card-header">
            <History class="card-icon" />
            <h2 class="card-title">Recent Stock In</h2>
          </div>
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Date</th>
                  <th>Item</th>
                  <th>Quantity</th>
                  <th>Reference</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="t in recentStockIns" :key="t.id">
                  <td class="date-col">{{ formatDate(t.created_at) }}</td>
                  <td>{{ t.item?.name || 'Unknown Item' }}</td>
                  <td class="qty-col positive">+{{ t.quantity }}</td>
                  <td>{{ t.reference || '-' }}</td>
                </tr>
                <tr v-if="recentStockIns.length === 0">
                  <td colspan="4" class="empty-state">No recent stock ins.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.inventory-dashboard {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-dark);
  margin-bottom: 4px;
}

.page-subtitle {
  color: var(--text-muted);
  font-size: 0.875rem;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}

.stat-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 20px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
  border: 1px solid var(--border-light);
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
}

.warning-bg {
  background: linear-gradient(to right, #fff, #fffaf0);
  border-color: #fed7aa;
}

.stat-icon-wrapper {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon-wrapper.blue {
  background: #eff6ff;
  color: #3b82f6;
}

.stat-icon-wrapper.orange {
  background: #fff7ed;
  color: #ea580c;
}

.stat-icon {
  width: 28px;
  height: 28px;
}

.stat-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-dark);
  line-height: 1;
}

/* Tables Grid */
.tables-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

@media (max-width: 1024px) {
  .tables-grid {
    grid-template-columns: 1fr;
  }
}

.card {
  background: white;
  border-radius: 16px;
  border: 1px solid var(--border-light);
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}

.card-header {
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-light);
  background: #fafafa;
  display: flex;
  align-items: center;
  gap: 12px;
}

.card-icon {
  width: 20px;
  height: 20px;
  color: var(--text-muted);
}

.card-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-dark);
}

.table-container {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.data-table th {
  padding: 12px 24px;
  background: white;
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-bottom: 1px solid var(--border-light);
}

.data-table td {
  padding: 16px 24px;
  font-size: 0.875rem;
  color: var(--text-dark);
  border-bottom: 1px solid var(--border-light);
}

.data-table tbody tr:last-child td {
  border-bottom: none;
}

.date-col {
  color: var(--text-muted) !important;
  font-size: 0.8125rem !important;
  white-space: nowrap;
}

.qty-col {
  font-weight: 600;
}

.qty-col.positive {
  color: #10b981;
}

.qty-col.negative {
  color: #ef4444;
}

.empty-state {
  text-align: center;
  padding: 32px !important;
  color: var(--text-muted) !important;
  font-style: italic;
}

/* Loading */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 0;
  color: var(--text-muted);
  gap: 16px;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #e2e8f0;
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
