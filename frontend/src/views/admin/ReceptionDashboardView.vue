<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import VueApexCharts from 'vue3-apexcharts'
import type { ApexOptions } from 'apexcharts'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import { reservationService } from '@/services/admin/reservationService'
import { roomService } from '@/services/roomService'
import { channelService } from '@/services/admin/channelService'
import { useAuthStore } from '@/stores/authStore'
import { useToastStore } from '@/stores/toastStore'
import type { Reservation, Room, Channel } from '@/types'
import {
  Calendar as CalendarIcon,
  ChevronRight,
  TrendingUp,
  DollarSign,
  Users,
  BedDouble,
  CheckCircle2,
  Clock,
  ArrowUpRight,
  Sparkles,
  BarChart3,
  Percent,
  Download,
  FileSpreadsheet,
  Building2,
  Layers,
  Activity,
  Compass,
  Radio,
  ExternalLink,
  SlidersHorizontal,
  Plus,
  RefreshCw,
} from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()
const toastStore = useToastStore()

// State from Database
const loading = ref(false)
const reservations = ref<Reservation[]>([])
const rooms = ref<Room[]>([])
const channels = ref<Channel[]>([])

// Dashboard Active Tab
const activeTab = ref<'all' | 'operations' | 'revenue' | 'channels'>('all')

// Month & Year Filter for Revenue Analytics
const selectedMonth = ref(new Date().getMonth()) // 0-indexed
const selectedYear = ref(new Date().getFullYear())
const selectedProperty = ref('all')
const selectedChannel = ref('all')

const monthsList = [
  { value: 0, label: 'Januari' },
  { value: 1, label: 'Februari' },
  { value: 2, label: 'Maret' },
  { value: 3, label: 'April' },
  { value: 4, label: 'Mei' },
  { value: 5, label: 'Juni' },
  { value: 6, label: 'Juli' },
  { value: 7, label: 'Agustus' },
  { value: 8, label: 'September' },
  { value: 9, label: 'Oktober' },
  { value: 10, label: 'November' },
  { value: 11, label: 'Desember' },
]

const yearsList = [2024, 2025, 2026, 2027]

// Greeting
const now = new Date()
const hour = now.getHours()
const greeting = computed(() => {
  if (hour < 12) return 'Selamat Pagi'
  if (hour < 15) return 'Selamat Siang'
  if (hour < 18) return 'Selamat Sore'
  return 'Selamat Malam'
})

const userName = computed(() => {
  return authStore.user?.first_name || 'Admin'
})

const dateSubtitle = computed(() => {
  const options: Intl.DateTimeFormatOptions = {
    weekday: 'long',
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  }
  return now.toLocaleDateString('id-ID', options)
})

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

const formatYMD = (d: Date): string => {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

const todayStr = formatYMD(now)

const formatIDR = (val: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(val || 0)
}

// Dynamic Channel Color Palette
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

// Identify channel directly from Database
const getBookingChannelName = (res: Reservation): string => {
  if (res.channel?.name) return res.channel.name
  if (res.channel_id) {
    const ch = channels.value.find((c) => c.id === res.channel_id)
    if (ch) return ch.name
  }
  return 'Direct Website'
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
  }
}

const getInitials = (name?: string) => {
  if (!name) return '?'
  const parts = name.trim().split(' ')
  if (parts.length >= 2 && parts[0] && parts[1]) {
    const f = parts[0][0] || ''
    const s = parts[1][0] || ''
    return (f + s).toUpperCase()
  }
  return (name[0] || '?').toUpperCase()
}

// Effective Rooms directly from DB
const effectiveRooms = computed(() => {
  if (selectedProperty.value === 'all') return rooms.value
  return rooms.value.filter((r) => r.id === selectedProperty.value)
})

// --- 1. OPERATIONS DASHBOARD COMPUTATIONS (100% PURE DATABASE DATA) ---
const todayArrivals = computed(() =>
  reservations.value.filter(
    (r) =>
      parseDateToISO(r.checkin_date) === todayStr &&
      ['confirmed', 'approved', 'pending'].includes(r.status),
  ),
)

const todayDepartures = computed(() =>
  reservations.value.filter(
    (r) =>
      parseDateToISO(r.checkout_date) === todayStr &&
      ['checked-in', 'confirmed', 'approved'].includes(r.status),
  ),
)

const inHouseTonight = computed(() =>
  reservations.value.filter((r) => r.status === 'checked-in'),
)

const thisWeekArrivals = computed(() => {
  const next7Days = new Date()
  next7Days.setDate(next7Days.getDate() + 7)
  const next7Str = formatYMD(next7Days)
  return reservations.value.filter(
    (r) => {
      const checkinISO = parseDateToISO(r.checkin_date)
      return (
        checkinISO >= todayStr &&
        checkinISO <= next7Str &&
        !['cancel', 'rejected'].includes(r.status)
      )
    },
  )
})

const upcomingArrivalsGrouped = computed(() => {
  const next60 = new Date()
  next60.setDate(next60.getDate() + 60)
  const maxDateStr = formatYMD(next60)

  const upcoming = reservations.value.filter(
    (r) => {
      const checkinISO = parseDateToISO(r.checkin_date)
      const notCancelled = !['cancel', 'rejected'].includes(r.status)
      const inDateRange = checkinISO >= todayStr && checkinISO <= maxDateStr
      const matchesChannel = selectedChannel.value === 'all' || getBookingChannelKey(r) === selectedChannel.value
      return notCancelled && inDateRange && matchesChannel
    },
  )

  const groups: Record<string, Reservation[]> = {}
  upcoming.forEach((r) => {
    const d = parseDateToISO(r.checkin_date)
    if (!groups[d]) groups[d] = []
    groups[d]!.push(r)
  })

  const sortedDates = Object.keys(groups).sort()
  return sortedDates.map((dateStr) => {
    const d = new Date(dateStr)
    const tomorrow = new Date()
    tomorrow.setDate(tomorrow.getDate() + 1)
    const isTomorrow = dateStr === formatYMD(tomorrow)

    let title = isTomorrow
      ? 'Besok'
      : d.toLocaleDateString('id-ID', { weekday: 'short', month: 'short', day: 'numeric' })

    const items = groups[dateStr] || []
    return {
      dateStr,
      title: `${title} · ${items.length} ${items.length > 1 ? 'kedatangan' : 'kedatangan'}`,
      items,
    }
  })
})

// --- 2. MONTHLY REVENUE & PERFORMANCE ANALYTICS (100% PURE DATABASE DATA) ---
const daysInSelectedMonth = computed(() => {
  return new Date(selectedYear.value, selectedMonth.value + 1, 0).getDate()
})

const monthTitle = computed(() => {
  const m = monthsList[selectedMonth.value]?.label || ''
  return `${m} ${selectedYear.value}`
})

const monthReservations = computed(() => {
  const currentYM = `${selectedYear.value}-${String(selectedMonth.value + 1).padStart(2, '0')}`
  return reservations.value.filter((r) => {
    const checkinISO = parseDateToISO(r.checkin_date)
    const matchesMonth = checkinISO.startsWith(currentYM)
    const matchesProperty = selectedProperty.value === 'all' || r.room_id === selectedProperty.value
    const matchesChannel = selectedChannel.value === 'all' || getBookingChannelKey(r) === selectedChannel.value
    
    // STRICT: Hanya transaksi yang sudah disetujui tim Finance/Admin ('approved', 'confirmed', 'checked-in', 'checked-out')
    // Status 'pending' TIDAK dimasukkan ke transaksi pendapatan/omzet sebelum disetujui.
    const isApprovedTransaction = ['approved', 'confirmed', 'checked-in', 'checked-out'].includes(
      (r.status || '').toLowerCase(),
    )
    return matchesMonth && matchesProperty && matchesChannel && isApprovedTransaction
  })
})

const totalReservations = computed(() => monthReservations.value.length)
const totalNights = computed(() => monthReservations.value.reduce((acc, r) => acc + (r.total_night || 1), 0))
const totalRevenue = computed(() => monthReservations.value.reduce((acc, r) => acc + (r.total_price || 0), 0))
const totalRoomsCount = computed(() => Math.max(1, rooms.value.length))

