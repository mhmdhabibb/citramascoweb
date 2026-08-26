<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { serviceRequestService } from '@/services/admin/serviceRequestService'
import { userService } from '@/services/admin/userService'
import { useToastStore } from '@/stores/toastStore'
import type { ServiceRequest, User } from '@/types'
import VueApexCharts from 'vue3-apexcharts'
import type { ApexOptions } from 'apexcharts'
import { Users, Sparkles, Clock, RotateCcw } from 'lucide-vue-next'

const toastStore = useToastStore()
const loading = ref(true)
const chartKey = ref(0)

const requests = ref<ServiceRequest[]>([])
const housekeepers = ref<User[]>([])

// Period filter
const periodOptions = [
  { label: 'Last 7 Days', value: 7 },
  { label: 'Last 30 Days', value: 30 },
  { label: 'Last 90 Days', value: 90 },
]
const selectedPeriod = ref(30)
const selectedStaffId = ref('all')

// SLA targets in minutes per priority
const SLA_TARGETS: Record<string, number> = {
  urgent: 30,
  high: 60,
  medium: 180,
  low: 360,
}

let pollingTimer: any = null

// ---------- DATA FETCHING ----------

const fetchData = async (silent = false) => {
  try {
    if (!silent) loading.value = true
    const [reqData, hkData] = await Promise.all([
      serviceRequestService.getAll(),
      userService.getByRole('housekeeping'),
    ])
    requests.value = reqData || []
    housekeepers.value = hkData || []
  } catch (e: any) {
    if (!silent) toastStore.error(e.message || 'Failed to load housekeeping performance data')
  } finally {
    if (!silent) loading.value = false
  }
}

onMounted(() => {
  fetchData(false)
  pollingTimer = setInterval(() => fetchData(true), 15000)
})

onUnmounted(() => {
  if (pollingTimer) clearInterval(pollingTimer)
})

// ---------- FILTERED DATA ----------

const filteredRequests = computed(() => {
  const cutoff = new Date()
  cutoff.setDate(cutoff.getDate() - selectedPeriod.value)
  return requests.value.filter((r) => {
    const matchesPeriod = new Date(r.created_at) >= cutoff
    const matchesStaff = selectedStaffId.value === 'all' || r.assigned_to_user_id === selectedStaffId.value
    return matchesPeriod && matchesStaff
  })
})

const completedRequests = computed(() =>
  filteredRequests.value.filter((r) => r.status === 'completed' && r.completed_at),
)

// ---------- SLA CALCULATIONS ----------

function getResolutionMinutes(r: ServiceRequest): number {
  if (!r.completed_at) return Infinity
  return (new Date(r.completed_at).getTime() - new Date(r.created_at).getTime()) / 60000
}

function isWithinSLA(r: ServiceRequest): boolean {
  const target = SLA_TARGETS[r.priority] || 180
  return getResolutionMinutes(r) <= target
}

// Team SLA compliance %
const slaComplianceRate = computed(() => {
  if (completedRequests.value.length === 0) return 0
  const onTime = completedRequests.value.filter(isWithinSLA).length
  return Math.round((onTime / completedRequests.value.length) * 100)
})

// Average resolution time in minutes
const avgResolutionMinutes = computed(() => {
  if (completedRequests.value.length === 0) return 0
  const total = completedRequests.value.reduce((sum, r) => sum + getResolutionMinutes(r), 0)
  return Math.round(total / completedRequests.value.length)
})

// Format minutes for display
function formatDuration(mins: number): string {
  if (mins < 60) return `${mins}m`
  const h = Math.floor(mins / 60)
  const m = mins % 60
  return m > 0 ? `${h}h ${m}m` : `${h}h`
}

// Avg resolution as percentage of a 6h benchmark (for radial display)
const avgResolutionPercent = computed(() => {
  if (avgResolutionMinutes.value === 0) return 0
  return Math.max(0, Math.min(100, Math.round((1 - avgResolutionMinutes.value / 360) * 100)))
})

