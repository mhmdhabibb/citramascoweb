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

// Drawer & Modal State
const selectedRoom = ref(null)
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
  status: 'Available',
})

// Stats
const totalRooms = computed(() => rooms.value.length)
const availableRooms = computed(
  () =>
    rooms.value.filter(
      (r) =>
        (r.status || 'Available').toLowerCase() === 'active' ||
        (r.status || 'Available').toLowerCase() === 'available',
    ).length,
)
const occupiedRooms = computed(() => rooms.value.filter((r) => r.status === 'Occupied').length)
const maintenanceRooms = computed(
  () => rooms.value.filter((r) => r.status?.toLowerCase() === 'maintenance').length,
)

// Filter Logic
const filteredRooms = computed(() => {
  return rooms.value.filter((room) => {
    const code = room.code || ''
    const name = room.name || ''
    const typeName = room.type?.name || ''
    const categoryName = room.category?.name || ''

    const matchesSearch =
      code.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      typeName.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      categoryName.toLowerCase().includes(searchQuery.value.toLowerCase())

    const status = room.status || 'Available'
    const matchesStatus =
      statusFilter.value === 'All' || status.toLowerCase() === statusFilter.value.toLowerCase()
    return matchesSearch && matchesStatus
  })
})

const handleImageUpload = (event) => {
  imageFile.value = event.target.files[0]
}

const selectRow = (room) => {
  selectedRoom.value = room
}

const closeDrawer = () => {
  selectedRoom.value = null
}

// --- AKSI BARU: Mengubah Status Fisik Kamar (Active, Maintenance, Dirty) ---
const mutateRoomStatus = async (id, targetStatus) => {
  try {
    loading.value = true
    const msg = await roomService.updateStatus(id, targetStatus)
    toastStore.success(msg || `Status kamar berhasil diubah ke ${targetStatus}`)

    // Refresh Stream Data
    rooms.value = await roomService.getAll()

    // Sync status terbaru di dalam Drawer peninjauan aktif
    if (selectedRoom.value && selectedRoom.value.id === id) {
      selectedRoom.value.status = targetStatus
    }
  } catch (error) {
    console.error(error)
    toastStore.error(error.message || 'Gagal memperbarui operasional kamar')
  } finally {
    loading.value = false
  }
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
    status: 'Available',
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
    status: item.status || 'Available',
  }
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
}

const saveRoom = async () => {
  if (!form.value.category_id) {
    toastStore.warning('Mohon pilih kategori kamar.')
    return
  }
  if (!form.value.type_id) {
    toastStore.warning('Mohon pilih tipe kamar.')
    return
  }

  const selectedCategory = roomCategories.value.find((c) => c.id === form.value.category_id)
  const roomName = selectedCategory ? selectedCategory.name : ''
  form.value.name = roomName

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
    formData.append('status', form.value.status)

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

    rooms.value = await roomService.getAll()
    closeModal()
    closeDrawer()
  } catch (error) {
    console.error(error)
    toastStore.error(
      error.response?.data?.message || error.message || 'Gagal menyimpan detail kamar',
    )
  } finally {
    loading.value = false
  }
}

