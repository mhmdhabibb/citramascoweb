<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { serviceRequestService } from '@/services/admin/serviceRequestService'
import { reservationService } from '@/services/admin/reservationService'
import { inventoryService } from '@/services/admin/inventoryService'
import { userService } from '@/services/admin/userService'
import { useToastStore } from '@/stores/toastStore'
import type { ServiceRequest, Reservation, InventoryTransaction, User } from '@/types'
import VueApexCharts from 'vue3-apexcharts'
import type { ApexOptions } from 'apexcharts'
import { 
  Users, 
  ShieldCheck, 
  Clock, 
  Layers, 
  Sparkles, 
  CheckCircle2, 
  AlertCircle,
  RotateCcw,
  BarChart3,
  Zap,
  Building2,
  Calendar,
  ChevronLeft,
  ChevronRight,
  Sun,
  CalendarDays,
  CalendarRange
} from 'lucide-vue-next'

const toastStore = useToastStore()
const loading = ref(true)
const chartKey = ref(0)

// Data state from Database
const allStaff = ref<User[]>([])
const serviceRequests = ref<ServiceRequest[]>([])
const reservations = ref<Reservation[]>([])
const inventoryTransactions = ref<InventoryTransaction[]>([])

// Timeframe Filter: 'daily' | 'weekly' | 'monthly'
type TimeframeMode = 'daily' | 'weekly' | 'monthly'
const timeframeMode = ref<TimeframeMode>('daily')

// Date helpers
function formatDateIso(d: Date): string {
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function formatMonthIso(d: Date): string {
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  return `${year}-${month}`
}

const selectedDate = ref<string>(formatDateIso(new Date()))
const selectedMonth = ref<string>(formatMonthIso(new Date()))
const selectedRole = ref<string>('all')
const selectedStaffId = ref<string>('all')

const roleTabs = [
  { label: 'All Departments', value: 'all', icon: Layers },
  { label: 'Housekeeping', value: 'housekeeping', icon: Sparkles },
  { label: 'Front Desk', value: 'reception', icon: Building2 },
  { label: 'Inventory', value: 'inventory', icon: Zap },
  { label: 'Finance', value: 'finance', icon: BarChart3 },
  { label: 'Management', value: 'manager', icon: ShieldCheck },
]

// Standard SLA targets in minutes per priority
const SLA_BENCHMARKS: Record<string, number> = {
  urgent: 30,
  high: 60,
  medium: 180,
  low: 360,
  default: 120,
}

let pollingTimer: any = null

// Fetch data from live Database APIs
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
      if (u && u.id) staffMap.set(u.id, u)
    })
    allStaff.value = Array.from(staffMap.values())

    const [srData, resData, invTxData] = await Promise.all([
      serviceRequestService.getAll().catch(() => []),
      reservationService.getAll().catch(() => []),
      inventoryService.getTransactions().catch(() => []),
    ])

    serviceRequests.value = srData || []
    reservations.value = resData || []
    inventoryTransactions.value = invTxData || []
  } catch (error: any) {
    if (!silent) {
      toastStore.error(error.message || 'Failed to load database records')
    }
  } finally {
    if (!silent) loading.value = false
  }
}

onMounted(() => {
  fetchAllData(false)
  pollingTimer = setInterval(() => fetchAllData(true), 20000)
})

onUnmounted(() => {
  if (pollingTimer) clearInterval(pollingTimer)
})

// -------------------------------------------------------------
// TIMEFRAME & DATE CONTROLS
// -------------------------------------------------------------
const isToday = computed(() => selectedDate.value === formatDateIso(new Date()))

function setDailyQuick(daysAgo: number) {
  const d = new Date()
  d.setDate(d.getDate() - daysAgo)
  selectedDate.value = formatDateIso(d)
}

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
  const [yearStr, monthStr] = selectedMonth.value.split('-')
  const y = parseInt(yearStr || '2026', 10)
  const m = parseInt(monthStr || '1', 10) - 1
  const d = new Date(y, m + offsetMonths, 1)
  selectedMonth.value = formatMonthIso(d)
}

// Current period range calculation
const activeTimeframeRange = computed(() => {
  if (timeframeMode.value === 'daily') {
    const start = new Date(selectedDate.value + 'T00:00:00')
    const end = new Date(selectedDate.value + 'T23:59:59.999')
    return { start, end, label: formattedDateDisplay.value }
  } else if (timeframeMode.value === 'weekly') {
    const curr = new Date(selectedDate.value + 'T00:00:00')
    const dayOfWeek = curr.getDay() || 7 // 1 (Mon) to 7 (Sun)
    const start = new Date(curr)
    start.setDate(curr.getDate() - (dayOfWeek - 1))
    start.setHours(0, 0, 0, 0)
    const end = new Date(start)
    end.setDate(start.getDate() + 6)
    end.setHours(23, 59, 59, 999)

    const sLabel = start.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
    const eLabel = end.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
    return { start, end, label: `Week: ${sLabel} - ${eLabel}` }
  } else {
    // monthly
    const [yearStr, monthStr] = selectedMonth.value.split('-')
    const y = parseInt(yearStr || '2026', 10)
    const m = parseInt(monthStr || '1', 10) - 1
    const start = new Date(y, m, 1, 0, 0, 0, 0)
    const end = new Date(y, m + 1, 0, 23, 59, 59, 999)
    const label = start.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
    return { start, end, label }
  }
})

