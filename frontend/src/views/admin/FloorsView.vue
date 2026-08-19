<script setup lang="ts">
import { floorService } from '@/services/admin/floorService'
import { useAuthStore } from '@/stores/authStore'
import { canMutateData } from '@/config/access'
import { useToastStore } from '@/stores/toastStore'
import type { Floor } from '@/types'
import { computed, onMounted, ref } from 'vue'

const authStore = useAuthStore()

const floors = ref<Floor[]>([])
const loading = ref(false)
const toastStore = useToastStore()
const searchQuery = ref('')

// Modal state
const isModalOpen = ref(false)
const isEditing = ref(false)
const editingId = ref('')
const form = ref({
  name: '',
  floor_number: 1,
  description: '',
})

// Stats
const totalFloors = computed(() => floors.value.length)
const maxFloor = computed(() => {
  if (floors.value.length === 0) return 0
  return Math.max(...floors.value.map((f) => f.floor_number))
})

const fetchFloors = async () => {
  try {
    loading.value = true
    floors.value = await floorService.getAll()
  } catch (error: any) {
    console.error('Error fetching floors:', error)
    toastStore.error(error.message || 'Gagal memuat data lantai')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchFloors()
})

const filteredFloors = computed(() => {
  if (!searchQuery.value) return floors.value
  const query = searchQuery.value.toLowerCase()
  return floors.value.filter(
    (f) =>
      f.name.toLowerCase().includes(query) ||
      (f.description && f.description.toLowerCase().includes(query)) ||
      String(f.floor_number).includes(query),
  )
})

const openCreateModal = () => {
  isEditing.value = false
  editingId.value = ''
  form.value = {
    name: `Lantai ${maxFloor.value + 1}`,
    floor_number: maxFloor.value + 1,
    description: '',
  }
  isModalOpen.value = true
}

const openEditModal = (item: Floor) => {
  isEditing.value = true
  editingId.value = item.id
  form.value = {
    name: item.name,
    floor_number: item.floor_number,
    description: item.description || '',
  }
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
}

const saveFloor = async () => {
  if (!form.value.name.trim()) {
    toastStore.warning('Nama lantai tidak boleh kosong!')
    return
  }

  try {
    loading.value = true
    if (isEditing.value) {
      const msg = await floorService.update(editingId.value, form.value)
      toastStore.success(msg || 'Lantai berhasil diperbarui!')
    } else {
      const msg = await floorService.create(form.value)
      toastStore.success(msg || 'Lantai berhasil ditambahkan!')
    }
    await fetchFloors()
    closeModal()
  } catch (error: any) {
    console.error('Error saving floor:', error)
    toastStore.error(error.response?.data?.message || error.message || 'Gagal menyimpan lantai')
  } finally {
    loading.value = false
  }
}

