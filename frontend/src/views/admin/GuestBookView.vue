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

// --- Stats Getters ---
const totalGuests = computed(() => reservations.value.length)
const checkedInGuests = computed(
  () => reservations.value.filter((r) => r.status === 'checked-in').length,
)

// --- Filter & Search Logic ---
const filteredReservations = computed(() => {
  return reservations.value.filter((res) => {
    const guestName = res.full_name || ''
    const roomName = res.room?.name || ''

    const matchesSearch =
      guestName.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
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
</script>

<template>
  <div class="guestbook-view">
    <div class="header-section">
      <div>
        <h1 class="page-title">Guest Book</h1>
        <p class="subtitle">Complete log of all guest reservations and visits.</p>
      </div>
    </div>

    <div class="stats-grid">
      <div class="stat-card">
        <span class="stat-icon">👥</span>
        <div>
          <h3>Total Visitors Logged</h3>
          <p class="main-val">{{ totalGuests }}</p>
        </div>
      </div>
      <div class="stat-card">
        <span class="stat-icon text-success">🛎️</span>
        <div>
          <h3>Currently Checked-In</h3>
          <p class="main-val">{{ checkedInGuests }}</p>
        </div>
      </div>
    </div>

    <div class="table-pane-box">
      <div class="control-bar">
        <div class="search-box">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search by Guest Name or Room..."
            class="search-input"
          />
        </div>
        <div class="filter-box">
          <select v-model="statusFilter" class="filter-select">
            <option value="All">All Visitors</option>
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
                <th>Guest Name</th>
                <th>Room</th>
                <th>Check In Date</th>
                <th>Check Out Date</th>
                <th class="text-center">Nights</th>
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
                  <span class="room-pill">{{ res.room?.name || 'N/A' }}</span>
                </td>
                <td>
                  <span class="date-text">{{ res.checkin_date }}</span>
                </td>
                <td>
                  <span class="date-text">{{ res.checkout_date }}</span>
                </td>
                <td class="text-center">
                  <span class="night-count">{{ res.total_night }}</span>
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
</template>

<style scoped>
/* Main View Blueprint */
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
  align-items: flex-end;
}
.page-title {
  font-size: 1.8rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0 0 4px 0;
  letter-spacing: -0.02em;
}
.subtitle {
  color: #64748b;
  font-size: 0.95rem;
  margin: 0;
}
.btn-export {
  background: #ffffff;
  color: #0f172a;
  border: 1px solid #e2e8f0;
  padding: 10px 18px;
  border-radius: 10px;
  font-size: 0.9rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 8px;
}
.btn-export:hover {
  background: #f1f5f9;
  border-color: #cbd5e1;
}

/* Grid Card Metrics Modern */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
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
.table-pane-box {
  display: flex;
  flex-direction: column;
  gap: 16px;
  width: 100%;
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

.responsive-table-wrap {
  overflow-x: auto;
  overflow-y: auto;
  max-height: 600px;
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
  padding: 16px 20px;
  font-size: 0.875rem;
  color: #334155;
  border-bottom: 1px solid #f1f5f9;
}
.data-row {
  transition: background 0.15s ease;
}
.data-row:hover {
  background-color: #f8fafc;
}

.bold-name {
  font-weight: 700;
  color: #0f172a;
}
.guest-info {
  display: flex;
  align-items: center;
  gap: 12px;
}
.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #e4793b;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 0.85rem;
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
.night-count {
  font-weight: 800;
  color: #334155;
  background: #f1f5f9;
  padding: 4px 10px;
  border-radius: 6px;
}
.text-center {
  text-align: center;
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
.no-data {
  text-align: center;
  padding: 40px !important;
  color: #94a3b8 !important;
}
</style>
