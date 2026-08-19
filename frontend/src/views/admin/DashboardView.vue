<script setup>
import { ref, onMounted, watch, nextTick, computed } from 'vue'
import { useDashboardStore } from '@/stores/dashboardStore'
import { useAuthStore } from '@/stores/authStore'
import { isManagerRole } from '@/config/access'
import { reservationService } from '@/services/admin/reservationService'
import { roomUnitService } from '@/services/admin/roomUnitService'
import { floorService } from '@/services/admin/floorService'
import VueApexCharts from 'vue3-apexcharts'

const dashboardStore = useDashboardStore()
const authStore = useAuthStore()
const reservations = ref([])
const roomUnits = ref([])
const floors = ref([])

// Indikator loading halaman utama (Hanya aktif sekali di awal)
const isPageLoading = ref(true)

// Atribut Key unik untuk memaksa re-render jika terjadi glitch canvas ApexCharts
const chartKey = ref(0)

// --- State Data untuk ApexCharts ---
const chartSeries = ref([{ name: 'Revenue', data: [] }])
const chartOptions = ref({
  chart: {
    type: 'area',
    toolbar: { show: false },
    zoom: { enabled: false },
    fontFamily: 'inherit',
  },
  colors: ['#e4793b'],
  stroke: { curve: 'smooth', width: 3 },
  fill: {
    type: 'gradient',
    gradient: { shadeIntensity: 1, opacityFrom: 0.25, opacityTo: 0.01 },
  },
  dataLabels: { enabled: false },
  markers: {
    size: 5,
    colors: ['#e4793b'],
    strokeColors: '#fff',
    strokeWidth: 2,
    hover: { size: 7 },
  },
  xaxis: {
    categories: [],
    axisBorder: { show: false },
    axisTicks: { show: false },
    labels: { style: { colors: '#94a3b8', fontWeight: 500 } },
  },
  yaxis: {
    min: 0,
    tickAmount: 5,
    labels: {
      style: { colors: '#94a3b8', fontWeight: 500 },
      formatter: (value) => (value === 0 ? '0' : `Rp ${(value / 1000).toLocaleString('id-ID')}k`),
    },
  },
  grid: {
    borderColor: '#f1f5f9',
    strokeDashArray: 4,
    xaxis: { lines: { show: false } },
    yaxis: { lines: { show: true } },
  },
  tooltip: {
    theme: 'light',
    y: { formatter: (val) => `Rp ${val.toLocaleString('id-ID')}` },
  },
})

const breakdownSeries = ref([0, 0, 0])
const breakdownOptions = ref({
  chart: { type: 'donut' },
  labels: ['Booked', 'Canceled & Rejected', 'Pending Confirmation'],
  colors: ['#e4793b', '#fb7185', '#6366f1'],
  stroke: { show: false },
  dataLabels: { enabled: false },
  legend: { show: false },
  plotOptions: {
    pie: {
      donut: {
        size: '75%',
        labels: {
          show: true,
          name: { show: true, fontSize: '14px', color: '#64748b', offsetY: 20 },
          value: { show: true, fontSize: '24px', fontWeight: 700, color: '#1e293b', offsetY: -20 },
          total: {
            show: true,
            label: 'Total Bookings',
            formatter: () => '0',
          },
        },
      },
    },
  },
  tooltip: { theme: 'dark' },
})

// --- 1. Konfigurasi Radial Gauges (Top Cards) ---
const radialOptions = (color) => ({
  chart: { type: 'radialBar', sparkline: { enabled: true } },
  colors: [color],
  plotOptions: {
    radialBar: {
      startAngle: -90,
      endAngle: 90,
      track: { background: '#f1f5f9', strokeWidth: '85%' },
      dataLabels: { show: false },
    },
  },
})

