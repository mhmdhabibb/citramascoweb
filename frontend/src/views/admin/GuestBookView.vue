<script setup lang="ts">
import { reservationService } from '@/services/admin/reservationService'
import { authService } from '@/services/authService'
import { useToastStore } from '@/stores/toastStore'
import type { Reservation } from '@/types'
import { computed, onMounted, ref } from 'vue'

const reservations = ref<Reservation[]>([])
const loading = ref(false)
const currentUser = ref<any>(null)
const toastStore = useToastStore()

// Tab Switcher
const activeTab = ref<'daily-logs' | 'guest-history'>('daily-logs')

// Filters for Daily Logs
const logSearchQuery = ref('')
const selectedLogDate = ref('')
const actionFilter = ref('All')

// Filters for Guest History
const searchQuery = ref('')
const statusFilter = ref('All')

// --- Flat Log Items ---
interface FlatLogItem {
  id: string
  reservation_id: string
  guest_name: string
  room_name: string
  room_number: string
  floor_name: string
  action: string
  timestamp: string
  is_early: boolean
  notes: string
}

const allLogs = computed<FlatLogItem[]>(() => {
  const result: FlatLogItem[] = []
  reservations.value.forEach((res) => {
    if (res.logs && Array.isArray(res.logs)) {
      res.logs.forEach((log) => {
        result.push({
          id: log.id || `${res.id}-${log.timestamp}`,
          reservation_id: res.id,
          guest_name: res.full_name,
          room_name: res.room?.name || 'N/A',
          room_number: res.room_unit?.room_number || '',
          floor_name: res.room_unit?.floor?.name || '',
          action: log.action,
          timestamp: log.timestamp,
          is_early: !!log.is_early,
          notes: log.notes || '',
        })
      })
    }
  })
  return result.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
})

// Filtered Daily Logs
const filteredLogs = computed(() => {
  return allLogs.value.filter((item) => {
    const matchesSearch =
      item.guest_name.toLowerCase().includes(logSearchQuery.value.toLowerCase()) ||
      item.room_name.toLowerCase().includes(logSearchQuery.value.toLowerCase()) ||
      item.room_number.toLowerCase().includes(logSearchQuery.value.toLowerCase())

    let matchesDate = true
    if (selectedLogDate.value) {
      const itemDateStr = new Date(item.timestamp).toISOString().split('T')[0]
      matchesDate = itemDateStr === selectedLogDate.value
    }

    const matchesAction =
      actionFilter.value === 'All' || item.action.toLowerCase() === actionFilter.value.toLowerCase()

    return matchesSearch && matchesDate && matchesAction
  })
})

// Stats for Daily Logs
const totalDailyLogs = computed(() => filteredLogs.value.length)
const totalCheckIns = computed(() => filteredLogs.value.filter((l) => l.action === 'check_in').length)
const totalCheckOuts = computed(() => filteredLogs.value.filter((l) => l.action === 'check_out').length)
const totalEarlyCheckIns = computed(() => filteredLogs.value.filter((l) => l.is_early).length)

// Stats for Reservations
const totalGuests = computed(() => reservations.value.length)
const checkedInGuests = computed(
  () => reservations.value.filter((r) => r.status === 'checked-in').length,
)

// Filtered Reservations
const filteredReservations = computed(() => {
  return reservations.value.filter((res) => {
    const guestName = res.full_name || ''
    const roomName = res.room?.name || ''
    const roomNum = res.room_unit?.room_number || ''

    const matchesSearch =
      guestName.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      roomName.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      roomNum.toLowerCase().includes(searchQuery.value.toLowerCase())

    const matchesStatus =
      statusFilter.value === 'All' || res.status?.toLowerCase() === statusFilter.value.toLowerCase()
    return matchesSearch && matchesStatus
  })
})

const refreshData = async () => {
  try {
    const data = await reservationService.getAll()
    if (data) {
      reservations.value = data
    }
  } catch (error) {
    console.error('Refresh data error:', error)
  }
}

