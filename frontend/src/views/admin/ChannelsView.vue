<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useChannelStore } from '@/stores/channelStore'
import { useToastStore } from '@/stores/toastStore'
import type { Channel } from '@/types'

const channelStore = useChannelStore()
const toastStore = useToastStore()

const searchQuery = ref('')
const isModalOpen = ref(false)
const isEditing = ref(false)
const editingId = ref('')
const submitting = ref(false)
const deleteModalOpen = ref(false)
const channelToDelete = ref<Channel | null>(null)

const form = ref({
  name: '',
})

// Suggested popular channels for quick selection
const popularChannels = [
  { name: 'Traveloka', tag: 'OTA' },
  { name: 'Booking.com', tag: 'OTA' },
  { name: 'Agoda', tag: 'OTA' },
  { name: 'Tiket.com', tag: 'OTA' },
  { name: 'Direct Website', tag: 'Direct' },
  { name: 'Walk-In / Front Desk', tag: 'Offline' },
  { name: 'Airbnb', tag: 'OTA' },
  { name: 'Corporate Partner', tag: 'B2B' },
]

onMounted(async () => {
  await loadChannels()
})

const loadChannels = async () => {
  try {
    await channelStore.fetchChannels()
  } catch (error: any) {
    console.error('Error loading channels:', error)
  }
}

// Computed Filtered List
const filteredChannels = computed(() => {
  if (!searchQuery.value) return channelStore.channels
  const q = searchQuery.value.toLowerCase().trim()
  return channelStore.channels.filter(
    (c) =>
      (c.name || '').toLowerCase().includes(q) ||
      (c.code || '').toLowerCase().includes(q) ||
      (c.id || '').toLowerCase().includes(q),
  )
})

// KPI Stats
const totalChannels = computed(() => channelStore.channels.length)
const latestChannel = computed(() => {
  if (!channelStore.channels.length) return null
  return [...channelStore.channels].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
  )[0]
})

// Quick helper to categorize channel
const getChannelCategory = (name: string) => {
  const n = (name || '').toLowerCase()
  if (n.includes('walk-in') || n.includes('offline') || n.includes('resepsionis') || n.includes('front desk')) {
    return { label: 'Offline / Front Desk', color: 'slate' }
  }
  if (n.includes('direct') || n.includes('website') || n.includes('web') || n.includes('whatsapp') || n.includes('telepon')) {
    return { label: 'Direct Booking', color: 'emerald' }
  }
  if (n.includes('corporate') || n.includes('b2b') || n.includes('agent')) {
    return { label: 'Corporate / B2B', color: 'indigo' }
  }
  return { label: 'Online Travel Agent (OTA)', color: 'blue' }
}

const getChannelIcon = (name: string) => {
  const n = (name || '').toLowerCase()
  if (n.includes('walk-in') || n.includes('offline')) return '🏨'
  if (n.includes('website') || n.includes('direct') || n.includes('web')) return '🌐'
  if (n.includes('traveloka')) return '🐦'
  if (n.includes('booking')) return '🅱️'
  if (n.includes('agoda')) return '🟢'
  if (n.includes('tiket')) return '🟡'
  if (n.includes('airbnb')) return '🏠'
  return '📡'
}

const formatDate = (dateStr?: string) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// Modal Actions
const openCreateModal = () => {
  isEditing.value = false
  editingId.value = ''
  form.value = { name: '' }
  isModalOpen.value = true
}

const openEditModal = (channel: Channel) => {
  isEditing.value = true
  editingId.value = channel.id
  form.value = { name: channel.name }
  isModalOpen.value = true
}

const selectSuggestion = (name: string) => {
  form.value.name = name
}

const closeModal = () => {
  isModalOpen.value = false
  form.value = { name: '' }
}

const saveChannel = async () => {
  const trimmedName = form.value.name.trim()
  if (!trimmedName) {
    toastStore.warning('Channel name cannot be empty!')
    return
  }

  try {
    submitting.value = true
    if (isEditing.value) {
      const msg = await channelStore.update(editingId.value, { name: trimmedName })
      toastStore.success(msg || 'Channel updated successfully!')
    } else {
      const msg = await channelStore.store({ name: trimmedName })
      toastStore.success(msg || 'Channel added successfully!')
    }
    await channelStore.fetchChannels()
    closeModal()
  } catch (err: any) {
    const errorMsg =
      err.response?.data?.message || err.message || 'Failed to save channel'
    toastStore.error(errorMsg)
  } finally {
    submitting.value = false
  }
}

