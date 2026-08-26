<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { serviceRequestService } from '@/services/admin/serviceRequestService'
import { reservationService } from '@/services/admin/reservationService'
import { inventoryService } from '@/services/admin/inventoryService'
import { userService } from '@/services/admin/userService'
import { roomUnitService } from '@/services/admin/roomUnitService'
import { useToastStore } from '@/stores/toastStore'
import type { ServiceRequest, Reservation, User, RoomUnit, InventoryItem } from '@/types'
import VueApexCharts from 'vue3-apexcharts'
import type { ApexOptions } from 'apexcharts'
import {
  Users, ShieldCheck, Clock, Layers, Sparkles, CheckCircle2, AlertCircle,
  RotateCcw, BarChart3, Zap, Building2, Calendar, ChevronLeft, ChevronRight,
  Sun, CalendarDays, CalendarRange, BedDouble, Package, Target, PieChart,
  ArrowUpRight, ArrowDownRight, AlertTriangle, Hotel
} from 'lucide-vue-next'

const toastStore = useToastStore()
const loading = ref(true)
const chartKey = ref(0)

const allStaff = ref<User[]>([])
const serviceRequests = ref<ServiceRequest[]>([])
const reservations = ref<Reservation[]>([])
const roomUnits = ref<RoomUnit[]>([])
const inventoryItems = ref<InventoryItem[]>([])

type TimeframeMode = 'daily' | 'weekly' | 'monthly'
const timeframeMode = ref<TimeframeMode>('daily')