onMounted(async () => {
  try {
    loading.value = true
    const token = localStorage.getItem('token')
    if (token === 'mock-developer-token-citramas') {
      currentUser.value = { role: 'admin', first_name: 'Developer' }
    } else {
      currentUser.value = await authService.getProfile()
    }
    await refreshData()
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
})

const setTodayFilter = () => {
  const today = new Date().toISOString().split('T')[0] || ''
  selectedLogDate.value = today
}

const clearDateFilter = () => {
  selectedLogDate.value = ''
}
</script>

<template>
  <div class="guestbook-view">
    <!-- Header Section -->
    <div class="header-section">
      <div>
        <h1 class="page-title">Buku Tamu & Log Aktivitas Harian</h1>
        <p class="subtitle">Rekap riwayat check-in & check-out harian serta catatan kedatangan tamu.</p>
      </div>

      <!-- Tab Navigation -->
      <div class="tab-pill-group">
        <button
          @click="activeTab = 'daily-logs'"
          class="tab-pill-btn"
          :class="{ 'tab-pill-active': activeTab === 'daily-logs' }"
        >
          <span>📋</span> Log In / Out Harian
        </button>
        <button
          @click="activeTab = 'guest-history'"
          class="tab-pill-btn"
          :class="{ 'tab-pill-active': activeTab === 'guest-history' }"
        >
          <span>👥</span> Buku Tamu Reservasi
        </button>
      </div>
    </div>

    <!-- TAB 1: DAILY IN/OUT LOGS -->
    <div v-if="activeTab === 'daily-logs'" class="space-y-4 animate-fade-in">
      <!-- Stats Grid -->
      <div class="stats-grid">
        <div class="stat-card">
          <span class="stat-icon bg-indigo-50 text-indigo-600">📊</span>
          <div>
            <h3>Total Aktivitas Log</h3>
            <p class="main-val">{{ totalDailyLogs }}</p>
          </div>
        </div>
        <div class="stat-card">
          <span class="stat-icon bg-emerald-50 text-emerald-600">🛎️</span>
          <div>
            <h3>Check In (Masuk)</h3>
            <p class="main-val text-emerald-600">{{ totalCheckIns }}</p>
          </div>
        </div>
        <div class="stat-card">
          <span class="stat-icon bg-slate-100 text-slate-700">🚪</span>
          <div>
            <h3>Check Out (Keluar)</h3>
            <p class="main-val text-slate-800">{{ totalCheckOuts }}</p>
          </div>
        </div>
        <div class="stat-card">
          <span class="stat-icon bg-amber-50 text-amber-600">🌅</span>
          <div>
            <h3>Early Check-In</h3>
            <p class="main-val text-amber-700">{{ totalEarlyCheckIns }}</p>
          </div>
        </div>
      </div>

      <!-- Control Bar -->
      <div class="table-pane-box">
        <div class="control-bar flex-wrap items-center justify-between gap-3">
          <div class="flex items-center gap-2 flex-wrap flex-1">
            <div class="search-box">
              <input
                v-model="logSearchQuery"
                type="text"
                placeholder="Cari nama tamu, no. kamar..."
                class="search-input"
              />
            </div>

            <!-- Date Picker Filter -->
            <div class="flex items-center gap-1.5 bg-white border border-slate-200 rounded-xl px-3 py-1.5 shadow-xs">
              <span class="text-xs text-slate-500 font-semibold">📅 Tanggal:</span>
              <input
                v-model="selectedLogDate"
                type="date"
                class="text-xs font-bold text-slate-800 outline-none bg-transparent cursor-pointer"
              />
              <button
                v-if="selectedLogDate"
                @click="clearDateFilter"
                class="text-xs text-slate-400 hover:text-rose-600 ml-1 font-bold"
                title="Hapus filter tanggal"
              >
                ✕
              </button>
            </div>

            <button
              @click="setTodayFilter"
              class="px-3 py-2 text-xs font-bold rounded-xl bg-indigo-50 hover:bg-indigo-100 text-indigo-700 border border-indigo-200 transition-colors cursor-pointer"
            >
              Hari Ini
            </button>

            <!-- Action Filter -->
            <select v-model="actionFilter" class="filter-select">
              <option value="All">Semua Aktivitas</option>
              <option value="check_in">Check In Saja</option>
              <option value="check_out">Check Out Saja</option>
            </select>
          </div>
        </div>

        <!-- Table Log Feed -->
        <div class="dashboard-card table-card">
          <div class="responsive-table-wrap">
            <table class="premium-table">
              <thead>
                <tr>
                  <th class="w-44">Waktu (WIB)</th>
                  <th>Nama Tamu</th>
                  <th>Unit & Tipe Kamar</th>
                  <th class="text-center w-36">Aktivitas</th>
                  <th>Catatan</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="log in filteredLogs" :key="log.id" class="data-row">
                  <td>
                    <div class="font-bold text-slate-900 text-xs">
                      {{ new Date(log.timestamp).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) }} WIB
                    </div>
                    <div class="text-[11px] text-slate-500 font-medium">
                      {{ new Date(log.timestamp).toLocaleDateString('id-ID', { weekday: 'short', day: '2-digit', month: 'short', year: 'numeric' }) }}
                    </div>
                  </td>
                  <td class="bold-name">
                    <div class="flex items-center gap-2">
                      <div class="w-7 h-7 rounded-full bg-slate-100 text-slate-700 font-bold text-xs flex items-center justify-center">
                        {{ log.guest_name?.charAt(0).toUpperCase() }}
                      </div>
                      <span class="text-xs font-bold text-slate-900">{{ log.guest_name }}</span>
                    </div>
                  </td>
                  <td>
                    <div class="flex items-center gap-1.5 flex-wrap">
                      <span
                        v-if="log.room_number"
                        class="px-2 py-0.5 rounded text-xs font-extrabold bg-indigo-50 text-indigo-700 border border-indigo-100"
                      >
                        🚪 No. {{ log.room_number }} {{ log.floor_name ? `(${log.floor_name})` : '' }}
                      </span>
                      <span class="text-xs font-semibold text-slate-700">
                        {{ log.room_name }}
                      </span>
                    </div>
                  </td>
                  <td class="text-center">
                    <span
                      v-if="log.action === 'check_in'"
                      class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-bold bg-emerald-50 text-emerald-700 border border-emerald-200"
                    >
                      <span>🛎️</span>
                      <span>{{ log.is_early ? 'Early Check-In' : 'Check In' }}</span>
                    </span>
                    <span
                      v-else
                      class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-bold bg-slate-100 text-slate-700 border border-slate-200"
                    >
                      <span>🚪</span>
                      <span>Check Out</span>
                    </span>
                  </td>
                  <td class="text-xs text-slate-500">
                    <span v-if="log.is_early" class="text-amber-800 font-semibold bg-amber-50 px-2 py-0.5 rounded border border-amber-200">
                      🌅 Tiba sebelum jam 14:00
                    </span>
                    <span v-else>{{ log.notes || 'Aktivitas reguler' }}</span>
                  </td>
                </tr>
                <tr v-if="filteredLogs.length === 0">
                  <td colspan="5" class="no-data">
                    <div class="py-6 text-center">
                      <p class="text-2xl mb-1">📭</p>
                      <p class="text-slate-500 font-medium text-xs">Belum ada riwayat aktivitas check-in / check-out pada filter ini.</p>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 2: GUEST RESERVATION HISTORY -->
    <div v-else class="space-y-4 animate-fade-in">
      <div class="stats-grid">
        <div class="stat-card">
          <span class="stat-icon">👥</span>
          <div>
            <h3>Total Tamu Terdaftar</h3>
            <p class="main-val">{{ totalGuests }} Tamu</p>
          </div>
        </div>
        <div class="stat-card">
          <span class="stat-icon text-success">🛎️</span>
          <div>
            <h3>Sedang Menginap</h3>
            <p class="main-val">{{ checkedInGuests }} Tamu</p>
          </div>
        </div>
      </div>

      <div class="table-pane-box">
        <div class="control-bar">
          <div class="search-box">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Cari nama tamu, no. kamar, tipe..."
              class="search-input"
            />
          </div>
          <div class="filter-box">
            <select v-model="statusFilter" class="filter-select">
              <option value="All">Semua Status</option>
              <option value="checked-in">Checked In</option>
              <option value="checked-out">Checked Out</option>
              <option value="pending">Pending</option>
              <option value="approve">Approved</option>
              <option value="cancel">Cancelled</option>
            </select>
          </div>
        </div>

        <div class="dashboard-card table-card">
          <div class="responsive-table-wrap">
            <table class="premium-table">
              <thead>
                <tr>
                  <th>Nama Tamu</th>
                  <th>Kamar & No. Fisik</th>
                  <th>Tgl Check In</th>
                  <th>Tgl Check Out</th>
                  <th class="text-center">Malam</th>
                  <th class="text-center">Status</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="res in filteredReservations"
                  :key="res.id"
                  class="data-row"
                >
                  <td class="bold-name">
                    <div class="guest-info">
                      <div class="avatar">{{ res.full_name?.charAt(0).toUpperCase() }}</div>
                      <span>{{ res.full_name }}</span>
                    </div>
                  </td>
                  <td>
                    <div>
                      <span class="room-pill">{{ res.room?.name || 'N/A' }}</span>
                      <span
                        v-if="res.room_unit?.room_number"
                        class="block text-[11px] font-extrabold text-indigo-700 bg-indigo-50 border border-indigo-200 rounded px-1.5 py-0.5 mt-1 w-fit"
                      >
                        🚪 No. {{ res.room_unit.room_number }}
                      </span>
                    </div>
                  </td>
                  <td>
                    <span class="date-text">{{ res.checkin_date }}</span>
                  </td>
                  <td>
                    <span class="date-text">{{ res.checkout_date }}</span>
                  </td>
                  <td class="text-center">
                    <span class="night-count">{{ res.total_night }} Malam</span>
                  </td>
                  <td class="text-center">
                    <span
                      class="status-dot-badge"
                      :class="{
                        'status-pending': res.status === 'pending',
                        'status-approved': res.status === 'approve' || res.status === 'approved',
                        'status-checkedin': res.status === 'checked-in',
                        'status-checkedout': res.status === 'checked-out',
                        'status-cancel': res.status === 'cancel' || res.status === 'rejected',
                      }"
                    >
                      {{ res.status }}
                    </span>
                  </td>
                </tr>
                <tr v-if="filteredReservations.length === 0">
                  <td colspan="6" class="no-data">Tidak ditemukan data buku tamu yang cocok.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.guestbook-view {
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

.tab-pill-group {
  display: inline-flex;
  background: #f1f5f9;
  padding: 4px;
  border-radius: 14px;
  gap: 4px;
  border: 1px solid #e2e8f0;
}

.tab-pill-btn {
  padding: 8px 16px;
  font-size: 13px;
  font-weight: 700;
  border-radius: 10px;
  border: none;
  background: transparent;
  color: #64748b;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: all 0.15s;
}

.tab-pill-btn:hover {
  color: #0f172a;
}

.tab-pill-active {
  background: #ffffff;
  color: #4f46e5;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
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
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
}

.stat-card h3 {
  font-size: 11px;
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

.guest-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #eef2ff;
  color: #4f46e5;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 12px;
}

.room-pill {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 6px;
  background: #f1f5f9;
  color: #475569;
  font-weight: 600;
  font-size: 12px;
}

.status-dot-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 9999px;
  font-size: 11px;
  font-weight: 700;
}

.status-pending {
  background: #fffbeb;
  color: #d97706;
}

.status-approved {
  background: #ecfdf5;
  color: #059669;
}

.status-checkedin {
  background: #eef2ff;
  color: #4f46e5;
}

.status-checkedout {
  background: #f1f5f9;
  color: #64748b;
}

.status-cancel {
  background: #fef2f2;
  color: #e11d48;
}

.no-data {
  text-align: center;
  padding: 40px;
  color: #94a3b8;
}
</style>