// --- WATCHER: Sinkronisasi Data Reaktif Dan ApexCharts ---
watch(
  () => dashboardStore.dashboardData,
  (newData) => {
    if (newData) {
      chartSeries.value = [
        {
          name: 'Revenue',
          data: newData.revenue_chart?.map((item) => item.income),
        },
      ]

      chartOptions.value = {
        ...chartOptions.value,
        xaxis: {
          ...chartOptions.value.xaxis,
          categories: newData.revenue_chart?.map((item) => item.month),
        },
      }

      const bd = newData.booking_breakdown
      if (bd) {
        breakdownSeries.value = [bd.booked, bd.canceled, bd.pending_confirmation]
        const currentTotal = bd.total_bookings

        breakdownOptions.value = {
          ...breakdownOptions.value,
          plotOptions: {
            ...breakdownOptions.value.plotOptions,
            pie: {
              ...breakdownOptions.value.plotOptions?.pie,
              donut: {
                ...breakdownOptions.value.plotOptions?.pie?.donut,
                labels: {
                  ...breakdownOptions.value.plotOptions?.pie?.donut?.labels,
                  total: {
                    ...breakdownOptions.value.plotOptions?.pie?.donut?.labels?.total,
                    formatter: () => currentTotal?.toLocaleString('id-ID') || '0',
                  },
                },
              },
            },
          },
        }
      }

      chartKey.value++
    }
  },
  { deep: true, immediate: true },
)

const changeRange = async (range) => {
  await dashboardStore.fetchDashboard(range, false)
}

// --- Live Operational Computed Properties ---
const availableUnitsCount = computed(() => roomUnits.value.filter((u) => u.status === 'available').length)
const occupiedUnitsCount = computed(() => roomUnits.value.filter((u) => u.status === 'occupied').length)
const dirtyUnitsCount = computed(() => roomUnits.value.filter((u) => u.status === 'dirty').length)
const maintenanceUnitsCount = computed(() => roomUnits.value.filter((u) => u.status === 'maintenance').length)

// Group room units by floor for live room rack matrix
const unitsByFloor = computed(() => {
  const map = {}
  floors.value.forEach((f) => {
    map[f.id] = {
      floor: f,
      units: [],
    }
  })

  // Fallback for units without floor or unmapped
  const noFloorUnits = []

  roomUnits.value.forEach((unit) => {
    if (unit.floor_id && map[unit.floor_id]) {
      map[unit.floor_id].units.push(unit)
    } else {
      noFloorUnits.push(unit)
    }
  })

  const result = Object.values(map)
    .filter((g) => g.units.length > 0)
    .sort((a, b) => (a.floor.floor_number || 0) - (b.floor.floor_number || 0))

  if (noFloorUnits.length > 0) {
    result.push({
      floor: { id: 'other', name: 'Lainnya', floor_number: 99 },
      units: noFloorUnits,
    })
  }

  // Sort units inside each floor
  result.forEach((g) => {
    g.units.sort((a, b) => a.room_number.localeCompare(b.room_number, undefined, { numeric: true }))
  })

  return result
})

// Recent Activity Log Stream (Check In & Check Out)
const recentActivityLogs = computed(() => {
  const logs = []
  reservations.value.forEach((res) => {
    if (res.logs && Array.isArray(res.logs)) {
      res.logs.forEach((log) => {
        logs.push({
          id: log.id || `${res.id}-${log.timestamp}`,
          guest_name: res.full_name,
          room_name: res.room?.name || 'Standard Room',
          room_number: res.room_unit?.room_number || '',
          action: log.action,
          timestamp: log.timestamp,
          is_early: !!log.is_early,
          notes: log.notes,
        })
      })
    }
  })
  return logs.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()).slice(0, 5)
})

onMounted(async () => {
  isPageLoading.value = true
  try {
    const [_, resData, unitsData, floorsData] = await Promise.all([
      dashboardStore.fetchDashboard(undefined, true),
      reservationService.getAll(),
      roomUnitService.getAll(),
      floorService.getAll(),
    ])
    reservations.value = resData || []
    roomUnits.value = unitsData || []
    floors.value = floorsData || []
  } catch (error) {
    console.error('Gagal memuat dashboard:', error)
  } finally {
    await nextTick()
    isPageLoading.value = false
  }
})
</script>

