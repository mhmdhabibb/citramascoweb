<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { inventoryService } from '@/services/admin/inventoryService'
import { useToastStore } from '@/stores/toastStore'
import type { InventoryItem, InventoryTransaction, Reservation } from '@/types'
import { reservationService } from '@/services/admin/reservationService'

const toastStore = useToastStore()
const loading = ref(false)

const items = ref<InventoryItem[]>([])
const transactions = ref<InventoryTransaction[]>([])
const reservations = ref<Reservation[]>([])
const savingRecord = ref(false)

const recordForm = ref({
  item_id: '',
  quantity: 1,
  usage_type: 'guest',
  guest_name: '',
  room_name: '',
  reference_other: '',
  note: ''
})

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

const fetchTransactions = async () => {
  try {
    const data = await inventoryService.getTransactions()
    // only show usage
    transactions.value = data.filter(t => t.type === 'usage')
  } catch (error) {
    console.error('Failed to load transactions', error)
  }
}

const fetchReservations = async () => {
  try {
    const data = await reservationService.getAll()
    // Only show active/checked in reservations
    reservations.value = data.filter(r => r.status === 'Checked In' || r.status === 'Booked')
  } catch (error) {
    console.error('Failed to load reservations', error)
  }
}

onMounted(() => {
  fetchItems()
  fetchTransactions()
  fetchReservations()
})

const submitRecord = async () => {
  if (!recordForm.value.item_id) {
    toastStore.error('Please select an item')
    return
  }
  const qty = Number(recordForm.value.quantity)
  if (!qty || qty <= 0) {
    toastStore.error('Quantity must be greater than zero')
    return
  }
  
  try {
    savingRecord.value = true
    
    let finalReference = ''
    if (recordForm.value.usage_type === 'guest') {
      const room = recordForm.value.room_name ? `Room ${recordForm.value.room_name}` : ''
      const guest = recordForm.value.guest_name || ''
      finalReference = [room, guest].filter(Boolean).join(' - ')
    } else {
      finalReference = recordForm.value.reference_other
    }

    const msg = await inventoryService.usage(recordForm.value.item_id, {
      quantity: qty,
      reference: finalReference,
      note: recordForm.value.note,
    })
    
    toastStore.success(msg || 'Recorded successfully')
    
    // reset form
    recordForm.value = {
      item_id: '',
      quantity: 1,
      usage_type: 'guest',
      guest_name: '',
      room_name: '',
      reference_other: '',
      note: ''
    }
    
    await fetchItems()
    await fetchTransactions()
  } catch (error) {
    toastStore.error((error as Error).message || 'Failed to record')
  } finally {
    savingRecord.value = false
  }
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric', month: 'short', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  })
}

const selectedItem = computed(() => {
  return items.value.find(i => i.id === recordForm.value.item_id)
})

const uniqueRooms = computed(() => {
  const rooms = reservations.value.map(r => r.room?.name).filter(Boolean) as string[]
  return [...new Set(rooms)]
})

const onGuestChange = () => {
  const selectedRes = reservations.value.find(r => r.full_name === recordForm.value.guest_name)
  if (selectedRes && selectedRes.room?.name) {
    recordForm.value.room_name = selectedRes.room.name
  }
}
</script>