function formatDateIso(d: Date): string {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

function formatMonthIso(d: Date): string {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  return `${y}-${m}`
}

const selectedDate = ref(formatDateIso(new Date()))
const selectedMonth = ref(formatMonthIso(new Date()))
const selectedRole = ref('all')
const selectedStaffId = ref('all')

const roleTabs = [
  { label: 'All Departments', value: 'all', icon: Layers },
  { label: 'Housekeeping', value: 'housekeeping', icon: Sparkles },
  { label: 'Front Desk', value: 'reception', icon: Building2 },
  { label: 'Inventory', value: 'inventory', icon: Zap },
  { label: 'Finance', value: 'finance', icon: BarChart3 },
  { label: 'Management', value: 'manager', icon: ShieldCheck },
]

const SLA_BENCHMARKS: Record<string, number> = {
  urgent: 30,
  high: 60,
  medium: 180,
  low: 360,
  default: 120,
}

let pollingTimer: ReturnType<typeof setInterval> | null = null

const fetchAllData = async (silent = false) => {
  try {
    if (!silent) loading.value = true
    const [hkUsers, recUsers, invUsers, finUsers, admUsers, mgrUsers] = await Promise.all([
      userService.getByRole('housekeeping').catch(() => []),
      userService.getByRole('reception').catch(() => []),
      userService.getByRole('inventory').catch(() => []),
      userService.getByRole('finance').catch(() => []),
      userService.getByRole('admin').catch(() => []),
      userService.getByRole('manager').catch(() => []),
    ])
    const staffMap = new Map<string, User>()
      ;[...hkUsers, ...recUsers, ...invUsers, ...finUsers, ...admUsers, ...mgrUsers].forEach((u) => {
        if (u?.id) staffMap.set(u.id, u)
      })
    allStaff.value = Array.from(staffMap.values())

    const [srData, resData, ruData, invItemsData] = await Promise.all([
      serviceRequestService.getAll().catch(() => []),
      reservationService.getAll().catch(() => []),
      roomUnitService.getAll().catch(() => []),
      inventoryService.getItems().catch(() => []),
    ])
    serviceRequests.value = srData || []
    reservations.value = resData || []
    roomUnits.value = ruData || []
    inventoryItems.value = invItemsData || []
  } catch (error: unknown) {
    if (!silent) toastStore.error(error instanceof Error ? error.message : 'Failed to load data')
  } finally {
    if (!silent) loading.value = false
  }
}

onMounted(() => {
  fetchAllData(false)
  pollingTimer = setInterval(() => fetchAllData(true), 20000)
})
onUnmounted(() => { if (pollingTimer) clearInterval(pollingTimer) })

function shiftDaily(offset: number) {
  const d = new Date(selectedDate.value + 'T00:00:00')
  d.setDate(d.getDate() + offset)
  selectedDate.value = formatDateIso(d)
}
function shiftWeekly(offsetWeeks: number) {
  const d = new Date(selectedDate.value + 'T00:00:00')
  d.setDate(d.getDate() + offsetWeeks * 7)
  selectedDate.value = formatDateIso(d)
}
function shiftMonthly(offsetMonths: number) {
  const [yStr, mStr] = selectedMonth.value.split('-')
  const y = parseInt(yStr || '2026', 10)
  const m = parseInt(mStr || '1', 10) - 1
  selectedMonth.value = formatMonthIso(new Date(y, m + offsetMonths, 1))
}

const activeTimeframeRange = computed(() => {
  if (timeframeMode.value === 'daily') {
    const start = new Date(selectedDate.value + 'T00:00:00')
    const end = new Date(selectedDate.value + 'T23:59:59.999')
    return { start, end, label: formattedDateDisplay.value }
  } else if (timeframeMode.value === 'weekly') {
    const curr = new Date(selectedDate.value + 'T00:00:00')
    const dow = curr.getDay() || 7
    const start = new Date(curr)
    start.setDate(curr.getDate() - (dow - 1))
    start.setHours(0, 0, 0, 0)
    const end = new Date(start)
    end.setDate(start.getDate() + 6)
    end.setHours(23, 59, 59, 999)
    const s = start.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
    const e = end.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
    return { start, end, label: `${s} — ${e}` }
  } else {
    const [yStr, mStr] = selectedMonth.value.split('-')
    const y = parseInt(yStr || '2026', 10)
    const m = parseInt(mStr || '1', 10) - 1
    const start = new Date(y, m, 1)
    const end = new Date(y, m + 1, 0, 23, 59, 59, 999)
    return { start, end, label: start.toLocaleDateString('en-US', { month: 'long', year: 'numeric' }) }
  }
})

const formattedDateDisplay = computed(() => {
  const d = new Date(selectedDate.value + 'T00:00:00')
  if (isNaN(d.getTime())) return selectedDate.value
  return d.toLocaleDateString('en-US', { weekday: 'short', day: 'numeric', month: 'short', year: 'numeric' })
})

const availableStaffOptions = computed(() =>
  selectedRole.value === 'all' ? allStaff.value : allStaff.value.filter((s) => s.role === selectedRole.value)
)

watch([timeframeMode, selectedDate, selectedMonth, selectedRole, selectedStaffId], () => {
  nextTick(() => chartKey.value++)
})

function getDurationMinutes(createdAt: string, completedAt?: string | null): number {
  if (!completedAt) return 0
  return Math.max(1, Math.round((new Date(completedAt).getTime() - new Date(createdAt).getTime()) / 60000))
}

function formatDuration(mins: number): string {
  if (!mins || mins <= 0) return '0m'
  if (mins < 60) return `${mins}m`
  const h = Math.floor(mins / 60)
  const m = mins % 60
  return m > 0 ? `${h}h ${m}m` : `${h}h`
}

function getRoleLabel(role: string): string {
  const found = roleTabs.find((r) => r.value === role)
  return found ? found.label : role.charAt(0).toUpperCase() + role.slice(1)
}

function formatCurrency(val: number): string {
  if (val >= 1_000_000_000) return `Rp${(val / 1_000_000_000).toFixed(1)}B`
  if (val >= 1_000_000) return `Rp${(val / 1_000_000).toFixed(1)}M`
  if (val >= 1_000) return `Rp${(val / 1_000).toFixed(0)}K`
  return `Rp${val}`
}

function getPerformanceLabel(rate: number): string {
  if (rate >= 95) return 'Excellent'
  if (rate >= 85) return 'Good'
  if (rate >= 70) return 'Needs Attention'
  return 'Critical'
}

// -------------------------------------------------------------
// TASK RECORDS (from service requests)
// -------------------------------------------------------------
export interface TaskRecord {
  id: string
  title: string
  roomName: string
  department: string
  staffName: string
  staffId?: string
  priority: string
  createdAt: string
  completedAt?: string | null
  status: string
  durationMinutes: number
  targetMinutes: number
  elapsedMinutes: number
  slaPercent: number
  isSlaMet: boolean
  isSlaBreached: boolean
  dateObj: Date
}

const filteredDatabaseTasks = computed<TaskRecord[]>(() => {
  const { start, end } = activeTimeframeRange.value
  const result: TaskRecord[] = []
  serviceRequests.value.forEach((r) => {
    if (!r.created_at) return
    const d = new Date(r.created_at)
    if (d < start || d > end) return
    if (selectedStaffId.value !== 'all' && r.assigned_to_user_id !== selectedStaffId.value) return
    if (selectedRole.value !== 'all' && selectedRole.value !== 'housekeeping') return

    const assignedUser = allStaff.value.find((u) => u.id === r.assigned_to_user_id)
    const staffName = assignedUser ? `${assignedUser.first_name} ${assignedUser.last_name}` : 'Unassigned'
    const target = SLA_BENCHMARKS[r.priority] || 180
    const isCompleted = r.status === 'completed' && Boolean(r.completed_at)
    const durationMinutes = isCompleted ? getDurationMinutes(r.created_at, r.completed_at) : 0
    const elapsedMinutes = Math.max(1, Math.round((Date.now() - d.getTime()) / 60000))
    const slaPercent = Math.min(100, Math.round(((isCompleted ? durationMinutes : elapsedMinutes) / target) * 100))
    const isSlaMet = isCompleted ? durationMinutes <= target : elapsedMinutes <= target
    const isSlaBreached = !isCompleted && elapsedMinutes > target

    result.push({
      id: r.id,
      title: r.title || 'Room Service Task',
      roomName: r.room?.name || 'Room General',
      department: 'Housekeeping',
      staffName,
      staffId: r.assigned_to_user_id,
      priority: r.priority || 'medium',
      createdAt: r.created_at,
      completedAt: r.completed_at,
      status: r.status,
      durationMinutes,
      targetMinutes: target,
      elapsedMinutes,
      slaPercent,
      isSlaMet,
      isSlaBreached,
      dateObj: d,
    })
  })
  return result
})

// -------------------------------------------------------------
// PERIOD SUMMARY
// -------------------------------------------------------------
const periodSummary = computed(() => {
  const list = filteredDatabaseTasks.value
  const total = list.length
  const completed = list.filter((t) => t.status === 'completed').length
  const inProgress = list.filter((t) => t.status === 'in_progress' || t.status === 'assigned_to_housekeeping').length
  const pending = list.filter((t) => ['pending', 'assigned'].includes(t.status)).length
  const onTime = list.filter((t) => t.isSlaMet).length
  const breached = list.filter((t) => t.isSlaBreached).length
  const slaRate = total > 0 ? Math.round((onTime / total) * 100) : 100
  const avgSlaPercent = total > 0 ? Math.round(list.reduce((s, t) => s + t.slaPercent, 0) / total) : 0
  const totalDurations = list.filter((t) => t.status === 'completed').reduce((s, t) => s + t.durationMinutes, 0)
  const avgDuration = completed > 0 ? Math.round(totalDurations / completed) : 0
  return { total, completed, inProgress, pending, onTime, breached, slaRate, avgSlaPercent, avgDuration }
})

const completionRate = computed(() => periodSummary.value.total > 0 ? Math.round((periodSummary.value.completed / periodSummary.value.total) * 100) : 100)
const activeLoadRate = computed(() => periodSummary.value.total > 0 ? Math.round((periodSummary.value.inProgress / periodSummary.value.total) * 100) : 0)

const prioritySummary = computed(() => {
  const list = filteredDatabaseTasks.value
  return {
    urgent: list.filter((t) => t.priority === 'urgent').length,
    high: list.filter((t) => t.priority === 'high').length,
    medium: list.filter((t) => t.priority === 'medium').length,
    low: list.filter((t) => t.priority === 'low').length,
  }
})

const statusSummary = computed(() => ({
  completed: periodSummary.value.completed,
  inProgress: periodSummary.value.inProgress,
  pending: periodSummary.value.pending,
}))

const averageTarget = computed(() => {
  const completed = filteredDatabaseTasks.value.filter((t) => t.status === 'completed')
  if (!completed.length) return 0
  return Math.round(completed.reduce((s, t) => s + t.targetMinutes, 0) / completed.length)
})

const workloadMessage = computed(() => {
  if (activeLoadRate.value >= 40) return 'High active workload'
  if (activeLoadRate.value >= 20) return 'Moderate active workload'
  return 'Workload under control'
})

// -------------------------------------------------------------
// HEALTH SCORE
// -------------------------------------------------------------
const healthScore = computed(() => Math.round((periodSummary.value.slaRate + completionRate.value) / 2))
const healthColor = computed(() => {
  if (healthScore.value >= 90) return '#10b981'
  if (healthScore.value >= 75) return '#3b82f6'
  if (healthScore.value >= 60) return '#f59e0b'
  return '#ef4444'
})

// -------------------------------------------------------------
// RESERVATION INSIGHTS
// -------------------------------------------------------------
const filteredReservations = computed(() => {
  const { start, end } = activeTimeframeRange.value
  return reservations.value.filter((r) => {
    const checkin = r.actual_checkin_at || r.checkin_date
    if (!checkin) return false
    const d = new Date(checkin)
    return d >= start && d <= end
  })
})

const reservationStats = computed(() => {
  const list = filteredReservations.value
  return {
    total: list.length,
    checkedIn: list.filter((r) => r.status === 'checked_in' || r.status === 'checked_out').length,
    checkedOut: list.filter((r) => r.status === 'checked_out').length,
    pending: list.filter((r) => r.status === 'pending' || r.status === 'confirmed').length,
    totalRevenue: list.reduce((s, r) => s + (r.total_price || 0), 0),
  }
})

// -------------------------------------------------------------
// OCCUPANCY
// -------------------------------------------------------------
const roomStatusSummary = computed(() => {
  const total = roomUnits.value.length
  const occupied = roomUnits.value.filter((r) => r.status === 'occupied').length
  const available = roomUnits.value.filter((r) => r.status === 'available').length
  const dirty = roomUnits.value.filter((r) => r.status === 'dirty').length
  const maintenance = roomUnits.value.filter((r) => r.status === 'maintenance').length
  return { total, occupied, available, dirty, maintenance, occupancyRate: total > 0 ? Math.round((occupied / total) * 100) : 0 }
})

// -------------------------------------------------------------
// INVENTORY
// -------------------------------------------------------------
const inventoryAlerts = computed(() => inventoryItems.value.filter((i) => i.current_stock <= i.reorder_level).slice(0, 3))
const inventoryStats = computed(() => ({
  totalItems: inventoryItems.value.length,
  lowStock: inventoryItems.value.filter((i) => i.current_stock <= i.reorder_level).length,
}))

// -------------------------------------------------------------
// CHARTS
// -------------------------------------------------------------

// Timeframe Flow Bar Chart
const timeframeFlowData = computed(() => {
  const makeBins = (slots: string[], fn: (d: Date) => number) => {
    const onTime = new Array(slots.length).fill(0)
    const late = new Array(slots.length).fill(0)
    filteredDatabaseTasks.value.forEach((t) => {
      const idx = Math.min(slots.length - 1, Math.max(0, fn(t.dateObj)))
      if (t.isSlaMet) onTime[idx]++
      else if (t.isSlaBreached) late[idx]++
    })
    return { categories: slots, onTime, late }
  }

  if (timeframeMode.value === 'daily') {
    const slots = ['08:00', '10:00', '12:00', '14:00', '16:00', '18:00', '20:00', '22:00']
    return { ...makeBins(slots, (d) => Math.floor((d.getHours() - 8) / 2)), title: "Hourly Task Flow" }
  } else if (timeframeMode.value === 'weekly') {
    const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
    return { ...makeBins(days, (d) => (d.getDay() + 6) % 7), title: "Day-by-Day Task Flow" }
  } else {
    const weeks = ['Week 1', 'Week 2', 'Week 3', 'Week 4', 'Week 5']
    return { ...makeBins(weeks, (d) => Math.min(4, Math.floor((d.getDate() - 1) / 7))), title: "Week-by-Week Flow" }
  }
})

const flowChartOptions = computed<ApexOptions>(() => ({
  chart: { type: 'bar', toolbar: { show: false }, fontFamily: 'inherit' },
  colors: ['#10b981', '#ef4444'],
  plotOptions: { bar: { horizontal: false, columnWidth: '42%', borderRadius: 5 } },
  dataLabels: { enabled: false },
  stroke: { show: true, width: 2, colors: ['transparent'] },
  xaxis: { categories: timeframeFlowData.value.categories, labels: { style: { colors: '#64748b', fontSize: '10px', fontWeight: 600 } }, axisBorder: { show: false }, axisTicks: { show: false } },
  yaxis: { min: 0, labels: { style: { colors: '#94a3b8', fontSize: '10px' } } },
  grid: { borderColor: '#f1f5f9', strokeDashArray: 4, xaxis: { lines: { show: false } }, yaxis: { lines: { show: true } } },
  legend: { position: 'top', horizontalAlign: 'right', fontSize: '11px', fontWeight: 600, labels: { colors: '#475569' } },
  tooltip: { theme: 'light', y: { formatter: (v: number) => `${v} tasks` } }
}))
const flowSeries = computed(() => [
  { name: 'On Time', data: timeframeFlowData.value.onTime },
  { name: 'Breached', data: timeframeFlowData.value.late },
])

// Task Completion Trend Area Chart
const taskTrendData = computed(() => {
  const makeBins = (slots: string[], fn: (d: Date) => number) => {
    const counts = new Array(slots.length).fill(0)
    filteredDatabaseTasks.value.forEach((t) => {
      if (t.status === 'completed') {
        const idx = Math.min(slots.length - 1, Math.max(0, fn(t.dateObj)))
        counts[idx]++
      }
    })
    const cumulative: number[] = []
    counts.reduce((acc, val, i) => { cumulative[i] = acc + val; return acc + val }, 0)
    return { categories: slots, data: cumulative }
  }

  if (timeframeMode.value === 'daily') {
    return makeBins(['06:00', '08:00', '10:00', '12:00', '14:00', '16:00', '18:00', '20:00', '22:00'], (d) => Math.floor((d.getHours() - 6) / 2))
  } else if (timeframeMode.value === 'weekly') {
    return makeBins(['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'], (d) => (d.getDay() + 6) % 7)
  } else {
    return makeBins(['W1', 'W2', 'W3', 'W4', 'W5'], (d) => Math.min(4, Math.floor((d.getDate() - 1) / 7)))
  }
})

const trendChartOptions = computed<ApexOptions>(() => ({
  chart: { type: 'area', toolbar: { show: false }, fontFamily: 'inherit', sparkline: { enabled: false } },
  colors: ['#6366f1'],
  stroke: { curve: 'smooth', width: 2.5 },
  fill: { type: 'gradient', gradient: { shadeIntensity: 1, opacityFrom: 0.35, opacityTo: 0.05, stops: [0, 90, 100] } },
  xaxis: { categories: taskTrendData.value.categories, labels: { style: { colors: '#64748b', fontSize: '10px', fontWeight: 600 } }, axisBorder: { show: false }, axisTicks: { show: false } },
  yaxis: { min: 0, labels: { style: { colors: '#94a3b8', fontSize: '10px' } } },
  grid: { borderColor: '#f1f5f9', strokeDashArray: 4, xaxis: { lines: { show: false } }, yaxis: { lines: { show: true } } },
  tooltip: { theme: 'light', y: { formatter: (v: number) => `${v} completed` } },
  markers: { size: 4, strokeWidth: 2, strokeColors: '#fff', hover: { sizeOffset: 2 } }
}))
const trendSeries = computed(() => [{ name: 'Completed', data: taskTrendData.value.data }])
const trendDirection = computed(() => {
  const data: any = taskTrendData.value.data
  return data.length >= 2 ? data[data.length - 1] - (data[0] || 0) : 0
})

// Room Occupancy Donut
const occupancyDonutOptions = computed<ApexOptions>(() => ({
  chart: { type: 'donut', fontFamily: 'inherit' },
  colors: ['#3b82f6', '#10b981', '#f59e0b', '#6b7280'],
  labels: ['Occupied', 'Available', 'Dirty', 'Maintenance'],
  plotOptions: { pie: { donut: { size: '75%', labels: { show: true, name: { show: true, fontSize: '10px', color: '#94a3b8' }, value: { show: true, fontSize: '22px', fontWeight: 800, color: '#111827' }, total: { show: true, label: 'Rooms', fontSize: '10px', color: '#94a3b8' } } } } },
  dataLabels: { enabled: false },
  legend: { position: 'bottom', fontSize: '10px', fontWeight: 600, labels: { colors: '#64748b', useSeriesColors: false } },
  stroke: { width: 2, colors: ['#fff'] },
  tooltip: { theme: 'light' }
}))
const occupancyDonutSeries = computed(() => [
  roomStatusSummary.value.occupied, roomStatusSummary.value.available,
  roomStatusSummary.value.dirty, roomStatusSummary.value.maintenance,
])

// Service Category Donut (from task titles)
const categoryBreakdown = computed(() => {
  const c = { Cleaning: 0, Maintenance: 0, Amenities: 0, Incident: 0, Other: 0 }
  filteredDatabaseTasks.value.forEach((t) => {
    const lower = t.title.toLowerCase()
    if (lower.includes('clean')) c.Cleaning++
    else if (lower.includes('repair') || lower.includes('fix') || lower.includes('maintenance')) c.Maintenance++
    else if (lower.includes('amenit') || lower.includes('towel') || lower.includes('pillow')) c.Amenities++
    else if (lower.includes('broken') || lower.includes('damage') || lower.includes('incident')) c.Incident++
    else c.Other++
  })
  return c
})
const categoryDonutOptions = computed<ApexOptions>(() => ({
  chart: { type: 'donut', fontFamily: 'inherit' },
  colors: ['#3b82f6', '#f59e0b', '#10b981', '#ef4444', '#8b5cf6'],
  labels: Object.keys(categoryBreakdown.value),
  plotOptions: { pie: { donut: { size: '72%', labels: { show: true, name: { show: true, fontSize: '11px', color: '#94a3b8' }, value: { show: true, fontSize: '18px', fontWeight: 800, color: '#111827' }, total: { show: true, label: 'Total', fontSize: '10px', color: '#94a3b8' } } } } },
  dataLabels: { enabled: false },
  legend: { position: 'bottom', fontSize: '10px', fontWeight: 600, labels: { colors: '#64748b', useSeriesColors: false } },
  stroke: { width: 2, colors: ['#fff'] },
  tooltip: { theme: 'light', y: { formatter: (v: number) => `${v} tasks` } }
}))
const categoryDonutSeries = computed(() => Object.values(categoryBreakdown.value))

// Department Benchmark
const departmentScores = computed(() => {
  const deptMap: Record<string, { total: number; onTime: number }> = {}

  staffPerformanceList.value.forEach((sp: any) => {
    if (!deptMap[sp.role]) {
      deptMap[sp.role] = { total: 0, onTime: 0 }
    }

    // Simpan ke variabel lokal (TypeScript paham ini tidak undefined)
    const item: any = deptMap[sp.role]

    item.total += Number(sp.tasksTotal || 0)
    item.onTime += Number(sp.tasksDone || 0) * (Number(sp.slaRate || 0) / 100)
  })

  const labels: string[] = []
  const scores: number[] = []

  for (const [role, d] of Object.entries(deptMap)) {
    labels.push(getRoleLabel(role))
    scores.push(d.total > 0 ? Math.round((d.onTime / d.total) * 100) : 0)
  }

  return { categories: labels, data: scores }
})

const departmentChartOptions = computed<ApexOptions>(() => ({
  chart: { type: 'bar', toolbar: { show: false }, fontFamily: 'inherit' },
  colors: ['#10b981', '#3b82f6', '#eab308', '#8b5cf6', '#f43f5e', '#14b8a6'],
  plotOptions: { bar: { horizontal: true, borderRadius: 5, barHeight: '48%', distributed: true } },
  dataLabels: { enabled: true, formatter: (v: number) => `${v}%`, style: { fontSize: '10.5px', fontWeight: 700, colors: ['#fff'] }, offsetX: -6 },
  xaxis: { categories: departmentScores.value.categories, max: 100, labels: { show: false }, axisBorder: { show: false }, axisTicks: { show: false } },
  yaxis: { labels: { style: { colors: '#475569', fontSize: '11px', fontWeight: 600 } } },
  grid: { show: false },
  legend: { show: false },
  tooltip: { theme: 'light', y: { formatter: (v: number) => `${v}% SLA Met` } }
}))

// -------------------------------------------------------------
// STAFF PERFORMANCE
// -------------------------------------------------------------
interface StaffPerformanceItem {
  user: User; role: string; tasksDone: number; tasksTotal: number
  activeTasks: number; slaRate: number; avgMins: number; statusText: string
}

const staffPerformanceList = computed<StaffPerformanceItem[]>(() => {
  const pool = selectedStaffId.value === 'all'
    ? availableStaffOptions.value
    : availableStaffOptions.value.filter((s) => s.id === selectedStaffId.value)

  return pool.map((s) => {
    const userTasks = filteredDatabaseTasks.value.filter((t) => t.staffId === s.id)
    const tasksDone = userTasks.filter((t) => t.status === 'completed').length
    const onTime = userTasks.filter((t) => t.isSlaMet).length
    const activeTasks = userTasks.filter((t) => t.status === 'in_progress').length
    const tasksTotal = userTasks.length
    const slaRate = tasksTotal > 0 ? Math.round((onTime / tasksTotal) * 100) : 100
    const avgMins = tasksDone > 0 ? Math.round(userTasks.filter((t) => t.status === 'completed').reduce((s, t) => s + t.durationMinutes, 0) / tasksDone) : 0
    return { user: s, role: s.role, tasksDone, tasksTotal, activeTasks, slaRate, avgMins, statusText: activeTasks > 0 ? `${activeTasks} Active` : 'Available' }
  }).sort((a, b) => b.tasksDone - a.tasksDone || b.slaRate - a.slaRate)
})

const topPerformers = computed(() => staffPerformanceList.value.slice(0, 5))
const latestTasks = computed(() => [...filteredDatabaseTasks.value].sort((a, b) => b.dateObj.getTime() - a.dateObj.getTime()).slice(0, 12))

function getInitials(user: User): string {
  return ((user.first_name?.[0] || '') + (user.last_name?.[0] || '')).toUpperCase() || 'ST'
}
</script>

<template>
  <main class="dashboard">
    <!-- HEADER -->
    <section class="hero">
      <div class="hero-top">
        <div class="hero-copy">
          <div class="live-badge"><span class="live-dot"></span> LIVE OPERATIONAL CONTROL CENTER</div>
          <h1>{{ timeframeMode === 'daily' ? 'Daily' : timeframeMode === 'weekly' ? 'Weekly' : 'Monthly' }} Operations
            Overview</h1>
          <p>Monitor SLA compliance, task throughput, occupancy, reservations, inventory, and employee performance from
            one screen.</p>
        </div>
        <div class="hero-status">
          <span class="chip chip-green"><span class="live-dot"></span> Live</span>
          <span class="chip chip-muted">Auto refresh · 20s</span>
        </div>
      </div>

      <div class="controls">
        <div class="control-block">
          <span class="label">TIMEFRAME</span>
          <div class="segments">
            <button :class="{ active: timeframeMode === 'daily' }" @click="timeframeMode = 'daily'">
              <Sun /> Daily
            </button>
            <button :class="{ active: timeframeMode === 'weekly' }" @click="timeframeMode = 'weekly'">
              <CalendarDays /> Weekly
            </button>
            <button :class="{ active: timeframeMode === 'monthly' }" @click="timeframeMode = 'monthly'">
              <CalendarRange /> Monthly
            </button>
          </div>
        </div>
        <div class="control-block">
          <span class="label">PERIOD</span>
          <div class="period-row">
            <button class="nav-btn"
              @click="timeframeMode === 'daily' ? shiftDaily(-1) : timeframeMode === 'weekly' ? shiftWeekly(-1) : shiftMonthly(-1)">
              <ChevronLeft />
            </button>
            <div class="period-display">
              <Calendar />
              <strong>{{ activeTimeframeRange.label }}</strong>
              <input v-if="timeframeMode === 'daily'" v-model="selectedDate" type="date" />
              <input v-else v-model="selectedMonth" type="month" />
            </div>
            <button class="nav-btn"
              @click="timeframeMode === 'daily' ? shiftDaily(1) : timeframeMode === 'weekly' ? shiftWeekly(1) : shiftMonthly(1)">
              <ChevronRight />
            </button>
          </div>
        </div>
        <div class="control-block">
          <span class="label">EMPLOYEE</span>
          <div class="select-box">
            <Users />
            <select v-model="selectedStaffId">
              <option value="all">All Employees ({{ availableStaffOptions.length }})</option>
              <option v-for="s in availableStaffOptions" :key="s.id" :value="s.id">{{ s.first_name }} {{ s.last_name }}
              </option>
            </select>
          </div>
        </div>
      </div>

      <div class="dept-row">
        <span class="label">DEPARTMENT</span>
        <div class="dept-pills">
          <button v-for="tab in roleTabs" :key="tab.value" :class="{ active: selectedRole === tab.value }"
            @click="selectedRole = tab.value">
            <component :is="tab.icon" /> {{ tab.label }}
          </button>
        </div>
        <button v-if="selectedRole !== 'all' || selectedStaffId !== 'all'" class="reset-btn"
          @click="selectedRole = 'all'; selectedStaffId = 'all'">
          <RotateCcw /> Reset
        </button>
      </div>
    </section>

    <!-- LOADING -->
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <strong>Loading operational data</strong>
      <span>Synchronizing live database records...</span>
    </div>

    <template v-else>
      <!-- KPI ROW 1: Core SLA Metrics -->
      <section class="kpi-row">
        <article class="kpi kpi-green">
          <div class="kpi-head"><span>SLA COMPLIANCE</span>
            <ShieldCheck />
          </div>
          <div class="kpi-big"><strong>{{ periodSummary.slaRate }}%</strong><span>{{
            getPerformanceLabel(periodSummary.slaRate) }}</span></div>
          <div class="bar"><i :style="{ width: periodSummary.slaRate + '%' }"></i></div>
          <div class="kpi-foot"><span>{{ periodSummary.onTime }} on time</span><span>{{ periodSummary.total -
            periodSummary.onTime }} need attention</span></div>
        </article>
        <article class="kpi kpi-blue">
          <div class="kpi-head"><span>COMPLETION RATE</span>
            <CheckCircle2 />
          </div>
          <div class="kpi-big"><strong>{{ completionRate }}%</strong><span>{{ periodSummary.completed }} / {{
            periodSummary.total }}</span></div>
          <div class="bar bar-blue"><i :style="{ width: completionRate + '%' }"></i></div>
          <div class="kpi-foot"><span>{{ periodSummary.inProgress }} in progress</span><span>{{ periodSummary.pending }}
              pending</span></div>
        </article>
        <article class="kpi kpi-amber">
          <div class="kpi-head"><span>AVG. TURNAROUND</span>
            <Clock />
          </div>
          <div class="kpi-big"><strong>{{ formatDuration(periodSummary.avgDuration) }}</strong><span>per completed
              task</span></div>
          <div class="benchmark"><span>Target benchmark</span><b>{{ formatDuration(averageTarget) }}</b></div>
          <div class="kpi-foot"><span>Avg SLA used: {{ periodSummary.avgSlaPercent }}%</span><span>{{
            periodSummary.completed }} completed</span></div>
        </article>
        <article class="kpi kpi-red">
          <div class="kpi-head"><span>SLA BREACHED</span>
            <AlertCircle />
          </div>
          <div class="kpi-big"><strong>{{ periodSummary.breached }}</strong><span>tasks over deadline</span></div>
          <div class="bar bar-red"><span
              :style="{ width: (periodSummary.total > 0 ? Math.round(periodSummary.breached / periodSummary.total * 100) : 0) + '%' }"></span>
          </div>
          <div class="kpi-foot"><span>{{ workloadMessage }}</span><span>{{ activeLoadRate }}% active</span></div>
        </article>
      </section>

      <!-- KPI ROW 2: Business Metrics -->
      <section class="kpi-row">
        <article class="kpi kpi-teal">
          <div class="kpi-head"><span>RESERVATIONS</span>
            <Hotel />
          </div>
          <div class="kpi-big"><strong>{{ reservationStats.total }}</strong><span>this period</span></div>
          <div class="mini-row">
            <span class="mini"><i class="dot-blue"></i>{{ reservationStats.checkedIn }} checked in</span>
            <span class="mini"><i class="dot-green"></i>{{ reservationStats.checkedOut }} checked out</span>
            <span class="mini"><i class="dot-amber"></i>{{ reservationStats.pending }} pending</span>
          </div>
          <div class="kpi-foot"><span>Revenue</span><b class="revenue">{{ formatCurrency(reservationStats.totalRevenue)
          }}</b></div>
        </article>
        <article class="kpi kpi-indigo">
          <div class="kpi-head"><span>OCCUPANCY</span>
            <BedDouble />
          </div>
          <div class="kpi-big"><strong>{{ roomStatusSummary.occupancyRate }}%</strong><span>{{
            roomStatusSummary.occupied }}/{{ roomStatusSummary.total }} rooms</span></div>
          <div class="bar bar-indigo"><i :style="{ width: roomStatusSummary.occupancyRate + '%' }"></i></div>
          <div class="kpi-foot"><span>{{ roomStatusSummary.available }} available</span><span>{{ roomStatusSummary.dirty
            + roomStatusSummary.maintenance }} attention</span></div>
        </article>
        <article class="kpi kpi-yellow">
          <div class="kpi-head"><span>INVENTORY ALERTS</span>
            <Package />
          </div>
          <div class="kpi-big"><strong>{{ inventoryStats.lowStock }}</strong><span>low stock items</span></div>
          <div v-if="inventoryAlerts.length" class="alert-stack">
            <div v-for="item in inventoryAlerts" :key="item.id" class="alert-row">
              <AlertTriangle /><span><strong>{{ item.name }}</strong> — {{ item.current_stock }} {{ item.unit }}
                left</span>
            </div>
          </div>
          <div v-else class="kpi-foot"><span>All items stocked</span><span>{{ inventoryStats.totalItems }}
              tracked</span></div>
        </article>
        <article class="kpi kpi-rose">
          <div class="kpi-head"><span>STAFF ON DUTY</span>
            <Users />
          </div>
          <div class="kpi-big"><strong>{{ allStaff.length }}</strong><span>total employees</span></div>
          <div class="mini-row">
            <span class="mini"><i class="dot-amber"></i>{{staffPerformanceList.filter(s => s.activeTasks > 0).length}}
              active</span>
            <span class="mini"><i class="dot-green"></i>{{staffPerformanceList.filter(s => s.activeTasks === 0).length
            }} available</span>
          </div>
          <div class="kpi-foot"><span>{{ periodSummary.inProgress }} tasks assigned</span><span>All departments</span>
          </div>
        </article>
      </section>

      <!-- SNAPSHOT -->
      <section class="snapshot-row">
        <article class="panel">
          <div class="panel-head">
            <div><span class="eyebrow">EXECUTIVE SNAPSHOT</span>
              <h2>Operational Health</h2>
            </div>
            <BarChart3 />
          </div>
          <div class="health">
            <div class="ring"
              :style="{ background: `conic-gradient(${healthColor} 0 ${healthScore}%, #e2e8f0 ${healthScore}%)` }">
              <strong>{{ healthScore }}</strong><span>Health</span>
            </div>
            <div class="health-text">
              <strong>{{ getPerformanceLabel(healthScore) }} operational condition</strong>
              <p>Based on SLA compliance and task completion rate.</p>
            </div>
          </div>
          <div class="health-bars">
            <div><span>SLA</span><b>{{ periodSummary.slaRate }}%</b></div>
            <div><span>Completion</span><b>{{ completionRate }}%</b></div>
            <div><span>Occupancy</span><b>{{ roomStatusSummary.occupancyRate }}%</b></div>
          </div>
        </article>

        <article class="panel">
          <div class="panel-head">
            <div><span class="eyebrow">TASK PIPELINE</span>
              <h2>Current Workload</h2>
            </div>
            <Layers />
          </div>
          <template
            v-for="(item, key) in { Completed: statusSummary.completed, 'In Progress': statusSummary.inProgress, Pending: statusSummary.pending }"
            :key="key">
            <div class="pipe-row">
              <div class="pipe-label"><span :class="'dot-' + String(key).toLowerCase().replace(' ', '')"></span>{{ key
              }}</div><strong>{{ item }}</strong>
            </div>
            <div class="pipe-track"><i :class="'fill-' + String(key).toLowerCase().replace(' ', '')"
                :style="{ width: periodSummary.total ? (item / periodSummary.total * 100) + '%' : '0%' }"></i></div>
          </template>
        </article>

        <article class="panel">
          <div class="panel-head">
            <div><span class="eyebrow">RISK DISTRIBUTION</span>
              <h2>Priority Mix</h2>
            </div>
            <AlertCircle />
          </div>
          <div class="prio-grid">
            <div class="prio urgent"><strong>{{ prioritySummary.urgent }}</strong><span>Urgent</span></div>
            <div class="prio high"><strong>{{ prioritySummary.high }}</strong><span>High</span></div>
            <div class="prio medium"><strong>{{ prioritySummary.medium }}</strong><span>Medium</span></div>
            <div class="prio low"><strong>{{ prioritySummary.low }}</strong><span>Low</span></div>
          </div>
          <p class="panel-note">Counts from tasks in the active timeframe.</p>
        </article>
      </section>

      <!-- CHARTS ROW 1 -->
      <section class="charts-row charts-3">
        <article class="panel chart-panel">
          <div class="panel-head">
            <div><span class="eyebrow">TREND ANALYSIS</span>
              <h2>{{ timeframeFlowData.title }}</h2>
              <p>On-time vs breached tasks.</p>
            </div>
            <div class="legend"><span><i class="green"></i> On Time</span><span><i class="red"></i> Breached</span>
            </div>
          </div>
          <VueApexCharts :key="'flow-' + chartKey" type="bar" :options="flowChartOptions" :series="flowSeries"
            height="320" />
        </article>
        <article class="panel chart-panel">
          <div class="panel-head">
            <div><span class="eyebrow">CATEGORIES</span>
              <h2>Request Types</h2>
              <p>Service request breakdown.</p>
            </div>
            <PieChart />
          </div>
          <VueApexCharts :key="'cat-' + chartKey" type="donut" :options="categoryDonutOptions"
            :series="categoryDonutSeries" height="320" />
        </article>
        <article class="panel chart-panel">
          <div class="panel-head">
            <div><span class="eyebrow">OCCUPANCY</span>
              <h2>Room Status</h2>
              <p>Current distribution.</p>
            </div>
            <BedDouble />
          </div>
          <VueApexCharts :key="'occ-' + chartKey" type="donut" :options="occupancyDonutOptions"
            :series="occupancyDonutSeries" height="320" />
        </article>
      </section>

      <!-- CHARTS ROW 2 -->
      <section class="charts-row">
        <article class="panel chart-panel wide">
          <div class="panel-head">
            <div><span class="eyebrow">COMPLETION TREND</span>
              <h2>Cumulative Task Completion</h2>
              <p>Total completed tasks over time.</p>
            </div>
            <div class="trend-chip" :class="trendDirection >= 0 ? 'up' : 'down'">
              <component :is="trendDirection >= 0 ? ArrowUpRight : ArrowDownRight" />
              {{ Math.abs(trendDirection) }} tasks
            </div>
          </div>
          <VueApexCharts :key="'trend-' + chartKey" type="area" :options="trendChartOptions" :series="trendSeries"
            height="300" />
        </article>
        <article class="panel chart-panel">
          <div class="panel-head">
            <div><span class="eyebrow">DEPARTMENT BENCHMARK</span>
              <h2>SLA Scorecard</h2>
              <p>Comparison across departments.</p>
            </div>
          </div>
          <VueApexCharts :key="'dept-' + chartKey" type="bar" :options="departmentChartOptions"
            :series="[{ name: 'SLA Score', data: departmentScores.data }]" height="300" />
        </article>
      </section>

      <!-- PEOPLE TABLE -->
      <section class="panel table-panel">
        <div class="table-header">
          <div><span class="eyebrow">PEOPLE ANALYTICS</span>
            <h2>Employee Performance</h2>
            <p>Ranked by completed workload, then SLA compliance.</p>
          </div>
          <span class="badge">{{ staffPerformanceList.length }} Employees</span>
        </div>
        <div v-if="staffPerformanceList.length === 0" class="empty">
          <Users /><strong>No employee data</strong><span>Change the department or employee filter.</span>
        </div>
        <div v-else class="table-scroll">
          <table class="data-table">
            <thead>
              <tr>
                <th>#</th>
                <th>EMPLOYEE</th>
                <th>DEPARTMENT</th>
                <th>COMPLETED</th>
                <th>SLA</th>
                <th>AVG. SPEED</th>
                <th>ACTIVE</th>
                <th>RATING</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, i) in staffPerformanceList" :key="item.user.id">
                <td><span class="rank" :class="{ gold: i < 3 }">{{ i + 1 }}</span></td>
                <td>
                  <div class="emp"><span class="avatar" :class="'av-' + (i % 5)">{{ getInitials(item.user)
                  }}</span><span><strong>{{ item.user.first_name }} {{ item.user.last_name }}</strong><small>{{
                        item.tasksTotal }} task(s)</small></span></div>
                </td>
                <td><span class="dept-tag" :class="'dept-' + item.role">{{ getRoleLabel(item.role) }}</span></td>
                <td><strong>{{ item.tasksDone }}</strong> <small>/ {{ item.tasksTotal }}</small></td>
                <td>
                  <div class="sla-cell"><strong
                      :class="item.slaRate >= 95 ? 'c-green' : item.slaRate >= 85 ? 'c-blue' : 'c-red'">{{ item.slaRate
                      }}%</strong>
                    <div class="sla-bar"><i
                        :class="item.slaRate >= 95 ? 'b-green' : item.slaRate >= 85 ? 'b-blue' : 'b-red'"
                        :style="{ width: item.slaRate + '%' }"></i></div>
                  </div>
                </td>
                <td><strong>{{ formatDuration(item.avgMins) }}</strong></td>
                <td><span class="status-pill" :class="item.activeTasks ? 'busy' : 'free'"><i></i>{{ item.statusText
                }}</span>
                </td>
                <td><span class="rating"
                    :class="item.slaRate >= 95 ? 'excellent' : item.slaRate >= 85 ? 'good' : 'attention'">{{
                      getPerformanceLabel(item.slaRate) }}</span></td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <!-- BOTTOM ROW -->
      <section class="bottom-row">
        <article class="panel leaders-panel">
          <div class="panel-head">
            <div><span class="eyebrow">LEADERS</span>
              <h2>Top Performers</h2>
            </div>
            <Target />
          </div>
          <div v-if="topPerformers.length" class="leader-list">
            <div v-for="(item, i) in topPerformers" :key="item.user.id" class="leader" :class="{ highlight: i === 0 }">
              <span class="leader-rank" :class="{ gold: i === 0, silver: i === 1, bronze: i === 2 }">{{ i + 1 }}</span>
              <span class="avatar sm" :class="'av-' + (i % 5)">{{ getInitials(item.user) }}</span>
              <div class="leader-info"><strong>{{ item.user.first_name }} {{ item.user.last_name }}</strong><small>{{
                getRoleLabel(item.role) }} · {{ item.tasksDone }} tasks</small></div>
              <div class="leader-score"><strong>{{ item.slaRate }}%</strong><small>SLA</small></div>
            </div>
          </div>
          <div v-else class="empty-inline">No performance records available.</div>
        </article>

        <article class="panel tasks-panel">
          <div class="panel-head">
            <div><span class="eyebrow">RECENT ACTIVITY</span>
              <h2>Operational Task Log</h2>
              <p>Latest 12 records.</p>
            </div>
            <span class="badge">{{ filteredDatabaseTasks.length }} total</span>
          </div>
          <div v-if="latestTasks.length === 0" class="empty compact">
            <Layers /><strong>No tasks recorded</strong>
          </div>
          <div v-else class="task-list">
            <div v-for="task in latestTasks" :key="task.id" class="task">
              <div class="task-time">{{ task.dateObj.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })
              }}
              </div>
              <div class="task-room">{{ task.roomName }}</div>
              <div class="task-info"><strong>{{ task.title }}</strong><small>{{ task.staffName }} · {{ task.department
              }}</small></div>
              <span class="prio-tag" :class="task.priority">{{ task.priority }}</span>
              <div class="task-sla">
                <div class="sla-bar"><i :class="task.isSlaBreached ? 'b-red' : task.isSlaMet ? 'b-green' : 'b-amber'"
                    :style="{ width: Math.min(task.slaPercent, 100) + '%' }"></i></div>
                <span class="sla-text" :class="{ 'c-red': task.isSlaBreached }">{{ task.slaPercent }}% · {{ task.status
                  ===
                  'completed' ? formatDuration(task.durationMinutes) : formatDuration(task.elapsedMinutes) }}/{{
                    formatDuration(task.targetMinutes) }}</span>
              </div>
              <span class="sla-tag" :class="task.isSlaMet ? 'met' : task.isSlaBreached ? 'breached' : 'active'">{{
                task.isSlaMet ? 'Met' : task.isSlaBreached ? 'Breached' : 'In Progress' }}</span>
            </div>
          </div>
        </article>
      </section>

      <!-- FOOTER -->
      <section class="footer-bar">
        <div><span>Period</span><strong>{{ activeTimeframeRange.label }}</strong></div>
        <div><span>Total Tasks</span><strong>{{ periodSummary.total }}</strong></div>
        <div><span>Completed</span><strong>{{ periodSummary.completed }}</strong></div>
        <div><span>SLA Compliance</span><strong :class="{ 'c-red': periodSummary.slaRate < 80 }">{{
          periodSummary.slaRate
            }}%</strong></div>
        <div><span>Avg SLA Used</span><strong>{{ periodSummary.avgSlaPercent }}%</strong></div>
        <div><span>Breached</span><strong :class="{ 'c-red': periodSummary.breached > 0 }">{{ periodSummary.breached
        }}</strong></div>
        <div><span>Occupancy</span><strong>{{ roomStatusSummary.occupancyRate }}%</strong></div>
      </section>
    </template>
  </main>
