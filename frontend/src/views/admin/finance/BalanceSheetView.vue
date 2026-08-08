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

const assets = computed(() => {
  return coas.value
    .filter(c => c.type === 'Asset')
    .map(c => ({
      ...c,
      balance: accountBalances.value[c.code] || 0
    }))
    .filter(c => c.balance !== 0)
})

const liabilities = computed(() => {
  return coas.value
    .filter(c => c.type === 'Liability')
    .map(c => ({
      ...c,
      balance: accountBalances.value[c.code] || 0
    }))
    .filter(c => c.balance !== 0)
})

const equities = computed(() => {
  return coas.value
    .filter(c => c.type === 'Equity')
    .map(c => ({
      ...c,
      balance: accountBalances.value[c.code] || 0
    }))
    .filter(c => c.balance !== 0)
})

// Net Income for Retained Earnings (Laba Berjalan)
const netIncome = computed(() => {
  const revenues = coas.value.filter(c => c.type === 'Revenue').reduce((sum, c) => sum + (accountBalances.value[c.code] || 0), 0)
  const expenses = coas.value.filter(c => c.type === 'Expense').reduce((sum, c) => sum + (accountBalances.value[c.code] || 0), 0)
  return revenues - expenses
})

const totalAssets = computed(() => assets.value.reduce((sum, c) => sum + c.balance, 0))
const totalLiabilities = computed(() => liabilities.value.reduce((sum, c) => sum + c.balance, 0))
const totalEquity = computed(() => equities.value.reduce((sum, c) => sum + c.balance, 0) + netIncome.value)

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
        <h1>Neraca & Arus Kas</h1>
        <p class="subtitle">Balance Sheet</p>
      </div>
      <button class="btn-secondary" @click="fetchData">Refresh Data</button>
    </div>

    <div v-if="loading" class="text-center p-8 text-gray-500">
      Menghitung laporan neraca...
    </div>
    
    <div v-else class="report-container">
      
      <!-- ASSETS -->
      <div class="report-section">
        <h2 class="section-title">Aktiva (Assets)</h2>
        <table class="report-table">
          <tbody>
            <tr v-if="assets.length === 0">
              <td class="text-gray-400 italic">Belum ada data aktiva.</td>
              <td class="text-right">-</td>
            </tr>
            <tr v-for="asset in assets" :key="asset.code">
              <td>{{ asset.name }} <span class="text-xs text-gray-400">({{ asset.code }})</span></td>
              <td class="text-right">{{ formatIDR(asset.balance) }}</td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="total-row bg-blue-50">
              <td>Total Aktiva</td>
              <td class="text-right font-bold text-blue-800">{{ formatIDR(totalAssets) }}</td>
            </tr>
          </tfoot>
        </table>
      </div>

      <!-- LIABILITIES & EQUITY -->
      <div class="report-section mt-8">
        <h2 class="section-title">Kewajiban & Ekuitas (Liabilities & Equity)</h2>
        
        <h3 class="subsection-title">Kewajiban (Liabilities)</h3>
        <table class="report-table mb-4">
          <tbody>
            <tr v-if="liabilities.length === 0">
              <td class="text-gray-400 italic">Belum ada data kewajiban.</td>
              <td class="text-right">-</td>
            </tr>
            <tr v-for="liab in liabilities" :key="liab.code">
              <td>{{ liab.name }} <span class="text-xs text-gray-400">({{ liab.code }})</span></td>
              <td class="text-right">{{ formatIDR(liab.balance) }}</td>
            </tr>
          </tbody>
        </table>

        <h3 class="subsection-title">Ekuitas (Equity)</h3>
        <table class="report-table">
          <tbody>
            <tr v-for="eq in equities" :key="eq.code">
              <td>{{ eq.name }} <span class="text-xs text-gray-400">({{ eq.code }})</span></td>
              <td class="text-right">{{ formatIDR(eq.balance) }}</td>
            </tr>
            <tr>
              <td>Laba/Rugi Berjalan <span class="text-xs text-gray-400">(Net Income)</span></td>
              <td class="text-right font-bold" :class="netIncome >= 0 ? 'text-green-600' : 'text-red-600'">
                {{ formatIDR(netIncome) }}
              </td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="total-row bg-blue-50">
              <td>Total Kewajiban & Ekuitas</td>
              <td class="text-right font-bold text-blue-800">{{ formatIDR(totalLiabilities + totalEquity) }}</td>
            </tr>
          </tfoot>
        </table>
      </div>

      <!-- BALANCE CHECK -->
      <div class="balance-check mt-8" :class="totalAssets === (totalLiabilities + totalEquity) ? 'bg-green-50 border-green-200 text-green-700' : 'bg-red-50 border-red-200 text-red-700'">
        <div class="font-bold">{{ totalAssets === (totalLiabilities + totalEquity) ? 'BALANCE' : 'UNBALANCED' }}</div>
        <div>Selisih: {{ formatIDR(Math.abs(totalAssets - (totalLiabilities + totalEquity))) }}</div>
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
  font-size: 1.2rem;
  font-weight: 700;
  color: #1e293b;
  border-bottom: 2px solid #e2e8f0;
  padding-bottom: 8px;
  margin-bottom: 12px;
}
.subsection-title {
  font-size: 1rem;
  font-weight: 600;
  color: #475569;
  margin: 16px 0 8px 0;
}
.report-table {
  width: 100%;
  border-collapse: collapse;
}
.report-table td {
  padding: 8px 12px;
  color: #475569;
}
.report-table tfoot td {
  padding-top: 12px;
  padding-bottom: 12px;
  border-top: 1px solid #cbd5e1;
  color: #0f172a;
}
.total-row {
  font-weight: 700;
}

.balance-check {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-radius: 8px;
  border: 2px solid;
  font-size: 1.1rem;
}

.mt-8 { margin-top: 32px; }
.mb-4 { margin-bottom: 16px; }
.text-right { text-align: right; }
.text-center { text-align: center; }
.p-8 { padding: 32px; }
.text-xs { font-size: 0.75rem; }
.text-gray-400 { color: #94a3b8; }
.text-gray-500 { color: #64748b; }
.italic { font-style: italic; }
.font-bold { font-weight: 700; }
.bg-blue-50 { background-color: #eff6ff; }
.text-blue-800 { color: #1e40af; }
.bg-green-50 { background-color: #f0fdf4; }
.border-green-200 { border-color: #bbf7d0; }
.text-green-600 { color: #16a34a; }
.text-green-700 { color: #15803d; }
.bg-red-50 { background-color: #fef2f2; }
.border-red-200 { border-color: #fecaca; }
.text-red-600 { color: #dc2626; }
.text-red-700 { color: #b91c1c; }
</style>
