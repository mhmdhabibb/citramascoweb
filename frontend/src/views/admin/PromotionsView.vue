<script setup>
import { ref, computed, onMounted } from 'vue'
import { useOfferStore } from '@/stores/offerStore'
import { useToastStore } from '@/stores/toastStore'

const offerStore = useOfferStore()
const toastStore = useToastStore()
const promotions = ref([])
const loading = ref(false)

const searchQuery = ref('')
const statusFilter = ref('All')

// Modal State
const isModalOpen = ref(false)
const isEditing = ref(false)
const editingId = ref('')
const imageFile = ref(null)

const form = ref({
  name: '',
  code: '',
  type: 'Percentage',
  value: 10,
  startDate: '',
  endDate: '',
  max_quota: 100,
  description: '',
  image: ''
})

// Stats
const totalPromotions = computed(() => promotions.value.length)
const activePromotions = computed(() => promotions.value.filter(p => p.status === 'active').length)
const expiredPromotions = computed(() => promotions.value.filter(p => p.status !== 'active').length)

// Date conversion helpers
const formatDateToInput = (dateStr) => {
  if (!dateStr) return ''
  if (/^\d{4}-\d{2}-\d{2}$/.test(dateStr)) return dateStr
  const parts = dateStr.split('-')
  if (parts.length === 3) {
    if (parts[0].length === 4) {
      return `${parts[0]}-${parts[1]}-${parts[2]}`
    }
    return `${parts[2]}-${parts[1]}-${parts[0]}`
  }
  return dateStr
}

const formatDateToApi = (dateStr) => {
  if (!dateStr) return ''
  const parts = dateStr.split('-')
  if (parts.length === 3) {
    if (parts[0].length === 4) {
      return `${parts[2]}-${parts[1]}-${parts[0]}`
    }
  }
  return dateStr
}

// Fetch list
const loadPromotions = async () => {
  try {
    loading.value = true
    await offerStore.fetchAdminOffers()
    promotions.value = offerStore.offers
  } catch (error) {
    console.error('Failed to load offers:', error)
  } finally {
    loading.value = false
  }
}

onMounted(loadPromotions)

// Filter Logic
const filteredPromotions = computed(() => {
  return promotions.value.filter(p => {
    const titleText = p.title || ''
    const codeText = p.code || ''
    const matchesSearch = titleText.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                          codeText.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesStatus = statusFilter.value === 'All' || 
                          (statusFilter.value === 'Active' && p.status === 'active') ||
                          (statusFilter.value === 'Expired' && p.status !== 'active')
    return matchesSearch && matchesStatus
  })
})

// Modal Open/Close Actions
const openCreateModal = () => {
  isEditing.value = false
  editingId.value = ''
  imageFile.value = null
  form.value = {
    name: '',
    code: '',
    type: 'Percentage',
    value: 10,
    startDate: new Date().toISOString().split('T')[0],
    endDate: '',
    max_quota: 100,
    description: '',
    image: ''
  }
  isModalOpen.value = true
}

const openEditModal = (item) => {
  isEditing.value = true
  editingId.value = item.id || ''
  imageFile.value = null
  
  // Determine type and value
  let discType = 'Percentage'
  let discValue = 0
  if (item.discount !== undefined && item.discount !== null && item.discount > 0) {
    discType = 'Percentage'
    discValue = item.discount
  } else if (item.price !== undefined && item.price !== null && item.price > 0) {
    discType = 'Fixed Amount'
    discValue = item.price
  }

  form.value = {
    name: item.title || '',
    code: item.code || '',
    type: discType,
    value: discValue,
    startDate: formatDateToInput(item.valid_start),
    endDate: formatDateToInput(item.valid_end),
    status: item.status === 'active' ? 'Active' : 'Expired',
    max_quota: item.max_quota || 0,
    description: item.description || '',
    image: item.image || ''
  }
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
}

const handleImageUpload = (event) => {
  const file = event.target.files[0]
  if (file) {
    imageFile.value = file
  }
}