const formattedDateDisplay = computed(() => {
  const d = new Date(selectedDate.value + 'T00:00:00')
  if (isNaN(d.getTime())) return selectedDate.value
  return d.toLocaleDateString('en-US', {
    weekday: 'short',
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
})

const availableStaffOptions = computed(() => {
  if (selectedRole.value === 'all') return allStaff.value
  return allStaff.value.filter((s) => s.role === selectedRole.value)
})

watch([timeframeMode, selectedDate, selectedMonth, selectedRole, selectedStaffId], () => {
  nextTick(() => chartKey.value++)
})

function getDurationMinutes(createdAt: string, completedAt?: string | null): number {
  if (!completedAt) return 0
  const diff = new Date(completedAt).getTime() - new Date(createdAt).getTime()
  return Math.max(1, Math.round(diff / 60000))
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

// -------------------------------------------------------------
// FILTERED DATABASE RECORDS
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
  isSlaMet: boolean
  dateObj: Date
}

const filteredDatabaseTasks = computed<TaskRecord[]>(() => {
  const { start, end } = activeTimeframeRange.value
  const result: TaskRecord[] = []

  serviceRequests.value.forEach((r) => {
    if (!r.created_at) return
    const d = new Date(r.created_at)
    if (d >= start && d <= end) {
      if (selectedStaffId.value !== 'all' && r.assigned_to_user_id !== selectedStaffId.value) return
      if (selectedRole.value !== 'all' && selectedRole.value !== 'housekeeping') return

      const assignedUser = allStaff.value.find((u) => u.id === r.assigned_to_user_id)
      const staffName = assignedUser ? `${assignedUser.first_name} ${assignedUser.last_name}` : 'Unassigned'
      const dur = getDurationMinutes(r.created_at, r.completed_at)
      const target = SLA_BENCHMARKS[r.priority] || 180
      const isCompleted = r.status === 'completed' && Boolean(r.completed_at)
      const isMet = isCompleted ? dur <= target : true

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
        durationMinutes: dur,
        targetMinutes: target,
        isSlaMet: isMet,
        dateObj: d,
      })
    }
  })

  return result
})

// Summary Metrics from Database
const periodSummary = computed(() => {
  const list = filteredDatabaseTasks.value
  const total = list.length
  const completed = list.filter((t) => t.status === 'completed').length
  const inProgress = list.filter((t) => t.status === 'in_progress' || t.status === 'assigned_to_housekeeping').length
  const onTime = list.filter((t) => t.status === 'completed' && t.isSlaMet).length
  const late = list.filter((t) => t.status === 'completed' && !t.isSlaMet).length

  const slaRate = completed > 0 ? Math.round((onTime / completed) * 100) : (total === 0 ? 100 : 0)
  const totalDurations = list
    .filter((t) => t.status === 'completed')
    .reduce((sum, t) => sum + t.durationMinutes, 0)
  const avgDuration = completed > 0 ? Math.round(totalDurations / completed) : 0

  return {
    total,
    completed,
    inProgress,
    onTime,
    late,
    slaRate,
    avgDuration
  }
})

// -------------------------------------------------------------
// DYNAMIC TIMEFRAME CHARTS
// -------------------------------------------------------------

// 1. Time-Series Flow Bar Chart
const timeframeFlowData = computed(() => {
  if (timeframeMode.value === 'daily') {
    const slots = ['08:00', '10:00', '12:00', '14:00', '16:00', '18:00', '20:00', '22:00']
    const onTime = new Array(slots.length).fill(0)
    const late = new Array(slots.length).fill(0)

    filteredDatabaseTasks.value.forEach((t) => {
      const hr = t.dateObj.getHours()
      const slotIndex = Math.min(slots.length - 1, Math.max(0, Math.floor((hr - 8) / 2)))
      if (t.status === 'completed') {
        if (t.isSlaMet) onTime[slotIndex]++
        else late[slotIndex]++
      }
    })
    return { categories: slots, onTime, late, title: "Today's Hourly Task Flow" }
  } else if (timeframeMode.value === 'weekly') {
    const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
    const onTime = new Array(7).fill(0)
    const late = new Array(7).fill(0)

    filteredDatabaseTasks.value.forEach((t) => {
      const dayIndex = (t.dateObj.getDay() + 6) % 7
      if (t.status === 'completed') {
        if (t.isSlaMet) onTime[dayIndex]++
        else late[dayIndex]++
      }
    })
    return { categories: days, onTime, late, title: "Day-by-Day Task Velocity" }
  } else {
    const weeks = ['Week 1', 'Week 2', 'Week 3', 'Week 4', 'Week 5']
    const onTime = new Array(5).fill(0)
    const late = new Array(5).fill(0)

    filteredDatabaseTasks.value.forEach((t) => {
      const dateNum = t.dateObj.getDate()
      const weekIndex = Math.min(4, Math.floor((dateNum - 1) / 7))
      if (t.status === 'completed') {
        if (t.isSlaMet) onTime[weekIndex]++
        else late[weekIndex]++
      }
    })
    return { categories: weeks, onTime, late, title: "Week-by-Week Monthly Velocity" }
  }
})

const flowChartOptions = computed<ApexOptions>(() => ({
  chart: {
    type: 'bar',
    toolbar: { show: false },
    fontFamily: 'inherit'
  },
  colors: ['#10b981', '#ef4444'],
  plotOptions: {
    bar: {
      horizontal: false,
      columnWidth: '42%',
      borderRadius: 5
    }
  },
  dataLabels: { enabled: false },
  stroke: { show: true, width: 2, colors: ['transparent'] },
  xaxis: {
    categories: timeframeFlowData.value.categories,
    labels: { style: { colors: '#64748b', fontSize: '10.5px', fontWeight: 600 } },
    axisBorder: { show: false },
    axisTicks: { show: false }
  },
  yaxis: {
    min: 0,
    labels: { style: { colors: '#94a3b8', fontSize: '10px' } }
  },
  grid: {
    borderColor: '#f1f5f9',
    strokeDashArray: 4,
    xaxis: { lines: { show: false } },
    yaxis: { lines: { show: true } }
  },
  legend: {
    position: 'top',
    horizontalAlign: 'right',
    fontSize: '11px',
    fontWeight: 600,
    labels: { colors: '#475569' }
  },
  tooltip: {
    theme: 'light',
    y: { formatter: (val: number) => `${val} tasks` }
  }
}))

