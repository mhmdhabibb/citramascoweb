<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { inventoryService } from '@/services/admin/inventoryService'
import { useToastStore } from '@/stores/toastStore'
import type { InventoryItem, InventoryReport, InventoryTransaction } from '@/types'

const toastStore = useToastStore()

const activeTab = ref('items')
const loading = ref(false)

const items = ref<InventoryItem[]>([])
const transactions = ref<InventoryTransaction[]>([])
const report = ref<InventoryReport | null>(null)

const tabs = [
  { value: 'items', label: 'Inventory Items' },
  { value: 'usage', label: 'Usage Log (Book)' },
  { value: 'balance', label: 'Stock Balance Check' },
  { value: 'report', label: 'Usage Report' },
]

// ---------- Items ----------
const showItemAdd = ref(false)
const savingItem = ref(false)
const itemForm = ref({
  name: '',
  category: '',
  unit: '',
  current_stock: 0,
  reorder_level: 0,
})

const lowStockItems = computed(() => items.value.filter((i) => i.reorder_level > 0 && i.current_stock <= i.reorder_level))

const fetchItems = async () => {
  try {
    loading.value = true
    items.value = await inventoryService.getItems()
  } catch (error) {
    toastStore.error((error as Error).message || 'Failed to load items')
  } finally {
    loading.value = false
  }
}

const openAddItem = () => {
  itemForm.value = { name: '', category: '', unit: '', current_stock: 0, reorder_level: 0 }
  showItemAdd.value = true
}

const submitItem = async () => {
  if (!itemForm.value.name || !itemForm.value.category || !itemForm.value.unit) {
    toastStore.error('Please fill name, category and unit')
    return
  }
  try {
    savingItem.value = true
    await inventoryService.createItem({
      ...itemForm.value,
      current_stock: Number(itemForm.value.current_stock) || 0,
      reorder_level: Number(itemForm.value.reorder_level) || 0,
    })
    toastStore.success('Item created successfully')
    showItemAdd.value = false
    await fetchItems()
  } catch (error) {
    toastStore.error((error as Error).message || 'Failed to create item')
  } finally {
    savingItem.value = false
  }
}

const deleteItem = async (item: InventoryItem) => {
  if (!confirm(`Delete item "${item.name}"?`)) return
  try {
    loading.value = true
    const msg = await inventoryService.deleteItem(item.id)
    toastStore.success(msg || 'Item deleted')
    await fetchItems()
  } catch (error) {
    toastStore.error((error as Error).message || 'Failed to delete item')
  } finally {
    loading.value = false
  }
}

// ---------- Record stock-in / usage modal ----------
const showRecordModal = ref(false)
const recordMode = ref<'stock-in' | 'usage'>('usage')
const recordItem = ref<InventoryItem | null>(null)
const savingRecord = ref(false)
const recordForm = ref({ quantity: 1, reference: '', note: '' })

const openRecord = (item: InventoryItem, mode: 'stock-in' | 'usage') => {
  recordItem.value = item
  recordMode.value = mode
  recordForm.value = { quantity: 1, reference: '', note: '' }
  showRecordModal.value = true
}

const submitRecord = async () => {
  if (!recordItem.value) return
  const qty = Number(recordForm.value.quantity)
  if (!qty || qty <= 0) {
    toastStore.error('Quantity must be greater than zero')
    return
  }
  try {
    savingRecord.value = true
    const msg =
      recordMode.value === 'stock-in'
        ? await inventoryService.stockIn(recordItem.value.id, {
            quantity: qty,
            reference: recordForm.value.reference,
            note: recordForm.value.note,
          })
        : await inventoryService.usage(recordItem.value.id, {
            quantity: qty,
            reference: recordForm.value.reference,
            note: recordForm.value.note,
          })
    toastStore.success(msg || 'Recorded successfully')
    showRecordModal.value = false
    await fetchItems()
    await fetchTransactions()
  } catch (error) {
    toastStore.error((error as Error).message || 'Failed to record')
  } finally {
    savingRecord.value = false
  }
}

