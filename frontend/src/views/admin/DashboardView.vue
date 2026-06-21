<script setup>
import { ref, onMounted, watch, nextTick } from 'vue'
import { useDashboardStore } from '@/stores/dashboardStore'
import { reservationService } from '@/services/admin/reservationService'
import VueApexCharts from 'vue3-apexcharts'

const dashboardStore = useDashboardStore()
const reservations = ref([])

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

const breakdownSeries = ref([0, 0, 0]) // Menggunakan 3 elemen sesuai agregasi backend baru
const breakdownOptions = ref({
  chart: { type: 'donut' },
  labels: ['Booked', 'Canceled & Rejected', 'Pending Confirmation'], // Sesuaikan dengan logika repo backend baru
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
            formatter: () => '0', // Akan di-override secara aman di dalam watcher
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

// --- WATCHER: Sinkronisasi Data Reaktif Dan Penyelamatan Render Donut ---
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
        // PERBAIKAN UTAMA: Pastikan total elemen array series ini sama panjang (3) dengan label opsi donut di atas
        breakdownSeries.value = [bd.booked, bd.canceled, bd.pending_confirmation]

        const currentTotal = bd.total_bookings

        // Amankan fungsionalitas total label formatter dari glitch internal ApexCharts
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

// Mengubah range filter secara silent di background tanpa memicu skeleton loader harian
const changeRange = async (range) => {
  await dashboardStore.fetchDashboard(range, false)
}

onMounted(async () => {
  isPageLoading.value = true
  try {
    await Promise.all([
      dashboardStore.fetchDashboard(undefined, true),
      reservationService.getAll().then((data) => (reservations.value = data)),
    ])
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
    <div class="dashboard-header">
      <div>
        <h1>Hotel Overview</h1>
        <p class="subtitle">Welcome back! Here's what's happening today.</p>
      </div>
      <div class="header-actions" v-if="!isPageLoading">
        <div class="time-filters">
          <button
            class="filter-btn"
            :class="{ active: dashboardStore.currentRange === 'today' }"
            @click="changeRange('today')"
          >
            Today
          </button>
          <button
            class="filter-btn"
            :class="{ active: dashboardStore.currentRange === 'weekly' }"
            @click="changeRange('weekly')"
          >
            This Week
          </button>
          <button
            class="filter-btn"
            :class="{ active: dashboardStore.currentRange === 'monthly' }"
            @click="changeRange('monthly')"
          >
            This Month
          </button>
        </div>
        <button class="btn-primary">Export Report</button>
      </div>
    </div>

    <div v-if="isPageLoading" class="skeleton-container">
      <div class="skeleton-grid-3">
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
        <div class="stats-grid">
          <div class="stat-card unique-purple">
            <div class="card-info">
              <span class="card-icon">👥</span>
              <h3>Occupancy Rate</h3>
              <p class="main-val">{{ Math.round(dashboardStore.occupancyRate) }}%</p>
            </div>
            <div class="card-chart">
              <VueApexCharts
                type="radialBar"
                height="130"
                :options="radialOptions('#6366f1')"
                :series="[dashboardStore.occupancyRate]"
              />
            </div>
          </div>

          <div class="stat-card unique-green">
            <div class="card-info">
              <span class="card-icon">🚪</span>
              <h3>Available Rooms</h3>
              <p class="main-val">
                {{ dashboardStore.availableRooms }}
                <span class="sub-label">/ {{ dashboardStore.totalRooms }}</span>
              </p>
            </div>
            <div class="card-chart">
              <VueApexCharts
                type="radialBar"
                height="130"
                :options="radialOptions('#10b981')"
                :series="[
                  dashboardStore.totalRooms > 0
                    ? (dashboardStore.availableRooms / dashboardStore.totalRooms) * 100
                    : 0,
                ]"
              />
            </div>
          </div>

          <div class="stat-card unique-orange">
            <div class="card-info">
              <span class="card-icon">💰</span>
              <h3>Total Revenue</h3>
              <p class="main-val" style="font-size: 1.5rem">
                {{ dashboardStore.formattedRevenue }}
              </p>
            </div>
            <div class="card-chart">
              <VueApexCharts
                type="radialBar"
                height="130"
                :options="radialOptions('#e4793b')"
                :series="[dashboardStore.totalRevenue > 0 ? 100 : 0]"
              />
            </div>
          </div>
        </div>

        <div class="main-split-grid">
          <div class="dashboard-card chart-main-box">
            <div class="card-header">
              <div class="header-left">
                <h2>Total Revenue</h2>
                <select
                  class="select-dropdown"
                  :value="dashboardStore.currentRange"
                  @change="(e) => changeRange(e.target.value)"
                >
                  <option value="today">Today</option>
                  <option value="weekly">This Week</option>
                  <option value="monthly">This Month</option>
                </select>
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
                height="320"
                :options="chartOptions"
                :series="chartSeries"
              />
            </div>
          </div>

          <div class="dashboard-card breakdown-chart">
            <div class="card-header">
              <h2>Booking Breakdown</h2>
              <select
                class="select-dropdown"
                :value="dashboardStore.currentRange"
                @change="(e) => changeRange(e.target.value)"
              >
                <option value="today">Today</option>
                <option value="weekly">This Week</option>
                <option value="monthly">This Month</option>
              </select>
            </div>

            <div class="donut-wrapper">
              <VueApexCharts
                :key="'break-' + chartKey"
                type="donut"
                height="250"
                :options="breakdownOptions"
                :series="breakdownSeries"
              />
            </div>

            <div class="breakdown-legend-list" v-if="dashboardStore.bookingBreakdown">
              <div class="legend-row">
                <div class="legend-label-group">
                  <span class="legend-dot" style="background-color: #e4793b"></span
                  ><span class="label-text">Booked</span>
                </div>
                <div class="legend-value-group">
                  <span class="val-num">{{ dashboardStore.bookingBreakdown.booked }}</span>
                  <span class="val-pct"
                    >{{ dashboardStore.bookingBreakdown.booked_percentage.toFixed(0) }}%</span
                  >
                </div>
              </div>
              <div class="legend-row">
                <div class="legend-label-group">
                  <span class="legend-dot" style="background-color: #fb7185"></span
                  ><span class="label-text">Canceled & Rejected</span>
                </div>
                <div class="legend-value-group">
                  <span class="val-num">{{ dashboardStore.bookingBreakdown.canceled }}</span>
                  <span class="val-pct" style="color: #fb7185"
                    >{{ dashboardStore.bookingBreakdown.canceled_percentage.toFixed(0) }}%</span
                  >
                </div>
              </div>
              <div class="legend-row">
                <div class="legend-label-group">
                  <span class="legend-dot" style="background-color: #6366f1"></span
                  ><span class="label-text">Pending Confirmation</span>
                </div>
                <div class="legend-value-group">
                  <span class="val-num">{{
                    dashboardStore.bookingBreakdown.pending_confirmation
                  }}</span>
                  <span class="val-pct" style="color: #6366f1"
                    >{{ dashboardStore.bookingBreakdown.pending_percentage.toFixed(0) }}%</span
                  >
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>

      <div class="dashboard-card table-card">
        <div class="card-header">
          <h2>Recently Booked Status</h2>
        </div>
        <div class="responsive-table-wrap">
          <table class="premium-table">
            <thead>
              <tr>
                <th>Guest Name</th>
                <th>Room</th>
                <th>Guests</th>
                <th>Check In</th>
                <th>Check Out</th>
                <th class="text-center">Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="res in reservations.slice(0, 4)" :key="res.id">
                <td class="bold-name">{{ res.full_name }}</td>
                <td>
                  <span>{{ res.room?.name || 'N/A' }}</span>
                </td>
                <td>{{ res.number_of_guest || 0 }} Users</td>
                <td>
                  <span class="date-badge-rect">{{ res.checkin_date }}</span>
                </td>
                <td>
                  <span class="date-badge-rect">{{ res.checkout_date }}</span>
                </td>
                <td class="text-center">
                  <span
                    class="status-dot-badge"
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
    </template>
  </div>
</template>

<style scoped>
/* (Seluruh CSS Premium milikmu tetap sama persis seperti sebelumnya) */
.skeleton-container {
  display: flex;
  flex-direction: column;
  gap: 28px;
  width: 100%;
}
.skeleton-grid-3 {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}
.skeleton-split {
  display: grid;
  grid-template-columns: 1.5fr 1fr;
  gap: 24px;
}
.skeleton-card {
  background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%);
  background-size: 200% 100%;
  animation: loading-shimmer 1.4s infinite linear;
  border-radius: 20px;
}
.card-shim {
  height: 132px;
}
.body-large-shim {
  height: 395px;
}
.body-small-shim {
  height: 395px;
}
.table-shim {
  height: 250px;
}

@keyframes loading-shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

.premium-dashboard {
  background-color: #f8fafc;
  padding: 32px;
  display: flex;
  flex-direction: column;
  gap: 28px;
  font-family: 'Plus Jakarta Sans', sans-serif;
}
.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.dashboard-header h1 {
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
.time-filters {
  background: #e2e8f0;
  padding: 4px;
  border-radius: 10px;
  display: inline-flex;
  margin-right: 12px;
}
.filter-btn {
  border: none;
  background: transparent;
  padding: 8px 16px;
  font-size: 0.85rem;
  font-weight: 600;
  border-radius: 8px;
  cursor: pointer;
  color: #64748b;
}
.filter-btn.active {
  background: #ffffff;
  color: #6366f1;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
}
.btn-primary {
  background: #6366f1;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
}
.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}
.main-split-grid {
  display: grid;
  grid-template-columns: 1.5fr 1fr;
  gap: 24px;
  align-items: start;
}
.dashboard-card {
  background: #ffffff;
  border-radius: 20px;
  padding: 24px;
  border: 1px solid rgba(230, 230, 230, 0.7);
}
.stat-card {
  background: #ffffff;
  border-radius: 20px;
  padding: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border: 1px solid #edf2f7;
}
.main-val {
  font-size: 1.75rem;
  font-weight: 700;
  color: #0f172a;
}
.sub-label {
  font-size: 0.8rem;
  color: #94a3b8;
  font-weight: 400;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}
.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}
.card-header h2 {
  font-size: 1.25rem;
  font-weight: 700;
  color: #1e293b;
}
.select-dropdown {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  padding: 8px 32px 8px 16px;
  border-radius: 10px;
  color: #4a5568;
  font-size: 0.9rem;
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' stroke='%234a5568'%3E%3Cpath stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='M19 9l-7 7-7-7'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
  background-size: 14px;
}
.donut-wrapper {
  display: flex;
  justify-content: center;
  margin: 16px 0;
}
.breakdown-legend-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 20px;
}
.legend-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.legend-label-group {
  display: flex;
  align-items: center;
  gap: 10px;
}
.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 3px;
}
.label-text {
  font-size: 0.85rem;
  font-weight: 600;
  color: #64748b;
}
.legend-value-group {
  display: flex;
  align-items: center;
  gap: 12px;
}
.val-num {
  font-size: 0.85rem;
  font-weight: 700;
  color: #1e293b;
}
.val-pct {
  font-size: 0.85rem;
  font-weight: 600;
  color: #10b981;
}
.responsive-table-wrap {
  overflow-x: auto;
}
.premium-table {
  width: 100%;
  border-collapse: collapse;
}
.premium-table th {
  padding: 14px 16px;
  color: #64748b;
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
  border-bottom: 1px solid #f1f5f9;
  background: #f8fafc;
}
.premium-table td {
  padding: 16px;
  font-size: 0.875rem;
  color: #334155;
  border-bottom: 1px solid #f1f5f9;
}
.bold-name {
  font-weight: 600;
  color: #0f172a;
}
.date-badge-rect {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 0.8rem;
  color: #475569;
}
.text-center {
  text-align: center;
}
.status-dot-badge {
  display: inline-flex;
  align-items: center;
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 700;
  text-transform: capitalize;
}
.status-approved {
  background-color: #d1fae5;
  color: #065f46;
}
.status-pending {
  background-color: #fef3c7;
  color: #92400e;
}
.status-cancel {
  background-color: #fee2e2;
  color: #991b1b;
}
</style>