// ---------- STATUS DISTRIBUTION ----------

const statusCounts = computed(() => {
  const map: Record<string, number> = {
    pending_reception: 0,
    assigned_to_housekeeping: 0,
    in_progress: 0,
    completed: 0,
    cancelled: 0,
  }
  filteredRequests.value.forEach((r) => {
    if (r.status in map) {
      map[r.status] = (map[r.status] ?? 0) + 1
    }
  })
  return map
})

// ---------- WEEKLY TREND ----------

const weeklyTrendData = computed(() => {
  const weeks: number[] = []
  const labels: string[] = []
  const now = new Date()

  for (let i = 5; i >= 0; i--) {
    const weekStart = new Date(now)
    weekStart.setDate(now.getDate() - i * 7)
    const weekEnd = new Date(weekStart)
    weekEnd.setDate(weekStart.getDate() + 7)

    const count = requests.value.filter((r) => {
      const d = new Date(r.created_at)
      return d >= weekStart && d < weekEnd
    }).length

    weeks.push(count)
    labels.push(
      weekStart.toLocaleDateString('en-US', { day: 'numeric', month: 'short' }),
    )
  }

  return { series: weeks, labels }
})

// ---------- CATEGORY BREAKDOWN ----------

const categoryLabels: Record<string, string> = {
  incident_broken_item: 'Broken Items',
  extra_cleaning: 'Extra Cleaning',
  amenities_request: 'Amenities',
  maintenance_repair: 'Maintenance',
  other: 'Other Services',
}

const categoryAvgTimes = computed(() => {
  const map: Record<string, { total: number; count: number }> = {}
  completedRequests.value.forEach((r) => {
    const cat = r.category || 'other'
    const current = map[cat] || { total: 0, count: 0 }
    current.total += getResolutionMinutes(r)
    current.count++
    map[cat] = current
  })

  const cats = Object.keys(map)
  return {
    categories: cats.map((c) => categoryLabels[c] || c),
    values: cats.map((c) => {
      const item = map[c]
      return item && item.count > 0 ? Math.round(item.total / item.count) : 0
    }),
  }
})

// ---------- INDIVIDUAL PERFORMANCE ----------

interface StaffPerformance {
  user: User
  totalTasks: number
  completedTasks: number
  activeTasks: number
  slaRate: number
  avgTime: number
}

const staffPerformance = computed<StaffPerformance[]>(() => {
  return housekeepers.value
    .map((hk) => {
      const tasks = filteredRequests.value.filter(
        (r) => r.assigned_to_user_id === hk.id,
      )
      const completed = tasks.filter(
        (r) => r.status === 'completed' && r.completed_at,
      )
      const active = tasks.filter(
        (r) => r.status === 'in_progress' || r.status === 'assigned_to_housekeeping',
      )
      const onTime = completed.filter(isWithinSLA).length
      const slaRate = completed.length > 0 ? Math.round((onTime / completed.length) * 100) : 0
      const avgTime =
        completed.length > 0
          ? Math.round(
              completed.reduce((s, r) => s + getResolutionMinutes(r), 0) / completed.length,
            )
          : 0

      return {
        user: hk,
        totalTasks: tasks.length,
        completedTasks: completed.length,
        activeTasks: active.length,
        slaRate,
        avgTime,
      }
    })
    .sort((a, b) => b.slaRate - a.slaRate)
})

// Staff ranking chart data
const staffRankingData = computed(() => ({
  names: staffPerformance.value.map(
    (s) => `${s.user.first_name} ${s.user.last_name?.charAt(0) || ''}.`,
  ),
  rates: staffPerformance.value.map((s) => s.slaRate),
}))

// ---------- APEX CHART CONFIGS ----------

