<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { serviceRequestService, type AssignTaskPayload } from '@/services/admin/serviceRequestService'
import { userService } from '@/services/admin/userService'
import { useToastStore } from '@/stores/toastStore'
import { useAuthStore } from '@/stores/authStore'
import { roomService } from '@/services/roomService'
import type { ServiceRequest, User, Room, ServiceRequestCategory, ServiceRequestPriority } from '@/types'

const toastStore = useToastStore()
const authStore = useAuthStore()

const requests = ref<ServiceRequest[]>([])
const housekeepers = ref<User[]>([])
const rooms = ref<Room[]>([])
const loading = ref(false)
const statusFilter = ref('all')
const searchQuery = ref('')
const selectedRequest = ref<ServiceRequest | null>(null)

// Modal Catat Permintaan Tamu Baru (Input Resepsionis via Telepon/Walk-in)
const isCreateModalOpen = ref(false)
const createForm = ref({
  room_id: '',
  guest_name: '',
  guest_phone: '',
  source: 'phone_ext0',
  category: 'incident_broken_item' as ServiceRequestCategory,
  title: '',
  description: '',
  priority: 'urgent' as ServiceRequestPriority,
  direct_assign_housekeeper_id: '',
})

// Modal Tugaskan ke Housekeeping
const isAssignModalOpen = ref(false)
const assignForm = ref<AssignTaskPayload>({
  assigned_to_user_id: '',
  notes: '',
  priority: 'high',
})

// Modal Selesai (Complete Task)
const isCompleteModalOpen = ref(false)
const completeDamageCharge = ref(0)

let pollingTimer: any = null

const totalCount = computed(() => requests.value.length)
const pendingCount = computed(() => requests.value.filter(r => r.status === 'pending_reception').length)
const assignedCount = computed(() => requests.value.filter(r => r.status === 'assigned_to_housekeeping' || r.status === 'in_progress').length)
const completedCount = computed(() => requests.value.filter(r => r.status === 'completed').length)

const filteredRequests = computed(() => {
  return requests.value.filter((r) => {
    const matchesSearch =
      r.guest_name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      r.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      (r.room?.name || '').toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      (r.room?.code || '').toLowerCase().includes(searchQuery.value.toLowerCase())

    const matchesStatus = statusFilter.value === 'all' || r.status === statusFilter.value
    return matchesSearch && matchesStatus
  })
})

const fetchRequests = async (isSilent: boolean = false) => {
  try {
    if (!isSilent) loading.value = true
    const data = await serviceRequestService.getAll()
    requests.value = data || []
    if (selectedRequest.value) {
      const updated = data.find(r => r.id === selectedRequest.value?.id)
      if (updated) selectedRequest.value = updated
    }
  } catch (error: any) {
    if (!isSilent) toastStore.error(error.message || 'Failed to load service requests')
  } finally {
    if (!isSilent) loading.value = false
  }
}

const fetchHousekeepers = async () => {
  try {
    const data = await userService.getByRole('housekeeping')
    housekeepers.value = data || []
    if (housekeepers.value.length === 0) {
      const adminUsers = await userService.getByRole('admin')
      housekeepers.value = adminUsers || []
    }
  } catch (error) {
    console.error('Failed to fetch housekeepers', error)
  }
}

const fetchRooms = async () => {
  try {
    const data = await roomService.getAll()
    rooms.value = data || []
  } catch (error) {
    console.error('Failed to fetch rooms', error)
  }
}

onMounted(() => {
  fetchRequests(false)
  fetchHousekeepers()
  fetchRooms()
  pollingTimer = setInterval(() => {
    fetchRequests(true)
  }, 5000)
})

const openCreateModal = () => {
  createForm.value = {
    room_id: rooms.value[0]?.id || '',
    guest_name: 'Room Guest',
    guest_phone: '',
    source: 'phone_ext0',
    category: 'incident_broken_item',
    title: 'Broken Glass in Room',
    description: 'Guest called via Ext 0 reporting broken glass on the room floor, please sweep and replace immediately.',
    priority: 'urgent',
    direct_assign_housekeeper_id: housekeepers.value[0]?.id || '',
  }
  isCreateModalOpen.value = true
}

