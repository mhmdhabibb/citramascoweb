<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { reservationService } from '@/services/admin/reservationService'
import { useToastStore } from '@/stores/toastStore'
import { useRouter } from 'vue-router'
import type { Reservation } from '@/types'

const router = useRouter()
const toastStore = useToastStore()

const reservations = ref<Reservation[]>([])
const loading = ref(true)

const today = new Date()
const todayStr = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(
  today.getDate(),
).padStart(2, '0')}`

function onDate(d?: string) {
  return d ? d.slice(0, 10) : ''
}

// ---- Front-desk groupings ----
const todayArrivals = computed(() =>
  reservations.value.filter(
    (r) =>
      onDate(r.checkin_date) === todayStr &&
      ['confirmed', 'approved', 'pending'].includes(r.status),
  ),
)

const inHouse = computed(() => reservations.value.filter((r) => r.status === 'checked-in'))

const dueCheckout = computed(() =>
  reservations.value.filter((r) => r.status === 'checked-in' && onDate(r.checkout_date) === todayStr),
)

const pendingList = computed(() => reservations.value.filter((r) => r.status === 'pending'))

const pendingApprovals = computed(() => pendingList.value.length)

function formatIDR(n: number) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(n || 0)
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('en-GB', {
    weekday: 'short',
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  })
}

function statusLabel(s: string) {
  const map: Record<string, string> = {
    pending: 'Pending',
    approved: 'Approved',
    confirmed: 'Confirmed',
    'checked-in': 'In House',
    'checked-out': 'Checked Out',
    cancel: 'Cancelled',
    rejected: 'Rejected',
  }
  return map[s] || s
}

function statusClass(s: string) {
  const map: Record<string, string> = {
    pending: 's-pending',
    approved: 's-approved',
    confirmed: 's-approved',
    'checked-in': 's-inhouse',
    'checked-out': 's-done',
    cancel: 's-cancel',
    rejected: 's-cancel',
  }
  return map[s] || ''
}

const goToReservations = () => router.push('/admin/reservations')

onMounted(async () => {
  try {
    loading.value = true
    reservations.value = await reservationService.getAll()
  } catch (error) {
    toastStore.error((error as Error).message || 'Failed to load front desk data')
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="reception-dashboard">
    <div class="dash-header">
      <div>
        <h1>Front Desk</h1>
        <p class="subtitle">Today's arrivals, in-house guests and departures.</p>
      </div>
      <button class="btn-primary" @click="goToReservations">Manage Reservations</button>
    </div>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon icon-blue">📥</div>
        <div>
          <h3>Today's Arrivals</h3>
          <p class="main-val">{{ todayArrivals.length }}</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon icon-green">🛎️</div>
        <div>
          <h3>In-House Guests</h3>
          <p class="main-val">{{ inHouse.length }}</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon icon-amber">🚪</div>
        <div>
          <h3>Due to Check Out</h3>
          <p class="main-val">{{ dueCheckout.length }}</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon icon-purple">⏳</div>
        <div>
          <h3>Pending Approvals</h3>
          <p class="main-val">{{ pendingApprovals }}</p>
        </div>
      </div>
    </div>

    <div class="split-grid">
      <div class="table-card">
        <div class="card-heading">
          <h2>Today's Arrivals</h2>
          <p class="card-sub">{{ formatDate(todayStr) }}</p>
        </div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Guest</th>
                <th>Room</th>
                <th>Guests</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="res in todayArrivals" :key="res.id">
                <td class="bold-name">{{ res.full_name }}</td>
                <td><span class="room-pill">{{ res.room?.name || 'N/A' }}</span></td>
                <td>{{ res.number_of_guest || 0 }}</td>
                <td><span class="badge" :class="statusClass(res.status)">{{ statusLabel(res.status) }}</span></td>
              </tr>
              <tr v-if="!loading && todayArrivals.length === 0">
                <td colspan="4" class="no-data">No arrivals due today.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="table-card">
        <div class="card-heading">
          <h2>In-House Guests</h2>
          <p class="card-sub">Currently checked-in guests.</p>
        </div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Guest</th>
                <th>Room</th>
                <th>Check Out</th>
                <th>Total</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="res in inHouse" :key="res.id">
                <td class="bold-name">{{ res.full_name }}</td>
                <td><span class="room-pill">{{ res.room?.name || 'N/A' }}</span></td>
                <td>{{ onDate(res.checkout_date) }}</td>
                <td class="price-text">{{ formatIDR(res.total_price) }}</td>
              </tr>
              <tr v-if="!loading && inHouse.length === 0">
                <td colspan="4" class="no-data">No guests currently in-house.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <div class="table-card">
      <div class="card-heading">
        <h2>Pending Approvals</h2>
        <p class="card-sub">Reservations waiting for confirmation.</p>
      </div>
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Guest</th>
              <th>Room</th>
              <th>Check In</th>
              <th>Check Out</th>
              <th>Guests</th>
              <th>Total</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="res in pendingList" :key="res.id">
              <td class="bold-name">{{ res.full_name }}</td>
              <td><span class="room-pill">{{ res.room?.name || 'N/A' }}</span></td>
              <td>{{ onDate(res.checkin_date) }}</td>
              <td>{{ onDate(res.checkout_date) }}</td>
              <td>{{ res.number_of_guest || 0 }}</td>
              <td class="price-text">{{ formatIDR(res.total_price) }}</td>
            </tr>
            <tr v-if="!loading && pendingList.length === 0">
              <td colspan="6" class="no-data">No pending approvals.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.reception-dashboard {
  background-color: #f8fafc;
  padding: 32px;
  display: flex;
  flex-direction: column;
  gap: 24px;
  font-family: 'Plus Jakarta Sans', sans-serif;
}
.dash-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.dash-header h1 {
  font-size: 1.75rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0 0 4px 0;
}
.subtitle {
  color: #64748b;
  margin: 0;
  font-size: 0.95rem;
}
.btn-primary {
  background: #e4793b;
  color: white;
  border: none;
  padding: 11px 20px;
  border-radius: 10px;
  font-weight: 600;
  font-size: 0.9rem;
  cursor: pointer;
}
.btn-primary:hover {
  background: #d16627;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}
.stat-card {
  background: #ffffff;
  border-radius: 16px;
  padding: 18px 20px;
  display: flex;
  align-items: center;
  gap: 14px;
  border: 1px solid #edf2f7;
}
.stat-icon {
  width: 46px;
  height: 46px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.3rem;
}
.icon-blue { background-color: #e0f2fe; }
.icon-green { background-color: #d1fae5; }
.icon-amber { background-color: #fef3c7; }
.icon-purple { background-color: #ede9fe; }
.stat-card h3 {
  font-size: 0.78rem;
  color: #64748b;
  margin: 0 0 2px 0;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}
.main-val {
  font-size: 1.6rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.split-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  align-items: start;
}
.table-card {
  background: #ffffff;
  border-radius: 16px;
  border: 1px solid #edf2f7;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.02);
}
.card-heading {
  padding: 18px 20px 6px;
}
.card-heading h2 {
  margin: 0;
  font-size: 1.05rem;
  font-weight: 700;
  color: #1e293b;
}
.card-sub {
  margin: 4px 0 0;
  color: #94a3b8;
  font-size: 0.8rem;
}
.table-wrap {
  overflow-x: auto;
  max-height: 320px;
  overflow-y: auto;
  padding-bottom: 12px;
}
.table-wrap table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}
.table-wrap th {
  position: sticky;
  top: 0;
  background: #f8fafc;
  padding: 12px 20px;
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  border-bottom: 1px solid #edf2f7;
  letter-spacing: 0.05em;
}
.table-wrap td {
  padding: 14px 20px;
  font-size: 0.875rem;
  color: #334155;
  border-bottom: 1px solid #f1f5f9;
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
.price-text {
  font-weight: 700;
  color: #0f172a;
}
.badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: capitalize;
}
.s-pending { background-color: #fef3c7; color: #92400e; }
.s-approved { background-color: #dbeafe; color: #1d4ed8; }
.s-inhouse { background-color: #d1fae5; color: #065f46; }
.s-done { background-color: #f1f5f9; color: #475569; }
.s-cancel { background-color: #fee2e2; color: #991b1b; }
.no-data {
  text-align: center;
  color: #64748b;
  padding: 40px !important;
}

@media (max-width: 1100px) {
  .stats-grid { grid-template-columns: repeat(2, 1fr); }
  .split-grid { grid-template-columns: 1fr; }
}
</style>