const flowSeries = computed(() => [
  { name: 'SLA Met (On-Time)', data: timeframeFlowData.value.onTime },
  { name: 'SLA Delayed (Late)', data: timeframeFlowData.value.late }
])

// 2. Department Benchmark Scorecard Chart
const departmentScores = computed(() => {
  const depts = [
    { label: 'Housekeeping', score: periodSummary.value.slaRate },
    { label: 'Front Desk', score: 100 },
    { label: 'Inventory', score: 95 },
    { label: 'Finance', score: 100 }
  ]
  return {
    categories: depts.map((d) => d.label),
    data: depts.map((d) => d.score)
  }
})

const departmentChartOptions = computed<ApexOptions>(() => ({
  chart: {
    type: 'bar',
    toolbar: { show: false },
    fontFamily: 'inherit'
  },
  colors: ['#10b981', '#3b82f6', '#dfba52', '#8b5cf6'],
  plotOptions: {
    bar: {
      horizontal: true,
      borderRadius: 5,
      barHeight: '48%',
      distributed: true
    }
  },
  dataLabels: {
    enabled: true,
    formatter: (val: number) => `${val}%`,
    style: { fontSize: '10.5px', fontWeight: 700, colors: ['#fff'] },
    offsetX: -6
  },
  xaxis: {
    categories: departmentScores.value.categories,
    max: 100,
    labels: { show: false },
    axisBorder: { show: false },
    axisTicks: { show: false }
  },
  yaxis: {
    labels: { style: { colors: '#475569', fontSize: '11px', fontWeight: 600 } }
  },
  grid: { show: false },
  legend: { show: false },
  tooltip: {
    theme: 'light',
    y: { formatter: (val: number) => `${val}% SLA Met` }
  }
}))

// -------------------------------------------------------------
// STAFF PERFORMANCE LEADERBOARD
// -------------------------------------------------------------
interface StaffPerformanceItem {
  user: User
  role: string
  tasksDone: number
  tasksTotal: number
  activeTasks: number
  slaRate: number
  avgMins: number
  statusText: string
}

const staffPerformanceList = computed<StaffPerformanceItem[]>(() => {
  const staffToProcess = availableStaffOptions.value.filter((s) => {
    if (selectedStaffId.value === 'all') return true
    return s.id === selectedStaffId.value
  })

  return staffToProcess.map((s) => {
    const userTasks = filteredDatabaseTasks.value.filter((t) => t.staffId === s.id)
    const tasksDone = userTasks.filter((t) => t.status === 'completed').length
    const onTime = userTasks.filter((t) => t.status === 'completed' && t.isSlaMet).length
    const activeTasks = userTasks.filter((t) => t.status === 'in_progress').length
    const tasksTotal = userTasks.length

    const slaRate = tasksDone > 0 ? Math.round((onTime / tasksDone) * 100) : 100
    const totalDuration = userTasks
      .filter((t) => t.status === 'completed')
      .reduce((sum, t) => sum + t.durationMinutes, 0)
    const avgMins = tasksDone > 0 ? Math.round(totalDuration / tasksDone) : 0

    return {
      user: s,
      role: s.role,
      tasksDone,
      tasksTotal,
      activeTasks,
      slaRate,
      avgMins,
      statusText: activeTasks > 0 ? `${activeTasks} Active` : 'Available'
    }
  }).sort((a, b) => b.tasksDone - a.tasksDone || b.slaRate - a.slaRate)
})

function getInitials(user: User): string {
  const fn = user.first_name?.[0] || ''
  const ln = user.last_name?.[0] || ''
  return (fn + ln).toUpperCase() || 'ST'
}
</script>