// 1) SLA Gauge
const slaGaugeOptions = computed<ApexOptions>(() => ({
  chart: { type: 'radialBar', sparkline: { enabled: true } },
  colors: [slaComplianceRate.value >= 80 ? '#10b981' : slaComplianceRate.value >= 60 ? '#f59e0b' : '#ef4444'],
  plotOptions: {
    radialBar: {
      startAngle: -135,
      endAngle: 135,
      hollow: { size: '62%' },
      track: { background: '#f1f5f9', strokeWidth: '100%' },
      dataLabels: {
        name: { show: true, fontSize: '12px', color: '#94a3b8', offsetY: 22, fontWeight: 500 },
        value: {
          show: true,
          fontSize: '28px',
          fontWeight: 700,
          color: '#1e293b',
          offsetY: -14,
          formatter: (val: number) => `${val}%`,
        },
      },
    },
  },
  labels: ['On-Target'],
  stroke: { lineCap: 'round' },
}))
const slaGaugeSeries = computed(() => [slaComplianceRate.value])

// 2) Avg Resolution Gauge
const avgGaugeOptions = computed<ApexOptions>(() => ({
  chart: { type: 'radialBar', sparkline: { enabled: true } },
  colors: [avgResolutionPercent.value >= 70 ? '#10b981' : avgResolutionPercent.value >= 40 ? '#f59e0b' : '#ef4444'],
  plotOptions: {
    radialBar: {
      startAngle: -135,
      endAngle: 135,
      hollow: { size: '62%' },
      track: { background: '#f1f5f9', strokeWidth: '100%' },
      dataLabels: {
        name: { show: true, fontSize: '12px', color: '#94a3b8', offsetY: 22, fontWeight: 500 },
        value: {
          show: true,
          fontSize: '24px',
          fontWeight: 700,
          color: '#1e293b',
          offsetY: -14,
          formatter: () => formatDuration(avgResolutionMinutes.value),
        },
      },
    },
  },
  labels: ['Average Time'],
  stroke: { lineCap: 'round' },
}))
const avgGaugeSeries = computed(() => [avgResolutionPercent.value])

// 3) Status Donut
const statusDonutOptions = computed<ApexOptions>(() => ({
  chart: { type: 'donut' },
  labels: ['Pending', 'Assigned', 'In Progress', 'Completed', 'Cancelled'],
  colors: ['#94a3b8', '#6366f1', '#f59e0b', '#10b981', '#ef4444'],
  stroke: { show: false },
  dataLabels: { enabled: false },
  legend: { show: false },
  plotOptions: {
    pie: {
      donut: {
        size: '72%',
        labels: {
          show: true,
          name: { show: true, fontSize: '11px', color: '#94a3b8', offsetY: 18 },
          value: { show: true, fontSize: '22px', fontWeight: 700, color: '#1e293b', offsetY: -12 },
          total: {
            show: true,
            label: 'Total Tasks',
            fontSize: '11px',
            color: '#94a3b8',
            formatter: () => String(filteredRequests.value.length),
          },
        },
      },
    },
  },
  tooltip: { theme: 'light' },
}))
const statusDonutSeries = computed(() => [
  statusCounts.value.pending_reception,
  statusCounts.value.assigned_to_housekeeping,
  statusCounts.value.in_progress,
  statusCounts.value.completed,
  statusCounts.value.cancelled,
])

// 4) Weekly Trend Sparkline
const trendOptions = computed<ApexOptions>(() => ({
  chart: {
    type: 'area',
    sparkline: { enabled: false },
    toolbar: { show: false },
    zoom: { enabled: false },
    fontFamily: 'inherit',
  },
  colors: ['#6366f1'],
  stroke: { curve: 'smooth', width: 2.5 },
  fill: {
    type: 'gradient',
    gradient: { shadeIntensity: 1, opacityFrom: 0.25, opacityTo: 0.02 },
  },
  dataLabels: { enabled: false },
  markers: {
    size: 4,
    colors: ['#6366f1'],
    strokeColors: '#fff',
    strokeWidth: 2,
    hover: { size: 6 },
  },
  xaxis: {
    categories: weeklyTrendData.value.labels,
    axisBorder: { show: false },
    axisTicks: { show: false },
    labels: { style: { colors: '#94a3b8', fontSize: '10px' } },
  },
  yaxis: {
    min: 0,
    labels: { show: false },
  },
  grid: {
    borderColor: '#f1f5f9',
    strokeDashArray: 4,
    xaxis: { lines: { show: false } },
    yaxis: { lines: { show: true } },
    padding: { left: 4, right: 4, top: 0, bottom: 0 },
  },
  tooltip: {
    theme: 'light',
    y: { formatter: (val: number) => `${val} tasks` },
  },
}))
const trendSeries = computed(() => [
  { name: 'Incoming Tasks', data: weeklyTrendData.value.series },
])