// Delete Flow
const confirmDelete = (channel: Channel) => {
  channelToDelete.value = channel
  deleteModalOpen.value = true
}

const closeDeleteModal = () => {
  deleteModalOpen.value = false
  channelToDelete.value = null
}

const executeDelete = async () => {
  if (!channelToDelete.value) return
  try {
    submitting.value = true
    const msg = await channelStore.destroy(channelToDelete.value.id)
    toastStore.success(msg || 'Channel deleted successfully!')
    await channelStore.fetchChannels()
    closeDeleteModal()
  } catch (err: any) {
    const errorMsg =
      err.response?.data?.message || err.message || 'Failed to delete channel'
    toastStore.error(errorMsg)
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="channels-view">
    <!-- Header Section -->
    <div class="page-header">
      <div class="header-content">
        <h2 class="title">Manage Channels</h2>
        <p class="subtitle">
          Manage booking distribution channel names (Online Travel Agents, Direct Website, Walk-in, etc).
        </p>
      </div>
      <button @click="openCreateModal" class="btn btn-primary btn-add">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <line x1="12" y1="5" x2="12" y2="19"></line>
          <line x1="5" y1="12" x2="19" y2="12"></line>
        </svg>
        <span>Add Channel</span>
      </button>
    </div>

    <!-- Stats KPI Cards -->
    <div class="kpi-grid">
      <div class="kpi-card">
        <div class="kpi-icon-box bg-orange">
          <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="2"></circle>
            <path d="M16.24 7.76a6 6 0 0 1 0 8.49m-8.48-.01a6 6 0 0 1 0-8.49m11.31-2.82a10 10 0 0 1 0 14.14m-14.14 0a10 10 0 0 1 0-14.14"></path>
          </svg>
        </div>
        <div class="kpi-info">
          <span class="kpi-label">Total Active Channels</span>
          <span class="kpi-value">{{ totalChannels }}</span>
        </div>
      </div>

      <div class="kpi-card">
        <div class="kpi-icon-box bg-blue">
          <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="2" y1="12" x2="22" y2="12"></line>
            <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path>
          </svg>
        </div>
        <div class="kpi-info">
          <span class="kpi-label">Latest Channel</span>
          <span class="kpi-value-sm">{{ latestChannel ? latestChannel.name : '-' }}</span>
        </div>
      </div>

      <div class="kpi-card">
        <div class="kpi-icon-box bg-emerald">
          <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
            <polyline points="22 4 12 14.01 9 11.01"></polyline>
          </svg>
        </div>
        <div class="kpi-info">
          <span class="kpi-label">Latest Code</span>
          <span class="kpi-value-sm font-mono">{{ latestChannel ? latestChannel.code : 'C-1' }}</span>
        </div>
      </div>
    </div>

    <!-- Filter & Control Bar -->
    <div class="control-bar">
      <div class="search-box">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="search-icon">
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
        </svg>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search by channel name or code (e.g.: Traveloka, C-1)..."
          class="search-input"
        />
        <button v-if="searchQuery" @click="searchQuery = ''" class="clear-search-btn" title="Clear">
          &times;
        </button>
      </div>

      <button @click="loadChannels" class="btn btn-outline refresh-btn" :disabled="channelStore.loading" title="Refresh data">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" :class="{ 'spin-anim': channelStore.loading }">
          <polyline points="23 4 23 10 17 10"></polyline>
          <polyline points="1 20 1 14 7 14"></polyline>
          <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path>
        </svg>
        <span>Refresh</span>
      </button>
    </div>

    <!-- Table Card -->
    <div class="table-card">
      <div v-if="channelStore.loading && channelStore.channels.length === 0" class="loading-state">
        <div class="spinner"></div>
        <p>Loading channel list...</p>
      </div>

      <div v-else class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th style="width: 100px;">Code</th>
              <th>Channel Name</th>
              <th>Channel Category</th>
              <th>Created At</th>
              <th style="text-align: right; width: 140px;">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="channel in filteredChannels" :key="channel.id" class="table-row">
              <td>
                <div class="code-badge">
                  {{ channel.code }}
                </div>
              </td>
              <td>
                <div class="channel-name-cell">
                  <span class="channel-avatar">{{ getChannelIcon(channel.name) }}</span>
                  <span class="channel-name-text">{{ channel.name }}</span>
                </div>
              </td>
              <td>
                <span class="type-tag" :class="`type-${getChannelCategory(channel.name).color}`">
                  {{ getChannelCategory(channel.name).label }}
                </span>
              </td>
              <td class="date-cell">
                {{ formatDate(channel.created_at) }}
              </td>
              <td>
                <div class="action-buttons">
                  <button
                    @click="openEditModal(channel)"
                    class="action-btn edit-btn"
                    title="Edit Channel"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                      <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                    </svg>
                    <span>Edit</span>
                  </button>
                  <button
                    @click="confirmDelete(channel)"
                    class="action-btn delete-btn"
                    title="Delete Channel"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <polyline points="3 6 5 6 21 6"></polyline>
                      <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                    </svg>
                    <span>Hapus</span>
                  </button>
                </div>
              </td>
            </tr>

            <tr v-if="filteredChannels.length === 0 && !channelStore.loading">
              <td colspan="5" class="empty-state">
                <div class="empty-content">
                  <div class="empty-icon">📡</div>
                  <h4>No channels found</h4>
                  <p v-if="searchQuery">No channel matches the search "{{ searchQuery }}".</p>
                  <p v-else>No channel data registered yet.</p>
                  <button v-if="!searchQuery" @click="openCreateModal" class="btn btn-primary btn-sm mt-3">
                    + Add First Channel
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create / Edit Modal -->
    <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
      <div class="modal-card">
        <div class="modal-header">
          <div class="modal-title-wrap">
            <div class="modal-icon-badge">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="2"></circle>
                <path d="M16.24 7.76a6 6 0 0 1 0 8.49m-8.48-.01a6 6 0 0 1 0-8.49m11.31-2.82a10 10 0 0 1 0 14.14m-14.14 0a10 10 0 0 1 0-14.14"></path>
              </svg>
            </div>
            <div>
              <h3>{{ isEditing ? 'Edit Channel' : 'Add New Channel' }}</h3>
              <p class="modal-subtitle">Manage booking distribution channel names</p>
            </div>
          </div>
          <button @click="closeModal" class="close-btn" aria-label="Close modal">&times;</button>
        </div>

        <div class="modal-body">
          <!-- Form Input -->
          <div class="form-group">
            <label class="form-label">
              Channel Name <span class="required">*</span>
            </label>
            <input
              v-model="form.name"
              type="text"
              placeholder="e.g.: Traveloka, Booking.com, Direct Website..."
              class="form-input"
              @keyup.enter="saveChannel"
              autofocus
            />
          </div>

          <!-- Quick Suggestions Chips (only in create mode) -->
          <div v-if="!isEditing" class="suggestions-section">
            <span class="suggestions-label">Quick Picks:</span>
            <div class="chips-container">
              <button
                v-for="item in popularChannels"
                :key="item.name"
                type="button"
                @click="selectSuggestion(item.name)"
                class="chip-btn"
                :class="{ active: form.name === item.name }"
              >
                {{ item.name }}
              </button>
            </div>
          </div>

          <!-- Live Preview Card -->
          <div v-if="form.name.trim()" class="preview-box">
            <span class="preview-title">Preview:</span>
            <div class="preview-item">
              <span class="channel-avatar">{{ getChannelIcon(form.name) }}</span>
              <div class="preview-details">
                <span class="preview-name">{{ form.name }}</span>
                <span class="type-tag-sm" :class="`type-${getChannelCategory(form.name).color}`">
                  {{ getChannelCategory(form.name).label }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button @click="closeModal" class="btn btn-secondary" :disabled="submitting">
            Cancel
          </button>
          <button @click="saveChannel" class="btn btn-primary" :disabled="submitting">
            <span v-if="submitting">Saving...</span>
            <span v-else>{{ isEditing ? 'Save Changes' : 'Add Channel' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="deleteModalOpen" class="modal-overlay" @click.self="closeDeleteModal">
      <div class="modal-card modal-delete">
        <div class="delete-icon-wrap">
          <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#ef4444" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
            <line x1="12" y1="9" x2="12" y2="13"></line>
            <line x1="12" y1="17" x2="12.01" y2="17"></line>
          </svg>
        </div>
        <h3 class="delete-title">Delete Channel?</h3>
        <p class="delete-desc">
          Are you sure you want to delete this channel
          <strong>"{{ channelToDelete?.name }}"</strong> ({{ channelToDelete?.code }})?
          This action cannot be undone.
        </p>

        <div class="modal-footer delete-footer">
          <button @click="closeDeleteModal" class="btn btn-secondary" :disabled="submitting">
            Cancel
          </button>
          <button @click="executeDelete" class="btn btn-danger" :disabled="submitting">
            <span v-if="submitting">Deleting...</span>
            <span v-else>Yes, Delete</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.channels-view {
  display: flex;
  flex-direction: column;
  gap: 24px;
  max-width: 1320px;
  margin: 0 auto;
}

/* Page Header */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.title {
  font-size: 1.5rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
  letter-spacing: -0.02em;
}

.subtitle {
  font-size: 0.875rem;
  color: #64748b;
  margin: 4px 0 0 0;
}

.btn-add {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  font-weight: 600;
  box-shadow: 0 4px 12px rgba(225, 91, 43, 0.25);
  transition: all 0.2s;
}

.btn-add:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(225, 91, 43, 0.35);
}

/* KPI Grid */
.kpi-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
}

.kpi-card {
  background: #ffffff;
  border-radius: 16px;
  padding: 18px 20px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  display: flex;
  align-items: center;
  gap: 16px;
  transition: transform 0.2s, box-shadow 0.2s;
}

.kpi-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.06);
}

.kpi-icon-box {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.bg-orange {
  background: rgba(225, 91, 43, 0.1);
  color: #e15b2b;
}

.bg-blue {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
}

.bg-emerald {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
}

.kpi-info {
  display: flex;
  flex-direction: column;
}

.kpi-label {
  font-size: 0.8rem;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.kpi-value {
  font-size: 1.6rem;
  font-weight: 800;
  color: #0f172a;
  line-height: 1.2;
}

.kpi-value-sm {
  font-size: 1.15rem;
  font-weight: 700;
  color: #0f172a;
  line-height: 1.2;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180px;
}

.font-mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

/* Control Bar */
.control-bar {
  display: flex;
  gap: 12px;
  background: #ffffff;
  padding: 14px 18px;
  border-radius: 16px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  align-items: center;
}

.search-box {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 14px;
  color: #94a3b8;
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 10px 38px 10px 42px;
  border-radius: 10px;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
  outline: none;
  font-size: 0.875rem;
  color: #1e293b;
  transition: all 0.2s;
}

.search-input:focus {
  background: #ffffff;
  border-color: #e15b2b;
  box-shadow: 0 0 0 3px rgba(225, 91, 43, 0.12);
}

.clear-search-btn {
  position: absolute;
  right: 12px;
  background: none;
  border: none;
  font-size: 1.2rem;
  color: #94a3b8;
  cursor: pointer;
  padding: 2px 6px;
  border-radius: 4px;
}

.clear-search-btn:hover {
  color: #0f172a;
}

.refresh-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 10px 14px;
  font-size: 0.85rem;
}

.spin-anim {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  100% {
    transform: rotate(360deg);
  }
}

/* Table Card */
.table-card {
  background: #ffffff;
  border-radius: 16px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.table-container {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.data-table th {
  padding: 14px 20px;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.data-table td {
  padding: 16px 20px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 0.875rem;
  color: #334155;
  vertical-align: middle;
}

.table-row:hover {
  background: #fcfcfd;
}

.code-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 8px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.8rem;
  font-weight: 700;
  background: #f1f5f9;
  color: #0f172a;
  border: 1px solid #e2e8f0;
}

.channel-name-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.channel-avatar {
  font-size: 1.25rem;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px solid #f1f5f9;
}

.channel-name-text {
  font-weight: 700;
  color: #0f172a;
  font-size: 0.925rem;
}

.type-tag {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 600;
}

.type-blue {
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #dbeafe;
}

.type-emerald {
  background: #ecfdf5;
  color: #059669;
  border: 1px solid #d1fae5;
}

.type-slate {
  background: #f1f5f9;
  color: #475569;
  border: 1px solid #e2e8f0;
}

.type-indigo {
  background: #eef2ff;
  color: #4f46e5;
  border: 1px solid #e0e7ff;
}

.date-cell {
  color: #64748b;
  font-size: 0.8rem;
  white-space: nowrap;
}

.action-buttons {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 0.78rem;
  font-weight: 600;
  border: 1px solid transparent;
  cursor: pointer;
  transition: all 0.15s;
}

.edit-btn {
  background: #f8fafc;
  color: #3b82f6;
  border-color: #dbeafe;
}

.edit-btn:hover {
  background: #eff6ff;
  border-color: #bfdbfe;
}

.delete-btn {
  background: #fff5f5;
  color: #ef4444;
  border-color: #fee2e2;
}

.delete-btn:hover {
  background: #fee2e2;
  border-color: #fca5a5;
}

/* Loading & Empty State */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  gap: 12px;
  color: #64748b;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #e2e8f0;
  border-top-color: #e15b2b;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.empty-state {
  text-align: center;
  padding: 60px 20px !important;
}

.empty-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.empty-icon {
  font-size: 2.5rem;
  margin-bottom: 4px;
}

.empty-content h4 {
  margin: 0;
  font-size: 1.1rem;
  color: #0f172a;
}

.empty-content p {
  margin: 0;
  color: #64748b;
  font-size: 0.875rem;
}

.mt-3 {
  margin-top: 12px;
}

/* Buttons */
.btn {
  padding: 8px 16px;
  border-radius: 10px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-sm {
  padding: 6px 12px;
  font-size: 0.8rem;
}

.btn-primary {
  background-color: #e15b2b;
  color: #ffffff;
}

.btn-primary:hover:not(:disabled) {
  background-color: #c84e20;
}

.btn-secondary {
  background: #ffffff;
  border-color: #cbd5e1;
  color: #475569;
}

.btn-secondary:hover:not(:disabled) {
  background: #f8fafc;
}

.btn-outline {
  background: #ffffff;
  border-color: #e2e8f0;
  color: #475569;
}

.btn-outline:hover:not(:disabled) {
  background: #f8fafc;
  color: #0f172a;
}

.btn-danger {
  background: #ef4444;
  color: #ffffff;
}

.btn-danger:hover:not(:disabled) {
  background: #dc2626;
}

/* Modals */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(15, 23, 42, 0.4);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  animation: fadeIn 0.2s ease;
}

.modal-card {
  background: #ffffff;
  border-radius: 20px;
  width: 90%;
  max-width: 480px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  animation: slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-header {
  padding: 20px 24px;
  border-bottom: 1px solid #f1f5f9;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
}

.modal-title-wrap {
  display: flex;
  align-items: center;
  gap: 12px;
}

.modal-icon-badge {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: rgba(225, 91, 43, 0.1);
  color: #e15b2b;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.15rem;
  font-weight: 700;
  color: #0f172a;
}

.modal-subtitle {
  margin: 2px 0 0 0;
  font-size: 0.8rem;
  color: #64748b;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #94a3b8;
  cursor: pointer;
  padding: 0;
  line-height: 1;
}

.close-btn:hover {
  color: #0f172a;
}

.modal-body {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-size: 0.75rem;
  font-weight: 700;
  color: #475569;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.required {
  color: #ef4444;
}

.form-input {
  padding: 12px 16px;
  border-radius: 10px;
  border: 1px solid #cbd5e1;
  outline: none;
  font-size: 0.925rem;
  color: #1e293b;
  background: #ffffff;
  box-sizing: border-box;
  transition: all 0.2s;
}

.form-input:focus {
  border-color: #e15b2b;
  box-shadow: 0 0 0 3px rgba(225, 91, 43, 0.12);
}

.suggestions-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.suggestions-label {
  font-size: 0.75rem;
  font-weight: 600;
  color: #64748b;
}

.chips-container {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.chip-btn {
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
  color: #334155;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 9999px;
  cursor: pointer;
  transition: all 0.15s;
}

.chip-btn:hover {
  background: #e2e8f0;
  color: #0f172a;
}

.chip-btn.active {
  background: rgba(225, 91, 43, 0.12);
  border-color: #e15b2b;
  color: #e15b2b;
}

.preview-box {
  background: #f8fafc;
  border: 1px dashed #cbd5e1;
  border-radius: 12px;
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.preview-title {
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  color: #94a3b8;
  letter-spacing: 0.05em;
}

.preview-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.preview-details {
  display: flex;
  align-items: center;
  gap: 8px;
}

.preview-name {
  font-weight: 700;
  color: #0f172a;
  font-size: 0.9rem;
}

.type-tag-sm {
  font-size: 0.7rem;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 12px;
}

.modal-footer {
  padding: 16px 24px;
  background: #f8fafc;
  border-top: 1px solid #f1f5f9;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

/* Delete Modal */
.modal-delete {
  max-width: 400px;
  text-align: center;
  padding: 32px 24px 20px 24px;
}

.delete-icon-wrap {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: #fef2f2;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px auto;
}

.delete-title {
  margin: 0 0 8px 0;
  font-size: 1.25rem;
  font-weight: 800;
  color: #0f172a;
}

.delete-desc {
  margin: 0 0 20px 0;
  font-size: 0.875rem;
  color: #64748b;
  line-height: 1.5;
}

.delete-footer {
  background: transparent;
  border-top: none;
  padding: 0;
  justify-content: center;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