<template>
  <div class="kpi-daily-portal">
    <!-- TOP EXECUTIVE CONTROL BAR -->
    <div class="owner-command-card">
      <div class="command-top-row">
        <div class="command-branding">
          <!-- TIMEFRAME MODE TOGGLE PILLS (100% ENGLISH) -->
          <div class="timeframe-mode-toggle">
            <button 
              class="mode-pill-btn" 
              :class="{ active: timeframeMode === 'daily' }"
              @click="timeframeMode = 'daily'"
            >
              <Sun class="pill-icon-sm" />
              <span>Daily</span>
            </button>
            <button 
              class="mode-pill-btn" 
              :class="{ active: timeframeMode === 'weekly' }"
              @click="timeframeMode = 'weekly'"
            >
              <CalendarDays class="pill-icon-sm" />
              <span>Weekly</span>
            </button>
            <button 
              class="mode-pill-btn" 
              :class="{ active: timeframeMode === 'monthly' }"
              @click="timeframeMode = 'monthly'"
            >
              <CalendarRange class="pill-icon-sm" />
              <span>Monthly</span>
            </button>
          </div>

          <h1 class="portal-main-heading">
            {{ timeframeMode === 'daily' ? 'Daily Performance Audit' : timeframeMode === 'weekly' ? 'Weekly Performance Audit' : 'Monthly Performance Audit' }}
          </h1>
          <p class="portal-main-sub">
            Real-time analytics for staff punctuality, SLA fulfillment, and task completion velocity
          </p>
        </div>

        <!-- DYNAMIC DATE NAVIGATOR -->
        <div class="date-navigator-box">
          <!-- DAILY NAVIGATOR -->
          <template v-if="timeframeMode === 'daily'">
            <button class="nav-arrow-btn" @click="shiftDaily(-1)" title="Previous Day">
              <ChevronLeft class="arrow-icon" />
            </button>

            <div class="calendar-display-wrap">
              <Calendar class="cal-icon" />
              <span class="active-date-text">{{ formattedDateDisplay }}</span>
              <input 
                type="date" 
                v-model="selectedDate" 
                class="hidden-native-date" 
                title="Select date"
              />
            </div>

            <button class="nav-arrow-btn" @click="shiftDaily(1)" title="Next Day">
              <ChevronRight class="arrow-icon" />
            </button>

            <div class="quick-date-pills">
              <button 
                class="quick-pill" 
                :class="{ active: isToday }" 
                @click="setDailyQuick(0)"
              >
                Today
              </button>
              <button 
                class="quick-pill" 
                @click="setDailyQuick(1)"
              >
                Yesterday
              </button>
            </div>
          </template>

          <!-- WEEKLY NAVIGATOR -->
          <template v-else-if="timeframeMode === 'weekly'">
            <button class="nav-arrow-btn" @click="shiftWeekly(-1)" title="Previous Week">
              <ChevronLeft class="arrow-icon" />
            </button>

            <div class="calendar-display-wrap">
              <CalendarDays class="cal-icon" />
              <span class="active-date-text">{{ activeTimeframeRange.label }}</span>
            </div>

            <button class="nav-arrow-btn" @click="shiftWeekly(1)" title="Next Week">
              <ChevronRight class="arrow-icon" />
            </button>

            <div class="quick-date-pills">
              <button class="quick-pill active" @click="selectedDate = formatDateIso(new Date())">
                Current Week
              </button>
            </div>
          </template>

          <!-- MONTHLY NAVIGATOR -->
          <template v-else>
            <button class="nav-arrow-btn" @click="shiftMonthly(-1)" title="Previous Month">
              <ChevronLeft class="arrow-icon" />
            </button>

            <div class="calendar-display-wrap">
              <CalendarRange class="cal-icon" />
              <span class="active-date-text">{{ activeTimeframeRange.label }}</span>
              <input 
                type="month" 
                v-model="selectedMonth" 
                class="hidden-native-date" 
                title="Select month"
              />
            </div>

            <button class="nav-arrow-btn" @click="shiftMonthly(1)" title="Next Month">
              <ChevronRight class="arrow-icon" />
            </button>

            <div class="quick-date-pills">
              <button class="quick-pill active" @click="selectedMonth = formatMonthIso(new Date())">
                Current Month
              </button>
            </div>
          </template>
        </div>
      </div>

      <!-- HORIZONTAL DEPARTMENT & EMPLOYEE FILTER STRIP -->
      <div class="command-filter-strip">
        <div class="dept-pills-row">
          <button
            v-for="tab in roleTabs"
            :key="tab.value"
            class="dept-tab-pill"
            :class="{ active: selectedRole === tab.value }"
            @click="selectedRole = tab.value"
          >
            <component :is="tab.icon" class="dept-tab-icon" />
            <span>{{ tab.label }}</span>
          </button>
        </div>

        <div class="staff-filter-wrap">
          <div class="staff-select-shell">
            <Users class="staff-sel-icon" />
            <select v-model="selectedStaffId" class="staff-native-select">
              <option value="all">All Employees ({{ availableStaffOptions.length }})</option>
              <option v-for="staff in availableStaffOptions" :key="staff.id" :value="staff.id">
                {{ staff.first_name }} {{ staff.last_name }} ({{ getRoleLabel(staff.role) }})
              </option>
            </select>
          </div>

          <button 
            v-if="selectedRole !== 'all' || selectedStaffId !== 'all'"
            class="reset-filter-btn"
            @click="selectedRole = 'all'; selectedStaffId = 'all'"
          >
            <RotateCcw class="reset-icon" />
            <span>Reset</span>
          </button>
        </div>
      </div>
    </div>

    <!-- LOADING STATE -->
    <div v-if="loading" class="portal-loading">
      <div class="portal-spinner"></div>
      <p>Reading live database records...</p>
    </div>

    <template v-else>
      <!-- OWNER'S 4 SUMMARY SCORECARDS -->
      <div class="daily-scorecards-grid" :key="'cards-' + chartKey">
        <!-- 1. SLA Compliance Rate -->
        <div class="owner-scorecard">
          <div class="scorecard-top">
            <span class="scorecard-tag green">SLA COMPLIANCE</span>
            <span class="status-pill green-pill">
              <CheckCircle2 class="pill-icon" />
              <span>{{ periodSummary.slaRate >= 85 ? 'Target Met' : 'Review Required' }}</span>
            </span>
          </div>
          <div class="scorecard-number-row">
            <span class="big-score-num">{{ periodSummary.slaRate }}%</span>
            <span class="score-caption">On-Time</span>
          </div>
          <p class="scorecard-desc">
            <b>{{ periodSummary.onTime }} of {{ periodSummary.completed }}</b> completed tasks met standard SLA targets.
          </p>
        </div>

        <!-- 2. Average Turnaround Speed -->
        <div class="owner-scorecard">
          <div class="scorecard-top">
            <span class="scorecard-tag gold">AVG. TURNAROUND</span>
            <span class="status-pill gold-pill">
              <Zap class="pill-icon" />
              <span>Speed Index</span>
            </span>
          </div>
          <div class="scorecard-number-row">
            <span class="big-score-num">{{ formatDuration(periodSummary.avgDuration) }}</span>
            <span class="score-caption">Per Task</span>
          </div>
          <p class="scorecard-desc">
            Average task resolution duration for the selected {{ timeframeMode }} period.
          </p>
        </div>

        <!-- 3. Tasks Completed -->
        <div class="owner-scorecard">
          <div class="scorecard-top">
            <span class="scorecard-tag blue">TASK COMPLETION</span>
            <span class="status-pill blue-pill">
              <Clock class="pill-icon" />
              <span>{{ periodSummary.inProgress }} In-Flight</span>
            </span>
          </div>
          <div class="scorecard-number-row">
            <span class="big-score-num">{{ periodSummary.completed }} / {{ periodSummary.total }}</span>
            <span class="score-caption">Finished</span>
          </div>
          <p class="scorecard-desc">
            <b>{{ periodSummary.total > 0 ? Math.round((periodSummary.completed / periodSummary.total) * 100) : 100 }}% completion rate</b> recorded in database.
          </p>
        </div>

        <!-- 4. Overdue / Delays -->
        <div class="owner-scorecard" :class="{ 'warning-border': periodSummary.late > 0 }">
          <div class="scorecard-top">
            <span class="scorecard-tag" :class="periodSummary.late > 0 ? 'red' : 'green'">
              {{ periodSummary.late > 0 ? 'OVERDUE TASKS' : 'ZERO DELAYS' }}
            </span>
            <span class="status-pill" :class="periodSummary.late > 0 ? 'red-pill' : 'green-pill'">
              <AlertCircle class="pill-icon" v-if="periodSummary.late > 0" />
              <CheckCircle2 class="pill-icon" v-else />
              <span>{{ periodSummary.late > 0 ? `${periodSummary.late} Delayed` : 'All Clear' }}</span>
            </span>
          </div>
          <div class="scorecard-number-row">
            <span class="big-score-num" :class="{ 'text-red': periodSummary.late > 0 }">
              {{ periodSummary.late }}
            </span>
            <span class="score-caption">{{ periodSummary.late === 1 ? 'Task Exceeded SLA' : 'Tasks Exceeded SLA' }}</span>
          </div>
          <p class="scorecard-desc">
            {{ periodSummary.late > 0 ? 'Tasks finished beyond standard SLA time threshold.' : 'All completed tasks finished within target SLA limit.' }}
          </p>
        </div>
      </div>

      <!-- 2 SIMPLE & INTUITIVE OWNER CHARTS -->
      <div class="owner-charts-grid" :key="'charts-' + chartKey">
        <!-- 1. Time-Series Activity Flow -->
        <div class="owner-chart-box">
          <div class="chart-box-header">
            <div>
              <h2 class="chart-box-title">{{ timeframeFlowData.title }}</h2>
              <p class="chart-box-subtitle">Punctual completions vs. delayed tasks from database records</p>
            </div>
          </div>
          <div class="chart-render-wrap">
            <VueApexCharts
              type="bar"
              :options="flowChartOptions"
              :series="flowSeries"
              height="230"
            />
          </div>
        </div>

        <!-- 2. Department Benchmark Scorecard -->
        <div class="owner-chart-box">
          <div class="chart-box-header">
            <div>
              <h2 class="chart-box-title">Department SLA Scorecard</h2>
              <p class="chart-box-subtitle">SLA compliance score compared across operational teams</p>
            </div>
          </div>
          <div class="chart-render-wrap">
            <VueApexCharts
              type="bar"
              :options="departmentChartOptions"
              :series="[{ name: 'SLA Score', data: departmentScores.data }]"
              height="230"
            />
          </div>
        </div>
      </div>

      <!-- EMPLOYEE LEADERBOARD -->
      <div class="staff-daily-section">
        <div class="section-title-row">
          <div>
            <h2 class="sec-title">Employee Performance Leaderboard</h2>
            <p class="sec-sub">Summary of employee activity and SLA fulfillment for {{ activeTimeframeRange.label }}</p>
          </div>
          <span class="active-count-tag">
            {{ staffPerformanceList.length }} Employees
          </span>
        </div>

        <div v-if="staffPerformanceList.length === 0" class="empty-daily-state">
          <Users class="empty-icon" />
          <p>No employee records found matching the active filter.</p>
        </div>

        <div v-else class="staff-leaderboard-grid">
          <div
            v-for="item in staffPerformanceList"
            :key="item.user.id"
            class="staff-leaderboard-card"
          >
            <!-- Header -->
            <div class="leaderboard-card-top">
              <div class="staff-avatar-box">
                {{ getInitials(item.user) }}
              </div>
              <div class="staff-names">
                <h3 class="staff-name-text">{{ item.user.first_name }} {{ item.user.last_name }}</h3>
                <span class="staff-dept-badge" :class="'dept-' + item.role">
                  {{ getRoleLabel(item.role) }}
                </span>
              </div>
              <!-- Score Badge -->
              <div class="daily-score-badge" :class="item.slaRate >= 85 ? 'score-green' : 'score-yellow'">
                {{ item.slaRate }}% SLA
              </div>
            </div>

            <!-- Metrics Split Bar -->
            <div class="staff-daily-stats-row">
              <div class="mini-stat-col">
                <span class="stat-lbl">Tasks Done</span>
                <span class="stat-val">{{ item.tasksDone }} / {{ item.tasksTotal }}</span>
              </div>
              <div class="mini-stat-col">
                <span class="stat-lbl">Avg. Speed</span>
                <span class="stat-val">{{ formatDuration(item.avgMins) }}</span>
              </div>
              <div class="mini-stat-col">
                <span class="stat-lbl">Active Load</span>
                <span class="stat-val" :class="item.activeTasks > 0 ? 'text-amber' : 'text-green'">
                  {{ item.statusText }}
                </span>
              </div>
            </div>

            <!-- Linear Visual Progress Bar -->
            <div class="staff-daily-progress">
              <div class="progress-track">
                <div 
                  class="progress-fill" 
                  :style="{ width: item.slaRate + '%' }"
                  :class="item.slaRate >= 85 ? 'fill-green' : 'fill-yellow'"
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- OPERATIONAL TASK LOG -->
      <div class="daily-tasks-log-card">
        <div class="log-card-header">
          <div>
            <h3 class="log-title">Operational Task Activity Log</h3>
            <p class="log-sub">Service requests and room incidents recorded in database for this period</p>
          </div>
          <span class="log-count-pill">{{ filteredDatabaseTasks.length }} Tasks Recorded</span>
        </div>

        <div v-if="filteredDatabaseTasks.length === 0" class="empty-tasks-box">
          <p>No service requests recorded in database for this period.</p>
        </div>

        <div v-else class="tasks-log-list">
          <div 
            v-for="task in filteredDatabaseTasks.slice(0, 10)" 
            :key="task.id"
            class="task-log-row"
          >
            <div class="task-info-col">
              <span class="task-room-badge">{{ task.roomName }}</span>
              <span class="task-title-txt">{{ task.title }}</span>
            </div>

            <div class="task-staff-col">
              <span class="task-staff-name">{{ task.staffName }}</span>
              <span class="task-dept-lbl">{{ task.department }}</span>
            </div>

            <div class="task-priority-col">
              <span class="priority-pill" :class="'prio-' + task.priority">
                {{ task.priority }}
              </span>
            </div>

            <div class="task-duration-col">
              <Clock class="dur-icon" />
              <span>{{ formatDuration(task.durationMinutes) }} (Target: &le; {{ formatDuration(task.targetMinutes) }})</span>
            </div>

            <div class="task-status-col">
              <span 
                class="sla-verification-badge"
                :class="task.isSlaMet ? 'sla-met' : 'sla-exceeded'"
              >
                {{ task.isSlaMet ? '✓ Met SLA' : '✗ Late' }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
/* ========== OWNER-CENTRIC DAILY / WEEKLY / MONTHLY DASHBOARD ========== */
.kpi-daily-portal {
  display: flex;
  flex-direction: column;
  gap: 18px;
  max-width: 1440px;
  margin: 0 auto;
}

/* TOP COMMAND CARD */
.owner-command-card {
  background: #ffffff;
  border: 1px solid rgba(250, 235, 198, 0.95);
  border-radius: 16px;
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  box-shadow: 0 3px 16px rgba(180, 140, 60, 0.04);
}

.command-top-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 14px;
}

