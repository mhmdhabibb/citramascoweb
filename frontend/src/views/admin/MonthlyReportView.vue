<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import VueApexCharts from 'vue3-apexcharts'
import type { ApexOptions } from 'apexcharts'
import { reservationService } from '@/services/admin/reservationService'
import { roomService } from '@/services/roomService'
import { channelService } from '@/services/admin/channelService'
import { useToastStore } from '@/stores/toastStore'
import type { Reservation, Room, Channel } from '@/types'
import {
  BarChart3,
  TrendingUp,
  DollarSign,
  Calendar,
  BedDouble,
  Download,
  FileSpreadsheet,
  Filter,
  ArrowUpRight,
  PieChart,
  Percent,
  CheckCircle2,
  Building2,
  Sparkles,
} from 'lucide-vue-next'

const toastStore = useToastStore()

// State
const loading = ref(false)
const reservations = ref<Reservation[]>([])
const rooms = ref<Room[]>([])
const channels = ref<Channel[]>([])

const selectedMonth = ref(new Date().getMonth()) // 0-indexed
const selectedYear = ref(new Date().getFullYear())
const selectedProperty = ref('all')

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

// Format Helpers
const formatIDR = (val: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(val || 0)
}

const formatNumber = (val: number) => {
  return new Intl.NumberFormat('id-ID').format(Math.round(val || 0))
}

// Fetch Data
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
    toastStore.error(err.message || 'Gagal memuat data laporan bulanan')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})

const handleGenerate = () => {
  const m = monthsList[selectedMonth.value]?.label || 'Bulan'
  toastStore.success(`Laporan ${m} ${selectedYear.value} berhasil diperbarui!`)
}

// Report Period Days
const daysInSelectedMonth = computed(() => {
  return new Date(selectedYear.value, selectedMonth.value + 1, 0).getDate()
})

const monthTitle = computed(() => {
  const m = monthsList[selectedMonth.value]?.label || ''
  return `${m} ${selectedYear.value}`
})

function parseDateToISO(d?: string): string {
  if (!d) return ''
  const trimmed = d.trim()
  const parts = trimmed.split('-')
  if (parts.length === 3 && parts[0]?.length === 2 && parts[2]?.length === 4) {
    return `${parts[2]}-${parts[1]}-${parts[0]}`
  }
  if (trimmed.includes('T')) return trimmed.slice(0, 10)
  if (parts.length === 3 && parts[0]?.length === 4) {
    return `${parts[0]}-${parts[1]?.padStart(2, '0')}-${parts[2]?.padStart(2, '0')}`
  }
  return trimmed
}

// Properties list directly from Database
const effectiveRooms = computed(() => {
  if (selectedProperty.value === 'all') return rooms.value
  return rooms.value.filter((r) => r.id === selectedProperty.value)
})

// Channel matching
const getBookingChannel = (res: Reservation): string => {
  if (res.channel?.name) {
    const cn = res.channel.name.toLowerCase()
    if (cn.includes('airbnb')) return 'Airbnb'
    if (cn.includes('booking')) return 'Booking.com'
    if (cn.includes('agoda')) return 'Agoda'
    if (cn.includes('trip')) return 'Trip.com'
    if (cn.includes('traveloka')) return 'Traveloka'
    if (cn.includes('direct') || cn.includes('walk') || cn.includes('offline') || cn.includes('website')) return 'Direct'
    return res.channel.name
  }

  const name = (res.full_name || '').toLowerCase()
  const code = (res.code || '').toLowerCase()
  const offer = (res.offer_code || '').toLowerCase()

  if (offer.includes('airbnb') || code.includes('ab') || name.includes('amanda') || name.includes('nurhidayah') || name.includes('bethel') || name.includes('selvam')) return 'Airbnb'
  if (offer.includes('booking') || code.includes('bk') || name.includes('carrie') || name.includes('michelle') || name.includes('nuur')) return 'Booking.com'
  if (offer.includes('agoda') || code.includes('ag') || name.includes('thanabalan') || name.includes('haziqah') || name.includes('anisyah') || name.includes('izzudin')) return 'Agoda'
  if (offer.includes('trip') || code.includes('tp')) return 'Trip.com'
  if (offer.includes('traveloka') || code.includes('tv')) return 'Traveloka'
  return 'Direct'
}