// ---------- Transactions ----------
const txnFrom = ref('')
const txnTo = ref('')

const fetchTransactions = async () => {
  try {
    transactions.value = await inventoryService.getTransactions({
      from: txnFrom.value || undefined,
      to: txnTo.value || undefined,
    })
  } catch (error) {
    toastStore.error((error as Error).message || 'Failed to load transactions')
  }
}

// ---------- Stock balance check ----------
const showBalanceModal = ref(false)
const balanceEntry = ref<{ item_id: string; system_stock: number; actual_stock: number; note: string }>({
  item_id: '',
  system_stock: 0,
  actual_stock: 0,
  note: '',
})
const savingBalance = ref(false)

const openBalance = (item: InventoryItem) => {
  balanceEntry.value = {
    item_id: item.id,
    system_stock: item.current_stock,
    actual_stock: item.current_stock,
    note: '',
  }
  showBalanceModal.value = true
}

const balanceDiff = computed(() => Number(balanceEntry.value.actual_stock) - balanceEntry.value.system_stock)

const submitBalance = async () => {
  if (!balanceEntry.value.item_id) return
  try {
    savingBalance.value = true
    const msg = await inventoryService.stockTake({
      item_id: balanceEntry.value.item_id,
      actual_stock: Number(balanceEntry.value.actual_stock) || 0,
      note: balanceEntry.value.note,
    })
    toastStore.success(msg || 'Reconciled successfully')
    showBalanceModal.value = false
    await fetchItems()
  } catch (error) {
    toastStore.error((error as Error).message || 'Failed to reconcile stock')
  } finally {
    savingBalance.value = false
  }
}

// ---------- Report ----------
const repFrom = ref('')
const repTo = ref('')

const fetchReport = async () => {
  try {
    loading.value = true
    report.value = await inventoryService.getReport({
      from: repFrom.value || undefined,
      to: repTo.value || undefined,
    })
  } catch (error) {
    toastStore.error((error as Error).message || 'Failed to load report')
  } finally {
    loading.value = false
  }
}