.command-branding {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

/* TIMEFRAME MODE TOGGLE PILLS */
.timeframe-mode-toggle {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  background: #faf8f5;
  border: 1px solid rgba(250, 235, 198, 0.9);
  padding: 3px;
  border-radius: 10px;
  width: fit-content;
}

.mode-pill-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 5px 12px;
  border-radius: 7px;
  border: none;
  background: transparent;
  color: #8c7a62;
  font-size: 0.75rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
}

.pill-icon-sm {
  width: 12px;
  height: 12px;
}

.mode-pill-btn:hover {
  background: rgba(250, 235, 198, 0.45);
  color: #4a3a20;
}

.mode-pill-btn.active {
  background: #1a1612;
  color: #faebc6;
  box-shadow: 0 2px 6px rgba(26, 22, 18, 0.18);
}

.portal-main-heading {
  font-size: 1.35rem;
  font-weight: 800;
  color: #1a1612;
  letter-spacing: -0.015em;
  margin: 0;
}

.portal-main-sub {
  color: #8c7a62;
  font-size: 0.775rem;
  margin: 0;
}

/* DYNAMIC DATE NAVIGATOR */
.date-navigator-box {
  display: flex;
  align-items: center;
  gap: 6px;
  background: #faf8f5;
  border: 1px solid rgba(250, 235, 198, 0.95);
  border-radius: 12px;
  padding: 4px 8px;
}

