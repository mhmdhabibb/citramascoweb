<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useToastStore } from '@/stores/toastStore'
import { useToastStore as useToast } from '@/stores/toastStore'
import { floorService } from '@/services/admin/floorService'
import { roomUnitService } from '@/services/admin/roomUnitService'
import { useRouter } from 'vue-router'
import { Monitor, Building, CalendarCheck } from 'lucide-vue-next'

const router = useRouter()

const toast = useToast()
const floors = ref([])
const units = ref([])
const loading = ref(false)
let timer: any = null

const fetchData = async () => {
  try {
    loading.value = true
    const [f, u] = await Promise.all([floorService.getAll(), roomUnitService.getAll()])
    floors.value = f
    units.value = u
  } catch (e: any) {
    toast.error(e.message || 'Failed load')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
  timer = setInterval(() => fetchData(), 10000)
})
onUnmounted(() => {
  if (timer) clearInterval(timer)
})

const unitsByFloor = computed(() => {
  const map: Record<string, any[]> = {}
  units.value.forEach((unit) => {
    const fid = unit.floor_id || 'unknown'
    if (!map[fid]) map[fid] = []
    map[fid].push(unit)
  })
  return map
})

const statusClass = (s: string) => {
  switch (s) {
    case 'dirty':
      return 'bg-amber-400 text-amber-900'
    case 'maintenance':
      return 'bg-rose-400 text-rose-900'
    case 'occupied':
      return 'bg-slate-700 text-white'
    default:
      return 'bg-emerald-400 text-emerald-900'
  }
}

const changeStatus = async (id: string, status: string) => {
  try {
    loading.value = true
    await roomUnitService.setStatusAndNotify(id, status)
    toast.success(`Room ${status} set`)
    await fetchData()
  } catch (e: any) {
    toast.error(e.message || 'Failed change')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="px-6 py-5 bg-slate-50 min-h-screen">
    <div class="flex items-center mb-6">
      <Monitor class="w-7 h-7 mr-2 text-indigo-600" />
      <h2 class="text-2xl font-bold text-slate-800">Front Desk – Room Plan</h2>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-indigo-600"></div>
    </div>

    <div v-else class="space-y-10">
      <div v-for="floor in floors" :key="floor.id" class="bg-white rounded-2xl shadow p-6">
        <div class="flex items-center mb-5">
          <Building class="w-5 h-5 mr-2 text-slate-500" />
          <h3 class="text-lg font-semibold text-slate-700">Lantai {{ floor.name }}</h3>
          <div class="ml-auto text-xs text-slate-400">{{ (unitsByFloor[floor.id] || []).length }} unit</div>
        </div>

        <!-- Lift / hallway visual -->
        <div class="flex items-center mb-3">
          <div class="flex-1 border-t-2 border-dashed border-slate-300"></div>
          <svg class="mx-2 w-4 h-4 text-slate-400" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2l3 5h-6l3-5zm0 20l-3-5h6l-3 5z"/></svg>
          <div class="flex-1 border-t-2 border-dashed border-slate-300"></div>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
          <div
            v-for="unit in unitsByFloor[floor.id] || []"
            :key="unit.id"
            class="rounded-xl border border-slate-200 overflow-hidden bg-white"
          >
            <div class="px-3 py-2 border-b flex items-center justify-between">
              <span class="font-mono text-sm text-slate-800">{{ unit.room_number }}</span>
              <span :class="['px-2 py-0.5 text-[10px] rounded-full', statusClass(unit.status)]">{{ unit.status }}</span>
            </div>
            <div class="p-2 flex flex-col gap-1">
              <button class="text-xs bg-amber-500 hover:bg-amber-600 text-white px-2 py-1 rounded" @click="changeStatus(unit.id, 'dirty')">Mark Dirty</button>
              <button class="text-xs bg-rose-500 hover:bg-rose-600 text-white px-2 py-1 rounded" @click="changeStatus(unit.id, 'maintenance')">Maintenance</button>
            </div>
          </div>
          <div v-if="!(unitsByFloor[floor.id] || []).length" class="col-span-full text-sm text-slate-400 text-center py-4">Tidak ada unit pada lantai ini</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>