const onCategoryChangeInCreate = () => {
  switch (createForm.value.category) {
    case 'incident_broken_item':
      createForm.value.title = 'Broken Glass / Items in Room'
      createForm.value.priority = 'urgent'
      createForm.value.description = 'Guest called via Ext 0 reporting broken items, please sweep & replace.'
      break
    case 'extra_cleaning':
      createForm.value.title = 'Extra Cleaning Request'
      createForm.value.priority = 'high'
      createForm.value.description = 'Tamu meminta pembersihan tambahan / ganti sprei kamar.'
      break
    case 'amenities_request':
      createForm.value.title = 'Additional Towels & Amenities Request'
      createForm.value.priority = 'medium'
      createForm.value.description = 'Guest requested extra towels / soap for the room.'
      break
    case 'maintenance_repair':
      createForm.value.title = 'Facility Damage Report'
      createForm.value.priority = 'high'
      createForm.value.description = 'Guest reported facility damage in the room.'
      break
    default:
      createForm.value.title = 'Guest Service Request'
      createForm.value.priority = 'medium'
      createForm.value.description = ''
  }
}

const submitCreate = async () => {
  if (!createForm.value.room_id || !createForm.value.title || !createForm.value.description) {
    toastStore.error('Please complete all required fields.')
    return
  }

  try {
    loading.value = true
    const newReq = await serviceRequestService.create({
      room_id: createForm.value.room_id,
      guest_name: createForm.value.guest_name,
      guest_phone: createForm.value.guest_phone,
      category: createForm.value.category,
      title: createForm.value.title,
      description: `[Sumber: ${createForm.value.source === 'phone_ext0' ? 'Telepon Ext 0' : createForm.value.source === 'whatsapp' ? 'WhatsApp' : 'Walk-in'}] ${createForm.value.description}`,
      priority: createForm.value.priority,
    })

    // Jika resepsionis memilih langsung tugaskan ke staf housekeeping:
    if (createForm.value.direct_assign_housekeeper_id) {
      await serviceRequestService.assignToHousekeeping(newReq.id, {
        assigned_to_user_id: createForm.value.direct_assign_housekeeper_id,
        notes: createForm.value.description,
        priority: createForm.value.priority,
      })
      toastStore.success('Guest request recorded & assigned to Housekeeping via FCM!')
    } else {
      toastStore.success('Guest request recorded to queue!')
    }

    isCreateModalOpen.value = false
    await fetchRequests(false)
  } catch (error: any) {
    toastStore.error(error.message || 'Failed to record request')
  } finally {
    loading.value = false
  }
}

onUnmounted(() => {
  if (pollingTimer) clearInterval(pollingTimer)
})

const openAssignModal = (req: ServiceRequest) => {
  selectedRequest.value = req
  assignForm.value = {
    assigned_to_user_id: housekeepers.value[0]?.id || '',
    notes: req.category === 'incident_broken_item' ? 'Please clean the broken glass in the room and deliver replacement items.' : '',
    priority: req.priority || 'high',
  }
  isAssignModalOpen.value = true
}

const submitAssign = async () => {
  if (!selectedRequest.value) return
  try {
    loading.value = true
    await serviceRequestService.assignToHousekeeping(selectedRequest.value.id, assignForm.value)
    toastStore.success('Task forwarded to Housekeeping team!')
    isAssignModalOpen.value = false
    await fetchRequests(false)
  } catch (error: any) {
    toastStore.error(error.message || 'Failed to assign to Housekeeping')
  } finally {
    loading.value = false
  }
}

const openCompleteModal = (req: ServiceRequest) => {
  selectedRequest.value = req
  completeDamageCharge.value = 0
  isCompleteModalOpen.value = true
}