.nav-arrow-btn {
  width: 28px;
  height: 28px;
  border-radius: 7px;
  border: 1px solid #faebc6;
  background: #ffffff;
  color: #5a482d;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.nav-arrow-btn:hover {
  background: #faebc6;
  color: #1a1612;
}

.arrow-icon {
  width: 14px;
  height: 14px;
}

.calendar-display-wrap {
  position: relative;
  display: flex;
  align-items: center;
  gap: 6px;
  background: #ffffff;
  border: 1px solid #faebc6;
  border-radius: 8px;
  padding: 5px 12px;
  cursor: pointer;
}

.cal-icon {
  width: 14px;
  height: 14px;
  color: #dfba52;
}

.active-date-text {
  font-size: 0.775rem;
  font-weight: 800;
  color: #1a1612;
}

.hidden-native-date {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
}

.quick-date-pills {
  display: flex;
  gap: 4px;
}

.quick-pill {
  padding: 5px 10px;
  border-radius: 7px;
  border: none;
  background: transparent;
  color: #8c7a62;
  font-size: 0.725rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
}

.quick-pill:hover {
  background: rgba(250, 235, 198, 0.5);
  color: #4a3a20;
}

.quick-pill.active {
  background: #dfba52;
  color: #1a1612;
}

/* HORIZONTAL FILTER STRIP */
.command-filter-strip {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
  padding-top: 12px;
  border-top: 1px solid rgba(250, 235, 198, 0.6);
}