const deleteRoom = async (id) => {
  if (confirm(`Apakah Anda yakin ingin menghapus kamar ini?`)) {
    try {
      loading.value = true
      const msg = await roomService.delete(id)
      toastStore.success(msg || 'Kamar berhasil dihapus!')
      rooms.value = await roomService.getAll()
      closeDrawer()
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
      typeService.getAll(),
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
   

    <div class="stats-grid">
      <div class="stat-card">
        <span class="stat-icon">🚪</span>
        <div>
          <h3>Total Rooms</h3>
          <p class="main-val">{{ totalRooms }}</p>
        </div>
      </div>
      <div class="stat-card">
        <span class="stat-icon text-success">🟢</span>
        <div>
          <h3>Available</h3>
          <p class="main-val text-success">{{ availableRooms }}</p>
        </div>
      </div>
      <div class="stat-card">
        <span class="stat-icon text-info">🔵</span>
        <div>
          <h3>Occupied</h3>
          <p class="main-val text-info">{{ occupiedRooms }}</p>
        </div>
      </div>
      <div class="stat-card">
        <span class="stat-icon text-warning">🛠️</span>
        <div>
          <h3>Maintenance</h3>
          <p class="main-val text-warning">{{ maintenanceRooms }}</p>
        </div>
      </div>
    </div>

    <div class="split-pane-layout">
      <div class="table-pane-box">
        <div class="control-bar">
          <div class="search-box">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Cari berdasarkan kode kamar, tipe, atau kategori..."
              class="search-input"
            />
          </div>
          <div class="filter-box">
            <select v-model="statusFilter" class="filter-select">
              <option value="All">Status</option>
              <option value="active">Active / Available</option>
              <option value="Occupied">Occupied</option>
              <option value="maintenance">Maintenance</option>
              <option value="dirty">Dirty</option>
            </select>
          </div>
          <button @click="openCreateModal" class="btn btn-primary">+ Add Room</button>
        </div>

        <div class="dashboard-card table-card">
          <div class="responsive-table-wrap">
            <table class="premium-table">
              <thead>
                <tr>
                  <th>Room Code</th>
                  <th>Category / Name</th>
                  <th>Type</th>
                  <th>Price / Night</th>
                  <th>Image</th>
                  <th class="text-center">Status</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="room in filteredRooms"
                  :key="room.id"
                  @click="selectRow(room)"
                  :class="{ 'selected-row': selectedRoom?.id === room.id }"
                  class="clickable-row"
                >
                  <td class="bold-code">{{ room.code }}</td>
                  <td>{{ room.name }}</td>
                  <td>
                    <span class="type-pill">{{ room.type?.name || 'N/A' }}</span>
                  </td>
                  <td class="price-text">
                    Rp {{ room.price ? room.price.toLocaleString('id-ID') : 0 }}
                  </td>
                  <td>
                    <img
                      v-if="room.image"
                      :src="room.image"
                      alt="Room Image"
                      class="room-thumbnail"
                    />
                    <span v-else class="text-muted">No Image</span>
                  </td>
                  <td class="text-center">
                    <span
                      class="status-dot-badge"
                      :class="{
                        'status-available':
                          (room.status || 'Available').toLowerCase() === 'active' ||
                          (room.status || 'Available').toLowerCase() === 'available',
                        'status-occupied': room.status === 'Occupied',
                        'status-maintenance': room.status?.toLowerCase() === 'maintenance',
                        'status-dirty': room.status?.toLowerCase() === 'dirty',
                      }"
                    >
                      {{ room.status || 'Available' }}
                    </span>
                  </td>
                </tr>
                <tr v-if="filteredRooms.length === 0">
                  <td colspan="6" class="no-data">Tidak ditemukan unit kamar yang cocok.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <transition name="drawer-slide">
        <div v-if="selectedRoom" class="action-drawer-pane">
          <div class="drawer-header">
            <div>
              <h3>Room Details</h3>
              <span class="drawer-id">Room Code: {{ selectedRoom.code }}</span>
            </div>
            <button @click="closeDrawer" class="close-drawer-btn">✕</button>
          </div>

          <div class="drawer-content">
            <div class="drawer-img-showcase">
              <img v-if="selectedRoom.image" :src="selectedRoom.image" alt="Room Preview" />
              <div v-else class="no-img-placeholder">🚫 No Room Image Uploaded</div>
            </div>

            <div class="info-block-card">
              <label>Category / Room Name</label>
              <p class="val-large">{{ selectedRoom.name }}</p>
            </div>

            <div class="info-grid-2">
              <div class="info-block-card">
                <label>Room Type</label>
                <p class="val-mid">{{ selectedRoom.type?.name || 'N/A' }}</p>
              </div>
              <div class="info-block-card">
                <label>Current Status</label>
                <div style="margin-top: 4px">
                  <span
                    class="status-dot-badge"
                    :class="{
                      'status-available':
                        (selectedRoom.status || 'Available').toLowerCase() === 'active' ||
                        (selectedRoom.status || 'Available').toLowerCase() === 'available',
                      'status-occupied': selectedRoom.status === 'Occupied',
                      'status-maintenance': selectedRoom.status?.toLowerCase() === 'maintenance',
                      'status-dirty': selectedRoom.status?.toLowerCase() === 'dirty',
                    }"
                    >{{ selectedRoom.status || 'Available' }}</span
                  >
                </div>
              </div>
            </div>

            <div class="info-block-card oper-card">
              <label>Ubah Status Operasional Fisik</label>
              <div class="oper-buttons-group">
                <button
                  @click="mutateRoomStatus(selectedRoom.id, 'active')"
                  class="op-btn op-active"
                  :disabled="selectedRoom.status?.toLowerCase() === 'active'"
                >
                  Set Active
                </button>
                <button
                  @click="mutateRoomStatus(selectedRoom.id, 'dirty')"
                  class="op-btn op-dirty"
                  :disabled="selectedRoom.status?.toLowerCase() === 'dirty'"
                >
                  Set Dirty
                </button>
                <button
                  @click="mutateRoomStatus(selectedRoom.id, 'maintenance')"
                  class="op-btn op-maint"
                  :disabled="selectedRoom.status?.toLowerCase() === 'maintenance'"
                >
                  Set Maintenance
                </button>
              </div>
            </div>

            <div class="info-grid-2">
              <div class="info-block-card">
                <label>Max Capacity</label>
                <p class="val-small">👥 {{ selectedRoom.capacity || 2 }} Persons</p>
              </div>
              <div class="info-block-card">
                <label>Room Size</label>
                <p class="val-small">📐 {{ selectedRoom.size || 30 }} m²</p>
              </div>
            </div>

            <div class="info-block-card">
              <label>Description / Features</label>
              <p class="val-desc">
                {{ selectedRoom.description || 'Tidak ada deskripsi tambahan.' }}
              </p>
            </div>

            <div class="info-block-card pricing-bg">
              <label>Rate per Night</label>
              <p class="val-price">Rp {{ selectedRoom.price?.toLocaleString('id-ID') || 0 }}</p>
            </div>

            <div class="drawer-actions-container">
              <button @click="openEditModal(selectedRoom)" class="btn btn-checkedin btn-block">
                Edit Room Specification
              </button>
              <button @click="deleteRoom(selectedRoom.id)" class="btn btn-danger-outline btn-block">
                Delete Unit From Inventory
              </button>
            </div>
          </div>
        </div>
      </transition>
    </div>

    <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
      <div class="modal-card">
        <div class="modal-header">
          <h3>{{ isEditing ? 'Edit Room Details' : 'Add New Room Unit' }}</h3>
          <button @click="closeModal" class="close-btn">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">Room Category</label>
            <select v-model="form.category_id" class="form-select">
              <option v-for="category in roomCategories" :key="category.id" :value="category.id">
                {{ category.name }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Room Type</label>
            <select v-model="form.type_id" class="form-select">
              <option v-for="type in roomTypes" :key="type.id" :value="type.id">
                {{ type.name }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Price per Night (Rp)</label>
            <input v-model.number="form.price" type="number" min="0" class="form-input" />
          </div>
          <div class="form-grid-2">
            <div class="form-group">
              <label class="form-label">Capacity (persons)</label>
              <input v-model.number="form.capacity" type="number" min="1" class="form-input" />
            </div>
            <div class="form-group">
              <label class="form-label">Room Size (m²)</label>
              <input v-model.number="form.size" type="number" min="1" class="form-input" />
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">Description</label>
            <textarea
              v-model="form.description"
              class="form-input"
              rows="2"
              placeholder="Tulis deskripsi fasilitas kamar..."
            ></textarea>
          </div>
          <div class="form-group">
            <label class="form-label">Room Image File</label>
            <input
              type="file"
              @change="handleImageUpload"
              class="form-input-file"
              accept="image/*"
            />
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeModal" class="btn btn-secondary">Cancel</button>
          <button @click="saveRoom" class="btn btn-save" :disabled="loading">Save Changes</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Main View Blueprint */
.rooms-view {
  display: flex;
  flex-direction: column;
  gap: 24px;
  background-color: #f8fafc;
  font-family: 'Plus Jakarta Sans', sans-serif;
}


.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}
.stat-card {
  background: #ffffff;
  border-radius: 16px;
  padding: 18px 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  border: 1px solid #e2e8f0;
}
.stat-card h3 {
  font-size: 0.8rem;
  color: #64748b;
  margin: 0 0 2px 0;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}
.stat-icon {
  font-size: 1.3rem;
  padding: 10px;
  background: #f8fafc;
  border-radius: 12px;
}
.main-val {
  font-size: 1.6rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.split-pane-layout {
  display: flex;
  gap: 24px;
  align-items: start;
  position: relative;
}
.table-pane-box {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  width: 100%;
  min-width: 0;
}

.control-bar {
  display: flex;
  gap: 14px;
  background-color: #ffffff;
  padding: 14px;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  align-items: center;
}
.search-box {
  flex-grow: 1;
}
.search-input {
  width: 100%;
  padding: 10px 16px;
  border-radius: 10px;
  border: 1px solid #e2e8f0;
  outline: none;
  font-size: 0.875rem;
  color: #1e293b;
  background: #f8fafc;
  transition: all 0.2s;
  box-sizing: border-box;
}
.search-input:focus {
  border-color: #e4793b;
  background: #ffffff;
  box-shadow: 0 0 0 3px rgba(228, 121, 59, 0.08);
}
.filter-select {
  padding: 10px 36px 10px 16px;
  border-radius: 10px;
  border: 1px solid #e2e8f0;
  outline: none;
  font-size: 0.875rem;
  color: #475569;
  background-color: #ffffff;
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' stroke='%23475569'%3E%3Cpath stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='M19 9l-7 7-7-7'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 14px center;
  background-size: 14px;
}

.table-card {
  background: #ffffff;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.01);
}
.responsive-table-wrap {
  overflow-x: auto;
  overflow-y: auto;
  max-height: 520px;
}
.premium-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}
.premium-table th {
  position: sticky;
  top: 0;
  z-index: 2;
  padding: 14px 20px;
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  border-bottom: 1px solid #edf2f7;
  background: #f8fafc;
  letter-spacing: 0.05em;
}
.premium-table td {
  padding: 14px 20px;
  font-size: 0.875rem;
  color: #334155;
  border-bottom: 1px solid #f1f5f9;
  vertical-align: middle;
}
.clickable-row {
  cursor: pointer;
  transition: background 0.15s ease;
}
.clickable-row:hover {
  background-color: #f8fafc;
}
.selected-row {
  background-color: #fff7f2 !important;
}
.selected-row td:first-child {
  box-shadow: inset 4px 0 0 #e4793b;
}

.bold-code {
  font-family: monospace;
  font-weight: 700;
  color: #0f172a;
}
.type-pill {
  background: #f1f5f9;
  padding: 4px 10px;
  border-radius: 6px;
  font-weight: 700;
  font-size: 0.8rem;
  color: #475569;
}
.price-text {
  font-weight: 700;
  color: #0f172a;
}
.room-thumbnail {
  width: 54px;
  height: 36px;
  object-fit: cover;
  border-radius: 6px;
  border: 1px solid #e2e8f0;
}

/* Status Badges */
.status-dot-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: capitalize;
}
.status-available {
  background-color: #d1fae5;
  color: #065f46;
}
.status-occupied {
  background-color: #e0f2fe;
  color: #0369a1;
}
.status-maintenance {
  background-color: #fee2e2;
  color: #991b1b;
}
.status-dirty {
  background-color: #fef3c7;
  color: #92400e;
}

/* Action Drawer Sidebar Panel */
.action-drawer-pane {
  width: 380px;
  min-width: 380px;
  background: #ffffff;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  padding: 24px;
  position: sticky;
  top: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  box-sizing: border-box;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.04);
}
.drawer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #f1f5f9;
  padding-bottom: 16px;
}
.drawer-header h3 {
  font-size: 1.1rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
  letter-spacing: -0.01em;
}
.drawer-id {
  font-size: 0.75rem;
  color: #94a3b8;
  font-weight: 600;
  display: block;
  margin-top: 4px;
}
.close-drawer-btn {
  background: #f1f5f9;
  border: none;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  font-size: 0.8rem;
  color: #64748b;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}
.close-drawer-btn:hover {
  background: #e2e8f0;
  color: #0f172a;
}

.drawer-content {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.drawer-img-showcase {
  width: 100%;
  height: 160px;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
}
.drawer-img-showcase img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.no-img-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8fafc;
  color: #94a3b8;
  font-size: 0.85rem;
  font-weight: 500;
  font-style: italic;
}

.info-block-card {
  background: #f8fafc;
  border: 1px solid #f1f5f9;
  padding: 14px 16px;
  border-radius: 12px;
}
.info-block-card label {
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  color: #94a3b8;
  letter-spacing: 0.04em;
  display: block;
  margin-bottom: 4px;
}
.val-large {
  font-size: 1.25rem;
  font-weight: 700;
  color: #0f172a;
  margin: 0;
}
.info-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}
.val-mid {
  font-size: 0.95rem;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}