<template>
  <div class="premium-dashboard">
    <!-- Header -->
    <div class="dashboard-header">
      <div>
        <h1 class="font-extrabold text-2xl text-slate-900 tracking-tight flex items-center gap-2.5">
          <span v-if="isManagerRole(authStore.role)" class="text-amber-500">👑</span>
          {{ isManagerRole(authStore.role) ? 'Executive Hotel Manager Overview' : 'Hotel Operations Dashboard' }}
        </h1>
        <p class="subtitle text-xs text-slate-500 font-medium">
          {{ isManagerRole(authStore.role) 
              ? 'Laporan Eksekutif KPI Hotel, Tingkat Hunian (Occupancy Rate), Revenue & Pengawasan Operasional Staf.' 
              : 'Monitoring real-time ketersediaan kamar, operasional check-in/out, dan pendapatan.' }}
        </p>
      </div>
      <div class="header-actions" v-if="!isPageLoading">
        <div class="time-filters">
          <button
            class="filter-btn"
            :class="{ active: dashboardStore.currentRange === 'today' }"
            @click="changeRange('today')"
          >
            Hari Ini
          </button>
          <button
            class="filter-btn"
            :class="{ active: dashboardStore.currentRange === 'weekly' }"
            @click="changeRange('weekly')"
          >
            Minggu Ini
          </button>
          <button
            class="filter-btn"
            :class="{ active: dashboardStore.currentRange === 'monthly' }"
            @click="changeRange('monthly')"
          >
            Bulan Ini
          </button>
        </div>
      </div>
    </div>

    <!-- Manager Supervisory Mode Banner -->
    <div
      v-if="!isPageLoading && isManagerRole(authStore.role)"
      class="p-4 bg-gradient-to-r from-slate-900 via-indigo-950 to-slate-900 text-white rounded-2xl flex items-center justify-between shadow-lg border border-indigo-900/40 animate-fade-in"
    >
      <div class="flex items-center gap-3.5">
        <div class="w-10 h-10 rounded-xl bg-amber-400/20 text-amber-300 flex items-center justify-center text-xl shrink-0 border border-amber-400/30">
          👑
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h3 class="font-extrabold text-sm text-white">Mode Pengawasan Hotel Manager (Read-Only)</h3>
            <span class="px-2 py-0.5 rounded-full text-[10px] font-black bg-amber-400/20 text-amber-300 border border-amber-400/30">
              Supervisor Access
            </span>
          </div>
          <p class="text-xs text-slate-300 mt-0.5 leading-relaxed">
            Anda memiliki akses monitoring penuh terhadap seluruh operasional front desk, matriks ketersediaan kamar, analitik tingkat hunian, dan laporan keuangan hotel.
          </p>
        </div>
      </div>
    </div>

    <!-- Skeleton Loader -->
    <div v-if="isPageLoading" class="skeleton-container">
      <div class="skeleton-grid-4">
        <div class="skeleton-card card-shim"></div>
        <div class="skeleton-card card-shim"></div>
        <div class="skeleton-card card-shim"></div>
        <div class="skeleton-card card-shim"></div>
      </div>
      <div class="skeleton-split">
        <div class="skeleton-card body-large-shim"></div>
        <div class="skeleton-card body-small-shim"></div>
      </div>
      <div class="skeleton-card table-shim"></div>
    </div>

    <template v-else>
      <div v-if="dashboardStore.errorMsg" class="error-state">
        {{ dashboardStore.errorMsg }}
      </div>

      <template v-else>
        <!-- Top Stats Row (4 Cards) -->
        <div class="stats-grid-4">
          <!-- 1. Occupancy Rate -->
          <div class="stat-card unique-purple">
            <div class="card-info">
              <span class="card-icon">👥</span>
              <h3>Tingkat Hunian (Occupancy)</h3>
              <p class="main-val">{{ Math.round(dashboardStore.occupancyRate) }}%</p>
              <span class="text-[11px] text-indigo-700 font-semibold mt-1">
                {{ occupiedUnitsCount }} terisi dari {{ roomUnits.length || dashboardStore.totalRooms }} unit
              </span>
            </div>
            <div class="card-chart">
              <VueApexCharts
                type="radialBar"
                height="125"
                :options="radialOptions('#6366f1')"
                :series="[dashboardStore.occupancyRate]"
              />
            </div>
          </div>

          <!-- 2. Available Rooms -->
          <div class="stat-card unique-green">
            <div class="card-info">
              <span class="card-icon">✨</span>
              <h3>Kamar Siap Pakai (Ready)</h3>
              <p class="main-val">
                {{ availableUnitsCount }}
                <span class="sub-label">/ {{ roomUnits.length }} Unit</span>
              </p>
              <span class="text-[11px] text-emerald-700 font-semibold mt-1">
                Tersedia untuk alokasi otomatis
              </span>
            </div>
            <div class="card-chart">
              <VueApexCharts
                type="radialBar"
                height="125"
                :options="radialOptions('#10b981')"
                :series="[
                  roomUnits.length > 0
                    ? (availableUnitsCount / roomUnits.length) * 100
                    : 0,
                ]"
              />
            </div>
          </div>

          <!-- 3. Housekeeping Alert -->
          <div class="stat-card unique-amber">
            <div class="card-info">
              <span class="card-icon">🧹</span>
              <h3>Housekeeping (Pembersihan)</h3>
              <p class="main-val text-amber-700">{{ dirtyUnitsCount }} <span class="sub-label">Kamar</span></p>
              <span class="text-[11px] text-amber-800 font-semibold mt-1">
                {{ maintenanceUnitsCount > 0 ? `${maintenanceUnitsCount} unit rusak (maintenance)` : 'Tidak ada kamar rusak' }}
              </span>
            </div>
            <div class="card-chart">
              <VueApexCharts
                type="radialBar"
                height="125"
                :options="radialOptions('#f59e0b')"
                :series="[roomUnits.length > 0 ? (dirtyUnitsCount / roomUnits.length) * 100 : 0]"
              />
            </div>
          </div>

          <!-- 4. Total Revenue -->
          <div class="stat-card unique-orange">
            <div class="card-info">
              <span class="card-icon">💰</span>
              <h3>Total Pendapatan (Revenue)</h3>
              <p class="main-val text-xl font-black">
                {{ dashboardStore.formattedRevenue }}
              </p>
              <span class="text-[11px] text-amber-900 font-semibold mt-1">
                Periode: {{ dashboardStore.currentRange }}
              </span>
            </div>
            <div class="card-chart">
              <VueApexCharts
                type="radialBar"
                height="125"
                :options="radialOptions('#e4793b')"
                :series="[dashboardStore.totalRevenue > 0 ? 100 : 0]"
              />
            </div>
          </div>
        </div>

        <!-- Live Physical Room Matrix (Denah Unit Kamar Real-time) -->
        <div class="dashboard-card p-5 space-y-4">
          <div class="flex items-center justify-between flex-wrap gap-2 pb-3 border-b border-slate-100">
            <div class="flex items-center gap-2.5">
              <div class="w-8 h-8 rounded-lg bg-indigo-50 text-indigo-600 flex items-center justify-center font-bold text-sm">
                🏢
              </div>
              <div>
                <h2 class="font-extrabold text-sm text-slate-900">Status Kamar Fisik Real-time (Live Room Rack)</h2>
                <p class="text-xs text-slate-400">Monitoring nomor kamar hotel per lantai secara visual.</p>
              </div>
            </div>

            <!-- Legend -->
            <div class="flex items-center gap-3 text-xs font-bold flex-wrap">
              <span class="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200">
                <span class="w-2 h-2 rounded-full bg-emerald-500"></span> Available ({{ availableUnitsCount }})
              </span>
              <span class="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-indigo-50 text-indigo-700 border border-indigo-200">
                <span class="w-2 h-2 rounded-full bg-indigo-500"></span> Occupied ({{ occupiedUnitsCount }})
              </span>
              <span class="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-amber-50 text-amber-700 border border-amber-200">
                <span class="w-2 h-2 rounded-full bg-amber-500"></span> Dirty ({{ dirtyUnitsCount }})
              </span>
              <span class="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-rose-50 text-rose-700 border border-rose-200">
                <span class="w-2 h-2 rounded-full bg-rose-500"></span> Maintenance ({{ maintenanceUnitsCount }})
              </span>
            </div>
          </div>

          <!-- Floor Rows -->
          <div v-if="unitsByFloor.length > 0" class="space-y-3">
            <div
              v-for="group in unitsByFloor"
              :key="group.floor.id"
              class="flex items-start gap-4 p-3 bg-slate-50/80 rounded-xl border border-slate-100"
            >
              <div class="w-24 shrink-0 pt-1">
                <span class="text-xs font-black text-slate-700 uppercase tracking-wider block">
                  {{ group.floor.name }}
                </span>
                <span class="text-[10px] text-slate-400 font-medium">Lantai {{ group.floor.floor_number }}</span>
              </div>

              <div class="flex flex-wrap gap-2 flex-1">
                <div
                  v-for="unit in group.units"
                  :key="unit.id"
                  class="px-3 py-1.5 rounded-xl text-xs font-black flex items-center gap-1.5 transition-all shadow-xs border"
                  :class="{
                    'bg-emerald-500 text-white border-emerald-600': unit.status === 'available',
                    'bg-indigo-600 text-white border-indigo-700': unit.status === 'occupied',
                    'bg-amber-400 text-amber-950 border-amber-500': unit.status === 'dirty',
                    'bg-rose-500 text-white border-rose-600': unit.status === 'maintenance',
                  }"
                  :title="`${unit.room_number} - ${unit.room?.name || ''} (${unit.status})`"
                >
                  <span>🚪</span>
                  <span>{{ unit.room_number }}</span>
                  <span class="text-[10px] opacity-85 font-semibold">
                    {{ unit.status === 'available' ? 'Ready' : unit.status === 'occupied' ? 'Terisi' : unit.status === 'dirty' ? 'Dirty' : 'Rusak' }}
                  </span>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="text-center py-4 text-xs text-slate-400">
            Belum ada data unit kamar fisik. Daftarkan di menu Room Units.
          </div>
        </div>

        <!-- Charts Split Row -->
        <div class="main-split-grid">
          <!-- Total Revenue Chart -->
          <div class="dashboard-card chart-main-box">
            <div class="card-header">
              <div class="header-left">
                <h2 class="font-extrabold text-sm text-slate-900">Grafik Pendapatan (Revenue)</h2>
              </div>
              <div class="custom-legend">
                <div class="legend-item">
                  <span class="dot income-dot"></span><span class="legend-text">Income</span>
                </div>
              </div>
            </div>
            <div class="chart-wrapper">
              <VueApexCharts
                :key="'rev-' + chartKey"
                type="area"
                height="300"
                :options="chartOptions"
                :series="chartSeries"
              />
            </div>
          </div>

          <!-- Booking Breakdown Donut -->
          <div class="dashboard-card breakdown-chart">
            <div class="card-header">
              <h2 class="font-extrabold text-sm text-slate-900">Distribusi Status Booking</h2>
            </div>

            <div class="donut-wrapper">
              <VueApexCharts
                :key="'break-' + chartKey"
                type="donut"
                height="240"
                :options="breakdownOptions"
                :series="breakdownSeries"
              />
            </div>

            <div class="breakdown-legend-list" v-if="dashboardStore.bookingBreakdown">
              <div class="legend-row">
                <div class="legend-label-group">
                  <span class="legend-dot" style="background-color: #e4793b"></span>
                  <span class="label-text">Booked & Confirmed</span>
                </div>
                <div class="legend-value-group">
                  <span class="val-num">{{ dashboardStore.bookingBreakdown.booked }}</span>
                  <span class="val-pct">{{ dashboardStore.bookingBreakdown.booked_percentage.toFixed(0) }}%</span>
                </div>
              </div>
              <div class="legend-row">
                <div class="legend-label-group">
                  <span class="legend-dot" style="background-color: #fb7185"></span>
                  <span class="label-text">Canceled & Rejected</span>
                </div>
                <div class="legend-value-group">
                  <span class="val-num">{{ dashboardStore.bookingBreakdown.canceled }}</span>
                  <span class="val-pct" style="color: #fb7185">{{ dashboardStore.bookingBreakdown.canceled_percentage.toFixed(0) }}%</span>
                </div>
              </div>
              <div class="legend-row">
                <div class="legend-label-group">
                  <span class="legend-dot" style="background-color: #6366f1"></span>
                  <span class="label-text">Pending Finance / Reception</span>
                </div>
                <div class="legend-value-group">
                  <span class="val-num">{{ dashboardStore.bookingBreakdown.pending_confirmation }}</span>
                  <span class="val-pct" style="color: #6366f1">{{ dashboardStore.bookingBreakdown.pending_percentage.toFixed(0) }}%</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Split Row: Recent Activities & Recently Booked -->
        <div class="main-split-grid">
          <!-- Recent Check-In / Check-Out Log Feed -->
          <div class="dashboard-card p-5 space-y-3">
            <div class="flex items-center justify-between pb-2 border-b border-slate-100">
              <div class="flex items-center gap-2">
                <span class="text-base">📋</span>
                <h2 class="font-extrabold text-sm text-slate-900">Aktivitas In / Out Terkini</h2>
              </div>
              <router-link to="/admin/guestbook" class="text-xs font-bold text-indigo-600 hover:text-indigo-800">
                Lihat Semua →
              </router-link>
            </div>

            <div v-if="recentActivityLogs.length > 0" class="space-y-2.5">
              <div
                v-for="log in recentActivityLogs"
                :key="log.id"
                class="p-2.5 bg-slate-50 hover:bg-slate-100/80 rounded-xl transition-colors flex items-center justify-between gap-3 text-xs"
              >
                <div class="flex items-center gap-2.5">
                  <span class="text-base leading-none">{{ log.action === 'check_in' ? '🛎️' : '🚪' }}</span>
                  <div>
                    <strong class="text-slate-900 font-bold block">{{ log.guest_name }}</strong>
                    <span class="text-[11px] text-slate-500">
                      {{ log.room_number ? `Kamar ${log.room_number} • ` : '' }}{{ log.room_name }}
                    </span>
                  </div>
                </div>

                <div class="text-right shrink-0">
                  <span
                    class="px-2 py-0.5 rounded-full text-[10px] font-black"
                    :class="log.action === 'check_in' ? 'bg-emerald-100 text-emerald-800' : 'bg-slate-200 text-slate-700'"
                  >
                    {{ log.is_early ? 'Early Check-In' : log.action === 'check_in' ? 'Check In' : 'Check Out' }}
                  </span>
                  <div class="text-[10px] text-slate-400 mt-0.5">
                    {{ new Date(log.timestamp).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) }} WIB
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="text-center py-6 text-xs text-slate-400">
              Belum ada log check-in/out terbaru.
            </div>
          </div>

          <!-- Recently Booked Status -->
          <div class="dashboard-card p-5 space-y-3">
            <div class="flex items-center justify-between pb-2 border-b border-slate-100">
              <div class="flex items-center gap-2">
                <span class="text-base">📅</span>
                <h2 class="font-extrabold text-sm text-slate-900">Reservasi Terbaru</h2>
              </div>
              <router-link to="/admin/reservations" class="text-xs font-bold text-indigo-600 hover:text-indigo-800">
                Lihat Semua →
              </router-link>
            </div>

            <div class="responsive-table-wrap">
              <table class="premium-table">
                <thead>
                  <tr>
                    <th>Tamu</th>
                    <th>Kamar</th>
                    <th>Tgl In - Out</th>
                    <th class="text-center">Status</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="res in reservations.slice(0, 4)" :key="res.id">
                    <td class="bold-name text-xs">{{ res.full_name }}</td>
                    <td class="text-xs">
                      <div>{{ res.room?.name || 'N/A' }}</div>
                      <span v-if="res.room_unit?.room_number" class="text-[10px] font-black text-indigo-600">
                        🚪 No. {{ res.room_unit.room_number }}
                      </span>
                    </td>
                    <td class="text-[11px] text-slate-600 font-medium">
                      {{ res.checkin_date }} s/d {{ res.checkout_date }}
                    </td>
                    <td class="text-center">
                      <span
                        class="status-dot-badge text-[10px]"
                        :class="{
                          'status-approved':
                            res.status === 'approve' ||
                            res.status === 'approved' ||
                            res.status === 'checked-in' ||
                            res.status === 'checked-out',
                          'status-pending': res.status === 'pending',
                          'status-cancel': res.status === 'cancel' || res.status === 'rejected',
                        }"
                      >
                        {{ res.status }}
                      </span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </template>
    </template>
  </div>