</template>

<style scoped>
.dashboard {
  max-width: 1500px;
  margin: 0 auto;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  color: #172033;
  background: #f6f7f9;
  min-height: 100%;
  font-family: inherit
}

/* Shared */
.panel,
.kpi,
.hero,
.loading,
.footer-bar {
  background: #fff;
  border: 1px solid #e5e8ee;
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(15, 23, 42, .03)
}

.label,
.eyebrow {
  font-size: 9px;
  font-weight: 800;
  letter-spacing: .12em;
  color: #94a3b8;
  text-transform: uppercase
}

/* Hero */
.hero {
  overflow: hidden
}

.hero-top {
  padding: 24px 24px 16px;
  display: flex;
  justify-content: space-between;
  gap: 20px
}

.hero-copy {
  display: flex;
  flex-direction: column;
  gap: 6px
}

.live-badge {
  display: flex;
  align-items: center;
  gap: 7px;
  font-size: 9px;
  font-weight: 800;
  letter-spacing: .1em;
  color: #059669
}

.live-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #10b981;
  box-shadow: 0 0 0 4px #d1fae5;
  animation: pulse 2s ease-in-out infinite
}

@keyframes pulse {

  0%,
  100% {
    box-shadow: 0 0 0 4px #d1fae5
  }

  50% {
    box-shadow: 0 0 0 7px transparent
  }
}