.val-small {
  font-size: 0.85rem;
  font-weight: 600;
  color: #334155;
  margin: 0;
}
.val-desc {
  font-size: 0.85rem;
  color: #475569;
  margin: 0;
  line-height: 1.5;
}

/* CSS BARU KELOMPOK TOMBOL MUTASI STATUS OPERASIONAL */
.oper-card {
  border: 1px solid #e2e8f0;
  background: #ffffff !important;
}
.oper-buttons-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 8px;
}
.op-btn {
  width: 100%;
  padding: 8px;
  border-radius: 8px;
  font-size: 0.8rem;
  font-weight: 700;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s;
}
.op-active {
  background-color: #d1fae5;
  color: #065f46;
}
.op-active:hover:not(:disabled) {
  background-color: #a7f3d0;
}
.op-dirty {
  background-color: #fef3c7;
  color: #92400e;
}
.op-dirty:hover:not(:disabled) {
  background-color: #fde68a;
}
.op-maint {
  background-color: #fee2e2;
  color: #991b1b;
}
.op-maint:hover:not(:disabled) {
  background-color: #fca5a5;
}
.op-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  border: 1px dashed #cbd5e1;
}

.pricing-bg {
  background-color: #fff7ed;
  border-color: #ffedd5;
}
.pricing-bg label {
  color: #c2410c;
}
.val-price {
  font-size: 1.35rem;
  font-weight: 800;
  color: #ea580c;
  margin: 0;
}

