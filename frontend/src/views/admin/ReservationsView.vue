<script setup>
import { ref, computed, onMounted } from 'vue'
import { reservationService } from '@/services/admin/reservationService'
import { authService } from '@/services/authService'
import { useToastStore } from '@/stores/toastStore'

const reservations = ref([])
const searchQuery = ref('')
const statusFilter = ref('All')
const loading = ref(false)
const currentUser = ref(null)
const toastStore = useToastStore()

// State melacak baris yang sedang dipilih/diklik
const selectedReservation = ref(null)

// --- Stats Getters ---
const totalBookings = computed(() => reservations.value.length)
const activeStays = computed(
  () => reservations.value.filter((r) => r.status === 'checked-in').length,
)
const pendingApprovals = computed(
  () => reservations.value.filter((r) => r.status === 'pending').length,
)

// --- Filter & Search Logic ---
const filteredReservations = computed(() => {
  return reservations.value.filter((res) => {
    const guestName = res.full_name || ''
    const resId = res.id || ''
    const roomName = res.room?.name || ''

    const matchesSearch =
      guestName.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      resId.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      roomName.toLowerCase().includes(searchQuery.value.toLowerCase())

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
      if (selectedReservation.value) {
        const updated = data.find((r) => r.id === selectedReservation.value.id)
        if (updated) selectedReservation.value = updated
      }
    }
  } catch (error) {
    console.error('Refresh data error:', error)
  }
}

const selectRow = (res) => {
  selectedReservation.value = res
}

const closeDrawer = () => {
  selectedReservation.value = null
}

// --- Action API Handlers ---
const approveReservation = async (id) => {
  try {
    loading.value = true
    const msg = await reservationService.approve(id)
    toastStore.success(msg || 'Reservasi berhasil disetujui!')
    await refreshData()
  } catch (error) {
    toastStore.error(error.message || 'Gagal menyetujui reservasi')
  } finally {
    loading.value = false
  }
}

const rejectReservation = async (id) => {
  try {
    loading.value = true
    const msg = await reservationService.reject(id)
    toastStore.success(msg || 'Reservasi berhasil ditolak!')
    await refreshData()
  } catch (error) {
    toastStore.error(error.message || 'Gagal menolak reservasi')
  } finally {
    loading.value = false
  }
}

const handleCheckIn = async (id) => {
  try {
    loading.value = true
    const msg = await reservationService.checkIn(id)
    toastStore.success(msg || 'Tamu berhasil Check In!')
    await refreshData()
  } catch (error) {
    toastStore.error(error.message || 'Gagal melakukan Check In')
  } finally {
    loading.value = false
  }
}

const handleCheckOut = async (id) => {
  if (confirm('Selesaikan masa inap kamar tamu ini (Check Out)?')) {
    try {
      loading.value = true
      const msg = await reservationService.checkOut(id)
      toastStore.success(msg || 'Tamu berhasil Check Out!')
      await refreshData()
    } catch (error) {
      toastStore.error(error.message || 'Gagal melakukan Check Out')
    } finally {
      loading.value = false
    }
  }
}

const cancelReservation = async (id) => {
  if (confirm('Batalkan reservasi ini?')) {
    try {
      loading.value = true
      const msg = await reservationService.cancel(id)
      toastStore.success(msg || 'Reservasi berhasil dibatalkan!')
      await refreshData()
    } catch (error) {
      toastStore.error(error.message || 'Gagal membatalkan reservasi')
    } finally {
      loading.value = false
    }
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
</script>

<template>
  <div class="reservations-view">
    <div class="stats-grid">
      <div class="stat-card">
        <span class="stat-icon">📋</span>
        <div>
          <h3>Total Bookings</h3>
          <p class="main-val">{{ totalBookings }}</p>
        </div>
      </div>
      <div class="stat-card">
        <span class="stat-icon text-success">🛎️</span>
        <div>
          <h3>Active Stays</h3>
          <p class="main-val">{{ activeStays }}</p>
        </div>
      </div>
      <div class="stat-card">
        <span class="stat-icon text-warning">⏳</span>
        <div>
          <h3>Pending Approvals</h3>
          <p class="main-val">{{ pendingApprovals }}</p>
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
              placeholder="Search..."
              class="search-input"
            />
          </div>
          <div class="filter-box">
            <select v-model="statusFilter" class="filter-select">
              <option value="All">Semua Status</option>
              <option value="pending">Pending</option>
              <option value="approve">Approved</option>
              <option value="checked-in">Checked In</option>
              <option value="checked-out">Checked Out</option>
              <option value="cancel">Cancelled</option>
              <option value="rejected">Rejected</option>
            </select>
          </div>
        </div>

        <div class="dashboard-card table-card">
          <div class="responsive-table-wrap">
            <table class="premium-table">
              <thead>
                <tr>
                  <th>Guest Name</th>
                  <th>Room</th>
                  <th>Check In</th>
                  <th>Check Out</th>
                  <th>Total Price</th>
                  <th class="text-center">Status</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="res in filteredReservations"
                  :key="res.id"
                  @click="selectRow(res)"
                  :class="{ 'selected-row': selectedReservation?.id === res.id }"
                  class="clickable-row"
                >
                  <td class="bold-name">{{ res.full_name }}</td>
                  <td>
                    <span class="room-pill">{{ res.room?.name || 'N/A' }}</span>
                  </td>
                  <td>
                    <span class="date-text">{{ res.checkin_date }}</span>
                  </td>
                  <td>
                    <span class="date-text">{{ res.checkout_date }}</span>
                  </td>
                  <td class="price-text">
                    Rp {{ res.total_price ? res.total_price.toLocaleString('id-ID') : 0 }}
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
                  <td colspan="6" class="no-data">Tidak ditemukan data reservasi yang cocok.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <transition name="drawer-slide">
        <div v-if="selectedReservation" class="action-drawer-pane">
          <div class="drawer-header">
            <div>
              <h3>Reservation Details</h3>
              <span class="drawer-id"
                >ID: #{{ selectedReservation.id?.slice(-6).toUpperCase() }}</span
              >
            </div>
            <button @click="closeDrawer" class="close-drawer-btn">✕</button>
          </div>

          <div class="drawer-content">
            <div class="info-block-card">
              <label>Guest Name</label>
              <p class="val-large">{{ selectedReservation.full_name }}</p>
            </div>

            <div class="info-grid-2">
              <div class="info-block-card">
                <label>Room Assigned</label>
                <p class="val-mid">{{ selectedReservation.room?.name || 'N/A' }}</p>
              </div>
              <div class="info-block-card">
                <label>Current Status</label>
                <div style="margin-top: 4px">
                  <span
                    class="status-dot-badge"
                    :class="{
                      'status-pending': selectedReservation.status === 'pending',
                      'status-approved':
                        selectedReservation.status === 'approve' ||
                        selectedReservation.status === 'approved',
                      'status-checkedin': selectedReservation.status === 'checked-in',
                      'status-checkedout': selectedReservation.status === 'checked-out',
                      'status-cancel':
                        selectedReservation.status === 'cancel' ||
                        selectedReservation.status === 'rejected',
                    }"
                    >{{ selectedReservation.status }}</span
                  >
                </div>
              </div>
            </div>

            <div class="info-block-card">
              <label>OFFER CODE</label>
              <p class="val-small">
                {{ selectedReservation.offer_code || 'No Offer Code' }}
              </p>
            </div>

            <div class="info-block-card">
              <label>Stay Schedule</label>
              <p class="val-small">
                📅 {{ selectedReservation.checkin_date }} — {{ selectedReservation.checkout_date }} |
                ({{ selectedReservation.total_night }}) nights
              </p>
            </div>

            <div class="info-block-card pricing-bg">
              <label>Billing Statement</label>
              <p class="val-price">
                Rp {{ selectedReservation.total_price?.toLocaleString('id-ID') || 0 }}
              </p>
            </div>

            <div v-if="currentUser?.role === 'admin'" class="drawer-actions-container">
              <div v-if="selectedReservation.status === 'pending'" class="btn-group-row">
                <button
                  @click="approveReservation(selectedReservation.id)"
                  class="btn btn-success flex-1"
                >
                  Approve
                </button>
                <button
                  @click="rejectReservation(selectedReservation.id)"
                  class="btn btn-warning flex-1"
                >
                  Reject
                </button>
              </div>

              <button
                v-if="
                  selectedReservation.status === 'approve' ||
                  selectedReservation.status === 'approved'
                "
                @click="handleCheckIn(selectedReservation.id)"
                class="btn btn-checkedin btn-block"
              >
                Check In Guest
              </button>

              <button
                v-if="selectedReservation.status === 'checked-in'"
                @click="handleCheckOut(selectedReservation.id)"
                class="btn btn-checkedout btn-block"
              >
                Process Check Out
              </button>

              <button
                v-if="
                  selectedReservation.status === 'pending' ||
                  selectedReservation.status === 'approve' ||
                  selectedReservation.status === 'approved'
                "
                @click="cancelReservation(selectedReservation.id)"
                class="btn btn-danger-outline btn-block"
              >
                Cancel Reservation
              </button>
            </div>
          </div>
        </div>
      </transition>
    </div>
  </div>
</template>

<style scoped>
/* Main View Blueprint */
.reservations-view {
  display: flex;
  flex-direction: column;
  gap: 24px;
  background-color: #f8fafc;
  font-family: 'Plus Jakarta Sans', sans-serif;
}
.view-header h1 {
  font-size: 1.6rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0 0 4px 0;
  letter-spacing: -0.02em;
}
.subtitle {
  color: #64748b;
  font-size: 0.9rem;
  margin: 0;
}

/* Grid Card Metrics Modern */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
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
  font-size: 1.4rem;
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

/* Workspace Panes Structure */
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

/* Filter Operations Bar */
.control-bar {
  display: flex;
  gap: 14px;
  background-color: #ffffff;
  padding: 14px;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
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

/* Premium Design Table Layout */
.table-card {
  background: #ffffff;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.01);
}

/* FIXED OVERFLOW & STICKY HEADER INTERNAL SCROLL ENGINE */
.responsive-table-wrap {
  overflow-x: auto;
  overflow-y: auto;
  max-height: 520px; /* Pengunci utama anti overflow 100+ data */
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
  background: #f8fafc; /* Backing color mutlak pelindung teks shadow */
  letter-spacing: 0.05em;
}
.premium-table td {
  padding: 16px 20px;
  font-size: 0.875rem;
  color: #334155;
  border-bottom: 1px solid #f1f5f9;
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

.bold-name {
  font-weight: 700;
  color: #0f172a;
}
.room-pill {
  background: #f1f5f9;
  padding: 4px 10px;
  border-radius: 6px;
  font-weight: 700;
  font-size: 0.8rem;
  color: #475569;
}
.date-text {
  color: #475569;
  font-size: 0.85rem;
  font-weight: 500;
}
.price-text {
  font-weight: 700;
  color: #0f172a;
}

/* Status Dot Badges Framework */
.status-dot-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: capitalize;
}
.status-pending {
  background-color: #fef3c7;
  color: #92400e;
}
.status-approved {
  background-color: #d1fae5;
  color: #065f46;
}
.status-checkedin {
  background-color: #e0f2fe;
  color: #0369a1;
}
.status-checkedout {
  background-color: #f1f5f9;
  color: #475569;
}
.status-cancel {
  background-color: #fee2e2;
  color: #991b1b;
}

/* RIGHT PANEL DRAWER (Minimalist Sidebar) */
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

/* Structural Data Content Blocks */
.drawer-content {
  display: flex;
  flex-direction: column;
  gap: 14px;
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

/* Drawer Action Button System */
.drawer-actions-container {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 10px;
}
.btn-group-row {
  display: flex;
  gap: 10px;
  width: 100%;
}
.flex-1 {
  flex: 1;
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
.btn-success {
  background-color: #10b981;
  color: white;
}
.btn-success:hover {
  background-color: #059669;
}
.btn-warning {
  background-color: #f59e0b;
  color: white;
}
.btn-warning:hover {
  background-color: #d97706;
}
.btn-checkedin {
  background-color: #3b82f6;
  color: white;
}
.btn-checkedin:hover {
  background-color: #2563eb;
}
.btn-checkedout {
  background-color: #475569;
  color: white;
}
.btn-checkedout:hover {
  background-color: #334155;
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

/* Webkit Scrollbar Customization */
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

/* Vue Standard Sliding Transitions */
.drawer-slide-enter-active,
.drawer-slide-leave-active {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}
.drawer-slide-enter-from,
.drawer-slide-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
</style>
