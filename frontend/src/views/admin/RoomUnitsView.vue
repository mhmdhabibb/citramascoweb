<script setup lang="ts">
import { floorService } from '@/services/admin/floorService'
import { roomUnitService } from '@/services/admin/roomUnitService'
import { roomService } from '@/services/roomService'
import { useAuthStore } from '@/stores/authStore'
import { canMutateData } from '@/config/access'
import { useToastStore } from '@/stores/toastStore'
import type { Floor, Room, RoomUnit } from '@/types'
import { computed, onMounted, ref } from 'vue'

const authStore = useAuthStore()

const units = ref<RoomUnit[]>([])
const floors = ref<Floor[]>([])
const rooms = ref<Room[]>([])
const loading = ref(false)
const toastStore = useToastStore()

// Search & Filter
const searchQuery = ref('')
const floorFilter = ref('All')
const typeFilter = ref('All')
const statusFilter = ref('All')

// Modal state
const isModalOpen = ref(false)
const isEditing = ref(false)
const editingId = ref('')
const form = ref({
  room_number: '',
  room_id: '',
  floor_id: '',
  status: 'available',
})

// Preview detail tipe kamar yang dipilih
const selectedRoomDetails = computed(() => {
  return rooms.value.find((r) => r.id === form.value.room_id)
})

// Stats
const totalUnits = computed(() => units.value.length)
const availableUnits = computed(() => units.value.filter((u) => u.status === 'available').length)
const occupiedUnits = computed(() => units.value.filter((u) => u.status === 'occupied').length)
const dirtyUnits = computed(() => units.value.filter((u) => u.status === 'dirty').length)
const maintenanceUnits = computed(() => units.value.filter((u) => u.status === 'maintenance').length)

const fetchData = async () => {
  try {
    loading.value = true
    const [uData, fData, rData] = await Promise.all([
      roomUnitService.getAll(),
      floorService.getAll(),
      roomService.getAll(),
    ])
    units.value = uData
    floors.value = fData
    rooms.value = rData
  } catch (error: any) {
    console.error('Error fetching room units data:', error)
    toastStore.error(error.message || 'Failed to load room unit data')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})

const filteredUnits = computed(() => {
  return units.value.filter((unit) => {
    const roomNum = unit.room_number || ''
    const roomName = unit.room?.name || ''
    const floorName = unit.floor?.name || ''

    const matchesSearch =
      roomNum.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      roomName.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      floorName.toLowerCase().includes(searchQuery.value.toLowerCase())

    const matchesFloor = floorFilter.value === 'All' || unit.floor_id === floorFilter.value
    const matchesType = typeFilter.value === 'All' || unit.room_id === typeFilter.value
    const matchesStatus = statusFilter.value === 'All' || unit.status === statusFilter.value

    return matchesSearch && matchesFloor && matchesType && matchesStatus
  })
})

const openCreateModal = () => {
  isEditing.value = false
  editingId.value = ''
  form.value = {
    room_number: '',
    room_id: rooms.value[0]?.id || '',
    floor_id: floors.value[0]?.id || '',
    status: 'available',
  }
  isModalOpen.value = true
}

const openEditModal = (item: RoomUnit) => {
  isEditing.value = true
  editingId.value = item.id
  form.value = {
    room_number: item.room_number,
    room_id: item.room_id,
    floor_id: item.floor_id,
    status: item.status || 'available',
  }
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
}

const saveRoomUnit = async () => {
  if (!form.value.room_number.trim()) {
    toastStore.warning('Room number cannot be empty!')
    return
  }
  if (!form.value.room_id) {
    toastStore.warning('Pilih tipe katalog kamar!')
    return
  }
  if (!form.value.floor_id) {
    toastStore.warning('Pilih lantai kamar!')
    return
  }

  try {
    loading.value = true
    if (isEditing.value) {
      const msg = await roomUnitService.update(editingId.value, form.value)
      toastStore.success(msg || 'Unit kamar berhasil diperbarui!')
    } else {
      const msg = await roomUnitService.create(form.value)
      toastStore.success(msg || 'Unit kamar berhasil ditambahkan!')
    }
    await fetchData()
    closeModal()
  } catch (error: any) {
    console.error('Error saving room unit:', error)
    toastStore.error(error.response?.data?.message || error.message || 'Failed to save room unit')
  } finally {
    loading.value = false
  }
}