const submitComplete = async () => {
  if (!selectedRequest.value) return
  try {
    loading.value = true
    await serviceRequestService.complete(selectedRequest.value.id, { damage_charge: completeDamageCharge.value })
    toastStore.success('Task completed successfully!')
    isCompleteModalOpen.value = false
    await fetchRequests(false)
  } catch (error: any) {
    toastStore.error(error.message || 'Failed to complete task')
  } finally {
    loading.value = false
  }
}

const handleCancel = async (id: string) => {
  if (confirm('Cancel this room service report?')) {
    try {
      loading.value = true
      await serviceRequestService.cancel(id)
      toastStore.success('Report cancelled.')
      await fetchRequests(false)
    } catch (error: any) {
      toastStore.error(error.message || 'Failed to cancel report')
    } finally {
      loading.value = false
    }
  }
}

const selectRow = (req: ServiceRequest) => {
  selectedRequest.value = req
}

const closeDrawer = () => {
  selectedRequest.value = null
}

const getCategoryBadgeClass = (category: string) => {
  switch (category) {
    case 'incident_broken_item': return 'bg-rose-100 text-rose-700 border-rose-200'
    case 'extra_cleaning': return 'bg-amber-100 text-amber-700 border-amber-200'
    case 'amenities_request': return 'bg-blue-100 text-blue-700 border-blue-200'
    case 'maintenance_repair': return 'bg-purple-100 text-purple-700 border-purple-200'
    default: return 'bg-slate-100 text-slate-700 border-slate-200'
  }
}

const getCategoryLabel = (category: string) => {
  switch (category) {
    case 'incident_broken_item': return '🍷 Broken Items (Incident)'
    case 'extra_cleaning': return '🧹 Extra Cleaning'
    case 'amenities_request': return '🪥 Amenities / Supplies'
    case 'maintenance_repair': return '🛠️ Facility Damage'
    default: return '❓ Special Assistance'
  }
}
</script>