.drawer-actions-container {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 10px;
}
.btn-block {
  width: 100%;
}
.btn {
  padding: 11px 16px;
  border-radius: 10px;
  font-size: 0.85rem;
  font-weight: 700;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s ease;
  text-align: center;
}
.btn-primary {
  background-color: #e4793b;
  color: white;
}
.btn-primary:hover {
  background-color: #c8632b;
}
.btn-checkedin {
  background-color: #3b82f6;
  color: white;
}
.btn-checkedin:hover {
  background-color: #2563eb;
}
.btn-secondary {
  border-color: #cbd5e1;
  color: #334155;
  background-color: #ffffff;
}
.btn-secondary:hover {
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
  padding: 48px !important;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(15, 23, 42, 0.2);
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
  max-width: 460px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.08);
  display: flex;
  flex-direction: column;
  max-height: 90vh;
  overflow-y: auto;
}
.modal-header {
  padding: 18px 24px;
  border-bottom: 1px solid #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.modal-header h3 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 800;
  color: #0f172a;
}
.close-btn {
  background: none;
  border: none;
  font-size: 1.4rem;
  color: #94a3b8;
  cursor: pointer;
}
.close-btn:hover {
  color: #0f172a;
}
.modal-body {
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.form-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}
.form-label {
  font-size: 0.72rem;
  font-weight: 700;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
.form-input,
.form-select {
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
.form-input:focus,
.form-select:focus {
  border-color: #e4793b;
  box-shadow: 0 0 0 3px rgba(228, 121, 59, 0.08);
}
.form-input-file {
  font-size: 0.8rem;
  color: #64748b;
}
.modal-footer {
  padding: 14px 24px;
  background-color: #f8fafc;
  border-top: 1px solid #f1f5f9;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
.btn-save {
  background-color: #e4793b;
  color: white;
}
.btn-save:hover {
  background-color: #c8632b;
}

.responsive-table-wrap::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}
.responsive-table-wrap::-webkit-scrollbar-track {
  background: #f1f5f9;
}
.responsive-table-wrap::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 4px;
}
.responsive-table-wrap::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

.drawer-slide-enter-active,
.drawer-slide-leave-active {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}
.drawer-slide-enter-from,
.drawer-slide-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}
</style>