</template>

<style scoped>
.premium-dashboard {
  display: flex;
  flex-direction: column;
  gap: 20px;
  width: 100%;
  font-family: 'Plus Jakarta Sans', sans-serif;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
}

.time-filters {
  display: inline-flex;
  background: rgba(250, 235, 198, 0.45);
  padding: 4px;
  border-radius: 14px;
  gap: 4px;
  border: 1px solid rgba(250, 235, 198, 0.95);
  box-shadow: 0 2px 8px rgba(180, 140, 60, 0.05);
}

.filter-btn {
  padding: 7px 14px;
  font-size: 12px;
  font-weight: 700;
  border-radius: 10px;
  border: none;
  background: transparent;
  color: #78623b;
  cursor: pointer;
  transition: all 0.2s ease;
}

.filter-btn:hover {
  color: #3a2806;
}

.filter-btn.active {
  background: #ffffff;
  color: #8c6716;
  box-shadow: 0 2px 8px rgba(180, 140, 60, 0.18);
}

.stats-grid-4 {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 18px;
}

.stat-card {
  background: #ffffff;
  padding: 20px 22px;
  border-radius: 22px;
  border: 1px solid rgba(250, 235, 198, 0.85);
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 
    0 8px 24px -6px rgba(180, 140, 60, 0.07),
    0 0 0 1px rgba(255, 255, 255, 0.8) inset;
  transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1), box-shadow 0.25s ease;
}