<template>
  <div class="service-requests-view">
    <div class="header-section">
      <div>
        <h1 class="page-title">🛎️ Room Service & Guest Requests</h1>
        <p class="page-subtitle">Central coordination for room incidents, cleaning requests, and Housekeeping assignments.</p>
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="stats-grid">
      <div class="stat-card">
        <span class="stat-icon">📋</span>
        <div>
          <h3>Total Requests</h3>
          <p class="main-val">{{ totalCount }}</p>
        </div>
      </div>
      <div class="stat-card">
        <span class="stat-icon text-amber-500">⏳</span>
        <div>
          <h3>Pending Reception</h3>
          <p class="main-val text-amber-600">{{ pendingCount }}</p>
        </div>
      </div>
      <div class="stat-card">
        <span class="stat-icon text-indigo-500">🧹</span>
        <div>
          <h3>In Progress</h3>
          <p class="main-val text-indigo-600">{{ assignedCount }}</p>
        </div>
      </div>
      <div class="stat-card">
        <span class="stat-icon text-emerald-500">✓</span>
        <div>
          <h3>Completed</h3>
          <p class="main-val text-emerald-600">{{ completedCount }}</p>
        </div>
      </div>
    </div>

    <!-- Main Table & Control Bar -->
    <div class="split-pane-layout">
      <div class="table-pane-box">
        <div class="control-bar">
          <div class="search-box">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search guest name, issue title, or room..."
              class="search-input"
            />
          </div>
          <div class="filter-box">
            <select v-model="statusFilter" class="filter-select">
              <option value="all">All Statuses</option>
              <option value="pending_reception">Pending Reception</option>
              <option value="assigned_to_housekeeping">Assigned to Housekeeping</option>
              <option value="completed">Completed</option>
              <option value="cancelled">Cancelled</option>
            </select>
          </div>
          <button @click="openCreateModal" class="btn btn-primary flex items-center gap-1.5 whitespace-nowrap shadow-sm">
            <span>📞</span>
            <span>+ Record Guest Report (Phone / Walk-in)</span>
          </button>
        </div>

        <div class="dashboard-card table-card">
          <div class="responsive-table-wrap">
            <table class="premium-table">
              <thead>
                <tr>
                  <th>Room Unit</th>
                  <th>Guest Name</th>
                  <th>Issue Category</th>
                  <th>Report Title</th>
                  <th class="text-center">Prioritas</th>
                  <th class="text-center">Status</th>
                  <th class="text-center">Petugas</th>
                  <th class="text-center">Aksi Cepat</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="req in filteredRequests"
                  :key="req.id"
                  @click="selectRow(req)"
                  :class="{ 'selected-row': selectedRequest?.id === req.id }"
                  class="clickable-row"
                >
                  <td class="font-bold text-slate-800">
                    {{ req.room?.name || 'N/A' }}
                    <span class="text-xs text-slate-400 block">({{ req.room?.code }})</span>
                  </td>
                  <td class="font-semibold text-slate-700">
                    {{ req.guest_name }}
                    <span v-if="req.guest_phone" class="text-xs text-slate-400 block">📞 {{ req.guest_phone }}</span>
                  </td>
                  <td>
                    <span class="px-2.5 py-1 text-xs font-bold rounded-full border" :class="getCategoryBadgeClass(req.category)">
                      {{ getCategoryLabel(req.category) }}
                    </span>
                  </td>
                  <td class="max-w-xs truncate">
                    <span class="font-bold text-slate-800">{{ req.title }}</span>
                    <p class="text-xs text-slate-500 truncate">{{ req.description }}</p>
                  </td>
                  <td class="text-center">
                    <span
                      class="px-2 py-0.5 text-xs font-bold rounded-full uppercase"
                      :class="{
                        'bg-rose-100 text-rose-700': req.priority === 'urgent',
                        'bg-amber-100 text-amber-700': req.priority === 'high',
                        'bg-blue-100 text-blue-700': req.priority === 'medium',
                        'bg-slate-100 text-slate-700': req.priority === 'low',
                      }"
                    >
                      {{ req.priority }}
                    </span>
                  </td>
                  <td class="text-center">
                    <span
                      v-if="req.status === 'pending_reception'"
                      class="px-2.5 py-1 text-xs font-bold rounded-full bg-amber-50 text-amber-700 border border-amber-300"
                    >
                      ⏳ Pending Frontdesk
                    </span>
                    <span
                      v-else-if="req.status === 'assigned_to_housekeeping' || req.status === 'in_progress'"
                      class="px-2.5 py-1 text-xs font-bold rounded-full bg-indigo-50 text-indigo-700 border border-indigo-300"
                    >
                      🧹 Di Housekeeping
                    </span>
                    <span
                      v-else-if="req.status === 'completed'"
                      class="px-2.5 py-1 text-xs font-bold rounded-full bg-emerald-50 text-emerald-700 border border-emerald-300"
                    >
                      ✓ Completed
                    </span>
                    <span v-else class="px-2.5 py-1 text-xs font-bold rounded-full bg-slate-100 text-slate-600">
                      {{ req.status }}
                    </span>
                  </td>
                  <td class="text-center text-xs font-medium text-slate-600">
                    {{ req.assigned_to_user ? `${req.assigned_to_user.first_name} ${req.assigned_to_user.last_name || ''}` : '-' }}
                  </td>
                  <td class="text-center" @click.stop>
                    <!-- Button: Assign to Housekeeping -->
                    <button
                      v-if="req.status === 'pending_reception'"
                      @click="openAssignModal(req)"
                      class="px-3 py-1 bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-xs rounded-lg shadow-sm transition-all flex items-center gap-1 mx-auto"
                      title="Forward task to Housekeeping"
                    >
                      <span>📋</span>
                      <span>Assign</span>
                    </button>

                    <!-- Button: Complete Task -->
                    <button
                      v-else-if="req.status === 'assigned_to_housekeeping' || req.status === 'in_progress'"
                      @click="openCompleteModal(req)"
                      class="px-3 py-1 bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs rounded-lg shadow-sm transition-all flex items-center gap-1 mx-auto"
                      title="Mark task as completed"
                    >
                      <span>✨</span>
                      <span>Completed</span>
                    </button>

                    <span v-else class="text-xs text-slate-400 font-semibold">Tuntas</span>
                  </td>
                </tr>
                <tr v-if="filteredRequests.length === 0">
                  <td colspan="8" class="no-data">No matching room service report found.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Detail Drawer Modal -->
      <transition name="drawer-slide">
        <div v-if="selectedRequest" class="action-drawer-pane">
          <div class="drawer-header">
            <div>
              <h3>Guest Report Details</h3>
              <span class="drawer-id">ID: #{{ selectedRequest.id.slice(-6).toUpperCase() }}</span>
            </div>
            <button @click="closeDrawer" class="close-drawer-btn">✕</button>
          </div>

          <div class="drawer-content">
            <div class="info-block-card">
              <label>Room & Guest</label>
              <p class="val-large">{{ selectedRequest.room?.name }} ({{ selectedRequest.room?.code }})</p>
              <p class="val-desc mt-1 font-semibold">👤 {{ selectedRequest.guest_name }} {{ selectedRequest.guest_phone ? `· 📞 ${selectedRequest.guest_phone}` : '' }}</p>
            </div>

            <div class="info-block-card">
              <label>Category & Priority</label>
              <div class="flex items-center gap-2 mt-1">
                <span class="px-2.5 py-1 text-xs font-bold rounded-full border" :class="getCategoryBadgeClass(selectedRequest.category)">
                  {{ getCategoryLabel(selectedRequest.category) }}
                </span>
                <span class="px-2 py-0.5 text-xs font-bold rounded-full bg-slate-100 text-slate-800 uppercase">
                  Prioritas: {{ selectedRequest.priority }}
                </span>
              </div>
            </div>

            <div class="info-block-card">
              <label>Issue Title & Description</label>
              <p class="val-mid font-bold text-slate-900">{{ selectedRequest.title }}</p>
              <p class="val-desc mt-1">{{ selectedRequest.description }}</p>
            </div>

            <div v-if="selectedRequest.notes_from_reception" class="info-block-card bg-indigo-50 border-indigo-100">
              <label class="text-indigo-900 font-bold">Catatan dari Resepsionis</label>
              <p class="val-desc text-indigo-800">{{ selectedRequest.notes_from_reception }}</p>
            </div>

            <div v-if="selectedRequest.damage_charge > 0" class="info-block-card bg-rose-50 border-rose-200">
              <label class="text-rose-900 font-bold">Biaya Ganti Rugi / Penggantian</label>
              <p class="val-price text-rose-700">Rp {{ selectedRequest.damage_charge.toLocaleString('id-ID') }}</p>
            </div>

            <div class="drawer-actions-container">
              <button
                v-if="selectedRequest.status === 'pending_reception'"
                @click="openAssignModal(selectedRequest)"
                class="btn btn-primary btn-block"
              >
                📋 Assign to Housekeeping
              </button>

              <button
                v-if="selectedRequest.status === 'assigned_to_housekeeping' || selectedRequest.status === 'in_progress'"
                @click="openCompleteModal(selectedRequest)"
                class="btn btn-checkedin btn-block"
              >
                ✨ Complete Task
              </button>

              <button
                v-if="selectedRequest.status !== 'completed' && selectedRequest.status !== 'cancelled'"
                @click="handleCancel(selectedRequest.id)"
                class="btn btn-danger-outline btn-block mt-2"
              >
                Cancel Report
              </button>
            </div>
          </div>
        </div>
      </transition>
    </div>

    <!-- Modal Form: Tugaskan ke Housekeeping -->
    <div v-if="isAssignModalOpen" class="modal-backdrop" @click.self="isAssignModalOpen = false">
      <div class="modal-dialog">
        <div class="modal-dialog-header">
          <h3>📋 Assign to Housekeeping Staff</h3>
          <button class="close-btn" @click="isAssignModalOpen = false">✕</button>
        </div>
        <form @submit.prevent="submitAssign" class="modal-dialog-body">
          <p class="text-xs text-slate-600 mb-3">
            Report for <strong>{{ selectedRequest?.room?.name }}</strong> ({{ selectedRequest?.title }}) will be forwarded via Push Notification to Housekeeping staff.
          </p>

          <div class="form-group mb-3">
            <label class="form-label">Select Housekeeping Staff</label>
            <select v-model="assignForm.assigned_to_user_id" class="form-input">
              <option value="">Semua Tim Housekeeping (Broadcast)</option>
              <option v-for="h in housekeepers" :key="h.id" :value="h.id">
                {{ h.first_name }} {{ h.last_name || '' }} ({{ h.role }})
              </option>
            </select>
          </div>

          <div class="form-group mb-3">
            <label class="form-label">Priority Level</label>
            <select v-model="assignForm.priority" class="form-input">
              <option value="urgent">Urgent (Immediate / Broken Glass)</option>
              <option value="high">High</option>
              <option value="medium">Medium (Standard)</option>
              <option value="low">Low</option>
            </select>
          </div>

          <div class="form-group mb-4">
            <label class="form-label">Additional Notes for Housekeeping</label>
            <textarea
              v-model="assignForm.notes"
              rows="3"
              class="form-input"
              placeholder="e.g. Please sweep broken glass under the table & deliver 1 new glass set..."
            ></textarea>
          </div>

          <div class="modal-dialog-footer">
            <button type="button" class="btn btn-secondary" @click="isAssignModalOpen = false">Cancel</button>
            <button type="submit" class="btn btn-primary" :disabled="loading">Send Assignment (FCM)</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Modal Form: Selesaikan Tugas (Complete Task) -->
    <div v-if="isCompleteModalOpen" class="modal-backdrop" @click.self="isCompleteModalOpen = false">
      <div class="modal-dialog">
        <div class="modal-dialog-header">
          <h3>✨ Complete Room Service Task</h3>
          <button class="close-btn" @click="isCompleteModalOpen = false">✕</button>
        </div>
        <form @submit.prevent="submitComplete" class="modal-dialog-body">
          <p class="text-xs text-slate-600 mb-3">
            Confirm that the issue at <strong>{{ selectedRequest?.room?.name }}</strong> has been resolved by the Housekeeping team.
          </p>

          <div class="form-group mb-4">
            <label class="form-label">Damage / Replacement Cost (IDR)</label>
            <input
              v-model.number="completeDamageCharge"
              type="number"
              min="0"
              step="5000"
              class="form-input"
              placeholder="0 if no damage cost"
            />
            <span class="text-xs text-slate-400 mt-1 block">Fill in if there is a replacement cost to be charged to the guest's bill.</span>
          </div>

          <div class="modal-dialog-footer">
            <button type="button" class="btn btn-secondary" @click="isCompleteModalOpen = false">Cancel</button>
            <button type="submit" class="btn btn-checkedin" :disabled="loading">✓ Mark Completed & Report to Frontdesk</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Modal Form: Catat Permintaan Tamu Baru (Input Resepsionis) -->
    <div v-if="isCreateModalOpen" class="modal-backdrop" @click.self="isCreateModalOpen = false">
      <div class="modal-dialog max-w-lg">
        <div class="modal-dialog-header">
          <h3>📞 Record Guest Report (Phone / Walk-in)</h3>
          <button class="close-btn" @click="isCreateModalOpen = false">✕</button>
        </div>
        <form @submit.prevent="submitCreate" class="modal-dialog-body">
          <div class="grid grid-cols-2 gap-3 mb-3">
            <div class="form-group">
              <label class="form-label">Room Unit *</label>
              <select v-model="createForm.room_id" class="form-input" required>
                <option v-for="r in rooms" :key="r.id" :value="r.id">
                  {{ r.name }} ({{ r.code }})
                </option>
              </select>
            </div>
            <div class="form-group">
              <label class="form-label">Report Source</label>
              <select v-model="createForm.source" class="form-input">
                <option value="phone_ext0">📞 Room Phone (Ext 0)</option>
                <option value="whatsapp">💬 Frontdesk WhatsApp</option>
                <option value="walk_in">🚶 Walk-in (Came to Lobby)</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3 mb-3">
            <div class="form-group">
              <label class="form-label">Guest / Reporter Name</label>
              <input
                v-model="createForm.guest_name"
                type="text"
                class="form-input"
                placeholder="e.g. John Doe"
                required
              />
            </div>
            <div class="form-group">
              <label class="form-label">Issue / Incident Type *</label>
              <select v-model="createForm.category" @change="onCategoryChangeInCreate" class="form-input" required>
                <option value="incident_broken_item">🍷 Broken Items / Glass</option>
                <option value="extra_cleaning">🧹 Extra Cleaning / Spills</option>
                <option value="amenities_request">🪥 Additional Towels & Amenities</option>
                <option value="maintenance_repair">🛠️ Facility Damage (AC/Lights)</option>
                <option value="other">❓ Special Assistance</option>
              </select>
            </div>
          </div>

          <div class="form-group mb-3">
            <label class="form-label">Request Title *</label>
            <input
              v-model="createForm.title"
              type="text"
              class="form-input"
              placeholder="e.g. Broken glass on the room floor"
              required
            />
          </div>

          <div class="form-group mb-3">
            <label class="form-label">Cleaning Instructions / Items to Deliver *</label>
            <textarea
              v-model="createForm.description"
              rows="3"
              class="form-input"
              placeholder="Describe the cleaning needs or replacement items to be delivered..."
              required
            ></textarea>
          </div>

          <div class="p-3 bg-indigo-50 border border-indigo-100 rounded-xl mb-3">
            <label class="form-label text-indigo-900 font-bold mb-1">🚀 Directly Forward to Housekeeping (Optional)</label>
            <select v-model="createForm.direct_assign_housekeeper_id" class="form-input bg-white">
              <option value="">Save to Queue Only (Unassigned)</option>
              <option v-for="h in housekeepers" :key="h.id" :value="h.id">
                Send Push Notification (FCM) to: {{ h.first_name }} {{ h.last_name || '' }} ({{ h.role }})
              </option>
            </select>
          </div>

          <div class="modal-dialog-footer">
            <button type="button" class="btn btn-secondary" @click="isCreateModalOpen = false">Cancel</button>
            <button type="submit" class="btn btn-primary" :disabled="loading">
              {{ createForm.direct_assign_housekeeper_id ? 'Save & Send Push to Housekeeping' : 'Save Report' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.service-requests-view {
  padding: 1.5rem;
  max-width: 1600px;
  margin: 0 auto;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 800;
  color: #0f172a;
}

.page-subtitle {
  font-size: 0.85rem;
  color: #64748b;
  margin-top: 2px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1rem;
  margin: 1.5rem 0;
}

.stat-card {
  background: white;
  padding: 1.25rem;
  border-radius: 1rem;
  border: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.stat-icon {
  font-size: 1.8rem;
}

.stat-card h3 {
  font-size: 0.8rem;
  color: #64748b;
  margin: 0;
  font-weight: 600;
}

.main-val {
  font-size: 1.4rem;
  font-weight: 800;
  margin: 0;
  color: #0f172a;
}

.split-pane-layout {
  display: flex;
  gap: 1.5rem;
  position: relative;
}

.table-pane-box {
  flex: 1;
  min-width: 0;
}

.control-bar {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}

.search-box {
  flex: 1;
}

.search-input, .filter-select {
  width: 100%;
  padding: 0.65rem 1rem;
  border: 1px solid #cbd5e1;
  border-radius: 0.75rem;
  font-size: 0.85rem;
  background: white;
  outline: none;
}

.dashboard-card {
  background: white;
  border-radius: 1rem;
  border: 1px solid #e2e8f0;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.premium-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.85rem;
}

.premium-table th {
  background: #f8fafc;
  padding: 0.85rem 1rem;
  font-weight: 700;
  color: #475569;
  border-bottom: 1px solid #e2e8f0;
}

.premium-table td {
  padding: 0.85rem 1rem;
  border-bottom: 1px solid #f1f5f9;
}

.clickable-row {
  cursor: pointer;
  transition: background 0.15s;
}

.clickable-row:hover {
  background: #f8fafc;
}

.selected-row {
  background: #eef2ff !important;
}

.no-data {
  text-align: center;
  padding: 3rem;
  color: #94a3b8;
  font-weight: 600;
}

/* Action Drawer */
.action-drawer-pane {
  width: 380px;
  background: white;
  border-radius: 1rem;
  border: 1px solid #e2e8f0;
  padding: 1.5rem;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  gap: 1rem;
  max-height: calc(100vh - 180px);
  overflow-y: auto;
}

.drawer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e2e8f0;
}

.drawer-header h3 {
  font-size: 1.1rem;
  font-weight: 800;
  margin: 0;
  color: #0f172a;
}

.drawer-id {
  font-size: 0.75rem;
  color: #64748b;
  font-weight: 600;
}

.close-drawer-btn {
  background: none;
  border: none;
  font-size: 1.1rem;
  cursor: pointer;
  color: #94a3b8;
}

.info-block-card {
  background: #f8fafc;
  padding: 0.85rem;
  border-radius: 0.75rem;
  border: 1px solid #e2e8f0;
}

.info-block-card label {
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  color: #64748b;
  display: block;
}

.val-large {
  font-size: 1.1rem;
  font-weight: 800;
  color: #0f172a;
  margin: 2px 0 0 0;
}

.val-desc {
  font-size: 0.8rem;
  color: #475569;
  line-height: 1.4;
  margin: 0;
}

.val-price {
  font-size: 1.25rem;
  font-weight: 800;
  margin: 0;
}

.drawer-actions-container {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.btn {
  padding: 0.65rem 1rem;
  border-radius: 0.6rem;
  font-weight: 700;
  font-size: 0.85rem;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-block {
  width: 100%;
}

.btn-primary {
  background: #6366f1;
  color: white;
}

.btn-primary:hover {
  background: #4f46e5;
}

.btn-checkedin {
  background: #059669;
  color: white;
}

.btn-checkedin:hover {
  background: #047857;
}

.btn-secondary {
  background: #f1f5f9;
  color: #475569;
}

.btn-danger-outline {
  background: transparent;
  border: 1px solid #fecdd3;
  color: #e11d48;
}

.btn-danger-outline:hover {
  background: #fff1f2;
}

/* Modals */
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(4px);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.modal-dialog {
  background: white;
  width: 100%;
  max-width: 500px;
  border-radius: 1rem;
  overflow: hidden;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

.modal-dialog-header {
  padding: 1.25rem 1.5rem;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-dialog-header h3 {
  font-size: 1rem;
  font-weight: 700;
  margin: 0;
  color: #0f172a;
}

.modal-dialog-body {
  padding: 1.5rem;
}

.form-label {
  display: block;
  font-size: 0.8rem;
  font-weight: 700;
  color: #334155;
  margin-bottom: 0.35rem;
}

.form-input {
  width: 100%;
  padding: 0.65rem 0.85rem;
  border: 1.5px solid #cbd5e1;
  border-radius: 0.6rem;
  font-size: 0.85rem;
  outline: none;
}

.form-input:focus {
  border-color: #6366f1;
}

.modal-dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1rem;
}
</style>
