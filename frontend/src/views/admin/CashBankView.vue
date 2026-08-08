<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { financeService, type CashTransaction, type ChartOfAccount } from '@/services/admin/financeService'
import { useToastStore } from '@/stores/toastStore'

const toastStore = useToastStore()
const transactions = ref<CashTransaction[]>([])
const coas = ref<ChartOfAccount[]>([])
const loading = ref(true)

const selectedAccountId = ref('')

const showModal = ref(false)
const form = ref({
  date: new Date().toISOString().split('T')[0] as string,
  description: '',
  type: 'out',
  amount: 0,
  account_id: ''
})

async function fetchData() {
  try {
    loading.value = true
    const [txs, allCoas] = await Promise.all([
      financeService.getCashTransactions(),
      financeService.getCOAs()
    ])
    transactions.value = txs
    // Only asset accounts for Cash/Bank
    coas.value = allCoas.filter(c => c.type === 'Asset')
    if (coas.value.length > 0) {
      selectedAccountId.value = coas.value[0]?.code || ''
    }
  } catch (error: any) {
    toastStore.error(error.message || 'Failed to fetch data')
  } finally {
    loading.value = false
  }
}

const filteredTransactions = computed(() => {
  return transactions.value.filter(t => t.account_id === selectedAccountId.value)
})

const totalBalance = computed(() => {
  return filteredTransactions.value.reduce((acc, t) => {
    return t.type === 'in' ? acc + (t.amount || 0) : acc - (t.amount || 0)
  }, 0)
})

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

function openModal() {
  form.value = {
    date: new Date().toISOString().split('T')[0] as string,
    description: '',
    type: 'out',
    amount: 0,
    account_id: selectedAccountId.value
  }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

async function saveTransaction() {
  try {
    await financeService.createCashTransaction({
      ...form.value,
      amount: Number(form.value.amount)
    })
    toastStore.success('Transaction saved')
    closeModal()
    fetchData()
  } catch (error: any) {
    toastStore.error(error.message || 'Failed to save transaction')
  }
}

async function deleteTransaction(id?: string) {
  if (!id) return
  if (confirm('Are you sure you want to delete this transaction?')) {
    try {
      await financeService.deleteCashTransaction(id)
      toastStore.success('Transaction deleted')
      fetchData()
    } catch (error: any) {
      toastStore.error(error.message || 'Failed to delete transaction')
    }
  }
}

onMounted(() => {
  fetchData()
})
</script>

<template>
  <div class="cash-page">
      <div class="header">
        <div>
          <h1>Kas & Bank</h1>
          <p class="subtitle">Manage cash flows and petty cash.</p>
        </div>
        <button class="btn-primary" @click="openModal()">+ Record Transaction</button>
      </div>

      <div class="filters">
        <label class="font-bold">Select Account:</label>
        <select v-model="selectedAccountId" class="account-select">
          <option v-for="c in coas" :key="c.code" :value="c.code">{{ c.code }} - {{ c.name }}</option>
        </select>
      </div>

      <div class="balance-card">
        <h3>Current Balance</h3>
        <h2 :class="totalBalance >= 0 ? 'text-green' : 'text-red'">{{ formatIDR(totalBalance) }}</h2>
      </div>

      <div class="table-container">
        <table>
          <thead>
            <tr>
              <th>Date</th>
              <th>Description</th>
              <th>In</th>
              <th>Out</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="5" class="text-center">Loading...</td>
            </tr>
            <tr v-else-if="filteredTransactions.length === 0">
              <td colspan="5" class="text-center">No transactions found.</td>
            </tr>
            <tr v-else v-for="t in filteredTransactions" :key="t.id">
              <td>{{ formatDate(t.date) }}</td>
              <td>{{ t.description }}</td>
              <td class="text-green font-bold">{{ t.type === 'in' ? formatIDR(t.amount) : '-' }}</td>
              <td class="text-red font-bold">{{ t.type === 'out' ? formatIDR(t.amount) : '-' }}</td>
              <td>
                <button class="btn-text text-danger" @click="deleteTransaction(t.id)">Delete</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Modal -->
      <div class="modal-overlay" v-if="showModal" @click.self="closeModal">
        <div class="modal">
          <h2>New Cash Transaction</h2>
          <form @submit.prevent="saveTransaction">
            <div class="form-group">
              <label>Account</label>
              <select v-model="form.account_id" required>
                <option v-for="c in coas" :key="c.code" :value="c.code">{{ c.name }}</option>
              </select>
            </div>
            <div class="form-group">
              <label>Date</label>
              <input type="date" v-model="form.date" required />
            </div>
            <div class="form-group">
              <label>Type</label>
              <div class="radio-group">
                <label><input type="radio" v-model="form.type" value="in" /> In (Pemasukan)</label>
                <label><input type="radio" v-model="form.type" value="out" /> Out (Pengeluaran)</label>
              </div>
            </div>
            <div class="form-group">
              <label>Amount (IDR)</label>
              <input type="number" v-model="form.amount" required min="1" />
            </div>
            <div class="form-group">
              <label>Description</label>
              <textarea v-model="form.description" required rows="3"></textarea>
            </div>
            <div class="modal-actions">
              <button type="button" class="btn-secondary" @click="closeModal">Cancel</button>
              <button type="submit" class="btn-primary">Save</button>
            </div>
          </form>
        </div>
      </div>
  </div>
</template>

<style scoped>
.cash-page {
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
.filters {
  display: flex;
  align-items: center;
  gap: 12px;
  background: white;
  padding: 16px;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
}
.account-select {
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 1rem;
  min-width: 250px;
}
.balance-card {
  background: white;
  padding: 24px;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  text-align: center;
}
.balance-card h3 {
  margin: 0 0 8px 0;
  color: #64748b;
  font-size: 0.9rem;
  text-transform: uppercase;
}
.balance-card h2 {
  margin: 0;
  font-size: 2.5rem;
  font-weight: 800;
}
.text-green { color: #16a34a; }
.text-red { color: #dc2626; }
.font-bold { font-weight: 600; }
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
.btn-text {
  background: none;
  border: none;
  font-weight: 600;
  cursor: pointer;
}
.text-danger { color: #ef4444; }

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
  width: 450px;
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
.form-group input, .form-group select, .form-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-family: inherit;
}
.radio-group {
  display: flex;
  gap: 16px;
}
.radio-group label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 400;
}
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}
</style>