const occupancyRate = computed(() => {
  const possible = totalRoomsCount.value * daysInSelectedMonth.value
  if (possible === 0) return 0
  return Math.round((totalNights.value / possible) * 1000) / 10
})

const adr = computed(() => {
  if (totalNights.value === 0) return 0
  return Math.round(totalRevenue.value / totalNights.value)
})

const revPAR = computed(() => {
  const possible = daysInSelectedMonth.value * totalRoomsCount.value
  if (possible === 0) return 0
  return Math.round(totalRevenue.value / possible)
})

// --- 3. CHANNEL PERFORMANCE DISTRIBUTION (100% PURE DATABASE DATA) ---
const channelBreakdown = computed(() => {
  const map: Record<string, { key: string; name: string; reservations: number; revenue: number; nights: number; color: string }> = {}

  // 1. Register all active channels from database
  channels.value.forEach((ch, idx) => {
    map[ch.id] = {
      key: ch.id,
      name: ch.name,
      reservations: 0,
      revenue: 0,
      nights: 0,
      color: getChannelColor(ch.name, idx),
    }
  })

  // Ensure default Direct Website channel exists
  const hasDirect = Object.values(map).some((c) => c.name.toLowerCase().includes('direct'))
  if (!hasDirect) {
    map['direct'] = {
      key: 'direct',
      name: 'Direct Website',
      reservations: 0,
      revenue: 0,
      nights: 0,
      color: '#16a34a',
    }
  }

  // 2. Aggregate actual reservations from DB
  monthReservations.value.forEach((r) => {
    let targetKey = ''
    if (r.channel_id && map[r.channel_id]) {
      targetKey = r.channel_id
    } else if (r.channel?.id && map[r.channel.id]) {
      targetKey = r.channel.id
    } else if (r.channel?.name) {
      const match = Object.values(map).find((c) => c.name.toLowerCase() === r.channel?.name.toLowerCase())
      if (match) targetKey = match.key
    }

    if (!targetKey) {
      const direct = Object.values(map).find((c) => c.name.toLowerCase().includes('direct'))
      targetKey = direct ? direct.key : 'direct'
    }

    if (map[targetKey]) {
      map[targetKey]!.reservations += 1
      map[targetKey]!.revenue += r.total_price || 0
      map[targetKey]!.nights += r.total_night || 1
    }
  })

  const revTotal = Math.max(1, totalRevenue.value)
  return Object.values(map)
    .map((c) => ({
      ...c,
      percent: Math.round((c.revenue / revTotal) * 1000) / 10,
    }))
    .sort((a, b) => b.revenue - a.revenue)
})

// Active channels with revenue > 0 for donut chart slices
const activeChannelsWithRevenue = computed(() => {
  return channelBreakdown.value.filter((c) => c.revenue > 0)
})

// Cancelled/Rejected Reservations in Selected Month (from DB)
const monthCancelledReservations = computed(() => {
  const currentYM = `${selectedYear.value}-${String(selectedMonth.value + 1).padStart(2, '0')}`
  return reservations.value.filter((r) => {
    const checkinISO = parseDateToISO(r.checkin_date)
    const matchesMonth = checkinISO.startsWith(currentYM)
    const matchesProperty = selectedProperty.value === 'all' || r.room_id === selectedProperty.value
    const isCancelled = r.status === 'cancel' || r.status === 'rejected'
    return matchesMonth && matchesProperty && isCancelled
  })
})

const totalCancelledReservations = computed(() => monthCancelledReservations.value.length)

// By Property Table Rows (100% from DB rooms & reservations)
const propertyRows = computed(() => {
  return rooms.value.map((room) => {
    const roomRes = monthReservations.value.filter((r) => r.room_id === room.id)
    const roomCancelled = monthCancelledReservations.value.filter((r) => r.room_id === room.id).length
    const reservationsCount = roomRes.length
    const nights = roomRes.reduce((acc, r) => acc + (r.total_night || 1), 0)
    const rev = roomRes.reduce((acc, r) => acc + (r.total_price || 0), 0)
    const occ = Math.round((nights / daysInSelectedMonth.value) * 1000) / 10
    const roomAdr = nights > 0 ? Math.round(rev / nights) : 0
    const roomRevPar = Math.round(rev / daysInSelectedMonth.value)

    return {
      id: room.id,
      name: room.name,
      reservations: reservationsCount,
      cancelled: roomCancelled,
      nights,
      occupancy: occ,
      revenue: rev,
      adr: roomAdr,
      revpar: roomRevPar,
    }
  })
})

// Daily Revenue & Occupancy Trend Chart Data (from DB)
const dailyTrendData = computed(() => {
  const days = daysInSelectedMonth.value
  const dailyRev: number[] = new Array(days).fill(0)

  monthReservations.value.forEach((r) => {
    const checkinISO = parseDateToISO(r.checkin_date)
    const checkoutISO = parseDateToISO(r.checkout_date)
    if (!checkinISO) return

    const startD = new Date(checkinISO)
    const endD = checkoutISO ? new Date(checkoutISO) : startD
    const pricePerNight = (r.total_price || 0) / Math.max(1, r.total_night || 1)

    for (let d = 1; d <= days; d++) {
      const cur = new Date(selectedYear.value, selectedMonth.value, d)
      if (cur >= startD && cur < endD) {
        dailyRev[d - 1] = (dailyRev[d - 1] || 0) + pricePerNight
      }
    }
  })

  return {
    categories: Array.from({ length: days }, (_, i) => String(i + 1)),
    revenueSeries: dailyRev,
  }
})

// ApexCharts Options for Modern Rounded Bar / Column Chart
const revenueChartOptions = computed<ApexOptions>(() => ({
  chart: {
    type: 'bar',
    height: 260,
    toolbar: { show: false },
    zoom: { enabled: false },
    fontFamily: 'Plus Jakarta Sans, sans-serif',
  },
  colors: ['#4f46e5'],
  plotOptions: {
    bar: {
      borderRadius: 5,
      borderRadiusApplication: 'end',
      columnWidth: '55%',
      colors: {
        ranges: [
          {
            from: 1,
            to: 1000000000,
            color: '#4f46e5',
          },
          {
            from: 0,
            to: 0,
            color: '#e2e8f0',
          },
        ],
      },
    },
  },
  dataLabels: { enabled: false },
  xaxis: {
    categories: dailyTrendData.value.categories,
    labels: {
      style: { colors: '#64748b', fontSize: '11px', fontWeight: 700 },
      rotate: 0,
      hideOverlappingLabels: true,
    },
    axisBorder: { show: false },
    axisTicks: { show: false },
    title: {
      text: `Tanggal (1 - ${daysInSelectedMonth.value} ${monthTitle.value})`,
      style: { color: '#94a3b8', fontSize: '11px', fontWeight: 600 },
      offsetY: 8,
    },
  },
  yaxis: {
    labels: {
      style: { colors: '#64748b', fontSize: '11px', fontWeight: 700 },
      formatter: (val: number) => {
        if (val === 0) return 'Rp 0'
        if (val >= 1000000) return `Rp ${(val / 1000000).toFixed(1)}jt`
        return `Rp ${Math.round(val / 1000)}k`
      },
    },
  },
  grid: {
    borderColor: '#f1f5f9',
    strokeDashArray: 4,
    yaxis: { lines: { show: true } },
    xaxis: { lines: { show: false } },
  },
  tooltip: {
    theme: 'dark',
    custom: function ({ series, seriesIndex, dataPointIndex }: any) {
      const day = dataPointIndex + 1
      const val = series[seriesIndex][dataPointIndex]
      const formatted = formatIDR(val)
      return `
        <div style="padding: 10px 14px; font-family: Plus Jakarta Sans, sans-serif; font-size: 12px; background: #0f172a; color: #fff; border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.3);">
          <div style="font-weight: 800; color: #94a3b8; margin-bottom: 4px;">Tanggal ${day} ${monthTitle.value}</div>
          <div style="font-size: 14px; font-weight: 800; color: #38bdf8;">${formatted}</div>
          <div style="font-size: 11px; color: ${val > 0 ? '#4ade80' : '#94a3b8'}; margin-top: 4px;">
            ${val > 0 ? '✓ Ada Pemesanan Menginap' : 'Tidak ada booking aktif'}
          </div>
        </div>
      `
    },
  },
}))

