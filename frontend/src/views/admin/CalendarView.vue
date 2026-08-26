<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { reservationService } from '@/services/admin/reservationService'
import { roomService } from '@/services/roomService'
import { channelService } from '@/services/admin/channelService'
import { useToastStore } from '@/stores/toastStore'
import type { Reservation, Room, Channel } from '@/types'
import {
  Calendar as CalendarIcon,
  ChevronLeft,
  ChevronRight,
  Filter,
  Users,
  BedDouble,
  DollarSign,
  TrendingUp,
  X,
  Layers,
  Search,
  CheckCircle2,
  Clock,
  ExternalLink,
  Radio,
  RefreshCw,
} from 'lucide-vue-next'

const router = useRouter()
const toastStore = useToastStore()

// State from Database
const loading = ref(false)
const reservations = ref<Reservation[]>([])
const rooms = ref<Room[]>([])
const channels = ref<Channel[]>([])
const selectedProperty = ref<string>('all')
const selectedChannel = ref<string>('all')
const showCancelled = ref<boolean>(false)
const viewMode = ref<'timeline' | 'month'>('timeline')
const searchQuery = ref<string>('')

// Date navigation
const currentDate = ref(new Date())
const currentYear = computed(() => currentDate.value.getFullYear())
const currentMonth = computed(() => currentDate.value.getMonth())

const monthNames = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember',
]

const dayNamesShort = ['Min', 'Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab']

// --- ROBUST DATE PARSING HELPERS (Handle DD-MM-YYYY from Go Backend & ISO) ---
function parseDateToISO(d?: string): string {
  if (!d) return ''
  const trimmed = d.trim()
  // Check if DD-MM-YYYY format (e.g. "14-08-2026")
  const parts = trimmed.split('-')
  if (parts.length === 3 && parts[0]?.length === 2 && parts[2]?.length === 4) {
    const day = parts[0]
    const month = parts[1]
    const year = parts[2]
    return `${year}-${month}-${day}`
  }
  if (trimmed.includes('T')) {
    return trimmed.slice(0, 10)
  }
  if (parts.length === 3 && parts[0]?.length === 4) {
    return `${parts[0]}-${parts[1]?.padStart(2, '0')}-${parts[2]?.padStart(2, '0')}`
  }
  return trimmed
}

function parseDateObj(d?: string): Date | null {
  const iso = parseDateToISO(d)
  if (!iso) return null
  const dt = new Date(iso)
  return isNaN(dt.getTime()) ? null : dt
}

// Dynamic Channel Palette from DB
const dynamicColorPalette = [
  '#16a34a', // Emerald
  '#2563eb', // Blue
  '#e11d48', // Rose
  '#7c3aed', // Purple
  '#0d9488', // Teal
  '#0284c7', // Sky
  '#d97706', // Amber
  '#ea580c', // Orange
  '#8b5cf6', // Violet
  '#ec4899', // Pink
]

const getChannelColor = (name?: string, index: number = 0): string => {
  if (!name) return '#16a34a'
  const lower = name.toLowerCase()
  if (lower.includes('direct') || lower.includes('website') || lower.includes('walk')) return '#16a34a'
  if (lower.includes('booking')) return '#2563eb'
  if (lower.includes('airbnb')) return '#e11d48'
  if (lower.includes('agoda')) return '#7c3aed'
  if (lower.includes('trip')) return '#0d9488'
  if (lower.includes('traveloka')) return '#0284c7'
  if (lower.includes('tiket')) return '#f59e0b'
  return dynamicColorPalette[index % dynamicColorPalette.length] || '#4f46e5'
}

const getBookingChannelKey = (res: Reservation): string => {
  if (res.channel?.id) return res.channel.id
  if (res.channel_id) return res.channel_id
  if (res.channel?.name) {
    const ch = channels.value.find((c) => c.name.toLowerCase() === res.channel?.name.toLowerCase())
    if (ch) return ch.id
  }
  return 'direct'
}

const getBookingChannelName = (res: Reservation): string => {
  if (res.channel?.name) return res.channel.name
  if (res.channel_id) {
    const ch = channels.value.find((c) => c.id === res.channel_id)
    if (ch) return ch.name
  }
  return 'Direct Website'
}

const getChannelStyle = (channelKeyOrName?: string) => {
  const name = channelKeyOrName || 'Direct Website'
  const ch = channels.value.find((c) => c.id === channelKeyOrName || c.name.toLowerCase() === name.toLowerCase())
  const label = ch ? ch.name : (name === 'direct' ? 'Direct Website' : name)
  const color = getChannelColor(label, 0)
  return {
    bg: color,
    text: '#ffffff',
    label,
    border: color,
    lightBg: `${color}15`,
  }
}

const formatIDR = (val: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(val || 0)
}