.hero-copy h1 {
  margin: 0;
  font-size: 26px;
  line-height: 1.15;
  letter-spacing: -.04em;
  color: #111827
}

.hero-copy p {
  margin: 0;
  color: #94a3b8;
  font-size: 12px;
  max-width: 700px
}

.hero-status {
  display: flex;
  align-items: flex-start;
  gap: 6px
}

.chip {
  padding: 6px 10px;
  border-radius: 8px;
  font-size: 10px;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 5px
}

.chip-green {
  background: #ecfdf5;
  color: #047857
}

.chip-muted {
  background: #f8fafc;
  border: 1px solid #e8ebf0;
  color: #94a3b8
}

/* Controls */
.controls {
  padding: 12px 24px;
  border-top: 1px solid #eef1f4;
  border-bottom: 1px solid #eef1f4;
  display: grid;
  grid-template-columns: 1fr 1.2fr 1fr;
  gap: 14px;
  background: #fafbfc
}

.control-block {
  display: flex;
  flex-direction: column;
  gap: 6px
}

.segments {
  display: flex;
  padding: 3px;
  background: #f0f2f5;
  border: 1px solid #e4e7ec;
  border-radius: 10px
}

.segments button {
  flex: 1;
  border: 0;
  background: transparent;
  padding: 7px 8px;
  border-radius: 7px;
  color: #64748b;
  font-size: 10px;
  font-weight: 800;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  transition: all .2s
}

