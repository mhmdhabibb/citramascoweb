<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { financeService, type GeneralJournal, type ChartOfAccount } from '@/services/admin/financeService'
import { useToastStore } from '@/stores/toastStore'

const toastStore = useToastStore()
const journals = ref<GeneralJournal[]>([])
const coas = ref<ChartOfAccount[]>([])
const loading = ref(true)
const selectedAccountId = ref('')

async function fetchData() {
  try {
    loading.value = true
    const [js, allCoas] = await Promise.all([
      financeService.getGeneralJournals(),
      financeService.getCOAs()
    ])
    journals.value = js.sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime())
    coas.value = allCoas
    if (coas.value.length > 0) {
      selectedAccountId.value = coas.value[0]?.code || ''
    }
  } catch (error: any) {
    toastStore.error(error.message || 'Failed to fetch data')
  } finally {
    loading.value = false
  }
}

const selectedAccount = computed(() => {
  return coas.value.find(c => c.code === selectedAccountId.value)
})

const ledgerEntries = computed(() => {
  if (!selectedAccount.value) return []
  const accType = selectedAccount.value.type

  // Filter journals for the selected account
  const entries = journals.value.filter(j => j.account_code === selectedAccountId.value)

  let runningBalance = 0
  return entries.map(j => {
    // Determine how balance is affected based on normal account balance rules
    // Asset & Expense -> Normal Debit
    // Liability, Equity, Revenue -> Normal Credit
    if (accType === 'Asset' || accType === 'Expense') {
      runningBalance += (j.debit || 0) - (j.credit || 0)
    } else {
      runningBalance += (j.credit || 0) - (j.debit || 0)
    }

    return {
      ...j,
      runningBalance
    }
  })
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
    month: 'short',
    year: 'numeric'
  })
}

onMounted(() => {
  fetchData()
})
</script>

<template>
  <div class="finance-page">
    <div class="header">
      <div>
        <h1>Buku Besar (General Ledger)</h1>
        <p class="subtitle">Laporan mutasi tiap akun secara detail</p>
      </div>
      <div class="filters">
        <label class="font-bold">Select Account:</label>
        <select v-model="selectedAccountId" class="account-select">
          <option v-for="c in coas" :key="c.code" :value="c.code">{{ c.code }} - {{ c.name }}</option>
        </select>
      </div>
    </div>

    <div class="stat-card" v-if="selectedAccount">
      <div>
        <h3>{{ selectedAccount.name }} ({{ selectedAccount.code }})</h3>
        <p class="text-sm text-gray-500">Account Type: {{ selectedAccount.type }}</p>
      </div>
      <div class="text-right">
        <h3>Current Balance</h3>
        <h2 class="text-blue">
          {{ ledgerEntries.length > 0 ? formatIDR(ledgerEntries[ledgerEntries.length - 1]?.runningBalance || 0) : formatIDR(0) }}
        </h2>
      </div>
    </div>

    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>Date</th>
            <th>Description</th>
            <th>Ref</th>
            <th class="text-right">Debit</th>
            <th class="text-right">Credit</th>
            <th class="text-right bg-blue-50">Balance</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td colspan="6" class="text-center">Loading...</td>
          </tr>
          <tr v-else-if="ledgerEntries.length === 0">
            <td colspan="6" class="text-center">No transactions found for this account.</td>
          </tr>
          <tr v-else v-for="entry in ledgerEntries" :key="entry.id">
            <td>{{ formatDate(entry.date) }}</td>
            <td>{{ entry.description }}</td>
            <td>{{ entry.reservation_id ? 'INV-' + entry.reservation_id.slice(-5).toUpperCase() : '-' }}</td>
            <td class="text-right">{{ (entry.debit || 0) > 0 ? formatIDR(entry.debit) : '-' }}</td>
            <td class="text-right">{{ (entry.credit || 0) > 0 ? formatIDR(entry.credit) : '-' }}</td>
            <td class="text-right bg-blue-50 font-bold">{{ formatIDR(entry.runningBalance) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.finance-page {
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
  padding: 12px 16px;
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
.stat-card {
  background: white;
  padding: 24px;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.stat-card h3 {
  margin: 0 0 8px 0;
  color: #64748b;
  font-size: 0.9rem;
  text-transform: uppercase;
}
.stat-card h2 {
  margin: 0;
  font-size: 2rem;
  font-weight: 800;
}
.text-blue { color: #0f172a; }
.text-gray-500 { color: #64748b; }
.text-sm { font-size: 0.875rem; }
.font-bold { font-weight: 600; }
.text-right { text-align: right !important; }
.text-center { text-align: center !important; }
.bg-blue-50 { background-color: #f8fafc; }

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
  background: #f1f5f9;
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
</style>
