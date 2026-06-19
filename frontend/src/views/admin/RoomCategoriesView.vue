<script setup>
import { useCategoryStore } from '@/stores/categoryStore'
import { useToastStore } from '@/stores/toastStore'
import { ref, computed, onMounted } from 'vue'

const categories = ref([])
const loading = ref(false)
const categoryStore = useCategoryStore()
const toastStore = useToastStore()

onMounted(async () => {
  try {
    loading.value = true 
    await categoryStore.fetchCategories()
    categories.value = categoryStore.categories
    console.log("Categories loaded:", categories.value)
  } catch (error) {
    console.error('Error fetching categories:', error)
  } finally {
    loading.value = false
  }
})

const searchQuery = ref('')

// Modal State
const isModalOpen = ref(false)
const isEditing = ref(false)
const form = ref({
  name: '',
  description: '',
  services: ''
})
const editingIndex = ref(-1)
const editingId = ref('')

// Computed List
const filteredRoomCategories = computed(() => {
  if (!searchQuery.value) return categories.value
  const query = searchQuery.value.toLowerCase()
  return categories.value.filter(c => 
    (c.name || '').toLowerCase().includes(query) ||
    (c.description || '').toLowerCase().includes(query)
  )
})

// Modal Open/Close Actions
const openCreateModal = () => {
  isEditing.value = false
  editingIndex.value = -1
  editingId.value = ''
  form.value = {
    name: '',
    description: '',
    services: ''
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

const saveRoomCategory = async () => {
  if(!form.value.name.trim()) {
    toastStore.warning('Nama kategori tidak boleh kosong!')
    return
  }

  try {
    loading.value = true 
    const payload = {
      name: form.value.name 
    }

    if (isEditing.value) {
      await categoryStore.update(editingId.value, payload)
      
      if (categoryStore.error) {
        toastStore.error(`Gagal mengubah kategori: ${categoryStore.error}`)
        return
      }
      toastStore.success('Kategori kamar berhasil diperbarui!')
    } else {
      await categoryStore.store(payload)

      if (categoryStore.error) {
        toastStore.error(`Gagal menambah kategori: ${categoryStore.error}`)
        return
      }
      toastStore.success('Kategori kamar berhasil ditambahkan!')
    }

    // Refresh categories in the view
    await categoryStore.fetchCategories()
    categories.value = categoryStore.categories
    
    closeModal()
  } catch (error) {
    console.error('Error saving category:', error)
    toastStore.error('Terjadi kesalahan saat menyimpan kategori')
  } finally {
    loading.value = false
  }
}

const deleteRoomCategory = async (index) => {
  const item = categories.value[index]
  if (confirm(`Apakah Anda yakin ingin menghapus Kategori: ${item.name}?`)) {
    try {
      loading.value = true
      await categoryStore.destroy(item.id)
      
      if (categoryStore.error) {
        toastStore.error(`Gagal menghapus kategori: ${categoryStore.error}`)
        return
      }

      toastStore.success('Kategori kamar berhasil dihapus!')
      await categoryStore.fetchCategories()
      categories.value = categoryStore.categories
    } catch (error) {
      console.error('Error deleting category:', error)
      toastStore.error('Terjadi kesalahan saat menghapus kategori')
    } finally {
      loading.value = false
    }
  }
}
</script>

<template>
  <div class="room-categories-view">
    <!-- Control Bar -->
    <div class="control-bar">
      <div class="search-box">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Search by category name or description..." 
          class="search-input"
        />
      </div>
      <button @click="openCreateModal" class="btn btn-primary add-btn">
        + Add Category
      </button>
    </div>

    <!-- Data Table Card -->
    <div class="table-card">
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>Category Name</th>
              <th>Slug</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(c, index) in filteredRoomCategories" :key="c.id || c.name">
              <td class="font-bold category-name">{{ c.name }}</td>
              <td class="desc-cell">{{ c.slug }}</td>
              <td>
                <div class="action-buttons">
                  <button 
                    @click="openEditModal(c, index)" 
                    class="btn btn-sm btn-success-outline"
                  >Edit</button>
                  <button 
                    @click="deleteRoomCategory(index)" 
                    class="btn btn-sm btn-danger-outline"
                  >Delete</button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredRoomCategories.length === 0">
              <td colspan="4" class="no-data">No categories found.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create/Edit Modal Overlay -->
    <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
      <div class="modal-card">
        <div class="modal-header">
          <h3>{{ isEditing ? 'Edit Category' : 'Add New Category' }}</h3>
          <button @click="closeModal" class="close-btn">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">Category Name</label>
            <input 
              v-model="form.name" 
              type="text" 
              placeholder="e.g. Co-working Loft" 
              class="form-input"
            />
          </div>
         
        </div>
        <div class="modal-footer">
          <button @click="closeModal" class="btn btn-secondary">Cancel</button>
          <button @click="saveRoomCategory" class="btn btn-primary">{{ isEditing ? 'Update' : 'Add' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.room-categories-view {
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

.category-name {
  min-width: 150px;
}

.desc-cell {
  max-width: 250px;
  line-height: 1.4;
}

.services-cell {
  max-width: 200px;
  line-height: 1.4;
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