// 5) Category Bar
const categoryBarOptions = computed<ApexOptions>(() => ({
  chart: {
    type: 'bar',
    toolbar: { show: false },
    fontFamily: 'inherit',
  },
  colors: ['#e4793b'],
  plotOptions: {
    bar: {
      horizontal: true,
      borderRadius: 6,
      barHeight: '55%',
    },
  },
  dataLabels: {
    enabled: true,
    formatter: (val: number) => formatDuration(val),
    style: { fontSize: '11px', fontWeight: 600, colors: ['#fff'] },
    offsetX: -4,
  },
  xaxis: {
    categories: categoryAvgTimes.value.categories,
    labels: {
      style: { colors: '#94a3b8', fontSize: '11px' },
      formatter: (val: number) => formatDuration(val),
    },
    axisBorder: { show: false },
    axisTicks: { show: false },
  },
  yaxis: {
    labels: { style: { colors: '#475569', fontSize: '12px', fontWeight: 500 } },
  },
  grid: {
    borderColor: '#f1f5f9',
    strokeDashArray: 4,
    xaxis: { lines: { show: true } },
    yaxis: { lines: { show: false } },
  },
  tooltip: {
    theme: 'light',
    y: { formatter: (val: number) => formatDuration(val) },
  },
}))
const categoryBarSeries = computed(() => [
  { name: 'Avg Resolution Time', data: categoryAvgTimes.value.values },
])

// 6) Staff Ranking Horizontal Bar
const staffBarOptions = computed<ApexOptions>(() => ({
  chart: {
    type: 'bar',
    toolbar: { show: false },
    fontFamily: 'inherit',
  },
  colors: ['#10b981'],
  plotOptions: {
    bar: {
      horizontal: true,
      borderRadius: 6,
      barHeight: '50%',
      distributed: true,
    },
  },
  dataLabels: {
    enabled: true,
    formatter: (val: number) => `${val}%`,
    style: { fontSize: '11px', fontWeight: 600, colors: ['#fff'] },
    offsetX: -4,
  },
  xaxis: {
    categories: staffRankingData.value.names,
    max: 100,
    labels: { show: false },
    axisBorder: { show: false },
    axisTicks: { show: false },
  },
  yaxis: {
    labels: { style: { colors: '#475569', fontSize: '12px', fontWeight: 500 } },
  },
  grid: { show: false },
  legend: { show: false },
  tooltip: {
    theme: 'light',
    y: { formatter: (val: number) => `${val}% SLA Compliance` },
  },
}))
const staffBarSeries = computed(() => [
  { name: 'SLA Rate', data: staffRankingData.value.rates },
])

// Individual mini gauges
function staffMiniGaugeOptions(rate: number): ApexOptions {
  return {
    chart: { type: 'radialBar', sparkline: { enabled: true } },
    colors: [rate >= 80 ? '#10b981' : rate >= 60 ? '#f59e0b' : '#ef4444'],
    plotOptions: {
      radialBar: {
        startAngle: -90,
        endAngle: 90,
        hollow: { size: '58%' },
        track: { background: '#f1f5f9', strokeWidth: '100%' },
        dataLabels: {
          name: { show: false },
          value: {
            show: true,
            fontSize: '15px',
            fontWeight: 700,
            color: '#1e293b',
            offsetY: -6,
            formatter: (val: number) => `${val}%`,
          },
        },
      },
    },
    stroke: { lineCap: 'round' },
  }
}