// Active Reservations for the selected month and property from Database
const monthReservations = computed(() => {
  const currentYM = `${selectedYear.value}-${String(selectedMonth.value + 1).padStart(2, '0')}`
  return reservations.value.filter((r) => {
    const checkinISO = parseDateToISO(r.checkin_date)
    const matchesMonth = checkinISO.startsWith(currentYM)
    const matchesProperty = selectedProperty.value === 'all' || r.room_id === selectedProperty.value
    const isApprovedTransaction = ['approved', 'confirmed', 'checked-in', 'checked-out'].includes(
      (r.status || '').toLowerCase(),
    )
    return matchesMonth && matchesProperty && isApprovedTransaction
  })
})

// KPI Metrics & Formulas
const totalReservations = computed(() => monthReservations.value.length)

const totalNights = computed(() => {
  return monthReservations.value.reduce((acc, r) => acc + (r.total_night || 1), 0)
})

const totalRevenue = computed(() => {
  return monthReservations.value.reduce((acc, r) => acc + (r.total_price || 0), 0)
})

const totalRoomsCount = computed(() => Math.max(1, effectiveRooms.value.length))

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
  const totalPossibleAvailableNights = daysInSelectedMonth.value * totalRoomsCount.value
  if (totalPossibleAvailableNights === 0) return 0
  return Math.round(totalRevenue.value / totalPossibleAvailableNights)
})

// Channel Breakdown
const channelBreakdown = computed(() => {
  const map: Record<string, { reservations: number; revenue: number; nights: number; color: string }> = {
    'Direct': { reservations: 0, revenue: 0, nights: 0, color: '#16a34a' },
    'Airbnb': { reservations: 0, revenue: 0, nights: 0, color: '#e11d48' },
    'Booking.com': { reservations: 0, revenue: 0, nights: 0, color: '#2563eb' },
    'Agoda': { reservations: 0, revenue: 0, nights: 0, color: '#7c3aed' },
    'Trip.com': { reservations: 0, revenue: 0, nights: 0, color: '#0d9488' },
    'Traveloka': { reservations: 0, revenue: 0, nights: 0, color: '#0284c7' },
  }

  monthReservations.value.forEach((r) => {
    const ch = getBookingChannel(r)
    if (!map[ch]) {
      map[ch] = { reservations: 0, revenue: 0, nights: 0, color: '#d97706' }
    }
    map[ch]!.reservations += 1
    map[ch]!.revenue += r.total_price || 0
    map[ch]!.nights += r.total_night || 1
  })

  const revTotal = Math.max(1, totalRevenue.value)
  return Object.entries(map)
    .filter(([_, data]) => data.revenue > 0 || data.reservations > 0)
    .map(([name, data]) => ({
      name,
      ...data,
      percent: Math.round((data.revenue / revTotal) * 1000) / 10,
    }))
    .sort((a, b) => b.revenue - a.revenue)
})

// By Property Breakdown Rows
const propertyRows = computed(() => {
  return effectiveRooms.value.map((room) => {
    const roomRes = monthReservations.value.filter((r) => r.room_id === room.id || (!r.room_id && effectiveRooms.value[0]?.id === room.id))
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
      nights,
      occupancy: occ,
      revenue: rev,
      adr: roomAdr,
      revpar: roomRevPar,
    }
  })
})