.segments button svg {
  width: 13px
}

.segments button.active {
  background: #111827;
  color: #fff;
  box-shadow: 0 2px 6px rgba(0, 0, 0, .12)
}

.period-row {
  display: flex;
  align-items: center;
  gap: 4px
}

.nav-btn {
  width: 32px;
  height: 32px;
  border: 1px solid #e1e5ea;
  background: #fff;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #64748b;
  cursor: pointer;
  transition: all .15s
}

.nav-btn:hover {
  background: #f8fafc;
  border-color: #cbd5e1
}

.nav-btn svg {
  width: 14px
}

.period-display {
  position: relative;
  flex: 1;
  height: 32px;
  border: 1px solid #e1e5ea;
  background: #fff;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0 10px
}

.period-display svg {
  width: 13px;
  color: #64748b
}

.period-display strong {
  font-size: 10px;
  color: #334155;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis
}

.period-display input {
  position: absolute;
  inset: 0;
  opacity: 0;
  cursor: pointer
}

.select-box {
  height: 32px;
  border: 1px solid #e1e5ea;
  background: #fff;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0 10px
}

.select-box svg {
  width: 13px;
  color: #64748b
}

.select-box select {
  border: 0;
  outline: 0;
  background: transparent;
  width: 100%;
  font-size: 10px;
  font-weight: 700;
  color: #334155
}