// Re-render charts when period or staff changes
watch([selectedPeriod, selectedStaffId], () => {
  nextTick(() => { chartKey.value++ })
})

// Also re-render when data changes
watch([requests, housekeepers], () => {
  nextTick(() => { chartKey.value++ })
}, { deep: true })

function getInitials(user: User): string {
  return ((user.first_name?.[0] || '') + (user.last_name?.[0] || '')).toUpperCase() || '?'
}
</script>

<template>
  <div class="hk-dashboard">
    <!-- Header with Elevated Toolbar -->
    <div class="hk-top-bar">
      <div>
        <h1 class="hk-title">Housekeeping SLA & Performance</h1>
        <p class="hk-subtitle">Team velocity, resolution benchmarks, and individual efficiency</p>
      </div>

      <div class="hk-toolbar">
        <!-- Staff Filter Dropdown -->
        <div class="hk-staff-select-wrapper">
          <Users class="staff-icon" />
          <select v-model="selectedStaffId" class="hk-staff-select">
            <option value="all">All Housekeeping Staff</option>
            <option v-for="hk in housekeepers" :key="hk.id" :value="hk.id">
              {{ hk.first_name }} {{ hk.last_name }}
            </option>
          </select>
        </div>

        <div class="hk-period-filter">
          <button
            v-for="opt in periodOptions"
            :key="opt.value"
            class="period-btn"
            :class="{ active: selectedPeriod === opt.value }"
            @click="selectedPeriod = opt.value"
          >
            {{ opt.label }}
          </button>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="hk-loading">
      <div class="hk-spinner" />
      <p>Loading performance data...</p>
    </div>

    <template v-else>
      <!-- Top 4 Metric Cards -->
      <div class="metrics-row" :key="'metrics-' + chartKey">
        <!-- SLA Compliance -->
        <div class="metric-card">
          <div class="metric-chart-wrap">
            <VueApexCharts
              type="radialBar"
              :options="slaGaugeOptions"
              :series="slaGaugeSeries"
              height="170"
            />
          </div>
          <p class="metric-label">SLA Compliance Rate</p>
        </div>

        <!-- Avg Resolution -->
        <div class="metric-card">
          <div class="metric-chart-wrap">
            <VueApexCharts
              type="radialBar"
              :options="avgGaugeOptions"
              :series="avgGaugeSeries"
              height="170"
            />
          </div>
          <p class="metric-label">Avg. Resolution Time</p>
        </div>

        <!-- Status Distribution -->
        <div class="metric-card">
          <div class="metric-chart-wrap">
            <VueApexCharts
              type="donut"
              :options="statusDonutOptions"
              :series="statusDonutSeries"
              height="170"
            />
          </div>
          <p class="metric-label">Status Breakdown</p>
        </div>

        <!-- Weekly Trend -->
        <div class="metric-card metric-card-wide">
          <p class="metric-label metric-label-top">Weekly Task Volume</p>
          <div class="metric-chart-wrap trend-wrap">
            <VueApexCharts
              type="area"
              :options="trendOptions"
              :series="trendSeries"
              height="140"
            />
          </div>
        </div>
      </div>

      <!-- Donut Legend Strip -->
      <div class="status-legend">
        <span class="legend-dot" style="--dot-color: #94a3b8">Pending</span>
        <span class="legend-dot" style="--dot-color: #6366f1">Assigned</span>
        <span class="legend-dot" style="--dot-color: #f59e0b">In Progress</span>
        <span class="legend-dot" style="--dot-color: #10b981">Completed</span>
        <span class="legend-dot" style="--dot-color: #ef4444">Cancelled</span>
      </div>

      <!-- Mid Charts Row -->
      <div class="charts-row" :key="'charts-' + chartKey">
        <!-- Category Avg Time -->
        <div class="chart-card">
          <div class="chart-card-head">
            <h2 class="chart-card-title">Resolution Time by Category</h2>
            <span class="chart-card-sub">Average completion duration</span>
          </div>
          <div class="chart-card-body">
            <VueApexCharts
              v-if="categoryAvgTimes.categories.length > 0"
              type="bar"
              :options="categoryBarOptions"
              :series="categoryBarSeries"
              height="220"
            />
            <div v-else class="chart-empty">No category data recorded</div>
          </div>
        </div>

        <!-- Staff Ranking -->
        <div class="chart-card">
          <div class="chart-card-head">
            <h2 class="chart-card-title">Staff Efficiency Ranking</h2>
            <span class="chart-card-sub">SLA Compliance Rate</span>
          </div>
          <div class="chart-card-body">
            <VueApexCharts
              v-if="staffRankingData.names.length > 0"
              type="bar"
              :options="staffBarOptions"
              :series="staffBarSeries"
              :height="Math.max(180, staffRankingData.names.length * 48)"
            />
            <div v-else class="chart-empty">No housekeeping staff registered</div>
          </div>
        </div>
      </div>

      <!-- Individual Section -->
      <div class="individuals-section">
        <h2 class="section-title">Individual Staff Performance</h2>
        <div v-if="staffPerformance.length === 0" class="chart-empty" style="padding: 40px">
          No housekeeping staff registered.
        </div>
        <div v-else class="staff-grid">
          <div
            v-for="sp in staffPerformance"
            :key="sp.user.id"
            class="staff-card"
          >
            <!-- Avatar & Name -->
            <div class="staff-identity">
              <div class="staff-avatar">{{ getInitials(sp.user) }}</div>
              <div>
                <p class="staff-name">{{ sp.user.first_name }} {{ sp.user.last_name }}</p>
                <p class="staff-meta">{{ sp.completedTasks }} completed · {{ sp.activeTasks }} active</p>
              </div>
            </div>

            <!-- Mini SLA Gauge -->
            <div class="staff-gauge-wrap">
              <VueApexCharts
                type="radialBar"
                :options="staffMiniGaugeOptions(sp.slaRate)"
                :series="[sp.slaRate]"
                height="110"
              />
              <span class="staff-gauge-label">SLA Rate</span>
            </div>

            <!-- Workload Bar -->
            <div class="staff-workload">
              <div class="workload-header">
                <span class="workload-label">Active Workload</span>
                <span class="workload-count">{{ sp.activeTasks }} active tasks</span>
              </div>
              <div class="workload-track">
                <div
                  class="workload-fill"
                  :style="{ width: Math.min(100, sp.activeTasks * 20) + '%' }"
                  :class="{
                    'fill-low': sp.activeTasks <= 2,
                    'fill-mid': sp.activeTasks > 2 && sp.activeTasks <= 4,
                    'fill-high': sp.activeTasks > 4,
                  }"
                />
              </div>
            </div>

            <!-- Avg Time Tag -->
            <div class="staff-avg-tag" v-if="sp.avgTime > 0">
              <Clock class="avg-icon" />
              <span>{{ formatDuration(sp.avgTime) }} avg</span>
            </div>
          </div>
        </div>
      </div>

      <!-- SLA Target Reference -->
      <div class="sla-reference">
        <div class="sla-ref-head">
          <Sparkles class="sla-icon" />
          <h3 class="sla-ref-title">Target SLA Benchmarks</h3>
        </div>
        <div class="sla-ref-items">
          <div class="sla-ref-item">
            <span class="sla-badge urgent">Urgent</span>
            <span class="sla-target">≤ 30 mins</span>
          </div>
          <div class="sla-ref-item">
            <span class="sla-badge high">High</span>
            <span class="sla-target">≤ 1 hour</span>
          </div>
          <div class="sla-ref-item">
            <span class="sla-badge medium">Medium</span>
            <span class="sla-target">≤ 3 hours</span>
          </div>
          <div class="sla-ref-item">
            <span class="sla-badge low">Low</span>
            <span class="sla-target">≤ 6 hours</span>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