<template>
  <div class="inventory-usage-view">
    <div class="page-header">
      <div>
        <h1 class="page-title">Add Inventory Usage</h1>
        <p class="page-subtitle">Record items consumed by guests or operational needs.</p>
      </div>
    </div>

    <div class="grid-container">
      <!-- Form Side -->
      <div class="form-section card">
        <div class="card-header">
          <h2 class="card-title">Usage Form</h2>
        </div>
        <div class="card-body">
          <div class="form-field">
            <label class="form-label">Select Item</label>
            <select v-model="recordForm.item_id" class="form-input">
              <option value="" disabled>-- Select Inventory Item --</option>
              <option v-for="item in items" :key="item.id" :value="item.id">
                {{ item.name }} (Stock: {{ item.current_stock }} {{ item.unit }})
              </option>
            </select>
            <span v-if="selectedItem && selectedItem.current_stock <= 0" class="error-text">Item is out of stock!</span>
          </div>

          <div class="form-grid">
            <div class="form-field">
              <label class="form-label">Quantity Used</label>
              <input type="number" v-model="recordForm.quantity" min="1" class="form-input" />
            </div>
            <div class="form-field full-width">
              <label class="form-label">Usage Type</label>
              <div class="usage-type-toggle">
                <label class="radio-label">
                  <input type="radio" v-model="recordForm.usage_type" value="guest" /> Guest Usage
                </label>
                <label class="radio-label">
                  <input type="radio" v-model="recordForm.usage_type" value="operational" /> Operational Usage
                </label>
              </div>
            </div>

            <template v-if="recordForm.usage_type === 'guest'">
              <div class="form-field">
                <label class="form-label">Guest Name</label>
                <select v-model="recordForm.guest_name" class="form-input" @change="onGuestChange">
                  <option value="">-- Select Guest --</option>
                  <option v-for="res in reservations" :key="res.id" :value="res.full_name">
                    {{ res.full_name }}
                  </option>
                </select>
              </div>
              <div class="form-field">
                <label class="form-label">Room</label>
                <select v-model="recordForm.room_name" class="form-input">
                  <option value="">-- Select Room --</option>
                  <!-- Getting unique rooms from reservations -->
                  <option v-for="room in uniqueRooms" :key="room" :value="room">
                    {{ room }}
                  </option>
                </select>
              </div>
            </template>

            <template v-else>
              <div class="form-field full-width">
                <label class="form-label">Operational Reference</label>
                <select v-model="recordForm.reference_other" class="form-input">
                  <option value="">-- Select Reference --</option>
                  <option value="Operational - Cleaning">Cleaning</option>
                  <option value="Operational - Maintenance">Maintenance</option>
                  <option value="Operational - Office">Office</option>
                  <option value="Other">Other</option>
                </select>
              </div>
            </template>
          </div>

          <div class="form-field">
            <label class="form-label">Note / Reason</label>
            <textarea v-model="recordForm.note" placeholder="Any additional notes..." class="form-input" rows="3"></textarea>
          </div>

          <div class="form-actions">
            <button 
              @click="submitRecord" 
              class="btn btn-primary" 
              :disabled="savingRecord || !recordForm.item_id || (selectedItem && selectedItem.current_stock < recordForm.quantity)"
            >
              <span v-if="savingRecord">Saving...</span>
              <span v-else>Record Usage</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Recent History Side -->
      <div class="history-section card">
        <div class="card-header">
          <h2 class="card-title">Recent Usage Logs</h2>
        </div>
        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>Date</th>
                <th>Item</th>
                <th>Qty</th>
                <th>Guest/Ref</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="t in transactions.slice(0, 10)" :key="t.id">
                <td class="date-col">{{ formatDate(t.created_at) }}</td>
                <td>{{ t.item?.name || 'Unknown' }}</td>
                <td class="qty-col negative">-{{ t.quantity }}</td>
                <td>{{ t.reference || '-' }}</td>
              </tr>
              <tr v-if="transactions.length === 0">
                <td colspan="4" class="empty-state">No usage recorded yet.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.inventory-usage-view {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header {
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

.grid-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

@media (max-width: 1024px) {
  .grid-container {
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
}

.card-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-dark);
  margin: 0;
}

.card-body {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-field { display: flex; flex-direction: column; gap: 6px; }
.form-label { font-size: 0.875rem; font-weight: 600; color: #334155; }
.form-input {
  padding: 12px 16px; border-radius: 10px; border: 1px solid #cbd5e1;
  outline: none; font-size: 0.875rem; color: #334155; background: #fff; box-sizing: border-box; width: 100%;
}
.form-input:focus { border-color: #e15b2b; box-shadow: 0 0 0 2px rgba(225, 91, 43, 0.1); }
.form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.full-width { grid-column: 1 / -1; }

.usage-type-toggle {
  display: flex;
  gap: 16px;
  padding: 8px 0;
}
.radio-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.875rem;
  color: var(--text-dark);
  cursor: pointer;
}

.error-text {
  color: #ef4444;
  font-size: 0.75rem;
  font-weight: 600;
}

.form-actions {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.btn {
  padding: 10px 20px; border-radius: 10px; font-size: 0.875rem;
  font-weight: 600; cursor: pointer; border: 1px solid transparent; transition: all 0.2s;
}
.btn:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-primary { background-color: #e15b2b; color: #ffffff; }
.btn-primary:hover:not(:disabled) { background-color: #c84e20; }

.table-container { overflow-x: auto; }
.data-table { width: 100%; border-collapse: collapse; text-align: left; }
.data-table th {
  padding: 12px 24px; background: white; font-size: 0.75rem; font-weight: 600;
  color: var(--text-muted); text-transform: uppercase; border-bottom: 1px solid var(--border-light);
}
.data-table td { padding: 16px 24px; font-size: 0.875rem; color: var(--text-dark); border-bottom: 1px solid var(--border-light); }
.date-col { color: var(--text-muted) !important; font-size: 0.8125rem !important; white-space: nowrap; }
.qty-col { font-weight: 600; }
.qty-col.negative { color: #ef4444; }
.empty-state { text-align: center; padding: 32px !important; color: var(--text-muted) !important; font-style: italic; }
</style>