const savePromotion = async () => {
  if (!form.value.name.trim()) {
    toastStore.warning("Mohon masukkan nama promosi.")
    return
  }
  if (!form.value.code.trim()) {
    toastStore.warning("Mohon masukkan kode promo.")
    return
  }

  try {
    loading.value = true
    const formData = new FormData()
    formData.append('title', form.value.name.trim())
    formData.append('code', form.value.code.trim().toUpperCase())
    
    // Fixed amount vs percentage mapping
    if (form.value.type === 'Percentage') {
      formData.append('discount', String(form.value.value))
      formData.append('price', '0')
    } else {
      formData.append('price', String(form.value.value))
      formData.append('discount', '0')
    }

    if (form.value.startDate) {
      formData.append('valid_start', formatDateToApi(form.value.startDate))
    }
    if (form.value.endDate) {
      formData.append('valid_end', formatDateToApi(form.value.endDate))
    }

    formData.append('max_quota', String(form.value.max_quota))
    formData.append('description', form.value.description)
    
    // Status mapping
    const apiStatus = form.value.status === 'Active' ? 'active' : 'archieved'
    formData.append('status', apiStatus)

    if (imageFile.value) {
      formData.append('image', imageFile.value)
    }

    if (isEditing.value) {
      const msg = await offerStore.update(editingId.value, formData)
      toastStore.success(msg || 'Penawaran berhasil diperbarui!')
    } else {
      const msg = await offerStore.store(formData)
      toastStore.success(msg || 'Penawaran berhasil ditambahkan!')
    }

    await loadPromotions()
    closeModal()
  } catch (error) {
    console.error(error)
    toastStore.error(error.response?.data?.message || error.message || 'Terjadi kesalahan saat menyimpan penawaran')
  } finally {
    loading.value = false
  }
}

const deletePromotion = async (id) => {
  if (confirm(`Are you sure you want to delete this offer/promotion?`)) {
    try {
      loading.value = true
      const msg = await offerStore.destroy(id)
      toastStore.success(msg || 'Penawaran berhasil dihapus!')
      await loadPromotions()
    } catch (error) {
      console.error(error)
      toastStore.error(error.response?.data?.message || error.message || 'Terjadi kesalahan saat menghapus penawaran')
    } finally {
      loading.value = false
    }
  }
}
</script>