/* ========== LAYOUT ========== */
.hk-dashboard {
  display: flex;
  flex-direction: column;
  gap: 22px;
  max-width: 1440px;
}

.hk-top-bar {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  flex-wrap: wrap;
  gap: 16px;
}

.hk-title {
  font-size: 1.6rem;
  font-weight: 800;
  color: #1a1612;
  letter-spacing: -0.015em;
  margin: 0 0 4px 0;
}

.hk-subtitle {
  color: #8c7a62;
  font-size: 0.875rem;
  margin: 0;
}

.hk-toolbar {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.hk-staff-select-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #ffffff;
  border: 1px solid rgba(250, 235, 198, 0.95);
  border-radius: 12px;
  padding: 2px 12px;
  box-shadow: 0 2px 8px rgba(180, 140, 60, 0.05);
}

.staff-icon {
  width: 16px;
  height: 16px;
  color: #8c6a22;
}

.hk-staff-select {
  padding: 8px 12px 8px 0;
  border: none;
  background: transparent;
  color: #1a1612;
  font-size: 0.8125rem;
  font-weight: 600;
  outline: none;
  cursor: pointer;
}

/* Period Filter */
.hk-period-filter {
  display: flex;
  gap: 6px;
  background: #ffffff;
  border: 1px solid rgba(250, 235, 198, 0.95);
  border-radius: 12px;
  padding: 4px;
  box-shadow: 0 2px 8px rgba(180, 140, 60, 0.05);
}