const revenueChartSeries = computed(() => [
  {
    name: 'Pendapatan Database Harian',
    data: dailyTrendData.value.revenueSeries,
  },
])

const channelChartOptions = computed<ApexOptions>(() => ({
  chart: {
    type: 'donut',
    height: 220,
    fontFamily: 'Plus Jakarta Sans, sans-serif',
  },
  labels: activeChannelsWithRevenue.value.length > 0 ? activeChannelsWithRevenue.value.map((c) => c.name) : ['Belum ada data'],
  colors: activeChannelsWithRevenue.value.length > 0 ? activeChannelsWithRevenue.value.map((c) => c.color) : ['#cbd5e1'],
  stroke: { width: 3, colors: ['#ffffff'] },
  legend: {
    show: false,
  },
  dataLabels: {
    enabled: true,
    formatter: (val: number) => `${Math.round(val)}%`,
    style: { fontSize: '11px', fontWeight: 700 },
    dropShadow: { enabled: false },
  },
  plotOptions: {
    pie: {
      donut: {
        size: '72%',
        labels: {
          show: true,
          total: {
            show: true,
            label: 'Total Revenue',
            fontSize: '10px',
            fontWeight: 700,
            color: '#64748b',
            formatter: () => formatIDR(totalRevenue.value),
          },
        },
      },
    },
  },
  tooltip: {
    theme: 'dark',
    y: {
      formatter: (val: number) => `${val}% (${formatIDR((val / 100) * totalRevenue.value)})`,
    },
  },
}))

const channelChartSeries = computed(() => {
  if (activeChannelsWithRevenue.value.length === 0) return [100]
  return activeChannelsWithRevenue.value.map((c) => c.percent)
})

// Daily Occupancy Heatmap
const dailyOccupancyMap = computed(() => {
  const days = daysInSelectedMonth.value
  const result = []
  for (let i = 1; i <= days; i++) {
    const dStr = `${selectedYear.value}-${String(selectedMonth.value + 1).padStart(2, '0')}-${String(i).padStart(2, '0')}`
    const dateObj = new Date(selectedYear.value, selectedMonth.value, i)
    const matching = monthReservations.value.filter((r) => {
      const inISO = parseDateToISO(r.checkin_date)
      const outISO = parseDateToISO(r.checkout_date)
      if (!inISO || !outISO) return false
      return dStr >= inISO && dStr < outISO
    })
    const occ = Math.round((matching.length / totalRoomsCount.value) * 100)
    result.push({
      day: i,
      dayName: ['Min', 'Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab'][dateObj.getDay()],
      dateStr: dStr,
      occupied: matching.length,
      occupancyPercent: occ,
    })
  }
  return result
})

// Navigation Actions
const goToReservations = () => router.push('/admin/reservations')
const goToCalendar = () => router.push('/admin/calendar')
const goToChannels = () => router.push('/admin/channels')

let pollingTimer: any = null

const refreshData = async (isSilent: boolean = false) => {
  try {
    if (!isSilent) {
      loading.value = true
    }
    const [resList, rList, cList] = await Promise.all([
      reservationService.getAll(),
      roomService.getAll().catch(() => []),
      channelService.getAll().catch(() => []),
    ])
    reservations.value = resList
    rooms.value = rList
    channels.value = cList
  } catch (error: any) {
    if (!isSilent) {
      toastStore.error(error.message || 'Gagal memuat data dari database')
    }
  } finally {
    if (!isSilent) {
      loading.value = false
    }
  }
}

const exportPDF = () => {
  try {
    const doc = new jsPDF('p', 'mm', 'a4')
    const period = `${monthsList[selectedMonth.value]?.label} ${selectedYear.value}`
    const timestamp = new Date().toLocaleString('id-ID')

    // Header Title
    doc.setFontSize(16)
    doc.setTextColor(15, 23, 42)
    doc.text('CM LIVING HOTEL & SUITES', 14, 18)

    doc.setFontSize(10)
    doc.setTextColor(71, 85, 105)
    doc.text(`Laporan Keuangan & Performa Reservasi · Periode: ${period}`, 14, 25)
    doc.setFontSize(8)
    doc.setTextColor(148, 163, 184)
    doc.text(`Waktu Ekspor: ${timestamp} · Data Database Asli`, 14, 30)

    // 1. KPI Summary Table
    autoTable(doc, {
      startY: 34,
      head: [['Total Booking', 'Malam Terjual', 'Occupancy (%)', 'Total Revenue (IDR)', 'ADR (IDR)', 'RevPAR (IDR)']],
      body: [[
        totalReservations.value.toString(),
        totalNights.value.toString(),
        `${occupancyRate.value}%`,
        formatIDR(totalRevenue.value),
        formatIDR(adr.value),
        formatIDR(revPAR.value),
      ]],
      headStyles: { fillColor: [79, 70, 229] },
      theme: 'grid',
      styles: { fontSize: 8.5, halign: 'center' },
    })

    // 2. Property Performance Table
    const propRowsData: any[] = propertyRows.value.map((p) => [
      p.name,
      p.reservations.toString(),
      p.cancelled.toString(),
      p.nights.toString(),
      `${p.occupancy}%`,
      formatIDR(p.revenue),
      formatIDR(p.adr),
      formatIDR(p.revpar),
    ])
    propRowsData.push([
      'TOTAL KESELURUHAN',
      totalReservations.value.toString(),
      totalCancelledReservations.value.toString(),
      totalNights.value.toString(),
      `${occupancyRate.value}%`,
      formatIDR(totalRevenue.value),
      formatIDR(adr.value),
      formatIDR(revPAR.value),
    ])

    autoTable(doc, {
      startY: (doc as any).lastAutoTable.finalY + 8,
      head: [['Unit Kamar', 'Booking', 'Dibatalkan', 'Malam', 'Okupansi', 'Total Revenue (IDR)', 'ADR (IDR)', 'RevPAR (IDR)']],
      body: propRowsData.length > 0 ? propRowsData : [['Tidak ada kamar di database', '-', '-', '-', '-', '-', '-', '-']],
      headStyles: { fillColor: [15, 23, 42] },
      theme: 'striped',
      styles: { fontSize: 8 },
    })

    // 3. Channel Distribution Table
    const channelRowsData: any[] = channelBreakdown.value.map((c) => [
      c.name,
      c.reservations.toString(),
      c.nights.toString(),
      `${c.percent}%`,
      formatIDR(c.revenue),
    ])

    autoTable(doc, {
      startY: (doc as any).lastAutoTable.finalY + 8,
      head: [['Channel OTA / Direct', 'Jumlah Booking', 'Malam Terjual', 'Pangsa Pasar (%)', 'Total Pendapatan (IDR)']],
      body: channelRowsData.length > 0 ? channelRowsData : [['Belum ada transaksi channel di database', '-', '-', '-', '-']],
      headStyles: { fillColor: [16, 185, 129] },
      theme: 'striped',
      styles: { fontSize: 8 },
    })

    // 4. Detail Transaksi Reservasi
    const resDetailRows: any[] = monthReservations.value.map((r) => [
      r.code || '-',
      r.full_name || '-',
      r.room?.name || 'Unit Kamar',
      r.checkin_date || '-',
      r.checkout_date || '-',
      `${r.total_night || 1} mlm`,
      getBookingChannelKey(r).toUpperCase(),
      r.status,
      formatIDR(r.total_price || 0),
    ])

    if (resDetailRows.length > 0) {
      autoTable(doc, {
        startY: (doc as any).lastAutoTable.finalY + 8,
        head: [['Kode', 'Nama Tamu', 'Kamar', 'Check-In', 'Check-Out', 'Durasi', 'Channel', 'Status', 'Total Biaya']],
        body: resDetailRows,
        headStyles: { fillColor: [67, 56, 202] },
        theme: 'striped',
        styles: { fontSize: 7.5 },
      })
    }

    doc.save(`Laporan_Keuangan_CMLiving_${selectedYear.value}_${selectedMonth.value + 1}.pdf`)
    toastStore.success('Dokumen PDF data performa & transaksi berhasil diunduh!')
  } catch (error: any) {
    toastStore.error(error.message || 'Gagal mengekspor dokumen PDF')
  }
}