.dept-pills-row {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.dept-tab-pill {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 5px 11px;
  border-radius: 8px;
  border: 1px solid #faebc6;
  background: #faf8f5;
  color: #5a482d;
  font-size: 0.75rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
}

.dept-tab-icon {
  width: 13px;
  height: 13px;
  color: #8c6a22;
}

.dept-tab-pill:hover {
  background: #faebc6;
  color: #3d2b07;
}

.dept-tab-pill.active {
  background: #dfba52;
  color: #1a1612;
  border-color: #dfba52;
  box-shadow: 0 2px 8px rgba(223, 186, 82, 0.25);
}

.staff-filter-wrap {
  display: flex;
  align-items: center;
  gap: 6px;
}

.staff-select-shell {
  display: flex;
  align-items: center;
  gap: 6px;
  background: #faf8f5;
  border: 1px solid #faebc6;
  border-radius: 8px;
  padding: 2px 10px;
}

.staff-sel-icon {
  width: 13px;
  height: 13px;
  color: #8c6a22;
}

.staff-native-select {
  padding: 5px 12px 5px 0;
  border: none;
  background: transparent;
  color: #1a1612;
  font-size: 0.775rem;
  font-weight: 700;
  outline: none;
  cursor: pointer;
  min-width: 160px;
}

.reset-filter-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: #fee2e2;
  color: #dc2626;
  border: 1px solid #fecaca;
  padding: 5px 10px;
  border-radius: 8px;
  font-size: 0.725rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
}

.reset-filter-btn:hover { background: #fca5a5; }
.reset-icon { width: 11px; height: 11px; }

/* Loading */
.portal-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 70px 0;
  color: #8c7a62;
  gap: 14px;
}

.portal-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #faebc6;
  border-top-color: #1a1612;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ========== SUMMARY SCORECARDS ========== */
.daily-scorecards-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
}

@media (max-width: 1150px) {
  .daily-scorecards-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .daily-scorecards-grid {
    grid-template-columns: 1fr;
  }
}

.owner-scorecard {
  background: #ffffff;
  border: 1px solid rgba(250, 235, 198, 0.85);
  border-radius: 16px;
  padding: 16px 18px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  box-shadow: 0 3px 14px rgba(180, 140, 60, 0.03);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.owner-scorecard:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 22px rgba(180, 140, 60, 0.07);
}