.period-btn {
  padding: 7px 16px;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: #8c7a62;
  font-size: 0.8125rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.period-btn:hover {
  background: rgba(250, 235, 198, 0.45);
  color: #4a3a20;
}

.period-btn.active {
  background: #1a1612;
  color: #faebc6;
  box-shadow: 0 2px 10px rgba(26, 22, 18, 0.2);
}

/* Loading */
.hk-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 0;
  color: #8c7a62;
  gap: 16px;
}

.hk-spinner {
  width: 36px;
  height: 36px;
  border: 3px solid #faebc6;
  border-top-color: #1a1612;
  border-radius: 50%;
  animation: hk-spin 0.8s linear infinite;
}

@keyframes hk-spin {
  to { transform: rotate(360deg); }
}

/* ========== TOP METRICS ========== */
.metrics-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 18px;
}

@media (max-width: 1100px) {
  .metrics-row {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 600px) {
  .metrics-row {
    grid-template-columns: 1fr;
  }
}

.metric-card {
  background: #ffffff;
  border: 1px solid rgba(250, 235, 198, 0.85);
  border-radius: 18px;
  padding: 20px 16px 14px;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 4px 14px rgba(180, 140, 60, 0.04);
  transition: transform 0.2s, box-shadow 0.2s;
}

.metric-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 24px rgba(180, 140, 60, 0.08);
}

.metric-card-wide {
  align-items: stretch;
}

.metric-chart-wrap {
  display: flex;
  justify-content: center;
  align-items: center;
}

.trend-wrap {
  margin: 0 -8px;
}

.metric-label {
  text-align: center;
  font-size: 0.78rem;
  font-weight: 700;
  color: #785a21;
  text-transform: uppercase;
  letter-spacing: 0.4px;
  margin-top: 4px;
}

.metric-label-top {
  margin-top: 0;
  margin-bottom: 6px;
}

/* ========== STATUS LEGEND ========== */
.status-legend {
  display: flex;
  justify-content: center;
  gap: 20px;
  flex-wrap: wrap;
  margin-top: -8px;
}

.legend-dot {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.75rem;
  font-weight: 600;
  color: #64748b;
}

.legend-dot::before {
  content: '';
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--dot-color);
  flex-shrink: 0;
}

/* ========== MID CHARTS ========== */
.charts-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

@media (max-width: 900px) {
  .charts-row {
    grid-template-columns: 1fr;
  }
}

.chart-card {
  background: #ffffff;
  border: 1px solid rgba(250, 235, 198, 0.85);
  border-radius: 18px;
  overflow: hidden;
  box-shadow: 0 4px 14px rgba(180, 140, 60, 0.04);
}