.stat-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 14px 28px -6px rgba(180, 140, 60, 0.15);
}

.stat-card.unique-purple {
  background: linear-gradient(135deg, rgba(250, 235, 198, 0.4) 0%, #ffffff 100%);
  border-color: rgba(250, 235, 198, 0.95);
}

.stat-card.unique-green {
  background: linear-gradient(135deg, rgba(236, 253, 245, 0.7) 0%, #ffffff 100%);
  border-color: rgba(16, 185, 129, 0.3);
}

.stat-card.unique-amber {
  background: linear-gradient(135deg, rgba(254, 243, 199, 0.6) 0%, #ffffff 100%);
  border-color: rgba(245, 158, 11, 0.3);
}

.stat-card.unique-orange {
  background: linear-gradient(135deg, #faebc6 0%, #f5db97 100%);
  border-color: rgba(180, 140, 50, 0.45);
  box-shadow: 0 8px 28px -4px rgba(180, 140, 60, 0.22);
}

.card-info {
  display: flex;
  flex-direction: column;
}

.card-icon {
  font-size: 22px;
  margin-bottom: 8px;
}

.card-info h3 {
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #7c6848;
  margin: 0;
}

.card-info .main-val {
  font-size: 24px;
  font-weight: 900;
  color: #1a1612;
  margin: 4px 0 0 0;
  line-height: 1.2;
}

.sub-label {
  font-size: 13px;
  font-weight: 600;
  color: #8c7a62;
}

.card-chart {
  width: 90px;
  height: 90px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.main-split-grid {
  display: grid;
  grid-template-columns: 1.6fr 1fr;
  gap: 18px;
}

@media (max-width: 1024px) {
  .main-split-grid {
    grid-template-columns: 1fr;
  }
}

.dashboard-card {
  background: #ffffff;
  border-radius: 22px;
  border: 1px solid rgba(250, 235, 198, 0.85);
  box-shadow: 
    0 8px 24px -6px rgba(180, 140, 60, 0.06),
    0 0 0 1px rgba(255, 255, 255, 0.8) inset;
}

.card-header {
  padding: 16px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid rgba(250, 235, 198, 0.6);
  background: rgba(250, 235, 198, 0.12);
  border-radius: 22px 22px 0 0;
}

.chart-wrapper {
  padding: 12px 16px;
}

.donut-wrapper {
  padding: 12px;
}

.breakdown-legend-list {
  padding: 0 20px 18px 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.legend-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.legend-label-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.legend-dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
}

.val-num {
  font-weight: 700;
  color: #1a1612;
  margin-right: 6px;
}

.val-pct {
  font-weight: 700;
  color: #7c6848;
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
  padding: 12px 16px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  color: #735928;
  background: rgba(250, 235, 198, 0.3);
  border-bottom: 1px solid rgba(250, 235, 198, 0.7);
}

.premium-table td {
  padding: 12px 16px;
  font-size: 12px;
  border-bottom: 1px solid rgba(250, 235, 198, 0.35);
  color: #1e1711;
}

.status-dot-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 9999px;
  font-weight: 700;
}

.status-approved {
  background: #ecfdf5;
  color: #059669;
  border: 1px solid rgba(5, 150, 105, 0.2);
}

.status-pending {
  background: #faebc6;
  color: #8c6716;
  border: 1px solid rgba(140, 103, 22, 0.25);
}

.status-cancel {
  background: #fef2f2;
  color: #e11d48;
  border: 1px solid rgba(225, 29, 72, 0.2);
}

/* Skeleton Loaders */
.skeleton-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.skeleton-grid-4 {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.skeleton-split {
  display: grid;
  grid-template-columns: 1.6fr 1fr;
  gap: 16px;
}

.skeleton-card {
  background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%);
  background-size: 200% 100%;
  animation: loading-shimmer 1.4s infinite linear;
  border-radius: 20px;
}

.card-shim {
  height: 120px;
}

.body-large-shim {
  height: 320px;
}

.body-small-shim {
  height: 320px;
}

.table-shim {
  height: 200px;
}

@keyframes loading-shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>
