<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { financeService, type Invoice, type ChartOfAccount } from '@/services/admin/financeService'
import { useToastStore } from '@/stores/toastStore'

const toastStore = useToastStore()
const invoices = ref<Invoice[]>([])
const coas = ref<ChartOfAccount[]>([])
const loading = ref(true)

const activeTab = ref('AP') // AP or AR
const showModal = ref(false)
const showPayModal = ref(false)

const form = ref<Invoice>({
  type: 'AP',
  partner_name: '',
  date: new Date().toISOString().split('T')[0] || '',
  due_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0] || '',
  amount: 0,
  description: ''
})

const payForm = ref({
  invoice_id: '',
  account_id: ''
})

async function fetchData() {
  try {
    loading.value = true
    const [invs, allCoas] = await Promise.all([
      financeService.getInvoices(),
      financeService.getCOAs()
    ])
    invoices.value = invs
    // Hanya akun Kas/Bank (Asset) untuk pembayaran
    coas.value = allCoas.filter(c => c.type === 'Asset')
  } catch (error: any) {
    toastStore.error(error.message || 'Failed to fetch data')
  } finally {
    loading.value = false
  }
}

const filteredInvoices = computed(() => {
  return invoices.value.filter(i => i.type === activeTab.value)
})

const totalUnpaid = computed(() => {
  return filteredInvoices.value
    .filter(i => i.status === 'Unpaid')
    .reduce((sum, i) => sum + (i.amount || 0), 0)
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

function openModal() {
  form.value = {
    type: activeTab.value,
    partner_name: '',
    date: new Date().toISOString().split('T')[0] || '',
    due_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0] || '',
    amount: 0,
    description: ''
  }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

async function saveInvoice() {
  try {
    await financeService.createInvoice({
      ...form.value,
      amount: Number(form.value.amount)
    })
    toastStore.success('Invoice saved successfully')
    closeModal()
    fetchData()
  } catch (error: any) {
    toastStore.error(error.message || 'Failed to save invoice')
  }
}

function openPayModal(inv: Invoice) {
  if (!inv.id) return
  payForm.value = {
    invoice_id: inv.id,
    account_id: coas.value.length > 0 ? (coas.value[0]?.code || '') : ''
  }
  showPayModal.value = true
}

function closePayModal() {
  showPayModal.value = false
}

async function payInvoice() {
  try {
    await financeService.payInvoice(payForm.value.invoice_id, payForm.value.account_id)
    toastStore.success('Payment recorded successfully')
    closePayModal()
    fetchData()
  } catch (error: any) {
    toastStore.error(error.message || 'Failed to process payment')
  }
}

onMounted(() => {
  fetchData()
})
</script>

<template>
  <div class="finance-page">
    <div class="header">
      <div>
        <h1>Hutang & Piutang (AP & AR)</h1>
        <p class="subtitle">Kelola Accounts Payable & Accounts Receivable</p>
      </div>
      <button class="btn-primary" @click="openModal()">+ Add {{ activeTab }}</button>
    </div>

    <div class="tabs">
      <button :class="{ active: activeTab === 'AP' }" @click="activeTab = 'AP'">Hutang (Accounts Payable)</button>
      <button :class="{ active: activeTab === 'AR' }" @click="activeTab = 'AR'">Piutang (Accounts Receivable)</button>
    </div>

    <div class="stat-card">
      <h3>Total Unpaid {{ activeTab }}</h3>
      <h2 class="text-blue">{{ formatIDR(totalUnpaid) }}</h2>
    </div>

    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>Invoice No</th>
            <th>Partner</th>
            <th>Date</th>
            <th>Due Date</th>
            <th>Description</th>
            <th class="text-right">Amount</th>
            <th class="text-center">Status</th>
            <th class="text-center">Action</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td colspan="8" class="text-center">Loading...</td>
          </tr>
          <tr v-else-if="filteredInvoices.length === 0">
            <td colspan="8" class="text-center">No invoices found.</td>
          </tr>
          <tr v-else v-for="inv in filteredInvoices" :key="inv.id">
            <td class="font-bold">{{ inv.invoice_number }}</td>
            <td>{{ inv.partner_name }}</td>
            <td>{{ formatDate(inv.date) }}</td>
            <td :class="{ 'text-red font-bold': new Date(inv.due_date) < new Date() && inv.status === 'Unpaid' }">
              {{ formatDate(inv.due_date) }}
            </td>
            <td>{{ inv.description }}</td>
            <td class="text-right font-bold">{{ formatIDR(inv.amount) }}</td>
            <td class="text-center">
              <span class="badge" :class="inv.status === 'Paid' ? 'badge-success' : 'badge-warning'">
                {{ inv.status }}
              </span>
            </td>
            <td class="text-center">
              <button v-if="inv.status === 'Unpaid'" class="btn-pay" @click="openPayModal(inv)">Pay</button>
              <span v-else class="text-gray-400">-</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create Invoice Modal -->
    <div class="modal-overlay" v-if="showModal" @click.self="closeModal">
      <div class="modal">
        <h2>New {{ activeTab === 'AP' ? 'Hutang (AP)' : 'Piutang (AR)' }}</h2>
        <form @submit.prevent="saveInvoice">
          <div class="form-group">
            <label>Partner Name (Supplier / Customer)</label>
            <input type="text" v-model="form.partner_name" required placeholder="e.g. PT Jaya Abadi" />
          </div>
          <div class="split-group">
            <div class="form-group w-full">
              <label>Date</label>
              <input type="date" v-model="form.date" required />
            </div>
            <div class="form-group w-full">
              <label>Due Date</label>
              <input type="date" v-model="form.due_date" required />
            </div>
          </div>
          <div class="form-group">
            <label>Amount (IDR)</label>
            <input type="number" v-model="form.amount" required min="1" />
          </div>
          <div class="form-group">
            <label>Description</label>
            <textarea v-model="form.description" required rows="2" placeholder="e.g. Pembelian perlengkapan"></textarea>
          </div>
          <div class="modal-actions">
            <button type="button" class="btn-secondary" @click="closeModal">Cancel</button>
            <button type="submit" class="btn-primary">Save Invoice</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Pay Invoice Modal -->
    <div class="modal-overlay" v-if="showPayModal" @click.self="closePayModal">
      <div class="modal">
        <h2>Process Payment</h2>
        <p class="text-sm text-gray-500 mb-4">Membayar tagihan ini akan secara otomatis membuat Jurnal Umum baru.</p>
        <form @submit.prevent="payInvoice">
          <div class="form-group">
            <label>Pay From / Receive To Account</label>
            <select v-model="payForm.account_id" required>
              <option v-for="c in coas" :key="c.code" :value="c.code">{{ c.code }} - {{ c.name }}</option>
            </select>
          </div>
          <div class="modal-actions">
            <button type="button" class="btn-secondary" @click="closePayModal">Cancel</button>
            <button type="submit" class="btn-primary">Confirm Payment</button>
          </div>
        </form>
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
.tabs {
  display: flex;
  gap: 12px;
  border-bottom: 2px solid #e2e8f0;
}
.tabs button {
  background: none;
  border: none;
  padding: 12px 24px;
  font-weight: 600;
  color: #64748b;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  margin-bottom: -2px;
  transition: all 0.2s;
}
.tabs button.active {
  color: #e4793b;
  border-bottom-color: #e4793b;
}
.tabs button:hover:not(.active) {
  color: #475569;
}
.stat-card {
  background: white;
  padding: 24px;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
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
.text-red { color: #dc2626; }
.text-gray-400 { color: #94a3b8; }
.text-gray-500 { color: #64748b; }
.font-bold { font-weight: 600; }
.text-right { text-align: right !important; }
.text-center { text-align: center !important; }
.mb-4 { margin-bottom: 16px; }
.text-sm { font-size: 0.875rem; }

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
.badge {
  padding: 6px 12px;
  border-radius: 99px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}
.badge-warning {
  background: #fef08a;
  color: #854d0e;
}
.badge-success {
  background: #bbf7d0;
  color: #166534;
}
.btn-pay {
  background: #10b981;
  color: white;
  border: none;
  padding: 6px 16px;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.2s;
}
.btn-pay:hover {
  opacity: 0.9;
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
.form-group input, .form-group select, .form-group textarea {
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
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}
</style>
