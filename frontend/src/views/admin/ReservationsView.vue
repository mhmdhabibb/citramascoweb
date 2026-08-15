<script setup>
import { financeService } from '@/services/admin/financeService'
import { reservationService } from '@/services/admin/reservationService'
import { authService } from '@/services/authService'
import { useToastStore } from '@/stores/toastStore'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const reservations = ref([])
const searchQuery = ref('')
const statusFilter = ref('All')
const loading = ref(false)
const currentUser = ref(null)
const toastStore = useToastStore()

// Receptionists (and above) may operate the front-desk reservation actions
const canManage = computed(() =>
  ['admin', 'manager', 'reception'].includes(currentUser.value?.role),
)

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

let pollingTimer = null

const refreshData = async (isSilent = false) => {
  try {
    if (!isSilent) loading.value = true
    const data = await reservationService.getAll()
    if (data) {
      reservations.value = data
      if (selectedReservation.value) {
        const updated = data.find((r) => r.id === selectedReservation.value.id)
        if (updated) selectedReservation.value = updated
      }
    }
  } catch (error) {
    if (!isSilent) console.error('Refresh data error:', error)
  } finally {
    if (!isSilent) loading.value = false
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

const verifyFinancePayment = async (id, status = 'confirmed') => {
  try {
    loading.value = true
    await financeService.verifyPayment(id, status)
    toastStore.success(status === 'confirmed' ? 'Pembayaran berhasil diverifikasi oleh Tim Finance (Status: Dibayar)!' : 'Pembayaran ditolak')
    await refreshData()
  } catch (error) {
    toastStore.error(error.message || 'Gagal memverifikasi pembayaran')
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
    await refreshData(false)
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }

  // Silent auto-reload polling every 5 seconds
  pollingTimer = setInterval(() => {
    refreshData(true)
  }, 5000)
})

onUnmounted(() => {
  if (pollingTimer) {
    clearInterval(pollingTimer)
    pollingTimer = null
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
          <button v-if="canManage" class="btn-new" @click="router.push('/admin/reservations/new')">
            + New Reservation
          </button>
        </div>

        <div class="dashboard-card table-card">
          <div class="responsive-table-wrap">
            <table class="premium-table">
              <thead>
                <tr>
                  <th>Guest Name</th>
                  <th>Channel</th>
                  <th>Room</th>
                  <th>Check In</th>
                  <th>Check Out</th>
                  <th>Total Price</th>
                  <th class="text-center">Pembayaran</th>
                  <th class="text-center">Status</th>
                  <th class="text-center">Aksi Resepsionis</th>
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
                    <span class="channel-badge-pill">
                      {{ res.channel?.name || 'Direct' }}
                    </span>
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
                  <td class="price-text">
                    Rp {{ res.total_price ? res.total_price.toLocaleString('id-ID') : 0 }}
                  </td>
                  <td class="text-center">
                    <span
                      v-if="res.transaction_status === 'paid' || (res.deposit >= res.total_price && res.total_price > 0)"
                      class="px-2.5 py-1 text-xs font-bold rounded-full bg-emerald-50 text-emerald-700 border border-emerald-300 inline-flex items-center gap-1"
                    >
                      <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
                      Dibayar
                    </span>
                    <span
                      v-else-if="res.transaction_status === 'down_payment' || (res.deposit > 0)"
                      class="px-2.5 py-1 text-xs font-bold rounded-full bg-blue-50 text-blue-700 border border-blue-300 inline-flex items-center gap-1"
                    >
                      <span class="w-1.5 h-1.5 rounded-full bg-blue-500"></span>
                      DP Rp {{ (res.deposit || 0).toLocaleString('id-ID') }}
                    </span>
                    <span
                      v-else
                      class="px-2.5 py-1 text-xs font-bold rounded-full bg-amber-50 text-amber-700 border border-amber-300 inline-flex items-center gap-1"
                    >
                      <span class="w-1.5 h-1.5 rounded-full bg-amber-500"></span>
                      Belum Bayar
                    </span>
                  </td>
                  <td class="text-center">
                    <span
                      class="status-dot-badge"
                      :class="{
                        'status-pending': res.status === 'pending',
                        'status-approved': res.status === 'approve' || res.status === 'approved' || res.status === 'confirmed',
                        'status-checkedin': res.status === 'checked-in',
                        'status-checkedout': res.status === 'checked-out',
                        'status-cancel': res.status === 'cancel' || res.status === 'rejected',
                      }"
                    >
                      {{ res.status === 'pending' ? 'Pending' : (res.status === 'approved' || res.status === 'confirmed' ? 'Confirmed' : res.status) }}
                    </span>
                  </td>
                  <td class="text-center" @click.stop>
                    <!-- KONDISI 1: PENDING & SUDAH DIBAYAR / CASH -> MUNCUL TOMBOL APPROVE RESEPSIONIS -->
                    <div
                      v-if="
                        res.status === 'pending' &&
                        (res.transaction_status === 'paid' ||
                          res.transaction_status === 'down_payment' ||
                          res.payment_method === 'cash' ||
                          (res.deposit && res.deposit > 0))
                      "
                      class="flex items-center justify-center gap-1"
                    >
                      <button
                        @click="approveReservation(res.id)"
                        class="px-2.5 py-1 bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs rounded-lg shadow-sm transition-all flex items-center gap-1"
                        title="Klik untuk Approve Reservasi & Kunci Kamar"
                      >
                        <span>✓</span>
                        <span>Approve</span>
                      </button>
                      <button
                        @click="rejectReservation(res.id)"
                        class="px-2 py-1 bg-rose-50 hover:bg-rose-100 text-rose-600 font-semibold text-xs rounded-lg border border-rose-200"
                        title="Tolak Reservasi"
                      >
                        ✕
                      </button>
                    </div>

                    <!-- KONDISI 2: PENDING & BELUM BAYAR (TRANSFER) -> MUNCUL TOMBOL FINANCE CONFIRM -->
                    <div
                      v-else-if="res.status === 'pending'"
                      class="flex items-center justify-center"
                    >
                      <button
                        @click="verifyFinancePayment(res.id, 'confirmed')"
                        class="px-2 py-1 bg-amber-50 hover:bg-amber-100 text-amber-800 font-bold text-xs rounded-lg border border-amber-300 transition-all flex items-center gap-1"
                        title="Klik untuk konfirmasi bahwa pembayaran sudah masuk rekening"
                      >
                        <span>💳</span>
                        <span>Finance Confirm</span>
                      </button>
                    </div>

                    <!-- KONDISI 3: CONFIRMED / APPROVED -> MUNCUL TOMBOL CHECK IN -->
                    <div
                      v-else-if="res.status === 'approved' || res.status === 'confirmed' || res.status === 'approve'"
                      class="flex items-center justify-center"
                    >
                      <button
                        @click="handleCheckIn(res.id)"
                        class="px-2.5 py-1 bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-xs rounded-lg shadow-sm transition-all flex items-center gap-1"
                        title="Proses Check In Tamu"
                      >
                        <span>🛎️</span>
                        <span>Check In</span>
                      </button>
                    </div>

                    <!-- KONDISI 4: CHECKED-IN -> MUNCUL TOMBOL CHECK OUT -->
                    <div
                      v-else-if="res.status === 'checked-in'"
                      class="flex items-center justify-center"
                    >
                      <button
                        @click="handleCheckOut(res.id)"
                        class="px-2.5 py-1 bg-slate-700 hover:bg-slate-800 text-white font-bold text-xs rounded-lg shadow-sm transition-all flex items-center gap-1"
                        title="Proses Check Out Tamu"
                      >
                        <span>🚪</span>
                        <span>Check Out</span>
                      </button>
                    </div>

                    <!-- KONDISI LAIN: SELESAI / CANCEL -->
                    <span v-else class="text-xs text-slate-400 font-medium">Selesai</span>
                  </td>
                </tr>
                <tr v-if="filteredReservations.length === 0">
                  <td colspan="9" class="no-data">Tidak ditemukan data reservasi yang cocok.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Detail Drawer Modal -->
      <transition name="drawer">
  <div v-if="selectedReservation" class="details-drawer w-full max-w-md bg-slate-50/50 border-l border-slate-200 shadow-2xl flex flex-col h-full overflow-hidden">
    <!-- Header with Gradient Accent -->
    <div class="drawer-header bg-white p-6 border-b border-slate-100 flex items-start justify-between shadow-xs">
      <div>
        <div class="flex items-center gap-2 mb-1">
          <span class="px-2.5 py-0.5 rounded-full text-xs font-semibold tracking-wider bg-blue-50 text-blue-700 border border-blue-100">
            #{{ selectedReservation.code }}
          </span>
          <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-medium bg-slate-100 text-slate-600">
            <span class="w-1.5 h-1.5 rounded-full bg-slate-400"></span>
            {{ selectedReservation.channel?.name || 'Direct Booking' }}
          </span>
        </div>
        <h2 class="text-xl font-extrabold text-slate-900 tracking-tight">{{ selectedReservation.full_name }}</h2>
      </div>
      <button 
        @click="closeDrawer" 
        class="w-8 h-8 flex items-center justify-center rounded-full bg-slate-100 hover:bg-slate-200 text-slate-500 hover:text-slate-800 transition-colors duration-150 cursor-pointer"
      >
        ✕
      </button>
    </div>

    <div class="drawer-body p-6 space-y-4 overflow-y-auto flex-1">
      <!-- Contact & Room Card -->
      <div class="bg-white rounded-2xl p-4 border border-slate-100 shadow-xs hover:border-slate-200 transition-all">
        <span class="text-[11px] font-bold tracking-wider text-slate-400 uppercase">Kontak & Kamar</span>
        <div class="flex items-center gap-2 mt-2">
          <div class="w-8 h-8 rounded-lg bg-indigo-50 flex items-center justify-center text-indigo-600 text-sm">
            ✉️
          </div>
          <span class="font-medium text-slate-700 text-sm truncate">{{ selectedReservation.email || '-' }}</span>
        </div>
        <div class="mt-3 flex items-center gap-2 pt-3 border-t border-slate-50">
          <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg text-xs font-semibold bg-indigo-50 text-indigo-700 border border-indigo-100">
            🛏️ {{ selectedReservation.room?.name || 'Standard Room' }}
          </span>
        </div>
      </div>

      <!-- Stay Schedule & Guests Card -->
      <div class="bg-white rounded-2xl p-4 border border-slate-100 shadow-xs hover:border-slate-200 transition-all">
        <span class="text-[11px] font-bold tracking-wider text-slate-400 uppercase">Jadwal & Tamu</span>
        
        <!-- Date Timeline -->
        <div class="mt-2 flex items-center justify-between bg-slate-50 rounded-xl p-3 border border-slate-100">
          <div>
            <div class="text-[10px] text-slate-400 font-semibold uppercase">Check-in</div>
            <div class="text-xs font-bold text-slate-800">{{ selectedReservation.checkin_date }}</div>
          </div>
          <div class="flex flex-col items-center px-2">
            <span class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-white text-indigo-600 shadow-xs border border-slate-100">
              {{ selectedReservation.total_night }} Malam
            </span>
            <div class="w-12 h-0.5 bg-slate-200 my-1 relative">
              <span class="absolute -right-0.5 -top-0.5 w-1.5 h-1.5 rounded-full bg-slate-400"></span>
            </div>
          </div>
          <div class="text-right">
            <div class="text-[10px] text-slate-400 font-semibold uppercase">Check-out</div>
            <div class="text-xs font-bold text-slate-800">{{ selectedReservation.checkout_date }}</div>
          </div>
        </div>

        <!-- Guest Badges -->
        <div class="mt-3 grid grid-cols-2 gap-2">
          <div class="flex items-center gap-2.5 p-2 rounded-xl bg-amber-50/60 border border-amber-100/80">
            <div class="w-7 h-7 rounded-lg bg-amber-100 flex items-center justify-center text-xs">👤</div>
            <div>
              <div class="text-[10px] text-amber-700/70 font-semibold">Dewasa</div>
              <div class="text-xs font-extrabold text-amber-900">
                {{ selectedReservation.adults || selectedReservation.number_of_adult || selectedReservation.adult || 1 }} Orang
              </div>
            </div>
          </div>
          <div class="flex items-center gap-2.5 p-2 rounded-xl bg-sky-50/60 border border-sky-100/80">
            <div class="w-7 h-7 rounded-lg bg-sky-100 flex items-center justify-center text-xs">🧒</div>
            <div>
              <div class="text-[10px] text-sky-700/70 font-semibold">Anak-anak</div>
              <div class="text-xs font-extrabold text-sky-900">
                {{ selectedReservation.children || selectedReservation.number_of_children || selectedReservation.child || 0 }} Anak
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Payment Method & Status Card -->
      <div class="bg-white rounded-2xl p-4 border border-slate-100 shadow-xs flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-9 h-9 rounded-xl bg-slate-100 flex items-center justify-center text-slate-700">
            💳
          </div>
          <div>
            <div class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Metode Bayar</div>
            <div class="text-xs font-bold text-slate-800">
              {{ selectedReservation.payment_method ? selectedReservation.payment_method.toUpperCase().replace('_', ' ') : 'BANK TRANSFER' }}
            </div>
          </div>
        </div>

        <div>
          <span
            v-if="selectedReservation.transaction_status === 'paid' || (selectedReservation.deposit >= selectedReservation.total_price && selectedReservation.total_price > 0)"
            class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold bg-emerald-50 text-emerald-600 border border-emerald-200"
          >
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
            Dibayar
          </span>
          <span
            v-else
            class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold bg-amber-50 text-amber-700 border border-amber-200"
          >
            <span class="w-1.5 h-1.5 rounded-full bg-amber-500"></span>
            Belum Bayar
          </span>
        </div>
      </div>

      <!-- Billing Statement Highlight Card -->
      <div class="relative overflow-hidden rounded-2xl bg-gradient-to-br from-amber-500 via-orange-500 to-rose-500 p-5 text-white shadow-lg shadow-orange-500/20">
        <div class="absolute -right-4 -bottom-4 w-24 h-24 rounded-full bg-white/10 blur-xl pointer-events-none"></div>
        <div class="text-xs font-medium text-orange-100 tracking-wider uppercase">Total Pembayaran</div>
        <div class="mt-1 flex items-baseline gap-1">
          <span class="text-sm font-semibold opacity-90">Rp</span>
          <span class="text-2xl font-black tracking-tight">
            {{ selectedReservation.total_price?.toLocaleString('id-ID') || 0 }}
          </span>
        </div>
      </div>

      <!-- Actions Container -->
      <div v-if="canManage" class="pt-2">
        <div v-if="selectedReservation.status === 'pending'" class="space-y-3">
          <!-- Paid Confirmation -->
          <div
            v-if="
              selectedReservation.transaction_status === 'paid' ||
              selectedReservation.transaction_status === 'down_payment' ||
              selectedReservation.payment_method === 'cash' ||
              (selectedReservation.deposit && selectedReservation.deposit > 0)
            "
            class="space-y-2.5"
          >
            <div class="p-3 bg-emerald-50/80 border border-emerald-200/70 rounded-xl text-xs text-emerald-900 flex items-start gap-2.5">
              <span class="text-base leading-none">✅</span>
              <div>
                <strong class="font-bold">Pembayaran Diterima:</strong> Kamar siap dikonfirmasi oleh resepsionis.
              </div>
            </div>
            <div class="flex gap-2">
              <button
                @click="approveReservation(selectedReservation.id)"
                class="flex-1 py-2.5 px-4 rounded-xl font-bold text-xs text-white bg-emerald-600 hover:bg-emerald-700 shadow-md shadow-emerald-600/20 transition-all cursor-pointer"
              >
                ✓ Approve (Receptionist)
              </button>
              <button
                @click="rejectReservation(selectedReservation.id)"
                class="py-2.5 px-4 rounded-xl font-bold text-xs text-slate-700 bg-white border border-slate-200 hover:bg-slate-50 transition-all cursor-pointer"
              >
                Tolak
              </button>
            </div>
          </div>

          <!-- Pending Finance Verification -->
          <div v-else class="space-y-2.5">
            <div class="p-3 bg-amber-50 border border-amber-200/80 rounded-xl text-xs text-amber-900 flex items-start gap-2.5">
              <span class="text-base leading-none">⏳</span>
              <div class="leading-relaxed">
                <strong class="font-bold">Menunggu Tim Finance:</strong> Verifikasi mutasi transfer untuk mengaktifkan persetujuan.
              </div>
            </div>

            <div class="flex gap-2">
              <button
                @click="verifyFinancePayment(selectedReservation.id, 'confirmed')"
                class="flex-1 py-2.5 px-4 rounded-xl font-bold text-xs text-white bg-blue-600 hover:bg-blue-700 shadow-md shadow-blue-600/20 transition-all cursor-pointer"
              >
                💳 Confirm Bayar (Finance)
              </button>
              <button
                @click="rejectReservation(selectedReservation.id)"
                class="py-2.5 px-4 rounded-xl font-bold text-xs text-slate-700 bg-white border border-slate-200 hover:bg-slate-50 transition-all cursor-pointer"
              >
                Tolak
              </button>
            </div>
          </div>
        </div>

        <!-- Confirmed Actions -->
        <button
          v-if="
            selectedReservation.status === 'approve' ||
            selectedReservation.status === 'approved' ||
            selectedReservation.status === 'confirmed'
          "
          @click="handleCheckIn(selectedReservation.id)"
          class="w-full py-3 px-4 rounded-xl font-bold text-sm text-white bg-emerald-600 hover:bg-emerald-700 shadow-lg shadow-emerald-600/25 transition-all cursor-pointer"
        >
          🛎️ Check In Tamu
        </button>

        <button
          v-if="selectedReservation.status === 'checked-in'"
          @click="handleCheckOut(selectedReservation.id)"
          class="w-full py-3 px-4 rounded-xl font-bold text-sm text-white bg-indigo-600 hover:bg-indigo-700 shadow-lg shadow-indigo-600/25 transition-all cursor-pointer"
        >
          🚪 Proses Check Out
        </button>

        <!-- Cancel Action -->
        <button
          v-if="
            selectedReservation.status === 'pending' ||
            selectedReservation.status === 'approve' ||
            selectedReservation.status === 'approved' ||
            selectedReservation.status === 'confirmed'
          "
          @click="cancelReservation(selectedReservation.id)"
          class="w-full mt-2 py-2 px-4 rounded-xl font-semibold text-xs text-rose-600 hover:bg-rose-50 border border-transparent hover:border-rose-100 transition-all cursor-pointer"
        >
          Batalkan Reservasi
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
.btn-new {
  background: #e4793b;
  color: white;
  border: none;
  padding: 10px 18px;
  border-radius: 10px;
  font-size: 0.85rem;
  font-weight: 700;
  cursor: pointer;
  white-space: nowrap;
  transition: background 0.2s;
}
.btn-new:hover {
  background: #d16627;
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
.channel-badge-pill {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  border-radius: 9999px;
  background: #f1f5f9;
  color: #334155;
  font-size: 0.75rem;
  font-weight: 700;
  border: 1px solid #e2e8f0;
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