/* Department Pills */
.dept-row {
  padding: 12px 24px;
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap
}

.dept-pills {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
  flex: 1
}

.dept-pills button {
  display: flex;
  align-items: center;
  gap: 4px;
  border: 1px solid #e4e7ec;
  background: #fff;
  color: #64748b;
  border-radius: 8px;
  padding: 6px 10px;
  font-size: 10px;
  font-weight: 700;
  cursor: pointer;
  transition: all .2s
}

.dept-pills button svg {
  width: 12px
}

.dept-pills button:hover {
  border-color: #c7d2fe;
  color: #4338ca
}

.dept-pills button.active {
  background: #eef2ff;
  color: #4338ca;
  border-color: #c7d2fe
}

.reset-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  border: 1px solid #fecaca;
  background: #fff1f2;
  color: #dc2626;
  border-radius: 8px;
  padding: 6px 10px;
  font-size: 10px;
  font-weight: 800;
  cursor: pointer;
  transition: all .15s
}

.reset-btn:hover {
  background: #fee2e2
}

.reset-btn svg {
  width: 11px
}

/* Loading */
.loading {
  min-height: 300px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  color: #94a3b8
}

.loading strong {
  font-size: 13px;
  color: #334155
}

.loading span {
  font-size: 10px
}

.spinner {
  width: 28px;
  height: 28px;
  border: 3px solid #e2e8f0;
  border-top-color: #111827;
  border-radius: 50%;
  animation: spin .8s linear infinite
}

