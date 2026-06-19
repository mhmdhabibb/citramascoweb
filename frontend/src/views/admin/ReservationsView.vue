<script setup>
import { reservationService } from '@/services/admin/reservationService'
import { authService } from '@/services/authService'
import { ref, computed, onMounted } from 'vue'

const reservations = ref([])
const searchQuery = ref('')
const statusFilter = ref('All')
const loading = ref(false)
const currentUser = ref(null)

// Stats
const totalBookings = computed(() => reservations.value.length)
const activeStays = computed(() => reservations.value.filter(r => r.status === 'approved' || r.status === 'checked-in').length)
const pendingApprovals = computed(() => reservations.value.filter(r => r.status === 'pending').length)

// Filter Logic
const filteredReservations = computed(() => {
  return reservations.value.filter(res => {
    const guestName = res.full_name || ''
    const resId = res.id || ''
    const roomName = res.room?.name || ''
    
    const matchesSearch = guestName.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                          resId.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                          roomName.toLowerCase().includes(searchQuery.value.toLowerCase())
                          
    const matchesStatus = statusFilter.value === 'All' || 
                          res.status?.toLowerCase() === statusFilter.value.toLowerCase()
    return matchesSearch && matchesStatus
  })
})

// Actions
const approveReservation = async (id) => {
  try {
    loading.value = true
    await reservationService.approve(id)
    const data = await reservationService.getAll()
    reservations.value = data
  } catch (error) {
    alert(error.message || 'Failed to approve reservation')
  } finally {
    loading.value = false
  }
}

const rejectReservation = async (id) => {
  try {
    loading.value = true
    await reservationService.reject(id)
    const data = await reservationService.getAll()
    reservations.value = data
  } catch (error) {
    alert(error.message || 'Failed to reject reservation')
  } finally {
    loading.value = false
  }
}

const cancelReservation = async (id) => {
  if (confirm('Are you sure you want to cancel this reservation?')) {
    try {
      loading.value = true
      await reservationService.cancel(id)
      const data = await reservationService.getAll()
      reservations.value = data
    } catch (error) {
      alert(error.message || 'Failed to cancel reservation')
    } finally {
      loading.value = false
    }
  }
}

onMounted(async () => {
  try {
    loading.value = true 
    
    // Get logged in user role
    const token = localStorage.getItem('token')
    if (token === 'mock-developer-token-citramas') {
      currentUser.value = { role: 'admin', first_name: 'Developer' }
    } else {
      currentUser.value = await authService.getProfile()
    }
    
    const data = await reservationService.getAll()
    if(data) {
      reservations.value = data
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="reservations-view">
    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <h3 class="stat-title">Total Bookings</h3>
        <p class="stat-value">{{ reservations.length }}</p>
      </div>
      <div class="stat-card">
        <h3 class="stat-title">Active Stays</h3>
        <p class="stat-value text-success">{{ activeStays }}</p>
      </div>
      <div class="stat-card">
        <h3 class="stat-title">Pending Approvals</h3>
        <p class="stat-value text-warning">{{ pendingApprovals }}</p>
      </div>
    </div>

    <!-- Filter Control Bar -->
    <div class="control-bar">
      <div class="search-box">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Search by name, ID or room..." 
          class="search-input"
        />
      </div>
      <div class="filter-box">
        <select v-model="statusFilter" class="filter-select">
          <option value="All">All Statuses</option>
          <option value="approved">Approved</option>
          <option value="pending">Pending</option>
          <option value="rejected">Rejected</option>
          <option value="cancel">Cancelled</option>
          <option value="checked-in">Checked In</option>
          <option value="checked-out">Checked Out</option>
        </select>
      </div>
    </div>

    <!-- Data Table Card -->
    <div class="table-card">
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>Guest Name</th>
              <th>Room</th>
              <th>Check-in</th>
              <th>Check-out</th>
              <th>Price</th>
              <th>Status</th>
              <th v-if="currentUser?.role === 'admin'">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="res in filteredReservations" :key="res.id">
              <td>{{ res.full_name }}</td>
              <td>{{ res.room?.name || 'Room' }}</td>
              <td>{{ res.checkin_date }}</td>
              <td>{{ res.checkout_date }}</td>
              <td>Rp{{ res.total_price ? res.total_price.toLocaleString('id-ID') : 0 }}</td>
              <td>
                <span class="badge" :class="{
                  'badge-success': res.status === 'approved' || res.status === 'checked-in',
                  'badge-warning': res.status === 'pending',
                  'badge-danger': res.status === 'rejected' || res.status === 'cancel',
                  'badge-secondary': res.status === 'checked-out'
                }">{{ res.status ? res.status.charAt(0).toUpperCase() + res.status.slice(1) : '' }}</span>
              </td>
              <td v-if="currentUser?.role === 'admin'">
                <div class="action-buttons">
                  <button 
                    v-if="res.status === 'pending'" 
                    @click="approveReservation(res.id)" 
                    class="btn btn-sm btn-success"
                  >Approve</button>
                  <button 
                    v-if="res.status === 'pending'" 
                    @click="rejectReservation(res.id)" 
                    class="btn btn-sm btn-warning"
                  >Reject</button>
                  <button 
                    v-if="res.status === 'pending'" 

                    @click="cancelReservation(res.id)" 
                    class="btn btn-sm btn-danger-outline"
                  >Cancel</button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredReservations.length === 0">
              <td :colspan="currentUser?.role === 'admin' ? 7 : 6" class="no-data">No reservations found.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.reservations-view {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
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

.font-mono { font-family: monospace; }
.font-bold { font-weight: 600; }

.badge {
  padding: 4px 8px;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  display: inline-block;
}

.badge-success { background-color: #d1fae5; color: #065f46; }
.badge-warning { background-color: #fef3c7; color: #92400e; }
.badge-danger { background-color: #fee2e2; color: #991b1b; }

.action-buttons {
  display: flex;
  gap: 8px;
}

.btn {
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s;
}

.btn-sm {
  padding: 4px 8px;
}

.btn-success {
  background-color: #10b981;
  color: #ffffff;
}

.btn-warning {
  background-color: orange;
  color: white;
}

.btn-success:hover { background-color: #059669; }

.btn-warning {
  background-color: #f59e0b;
  color: #ffffff;
}
.btn-warning:hover { background-color: #d97706; }

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
</style>