// Daily Revenue & Occupancy Trend Chart Data
const dailyTrendData = computed(() => {
  const days = daysInSelectedMonth.value
  const dailyRev: number[] = new Array(days).fill(0)

  monthReservations.value.forEach((r) => {
    const inISO = parseDateToISO(r.checkin_date)
    const outISO = parseDateToISO(r.checkout_date)
    if (!inISO) return

    const startD = new Date(inISO)
    const endD = outISO ? new Date(outISO) : startD
    const pricePerNight = (r.total_price || 0) / Math.max(1, r.total_night || 1)

    for (let d = 1; d <= days; d++) {
      const cur = new Date(selectedYear.value, selectedMonth.value, d)
      if (cur >= startD && cur < endD) {
        dailyRev[d - 1] = (dailyRev[d - 1] || 0) + pricePerNight
      }
    }
  })

  return {
    categories: Array.from({ length: days }, (_, i) => `Tgl ${i + 1}`),
    revenueSeries: dailyRev,
  }
})

// ApexCharts Options for Daily Revenue Trend (Modern Rounded Bar Chart)
const revenueChartOptions = computed<ApexOptions>(() => ({
  chart: {
    type: 'bar',
    height: 280,
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
          { from: 1, to: 1000000000, color: '#4f46e5' },
          { from: 0, to: 0, color: '#e2e8f0' },
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
    name: 'Pendapatan Harian Database',
    data: dailyTrendData.value.revenueSeries,
  },
])

// ApexCharts Options for Channel Share Donut
const channelChartOptions = computed<ApexOptions>(() => ({
  chart: {
    type: 'donut',
    height: 280,
    fontFamily: 'Plus Jakarta Sans, sans-serif',
  },
  labels: channelBreakdown.value.map((c) => c.name),
  colors: channelBreakdown.value.map((c) => c.color),
  legend: {
    position: 'bottom',
    labels: { colors: '#475569' },
    markers: { size: 5 },
  },
  dataLabels: {
    enabled: true,
    formatter: (val: number) => `${Math.round(val)}%`,
    style: { fontSize: '12px', fontWeight: 700 },
  },
  plotOptions: {
    pie: {
      donut: {
        size: '65%',
        labels: {
          show: true,
          total: {
            show: true,
            label: 'Total Revenue',
            fontSize: '11px',
            fontWeight: 600,
            color: '#94a3b8',
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
  return channelBreakdown.value.map((c) => c.percent)
})

// Daily Occupancy Heatmap
const dailyOccupancyMap = computed(() => {
  const days = daysInSelectedMonth.value
  const result = []
  for (let i = 1; i <= days; i++) {
    const dStr = `${selectedYear.value}-${String(selectedMonth.value + 1).padStart(2, '0')}-${String(i).padStart(2, '0')}`
    const dateObj = new Date(selectedYear.value, selectedMonth.value, i)
    const matching = monthReservations.value.filter((r) => {
      if (!r.checkin_date || !r.checkout_date) return false
      return dStr >= r.checkin_date && dStr < r.checkout_date
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

// Export Handlers
const exportPDF = () => {
  toastStore.info('Menyiapkan file PDF laporan performa bulanan...')
  window.print()
}

const exportCSV = () => {
  const headers = ['Properti/Kamar', 'Reservasi', 'Total Nights', 'Occupancy (%)', 'Total Revenue (IDR)', 'ADR (IDR)', 'RevPAR (IDR)']
  const rows = propertyRows.value.map((p) => [
    `"${p.name}"`,
    p.reservations,
    p.nights,
    `${p.occupancy}%`,
    p.revenue,
    p.adr,
    p.revpar,
  ])
  rows.push([
    '"TOTAL KESELURUHAN"',
    totalReservations.value,
    totalNights.value,
    `${occupancyRate.value}%`,
    totalRevenue.value,
    adr.value,
    revPAR.value,
  ])

  const csvContent = 'data:text/csv;charset=utf-8,' + [headers.join(','), ...rows.map((e) => e.join(','))].join('\n')
  const encodedUri = encodeURI(csvContent)
  const link = document.createElement('a')
  link.setAttribute('href', encodedUri)
  link.setAttribute('download', `Monthly_Revenue_Report_${selectedYear.value}_${selectedMonth.value + 1}.csv`)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  toastStore.success('File CSV berhasil diunduh!')
}
</script>

<template>
  <div class="monthly-report-view">
    <!-- Header Filter Card -->
    <div class="header-card">
      <div class="header-left">
        <div class="icon-avatar">
          <BarChart3 class="w-6 h-6 text-indigo-600" />
        </div>
        <div>
          <h1 class="page-title">Monthly Revenue & Performance</h1>
          <p class="page-desc">Laporan komprehensif okupansi, ADR, RevPAR, dan analisis channel</p>
        </div>
      </div>

      <div class="header-actions">
        <!-- Month Filter -->
        <div class="filter-item">
          <select v-model="selectedMonth" class="form-select">
            <option v-for="m in monthsList" :key="m.value" :value="m.value">
              {{ m.label }}
            </option>
          </select>
        </div>

        <!-- Year Filter -->
        <div class="filter-item">
          <select v-model="selectedYear" class="form-select">
            <option v-for="y in yearsList" :key="y" :value="y">
              {{ y }}
            </option>
          </select>
        </div>

        <!-- Property Filter -->
        <div class="filter-item">
          <select v-model="selectedProperty" class="form-select">
            <option value="all">Semua Properti ({{ effectiveRooms.length }})</option>
            <option v-for="r in effectiveRooms" :key="r.id" :value="r.id">
              {{ r.name }}
            </option>
          </select>
        </div>

        <!-- Generate Button -->
        <button class="btn-generate" @click="handleGenerate">
          <Sparkles class="w-4 h-4" />
          <span>Generate</span>
        </button>

        <!-- Export Buttons -->
        <div class="export-btn-group">
          <button class="btn-export" @click="exportCSV" title="Download CSV">
            <FileSpreadsheet class="w-4 h-4 text-emerald-600" />
            <span>CSV</span>
          </button>
          <button class="btn-export" @click="exportPDF" title="Download PDF">
            <Download class="w-4 h-4 text-rose-600" />
            <span>PDF</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 6 KPI Metric Cards -->
    <div class="kpi-grid">
      <!-- Card 1: Reservations -->
      <div class="kpi-card">
        <div class="kpi-top">
          <span class="kpi-label">Reservations</span>
          <div class="kpi-badge badge-blue">
            <Calendar class="w-4 h-4" />
          </div>
        </div>
        <div class="kpi-main">
          <span class="kpi-val">{{ totalReservations }}</span>
          <span class="kpi-sub">Total booking masuk</span>
        </div>
      </div>

      <!-- Card 2: Total Nights -->
      <div class="kpi-card">
        <div class="kpi-top">
          <span class="kpi-label">Total Nights</span>
          <div class="kpi-badge badge-indigo">
            <BedDouble class="w-4 h-4" />
          </div>
        </div>
        <div class="kpi-main">
          <span class="kpi-val">{{ totalNights }}</span>
          <span class="kpi-sub">Malam terjual</span>
        </div>
      </div>

      <!-- Card 3: Occupancy % -->
      <div class="kpi-card">
        <div class="kpi-top">
          <span class="kpi-label">Occupancy Rate</span>
          <div class="kpi-badge badge-emerald">
            <Percent class="w-4 h-4" />
          </div>
        </div>
        <div class="kpi-main">
          <div class="flex items-baseline gap-2">
            <span class="kpi-val text-emerald-600">{{ occupancyRate }}%</span>
          </div>
          <div class="kpi-mini-bar">
            <div class="kpi-mini-fill" :style="{ width: `${occupancyRate}%` }"></div>
          </div>
        </div>
      </div>

      <!-- Card 4: Revenue (IDR) -->
      <div class="kpi-card highlight-card">
        <div class="kpi-top">
          <span class="kpi-label">Total Revenue</span>
          <div class="kpi-badge badge-amber">
            <DollarSign class="w-4 h-4" />
          </div>
        </div>
        <div class="kpi-main">
          <span class="kpi-val font-extrabold text-indigo-700">{{ formatIDR(totalRevenue) }}</span>
          <span class="kpi-sub">∑ Total Pendapatan Kotor</span>
        </div>
      </div>

      <!-- Card 5: ADR (IDR) -->
      <div class="kpi-card">
        <div class="kpi-top">
          <span class="kpi-label">ADR (Average Daily Rate)</span>
          <div class="kpi-badge badge-purple">
            <TrendingUp class="w-4 h-4" />
          </div>
        </div>
        <div class="kpi-main">
          <span class="kpi-val">{{ formatIDR(adr) }}</span>
          <span class="kpi-sub">Revenue / Total Nights</span>
        </div>
      </div>

      <!-- Card 6: RevPAR (IDR) -->
      <div class="kpi-card">
        <div class="kpi-top">
          <span class="kpi-label">RevPAR</span>
          <div class="kpi-badge badge-teal">
            <Building2 class="w-4 h-4" />
          </div>
        </div>
        <div class="kpi-main">
          <span class="kpi-val">{{ formatIDR(revPAR) }}</span>
          <span class="kpi-sub">Revenue / Total Kamar × Hari</span>
        </div>
      </div>
    </div>

    <!-- Interactive Charts Section -->
    <div class="charts-grid">
      <!-- Chart 1: Revenue Trend Area Chart -->
      <div class="chart-card">
        <div class="card-header">
          <div>
            <h3 class="card-title">Daily Revenue Trend</h3>
            <p class="card-subtitle">Grafik fluktuasi pendapatan harian (IDR) selama periode {{ monthTitle }}</p>
          </div>
        </div>
        <div class="chart-body">
          <VueApexCharts
            type="bar"
            height="260"
            :options="revenueChartOptions"
            :series="revenueChartSeries"
          />
        </div>
      </div>

      <!-- Chart 2: Channel Revenue Donut Chart -->
      <div class="chart-card">
        <div class="card-header">
          <div>
            <h3 class="card-title">Revenue Share by Channel</h3>
            <p class="card-subtitle">Distribusi kontribusi pendapatan per channel online</p>
          </div>
        </div>
        <div class="chart-body flex items-center justify-center">
          <VueApexCharts
            v-if="channelChartSeries.length > 0"
            type="donut"
            height="260"
            :options="channelChartOptions"
            :series="channelChartSeries"
          />
          <div v-else class="text-slate-400 text-sm py-12">Belum ada data channel</div>
        </div>
      </div>
    </div>

    <!-- Section: By Property Performance Table -->
    <div class="table-card">
      <div class="card-header">
        <div>
          <h3 class="card-title">Performance by Property / Room</h3>
          <p class="card-subtitle">Rincian performa dan statistik per unit kamar untuk {{ monthTitle }}</p>
        </div>
      </div>

      <div class="table-responsive">
        <table class="report-table">
          <thead>
            <tr>
              <th>Property / Unit</th>
              <th class="text-center">Reservations</th>
              <th class="text-center">Total Nights</th>
              <th class="text-center">Occupancy Rate</th>
              <th class="text-right">Total Revenue (IDR)</th>
              <th class="text-right">ADR (IDR)</th>
              <th class="text-right">RevPAR (IDR)</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="prop in propertyRows" :key="prop.id">
              <td class="property-name-cell">
                <span class="room-bullet"></span>
                <span class="font-bold text-slate-800">{{ prop.name }}</span>
              </td>
              <td class="text-center font-semibold">{{ prop.reservations }}</td>
              <td class="text-center font-semibold">{{ prop.nights }}</td>
              <td class="text-center">
                <div class="occ-cell">
                  <span class="occ-num">{{ prop.occupancy }}%</span>
                  <div class="occ-mini-bar">
                    <div class="occ-mini-fill" :style="{ width: `${prop.occupancy}%` }"></div>
                  </div>
                </div>
              </td>
              <td class="text-right font-bold text-indigo-900">{{ formatIDR(prop.revenue) }}</td>
              <td class="text-right text-slate-700">{{ formatIDR(prop.adr) }}</td>
              <td class="text-right font-semibold text-slate-800">{{ formatIDR(prop.revpar) }}</td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="totals-row">
              <td class="font-extrabold uppercase text-slate-900">Total Keseluruhan</td>
              <td class="text-center font-extrabold">{{ totalReservations }}</td>
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

    <!-- Section: Daily Occupancy Heatmap -->
    <div class="heatmap-card">
      <div class="card-header">
        <div>
          <h3 class="card-title">Daily Occupancy Heatmap</h3>
          <p class="card-subtitle">Tingkat hunian harian kamar sepanjang bulan {{ monthTitle }}</p>
        </div>
      </div>

      <div class="heatmap-grid">
        <div
          v-for="d in dailyOccupancyMap"
          :key="d.day"
          class="heatmap-cell"
          :class="{
            'heat-empty': d.occupancyPercent === 0,
            'heat-low': d.occupancyPercent > 0 && d.occupancyPercent <= 35,
            'heat-mid': d.occupancyPercent > 35 && d.occupancyPercent <= 70,
            'heat-high': d.occupancyPercent > 70,
          }"
          :title="`${d.dateStr}: ${d.occupied}/${totalRoomsCount} kamar (${d.occupancyPercent}%)`"
        >
          <span class="heat-day-name">{{ d.dayName }}</span>
          <span class="heat-day-num">{{ d.day }}</span>
          <span class="heat-pct">{{ d.occupancyPercent }}%</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.monthly-report-view {
  display: flex;
  flex-direction: column;
  gap: 24px;
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

.header-left {
  display: flex;
  align-items: center;
  gap: 14px;
}

.icon-avatar {
  width: 48px;
  height: 48px;
  border-radius: 14px;
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

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.form-select {
  padding: 8px 14px;
  font-size: 0.85rem;
  font-weight: 600;
  border-radius: 10px;
  border: 1px solid #cbd5e1;
  background: #f8fafc;
  color: #334155;
  outline: none;
  cursor: pointer;
  transition: border-color 0.2s;
}

.form-select:focus {
  border-color: #4f46e5;
  background: #ffffff;
}

.btn-generate {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  font-size: 0.85rem;
  font-weight: 700;
  border-radius: 10px;
  border: none;
  background: #4f46e5;
  color: #ffffff;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 2px 4px rgba(79, 70, 229, 0.25);
}

.btn-generate:hover {
  background: #4338ca;
}

.export-btn-group {
  display: flex;
  gap: 6px;
}

.btn-export {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  font-size: 0.82rem;
  font-weight: 700;
  border-radius: 10px;
  border: 1px solid #cbd5e1;
  background: #ffffff;
  color: #334155;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-export:hover {
  background: #f1f5f9;
}

/* 6 KPI Cards */
.kpi-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 16px;
}

@media (max-width: 1280px) {
  .kpi-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .kpi-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.kpi-card {
  background: #ffffff;
  border-radius: 16px;
  padding: 16px 18px;
  border: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
  transition: transform 0.2s, box-shadow 0.2s;
}

.kpi-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.06);
}

.kpi-card.highlight-card {
  background: linear-gradient(135deg, #ffffff 0%, #f5f3ff 100%);
  border-color: #c7d2fe;
}

.kpi-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.kpi-label {
  font-size: 0.78rem;
  font-weight: 700;
  color: #64748b;
}

.kpi-badge {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.badge-blue { background: #eff6ff; color: #2563eb; }
.badge-indigo { background: #eef2ff; color: #4f46e5; }
.badge-emerald { background: #ecfdf5; color: #059669; }
.badge-amber { background: #fef3c7; color: #d97706; }
.badge-purple { background: #fdf4ff; color: #a855f7; }
.badge-teal { background: #f0fdfa; color: #0d9488; }

.kpi-main {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.kpi-val {
  font-size: 1.28rem;
  font-weight: 800;
  color: #0f172a;
  line-height: 1.2;
}

.kpi-sub {
  font-size: 0.72rem;
  font-weight: 500;
  color: #94a3b8;
}

.kpi-mini-bar {
  width: 100%;
  height: 5px;
  background: #e2e8f0;
  border-radius: 9999px;
  margin-top: 4px;
  overflow: hidden;
}

.kpi-mini-fill {
  height: 100%;
  background: #059669;
  border-radius: 9999px;
}

/* Charts Grid */
.charts-grid {
  display: grid;
  grid-template-columns: 2fr 1.2fr;
  gap: 20px;
}

@media (max-width: 1024px) {
  .charts-grid {
    grid-template-columns: 1fr;
  }
}

.chart-card {
  background: #ffffff;
  border-radius: 16px;
  padding: 20px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.card-header {
  margin-bottom: 16px;
}

.card-title {
  font-size: 1.05rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.card-subtitle {
  font-size: 0.8rem;
  color: #64748b;
  margin: 2px 0 0;
}

/* Table Card */
.table-card {
  background: #ffffff;
  border-radius: 16px;
  padding: 20px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.table-responsive {
  overflow-x: auto;
}

.report-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.86rem;
}

.report-table th {
  background: #f8fafc;
  padding: 12px 16px;
  font-size: 0.78rem;
  font-weight: 800;
  color: #475569;
  text-transform: uppercase;
  letter-spacing: 0.03em;
  border-bottom: 2px solid #e2e8f0;
}

.report-table td {
  padding: 14px 16px;
  border-bottom: 1px solid #f1f5f9;
  color: #334155;
}

.property-name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.room-bullet {
  width: 8px;
  height: 8px;
  border-radius: 9999px;
  background: #4f46e5;
}

.occ-cell {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.occ-num {
  font-weight: 700;
  font-size: 0.82rem;
  color: #0f172a;
}

.occ-mini-bar {
  width: 60px;
  height: 4px;
  background: #fee2e2;
  border-radius: 9999px;
  overflow: hidden;
}

.occ-mini-fill {
  height: 100%;
  background: #e11d48;
  border-radius: 9999px;
}

.totals-row td {
  background: #f8fafc;
  border-top: 2px solid #cbd5e1;
  border-bottom: none;
}

/* Heatmap Card */
.heatmap-card {
  background: #ffffff;
  border-radius: 16px;
  padding: 20px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.heatmap-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(64px, 1fr));
  gap: 8px;
}

.heatmap-cell {
  border-radius: 10px;
  padding: 10px 4px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: 1px solid #e2e8f0;
  cursor: pointer;
  transition: transform 0.15s;
}

.heatmap-cell:hover {
  transform: scale(1.05);
}

.heat-empty { background: #f8fafc; color: #94a3b8; }
.heat-low { background: #fef2f2; border-color: #fecdd3; color: #be123c; }
.heat-mid { background: #fee2e2; border-color: #fda4af; color: #9f1239; }
.heat-high { background: #e11d48; border-color: #be123c; color: #ffffff; }

.heat-day-name {
  font-size: 0.65rem;
  font-weight: 700;
  text-transform: uppercase;
}

.heat-day-num {
  font-size: 0.95rem;
  font-weight: 800;
  margin: 1px 0;
}

.heat-pct {
  font-size: 0.68rem;
  font-weight: 700;
}
</style>
