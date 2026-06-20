<script setup>
import { ref, computed, onMounted } from 'vue'
import { roomService } from '@/services/roomService'
import { categoryService } from '@/services/categoryService'
import { typeService } from '@/services/typeService'
import { useToastStore } from '@/stores/toastStore'

const rooms = ref([])
const roomCategories = ref([])
const roomTypes = ref([])
const roomStatusesOptions = ['Available', 'Occupied', 'Maintenance']
const toastStore = useToastStore()

const searchQuery = ref('')
const statusFilter = ref('All')
const loading = ref(false)

// Modal State
const isModalOpen = ref(false)
const isEditing = ref(false)
const editingId = ref('')
const imageFile = ref(null)

const form = ref({
  name: '',
  category_id: '',
  type_id: '',
  price: 0,
  capacity: 2,
  size: 30,
  description: '',
  status: 'Available'
})

// Stats
const totalRooms = computed(() => rooms.value.length)
const availableRooms = computed(() => rooms.value.filter(r => (r.status || 'Available') === 'Available').length)
const occupiedRooms = computed(() => rooms.value.filter(r => r.status === 'Occupied').length)
const maintenanceRooms = computed(() => rooms.value.filter(r => r.status === 'Maintenance').length)

// Filter Logic
const filteredRooms = computed(() => {
  return rooms.value.filter(room => {
    const code = room.code || ''
    const name = room.name || ''
    const typeName = room.type?.name || ''
    const categoryName = room.category?.name || ''
    
    const matchesSearch = code.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                          name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                          typeName.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                          categoryName.toLowerCase().includes(searchQuery.value.toLowerCase())
                          
    const status = room.status || 'Available'
    const matchesStatus = statusFilter.value === 'All' || 
                          status.toLowerCase() === statusFilter.value.toLowerCase()
    return matchesSearch && matchesStatus
  })
})

const handleImageUpload = (event) => {
  imageFile.value = event.target.files[0]
}

// Modal Actions
const openCreateModal = () => {
  isEditing.value = false
  editingId.value = ''
  imageFile.value = null
  form.value = {
    name: '',
    category_id: roomCategories.value[0]?.id || '',
    type_id: roomTypes.value[0]?.id || '',
    price: 350000,
    capacity: 2,
    size: 30,
    description: '',
  }
  isModalOpen.value = true
}

const openEditModal = (item) => {
  isEditing.value = true
  editingId.value = item.id
  imageFile.value = null
  form.value = {
    name: item.name || '',
    category_id: item.category_id || '',
    type_id: item.type_id || '',
    price: item.price || 0,
    capacity: item.capacity || 2,
    size: item.size || 30,
    description: item.description || '',
  }
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
}

const saveRoom = async () => {
  if (!form.value.category_id) {
    toastStore.warning("Mohon pilih kategori kamar.")
    return
  }
  if (!form.value.type_id) {
    toastStore.warning("Mohon pilih tipe kamar.")
    return
  }

  // Set room name to category name automatically
  const selectedCategory = roomCategories.value.find(c => c.id === form.value.category_id)
  const roomName = selectedCategory ? selectedCategory.name : ''
  form.value.name = roomName

  if (!roomName.trim()) {
    toastStore.warning("Mohon pilih kategori kamar.")
    return
  }

  try {
    loading.value = true
    const formData = new FormData()
    formData.append('name', roomName.trim())
    formData.append('category_id', form.value.category_id)
    formData.append('type_id', form.value.type_id)
    formData.append('price', String(form.value.price))
    formData.append('capacity', String(form.value.capacity))
    formData.append('size', String(form.value.size))
    formData.append('description', form.value.description)
    
    if (imageFile.value) {
      formData.append('image', imageFile.value)
    }

    if (isEditing.value) {
      const msg = await roomService.update(editingId.value, formData)
      toastStore.success(msg || 'Kamar berhasil diperbarui!')
    } else {
      const msg = await roomService.create(formData)
      toastStore.success(msg || 'Kamar baru berhasil ditambahkan!')
    }

    // Refresh rooms
    rooms.value = await roomService.getAll()
    closeModal()
  } catch (error) {
    console.error(error)
    toastStore.error(error.response?.data?.message || error.message || 'Gagal menyimpan detail kamar')
  } finally {
    loading.value = false
  }
}

// Inline list updates
const setStatus = (roomCode, newStatus) => {
  const room = rooms.value.find(r => r.code === roomCode)
  if (room) {
    room.status = newStatus
  }
}

const deleteRoom = async (id) => {
  if (confirm(`Are you sure you want to delete this room?`)) {
    try {
      loading.value = true
      const msg = await roomService.delete(id)
      toastStore.success(msg || 'Kamar berhasil dihapus!')
      rooms.value = await roomService.getAll()
    } catch (error) {
      console.error(error)
      toastStore.error(error.response?.data?.message || error.message || 'Gagal menghapus kamar')
    } finally {
      loading.value = false
    }
  }
}

