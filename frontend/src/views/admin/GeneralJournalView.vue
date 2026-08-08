<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { financeService, type GeneralJournal, type ChartOfAccount } from '@/services/admin/financeService'
import { useToastStore } from '@/stores/toastStore'

const toastStore = useToastStore()
const journals = ref<GeneralJournal[]>([])
const coas = ref<ChartOfAccount[]>([])
const loading = ref(true)

const showModal = ref(false)
const form = ref({
  date: new Date().toISOString().split('T')[0] as string,
  description: '',
  account_code: '',
  debit: 0,
  credit: 0
})

async function fetchData() {
  try {
    loading.value = true
    const [js, allCoas] = await Promise.all([
      financeService.getGeneralJournals(),
      financeService.getCOAs()
    ])
    journals.value = js
    coas.value = allCoas
  } catch (error: any) {
    toastStore.error(error.message || 'Failed to fetch data')
  } finally {
    loading.value = false
  }
}

function formatIDR(n: number) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(n || 0)
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'long',
    year: 'numeric'
  })
}

const totalDebit = computed(() => journals.value.reduce((sum, j) => sum + (j.debit || 0), 0))
const totalCredit = computed(() => journals.value.reduce((sum, j) => sum + (j.credit || 0), 0))

function openModal() {
  form.value = {
    date: new Date().toISOString().split('T')[0] as string,
    description: '',
    account_code: coas.value.length > 0 ? (coas.value[0]?.code || '') : '',
    debit: 0,
    credit: 0
  }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

async function saveJournal() {
  try {
    await financeService.createGeneralJournal({
      ...form.value,
      debit: Number(form.value.debit),
      credit: Number(form.value.credit)
    })
    toastStore.success('Journal entry saved')
    closeModal()
    fetchData()
  } catch (error: any) {
    toastStore.error(error.message || 'Failed to save journal')
  }
}

onMounted(() => {
  fetchData()
})
</script>

<template>
  <div class="journal-page">
      <div class="header">
        <div>
          <h1>Jurnal Umum</h1>
          <p class="subtitle">General Journal and Ledger Entries.</p>
        </div>
        <button class="btn-primary" @click="openModal()">+ Add Entry</button>
      </div>

      <div class="stats-grid">
        <div class="stat-card">
          <h3>Total Debit</h3>
          <h2 class="text-blue">{{ formatIDR(totalDebit) }}</h2>
        </div>
        <div class="stat-card">
          <h3>Total Credit</h3>
          <h2 class="text-blue">{{ formatIDR(totalCredit) }}</h2>
        </div>
        <div class="stat-card">
          <h3>Balance</h3>
          <h2 :class="totalDebit === totalCredit ? 'text-green' : 'text-red'">
            {{ totalDebit === totalCredit ? 'Balanced' : 'Unbalanced' }}
          </h2>
        </div>
      </div>

      <div class="table-container">
        <table>
          <thead>
            <tr>
              <th>Date</th>
              <th>Description</th>
              <th>Account</th>
              <th>Ref</th>
              <th class="text-right">Debit</th>
              <th class="text-right">Credit</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="6" class="text-center">Loading...</td>
            </tr>
            <tr v-else-if="journals.length === 0">
              <td colspan="6" class="text-center">No journal entries found.</td>
            </tr>
            <tr v-else v-for="j in journals" :key="j.id">
              <td>{{ formatDate(j.date) }}</td>
              <td>{{ j.description }}</td>
              <td>{{ j.account?.name || j.account_code }} <span class="text-xs text-gray-500">({{ j.account_code }})</span></td>
              <td>{{ j.reservation_id ? 'INV-' + j.reservation_id.slice(-5).toUpperCase() : '-' }}</td>
              <td class="text-right">{{ j.debit > 0 ? formatIDR(j.debit) : '-' }}</td>
              <td class="text-right">{{ j.credit > 0 ? formatIDR(j.credit) : '-' }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Modal -->
      <div class="modal-overlay" v-if="showModal" @click.self="closeModal">
        <div class="modal">
          <h2>New Journal Entry</h2>
          <form @submit.prevent="saveJournal">
            <div class="form-group">
              <label>Date</label>
              <input type="date" v-model="form.date" required />
            </div>
            <div class="form-group">
              <label>Account</label>
              <select v-model="form.account_code" required>
                <option v-for="c in coas" :key="c.code" :value="c.code">{{ c.code }} - {{ c.name }}</option>
              </select>
            </div>
            <div class="form-group">
              <label>Description</label>
              <input type="text" v-model="form.description" required />
            </div>
            <div class="split-group">
              <div class="form-group w-full">
                <label>Debit (IDR)</label>
                <input type="number" v-model="form.debit" min="0" />
              </div>
              <div class="form-group w-full">
                <label>Credit (IDR)</label>
                <input type="number" v-model="form.credit" min="0" />
              </div>
            </div>
            <p class="help-text">Please ensure you enter either debit or credit, not both in a single line if it's a simple entry.</p>
            <div class="modal-actions">
              <button type="button" class="btn-secondary" @click="closeModal">Cancel</button>
              <button type="submit" class="btn-primary">Save Entry</button>
            </div>
          </form>
        </div>
      </div>
  </div>
</template>

<style scoped>
.journal-page {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header h1 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}
.subtitle {
  color: #64748b;
  margin: 4px 0 0 0;
}
.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}
.stat-card {
  background: white;
  padding: 20px;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
}
.stat-card h3 {
  margin: 0 0 8px 0;
  color: #64748b;
  font-size: 0.85rem;
  text-transform: uppercase;
}
.stat-card h2 {
  margin: 0;
  font-size: 1.8rem;
  font-weight: 800;
}
.text-blue { color: #0f172a; }
.text-green { color: #16a34a; }
.text-red { color: #dc2626; }
.text-xs { font-size: 0.75rem; }
.text-gray-500 { color: #64748b; }
.text-right { text-align: right !important; }

.table-container {
  background: white;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  overflow: hidden;
}
table {
  width: 100%;
  border-collapse: collapse;
}
th {
  background: #f8fafc;
  padding: 12px 16px;
  text-align: left;
  font-size: 0.85rem;
  color: #64748b;
  text-transform: uppercase;
  border-bottom: 1px solid #e2e8f0;
}
td {
  padding: 16px;
  border-bottom: 1px solid #f1f5f9;
  color: #334155;
}
.text-center {
  text-align: center;
  padding: 32px !important;
  color: #94a3b8;
}
.btn-primary {
  background: #e4793b;
  color: white;
  padding: 10px 18px;
  border-radius: 8px;
  border: none;
  cursor: pointer;
  font-weight: 600;
}
.btn-secondary {
  background: #f1f5f9;
  color: #475569;
  padding: 10px 18px;
  border-radius: 8px;
  border: none;
  cursor: pointer;
  font-weight: 600;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 50;
}
.modal {
  background: white;
  width: 500px;
  border-radius: 12px;
  padding: 24px;
}
.modal h2 {
  margin: 0 0 20px 0;
}
.form-group {
  margin-bottom: 16px;
}
.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
}
.form-group input, .form-group select {
  width: 100%;
  padding: 10px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-family: inherit;
}
.split-group {
  display: flex;
  gap: 16px;
}
.w-full {
  width: 100%;
}
.help-text {
  font-size: 0.8rem;
  color: #64748b;
  margin: 0;
}
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}
</style>
