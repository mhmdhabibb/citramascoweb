<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { reservationService } from '@/services/admin/reservationService'
import { useToastStore } from '@/stores/toastStore'
import type { Reservation } from '@/types'

const toastStore = useToastStore()

const reservations = ref<Reservation[]>([])
const loading = ref(false)

const nonCancelled = computed(() =>
  reservations.value.filter((r) => r.status !== 'cancel'),
)

const totalRevenue = computed(() =>
  nonCancelled.value.reduce((sum, r) => sum + (r.total_price || 0), 0),
)

const totalBookings = computed(() => reservations.value.length)

const confirmed = computed(
  () =>
    reservations.value.filter((r) =>
      ['confirmed', 'approved', 'checked-in', 'checked-out'].includes(r.status),
    ).length,
)

const checkedIn = computed(
  () => reservations.value.filter((r) => r.status === 'checked-in').length,
)

const rows = computed(() =>
  nonCancelled.value.map((r) => ({
    id: r.id,
    code: r.code,
    guest: r.full_name,
    room: r.room?.name || '',
    checkin: formatDate(r.checkin_date),
    checkout: formatDate(r.checkout_date),
    status: r.status,
    total: r.total_price || 0,
  })),
)

function formatIDR(n: number) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(n)
}

function formatDate(d?: string) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('en-GB', {
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
    'checked-in': 'Checked In',
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
    'checked-in': 's-active',
    'checked-out': 's-done',
    cancel: 's-cancel',
    rejected: 's-cancel',
  }
  return map[s] || ''
}

const fetchAll = async () => {
  try {
    loading.value = true
    reservations.value = await reservationService.getAll()
  } catch (error) {
    toastStore.error((error as Error).message || 'Failed to load finance data')
  } finally {
    loading.value = false
  }
}

onMounted(fetchAll)
</script>

<template>
  <div class="finance-view">
    <!-- Revenue summary -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon stat-icon-green">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="1" x2="12" y2="23" /><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" /></svg>
        </div>
        <div class="stat-info"><span class="stat-label">Total Revenue</span><span class="stat-value">{{ formatIDR(totalRevenue) }}</span></div>
      </div>
      <div class="stat-card">
        <div class="stat-icon stat-icon-primary">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" /><circle cx="9" cy="7" r="4" /><path d="M22 21v-2a4 4 0 0 0-3-3.87" /><path d="M16 3.13a4 4 0 0 1 0 7.75" /></svg>
        </div>
        <div class="stat-info"><span class="stat-label">Total Bookings</span><span class="stat-value">{{ totalBookings }}</span></div>
      </div>
      <div class="stat-card">
        <div class="stat-icon stat-icon-purple">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" /><line x1="16" y1="2" x2="16" y2="6" /><line x1="8" y1="2" x2="8" y2="6" /><line x1="3" y1="10" x2="21" y2="10" /></svg>
        </div>
        <div class="stat-info"><span class="stat-label">Confirmed</span><span class="stat-value">{{ confirmed }}</span></div>
      </div>
      <div class="stat-card">
        <div class="stat-icon stat-icon-amber">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3" /></svg>
        </div>
        <div class="stat-info"><span class="stat-label">Checked In</span><span class="stat-value">{{ checkedIn }}</span></div>
      </div>
    </div>

    <!-- Bookings / revenue table -->
    <div class="table-card">
      <div class="card-heading">
        <h3 class="card-title">Booking Revenue</h3>
        <p class="card-sub">Customer bookings and their total price (excluding cancelled).</p>
      </div>
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>Guest</th>
              <th>Code</th>
              <th>Room</th>
              <th>Check-in</th>
              <th>Check-out</th>
              <th>Status</th>
              <th>Total</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in rows" :key="row.id">
              <td>{{ row.guest }}</td>
              <td class="font-mono">{{ row.code }}</td>
              <td>{{ row.room }}</td>
              <td>{{ row.checkin }}</td>
              <td>{{ row.checkout }}</td>
              <td><span class="badge" :class="statusClass(row.status)">{{ statusLabel(row.status) }}</span></td>
              <td class="bold">{{ formatIDR(row.total) }}</td>
            </tr>
            <tr v-if="loading"><td colspan="7" class="no-data">Loading...</td></tr>
            <tr v-else-if="rows.length === 0"><td colspan="7" class="no-data">No bookings yet.</td></tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.finance-view {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 14px;
  background-color: #ffffff;
  padding: 18px;
  border-radius: 16px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.stat-icon svg { width: 22px; height: 22px; }
.stat-icon-primary { background-color: #fff1ea; color: #e15b2b; }
.stat-icon-green { background-color: #d1fae5; color: #059669; }
.stat-icon-purple { background-color: #ede9fe; color: #7c3aed; }
.stat-icon-amber { background-color: #fef3c7; color: #b45309; }

.stat-info { display: flex; flex-direction: column; }
.stat-label {
  font-size: 0.75rem; font-weight: 600; color: #64748b;
  text-transform: uppercase; letter-spacing: 0.04em;
}
.stat-value { font-size: 1.5rem; font-weight: 700; color: #0f172a; }

.table-card {
  background-color: #ffffff;
  border-radius: 16px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}
.card-heading { padding: 18px 20px 0; }
.card-title { margin: 0; font-size: 1rem; color: #0f172a; }
.card-sub { margin: 4px 0 0; color: #64748b; font-size: 0.82rem; }

.table-container { overflow-x: auto; }
.data-table { width: 100%; border-collapse: collapse; text-align: left; }
.data-table th {
  padding: 14px 20px; border-bottom: 1px solid #e2e8f0;
  color: #64748b; font-size: 0.75rem; font-weight: 700;
  text-transform: uppercase; letter-spacing: 0.05em;
}
.data-table td {
  padding: 16px 20px; border-bottom: 1px solid #f1f5f9;
  font-size: 0.875rem; color: #334155; vertical-align: middle;
}
.font-mono { font-family: monospace; }
.bold { font-weight: 700; }

.badge {
  padding: 4px 10px; border-radius: 9999px;
  font-size: 0.75rem; font-weight: 600; display: inline-block;
}
.s-pending { background-color: #fef3c7; color: #92400e; }
.s-approved { background-color: #dbeafe; color: #1d4ed8; }
.s-active { background-color: #d1fae5; color: #065f46; }
.s-done { background-color: #e2e8f0; color: #334155; }
.s-cancel { background-color: #fee2e2; color: #991b1b; }

.no-data { text-align: center; color: #64748b; padding: 32px !important; }
</style>