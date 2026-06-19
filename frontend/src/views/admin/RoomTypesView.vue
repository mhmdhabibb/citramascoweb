<script setup>
import { useTypeStore } from '@/stores/typeStore'
import { useToastStore } from '@/stores/toastStore'
import { ref, computed, onMounted } from 'vue'

const types = ref([])
const typeStore = useTypeStore()
const toastStore = useToastStore()
const loading = ref(false)
const editingIndex = ref(-1)
const editingId = ref('')
const searchQuery = ref('')

// Modal State
const isModalOpen = ref(false)
const isEditing = ref(false)
const form = ref({
  name: ''
})

onMounted(async () => {
  try {
    loading.value = true 
    await typeStore.fetchTypes()
    types.value = typeStore.types
    console.log("Types loaded:", types.value)
  } catch (error) {
    console.error('Error fetching types:', error)
  } finally {
    loading.value = false
  }
})

// Computed List
const filteredRoomTypes = computed(() => {
  if (!searchQuery.value) return types.value
  const query = searchQuery.value.toLowerCase()
  return types.value.filter(t => 
    (t.name || '').toLowerCase().includes(query) ||
    (t.slug || '').toLowerCase().includes(query)
  )
})

// Modal Open/Close Actions
const openCreateModal = () => {
  isEditing.value = false
  editingIndex.value = -1
  editingId.value = ''
  form.value = {
    name: ''
  }
  isModalOpen.value = true
}

const openEditModal = (item, index) => {
  isEditing.value = true
  editingIndex.value = index
  editingId.value = item.id || ''
  form.value = { ...item }
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
}

const saveRoomType = async () => {
  if (!form.value.name.trim()) {
    toastStore.warning("Mohon masukkan nama tipe kamar.")
    return
  }

  try {
    loading.value = true
    const payload = {
      name: form.value.name.trim()
    }

    if (isEditing.value) {
      await typeStore.update(editingId.value, payload)
      if (typeStore.error) {
        toastStore.error(`Gagal mengubah tipe: ${typeStore.error}`)
        return
      }
      toastStore.success('Tipe kamar berhasil diperbarui!')
    } else {
      await typeStore.store(payload)
      if (typeStore.error) {
        toastStore.error(`Gagal menambah tipe: ${typeStore.error}`)
        return
      }
      toastStore.success('Tipe kamar berhasil ditambahkan!')
    }

    // Refresh list
    await typeStore.fetchTypes()
    types.value = typeStore.types
    closeModal()
  } catch (error) {
    console.error(error)
    toastStore.error('Terjadi kesalahan saat menyimpan tipe kamar')
  } finally {
    loading.value = false
  }
}

const deleteRoomType = async (index) => {
  const item = types.value[index]
  if (confirm(`Apakah Anda yakin ingin menghapus Tipe: ${item.name}?`)) {
    try {
      loading.value = true
      await typeStore.destroy(item.id)
      if (typeStore.error) {
        toastStore.error(`Gagal menghapus tipe: ${typeStore.error}`)
        return
      }
      toastStore.success('Tipe kamar berhasil dihapus!')
      await typeStore.fetchTypes()
      types.value = typeStore.types
    } catch (error) {
      console.error(error)
      toastStore.error('Terjadi kesalahan saat menghapus tipe kamar')
    } finally {
      loading.value = false
    }
  }
}
</script>

<template>
  <div class="room-types-view">
    <!-- Control Bar -->
    <div class="control-bar">
      <div class="search-box">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Search by Room Type or amenities..." 
          class="search-input"
        />
      </div>
      <button @click="openCreateModal" class="btn btn-primary add-btn">
        + Add Room Type
      </button>
    </div>

    <!-- Data Table Card -->
    <div class="table-card">
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>Type Name</th>
              <th>Slug</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(t, index) in filteredRoomTypes" :key="t.id || t.name">
              <td class="font-bold">{{ t.name }}</td>
              <td>{{ t.slug }}</td>
              <td>
                <div class="action-buttons">
                  <button 
                    @click="openEditModal(t, index)" 
                    class="btn btn-sm btn-success-outline"
                  >Edit</button>
                  <button 
                    @click="deleteRoomType(index)" 
                    class="btn btn-sm btn-danger-outline"
                  >Delete</button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredRoomTypes.length === 0">
              <td colspan="3" class="no-data">No room types found.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create/Edit Modal Overlay -->
    <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
      <div class="modal-card">
        <div class="modal-header">
          <h3>{{ isEditing ? 'Edit Room Type' : 'Add New Room Type' }}</h3>
          <button @click="closeModal" class="close-btn">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">Type Name</label>
            <input 
              v-model="form.name" 
              type="text" 
              placeholder="e.g. Double Deluxe" 
              class="form-input"
            />
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeModal" class="btn btn-secondary">Cancel</button>
          <button @click="saveRoomType" class="btn btn-primary" :disabled="loading">Save Changes</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.room-types-view {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.control-bar {
  display: flex;
  gap: 16px;
  background-color: #ffffff;
  padding: 16px;
  border-radius: 16px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  align-items: center;
}

.search-box {
  flex-grow: 1;
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

.amenities-cell {
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

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
  max-width: 440px;
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

.form-label {
  font-size: 0.75rem;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.form-input, .form-textarea {
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

.form-input:focus, .form-textarea:focus {
  border-color: #e15b2b;
  box-shadow: 0 0 0 2px rgba(225, 91, 43, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
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
</style>