const exportCSV = () => {
  try {
    const period = `${monthsList[selectedMonth.value]?.label} ${selectedYear.value}`
    const timestamp = new Date().toLocaleString('id-ID')
    const lines: string[] = []

    // 1. Header Metadata
    lines.push(`LAPORAN KEUANGAN & PERFORMA RESERVASI CM LIVING`)
    lines.push(`Periode: ${period}`)
    lines.push(`Waktu Ekspor: ${timestamp}`)
    lines.push(`Total Unit Kamar: ${rooms.value.length}`)
    lines.push(``)

    // 2. Summary KPI Metrics
    lines.push(`=== RINGKASAN METRIK UTAMA ===`)
    lines.push(`Total Booking,Total Malam Terjual,Tingkat Okupansi (%),Total Pendapatan (IDR),ADR (IDR),RevPAR (IDR)`)
    lines.push(`${totalReservations.value},${totalNights.value},${occupancyRate.value}%,${totalRevenue.value},${adr.value},${revPAR.value}`)
    lines.push(``)

    // 3. Property Performance Breakdown
    lines.push(`=== RINCIAN PERFORMA PER UNIT KAMAR ===`)
    lines.push(`Nama Kamar,Jumlah Reservasi,Dibatalkan,Total Malam,Occupancy Rate (%),Total Pendapatan (IDR),ADR (IDR),RevPAR (IDR)`)
    propertyRows.value.forEach((p) => {
      lines.push(`"${p.name}",${p.reservations},${p.cancelled},${p.nights},${p.occupancy}%,${p.revenue},${p.adr},${p.revpar}`)
    })
    lines.push(`"TOTAL KESELURUHAN",${totalReservations.value},${totalCancelledReservations.value},${totalNights.value},${occupancyRate.value}%,${totalRevenue.value},${adr.value},${revPAR.value}`)
    lines.push(``)

    // 4. Channel Breakdown
    lines.push(`=== DISTRIBUSI CHANNEL PEMESANAN ===`)
    lines.push(`Channel OTA / Direct,Jumlah Booking,Total Malam,Pangsa Pasar (%),Total Pendapatan (IDR)`)
    channelBreakdown.value.forEach((c) => {
      lines.push(`"${c.name}",${c.reservations},${c.nights},${c.percent}%,${c.revenue}`)
    })
    lines.push(``)

    // 5. Individual Reservations Detail
    lines.push(`=== DAFTAR TRANSAKSI RESERVASI (${period}) ===`)
    lines.push(`Kode Booking,Nama Tamu,Kamar,Check-In,Check-Out,Malam,Tamu,Channel,Status,Total Harga (IDR)`)
    monthReservations.value.forEach((r) => {
      lines.push(`"${r.code || '-'}","${r.full_name || '-'}","${r.room?.name || '-'}","${r.checkin_date || '-'}","${r.checkout_date || '-'}","${r.total_night || 1}","${(r.number_of_adult || 1) + (r.number_of_children || 0)}","${getBookingChannelKey(r)}","${r.status}","${r.total_price || 0}"`)
    })

    // UTF-8 BOM so Excel opens with proper Indonesian formatting
    const csvContent = '\uFEFF' + lines.join('\r\n')
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.setAttribute('href', url)
    link.setAttribute('download', `Laporan_Performa_CMLiving_${selectedYear.value}_${selectedMonth.value + 1}.csv`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    toastStore.success('File CSV data performa & transaksi berhasil diunduh!')
  } catch (error: any) {
    toastStore.error(error.message || 'Gagal mengekspor file CSV')
  }
}

onMounted(() => {
  refreshData(false)
  // Silent auto-reload polling every 5 seconds without triggering loading state or screen flickering
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
  <div class="master-dashboard">
    <!-- Top Header Card -->
    <div class="top-header-card">
      <div class="greeting-left">
        <div class="greeting-badge">
          <Sparkles class="w-4 h-4 text-emerald-600" />
          <span>Live Database Connected · {{ reservations.length }} Total Reservasi di DB</span>
        </div>
        <h1 class="greeting-title">{{ greeting }}, {{ userName }}</h1>
        <p class="greeting-subtitle">{{ dateSubtitle }} · {{ rooms.length }} Unit Kamar Terdaftar di Database</p>
      </div>

      <!-- Quick Global Controls & Tab Switcher -->
      <div class="header-right-controls">
        <button @click="() => refreshData(false)" :disabled="loading" class="btn-refresh-db" title="Segarkan Data Database">
          <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
          <span>{{ loading ? 'Memuat...' : 'Sync DB' }}</span>
        </button>

        <div class="tab-switcher">
          <button
            class="tab-btn"
            :class="{ active: activeTab === 'all' }"
            @click="activeTab = 'all'"
          >
            <Compass class="w-4 h-4" />
            <span>Semua</span>
          </button>
          <button
            class="tab-btn"
            :class="{ active: activeTab === 'operations' }"
            @click="activeTab = 'operations'"
          >
            <Activity class="w-4 h-4" />
            <span>Operasional</span>
          </button>
          <button
            class="tab-btn"
            :class="{ active: activeTab === 'revenue' }"
            @click="activeTab = 'revenue'"
          >
            <BarChart3 class="w-4 h-4" />
            <span>Revenue (IDR)</span>
          </button>
          <button
            class="tab-btn"
            :class="{ active: activeTab === 'channels' }"
            @click="activeTab = 'channels'"
          >
            <Radio class="w-4 h-4" />
            <span>Channel OTA</span>
          </button>
        </div>

        <button @click="goToChannels" class="btn-channel-manage" title="Kelola Master Channel OTA">
          <Radio class="w-4 h-4" />
          <span>Kelola Master Channel</span>
        </button>
      </div>
    </div>

    <!-- ============================================== -->
    <!-- SECTION 1: OPERATIONS DASHBOARD (Daily Stays)  -->
    <!-- ============================================== -->
    <div v-if="activeTab === 'all' || activeTab === 'operations'" class="operations-section">
      <div class="section-title-wrap">
        <div class="section-icon-badge">
          <Activity class="w-4 h-4 text-indigo-600" />
        </div>
        <div>
          <h2 class="section-title">Operasional Resepsionis & Frontdesk Hari Ini</h2>
          <p class="section-desc">Data live dari tabel reservasi database CM Living</p>
        </div>
      </div>

      <!-- 5 KPI Cards -->
      <div class="kpi-5-grid">
        <div class="kpi-card" @click="goToReservations">
          <div class="kpi-top">
            <span class="kpi-label">Arrivals today</span>
            <Users class="w-4 h-4 text-blue-600" />
          </div>
          <span class="kpi-val">{{ todayArrivals.length }} <span class="kpi-unit">tamu</span></span>
        </div>

        <div class="kpi-card" @click="goToReservations">
          <div class="kpi-top">
            <span class="kpi-label">Departures today</span>
            <CheckCircle2 class="w-4 h-4 text-rose-600" />
          </div>
          <span class="kpi-val">{{ todayDepartures.length }} <span class="kpi-unit">tamu</span></span>
        </div>

        <div class="kpi-card" @click="goToReservations">
          <div class="kpi-top">
            <span class="kpi-label">In house tonight</span>
            <BedDouble class="w-4 h-4 text-emerald-600" />
          </div>
          <span class="kpi-val">{{ inHouseTonight.length }} <span class="kpi-unit">kamar</span></span>
        </div>

        <div class="kpi-card" @click="goToReservations">
          <div class="kpi-top">
            <span class="kpi-label">Arrivals this week</span>
            <CalendarIcon class="w-4 h-4 text-amber-600" />
          </div>
          <span class="kpi-val">{{ thisWeekArrivals.length }} <span class="kpi-unit">tamu</span></span>
        </div>

        <div class="kpi-card" @click="goToCalendar">
          <div class="kpi-top">
            <span class="kpi-label">Properties / Units</span>
            <Building2 class="w-4 h-4 text-purple-600" />
          </div>
          <span class="kpi-val">{{ rooms.length }} <span class="kpi-unit">kamar</span></span>
        </div>
      </div>

      <!-- 2x2 Operations Grid -->
      <div class="operations-2x2-grid">
        <div class="op-panel">
          <div class="panel-head">
            <h3>Arrivals · today</h3>
            <button @click="goToReservations" class="view-all-link">Lihat semua &rarr;</button>
          </div>
          <div v-if="todayArrivals.length > 0" class="panel-content">
            <div v-for="item in todayArrivals" :key="item.id" class="arrival-row">
              <div class="guest-info">
                <span class="guest-name">{{ item.full_name }}</span>
                <span class="unit-name">{{ item.room?.name || 'Unit Kamar' }} · {{ (item.number_of_adult || 1) + (item.number_of_children || 0) }} tamu</span>
              </div>
              <span
                class="channel-pill"
                :style="{
                  backgroundColor: getChannelStyle(getBookingChannelKey(item)).bg,
                  color: getChannelStyle(getBookingChannelKey(item)).text,
                }"
              >
                {{ getChannelStyle(getBookingChannelKey(item)).label }}
              </span>
            </div>
          </div>
          <div v-else class="empty-panel">
            <span>Tidak ada jadwal kedatangan tamu hari ini di database</span>
          </div>
        </div>

        <div class="op-panel">
          <div class="panel-head">
            <h3>Departures · today</h3>
            <button @click="goToReservations" class="view-all-link">Lihat semua &rarr;</button>
          </div>
          <div v-if="todayDepartures.length > 0" class="panel-content">
            <div v-for="item in todayDepartures" :key="item.id" class="arrival-row">
              <div class="guest-info">
                <span class="guest-name">{{ item.full_name }}</span>
                <span class="unit-name">{{ item.room?.name || 'Unit Kamar' }}</span>
              </div>
              <span class="badge-checkedout">Check-Out</span>
            </div>
          </div>
          <div v-else class="empty-panel">
            <span>Tidak ada kepulangan tamu hari ini di database</span>
          </div>
        </div>

        <div class="op-panel">
          <div class="panel-head">
            <h3>Open tasks · today</h3>
            <button class="view-all-link">Lihat semua &rarr;</button>
          </div>
          <div class="empty-panel">
            <span>Semua operasional hari ini sudah diselesaikan</span>
          </div>
        </div>

        <div class="op-panel">
          <div class="panel-head">
            <h3>Open tasks · tomorrow</h3>
            <button class="view-all-link">Lihat semua &rarr;</button>
          </div>
          <div class="empty-panel">
            <span>Tidak ada tugas tertunda untuk besok</span>
          </div>
        </div>
      </div>

      <!-- Upcoming Arrivals Next 60 Days -->
      <div class="upcoming-card">
        <div class="upcoming-head">
          <div>
            <h3>Upcoming arrivals · next 60 days</h3>
            <p class="text-xs text-slate-500 mt-1">Daftar tamu yang akan datang dari database (60 hari ke depan)</p>
          </div>
          <div class="flex items-center gap-3">
            <select v-model="selectedChannel" class="custom-select-sm">
              <option value="all">Semua Channel</option>
              <option value="booking">Booking.com</option>
              <option value="airbnb">Airbnb</option>
              <option value="agoda">Agoda</option>
              <option value="traveloka">Traveloka</option>
              <option value="direct">Direct Website</option>
            </select>
            <button @click="goToCalendar" class="btn-calendar-link">
              <CalendarIcon class="w-4 h-4 text-indigo-600" />
              <span>Buka PMS Calendar</span>
            </button>
          </div>
        </div>

        <div v-if="upcomingArrivalsGrouped.length > 0" class="upcoming-body">
          <div
            v-for="group in upcomingArrivalsGrouped"
            :key="group.dateStr"
            class="arrival-group"
          >
            <div class="group-header-row">
              <span class="group-title">{{ group.title }}</span>
            </div>

            <div class="group-items-list">
              <div
                v-for="item in group.items"
                :key="item.id"
                class="guest-arrival-card"
                @click="goToReservations"
              >
                <div class="guest-left">
                  <div class="guest-avatar">
                    {{ getInitials(item.full_name) }}
                  </div>
                  <div class="guest-meta">
                    <span class="guest-name">{{ item.full_name }}</span>
                    <span class="stay-info">
                      {{ item.room?.name || 'Unit Kamar' }} · {{ (item.number_of_adult || 1) + (item.number_of_children || 0) }} guests · {{ item.total_night || 1 }} nights
                    </span>
                  </div>
                </div>

                <div class="guest-right">
                  <span
                    class="channel-pill"
                    :style="{
                      backgroundColor: getChannelStyle(getBookingChannelKey(item)).bg,
                      color: getChannelStyle(getBookingChannelKey(item)).text,
                    }"
                  >
                    {{ getChannelStyle(getBookingChannelKey(item)).label }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="empty-state-box">
          <CalendarIcon class="w-10 h-10 text-slate-300 mx-auto mb-2" />
          <p class="font-bold text-slate-700">Belum ada jadwal kedatangan tamu di database</p>
          <p class="text-xs text-slate-400 mt-1">Gunakan menu 'Reservation' &rarr; 'New Reservation' untuk menambahkan pemesanan baru</p>
        </div>
      </div>
    </div>

    <!-- ======================================================== -->
    <!-- SECTION 2: DISTRIBUTION CHANNEL PERFORMANCE WIDGET       -->
    <!-- ======================================================== -->
    <div v-if="activeTab === 'all' || activeTab === 'channels'" class="channel-widget-section">
      <div class="channel-widget-card">
        <div class="widget-head">
          <div class="section-title-wrap">
            <div class="section-icon-badge bg-rose-50 text-rose-600">
              <Radio class="w-4 h-4" />
            </div>
            <div>
              <h2 class="section-title">Performa Distribusi Channel OTA & Direct</h2>
              <p class="section-desc">Pangsa pasar, total pemesanan, dan pendapatan per channel dari database</p>
            </div>
          </div>

          <button @click="goToChannels" class="btn-manage-accent">
            <Radio class="w-4 h-4" />
            <span>Kelola Master Channel (/admin/channels)</span>
          </button>
        </div>

        <div v-if="channelBreakdown.length > 0" class="channel-grid-cards">
          <div
            v-for="ch in channelBreakdown"
            :key="ch.key"
            class="channel-stat-box"
            :style="{ borderLeftColor: ch.color }"
          >
            <div class="ch-box-top">
              <span class="ch-badge" :style="{ backgroundColor: ch.color }">{{ ch.name }}</span>
              <span class="ch-share-pct">{{ ch.percent }}% Share</span>
            </div>
            <div class="ch-box-revenue">
              <span class="ch-rev-val">{{ formatIDR(ch.revenue) }}</span>
              <span class="ch-sub-info">{{ ch.reservations }} Booking · {{ ch.nights }} Malam</span>
            </div>
            <div class="ch-progress-bg">
              <div class="ch-progress-fill" :style="{ width: `${ch.percent}%`, backgroundColor: ch.color }"></div>
            </div>
          </div>
        </div>
        <div v-else class="empty-state-box">
          <Radio class="w-8 h-8 text-slate-300 mx-auto mb-2" />
          <p class="text-sm font-semibold text-slate-600">Belum ada transaksi channel untuk periode yang dipilih</p>
        </div>
      </div>
    </div>

    <!-- ======================================================== -->
    <!-- SECTION 3: MONTHLY REVENUE & PERFORMANCE REPORT (IDR)    -->
    <!-- ======================================================== -->
    <div v-if="activeTab === 'all' || activeTab === 'revenue'" class="revenue-analytics-section">
      <!-- Section Filter Bar -->
      <div class="revenue-filter-card">
        <div class="section-title-wrap">
          <div class="section-icon-badge bg-emerald-50 text-emerald-600">
            <BarChart3 class="w-4 h-4" />
          </div>
          <div>
            <h2 class="section-title">Laporan Keuangan & Performa Bulanan</h2>
            <p class="section-desc">Kalkulasi Revenue, ADR, Occupancy Rate, dan RevPAR dari database (IDR)</p>
          </div>
        </div>

        <div class="revenue-controls">
          <select v-model="selectedMonth" class="custom-select">
            <option v-for="m in monthsList" :key="m.value" :value="m.value">
              {{ m.label }}
            </option>
          </select>

          <select v-model="selectedYear" class="custom-select">
            <option v-for="y in yearsList" :key="y" :value="y">
              {{ y }}
            </option>
          </select>

          <select v-model="selectedProperty" class="custom-select">
            <option value="all">Semua Kamar ({{ rooms.length }})</option>
            <option v-for="r in rooms" :key="r.id" :value="r.id">
              {{ r.name }}
            </option>
          </select>

          <!-- Export Buttons -->
          <button class="btn-export-icon" @click="exportCSV" title="Unduh CSV">
            <FileSpreadsheet class="w-4 h-4 text-emerald-600" />
            <span>CSV</span>
          </button>
          <button class="btn-export-icon" @click="exportPDF" title="Cetak PDF">
            <Download class="w-4 h-4 text-rose-600" />
            <span>PDF</span>
          </button>
        </div>
      </div>

      <!-- 6 KPI Metric Cards in IDR -->
      <div class="kpi-6-grid">
        <div class="kpi-metric-card">
          <div class="metric-top">
            <span class="metric-label">Reservations</span>
            <CalendarIcon class="w-4 h-4 text-blue-600" />
          </div>
          <span class="metric-val">{{ totalReservations }}</span>
          <span class="metric-hint">Total booking di database</span>
        </div>

        <div class="kpi-metric-card">
          <div class="metric-top">
            <span class="metric-label">Total Nights</span>
            <BedDouble class="w-4 h-4 text-indigo-600" />
          </div>
          <span class="metric-val">{{ totalNights }}</span>
          <span class="metric-hint">Malam kamar terjual</span>
        </div>

        <div class="kpi-metric-card">
          <div class="metric-top">
            <span class="metric-label">Occupancy Rate</span>
            <Percent class="w-4 h-4 text-emerald-600" />
          </div>
          <span class="metric-val text-emerald-600">{{ occupancyRate }}%</span>
          <div class="metric-progress-bg">
            <div class="metric-progress-fill" :style="{ width: `${occupancyRate}%` }"></div>
          </div>
        </div>

        <div class="kpi-metric-card highlight-revenue">
          <div class="metric-top">
            <span class="metric-label">Total Revenue</span>
            <DollarSign class="w-4 h-4 text-amber-600" />
          </div>
          <span class="metric-val font-extrabold text-indigo-700">{{ formatIDR(totalRevenue) }}</span>
          <span class="metric-hint">∑ Total Pendapatan Kotor</span>
        </div>

        <div class="kpi-metric-card">
          <div class="metric-top">
            <span class="metric-label">ADR (Avg Daily Rate)</span>
            <TrendingUp class="w-4 h-4 text-purple-600" />
          </div>
          <span class="metric-val">{{ formatIDR(adr) }}</span>
          <span class="metric-hint">Revenue / Total Nights</span>
        </div>

        <div class="kpi-metric-card">
          <div class="metric-top">
            <span class="metric-label">RevPAR</span>
            <Building2 class="w-4 h-4 text-teal-600" />
          </div>
          <span class="metric-val">{{ formatIDR(revPAR) }}</span>
          <span class="metric-hint">Revenue / Total Kamar × Hari</span>
        </div>
      </div>

      <!-- Charts Row -->
      <div class="charts-2-grid">
        <div class="chart-box">
          <div class="chart-header">
            <h3>Daily Revenue Trend ({{ monthTitle }})</h3>
            <p>Grafik tren pendapatan harian dari database (IDR)</p>
          </div>
          <VueApexCharts
            type="bar"
            height="260"
            :options="revenueChartOptions"
            :series="revenueChartSeries"
          />
        </div>

        <div class="chart-box channel-split-chart-box">
          <div class="chart-header">
            <h3>Revenue Share by Channel ({{ monthTitle }})</h3>
            <p>Porsi kontribusi pendapatan OTA & Direct dari database</p>
          </div>
          
          <div class="donut-split-layout">
            <!-- Left: Donut Chart -->
            <div class="donut-chart-col">
              <VueApexCharts
                type="donut"
                height="220"
                :options="channelChartOptions"
                :series="channelChartSeries"
              />
            </div>

            <!-- Right: Channel Mini Breakdown Cards -->
            <div class="channel-mini-list-col">
              <div
                v-for="ch in channelBreakdown"
                :key="ch.key"
                class="ch-summary-card"
                :style="{ borderLeftColor: ch.color }"
              >
                <div class="ch-summary-top">
                  <div class="ch-name-group">
                    <span class="ch-indicator-dot" :style="{ backgroundColor: ch.color }"></span>
                    <span class="ch-summary-name">{{ ch.name }}</span>
                  </div>
                  <span class="ch-summary-pct" :style="{ color: ch.color }">{{ ch.percent }}%</span>
                </div>
                <div class="ch-summary-middle">
                  <span class="ch-summary-rev">{{ formatIDR(ch.revenue) }}</span>
                  <span class="ch-summary-sub">{{ ch.reservations }} Booking · {{ ch.nights }} Malam</span>
                </div>
                <div class="ch-summary-bar-bg">
                  <div class="ch-summary-bar-fill" :style="{ width: `${ch.percent}%`, backgroundColor: ch.color }"></div>
                </div>
              </div>

              <div v-if="channelBreakdown.length === 0" class="ch-empty-card">
                <span class="text-xs text-slate-400">Belum ada transaksi channel di database</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Property Performance Breakdown Table -->
      <div class="table-box">
        <div class="chart-header">
          <h3>Performance by Property / Room ({{ monthTitle }})</h3>
          <p>Tabel rincian performa okupansi, ADR, dan pendapatan tiap kamar dari database</p>
        </div>

        <div class="table-responsive">
          <table class="data-table">
            <thead>
              <tr>
                <th>Unit Kamar</th>
                <th class="text-center">Reservasi</th>
                <th class="text-center">Dibatalkan</th>
                <th class="text-center">Total Malam</th>
                <th class="text-center">Occupancy Rate</th>
                <th class="text-right">Total Revenue (IDR)</th>
                <th class="text-right">ADR (IDR)</th>
                <th class="text-right">RevPAR (IDR)</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="prop in propertyRows" :key="prop.id">
                <td class="font-bold text-slate-800">{{ prop.name }}</td>
                <td class="text-center font-bold text-slate-700">{{ prop.reservations }}</td>
                <td class="text-center font-bold" :class="prop.cancelled > 0 ? 'text-rose-600' : 'text-slate-400'">{{ prop.cancelled }}</td>
                <td class="text-center">{{ prop.nights }}</td>
                <td class="text-center font-semibold text-emerald-600">{{ prop.occupancy }}%</td>
                <td class="text-right font-bold text-indigo-900">{{ formatIDR(prop.revenue) }}</td>
                <td class="text-right">{{ formatIDR(prop.adr) }}</td>
                <td class="text-right font-semibold">{{ formatIDR(prop.revpar) }}</td>
              </tr>
              <tr v-if="propertyRows.length === 0">
                <td colspan="8" class="text-center py-6 text-slate-400 font-semibold">
                  Belum ada kamar terdaftar di database
                </td>
              </tr>
            </tbody>
            <tfoot>
              <tr class="total-row">
                <td class="font-extrabold uppercase">Total Keseluruhan</td>
                <td class="text-center font-extrabold">{{ totalReservations }}</td>
                <td class="text-center font-extrabold" :class="totalCancelledReservations > 0 ? 'text-rose-600' : 'text-slate-500'">{{ totalCancelledReservations }}</td>
                <td class="text-center font-extrabold">{{ totalNights }}</td>
                <td class="text-center font-extrabold text-emerald-700">{{ occupancyRate }}%</td>
                <td class="text-right font-extrabold text-indigo-700 text-base">{{ formatIDR(totalRevenue) }}</td>
                <td class="text-right font-extrabold">{{ formatIDR(adr) }}</td>
                <td class="text-right font-extrabold">{{ formatIDR(revPAR) }}</td>
              </tr>
            </tfoot>
          </table>
        </div>
      </div>

      <!-- Daily Occupancy Heatmap -->
      <div class="heatmap-box">
        <div class="chart-header">
          <h3>Daily Occupancy Heatmap</h3>
          <p>Tingkat hunian kamar harian pada bulan {{ monthTitle }} dari database</p>
        </div>

        <div class="heatmap-grid">
          <div
            v-for="d in dailyOccupancyMap"
            :key="d.day"
            class="heatmap-pill"
            :class="{
              'heat-empty': d.occupancyPercent === 0,
              'heat-low': d.occupancyPercent > 0 && d.occupancyPercent <= 35,
              'heat-mid': d.occupancyPercent > 35 && d.occupancyPercent <= 70,
              'heat-high': d.occupancyPercent > 70,
            }"
            :title="`${d.dateStr}: ${d.occupied}/${totalRoomsCount} kamar (${d.occupancyPercent}%)`"
          >
            <span class="heat-day">{{ d.dayName }}</span>
            <span class="heat-num">{{ d.day }}</span>
            <span class="heat-rate">{{ d.occupancyPercent }}%</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.master-dashboard {
  display: flex;
  flex-direction: column;
  gap: 24px;
  background-color: #f8fafc;
  font-family: 'Plus Jakarta Sans', system-ui, -apple-system, sans-serif;
  color: #0f172a;
}

/* Top Header Card */
.top-header-card {
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

.greeting-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: #ecfdf5;
  color: #065f46;
  border: 1px solid #a7f3d0;
  font-size: 0.75rem;
  font-weight: 700;
  padding: 3px 10px;
  border-radius: 9999px;
  margin-bottom: 8px;
}

.greeting-title {
  font-size: 1.4rem;
  font-weight: 800;
  color: #0f172a;
  letter-spacing: -0.02em;
  margin: 0;
}

.greeting-subtitle {
  font-size: 0.85rem;
  color: #64748b;
  margin: 4px 0 0;
}

.header-right-controls {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.btn-refresh-db {
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

.btn-refresh-db:hover {
  background: #f8fafc;
  border-color: #94a3b8;
}

.tab-switcher {
  display: flex;
  background: #f1f5f9;
  padding: 4px;
  border-radius: 12px;
  gap: 4px;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  font-size: 0.82rem;
  font-weight: 700;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn.active {
  background: #ffffff;
  color: #4f46e5;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.06);
}

.btn-channel-manage {
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

.btn-channel-manage:hover {
  background: #e0e7ff;
}

/* Section Common Styling */
.operations-section, .revenue-analytics-section, .channel-widget-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.section-title-wrap {
  display: flex;
  align-items: center;
  gap: 10px;
}

.section-icon-badge {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  background: #eef2ff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.section-title {
  font-size: 1.15rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.section-desc {
  font-size: 0.8rem;
  color: #64748b;
  margin: 2px 0 0;
}

/* 5 KPI Grid */
.kpi-5-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 14px;
}

@media (max-width: 1024px) {
  .kpi-5-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.kpi-card {
  background: #ffffff;
  padding: 16px 18px;
  border-radius: 14px;
  border: 1px solid #e2e8f0;
  cursor: pointer;
  transition: all 0.2s;
}

.kpi-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.05);
}

.kpi-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.kpi-label {
  font-size: 0.78rem;
  font-weight: 600;
  color: #64748b;
}

.kpi-val {
  font-size: 1.4rem;
  font-weight: 800;
  color: #0f172a;
}

.kpi-unit {
  font-size: 0.8rem;
  font-weight: 500;
  color: #94a3b8;
}

/* 2x2 Operations Grid */
.operations-2x2-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

@media (max-width: 768px) {
  .operations-2x2-grid {
    grid-template-columns: 1fr;
  }
}

.op-panel {
  background: #ffffff;
  border-radius: 14px;
  border: 1px solid #e2e8f0;
  padding: 16px 20px;
}

.panel-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.panel-head h3 {
  font-size: 0.95rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.view-all-link {
  background: transparent;
  border: none;
  font-size: 0.78rem;
  font-weight: 700;
  color: #4f46e5;
  cursor: pointer;
}

.empty-panel {
  padding: 24px 0;
  text-align: center;
  font-size: 0.85rem;
  color: #94a3b8;
}

.arrival-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f1f5f9;
}

.arrival-row:last-child {
  border-bottom: none;
}

.guest-info {
  display: flex;
  flex-direction: column;
}

.guest-name {
  font-size: 0.85rem;
  font-weight: 700;
  color: #1e293b;
}

.unit-name {
  font-size: 0.75rem;
  color: #64748b;
}

.channel-pill {
  padding: 3px 10px;
  border-radius: 9999px;
  font-size: 0.72rem;
  font-weight: 800;
}

.badge-checkedout {
  padding: 3px 10px;
  border-radius: 9999px;
  background: #f1f5f9;
  color: #475569;
  font-size: 0.72rem;
  font-weight: 700;
}

/* Upcoming Arrivals Next 60 Days */
.upcoming-card {
  background: #ffffff;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  padding: 20px;
}

.upcoming-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 12px;
}

.upcoming-head h3 {
  font-size: 1.05rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.custom-select-sm {
  padding: 6px 12px;
  font-size: 0.8rem;
  font-weight: 700;
  border-radius: 8px;
  border: 1px solid #cbd5e1;
  background: #f8fafc;
  outline: none;
}

.btn-calendar-link {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  font-size: 0.8rem;
  font-weight: 700;
  border-radius: 8px;
  border: 1px solid #cbd5e1;
  background: #ffffff;
  color: #334155;
  cursor: pointer;
}

.btn-calendar-link:hover {
  background: #f1f5f9;
}

.arrival-group {
  margin-bottom: 18px;
}

.group-title {
  font-size: 0.82rem;
  font-weight: 800;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.03em;
  margin-bottom: 8px;
  display: block;
}

.group-items-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.guest-arrival-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-radius: 12px;
  border: 1px solid #f1f5f9;
  background: #fafbff;
  cursor: pointer;
  transition: all 0.15s;
}

.guest-arrival-card:hover {
  border-color: #cbd5e1;
  background: #ffffff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.04);
}

.guest-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.guest-avatar {
  width: 36px;
  height: 36px;
  border-radius: 9999px;
  background: #e0e7ff;
  color: #4338ca;
  font-weight: 800;
  font-size: 0.85rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.guest-meta {
  display: flex;
  flex-direction: column;
}

.stay-info {
  font-size: 0.78rem;
  color: #64748b;
}

.empty-state-box {
  text-align: center;
  padding: 32px 16px;
  background: #f8fafc;
  border-radius: 12px;
  border: 1px dashed #cbd5e1;
}

/* Distribution Channel Widget Section */
.channel-widget-card {
  background: #ffffff;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  padding: 20px;
}

.widget-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 18px;
  flex-wrap: wrap;
  gap: 12px;
}

.btn-manage-accent {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  font-size: 0.82rem;
  font-weight: 700;
  border-radius: 10px;
  border: none;
  background: #4f46e5;
  color: #ffffff;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 2px 5px rgba(79, 70, 229, 0.2);
}

.btn-manage-accent:hover {
  background: #4338ca;
}

.channel-grid-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 14px;
}

.channel-stat-box {
  background: #f8fafc;
  border-radius: 12px;
  padding: 14px 16px;
  border: 1px solid #e2e8f0;
  border-left-width: 4px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.ch-box-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.ch-badge {
  padding: 2px 8px;
  border-radius: 6px;
  color: #ffffff;
  font-size: 0.72rem;
  font-weight: 800;
}

.ch-share-pct {
  font-size: 0.78rem;
  font-weight: 700;
  color: #64748b;
}

.ch-box-revenue {
  display: flex;
  flex-direction: column;
}

.ch-rev-val {
  font-size: 1.15rem;
  font-weight: 800;
  color: #0f172a;
}

.ch-sub-info {
  font-size: 0.72rem;
  color: #94a3b8;
}

.ch-progress-bg {
  width: 100%;
  height: 4px;
  background: #e2e8f0;
  border-radius: 9999px;
  overflow: hidden;
}

.ch-progress-fill {
  height: 100%;
  border-radius: 9999px;
}

/* Revenue Filter Card */
.revenue-filter-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 14px;
  background: #ffffff;
  padding: 16px 20px;
  border-radius: 14px;
  border: 1px solid #e2e8f0;
}

.revenue-controls {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.custom-select {
  padding: 6px 12px;
  font-size: 0.82rem;
  font-weight: 700;
  border-radius: 8px;
  border: 1px solid #cbd5e1;
  background: #f8fafc;
  color: #334155;
  outline: none;
}

.btn-export-icon {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 10px;
  font-size: 0.8rem;
  font-weight: 700;
  border-radius: 8px;
  border: 1px solid #cbd5e1;
  background: #ffffff;
  cursor: pointer;
}

/* 6 KPI Metric Cards in IDR */
.kpi-6-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 14px;
}

@media (max-width: 1280px) {
  .kpi-6-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .kpi-6-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.kpi-metric-card {
  background: #ffffff;
  border-radius: 14px;
  padding: 14px 16px;
  border: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.kpi-metric-card.highlight-revenue {
  background: linear-gradient(135deg, #ffffff 0%, #f5f3ff 100%);
  border-color: #c7d2fe;
}

.metric-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.metric-label {
  font-size: 0.74rem;
  font-weight: 700;
  color: #64748b;
}

.metric-val {
  font-size: 1.25rem;
  font-weight: 800;
  color: #0f172a;
}

.metric-hint {
  font-size: 0.68rem;
  color: #94a3b8;
  margin-top: 2px;
}

.metric-progress-bg {
  width: 100%;
  height: 4px;
  background: #e2e8f0;
  border-radius: 9999px;
  margin-top: 4px;
  overflow: hidden;
}

.metric-progress-fill {
  height: 100%;
  background: #059669;
}

/* Charts 2-Grid */
.charts-2-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 18px;
}

@media (max-width: 1024px) {
  .charts-2-grid {
    grid-template-columns: 1fr;
  }
}

.chart-box {
  background: #ffffff;
  border-radius: 16px;
  padding: 20px;
  border: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
}

.chart-header {
  margin-bottom: 12px;
}

.chart-header h3 {
  font-size: 1rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.chart-header p {
  font-size: 0.78rem;
  color: #64748b;
  margin: 2px 0 0;
}

/* Donut Split Layout */
.donut-split-layout {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  flex: 1;
  padding: 6px 0;
}

@media (max-width: 640px) {
  .donut-split-layout {
    flex-direction: column;
  }
}

.donut-chart-col {
  width: 190px;
  min-width: 190px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.channel-mini-list-col {
  flex: 1;
  max-width: 250px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  max-height: 220px;
  overflow-y: auto;
  scrollbar-width: thin;
  scrollbar-color: #cbd5e1 transparent;
  padding-right: 3px;
}

.channel-mini-list-col::-webkit-scrollbar {
  width: 4px;
}

.channel-mini-list-col::-webkit-scrollbar-track {
  background: transparent;
}

.channel-mini-list-col::-webkit-scrollbar-thumb {
  background-color: #cbd5e1;
  border-radius: 9999px;
}

.ch-summary-card {
  background: #f8fafc;
  border-radius: 8px;
  padding: 6px 10px;
  border: 1px solid #e2e8f0;
  border-left-width: 3px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  transition: all 0.15s;
}

.ch-summary-card:hover {
  background: #ffffff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.ch-summary-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.ch-name-group {
  display: flex;
  align-items: center;
  gap: 5px;
}

.ch-indicator-dot {
  width: 6px;
  height: 6px;
  border-radius: 9999px;
  flex-shrink: 0;
}

.ch-summary-name {
  font-size: 0.74rem;
  font-weight: 700;
  color: #1e293b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.ch-summary-pct {
  font-size: 0.74rem;
  font-weight: 800;
}

.ch-summary-middle {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

.ch-summary-rev {
  font-size: 0.78rem;
  font-weight: 800;
  color: #0f172a;
}

.ch-summary-sub {
  font-size: 0.65rem;
  color: #94a3b8;
}

.ch-summary-bar-bg {
  width: 100%;
  height: 3px;
  background: #e2e8f0;
  border-radius: 9999px;
  overflow: hidden;
  margin-top: 1px;
}

.ch-summary-bar-fill {
  height: 100%;
  border-radius: 9999px;
}

.ch-empty-card {
  text-align: center;
  padding: 16px 0;
}

/* Table Box */
.table-box {
  background: #ffffff;
  border-radius: 16px;
  padding: 18px;
  border: 1px solid #e2e8f0;
}

.table-responsive {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.85rem;
}

.data-table th {
  background: #f8fafc;
  padding: 10px 14px;
  font-size: 0.75rem;
  font-weight: 800;
  color: #475569;
  text-transform: uppercase;
  border-bottom: 2px solid #e2e8f0;
}

.data-table td {
  padding: 12px 14px;
  border-bottom: 1px solid #f1f5f9;
}

.total-row td {
  background: #f8fafc;
  border-top: 2px solid #cbd5e1;
  border-bottom: none;
}

/* Heatmap Box */
.heatmap-box {
  background: #ffffff;
  border-radius: 16px;
  padding: 18px;
  border: 1px solid #e2e8f0;
}

.heatmap-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(60px, 1fr));
  gap: 6px;
}

.heatmap-pill {
  border-radius: 8px;
  padding: 8px 4px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: 1px solid #e2e8f0;
}

.heat-empty { background: #f8fafc; color: #94a3b8; }
.heat-low { background: #fef2f2; border-color: #fecdd3; color: #be123c; }
.heat-mid { background: #fee2e2; border-color: #fda4af; color: #9f1239; }
.heat-high { background: #e11d48; border-color: #be123c; color: #ffffff; }

.heat-day { font-size: 0.62rem; font-weight: 700; text-transform: uppercase; }
.heat-num { font-size: 0.9rem; font-weight: 800; }
.heat-rate { font-size: 0.65rem; font-weight: 700; }

/* ======================================================== */
/* COMPREHENSIVE RESPONSIVE STYLING (MOBILE & TABLET)       */
/* ======================================================== */
@media (max-width: 1024px) {
  .master-dashboard {
    padding: 16px;
    gap: 16px;
  }
  .operations-2x2-grid {
    grid-template-columns: 1fr;
  }
  .charts-2-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .master-dashboard {
    padding: 12px;
    gap: 14px;
  }
  .top-header-card {
    padding: 16px;
    flex-direction: column;
    align-items: stretch;
    gap: 14px;
  }
  .header-actions {
    width: 100%;
    justify-content: space-between;
  }
  .tab-switcher {
    width: 100%;
    overflow-x: auto;
    padding-bottom: 2px;
  }
  .kpi-5-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
  }
  .kpi-6-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
  }
  .revenue-filter-card {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  .revenue-controls {
    width: 100%;
    justify-content: space-between;
  }
  .custom-select {
    flex: 1;
    min-width: 110px;
  }
  .donut-split-layout {
    flex-direction: column;
    align-items: center;
  }
  .channel-mini-list-col {
    max-width: 100%;
    width: 100%;
  }
  .heatmap-grid {
    grid-template-columns: repeat(auto-fill, minmax(44px, 1fr));
    gap: 4px;
  }
  .heatmap-pill {
    padding: 6px 2px;
  }
}

@media (max-width: 480px) {
  .kpi-5-grid {
    grid-template-columns: 1fr;
  }
  .kpi-6-grid {
    grid-template-columns: 1fr;
  }
  .channel-grid-cards {
    grid-template-columns: 1fr;
  }
  .greeting-title {
    font-size: 1.2rem;
  }
}
</style>