const updateStatusQuick = async (id: string, status: string) => {
  try {
    loading.value = true
    const msg = await roomUnitService.updateStatus(id, status)
    toastStore.success(msg || `Room status changed to ${status}!`)
    await fetchData()
  } catch (error: any) {
    console.error('Error updating status:', error)
    toastStore.error(error.response?.data?.message || error.message || 'Failed to change room status')
  } finally {
    loading.value = false
  }
}

const deleteRoomUnit = async (id: string, roomNumber: string) => {
  if (!confirm(`Are you sure you want to delete Room ${roomNumber}?`)) return

  try {
    loading.value = true
    const msg = await roomUnitService.delete(id)
    toastStore.success(msg || 'Unit kamar berhasil dihapus!')
    await fetchData()
  } catch (error: any) {
    console.error('Error deleting room unit:', error)
    toastStore.error(error.response?.data?.message || error.message || 'Failed to delete room unit')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="room-units-view">
    <!-- Header Section -->
    <div class="header-section">
      <div>
        <h1 class="page-title">Room Unit Master</h1>
        <p class="subtitle">Manage hotel physical room numbers, floor placement, and status.</p>
      </div>
      <button v-if="canMutateData(authStore.role)" @click="openCreateModal" class="btn-primary">
        <span>+</span> Add Room No.
      </button>
      <div v-else class="px-3.5 py-2 rounded-xl bg-slate-900 text-amber-300 text-xs font-bold border border-indigo-900/50 flex items-center gap-2 shadow-sm">
        <span>👑</span> Manager Review Mode (Read-Only)
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">🚪</div>
        <div>
          <h3>Total Room Units</h3>
          <p class="main-val">{{ totalUnits }} Unit</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon text-emerald-600 bg-emerald-50">✨</div>
        <div>
          <h3>Available</h3>
          <p class="main-val text-emerald-600">{{ availableUnits }}</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon text-indigo-600 bg-indigo-50">🛏️</div>
        <div>
          <h3>Occupied</h3>
          <p class="main-val text-indigo-600">{{ occupiedUnits }}</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon text-amber-600 bg-amber-50">🧹</div>
        <div>
          <h3>Dirty (Cleaning)</h3>
          <p class="main-val text-amber-600">{{ dirtyUnits }}</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon text-rose-600 bg-rose-50">🔧</div>
        <div>
          <h3>Maintenance (Broken)</h3>
          <p class="main-val text-rose-600">{{ maintenanceUnits }}</p>
        </div>
      </div>
    </div>

    <!-- Table Pane -->
    <div class="table-pane-box">
      <div class="control-bar flex-wrap">
        <div class="search-box">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search room no., type, floor..."
            class="search-input"
          />
        </div>

        <div class="filter-box flex items-center gap-2">
          <select v-model="floorFilter" class="filter-select">
            <option value="All">All Floors</option>
            <option v-for="fl in floors" :key="fl.id" :value="fl.id">{{ fl.name }}</option>
          </select>

          <select v-model="typeFilter" class="filter-select">
            <option value="All">All Room Types</option>
            <option v-for="rm in rooms" :key="rm.id" :value="rm.id">{{ rm.name }}</option>
          </select>

          <select v-model="statusFilter" class="filter-select">
            <option value="All">All Statuses</option>
            <option value="available">Available</option>
            <option value="occupied">Occupied</option>
            <option value="dirty">Dirty</option>
            <option value="maintenance">Maintenance</option>
          </select>
        </div>
      </div>

      <div class="dashboard-card table-card">
        <div class="responsive-table-wrap">
          <table class="premium-table">
            <thead>
              <tr>
                <th class="w-28">No. Kamar</th>
                <th>Floor</th>
                <th>Room Category</th>
                <th>Bed Type</th>
                <th>Price / Night</th>
                <th class="text-center">Status</th>
                <th v-if="canMutateData(authStore.role)" class="text-center w-48">Change Status</th>
                <th v-if="canMutateData(authStore.role)" class="text-center w-28">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="unit in filteredUnits" :key="unit.id" class="data-row">
                <td class="bold-name">
                  <span class="room-number-badge">
                    🚪 {{ unit.room_number }}
                  </span>
                </td>
                <td>
                  <span class="px-2 py-0.5 rounded text-xs font-semibold bg-slate-100 text-slate-700">
                    {{ unit.floor?.name || '-' }}
                  </span>
                </td>
                <td>
                  <span class="font-extrabold text-slate-900 text-xs">
                    {{ unit.room?.category?.name || unit.room?.name || '-' }}
                  </span>
                </td>
                <td>
                  <span class="px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-indigo-50 text-indigo-700 border border-indigo-100">
                    🛏️ {{ unit.room?.type?.name || 'Standard' }}
                  </span>
                </td>
                <td class="text-xs font-bold text-emerald-700">
                  Rp {{ (unit.room?.price || 0).toLocaleString('id-ID') }}
                </td>
                <td class="text-center">
                  <span
                    class="status-pill"
                    :class="{
                      'status-available': unit.status === 'available',
                      'status-occupied': unit.status === 'occupied',
                      'status-dirty': unit.status === 'dirty',
                      'status-maintenance': unit.status === 'maintenance',
                    }"
                  >
                    {{ unit.status?.toUpperCase() }}
                  </span>
                </td>
                <td v-if="canMutateData(authStore.role)" class="text-center">
                  <div class="flex items-center justify-center gap-1">
                    <button
                      @click="updateStatusQuick(unit.id, 'available')"
                      class="quick-btn"
                      :class="unit.status === 'available' ? 'quick-btn-active' : ''"
                      title="Set Available"
                    >
                      ✨
                    </button>
                    <button
                      @click="updateStatusQuick(unit.id, 'occupied')"
                      class="quick-btn"
                      :class="unit.status === 'occupied' ? 'quick-btn-active' : ''"
                      title="Set Occupied"
                    >
                      🛏️
                    </button>
                    <button
                      @click="updateStatusQuick(unit.id, 'dirty')"
                      class="quick-btn"
                      :class="unit.status === 'dirty' ? 'quick-btn-active' : ''"
                      title="Set Dirty (Cleaning)"
                    >
                      🧹
                    </button>
                    <button
                      @click="updateStatusQuick(unit.id, 'maintenance')"
                      class="quick-btn"
                      :class="unit.status === 'maintenance' ? 'quick-btn-active' : ''"
                      title="Set Maintenance"
                    >
                      🔧
                    </button>
                  </div>
                </td>
                <td v-if="canMutateData(authStore.role)" class="text-center">
                  <div class="flex items-center justify-center gap-1.5">
                    <button
                      @click="openEditModal(unit)"
                      class="px-2.5 py-1 text-xs font-semibold text-indigo-600 bg-indigo-50 hover:bg-indigo-100 rounded-lg transition-colors cursor-pointer"
                    >
                      Edit
                    </button>
                    <button
                      @click="deleteRoomUnit(unit.id, unit.room_number)"
                      class="px-2.5 py-1 text-xs font-semibold text-rose-600 bg-rose-50 hover:bg-rose-100 rounded-lg transition-colors cursor-pointer"
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="filteredUnits.length === 0">
                <td :colspan="canMutateData(authStore.role) ? 8 : 6" class="no-data">No matching room unit found.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Modal Form -->
    <div v-if="isModalOpen" class="modal-backdrop">
      <div class="modal-card">
        <div class="modal-header">
          <h3>{{ isEditing ? 'Edit Room Unit' : 'Add New Room Unit' }}</h3>
          <button @click="closeModal" class="close-btn">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>Physical Room Number *</label>
            <input
              v-model="form.room_number"
              type="text"
              placeholder="e.g.: 101, 102, 402"
              class="form-input font-bold text-base"
              required
            />
          </div>

          <div class="form-group">
            <label>Room Type & Category *</label>
            <select v-model="form.room_id" class="form-input" required>
              <option value="" disabled>Select Room Type & Category</option>
              <option v-for="rm in rooms" :key="rm.id" :value="rm.id">
                {{ rm.category?.name || rm.name }} — [{{ rm.type?.name || 'Standard Bed' }}] (Rp {{ (rm.price || 0).toLocaleString('id-ID') }})
              </option>
            </select>
          </div>

          <!-- Live Selected Room Details Card -->
          <div v-if="selectedRoomDetails" class="p-3 bg-gradient-to-br from-indigo-50/90 to-slate-50 border border-indigo-100 rounded-xl space-y-1.5 text-xs">
            <div class="flex justify-between items-center">
              <span class="text-slate-500 font-semibold">🏷️ Category:</span>
              <strong class="text-indigo-950 font-extrabold">{{ selectedRoomDetails.category?.name || selectedRoomDetails.name }}</strong>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-slate-500 font-semibold">🛏️ Bed Type:</span>
              <span class="font-bold text-indigo-700 bg-white px-2 py-0.5 rounded border border-indigo-100">
                {{ selectedRoomDetails.type?.name || 'Standard Bed' }}
              </span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-slate-500 font-semibold">👥 Kapasitas:</span>
              <span class="font-bold text-slate-700">
                {{ selectedRoomDetails.capacity || 2 }} Dewasa {{ selectedRoomDetails.child_capacity ? `+ ${selectedRoomDetails.child_capacity} Anak` : '' }}
              </span>
            </div>
            <div class="flex justify-between items-center pt-1.5 border-t border-indigo-100/70">
              <span class="text-slate-500 font-semibold">💰 Standard Price:</span>
              <strong class="text-emerald-700 font-black text-sm">
                Rp {{ (selectedRoomDetails.price || 0).toLocaleString('id-ID') }} <span class="text-[10px] font-normal text-slate-400">/ malam</span>
              </strong>
            </div>
          </div>

          <div class="form-group">
            <label>Floor Placement *</label>
            <select v-model="form.floor_id" class="form-input" required>
              <option value="" disabled>Select Floor</option>
              <option v-for="fl in floors" :key="fl.id" :value="fl.id">
                {{ fl.name }} (Floor {{ fl.floor_number }})
              </option>
            </select>
          </div>

          <div class="form-group">
            <label>Status Awal</label>
            <select v-model="form.status" class="form-input">
              <option value="available">Available</option>
              <option value="occupied">Occupied</option>
              <option value="dirty">Dirty (Needs Cleaning)</option>
              <option value="maintenance">Maintenance (Broken)</option>
            </select>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeModal" class="btn-secondary">Cancel</button>
          <button @click="saveRoomUnit" :disabled="loading" class="btn-primary">
            {{ loading ? 'Saving...' : 'Save' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.room-units-view {
  display: flex;
  flex-direction: column;
  gap: 24px;
  background-color: #f8fafc;
  font-family: 'Plus Jakarta Sans', sans-serif;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
}

.page-title {
  font-size: 22px;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.subtitle {
  font-size: 13px;
  color: #64748b;
  margin: 4px 0 0 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(170px, 1fr));
  gap: 14px;
}

.stat-card {
  background: #ffffff;
  padding: 14px 18px;
  border-radius: 16px;
  border: 1px solid #f1f5f9;
  display: flex;
  align-items: center;
  gap: 14px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.02);
}

.stat-icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: #f8fafc;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.stat-card h3 {
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  margin: 0;
}

.stat-card .main-val {
  font-size: 17px;
  font-weight: 800;
  color: #0f172a;
  margin: 2px 0 0 0;
}

.table-pane-box {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.control-bar {
  display: flex;
  gap: 12px;
}

.search-box {
  flex: 1;
  min-width: 260px;
}

.search-input,
.filter-select {
  padding: 10px 14px;
  font-size: 13px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  outline: none;
  transition: border-color 0.2s;
}

.search-input:focus,
.filter-select:focus {
  border-color: #6366f1;
}

.dashboard-card {
  background: #ffffff;
  border-radius: 16px;
  border: 1px solid #f1f5f9;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.02);
}

.responsive-table-wrap {
  overflow-x: auto;
}

.premium-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.premium-table th {
  padding: 14px 16px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #64748b;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
}

.premium-table td {
  padding: 14px 16px;
  font-size: 13px;
  border-bottom: 1px solid #f1f5f9;
}

.data-row:hover {
  background-color: #f8fafc;
}

.room-number-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-weight: 800;
  font-size: 13px;
  color: #0f172a;
  background: #f1f5f9;
  padding: 4px 8px;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
}

.status-pill {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 9999px;
  font-size: 11px;
  font-weight: 800;
}

.status-available {
  background: #ecfdf5;
  color: #059669;
  border: 1px solid #a7f3d0;
}

.status-occupied {
  background: #eef2ff;
  color: #4f46e5;
  border: 1px solid #c7d2fe;
}

.status-dirty {
  background: #fffbeb;
  color: #b45309;
  border: 1px solid #fde68a;
}

.status-maintenance {
  background: #fef2f2;
  color: #e11d48;
  border: 1px solid #fecdd3;
}

.quick-btn {
  padding: 4px 6px;
  font-size: 12px;
  border-radius: 6px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  cursor: pointer;
  transition: all 0.15s;
}

.quick-btn:hover {
  background: #e2e8f0;
}

.quick-btn-active {
  background: #4f46e5 !important;
  color: white;
  border-color: #4f46e5;
}

.no-data {
  text-align: center;
  padding: 40px;
  color: #94a3b8;
}

/* Button & Modal */
.btn-primary {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 10px 18px;
  font-size: 13px;
  font-weight: 700;
  color: #ffffff;
  background: #4f46e5;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(79, 70, 229, 0.2);
  transition: background 0.2s;
}

.btn-primary:hover {
  background: #4338ca;
}

.btn-secondary {
  padding: 10px 18px;
  font-size: 13px;
  font-weight: 600;
  color: #475569;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 12px;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-secondary:hover {
  background: #f8fafc;
}

.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 50;
  padding: 16px;
}

.modal-card {
  background: #ffffff;
  border-radius: 20px;
  width: 100%;
  max-width: 440px;
  overflow: hidden;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
  animation: popIn 0.2s ease-out;
}

@keyframes popIn {
  from {
    opacity: 0;
    transform: scale(0.96);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.modal-header {
  padding: 18px 24px;
  border-bottom: 1px solid #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-header h3 {
  font-size: 16px;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  font-size: 16px;
  color: #94a3b8;
  cursor: pointer;
}

.modal-body {
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 12px;
  font-weight: 700;
  color: #334155;
}

.form-input {
  padding: 10px 14px;
  font-size: 13px;
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  outline: none;
  font-family: inherit;
  transition: border-color 0.2s;
  background: white;
}

.form-input:focus {
  border-color: #6366f1;
}

.modal-footer {
  padding: 16px 24px;
  background: #f8fafc;
  border-top: 1px solid #f1f5f9;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