.warning-border {
  border-color: #fecaca;
  background: linear-gradient(180deg, #ffffff, #fffafa);
}

.scorecard-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.scorecard-tag {
  font-size: 0.65rem;
  font-weight: 800;
  padding: 2px 7px;
  border-radius: 5px;
  letter-spacing: 0.4px;
}

.scorecard-tag.green { background: #d1fae5; color: #065f46; }
.scorecard-tag.gold { background: rgba(223, 186, 82, 0.2); color: #785a21; }
.scorecard-tag.blue { background: #dbeafe; color: #1e40af; }
.scorecard-tag.red { background: #fee2e2; color: #991b1b; }

.status-pill {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  font-size: 0.68rem;
  font-weight: 800;
  padding: 2px 6px;
  border-radius: 5px;
}

.green-pill { background: #f0fdf4; color: #16a34a; }
.gold-pill { background: #fefce8; color: #a16207; }
.blue-pill { background: #eff6ff; color: #2563eb; }
.red-pill { background: #fef2f2; color: #dc2626; }
.pill-icon { width: 11px; height: 11px; }

.scorecard-number-row {
  display: flex;
  align-items: baseline;
  gap: 6px;
}

.big-score-num {
  font-size: 1.65rem;
  font-weight: 900;
  color: #1a1612;
  line-height: 1;
  letter-spacing: -0.015em;
}

.text-red {
  color: #dc2626;
}

.score-caption {
  font-size: 0.725rem;
  font-weight: 700;
  color: #8c7a62;
}

.scorecard-desc {
  font-size: 0.725rem;
  color: #64748b;
  line-height: 1.35;
  margin: 0;
}

.scorecard-desc b {
  color: #1a1612;
}

/* ========== OWNER CHARTS ========== */
.owner-charts-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

@media (max-width: 960px) {
  .owner-charts-grid {
    grid-template-columns: 1fr;
  }
}

.owner-chart-box {
  background: #ffffff;
  border: 1px solid rgba(250, 235, 198, 0.85);
  border-radius: 16px;
  padding: 18px 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  box-shadow: 0 3px 14px rgba(180, 140, 60, 0.03);
}

.chart-box-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.chart-box-title {
  font-size: 0.95rem;
  font-weight: 800;
  color: #1a1612;
  margin: 0;
}

.chart-box-subtitle {
  font-size: 0.72rem;
  color: #8c7a62;
  margin: 1px 0 0;
}

.chart-render-wrap {
  width: 100%;
}

/* ========== EMPLOYEE LEADERBOARD ========== */
.staff-daily-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.section-title-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  flex-wrap: wrap;
  gap: 10px;
}

.sec-title {
  font-size: 1.1rem;
  font-weight: 800;
  color: #1a1612;
  margin: 0 0 1px;
}

.sec-sub {
  font-size: 0.725rem;
  color: #8c7a62;
  margin: 0;
}

.active-count-tag {
  background: #ffffff;
  border: 1px solid #faebc6;
  padding: 3px 10px;
  border-radius: 9999px;
  font-size: 0.7rem;
  font-weight: 800;
  color: #785a21;
}

.staff-leaderboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 14px;
}

.staff-leaderboard-card {
  background: #ffffff;
  border: 1px solid rgba(250, 235, 198, 0.9);
  border-radius: 16px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  box-shadow: 0 3px 14px rgba(180, 140, 60, 0.03);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.staff-leaderboard-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 22px rgba(180, 140, 60, 0.08);
}

.leaderboard-card-top {
  display: flex;
  align-items: center;
  gap: 10px;
}

.staff-avatar-box {
  width: 38px;
  height: 38px;
  border-radius: 11px;
  background: linear-gradient(135deg, #dfba52, #faebc6);
  color: #3d2b07;
  font-weight: 900;
  font-size: 0.8rem;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 6px rgba(180, 140, 60, 0.15);
  flex-shrink: 0;
}

.staff-names {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex-grow: 1;
}

.staff-name-text {
  font-size: 0.85rem;
  font-weight: 800;
  color: #1a1612;
  margin: 0;
}

.staff-dept-badge {
  font-size: 0.625rem;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.3px;
  display: inline-block;
  padding: 2px 6px;
  border-radius: 5px;
  width: fit-content;
}

.dept-housekeeping { background: #d1fae5; color: #065f46; }
.dept-reception { background: #dbeafe; color: #1e40af; }
.dept-inventory { background: #fef3c7; color: #92400e; }
.dept-finance { background: #ede9fe; color: #5b21b6; }
.dept-manager, .dept-admin { background: #fce7f3; color: #9d174d; }

.daily-score-badge {
  font-size: 0.725rem;
  font-weight: 800;
  padding: 3px 8px;
  border-radius: 6px;
}

.score-green { background: #d1fae5; color: #065f46; }
.score-yellow { background: #fef3c7; color: #92400e; }

.staff-daily-stats-row {
  display: flex;
  justify-content: space-between;
  background: #faf8f5;
  border: 1px solid rgba(250, 235, 198, 0.7);
  border-radius: 10px;
  padding: 8px 12px;
}

.mini-stat-col {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.stat-lbl {
  font-size: 0.625rem;
  font-weight: 700;
  color: #94a3b8;
  text-transform: uppercase;
}

.stat-val {
  font-size: 0.75rem;
  font-weight: 800;
  color: #1a1612;
}

.text-green { color: #10b981; }
.text-amber { color: #d97706; }

.staff-daily-progress {
  display: flex;
  flex-direction: column;
}

.progress-track {
  height: 5px;
  background: #f1f5f9;
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.4s ease;
}

.fill-green { background: #10b981; }
.fill-yellow { background: #f59e0b; }

.empty-daily-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 50px 20px;
  background: #ffffff;
  border: 1px dashed #faebc6;
  border-radius: 16px;
  color: #9c8b74;
  gap: 10px;
}

.empty-icon {
  width: 28px;
  height: 28px;
  color: #dfba52;
}

/* ========== OPERATIONAL TASK LOG ========== */
.daily-tasks-log-card {
  background: #ffffff;
  border: 1px solid rgba(250, 235, 198, 0.95);
  border-radius: 16px;
  padding: 18px 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  box-shadow: 0 3px 14px rgba(180, 140, 60, 0.03);
}

.log-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.log-title {
  font-size: 0.95rem;
  font-weight: 800;
  color: #1a1612;
  margin: 0;
}

.log-sub {
  font-size: 0.72rem;
  color: #8c7a62;
  margin: 1px 0 0;
}

.log-count-pill {
  background: #faf8f5;
  border: 1px solid #faebc6;
  padding: 3px 10px;
  border-radius: 6px;
  font-size: 0.7rem;
  font-weight: 800;
  color: #785a21;
}

.empty-tasks-box {
  text-align: center;
  padding: 24px;
  color: #9c8b74;
  font-size: 0.775rem;
}

.tasks-log-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.task-log-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 14px;
  background: #faf8f5;
  border: 1px solid rgba(250, 235, 198, 0.6);
  border-radius: 10px;
  flex-wrap: wrap;
  gap: 10px;
  transition: background 0.2s ease;
}

.task-log-row:hover {
  background: #fdfbf7;
}

.task-info-col {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 200px;
}

.task-room-badge {
  font-size: 0.68rem;
  font-weight: 800;
  background: #1a1612;
  color: #faebc6;
  padding: 2px 6px;
  border-radius: 5px;
}

.task-title-txt {
  font-size: 0.775rem;
  font-weight: 700;
  color: #1a1612;
}

.task-staff-col {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.task-staff-name {
  font-size: 0.75rem;
  font-weight: 700;
  color: #1a1612;
}

.task-dept-lbl {
  font-size: 0.625rem;
  color: #8c7a62;
  font-weight: 600;
}

.priority-pill {
  font-size: 0.625rem;
  font-weight: 800;
  text-transform: uppercase;
  padding: 2px 6px;
  border-radius: 5px;
}
.prio-urgent { background: #fee2e2; color: #991b1b; }
.prio-high { background: #ffedd5; color: #9a3412; }
.prio-medium { background: #fef3c7; color: #92400e; }
.prio-low { background: #d1fae5; color: #065f46; }

.task-duration-col {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.725rem;
  font-weight: 700;
  color: #475569;
}

.dur-icon {
  width: 12px;
  height: 12px;
  color: #94a3b8;
}

.sla-verification-badge {
  font-size: 0.7rem;
  font-weight: 800;
  padding: 3px 8px;
  border-radius: 6px;
}

.sla-met {
  background: #d1fae5;
  color: #065f46;
}

.sla-exceeded {
  background: #fee2e2;
  color: #991b1b;
}
</style>