<template>
  <div class="promotions-view">
    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <h3 class="stat-title">Total Offers</h3>
        <p class="stat-value">{{ totalPromotions }}</p>
      </div>
      <div class="stat-card">
        <h3 class="stat-title">Active Offers</h3>
        <p class="stat-value text-success">{{ activePromotions }}</p>
      </div>
      <div class="stat-card">
        <h3 class="stat-title">Expired/Draft Offers</h3>
        <p class="stat-value text-danger">{{ expiredPromotions }}</p>
      </div>
    </div>

    <!-- Control Bar -->
    <div class="control-bar">
      <div class="search-box">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Search by Promo Name or Code..." 
          class="search-input"
        />
      </div>
      <div class="filter-box">
        <select v-model="statusFilter" class="filter-select">
          <option value="All">All Statuses</option>
          <option value="Active">Active</option>
          <option value="Expired">Expired/Draft</option>
        </select>
      </div>
      <button @click="openCreateModal" class="btn btn-primary add-btn">
        + Add Offer/Promo
      </button>
    </div>

    <!-- Data Table Card -->
    <div class="table-card">
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>Image</th>
              <th>Promo Name</th>
              <th>Code</th>
              <th>Discount Value</th>
              <th>Start Date</th>
              <th>End Date</th>
              <th>Quota</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="p in filteredPromotions" :key="p.id">
              <td>
                <img 
                  v-if="p.image" 
                  :src="p.image" 
                  alt="Promo Image" 
                  class="promo-thumbnail"
                />
                <span v-else class="text-gray-400 text-xs">No Cover</span>
              </td>
              <td class="font-bold">{{ p.title }}</td>
              <td class="font-mono code-cell">{{ p.code }}</td>
              <td>
                <span v-if="p.discount && p.discount > 0">{{ p.discount }}%</span>
                <span v-else-if="p.price && p.price > 0">Rp {{ p.price.toLocaleString('id-ID') }}</span>
                <span v-else>0</span>
              </td>
              <td>{{ p.valid_start || '-' }}</td>
              <td>{{ p.valid_end || 'Ongoing' }}</td>
              <td>{{ p.max_quota !== null && p.max_quota !== undefined ? p.max_quota : 'Unlimited' }}</td>
              <td>
                <span class="badge" :class="{
                  'badge-success': p.status === 'active',
                  'badge-danger': p.status !== 'active'
                }">{{ p.status === 'active' ? 'Active' : 'Expired/Draft' }}</span>
              </td>
              <td>
                <div class="action-buttons">
                  <button 
                    @click="openEditModal(p)" 
                    class="btn btn-sm btn-success-outline"
                  >Edit</button>
                  <button 
                    @click="deletePromotion(p.id)" 
                    class="btn btn-sm btn-danger-outline"
                  >Delete</button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredPromotions.length === 0">
              <td colspan="9" class="no-data">No promotions or offers found.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create/Edit Modal Overlay -->
    <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
      <div class="modal-card">
        <div class="modal-header">
          <h3>{{ isEditing ? 'Edit Offer/Promotion' : 'Add New Offer/Promotion' }}</h3>
          <button @click="closeModal" class="close-btn">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">Promotion Name</label>
            <input 
              v-model="form.name" 
              type="text" 
              placeholder="e.g. Summer Special Deal" 
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label class="form-label">Promo Code</label>
            <input 
              v-model="form.code" 
              type="text" 
              placeholder="e.g. SUMMER26" 
              class="form-input code-input"
            />
          </div>
          <div class="form-row">
            
            <div class="form-group flex-1">
              <label class="form-label">Discount Value</label>
              <input 
                v-model.number="form.value" 
                type="number" 
                min="0" 
                class="form-input"
              />
            </div>
          </div>
          <div class="form-row">
            <div class="form-group flex-1">
              <label class="form-label">Start Date</label>
              <input 
                v-model="form.startDate" 
                type="date" 
                class="form-input"
              />
            </div>
            <div class="form-group flex-1">
              <label class="form-label">End Date</label>
              <input 
                v-model="form.endDate" 
                type="date" 
                class="form-input"
              />
            </div>
          </div>
          <div class="form-row">
            <div class="form-group flex-1">
              <label class="form-label">Quota Limit</label>
              <input 
                v-model.number="form.max_quota" 
                type="number" 
                min="0" 
                class="form-input"
              />
            </div>
           
          </div>
          <div class="form-group">
            <label class="form-label">Description</label>
            <textarea 
              v-model="form.description" 
              placeholder="Provide offer description here..." 
              rows="3" 
              class="form-input resize-none"
            ></textarea>
          </div>
          <div class="form-group">
            <label class="form-label">Offer Cover Image</label>
            <input 
              type="file" 
              accept="image/*" 
              @change="handleImageUpload" 
              class="form-input"
            />
            <div v-if="form.image" class="current-image-preview mt-2">
              <img :src="form.image" alt="Current Image" class="w-20 h-20 object-cover rounded-lg border" />
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeModal" class="btn btn-secondary">Cancel</button>
          <button @click="savePromotion" class="btn btn-primary">Save Changes</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.promotions-view {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.stat-card {
  background-color: #ffffff;
  padding: 20px 24px;
  border-radius: 16px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.stat-title {
  font-size: 0.875rem;
  color: #64748b;
  margin: 0 0 6px 0;
  font-weight: 600;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: #0f172a;
  margin: 0;
}

.text-success { color: #10b981; }
.text-danger { color: #ef4444; }

.control-bar {
  display: flex;
  gap: 16px;
  background-color: #ffffff;
  padding: 16px;
  border-radius: 16px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  align-items: center;
  flex-wrap: wrap;
}

.search-box {
  flex-grow: 1;
  min-width: 250px;
}

.search-input {
  width: 100%;
  padding: 10px 14px;
  border-radius: 10px;
  border: 1px solid #cbd5e1;
  outline: none;
  font-size: 0.875rem;
  color: #334155;
  box-sizing: border-box;
}

.search-input:focus {
  border-color: #e15b2b;
  box-shadow: 0 0 0 2px rgba(225, 91, 43, 0.1);
}

.filter-select {
  padding: 10px 14px;
  border-radius: 10px;
  border: 1px solid #cbd5e1;
  outline: none;
  font-size: 0.875rem;
  color: #334155;
  background-color: #ffffff;
  cursor: pointer;
}

.filter-select:focus {
  border-color: #e15b2b;
}

.table-card {
  background-color: #ffffff;
  border-radius: 16px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  overflow: hidden;
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
  padding: 14px 20px;
  border-bottom: 1px solid #e2e8f0;
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.data-table td {
  padding: 16px 20px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 0.875rem;
  color: #334155;
  vertical-align: middle;
}

.font-bold { font-weight: 600; }
.font-mono { font-family: monospace; }

.code-cell {
  background-color: #f8fafc;
  padding: 4px 8px;
  border-radius: 6px;
  font-weight: 700;
  color: #0f172a;
  border: 1px solid #e2e8f0;
  display: inline-block;
}

.badge {
  padding: 4px 8px;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  display: inline-block;
}

.badge-success { background-color: #d1fae5; color: #065f46; }
.badge-danger { background-color: #fee2e2; color: #991b1b; }

.action-buttons {
  display: flex;
  gap: 8px;
}

.btn {
  padding: 8px 16px;
  border-radius: 10px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s;
}

.btn-sm {
  padding: 4px 8px;
  font-size: 0.75rem;
  border-radius: 6px;
}

.btn-primary {
  background-color: #e15b2b;
  color: #ffffff;
}
.btn-primary:hover { background-color: #c84e20; }

.btn-success-outline {
  border-color: #10b981;
  color: #10b981;
  background: transparent;
}
.btn-success-outline:hover {
  background-color: #e6fcf5;
}

.btn-danger-outline {
  border-color: #ef4444;
  color: #ef4444;
  background: transparent;
}
.btn-danger-outline:hover {
  background-color: #fee2e2;
}

.no-data {
  text-align: center;
  color: #64748b;
  padding: 32px !important;
}

/* Modal CSS styling */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(15, 23, 42, 0.3);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  animation: fadeIn 0.2s ease;
}

.modal-card {
  background-color: #ffffff;
  border-radius: 16px;
  width: 90%;
  max-width: 480px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  animation: slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-header {
  padding: 20px 24px;
  border-bottom: 1px solid #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 700;
  color: #0f172a;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #64748b;
  cursor: pointer;
  padding: 4px;
  line-height: 1;
}

.close-btn:hover {
  color: #0f172a;
}

.modal-body {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-row {
  display: flex;
  gap: 16px;
}

.flex-1 {
  flex: 1;
}

.form-label {
  font-size: 0.75rem;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.form-input, .form-select {
  padding: 10px 14px;
  border-radius: 8px;
  border: 1px solid #cbd5e1;
  outline: none;
  font-size: 0.875rem;
  color: #334155;
  background-color: #ffffff;
  width: 100%;
  box-sizing: border-box;
}

.code-input {
  text-transform: uppercase;
  font-family: monospace;
  font-weight: 700;
}

.form-input:focus, .form-select:focus {
  border-color: #e15b2b;
  box-shadow: 0 0 0 2px rgba(225, 91, 43, 0.1);
}

.modal-footer {
  padding: 16px 24px;
  background-color: #f8fafc;
  border-top: 1px solid #f1f5f9;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.btn-secondary {
  border-color: #cbd5e1;
  color: #334155;
  background-color: #ffffff;
}

.btn-secondary:hover {
  background-color: #f1f5f9;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.promo-thumbnail {
  width: 48px;
  height: 48px;
  object-fit: cover;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
}

.resize-none {
  resize: none;
}
</style>