const deleteFloor = async (id: string, name: string) => {
  if (!confirm(`Apakah Anda yakin ingin menghapus "${name}"?`)) return

  try {
    loading.value = true
    const msg = await floorService.delete(id)
    toastStore.success(msg || 'Lantai berhasil dihapus!')
    await fetchFloors()
  } catch (error: any) {
    console.error('Error deleting floor:', error)
    toastStore.error(error.response?.data?.message || error.message || 'Gagal menghapus lantai')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="floors-view">
    <!-- Header Section -->
    <div class="header-section">
      <div>
        <h1 class="page-title">Master Data Lantai (Floors)</h1>
        <p class="subtitle">Kelola daftar lantai gedung hotel untuk penempatan unit kamar.</p>
      </div>
      <button v-if="canMutateData(authStore.role)" @click="openCreateModal" class="btn-primary">
        <span>+</span> Tambah Lantai
      </button>
      <div v-else class="px-3.5 py-2 rounded-xl bg-slate-900 text-amber-300 text-xs font-bold border border-indigo-900/50 flex items-center gap-2 shadow-sm">
        <span>👑</span> Mode Tinjauan Manager (Read-Only)
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">🏢</div>
        <div>
          <h3>Total Lantai Terdaftar</h3>
          <p class="main-val">{{ totalFloors }} Lantai</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon text-indigo-600">🔝</div>
        <div>
          <h3>Lantai Tertinggi</h3>
          <p class="main-val">Lantai {{ maxFloor }}</p>
        </div>
      </div>
    </div>

    <!-- Table Pane -->
    <div class="table-pane-box">
      <div class="control-bar">
        <div class="search-box">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari nama lantai atau nomor..."
            class="search-input"
          />
        </div>
      </div>

      <div class="dashboard-card table-card">
        <div class="responsive-table-wrap">
          <table class="premium-table">
            <thead>
              <tr>
                <th class="w-24 text-center">No. Lantai</th>
                <th>Nama Lantai</th>
                <th>Deskripsi</th>
                <th v-if="canMutateData(authStore.role)" class="text-center w-32">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="fl in filteredFloors" :key="fl.id" class="data-row">
                <td class="text-center">
                  <span class="px-2.5 py-1 rounded-lg text-xs font-black bg-indigo-50 text-indigo-700 border border-indigo-100">
                    L{{ fl.floor_number }}
                  </span>
                </td>
                <td class="bold-name">{{ fl.name }}</td>
                <td class="text-slate-500 text-sm">{{ fl.description || '-' }}</td>
                <td v-if="canMutateData(authStore.role)" class="text-center">
                  <div class="flex items-center justify-center gap-1.5">
                    <button
                      @click="openEditModal(fl)"
                      class="px-2.5 py-1 text-xs font-semibold text-indigo-600 bg-indigo-50 hover:bg-indigo-100 rounded-lg transition-colors cursor-pointer"
                    >
                      Edit
                    </button>
                    <button
                      @click="deleteFloor(fl.id, fl.name)"
                      class="px-2.5 py-1 text-xs font-semibold text-rose-600 bg-rose-50 hover:bg-rose-100 rounded-lg transition-colors cursor-pointer"
                    >
                      Hapus
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="filteredFloors.length === 0">
                <td :colspan="canMutateData(authStore.role) ? 4 : 3" class="no-data">Tidak ada data lantai yang cocok.</td>
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
          <h3>{{ isEditing ? 'Edit Lantai' : 'Tambah Lantai Baru' }}</h3>
          <button @click="closeModal" class="close-btn">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>Nomor Urut Lantai *</label>
            <input
              v-model.number="form.floor_number"
              type="number"
              placeholder="Contoh: 1, 2, 3"
              class="form-input"
              required
            />
          </div>
          <div class="form-group">
            <label>Nama Lantai *</label>
            <input
              v-model="form.name"
              type="text"
              placeholder="Contoh: Lantai 1, Ground Floor"
              class="form-input"
              required
            />
          </div>
          <div class="form-group">
            <label>Deskripsi / Keterangan</label>
            <textarea
              v-model="form.description"
              rows="3"
              placeholder="Keterangan tambahan (opsional)..."
              class="form-textarea"
            ></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeModal" class="btn-secondary">Batal</button>
          <button @click="saveFloor" :disabled="loading" class="btn-primary">
            {{ loading ? 'Menyimpan...' : 'Simpan' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.floors-view {
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
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.stat-card {
  background: #ffffff;
  padding: 16px 20px;
  border-radius: 16px;
  border: 1px solid #f1f5f9;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.02);
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: #f8fafc;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
}

.stat-card h3 {
  font-size: 12px;
  font-weight: 600;
  color: #64748b;
  margin: 0;
}

.stat-card .main-val {
  font-size: 18px;
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
  max-width: 360px;
}

.search-input {
  width: 100%;
  padding: 10px 14px;
  font-size: 13px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  outline: none;
  transition: border-color 0.2s;
}

.search-input:focus {
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
  padding: 14px 18px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #64748b;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
}

.premium-table td {
  padding: 14px 18px;
  font-size: 13px;
  border-bottom: 1px solid #f1f5f9;
}

.data-row:hover {
  background-color: #f8fafc;
}

.bold-name {
  font-weight: 700;
  color: #0f172a;
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

.form-input,
.form-textarea {
  padding: 10px 14px;
  font-size: 13px;
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  outline: none;
  font-family: inherit;
  transition: border-color 0.2s;
}

.form-input:focus,
.form-textarea:focus {
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