onMounted(async () => {
  try {
    loading.value = true
    const [roomsData, categoriesData, typesData] = await Promise.all([
      roomService.getAll(),
      categoryService.getAll(),
      typeService.getAll()
    ])
    rooms.value = roomsData
    roomCategories.value = categoriesData
    roomTypes.value = typesData
  } catch (error) {
    console.error('Failed to initialize rooms manager:', error)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="rooms-view">
    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <h3 class="stat-title">Total Rooms</h3>
        <p class="stat-value">{{ totalRooms }}</p>
      </div>
      <div class="stat-card">
        <h3 class="stat-title">Available</h3>
        <p class="stat-value text-success">{{ availableRooms }}</p>
      </div>
      <div class="stat-card">
        <h3 class="stat-title">Occupied</h3>
        <p class="stat-value text-info">{{ occupiedRooms }}</p>
      </div>
      <div class="stat-card">
        <h3 class="stat-title">Maintenance</h3>
        <p class="stat-value text-warning">{{ maintenanceRooms }}</p>
      </div>
    </div>

    <!-- Control Bar -->
    <div class="control-bar">
      <div class="search-box">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Search by Room #, Type or Category..." 
          class="search-input"
        />
      </div>
      <!-- <div class="filter-box">
        <select v-model="statusFilter" class="filter-select">
          <option value="All">All Statuses</option>
          <option value="Available">Available</option>
          <option value="Occupied">Occupied</option>
          <option value="Maintenance">Maintenance</option>
        </select>
      </div> -->
      <button @click="openCreateModal" class="btn btn-primary add-btn">
        + Add Room
      </button>
    </div>

    <!-- Data Table Card -->
    <div class="table-card">
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>Room Code</th>
              <th>Room Name</th>
              <th>Type</th>
              <th>Price/Night</th>
              <th>Image</th>
              <th>Description</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="room in filteredRooms" :key="room.id">
              <td class="font-mono font-bold">{{ room.code }}</td>
              <td>{{ room.name }}</td>
              <td>{{ room.type?.name || 'N/A' }}</td>
              <td>Rp{{ room.price ? room.price.toLocaleString('id-ID') : 0 }}</td>
              <td>
                <img 
                  v-if="room.image" 
                  :src="room.image" 
                  alt="Room Image" 
                  class="room-thumbnail"
                />
                <span v-else class="text-muted">No Image</span>
              </td>
                           <td>{{ room.description }}</td>
              <td>
                <div class="action-buttons">
                  <button 
                    @click="openEditModal(room)" 
                    class="btn btn-sm btn-primary-outline"
                  >Edit</button>
                  <button 
                    @click="deleteRoom(room.id)" 
                    class="btn btn-sm btn-danger-outline"
                  >Delete</button>
                </div>
              </td>

      
            </tr>
            <tr v-if="filteredRooms.length === 0">
              <td colspan="7" class="no-data">No rooms found.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create/Edit Modal Overlay -->
    <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
      <div class="modal-card">
        <div class="modal-header">
          <h3>{{ isEditing ? 'Edit Room Details' : 'Add New Room' }}</h3>
          <button @click="closeModal" class="close-btn">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">Room Category</label>
            <select v-model="form.category_id" class="form-select">
              <option v-for="category in roomCategories" :key="category.id" :value="category.id">{{ category.name }}</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Room Type</label>
            <select v-model="form.type_id" class="form-select">
              <option v-for="type in roomTypes" :key="type.id" :value="type.id">{{ type.name }}</option>
            </select>
          </div>
       
          <div class="form-group">
            <label class="form-label">Price per Night (Rp)</label>
            <input 
              v-model.number="form.price" 
              type="number" 
              min="0" 
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label class="form-label">Capacity (persons)</label>
            <input 
              v-model.number="form.capacity" 
              type="number" 
              min="1" 
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label class="form-label">Room Size (sqm)</label>
            <input 
              v-model.number="form.size" 
              type="number" 
              min="1" 
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label class="form-label">Description</label>
            <textarea 
              v-model="form.description" 
              class="form-input" 
              rows="2" 
              placeholder="Describe this room..."
            ></textarea>
          </div>
          <div class="form-group">
            <label class="form-label">Room Image File</label>
            <input 
              type="file" 
              @change="handleImageUpload" 
              class="form-input" 
              accept="image/*"
            />
          </div>
          <!-- <div class="form-group">
            <label class="form-label">Status</label>
            <select v-model="form.status" class="form-select">
              <option v-for="stat in roomStatusesOptions" :key="stat" :value="stat">{{ stat }}</option>
            </select>
          </div> -->
        </div>
        <div class="modal-footer">
          <button @click="closeModal" class="btn btn-secondary">Cancel</button>
          <button @click="saveRoom" class="btn btn-primary" :disabled="loading">Save Changes</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.rooms-view {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
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
.text-info { color: #3b82f6; }
.text-warning { color: #f59e0b; }

.control-bar {
  display: flex;
  gap: 16px;
  background-color: #ffffff;
  padding: 16px;
  border-radius: 16px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  flex-wrap: wrap;
  align-items: center;
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

.badge {
  padding: 4px 8px;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  display: inline-block;
}

.badge-success { background-color: #d1fae5; color: #065f46; }
.badge-info { background-color: #dbeafe; color: #1e40af; }
.badge-warning { background-color: #fef3c7; color: #92400e; }

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

.btn-warning-outline {
  border-color: #f59e0b;
  color: #f59e0b;
  background: transparent;
}
.btn-warning-outline:hover {
  background-color: #fffbeb;
}

.btn-primary-outline {
  border-color: #cbd5e1;
  color: #334155;
  background: transparent;
}
.btn-primary-outline:hover {
  background-color: #f1f5f9;
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

.room-thumbnail {
  width: 60px;
  height: 40px;
  object-fit: cover;
  border-radius: 6px;
  border: 1px solid #cbd5e1;
}

.text-muted {
  color: #94a3b8;
  font-size: 0.8rem;
  font-style: italic;
}
</style>
