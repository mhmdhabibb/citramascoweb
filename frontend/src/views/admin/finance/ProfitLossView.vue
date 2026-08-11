<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { financeService, type GeneralJournal, type ChartOfAccount } from '@/services/admin/financeService'
import { useToastStore } from '@/stores/toastStore'

const toastStore = useToastStore()
const journals = ref<GeneralJournal[]>([])
const coas = ref<ChartOfAccount[]>([])
const loading = ref(true)

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

// Calculate balances for each account
const accountBalances = computed(() => {
  const balances: Record<string, number> = {}
  
  coas.value.forEach(c => {
    balances[c.code] = 0
  })

  journals.value.forEach(j => {
    if (!balances[j.account_code]) balances[j.account_code] = 0
    const accType = coas.value.find(c => c.code === j.account_code)?.type

    if (accType === 'Asset' || accType === 'Expense') {
      balances[j.account_code] = (balances[j.account_code] || 0) + (j.debit || 0) - (j.credit || 0)
    } else {
      balances[j.account_code] = (balances[j.account_code] || 0) + (j.credit || 0) - (j.debit || 0)
    }
  })

  return balances
})

const revenues = computed(() => {
  return coas.value
    .filter(c => c.type === 'Revenue')
    .map(c => ({
      ...c,
      balance: accountBalances.value[c.code] || 0
    }))
    .filter(c => c.balance !== 0) // Hide zero balance accounts
})

const expenses = computed(() => {
  return coas.value
    .filter(c => c.type === 'Expense')
    .map(c => ({
      ...c,
      balance: accountBalances.value[c.code] || 0
    }))
    .filter(c => c.balance !== 0)
})

const totalRevenue = computed(() => revenues.value.reduce((sum, c) => sum + c.balance, 0))
const totalExpense = computed(() => expenses.value.reduce((sum, c) => sum + c.balance, 0))
const netIncome = computed(() => totalRevenue.value - totalExpense.value)

function formatIDR(n: number) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(n || 0)
}

onMounted(() => {
  fetchData()
})
</script>

<template>
  <div class="finance-page">
    <div class="header">
      <div>
        <h1>Laporan Laba Rugi</h1>
        <p class="subtitle">Profit & Loss / Income Statement</p>
      </div>
      <button class="btn-secondary" @click="fetchData">Refresh Data</button>
    </div>

    <div v-if="loading" class="text-center p-8 text-gray-500">
      Menghitung laporan laba rugi...
    </div>
    
    <div v-else class="report-container">
      <!-- REVENUES -->
      <div class="report-section">
        <h2 class="section-title">Pendapatan (Revenues)</h2>
        <table class="report-table">
          <tbody>
            <tr v-if="revenues.length === 0">
              <td class="text-gray-400 italic">Belum ada data pendapatan.</td>
              <td class="text-right">-</td>
            </tr>
            <tr v-for="rev in revenues" :key="rev.code">
              <td>{{ rev.name }} <span class="text-xs text-gray-400">({{ rev.code }})</span></td>
              <td class="text-right">{{ formatIDR(rev.balance) }}</td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="total-row">
              <td>Total Pendapatan</td>
              <td class="text-right font-bold">{{ formatIDR(totalRevenue) }}</td>
            </tr>
          </tfoot>
        </table>
      </div>

      <!-- EXPENSES -->
      <div class="report-section mt-8">
        <h2 class="section-title">Biaya (Expenses)</h2>
        <table class="report-table">
          <tbody>
            <tr v-if="expenses.length === 0">
              <td class="text-gray-400 italic">Belum ada data biaya.</td>
              <td class="text-right">-</td>
            </tr>
            <tr v-for="exp in expenses" :key="exp.code">
              <td>{{ exp.name }} <span class="text-xs text-gray-400">({{ exp.code }})</span></td>
              <td class="text-right">{{ formatIDR(exp.balance) }}</td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="total-row">
              <td>Total Biaya</td>
              <td class="text-right font-bold">{{ formatIDR(totalExpense) }}</td>
            </tr>
          </tfoot>
        </table>
      </div>

      <!-- NET INCOME -->
      <div class="net-income-section mt-8" :class="netIncome >= 0 ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'">
        <h2>{{ netIncome >= 0 ? 'Laba Bersih (Net Income)' : 'Rugi Bersih (Net Loss)' }}</h2>
        <h1 :class="netIncome >= 0 ? 'text-green-700' : 'text-red-700'">{{ formatIDR(netIncome) }}</h1>
      </div>
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
.btn-secondary {
  background: #f1f5f9;
  color: #475569;
  padding: 10px 18px;
  border-radius: 8px;
  border: 1px solid #cbd5e1;
  cursor: pointer;
  font-weight: 600;
}

.report-container {
  background: white;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  padding: 32px;
  max-width: 800px;
  margin: 0 auto;
  width: 100%;
}
.report-section {
  margin-bottom: 24px;
}
.section-title {
  font-size: 1.1rem;
  color: #334155;
  border-bottom: 2px solid #e2e8f0;
  padding-bottom: 8px;
  margin-bottom: 12px;
}
.report-table {
  width: 100%;
  border-collapse: collapse;
}
.report-table td {
  padding: 8px 0;
  color: #475569;
}
.report-table tfoot td {
  padding-top: 16px;
  border-top: 1px dashed #cbd5e1;
  color: #0f172a;
}
.total-row {
  font-weight: 700;
}

.net-income-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px;
  border-radius: 8px;
  border: 2px solid;
}
.net-income-section h2 {
  margin: 0;
  font-size: 1.2rem;
  color: #334155;
}
.net-income-section h1 {
  margin: 0;
  font-size: 2rem;
  font-weight: 800;
}

.mt-8 { margin-top: 32px; }
.text-right { text-align: right; }
.text-center { text-align: center; }
.p-8 { padding: 32px; }
.text-xs { font-size: 0.75rem; }
.text-gray-400 { color: #94a3b8; }
.text-gray-500 { color: #64748b; }
.italic { font-style: italic; }
.font-bold { font-weight: 700; }
.bg-green-50 { background-color: #f0fdf4; }
.border-green-200 { border-color: #bbf7d0; }
.text-green-700 { color: #15803d; }
.bg-red-50 { background-color: #fef2f2; }
.border-red-200 { border-color: #fecaca; }
.text-red-700 { color: #b91c1c; }
</style>