.chart-card-head {
  padding: 18px 22px 14px;
  border-bottom: 1px solid rgba(250, 235, 198, 0.6);
  background: rgba(250, 235, 198, 0.15);
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

.chart-card-title {
  font-size: 1.05rem;
  font-weight: 700;
  color: #1a1612;
}

.chart-card-sub {
  font-size: 0.75rem;
  color: #8c7a62;
  font-weight: 500;
}

.chart-card-body {
  padding: 14px 18px 18px;
}

.chart-empty {
  text-align: center;
  padding: 40px;
  color: #9c8b74;
  font-size: 0.85rem;
  font-style: italic;
}

/* ========== INDIVIDUAL SECTION ========== */
.individuals-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 800;
  color: #1a1612;
}

.staff-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 18px;
}

.staff-card {
  background: #ffffff;
  border: 1px solid rgba(250, 235, 198, 0.85);
  border-radius: 18px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  box-shadow: 0 4px 14px rgba(180, 140, 60, 0.04);
  transition: transform 0.2s, box-shadow 0.2s;
}

.staff-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 24px rgba(180, 140, 60, 0.09);
}

.staff-identity {
  display: flex;
  align-items: center;
  gap: 12px;
}

.staff-avatar {
  width: 42px;
  height: 42px;
  border-radius: 13px;
  background: linear-gradient(135deg, #dfba52, #faebc6);
  color: #3d2b07;
  font-weight: 800;
  font-size: 0.85rem;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 6px rgba(180, 140, 60, 0.15);
  flex-shrink: 0;
}

.staff-name {
  font-size: 0.925rem;
  font-weight: 700;
  color: #1a1612;
  line-height: 1.2;
}

.staff-meta {
  font-size: 0.75rem;
  color: #8c7a62;
  margin-top: 2px;
}

.staff-gauge-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin: -4px 0;
}

.staff-gauge-label {
  font-size: 0.68rem;
  font-weight: 700;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-top: -8px;
}

/* Workload bar */
.staff-workload {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.workload-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.workload-label {
  font-size: 0.72rem;
  font-weight: 600;
  color: #64748b;
}

.workload-count {
  font-size: 0.7rem;
  color: #94a3b8;
  font-weight: 600;
}

.workload-track {
  height: 6px;
  background: #f1f5f9;
  border-radius: 3px;
  overflow: hidden;
}

.workload-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.4s ease;
}

.fill-low { background: #10b981; }
.fill-mid { background: #f59e0b; }
.fill-high { background: #ef4444; }

/* Avg time tag */
.staff-avg-tag {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 0.75rem;
  color: #475569;
  font-weight: 600;
  background: #f8f6f2;
  padding: 4px 8px;
  border-radius: 8px;
  align-self: flex-start;
}

.avg-icon {
  width: 12px;
  height: 12px;
}

/* ========== SLA REFERENCE ========== */
.sla-reference {
  background: #ffffff;
  border: 1px solid rgba(250, 235, 198, 0.85);
  border-radius: 18px;
  padding: 20px 24px;
  box-shadow: 0 4px 14px rgba(180, 140, 60, 0.03);
}

.sla-ref-head {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 14px;
}

.sla-icon {
  width: 18px;
  height: 18px;
  color: #dfba52;
}

.sla-ref-title {
  font-size: 0.85rem;
  font-weight: 700;
  color: #785a21;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin: 0;
}

.sla-ref-items {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.sla-ref-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.sla-badge {
  font-size: 0.7rem;
  font-weight: 700;
  padding: 3px 10px;
  border-radius: 6px;
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.sla-badge.urgent {
  background: #fef2f2;
  color: #dc2626;
}

.sla-badge.high {
  background: #fff7ed;
  color: #ea580c;
}

.sla-badge.medium {
  background: #fffbeb;
  color: #d97706;
}

.sla-badge.low {
  background: #f0fdf4;
  color: #16a34a;
}

.sla-target {
  font-size: 0.8rem;
  color: #475569;
  font-weight: 600;
}
</style>
