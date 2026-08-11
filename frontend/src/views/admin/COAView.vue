<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { financeService, type ChartOfAccount } from '@/services/admin/financeService'
import { useToastStore } from '@/stores/toastStore'

const toastStore = useToastStore()
const coas = ref<ChartOfAccount[]>([])
const loading = ref(true)

const showModal = ref(false)
const isEditing = ref(false)

const form = ref({
  code: '',
  name: '',
  type: 'Asset',
  parent_code: ''
})

const fileInput = ref<HTMLInputElement | null>(null)

const accountTypes = ['Asset', 'Liability', 'Equity', 'Revenue', 'Expense']

async function fetchCOAs() {
  try {
    loading.value = true
    coas.value = await financeService.getCOAs()
  } catch (error: any) {
    toastStore.error(error.message || 'Failed to fetch COA')
  } finally {
    loading.value = false
  }
}

function openModal(coa?: ChartOfAccount) {
  if (coa) {
    isEditing.value = true
    form.value = { ...coa, parent_code: coa.parent_code || '' }
  } else {
    isEditing.value = false
    form.value = { code: '', name: '', type: 'Asset', parent_code: '' }
  }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

async function saveCOA() {
  try {
    if (isEditing.value) {
      await financeService.updateCOA(form.value.code, {
        name: form.value.name,
        type: form.value.type,
        parent_code: form.value.parent_code || undefined,
      })
      toastStore.success('Account updated successfully')
    } else {
      await financeService.createCOA(form.value)
      toastStore.success('Account created successfully')
    }
    closeModal()
    fetchCOAs()
  } catch (error: any) {
    toastStore.error(error.message || 'Failed to save COA')
  }
}

async function deleteCOA(code: string) {
  if (confirm('Are you sure you want to delete this account?')) {
    try {
      await financeService.deleteCOA(code)
      toastStore.success('Account deleted successfully')
      fetchCOAs()
    } catch (error: any) {
      toastStore.error(error.message || 'Failed to delete COA')
    }
  }
}

onMounted(() => {
  fetchCOAs()
})

async function onFileSelected(event: Event) {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    const file = target.files[0]
    if (!file) return
    try {
      loading.value = true
      await financeService.importCOA(file)
      toastStore.success('Excel imported successfully')
      fetchCOAs()
    } catch (error: any) {
      toastStore.error(error.message || 'Failed to import Excel')
      loading.value = false
    } finally {
      if (fileInput.value) fileInput.value.value = ''
    }
  }
}

function triggerImport() {
  fileInput.value?.click()
}
</script>

<template>
  <div class="coa-page">
      <div class="header">
        <h1>Chart of Accounts</h1>
        <div class="actions">
          <input type="file" ref="fileInput" accept=".xlsx, .xls" style="display: none" @change="onFileSelected" />
          <button class="btn-secondary" @click="triggerImport()">Import Excel</button>
          <button class="btn-primary" @click="openModal()">+ New Account</button>
        </div>
      </div>

      <div class="table-container">
        <table>
          <thead>
            <tr>
              <th>Code</th>
              <th>Name</th>
              <th>Type</th>
              <th>Parent</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="5" class="text-center">Loading...</td>
            </tr>
            <tr v-else-if="coas.length === 0">
              <td colspan="5" class="text-center">No accounts found.</td>
            </tr>
            <tr v-else v-for="coa in coas" :key="coa.code">
              <td class="font-bold">{{ coa.code }}</td>
              <td>{{ coa.name }}</td>
              <td><span class="badge" :class="'badge-' + coa.type.toLowerCase()">{{ coa.type }}</span></td>
              <td>{{ coa.parent_code || '-' }}</td>
              <td>
                <button class="btn-text" @click="openModal(coa)">Edit</button>
                <button class="btn-text text-danger" @click="deleteCOA(coa.code)">Delete</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Modal -->
      <div class="modal-overlay" v-if="showModal" @click.self="closeModal">
        <div class="modal">
          <h2>{{ isEditing ? 'Edit Account' : 'New Account' }}</h2>
          <form @submit.prevent="saveCOA">
            <div class="form-group">
              <label>Account Code</label>
              <input type="text" v-model="form.code" required :disabled="isEditing" placeholder="e.g., 1-100" />
            </div>
            <div class="form-group">
              <label>Account Name</label>
              <input type="text" v-model="form.name" required placeholder="e.g., Kas Besar" />
            </div>
            <div class="form-group">
              <label>Parent Code (Optional)</label>
              <input type="text" v-model="form.parent_code" placeholder="e.g., 1-100" />
            </div>
            <div class="form-group">
              <label>Account Type</label>
              <select v-model="form.type" required>
                <option v-for="t in accountTypes" :key="t" :value="t">{{ t }}</option>
              </select>
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
.coa-page {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
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
.actions {
  display: flex;
  gap: 12px;
}
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
.font-bold {
  font-weight: 600;
}
.badge {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  background: #f1f5f9;
}
.badge-asset { background: #dbeafe; color: #1e40af; }
.badge-revenue { background: #dcfce7; color: #166534; }
.badge-expense { background: #fee2e2; color: #991b1b; }
.badge-liability { background: #ffedd5; color: #9a3412; }
.badge-equity { background: #f3e8ff; color: #6b21a8; }
.btn-primary {
  background: #e4793b;
  color: white;
  padding: 8px 16px;
  border-radius: 8px;
  border: none;
  cursor: pointer;
  font-weight: 600;
}
.btn-secondary {
  background: #f1f5f9;
  color: #475569;
  padding: 8px 16px;
  border-radius: 8px;
  border: none;
  cursor: pointer;
  font-weight: 600;
}
.btn-text {
  background: none;
  border: none;
  color: #e4793b;
  font-weight: 600;
  cursor: pointer;
  margin-right: 12px;
}
.text-danger {
  color: #ef4444;
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
  width: 400px;
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
}
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}
</style>