const formatYMD = (d: Date): string => {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

// Fetch Live Data
const fetchData = async () => {
  try {
    loading.value = true
    const [resList, rList, cList] = await Promise.all([
      reservationService.getAll(),
      roomService.getAll().catch(() => []),
      channelService.getAll().catch(() => []),
    ])
    reservations.value = resList
    rooms.value = rList
    channels.value = cList
  } catch (err: any) {
    toastStore.error(err.message || 'Failed to load calendar data')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})

// Effective Rooms directly from DB
const effectiveRooms = computed(() => {
  if (selectedProperty.value === 'all') return rooms.value
  return rooms.value.filter((r) => r.id === selectedProperty.value)
})

// Date calculations
const daysInCurrentMonth = computed(() => {
  return new Date(currentYear.value, currentMonth.value + 1, 0).getDate()
})

// Month Days Array for Timeline
const timelineDays = computed(() => {
  const year = currentYear.value
  const month = currentMonth.value
  const totalDays = daysInCurrentMonth.value
  const todayStr = formatYMD(new Date())

  const list = []
  for (let d = 1; d <= totalDays; d++) {
    const dateObj = new Date(year, month, d)
    const dateStr = formatYMD(dateObj)
    list.push({
      dayNumber: d,
      date: dateObj,
      dateStr,
      dayName: dayNamesShort[dateObj.getDay()] ?? '',
      isWeekend: dateObj.getDay() === 0 || dateObj.getDay() === 6,
      isToday: dateStr === todayStr,
    })
  }
  return list
})

// Month Grid Days (7 columns SUN-SAT)
interface CalendarDay {
  date: Date
  dateStr: string
  dayNumber: number
  isCurrentMonth: boolean
  isToday: boolean
}

const calendarDays = computed(() => {
  const year = currentYear.value
  const month = currentMonth.value

  const firstDayOfMonth = new Date(year, month, 1)
  const lastDayOfMonth = new Date(year, month + 1, 0)
  const days: CalendarDay[] = []
  const todayStr = formatYMD(new Date())

  const startDayOfWeek = firstDayOfMonth.getDay()
  for (let i = startDayOfWeek - 1; i >= 0; i--) {
    const d = new Date(year, month, -i)
    const dateStr = formatYMD(d)
    days.push({
      date: d,
      dateStr,
      dayNumber: d.getDate(),
      isCurrentMonth: false,
      isToday: dateStr === todayStr,
    })
  }

  for (let i = 1; i <= lastDayOfMonth.getDate(); i++) {
    const d = new Date(year, month, i)
    const dateStr = formatYMD(d)
    days.push({
      date: d,
      dateStr,
      dayNumber: i,
      isCurrentMonth: true,
      isToday: dateStr === todayStr,
    })
  }

  const remaining = (7 - (days.length % 7)) % 7
  for (let i = 1; i <= remaining; i++) {
    const d = new Date(year, month + 1, i)
    const dateStr = formatYMD(d)
    days.push({
      date: d,
      dateStr,
      dayNumber: i,
      isCurrentMonth: false,
      isToday: dateStr === todayStr,
    })
  }
  return days
})

const rangeDisplay = computed(() => {
  return `${monthNames[currentMonth.value]} ${currentYear.value}`
})

// Navigation methods
const prevMonth = () => {
  currentDate.value = new Date(currentYear.value, currentMonth.value - 1, 1)
}

const nextMonth = () => {
  currentDate.value = new Date(currentYear.value, currentMonth.value + 1, 1)
}

const goToToday = () => {
  currentDate.value = new Date()
}

// Filtered Bookings for the current month from Database
const activeBookings = computed(() => {
  return reservations.value.filter((b) => {
    if (!showCancelled.value && (b.status === 'cancel' || b.status === 'rejected')) {
      return false
    }
    if (selectedProperty.value !== 'all' && b.room_id !== selectedProperty.value) {
      return false
    }
    if (selectedChannel.value !== 'all' && getBookingChannelKey(b) !== selectedChannel.value) {
      return false
    }
    if (searchQuery.value.trim()) {
      const q = searchQuery.value.toLowerCase()
      const matchName = (b.full_name || '').toLowerCase().includes(q)
      const matchCode = (b.code || '').toLowerCase().includes(q)
      if (!matchName && !matchCode) return false
    }
    return true
  })
})

// KPI Metrics for Current Month (100% from Database)
const arrivalsCount = computed(() => {
  const currentYM = `${currentYear.value}-${String(currentMonth.value + 1).padStart(2, '0')}`
  return activeBookings.value.filter((b) => {
    const checkinISO = parseDateToISO(b.checkin_date)
    return checkinISO.startsWith(currentYM)
  }).length
})

const departuresCount = computed(() => {
  const currentYM = `${currentYear.value}-${String(currentMonth.value + 1).padStart(2, '0')}`
  return activeBookings.value.filter((b) => {
    const checkoutISO = parseDateToISO(b.checkout_date)
    return checkoutISO.startsWith(currentYM)
  }).length
})

const bookedNightsCount = computed(() => {
  return activeBookings.value.reduce((acc, b) => acc + (b.total_night || 1), 0)
})

const occupancyRate = computed(() => {
  const totalRooms = Math.max(1, rooms.value.length)
  const totalPossibleNights = totalRooms * daysInCurrentMonth.value
  const percent = (bookedNightsCount.value / totalPossibleNights) * 100
  return Math.min(100, Math.round(percent * 10) / 10)
})

// Timeline Position helper for a reservation in a specific room
const getTimelineBookingSpan = (b: Reservation) => {
  const inDate = parseDateObj(b.checkin_date)
  const outDate = parseDateObj(b.checkout_date)
  if (!inDate || !outDate) return null

  const daysInMonth = daysInCurrentMonth.value

  const startDay = inDate.getMonth() === currentMonth.value && inDate.getFullYear() === currentYear.value
    ? inDate.getDate()
    : 1

  const endDay = outDate.getMonth() === currentMonth.value && outDate.getFullYear() === currentYear.value
    ? outDate.getDate()
    : daysInMonth

  if (startDay > daysInMonth || endDay < 1) return null

  const span = Math.max(1, endDay - startDay)
  return {
    gridColumnStart: startDay,
    gridColumnEnd: `span ${span}`,
    startDay,
    span,
  }
}

// Get bookings for a specific room
const getBookingsForRoom = (roomId: string) => {
  return activeBookings.value.filter((b) => b.room_id === roomId)
}

// Get bookings for a specific calendar day
const getBookingsForDay = (dateStr: string) => {
  return activeBookings.value.filter((b) => {
    const inISO = parseDateToISO(b.checkin_date)
    const outISO = parseDateToISO(b.checkout_date)
    if (!inISO || !outISO) return false
    return dateStr >= inISO && dateStr <= outISO
  })
}

// Modal Detail State
const selectedBooking = ref<Reservation | null>(null)
const isDetailModalOpen = ref(false)

const openBookingDetail = (b: Reservation) => {
  selectedBooking.value = b
  isDetailModalOpen.value = true
}

const closeBookingDetail = () => {
  isDetailModalOpen.value = false
  selectedBooking.value = null
}

const goToChannels = () => router.push('/admin/channels')
</script>

<template>
  <div class="calendar-page">
    <!-- Top Header Bar -->
    <div class="header-card">
      <div class="header-left">
        <div class="header-title-wrap">
          <div class="icon-avatar">
            <CalendarIcon class="w-5 h-5 text-indigo-600" />
          </div>
          <div>
            <h1 class="page-title">Room & Booking Calendar</h1>
            <p class="page-desc">Room & Booking Calendar matrix timeline & channel distribution from database</p>
          </div>
        </div>
      </div>

      <div class="header-controls">
        <!-- Refresh Button -->
        <button @click="fetchData" :disabled="loading" class="btn-refresh-cal" title="Refresh Database Data">
          <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
          <span>{{ loading ? 'Loading...' : 'Sync DB' }}</span>
        </button>

        <!-- Search -->
        <div class="search-input-wrap">
          <Search class="search-icon w-4 h-4 text-slate-400" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search guest / booking code..."
            class="search-field"
          />
        </div>

        <!-- Channel Filter Dropdown -->
        <div class="filter-select-wrap">
          <Radio class="w-4 h-4 text-slate-400 select-icon" />
          <select v-model="selectedChannel" class="custom-select">
            <option value="all">All Channels ({{ channels.length }})</option>
            <option v-for="ch in channels" :key="ch.id" :value="ch.id">
              {{ ch.name }}
            </option>
            <option value="direct">Direct Website</option>
          </select>
        </div>

        <!-- Property Filter -->
        <div class="filter-select-wrap">
          <BedDouble class="w-4 h-4 text-slate-400 select-icon" />
          <select v-model="selectedProperty" class="custom-select">
            <option value="all">All Rooms ({{ rooms.length }})</option>
            <option v-for="r in rooms" :key="r.id" :value="r.id">
              {{ r.name }}
            </option>
          </select>
        </div>

        <!-- View Mode Switcher -->
        <div class="view-switcher">
          <button
            class="switch-btn"
            :class="{ active: viewMode === 'timeline' }"
            @click="viewMode = 'timeline'"
          >
            <Layers class="w-4 h-4" />
            <span>Timeline PMS</span>
          </button>
          <button
            class="switch-btn"
            :class="{ active: viewMode === 'month' }"
            @click="viewMode = 'month'"
          >
            <CalendarIcon class="w-4 h-4" />
            <span>Month Grid</span>
          </button>
        </div>

        <button @click="goToChannels" class="btn-channels-link" title="Manage Master Channel">
          <Radio class="w-4 h-4" />
          <span>Manage Channel</span>
        </button>
      </div>
    </div>

    <!-- 4 KPI Summary Cards -->
    <div class="kpi-grid">
      <div class="kpi-card">
        <div class="kpi-icon arrivals-icon">
          <Users class="w-5 h-5" />
        </div>
        <div class="kpi-content">
          <span class="kpi-label">Arrivals (Check-In)</span>
          <span class="kpi-val">{{ arrivalsCount }} <span class="kpi-unit">guests</span></span>
        </div>
      </div>

      <div class="kpi-card">
        <div class="kpi-icon departures-icon">
          <CheckCircle2 class="w-5 h-5" />
        </div>
        <div class="kpi-content">
          <span class="kpi-label">Departures (Check-Out)</span>
          <span class="kpi-val">{{ departuresCount }} <span class="kpi-unit">guests</span></span>
        </div>
      </div>

      <div class="kpi-card">
        <div class="kpi-icon nights-icon">
          <BedDouble class="w-5 h-5" />
        </div>
        <div class="kpi-content">
          <span class="kpi-label">Booked Nights</span>
          <span class="kpi-val">{{ bookedNightsCount }} <span class="kpi-unit">nights</span></span>
        </div>
      </div>

      <div class="kpi-card">
        <div class="kpi-icon occupancy-icon">
          <TrendingUp class="w-5 h-5" />
        </div>
        <div class="kpi-content">
          <div class="occ-header">
            <span class="kpi-label">Occupancy Rate</span>
            <span class="occ-badge">{{ occupancyRate }}%</span>
          </div>
          <div class="progress-bar-bg">
            <div class="progress-bar-fill" :style="{ width: `${occupancyRate}%` }"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Date Navigation Toolbar -->
    <div class="nav-toolbar">
      <div class="nav-left">
        <button @click="goToToday" class="today-btn">Today</button>
        <div class="nav-arrows">
          <button @click="prevMonth" class="nav-arrow-btn" title="Previous Month">
            <ChevronLeft class="w-4 h-4" />
          </button>
          <button @click="nextMonth" class="nav-arrow-btn" title="Next Month">
            <ChevronRight class="w-4 h-4" />
          </button>
        </div>
        <span class="month-title">{{ rangeDisplay }}</span>
      </div>

      <div class="nav-right">
        <label class="toggle-wrap">
          <input type="checkbox" v-model="showCancelled" />
          <span class="toggle-label">Show Cancelled</span>
        </label>
      </div>
    </div>

    <!-- MAIN VIEW 1: TIMELINE PMS VIEW (Room vs Days) -->
    <div v-if="viewMode === 'timeline'" class="timeline-container">
      <div class="timeline-scroll-wrap">
        <div class="timeline-table">
          <!-- Timeline Header Row: Days 1..N -->
          <div class="timeline-header-row">
            <div class="room-col-header">
              <span>Room Unit</span>
            </div>
            <div
              class="days-header-grid"
              :style="{ gridTemplateColumns: `repeat(${daysInCurrentMonth}, minmax(44px, 1fr))` }"
            >
              <div
                v-for="d in timelineDays"
                :key="d.dayNumber"
                class="day-col-head"
                :class="{ 'is-today': d.isToday, 'is-weekend': d.isWeekend }"
              >
                <span class="day-name-label">{{ d.dayName }}</span>
                <span class="day-num-label">{{ d.dayNumber }}</span>
              </div>
            </div>
          </div>

          <!-- Timeline Room Rows -->
          <div class="timeline-body">
            <div
              v-for="room in effectiveRooms"
              :key="room.id"
              class="timeline-room-row"
            >
              <!-- Room Label Column -->
              <div class="room-label-cell">
                <div class="room-cell-info">
                  <span class="room-title">{{ room.name }}</span>
                  <span class="room-price-tag">{{ formatIDR(room.price || 0) }} / night</span>
                </div>
              </div>

              <!-- Day Cells & Booking Bars Grid -->
              <div
                class="room-days-grid"
                :style="{ gridTemplateColumns: `repeat(${daysInCurrentMonth}, minmax(44px, 1fr))` }"
              >
                <!-- Background day column grid lines -->
                <div
                  v-for="d in timelineDays"
                  :key="d.dayNumber"
                  class="day-slot-cell"
                  :class="{ 'is-today-slot': d.isToday, 'is-weekend-slot': d.isWeekend }"
                ></div>

                <!-- Booking Bars Layer from Database -->
                <template v-for="b in getBookingsForRoom(room.id)" :key="b.id">
                  <div
                    v-if="getTimelineBookingSpan(b)"
                    class="timeline-bar"
                    :style="{
                      gridColumnStart: getTimelineBookingSpan(b)!.gridColumnStart,
                      gridColumnEnd: getTimelineBookingSpan(b)!.gridColumnEnd,
                      backgroundColor: getChannelStyle(getBookingChannelKey(b)).bg,
                      color: getChannelStyle(getBookingChannelKey(b)).text,
                    }"
                    @click="openBookingDetail(b)"
                  >
                    <span class="bar-channel-tag">
                      {{ getChannelStyle(getBookingChannelKey(b)).label }}
                    </span>
                    <span class="bar-guest-name">{{ b.full_name || 'Guest' }}</span>
                    <span class="bar-nights">({{ b.total_night || 1 }}n)</span>
                  </div>
                </template>
              </div>
            </div>

            <!-- Empty State for Timeline -->
            <div v-if="effectiveRooms.length === 0" class="empty-timeline-state">
              <BedDouble class="w-10 h-10 text-slate-300 mx-auto mb-2" />
              <p class="font-bold text-slate-700">No room units in database yet</p>
              <p class="text-xs text-slate-400 mt-1">Add room units via the 'Manage Rooms' &rarr; 'Rooms' menu</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- MAIN VIEW 2: MONTH GRID VIEW -->
    <div v-else class="month-grid-container">
      <div class="month-grid-header">
        <div v-for="dayName in ['MIN', 'SEN', 'SEL', 'RAB', 'KAM', 'JUM', 'SAB']" :key="dayName" class="grid-head-cell">
          {{ dayName }}
        </div>
      </div>

      <div class="month-grid-body">
        <div
          v-for="(cell, index) in calendarDays"
          :key="index"
          class="month-day-cell"
          :class="{
            'is-other-month': !cell.isCurrentMonth,
            'is-today-cell': cell.isToday,
          }"
        >
          <div class="day-cell-top">
            <span
              class="day-number"
              :class="{
                'sunday-num': cell.date.getDay() === 0,
                'today-badge': cell.isToday,
              }"
            >
              {{ cell.dayNumber }}
            </span>
          </div>

          <div class="day-bookings-list">
            <div
              v-for="b in getBookingsForDay(cell.dateStr)"
              :key="b.id"
              class="month-booking-pill"
              :style="{
                backgroundColor: getChannelStyle(getBookingChannelKey(b)).bg,
                color: getChannelStyle(getBookingChannelKey(b)).text,
              }"
              @click="openBookingDetail(b)"
              :title="`${b.full_name} (${b.room?.name || 'Room'}) - ${formatIDR(b.total_price)}`"
            >
              <span class="booking-inner-name">&lt; {{ b.full_name || 'Guest' }} &gt;</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Channel Legends Footer -->
    <div class="channel-legends-card">
      <div class="legends-title-wrap">
        <Radio class="w-4 h-4 text-indigo-600" />
        <span class="legends-title">Channel Legend:</span>
      </div>
      <div class="legends-flex">
        <div
          v-for="ch in channels"
          :key="ch.id"
          class="legend-item"
          :class="{ 'is-selected': selectedChannel === ch.id }"
          @click="selectedChannel = selectedChannel === ch.id ? 'all' : ch.id"
        >
          <span class="legend-dot" :style="{ backgroundColor: getChannelStyle(ch.id).bg }"></span>
          <span class="legend-text">{{ ch.name }}</span>
        </div>
        <div
          class="legend-item"
          :class="{ 'is-selected': selectedChannel === 'direct' }"
          @click="selectedChannel = selectedChannel === 'direct' ? 'all' : 'direct'"
        >
          <span class="legend-dot" style="background-color: #16a34a;"></span>
          <span class="legend-text">Direct Website</span>
        </div>
      </div>
    </div>

    <!-- Booking Detail Modal -->
    <div v-if="isDetailModalOpen && selectedBooking" class="modal-backdrop" @click.self="closeBookingDetail">
      <div class="modal-container">
        <div class="modal-top">
          <div class="modal-channel-header">
            <span
              class="modal-channel-pill"
              :style="{
                backgroundColor: getChannelStyle(getBookingChannelKey(selectedBooking)).bg,
                color: getChannelStyle(getBookingChannelKey(selectedBooking)).text,
              }"
            >
              {{ getChannelStyle(getBookingChannelKey(selectedBooking)).label }}
            </span>
            <span class="modal-code">#{{ selectedBooking.code }}</span>
          </div>
          <button @click="closeBookingDetail" class="modal-close-btn">&times;</button>
        </div>

        <div class="modal-body">
          <h2 class="guest-title">{{ selectedBooking.full_name }}</h2>
          <p class="guest-email">{{ selectedBooking.email || 'No email provided' }}</p>

          <div class="detail-divider"></div>

          <div class="modal-grid-info">
            <div class="info-item">
              <span class="info-label">Room / Unit</span>
              <span class="info-val font-semibold">{{ selectedBooking.room?.name || 'Room Unit' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Reservation Status</span>
              <span class="status-chip" :class="`status-${selectedBooking.status}`">
                {{ selectedBooking.status }}
              </span>
            </div>
            <div class="info-item">
              <span class="info-label">Check-In</span>
              <span class="info-val">{{ selectedBooking.checkin_date }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Check-Out</span>
              <span class="info-val">{{ selectedBooking.checkout_date }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Stay Duration</span>
              <span class="info-val">{{ selectedBooking.total_night || 1 }} Night(s)</span>
            </div>
            <div class="info-item">
              <span class="info-label">Number of Guests</span>
              <span class="info-val">{{ (selectedBooking.number_of_adult || 1) + (selectedBooking.number_of_children || 0) }} Person(s)</span>
            </div>
          </div>

          <div class="price-highlight-box">
            <div class="price-row">
              <span>Rate per Night</span>
              <span>{{ formatIDR(selectedBooking.price || 0) }}</span>
            </div>
            <div class="price-row total-row">
              <span>Total Cost</span>
              <span class="total-amount">{{ formatIDR(selectedBooking.total_price || 0) }}</span>
            </div>
            <div v-if="selectedBooking.deposit" class="price-row deposit-row">
              <span>Deposit Paid</span>
              <span class="text-emerald-600 font-semibold">{{ formatIDR(selectedBooking.deposit) }}</span>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button @click="closeBookingDetail" class="btn-close-modal">Close</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.calendar-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
  background-color: #f8fafc;
  font-family: 'Plus Jakarta Sans', system-ui, -apple-system, sans-serif;
  color: #0f172a;
}

/* Header Card */
.header-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
  background: #ffffff;
  padding: 20px 24px;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.header-title-wrap {
  display: flex;
  align-items: center;
  gap: 14px;
}

.icon-avatar {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: #eef2ff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.page-title {
  font-size: 1.35rem;
  font-weight: 800;
  color: #0f172a;
  letter-spacing: -0.02em;
  margin: 0;
}

.page-desc {
  font-size: 0.85rem;
  color: #64748b;
  margin: 2px 0 0;
}

.header-controls {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.btn-refresh-cal {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  font-size: 0.8rem;
  font-weight: 700;
  border-radius: 10px;
  border: 1px solid #cbd5e1;
  background: #ffffff;
  color: #334155;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-refresh-cal:hover {
  background: #f8fafc;
  border-color: #94a3b8;
}

.search-input-wrap {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
  pointer-events: none;
}

.search-field {
  padding: 8px 12px 8px 36px;
  font-size: 0.85rem;
  border-radius: 10px;
  border: 1px solid #cbd5e1;
  background: #f8fafc;
  outline: none;
  min-width: 200px;
  transition: all 0.2s;
}

.search-field:focus {
  border-color: #4f46e5;
  background: #ffffff;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

.filter-select-wrap {
  position: relative;
  display: flex;
  align-items: center;
}

.select-icon {
  position: absolute;
  left: 12px;
  pointer-events: none;
}

.custom-select {
  padding: 8px 16px 8px 34px;
  font-size: 0.85rem;
  font-weight: 600;
  border-radius: 10px;
  border: 1px solid #cbd5e1;
  background: #f8fafc;
  color: #334155;
  outline: none;
  cursor: pointer;
}

.view-switcher {
  display: flex;
  background: #f1f5f9;
  padding: 4px;
  border-radius: 12px;
  gap: 4px;
}

.switch-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  font-size: 0.82rem;
  font-weight: 700;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: #64748b;
  cursor: pointer;
  transition: all 0.18s;
}

.switch-btn.active {
  background: #ffffff;
  color: #4f46e5;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.06);
}

.btn-channels-link {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  font-size: 0.82rem;
  font-weight: 700;
  border-radius: 10px;
  border: 1px solid #e0e7ff;
  background: #eef2ff;
  color: #4338ca;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-channels-link:hover {
  background: #e0e7ff;
}

/* KPI Grid */
.kpi-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

@media (max-width: 1024px) {
  .kpi-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.kpi-card {
  background: #ffffff;
  border-radius: 16px;
  padding: 16px 18px;
  display: flex;
  align-items: center;
  gap: 14px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
}

.kpi-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.arrivals-icon { background: #eff6ff; color: #2563eb; }
.departures-icon { background: #f0fdf4; color: #16a34a; }
.nights-icon { background: #fdf4ff; color: #a855f7; }
.occupancy-icon { background: #fff1f2; color: #e11d48; }

.kpi-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.kpi-label {
  font-size: 0.78rem;
  font-weight: 600;
  color: #64748b;
}

.kpi-val {
  font-size: 1.3rem;
  font-weight: 800;
  color: #0f172a;
}

.kpi-unit {
  font-size: 0.8rem;
  font-weight: 500;
  color: #94a3b8;
}

.occ-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.occ-badge {
  font-size: 0.85rem;
  font-weight: 800;
  color: #e11d48;
}

.progress-bar-bg {
  width: 100%;
  height: 6px;
  background: #fee2e2;
  border-radius: 9999px;
  overflow: hidden;
}

.progress-bar-fill {
  height: 100%;
  background: #e11d48;
  border-radius: 9999px;
  transition: width 0.3s ease;
}

/* Nav Toolbar */
.nav-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #ffffff;
  padding: 12px 20px;
  border-radius: 14px;
  border: 1px solid #e2e8f0;
}

.nav-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.today-btn {
  padding: 6px 14px;
  font-size: 0.82rem;
  font-weight: 700;
  border-radius: 8px;
  border: 1px solid #cbd5e1;
  background: #ffffff;
  color: #334155;
  cursor: pointer;
  transition: all 0.2s;
}

.today-btn:hover {
  background: #f1f5f9;
}

.nav-arrows {
  display: flex;
  gap: 4px;
}

.nav-arrow-btn {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  border: 1px solid #cbd5e1;
  background: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #475569;
  cursor: pointer;
  transition: all 0.2s;
}

.nav-arrow-btn:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.month-title {
  font-size: 1.1rem;
  font-weight: 800;
  color: #0f172a;
  margin-left: 6px;
}

.toggle-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.82rem;
  font-weight: 600;
  color: #64748b;
  cursor: pointer;
}

/* TIMELINE PMS VIEW */
.timeline-container {
  background: #ffffff;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  overflow: hidden;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.03);
}

.timeline-scroll-wrap {
  overflow-x: auto;
  min-width: 100%;
}

.timeline-table {
  min-width: 1200px;
  display: flex;
  flex-direction: column;
}

.timeline-header-row {
  display: flex;
  background: #f8fafc;
  border-bottom: 2px solid #e2e8f0;
  position: sticky;
  top: 0;
  z-index: 10;
}

.room-col-header {
  width: 220px;
  min-width: 220px;
  padding: 14px 16px;
  font-size: 0.82rem;
  font-weight: 800;
  color: #475569;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  border-right: 1px solid #e2e8f0;
  background: #f8fafc;
  display: flex;
  align-items: center;
}

.days-header-grid {
  flex: 1;
  display: grid;
}

.day-col-head {
  padding: 10px 4px;
  text-align: center;
  border-right: 1px solid #f1f5f9;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.day-name-label {
  font-size: 0.68rem;
  font-weight: 700;
  color: #94a3b8;
  text-transform: uppercase;
}

.day-num-label {
  font-size: 0.88rem;
  font-weight: 800;
  color: #1e293b;
  margin-top: 2px;
}

.day-col-head.is-today {
  background: #e0e7ff;
}

.day-col-head.is-today .day-num-label {
  color: #4338ca;
}

.day-col-head.is-weekend {
  background: #fbfbfe;
}

.day-col-head.is-weekend .day-name-label {
  color: #f43f5e;
}

/* Timeline Room Rows */
.timeline-body {
  display: flex;
  flex-direction: column;
}

.timeline-room-row {
  display: flex;
  border-bottom: 1px solid #f1f5f9;
  min-height: 64px;
  position: relative;
}

.timeline-room-row:hover {
  background: #fafbff;
}

.room-label-cell {
  width: 220px;
  min-width: 220px;
  padding: 12px 16px;
  border-right: 1px solid #e2e8f0;
  background: #ffffff;
  display: flex;
  align-items: center;
}

.room-title {
  font-size: 0.85rem;
  font-weight: 700;
  color: #0f172a;
  display: block;
}

.room-price-tag {
  font-size: 0.72rem;
  font-weight: 600;
  color: #64748b;
  display: block;
  margin-top: 2px;
}

.room-days-grid {
  flex: 1;
  display: grid;
  position: relative;
}

.day-slot-cell {
  border-right: 1px solid #f8fafc;
  height: 100%;
}

.day-slot-cell.is-today-slot {
  background: rgba(224, 231, 255, 0.4);
}

.day-slot-cell.is-weekend-slot {
  background: rgba(248, 250, 252, 0.5);
}

/* Booking Bar on Timeline */
.timeline-bar {
  position: absolute;
  top: 12px;
  bottom: 12px;
  left: 2px;
  right: 2px;
  border-radius: 8px;
  padding: 4px 10px;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  z-index: 5;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.12);
  transition: transform 0.15s, box-shadow 0.15s;
  overflow: hidden;
  white-space: nowrap;
}

.timeline-bar:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.18);
}

.bar-channel-tag {
  font-size: 0.65rem;
  font-weight: 800;
  text-transform: uppercase;
  background: rgba(255, 255, 255, 0.25);
  padding: 1px 6px;
  border-radius: 4px;
  flex-shrink: 0;
}

.bar-guest-name {
  font-size: 0.78rem;
  font-weight: 700;
  overflow: hidden;
  text-overflow: ellipsis;
}

.bar-nights {
  font-size: 0.7rem;
  opacity: 0.85;
}

.empty-timeline-state {
  text-align: center;
  padding: 40px 16px;
}

/* MONTH GRID VIEW */
.month-grid-container {
  background: #ffffff;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  overflow: hidden;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.month-grid-header {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  background: #f8fafc;
  border-bottom: 2px solid #e2e8f0;
}

.grid-head-cell {
  padding: 12px;
  text-align: center;
  font-size: 0.78rem;
  font-weight: 800;
  color: #475569;
}

.month-grid-body {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
}

.month-day-cell {
  min-height: 110px;
  border-right: 1px solid #f1f5f9;
  border-bottom: 1px solid #f1f5f9;
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.month-day-cell:nth-child(7n) {
  border-right: none;
}

.month-day-cell.is-other-month {
  background: #f8fafc;
  opacity: 0.45;
}

.month-day-cell.is-today-cell {
  background: #faf5ff;
}

.day-cell-top {
  display: flex;
  justify-content: flex-end;
}

.day-number {
  font-size: 0.82rem;
  font-weight: 700;
  color: #334155;
}

.day-number.sunday-num {
  color: #f43f5e;
}

.day-number.today-badge {
  background: #4f46e5;
  color: #ffffff;
  width: 24px;
  height: 24px;
  border-radius: 9999px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.day-bookings-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.month-booking-pill {
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 700;
  cursor: pointer;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

/* Channel Legends Card */
.channel-legends-card {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
  background: #ffffff;
  padding: 14px 20px;
  border-radius: 14px;
  border: 1px solid #e2e8f0;
}

.legends-title-wrap {
  display: flex;
  align-items: center;
  gap: 6px;
}

.legends-title {
  font-size: 0.82rem;
  font-weight: 800;
  color: #475569;
}

.legends-flex {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.15s;
}

.legend-item:hover {
  background: #f1f5f9;
}

.legend-item.is-selected {
  background: #eef2ff;
  border: 1px solid #c7d2fe;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 9999px;
}

.legend-text {
  font-size: 0.78rem;
  font-weight: 600;
  color: #334155;
}

/* Modal Styling */
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999;
  padding: 20px;
}

.modal-container {
  background: #ffffff;
  width: 100%;
  max-width: 480px;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.15);
  animation: modalPop 0.2s ease-out;
}

@keyframes modalPop {
  from { opacity: 0; transform: scale(0.95); }
  to { opacity: 1; transform: scale(1); }
}

.modal-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f1f5f9;
}

.modal-channel-pill {
  padding: 4px 10px;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 800;
}

.modal-code {
  font-size: 0.8rem;
  font-weight: 700;
  color: #64748b;
  margin-left: 8px;
}

.modal-close-btn {
  background: transparent;
  border: none;
  font-size: 1.5rem;
  color: #94a3b8;
  cursor: pointer;
}

.modal-body {
  padding: 20px;
}

.guest-title {
  font-size: 1.25rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.guest-email {
  font-size: 0.85rem;
  color: #64748b;
  margin: 2px 0 0;
}

.detail-divider {
  height: 1px;
  background: #f1f5f9;
  margin: 16px 0;
}

.modal-grid-info {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 14px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.info-label {
  font-size: 0.72rem;
  font-weight: 600;
  color: #94a3b8;
  text-transform: uppercase;
}

.info-val {
  font-size: 0.88rem;
  color: #1e293b;
}

.status-chip {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: capitalize;
  width: fit-content;
}

.status-confirmed, .status-approved { background: #dcfce7; color: #15803d; }
.status-checked-in { background: #e0e7ff; color: #4338ca; }
.status-pending { background: #fef3c7; color: #b45309; }
.status-cancel, .status-rejected { background: #fee2e2; color: #b91c1c; }

.price-highlight-box {
  background: #f8fafc;
  border-radius: 12px;
  padding: 14px 16px;
  margin-top: 18px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  border: 1px solid #e2e8f0;
}

.price-row {
  display: flex;
  justify-content: space-between;
  font-size: 0.85rem;
  color: #64748b;
}

.price-row.total-row {
  font-size: 1rem;
  font-weight: 800;
  color: #0f172a;
  padding-top: 6px;
  border-top: 1px dashed #cbd5e1;
}

.total-amount {
  color: #4f46e5;
}

.modal-footer {
  padding: 14px 20px;
  background: #f8fafc;
  border-top: 1px solid #f1f5f9;
  display: flex;
  justify-content: flex-end;
}

.btn-close-modal {
  padding: 8px 18px;
  font-size: 0.85rem;
  font-weight: 700;
  border-radius: 10px;
  border: 1px solid #cbd5e1;
  background: #ffffff;
  color: #334155;
  cursor: pointer;
}

.btn-close-modal:hover {
  background: #f1f5f9;
}
</style>