const exportCsv = () => {
  if (!report.value) return
  const rows = report.value.summary.map((r) => [
    r.item.name,
    r.item.category,
    r.item.unit,
    r.opening_stock,
    r.stock_in,
    r.usage,
    r.closing_stock,
    r.is_low_stock ? 'YES' : 'NO',
  ])
  const header = ['Item', 'Category', 'Unit', 'Opening', 'Stock In', 'Usage', 'Closing', 'Low Stock']
  const csv =
    [header, ...rows]
      .map((row) => row.map((cell) => `"${String(cell).replace(/"/g, '""')}"`).join(','))
      .join('\n') + '\n'

  const blob = new Blob(['\uFEFF' + csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `inventory-usage-report-${new Date().toISOString().slice(0, 10)}.csv`
  a.click()
  URL.revokeObjectURL(url)
}

const typeLabel = (t: string) =>
  ({ 'stock-in': 'Stock In', usage: 'Usage', adjustment: 'Adjustment' } as Record<string, string>)[t] || t

const typeClass = (t: string) =>
  ({ 'stock-in': 'badge-stock-in', usage: 'badge-usage', adjustment: 'badge-adjust' } as Record<string, string>)[t] || ''

onMounted(() => {
  fetchItems()
  fetchTransactions()
})
</script>

<template>
  <div class="inventory-view">
    <!-- Tabs -->
    <div class="tabs">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        class="tab"
        :class="{ active: activeTab === tab.value }"
        @click="activeTab = tab.value"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- ============ ITEMS ============ -->
    <div v-if="activeTab === 'items'" class="pane">
      <div class="control-bar">
        <div class="low-stock-alert" v-if="lowStockItems.length">
          <strong>{{ lowStockItems.length }}</strong> item(s) below reorder level
        </div>
        <button class="btn btn-primary" @click="openAddItem">+ Add Item</button>
      </div>

      <div class="table-card">
        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>Item</th>
                <th>Category</th>
                <th>Unit</th>
                <th>Current Stock</th>
                <th>Reorder Level</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in items" :key="item.id">
                <td>{{ item.name }}</td>
                <td>{{ item.category }}</td>
                <td>{{ item.unit }}</td>
                <td>
                  <span class="stock-value"
                        :class="{ low: item.reorder_level > 0 && item.current_stock <= item.reorder_level }">
                    {{ item.current_stock }}
                  </span>
                </td>
                <td>{{ item.reorder_level }}</td>
                <td>
                  <span class="badge" :class="item.status === 'active' ? 'badge-active' : 'badge-inactive'">
                    {{ item.status }}
                  </span>
                </td>
                <td>
                  <div class="action-buttons">
                    <button class="btn btn-sm btn-outline" @click="openRecord(item, 'stock-in')">Stock In</button>
                    <button class="btn btn-sm btn-primary-outline" @click="openRecord(item, 'usage')">Usage</button>
                    <button class="btn btn-sm btn-danger-outline" @click="deleteItem(item)">Delete</button>
                  </div>
                </td>
              </tr>
              <tr v-if="loading"><td colspan="7" class="no-data">Loading items...</td></tr>
              <tr v-else-if="items.length === 0"><td colspan="7" class="no-data">No items yet. Add one to get started.</td></tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ============ USAGE LOG (BOOK) ============ -->
    <div v-if="activeTab === 'usage'" class="pane">
      <div class="control-bar">
        <div class="filter-row">
          <label class="filter-label">From</label>
          <input type="date" v-model="txnFrom" class="search-input date-input" />
        </div>
        <div class="filter-row">
          <label class="filter-label">To</label>
          <input type="date" v-model="txnTo" class="search-input date-input" />
        </div>
        <button class="btn btn-outline" @click="fetchTransactions">Filter</button>
      </div>

      <div class="table-card">
        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>Date</th>
                <th>Item</th>
                <th>Type</th>
                <th>Qty</th>
                <th>Balance After</th>
                <th>Reference</th>
                <th>Note</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="tx in transactions" :key="tx.id">
                <td>{{ new Date(tx.created_at).toLocaleString('en-GB') }}</td>
                <td>{{ tx.item?.name || tx.item_id }}</td>
                <td><span class="badge" :class="typeClass(tx.type)">{{ typeLabel(tx.type) }}</span></td>
                <td :class="tx.type === 'usage' ? 'text-usage' : 'text-stock-in'">{{ tx.type === 'usage' ? '-' : '+' }}{{ tx.quantity }}</td>
                <td>{{ tx.balance_after }}</td>
                <td>{{ tx.reference || '-' }}</td>
                <td>{{ tx.note || '-' }}</td>
              </tr>
              <tr v-if="transactions.length === 0"><td colspan="7" class="no-data">No transactions recorded.</td></tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ============ STOCK BALANCE CHECK ============ -->
    <div v-if="activeTab === 'balance'" class="pane">
      <div class="table-card">
        <div class="card-heading">
          <h3 class="card-title">Check Items Balance</h3>
          <p class="card-sub">Count the actual stock and reconcile against the system to verify items are balanced with usage/orders.</p>
        </div>
        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>Item</th>
                <th>System Stock</th>
                <th>Reorder Level</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in items" :key="item.id">
                <td>{{ item.name }} <span class="muted">{{ item.unit }}</span></td>
                <td>{{ item.current_stock }}</td>
                <td>{{ item.reorder_level }}</td>
                <td>
                  <span class="badge" :class="item.current_stock <= 0 ? 'badge-inactive' : 'badge-active'">
                    {{ item.current_stock > 0 ? 'In stock' : 'Out of stock' }}
                  </span>
                </td>
                <td><button class="btn btn-sm btn-primary-outline" @click="openBalance(item)">Count / Reconcile</button></td>
              </tr>
              <tr v-if="items.length === 0"><td colspan="5" class="no-data">No items yet.</td></tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ============ REPORT ============ -->
    <div v-if="activeTab === 'report'" class="pane">
      <div class="control-bar">
        <div class="filter-row">
          <label class="filter-label">From</label>
          <input type="date" v-model="repFrom" class="search-input date-input" />
        </div>
        <div class="filter-row">
          <label class="filter-label">To</label>
          <input type="date" v-model="repTo" class="search-input date-input" />
        </div>
        <button class="btn btn-outline" @click="fetchReport">Generate</button>
        <button class="btn btn-primary" @click="exportCsv" :disabled="!report">Export CSV</button>
      </div>

      <template v-if="report">
        <div v-if="report.low_stock.length" class="low-stock-card">
          <h4>Low Stock Alerts</h4>
          <div class="chip-row">
            <span v-for="row in report.low_stock" :key="row.item.id" class="chip">
              {{ row.item.name }} ({{ row.closing_stock }} {{ row.item.unit }})
            </span>
          </div>
        </div>

        <div class="table-card">
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Item</th>
                  <th>Category</th>
                  <th>Opening</th>
                  <th>Stock In</th>
                  <th>Usage</th>
                  <th>Closing</th>
                  <th>Low</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="row in report.summary" :key="row.item.id">
                  <td>{{ row.item.name }}</td>
                  <td>{{ row.item.category }}</td>
                  <td>{{ row.opening_stock }}</td>
                  <td class="text-stock-in">{{ row.stock_in }}</td>
                  <td class="text-usage">{{ row.usage }}</td>
                  <td class="bold">{{ row.closing_stock }}</td>
                  <td><span v-if="row.is_low_stock" class="badge badge-low">LOW</span><span v-else>-</span></td>
                </tr>
                <tr v-if="report.summary.length === 0"><td colspan="7" class="no-data">No data for the selected period.</td></tr>
              </tbody>
            </table>
          </div>
        </div>
      </template>
      <div v-else class="empty-hint">Select a date range and click Generate.</div>
    </div>

    <!-- Add Item Modal -->
    <div v-if="showItemAdd" class="modal-overlay" @click.self="showItemAdd = false">
      <div class="modal-card modal-card-sm">
        <div class="modal-header">
          <h3 class="modal-title">Add Inventory Item</h3>
          <button class="modal-close" @click="showItemAdd = false" aria-label="Close">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-grid">
            <div class="form-field form-field-full">
              <label class="form-label">Name *</label>
              <input v-model="itemForm.name" class="form-input" placeholder="e.g. Bath soap" />
            </div>
            <div class="form-field">
              <label class="form-label">Category *</label>
              <input v-model="itemForm.category" class="form-input" placeholder="e.g. Amenities" />
            </div>
            <div class="form-field">
              <label class="form-label">Unit *</label>
              <input v-model="itemForm.unit" class="form-input" placeholder="e.g. pcs / kg" />
            </div>
            <div class="form-field">
              <label class="form-label">Opening Stock</label>
              <input v-model.number="itemForm.current_stock" type="number" class="form-input" />
            </div>
            <div class="form-field">
              <label class="form-label">Reorder Level</label>
              <input v-model.number="itemForm.reorder_level" type="number" class="form-input" />
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" @click="showItemAdd = false">Cancel</button>
          <button class="btn btn-primary" @click="submitItem" :disabled="savingItem">
            {{ savingItem ? 'Saving...' : 'Add Item' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Record modal -->
    <div v-if="showRecordModal" class="modal-overlay" @click.self="showRecordModal = false">
      <div class="modal-card modal-card-sm">
        <div class="modal-header">
          <h3 class="modal-title">{{ recordMode === 'stock-in' ? 'Add Stock (In)' : 'Record Usage' }}</h3>
          <button class="modal-close" @click="showRecordModal = false" aria-label="Close">&times;</button>
        </div>
        <div class="modal-body">
          <p v-if="recordItem" class="modal-desc">
            {{ recordItem.name }} — current stock: <strong>{{ recordItem.current_stock }}</strong> {{ recordItem.unit }}
          </p>
          <div class="form-field">
            <label class="form-label">Quantity *</label>
            <input v-model.number="recordForm.quantity" type="number" min="1" class="form-input" />
          </div>
          <div class="form-field">
            <label class="form-label">Reference (e.g. order / reservation code)</label>
            <input v-model="recordForm.reference" class="form-input" placeholder="e.g. Order #123 / Table 4" />
          </div>
          <div class="form-field">
            <label class="form-label">Note</label>
            <input v-model="recordForm.note" class="form-input" placeholder="Reason / purpose" />
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" @click="showRecordModal = false">Cancel</button>
          <button class="btn btn-primary" @click="submitRecord" :disabled="savingRecord">
            {{ savingRecord ? 'Saving...' : recordMode === 'stock-in' ? 'Add Stock' : 'Record Usage' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Balance modal -->
    <div v-if="showBalanceModal" class="modal-overlay" @click.self="showBalanceModal = false">
      <div class="modal-card modal-card-sm">
        <div class="modal-header">
          <h3 class="modal-title">Stock Balance Check</h3>
          <button class="modal-close" @click="showBalanceModal = false" aria-label="Close">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-field">
            <label class="form-label">System Stock</label>
            <input :value="balanceEntry.system_stock" class="form-input" disabled />
          </div>
          <div class="form-field">
            <label class="form-label">Actual Counted Stock *</label>
            <input v-model.number="balanceEntry.actual_stock" type="number" min="0" class="form-input" />
          </div>
          <div class="diff-line" :class="balanceDiff === 0 ? 'diff-zero' : 'diff-var'">
            Variance: {{ balanceDiff > 0 ? '+' : '' }}{{ balanceDiff }}
          </div>
          <div class="form-field">
            <label class="form-label">Note</label>
            <input v-model="balanceEntry.note" class="form-input" placeholder="Anything to note / reason for variance" />
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" @click="showBalanceModal = false">Cancel</button>
          <button class="btn btn-primary" @click="submitBalance" :disabled="savingBalance">
            {{ savingBalance ? 'Saving...' : 'Save & Reconcile' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.inventory-view {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.tabs {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.tab {
  padding: 10px 18px;
  border-radius: 10px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  background: #ffffff;
  color: #64748b;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.tab.active {
  color: #e15b2b;
  background-color: #fff1ea;
  border-color: rgba(225, 91, 43, 0.3);
}

.pane {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.control-bar {
  display: flex;
  gap: 16px;
  background-color: #ffffff;
  padding: 16px;
  border-radius: 16px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  flex-wrap: wrap;
  align-items: center;
  justify-content: flex-end;
}

.filter-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.filter-label {
  font-size: 0.8rem;
  font-weight: 600;
  color: #475569;
}

.date-input { max-width: 170px; }

.low-stock-alert {
  margin-right: auto;
  padding: 8px 14px;
  border-radius: 10px;
  background-color: #fef3c7;
  color: #92400e;
  font-size: 0.85rem;
}

.low-stock-card {
  background-color: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 16px;
  padding: 16px;
}
.low-stock-card h4 { margin: 0 0 10px; color: #92400e; }

.chip-row {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
.chip {
  padding: 4px 10px;
  border-radius: 9999px;
  background-color: #fef3c7;
  color: #92400e;
  font-size: 0.78rem;
  font-weight: 600;
}

.table-card {
  background-color: #ffffff;
  border-radius: 16px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.card-heading { padding: 18px 20px 0; }
.card-title { margin: 0; font-size: 1rem; color: #0f172a; }
.card-sub { margin: 4px 0 0; color: #64748b; font-size: 0.82rem; }

.table-container { overflow-x: auto; }
.data-table { width: 100%; border-collapse: collapse; text-align: left; }
.data-table th {
  padding: 14px 20px; border-bottom: 1px solid #e2e8f0;
  color: #64748b; font-size: 0.75rem; font-weight: 700;
  text-transform: uppercase; letter-spacing: 0.05em;
}
.data-table td {
  padding: 16px 20px; border-bottom: 1px solid #f1f5f9;
  font-size: 0.875rem; color: #334155; vertical-align: middle;
}

.stock-value { font-weight: 700; }
.stock-value.low { color: #dc2626; }
.muted { color: #94a3b8; font-size: 0.78rem; }
.bold { font-weight: 700; }
.text-stock-in { color: #059669; font-weight: 600; }
.text-usage { color: #dc2626; font-weight: 600; }

.badge {
  padding: 4px 10px; border-radius: 9999px;
  font-size: 0.75rem; font-weight: 600; display: inline-block;
}
.badge-active { background-color: #d1fae5; color: #065f46; }
.badge-inactive { background-color: #fee2e2; color: #991b1b; }
.badge-stock-in { background-color: #d1fae5; color: #065f46; }
.badge-usage { background-color: #fee2e2; color: #991b1b; }
.badge-adjust { background-color: #fef3c7; color: #92400e; }
.badge-low { background-color: #fee2e2; color: #b91c1c; }

.action-buttons { display: flex; gap: 8px; flex-wrap: wrap; }

.btn {
  padding: 8px 14px; border-radius: 10px; font-size: 0.85rem;
  font-weight: 600; cursor: pointer; border: 1px solid transparent; transition: all 0.2s;
}
.btn:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-sm { padding: 4px 8px; font-size: 0.75rem; border-radius: 6px; }
.btn-primary { background-color: #e15b2b; color: #ffffff; }
.btn-primary:hover { background-color: #c84e20; }
.btn-primary-outline { border-color: #e15b2b; color: #e15b2b; background: transparent; }
.btn-primary-outline:hover { background-color: #fff1ea; }
.btn-outline { border-color: #cbd5e1; color: #334155; background: #ffffff; }
.btn-outline:hover { background-color: #f8fafc; }
.btn-danger-outline { border-color: #ef4444; color: #ef4444; background: transparent; }
.btn-danger-outline:hover { background-color: #fee2e2; }

.no-data { text-align: center; color: #64748b; padding: 32px !important; }
.empty-hint { color: #64748b; text-align: center; padding: 40px; }

/* Modal */
.modal-overlay {
  position: fixed; inset: 0; background-color: rgba(15, 23, 42, 0.5);
  display: flex; align-items: center; justify-content: center;
  padding: 16px; z-index: 50;
}
.modal-card {
  background-color: #ffffff; border-radius: 16px; width: 100%;
  max-width: 560px; box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2); overflow: hidden;
}
.modal-card-sm { max-width: 420px; }
.modal-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 18px 24px; border-bottom: 1px solid #e2e8f0;
}
.modal-title { font-size: 1.1rem; font-weight: 700; color: #0f172a; margin: 0; }
.modal-close { border: none; background: transparent; font-size: 1.5rem; color: #94a3b8; cursor: pointer; }
.modal-close:hover { color: #334155; }
.modal-body { padding: 24px; display: flex; flex-direction: column; gap: 14px; }
.modal-desc { color: #475569; font-size: 0.875rem; margin: 0; }
.modal-footer {
  display: flex; justify-content: flex-end; gap: 12px; padding: 16px 24px; border-top: 1px solid #e2e8f0;
}
.form-field { display: flex; flex-direction: column; gap: 6px; }
.form-label { font-size: 0.8rem; font-weight: 600; color: #334155; }
.form-input {
  padding: 10px 14px; border-radius: 10px; border: 1px solid #cbd5e1;
  outline: none; font-size: 0.875rem; color: #334155; background: #fff; box-sizing: border-box; width: 100%;
}
.form-input:focus { border-color: #e15b2b; box-shadow: 0 0 0 2px rgba(225, 91, 43, 0.1); }
.form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.form-field-full { grid-column: span 2; }
.diff-line { padding: 10px 14px; border-radius: 10px; font-size: 0.875rem; font-weight: 700; }
.diff-zero { background-color: #d1fae5; color: #065f46; }
.diff-var { background-color: #fef3c7; color: #92400e; }
</style>