@keyframes spin {
  to {
    transform: rotate(360deg)
  }
}

/* KPI Cards */
.kpi-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px
}

.kpi {
  padding: 16px;
  min-height: 160px;
  display: flex;
  flex-direction: column;
  position: relative;
  overflow: hidden;
  transition: transform .2s, box-shadow .2s
}

.kpi:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(15, 23, 42, .05)
}

.kpi::before {
  content: "";
  position: absolute;
  left: 0;
  top: 0;
  width: 3px;
  height: 100%
}

.kpi-green::before {
  background: #10b981
}

.kpi-blue::before {
  background: #3b82f6
}

.kpi-amber::before {
  background: #eab308
}

.kpi-red::before {
  background: #ef4444
}

.kpi-teal::before {
  background: #14b8a6
}

.kpi-indigo::before {
  background: #6366f1
}

.kpi-yellow::before {
  background: #f59e0b
}

.kpi-rose::before {
  background: #f43f5e
}

.kpi-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 9px;
  font-weight: 800;
  letter-spacing: .08em;
  color: #94a3b8
}

.kpi-head svg {
  width: 15px;
  height: 15px;
  color: #cbd5e1
}

.kpi-green .kpi-head svg {
  color: #10b981
}

.kpi-blue .kpi-head svg {
  color: #3b82f6
}

.kpi-amber .kpi-head svg {
  color: #eab308
}

.kpi-red .kpi-head svg {
  color: #ef4444
}

.kpi-teal .kpi-head svg {
  color: #14b8a6
}

.kpi-indigo .kpi-head svg {
  color: #6366f1
}

.kpi-yellow .kpi-head svg {
  color: #f59e0b
}

.kpi-rose .kpi-head svg {
  color: #f43f5e
}

.kpi-big {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-top: 12px
}

.kpi-big strong {
  font-size: 28px;
  line-height: 1;
  font-weight: 900;
  letter-spacing: -.04em;
  color: #111827
}

.kpi-big span {
  font-size: 9px;
  font-weight: 700;
  color: #94a3b8
}

.bar,
.sla-bar {
  height: 5px;
  background: #f1f5f9;
  border-radius: 99px;
  overflow: hidden;
  margin-top: auto
}

.bar i,
.sla-bar i {
  display: block;
  height: 100%;
  border-radius: 99px;
  transition: width .5s ease
}

.bar i {
  background: #10b981
}

.bar-blue i {
  background: #3b82f6
}

.bar-red span,
.bar-red i {
  background: #ef4444
}

.bar-indigo i {
  background: #6366f1
}

.b-green {
  background: #10b981
}

.b-blue {
  background: #3b82f6
}

.b-amber {
  background: #f59e0b
}

.b-red {
  background: #ef4444
}

.kpi-foot {
  display: flex;
  justify-content: space-between;
  margin-top: 8px;
  font-size: 9px;
  color: #94a3b8
}

.benchmark {
  margin-top: auto;
  padding: 6px 8px;
  border-radius: 8px;
  background: #fffbeb;
  display: flex;
  justify-content: space-between;
  font-size: 9px;
  color: #a16207
}

.benchmark b {
  color: #92400e
}

.revenue {
  color: #047857 !important
}

.mini-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: auto;
  padding-top: 6px
}

.mini {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 9px;
  color: #64748b
}

.dot-blue,
.dot-green,
.dot-amber {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  display: inline-block
}

.dot-blue {
  background: #3b82f6
}

.dot-green {
  background: #10b981
}

.dot-amber {
  background: #f59e0b
}

.alert-stack {
  display: flex;
  flex-direction: column;
  gap: 5px;
  margin-top: auto
}

.alert-row {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 4px 7px;
  background: #fffbeb;
  border-radius: 6px;
  font-size: 9px;
  color: #92400e
}

.alert-row svg {
  width: 11px;
  flex-shrink: 0;
  color: #f59e0b
}

/* Snapshot */
.snapshot-row {
  display: grid;
  grid-template-columns: 1.15fr 1fr 1fr;
  gap: 14px
}

.panel {
  padding: 18px
}

.panel-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px
}

.panel-head>svg {
  width: 16px;
  color: #cbd5e1
}

.panel-head h2 {
  margin: 3px 0 2px;
  font-size: 16px;
  letter-spacing: -.02em;
  color: #111827
}

.panel-head p {
  margin: 0;
  font-size: 10px;
  color: #94a3b8
}

.health {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-top: 16px
}

.ring {
  width: 82px;
  height: 82px;
  border-radius: 50%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
  flex-shrink: 0
}

.ring::after {
  content: "";
  position: absolute;
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: #fff
}

.ring strong,
.ring span {
  position: relative;
  z-index: 1
}

.ring strong {
  font-size: 18px;
  color: #111827
}

.ring span {
  font-size: 8px;
  color: #94a3b8
}

.health-text strong {
  font-size: 12px;
  color: #334155
}

.health-text p {
  font-size: 10px;
  line-height: 1.5;
  color: #94a3b8;
  margin: 3px 0 0
}

.health-bars {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  margin-top: 14px;
  border-top: 1px solid #eef1f4;
  padding-top: 10px
}

.health-bars div {
  display: flex;
  flex-direction: column;
  gap: 2px
}

.health-bars span {
  font-size: 8px;
  color: #94a3b8;
  text-transform: uppercase
}

.health-bars b {
  font-size: 11px;
  color: #334155
}

/* Pipeline */
.pipe-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 14px;
  font-size: 10px;
  color: #475569
}

.pipe-label {
  display: flex;
  align-items: center;
  gap: 5px
}

.pipe-row strong {
  font-size: 11px;
  color: #334155
}

.pipe-track {
  height: 4px;
  background: #f1f5f9;
  border-radius: 99px;
  overflow: hidden;
  margin-top: 4px
}

.pipe-track i {
  display: block;
  height: 100%;
  border-radius: 99px;
  transition: width .5s ease
}

.fill-completed {
  background: #10b981
}

.fill-inprogress {
  background: #3b82f6
}

.fill-pending {
  background: #f59e0b
}

.dot-completed {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #10b981
}

.dot-inprogress {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #3b82f6
}

.dot-pending {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #f59e0b
}

/* Priority */
.prio-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 6px;
  margin-top: 14px
}

.prio {
  padding: 10px;
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  gap: 1px;
  transition: transform .15s
}

.prio:hover {
  transform: scale(1.02)
}

.prio strong {
  font-size: 18px
}

.prio span {
  font-size: 9px;
  font-weight: 800
}

.prio.urgent {
  background: #fef2f2;
  color: #b91c1c
}

.prio.high {
  background: #fff7ed;
  color: #c2410c
}

.prio.medium {
  background: #fffbeb;
  color: #a16207
}

.prio.low {
  background: #ecfdf5;
  color: #047857
}

.panel-note {
  font-size: 9px;
  color: #94a3b8;
  line-height: 1.4;
  margin: 10px 0 0
}

/* Charts */
.charts-row {
  display: grid;
  grid-template-columns: 1.45fr 1fr;
  gap: 14px
}

.charts-row.charts-3 {
  grid-template-columns: 1.2fr 1fr 1fr
}

.chart-panel {
  min-width: 0
}

.legend {
  display: flex;
  gap: 10px;
  font-size: 9px;
  font-weight: 700;
  color: #64748b
}

.legend span {
  display: flex;
  gap: 4px;
  align-items: center
}

.legend i {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  display: inline-block
}

.legend .green {
  background: #10b981
}

.legend .red {
  background: #ef4444
}

.trend-chip {
  display: flex;
  align-items: center;
  gap: 3px;
  padding: 4px 8px;
  border-radius: 8px;
  font-size: 10px;
  font-weight: 800
}

.trend-chip svg {
  width: 13px
}

.trend-chip.up {
  background: #ecfdf5;
  color: #047857
}

.trend-chip.down {
  background: #fef2f2;
  color: #dc2626
}

/* People Table */
.table-panel {
  padding: 0;
  overflow: hidden
}

.table-header {
  padding: 18px 20px;
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  border-bottom: 1px solid #eef1f4
}

.badge {
  padding: 5px 10px;
  border-radius: 8px;
  font-size: 10px;
  font-weight: 800;
  background: #f8fafc;
  border: 1px solid #e8ebf0;
  color: #64748b
}

.table-scroll {
  overflow: auto
}

.data-table {
  width: 100%;
  min-width: 900px;
  border-collapse: collapse
}

.data-table th {
  padding: 9px 14px;
  background: #fafbfc;
  border-bottom: 1px solid #e8ebf0;
  text-align: left;
  font-size: 8px;
  letter-spacing: .08em;
  color: #94a3b8
}

.data-table td {
  padding: 11px 14px;
  border-bottom: 1px solid #f1f3f6;
  font-size: 10px;
  color: #475569
}

.data-table tr:hover td {
  background: #f8fafc
}

.data-table tr:last-child td {
  border-bottom: 0
}

.rank {
  width: 22px;
  height: 22px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  background: #f1f5f9;
  font-size: 9px;
  font-weight: 800
}

.rank.gold {
  background: #fef3c7;
  color: #92400e
}

.emp {
  display: flex;
  align-items: center;
  gap: 7px
}

.avatar {
  width: 28px;
  height: 28px;
  border-radius: 7px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 8px;
  font-weight: 900;
  flex-shrink: 0
}

.avatar.sm {
  width: 26px;
  height: 26px
}

.av-0 {
  background: #eef2ff;
  color: #4338ca
}

.av-1 {
  background: #fef3c7;
  color: #92400e
}

.av-2 {
  background: #ecfdf5;
  color: #047857
}

.av-3 {
  background: #fce7f3;
  color: #be185d
}

.av-4 {
  background: #f0f9ff;
  color: #0369a1
}

.emp>span:last-child {
  display: flex;
  flex-direction: column;
  gap: 1px
}

.emp strong {
  font-size: 10px;
  color: #1e293b
}

.emp small {
  font-size: 8px;
  color: #94a3b8
}

.dept-tag {
  padding: 3px 6px;
  border-radius: 5px;
  font-size: 8px;
  font-weight: 800
}

.dept-housekeeping {
  background: #ecfdf5;
  color: #047857
}

.dept-reception {
  background: #eef2ff;
  color: #4338ca
}

.dept-inventory {
  background: #fffbeb;
  color: #a16207
}

.dept-finance {
  background: #f0f9ff;
  color: #0369a1
}

.dept-manager {
  background: #fce7f3;
  color: #be185d
}

.dept-admin {
  background: #f1f5f9;
  color: #475569
}

.sla-cell {
  min-width: 80px
}

.sla-cell>div {
  height: 4px;
  background: #e2e8f0;
  border-radius: 99px;
  margin-top: 3px
}

.sla-cell i {
  display: block;
  height: 100%;
  border-radius: 99px
}

.c-green {
  color: #047857
}

.c-blue {
  color: #1d4ed8
}

.c-red {
  color: #dc2626
}

.b-green {
  background: #10b981
}

.b-blue {
  background: #3b82f6
}

.b-red {
  background: #ef4444
}

.status-pill,
.rating {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 6px;
  border-radius: 5px;
  font-size: 8px;
  font-weight: 800
}

.status-pill i {
  width: 5px;
  height: 5px;
  border-radius: 50%
}

.status-pill.busy {
  background: #fff7ed;
  color: #c2410c
}

.status-pill.busy i {
  background: #f59e0b
}

.status-pill.free {
  background: #ecfdf5;
  color: #047857
}

.status-pill.free i {
  background: #10b981
}

.rating.excellent {
  background: #ecfdf5;
  color: #047857
}

.rating.good {
  background: #eff6ff;
  color: #1d4ed8
}

.rating.attention {
  background: #fff7ed;
  color: #c2410c
}

/* Bottom Row */
.bottom-row {
  display: grid;
  grid-template-columns: .8fr 1.7fr;
  gap: 14px
}

.leader-list {
  margin-top: 10px
}

.leader {
  display: grid;
  grid-template-columns: 22px 26px 1fr auto;
  align-items: center;
  gap: 7px;
  padding: 9px 0;
  border-bottom: 1px solid #f1f3f6;
  transition: background .15s
}

.leader:hover {
  background: #f8fafc;
  border-radius: 8px
}

.leader:last-child {
  border-bottom: 0
}

.leader.highlight {
  background: #fffbeb;
  border-radius: 10px;
  padding: 9px 8px;
  margin: 2px 0
}

.leader-rank {
  font-size: 9px;
  font-weight: 900;
  color: #cbd5e1;
  text-align: center
}

.leader-rank.gold {
  color: #f59e0b
}

.leader-rank.silver {
  color: #94a3b8
}

.leader-rank.bronze {
  color: #d97706
}

.leader-info,
.leader-score {
  display: flex;
  flex-direction: column;
  gap: 1px
}

.leader-info strong {
  font-size: 10px;
  color: #334155
}

.leader-info small,
.leader-score small {
  font-size: 8px;
  color: #94a3b8
}

.leader-score {
  text-align: right
}

.leader-score strong {
  font-size: 11px;
  color: #059669
}

.task-list {
  margin-top: 10px;
  display: flex;
  flex-direction: column
}

.task {
  display: grid;
  grid-template-columns: 50px 60px minmax(140px, 1fr) 50px 170px 70px;
  align-items: center;
  gap: 8px;
  padding: 9px 0;
  border-bottom: 1px solid #f1f3f6;
  transition: background .1s
}

.task:hover {
  background: #f8fafc;
  border-radius: 4px
}

.task:last-child {
  border-bottom: 0
}

.task-time {
  font-size: 9px;
  font-weight: 800;
  color: #64748b
}

.task-room {
  padding: 3px 5px;
  background: #111827;
  color: #fff;
  border-radius: 4px;
  font-size: 8px;
  font-weight: 800;
  text-align: center
}

.task-info {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0
}

.task-info strong {
  font-size: 9px;
  color: #334155;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap
}

.task-info small {
  font-size: 8px;
  color: #94a3b8;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap
}

.prio-tag,
.sla-tag {
  display: inline-flex;
  justify-content: center;
  padding: 3px 5px;
  border-radius: 4px;
  font-size: 7px;
  font-weight: 900;
  text-transform: uppercase
}

.prio-tag.urgent {
  background: #fee2e2;
  color: #b91c1c
}

.prio-tag.high {
  background: #ffedd5;
  color: #c2410c
}

.prio-tag.medium {
  background: #fef3c7;
  color: #a16207
}

.prio-tag.low {
  background: #ecfdf5;
  color: #047857
}

.task-sla {
  display: flex;
  flex-direction: column;
  gap: 2px
}

.sla-text {
  font-size: 8px;
  font-weight: 700;
  color: #475569
}

.sla-tag.met {
  background: #ecfdf5;
  color: #047857
}

.sla-tag.breached {
  background: #fef2f2;
  color: #dc2626;
  animation: pulse-badge 2s infinite
}

.sla-tag.active {
  background: #fffbeb;
  color: #d97706
}

@keyframes pulse-badge {

  0%,
  100% {
    opacity: 1
  }

  50% {
    opacity: .7
  }
}

/* Empty */
.empty {
  min-height: 180px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 5px;
  color: #94a3b8
}

.empty svg {
  width: 24px
}

.empty strong {
  font-size: 11px;
  color: #64748b
}

.empty span {
  font-size: 9px
}

.empty.compact {
  min-height: 120px
}

.empty-inline {
  padding: 28px;
  text-align: center;
  color: #94a3b8;
  font-size: 10px
}

/* Footer */
.footer-bar {
  display: grid;
  grid-template-columns: 1.4fr repeat(6, 1fr);
  padding: 12px 16px;
  gap: 10px
}

.footer-bar div {
  display: flex;
  flex-direction: column;
  gap: 2px
}

.footer-bar span {
  font-size: 8px;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: .05em
}

.footer-bar strong {
  font-size: 11px;
  color: #334155
}

.footer-bar .c-red {
  color: #dc2626
}

/* Responsive */
@media(max-width:1200px) {
  .kpi-row {
    grid-template-columns: repeat(2, 1fr)
  }

  .snapshot-row {
    grid-template-columns: 1fr 1fr
  }

  .snapshot-row> :first-child {
    grid-column: span 2
  }

  .charts-row,
  .charts-row.charts-3 {
    grid-template-columns: 1fr
  }

  .bottom-row {
    grid-template-columns: 1fr
  }

  .controls {
    grid-template-columns: 1fr 1fr
  }

  .control-block:last-child {
    grid-column: span 2
  }
}

@media(max-width:760px) {
  .dashboard {
    padding: 12px;
    gap: 12px
  }

  .hero-top {
    flex-direction: column;
    padding: 16px
  }

  .hero-copy h1 {
    font-size: 20px
  }

  .hero-status {
    align-self: flex-start
  }

  .controls {
    grid-template-columns: 1fr;
    padding: 10px 16px
  }

  .control-block:last-child {
    grid-column: auto
  }

  .dept-row {
    padding: 10px 16px
  }

  .kpi-row,
  .snapshot-row {
    grid-template-columns: 1fr
  }

  .snapshot-row> :first-child {
    grid-column: auto
  }

  .table-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px
  }

  .footer-bar {
    grid-template-columns: 1fr 1fr
  }

  .task {
    min-width: 580px
  }

  .task-list {
    overflow-x: auto
  }

  .dept-pills {
    overflow-x: auto;
    flex-wrap: nowrap
  }

  .dept-pills button {
    white-space: nowrap
  }
}
</style>
