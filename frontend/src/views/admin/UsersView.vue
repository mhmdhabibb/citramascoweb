<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { userService } from '@/services/admin/userService'
import { useToastStore } from '@/stores/toastStore'
import type { User, UserRole } from '@/types'

const users = ref<User[]>([])
const searchQuery = ref('')
const roleFilter = ref('All')
const loading = ref(false)
const toastStore = useToastStore()

const roleLabels: Record<UserRole | 'All', string> = {
  All: 'All Roles',
  admin: 'Admin',
  manager: 'Manager',
  user: 'Customer',
  reception: 'Reception',
  finance: 'Finance',
  inventory: 'Inventory',
}

const roleOptions: Array<{ value: UserRole | 'All'; label: string }> = [
  { value: 'All', label: 'All Roles' },
  { value: 'user', label: 'Customer' },
  { value: 'admin', label: 'Admin' },
  { value: 'manager', label: 'Manager' },
  { value: 'reception', label: 'Reception' },
  { value: 'finance', label: 'Finance' },
  { value: 'inventory', label: 'Inventory' },
]

// --- Stats ---
const totalUsers = computed(() => users.value.length)
const customerCount = computed(() => users.value.filter((u) => u.role === 'user').length)
const adminCount = computed(() => users.value.filter((u) => u.role === 'admin').length)
const managerCount = computed(() => users.value.filter((u) => u.role === 'manager').length)

// --- Filter & Search ---
const filteredUsers = computed(() => {
  return users.value.filter((user) => {
    const fullName = `${user.first_name} ${user.last_name}`.toLowerCase()
    const matchesSearch =
      fullName.includes(searchQuery.value.toLowerCase()) ||
      (user.email || '').toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      (user.username || '').toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      (user.id || '').toLowerCase().includes(searchQuery.value.toLowerCase())

    const matchesRole = roleFilter.value === 'All' || user.role === roleFilter.value
    return matchesSearch && matchesRole
  })
})

const fetchUsers = async () => {
  try {
    loading.value = true
    if (roleFilter.value === 'All') {
      const roles: UserRole[] = ['user', 'admin', 'manager', 'reception', 'finance', 'inventory']
      const results = await Promise.all(roles.map((role) => userService.getByRole(role)))
      const merged = new Map<string, User>()
      for (const list of results) {
        for (const user of list) {
          merged.set(user.id, user)
        }
      }
      users.value = Array.from(merged.values())
    } else {
      users.value = await userService.getByRole(roleFilter.value as UserRole)
    }
  } catch (error) {
    toastStore.error((error as Error).message || 'Failed to fetch users')
  } finally {
    loading.value = false
  }
}

const handleDelete = async (user: User) => {
  const fullName = `${user.first_name} ${user.last_name}`.trim() || user.username
  if (!confirm(`Are you sure you want to delete user "${fullName}"?`)) return

  try {
    loading.value = true
    const msg = await userService.delete(user.id)
    toastStore.success(msg || 'User deleted successfully')
    await fetchUsers()
  } catch (error) {
    toastStore.error((error as Error).message || 'Failed to delete user')
  } finally {
    loading.value = false
  }
}

const formatDate = (date?: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('en-GB', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  })
}

const initials = (user: User) => {
  const first = user.first_name?.[0] || ''
  const last = user.last_name?.[0] || ''
  return (first + last).toUpperCase() || (user.username?.[0] || '?').toUpperCase()
}

// --- Create User Modal ---
const showCreateModal = ref(false)
const creating = ref(false)
const createForm = ref({
  first_name: '',
  last_name: '',
  username: '',
  password: '',
  phone: '',
  email: '',
  address: '',
  role: 'user' as UserRole,
})

const creatableRoles = computed(() => roleOptions.filter((o) => o.value !== 'All'))

const openCreateModal = () => {
  createForm.value = {
    first_name: '',
    last_name: '',
    username: '',
    password: '',
    phone: '',
    email: '',
    address: '',
    role: 'user',
  }
  showCreateModal.value = true
}

const submitCreate = async () => {
  const f = createForm.value
  if (
    !f.first_name ||
    !f.last_name ||
    !f.username ||
    !f.password ||
    !f.phone ||
    !f.email ||
    !f.role
  ) {
    toastStore.error('Please fill in all required fields')
    return
  }

  try {
    creating.value = true
    await userService.create(f)
    toastStore.success('User created successfully')
    showCreateModal.value = false
    await fetchUsers()
  } catch (error) {
    toastStore.error((error as Error).message || 'Failed to create user')
  } finally {
    creating.value = false
  }
}

// --- Assign Role Modal ---
const showRoleModal = ref(false)
const roleUser = ref<User | null>(null)
const selectedRole = ref<UserRole>('user')
const updatingRole = ref(false)

const openRoleModal = (user: User) => {
  roleUser.value = user
  selectedRole.value = user.role
  showRoleModal.value = true
}

const submitRole = async () => {
  if (!roleUser.value) return
  if (roleUser.value.role === selectedRole.value) {
    showRoleModal.value = false
    return
  }

  try {
    updatingRole.value = true
    const msg = await userService.updateRole(roleUser.value.id, selectedRole.value)
    toastStore.success(msg || 'Role updated successfully')
    showRoleModal.value = false
    await fetchUsers()
  } catch (error) {
    toastStore.error((error as Error).message || 'Failed to update role')
  } finally {
    updatingRole.value = false
  }
}

onMounted(fetchUsers)
</script>

<template>
  <div class="users-view">
    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon stat-icon-primary">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" /><circle cx="9" cy="7" r="4" /><path d="M22 21v-2a4 4 0 0 0-3-3.87" /><path d="M16 3.13a4 4 0 0 1 0 7.75" /></svg>
        </div>
        <div class="stat-info">
          <span class="stat-label">Total Users</span>
          <span class="stat-value">{{ totalUsers }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon stat-icon-green">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" /><circle cx="9" cy="7" r="4" /><path d="M22 21v-2a4 4 0 0 0-3-3.87" /><path d="M16 3.13a4 4 0 0 1 0 7.75" /></svg>
        </div>
        <div class="stat-info">
          <span class="stat-label">Customers</span>
          <span class="stat-value">{{ customerCount }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon stat-icon-blue">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" /><circle cx="9" cy="7" r="4" /><path d="M22 21v-2a4 4 0 0 0-3-3.87" /><path d="M16 3.13a4 4 0 0 1 0 7.75" /></svg>
        </div>
        <div class="stat-info">
          <span class="stat-label">Admins</span>
          <span class="stat-value">{{ adminCount }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon stat-icon-purple">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" /><circle cx="9" cy="7" r="4" /><path d="M22 21v-2a4 4 0 0 0-3-3.87" /><path d="M16 3.13a4 4 0 0 1 0 7.75" /></svg>
        </div>
        <div class="stat-info">
          <span class="stat-label">Managers</span>
          <span class="stat-value">{{ managerCount }}</span>
        </div>
      </div>
    </div>

    <!-- Control Bar -->
    <div class="control-bar">
      <div class="search-box">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search by name, email, username or ID..."
          class="search-input"
        />
      </div>
      <div class="filter-box">
        <select v-model="roleFilter" class="filter-select" @change="fetchUsers">
          <option v-for="opt in roleOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
      </div>
      <button @click="fetchUsers" class="btn btn-outline" :disabled="loading">
        {{ loading ? 'Loading...' : 'Refresh' }}
      </button>
      <button @click="openCreateModal" class="btn btn-primary">+ Add User</button>
    </div>

    <!-- Data Table Card -->
    <div class="table-card">
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>User</th>
              <th>Username</th>
              <th>Email</th>
              <th>Phone</th>
              <th>Role</th>
              <th>Joined</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in filteredUsers" :key="user.id">
              <td>
                <div class="user-cell">
                  <div class="user-avatar">{{ initials(user) }}</div>
                  <div class="user-meta">
                    <span class="user-fullname">
                      {{ `${user.first_name} ${user.last_name}`.trim() || '-' }}
                    </span>
                    <span class="user-id font-mono">{{ user.id }}</span>
                  </div>
                </div>
              </td>
              <td>{{ user.username || '-' }}</td>
              <td>{{ user.email || '-' }}</td>
              <td>{{ user.phone || '-' }}</td>
              <td>
                <span class="badge" :class="`badge-${user.role}`">
                  {{ roleLabels[user.role] || user.role }}
                </span>
              </td>
              <td>{{ formatDate(user.created_at) }}</td>
              <td>
                <div class="action-buttons">
                  <button
                    @click="openRoleModal(user)"
                    class="btn btn-sm btn-primary-outline"
                    :disabled="loading"
                  >Assign Role</button>
                  <button
                    @click="handleDelete(user)"
                    class="btn btn-sm btn-danger-outline"
                    :disabled="loading"
                  >Delete</button>
                </div>
              </td>
            </tr>
            <tr v-if="loading">
              <td colspan="7" class="no-data">Loading users...</td>
            </tr>
            <tr v-else-if="filteredUsers.length === 0">
              <td colspan="7" class="no-data">No users found.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create User Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal-card">
        <div class="modal-header">
          <h3 class="modal-title">Add New User</h3>
          <button class="modal-close" @click="showCreateModal = false" aria-label="Close">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-grid">
            <div class="form-field">
              <label class="form-label">First Name *</label>
              <input v-model="createForm.first_name" type="text" class="form-input" placeholder="e.g. Budi" />
            </div>
            <div class="form-field">
              <label class="form-label">Last Name *</label>
              <input v-model="createForm.last_name" type="text" class="form-input" placeholder="e.g. Santoso" />
            </div>
            <div class="form-field">
              <label class="form-label">Username *</label>
              <input v-model="createForm.username" type="text" class="form-input" placeholder="e.g. budi.santoso" />
            </div>
            <div class="form-field">
              <label class="form-label">Password *</label>
              <input v-model="createForm.password" type="password" class="form-input" placeholder="Min 8 characters" />
            </div>
            <div class="form-field">
              <label class="form-label">Email *</label>
              <input v-model="createForm.email" type="email" class="form-input" placeholder="e.g. budi@citramasco.com" />
            </div>
            <div class="form-field">
              <label class="form-label">Phone *</label>
              <input v-model="createForm.phone" type="tel" class="form-input" placeholder="e.g. 081234567890" />
            </div>
            <div class="form-field">
              <label class="form-label">Role *</label>
              <select v-model="createForm.role" class="form-input">
                <option v-for="opt in creatableRoles" :key="opt.value" :value="opt.value">
                  {{ opt.label }}
                </option>
              </select>
            </div>
            <div class="form-field form-field-full">
              <label class="form-label">Address</label>
              <input v-model="createForm.address" type="text" class="form-input" placeholder="e.g. Jakarta, Indonesia" />
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" @click="showCreateModal = false" :disabled="creating">Cancel</button>
          <button class="btn btn-primary" @click="submitCreate" :disabled="creating">
            {{ creating ? 'Creating...' : 'Create User' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Assign Role Modal -->
    <div v-if="showRoleModal" class="modal-overlay" @click.self="showRoleModal = false">
      <div class="modal-card modal-card-sm">
        <div class="modal-header">
          <h3 class="modal-title">Assign Role</h3>
          <button class="modal-close" @click="showRoleModal = false" aria-label="Close">&times;</button>
        </div>
        <div class="modal-body">
          <p v-if="roleUser" class="modal-desc">
            Update role for
            <strong>{{ `${roleUser.first_name} ${roleUser.last_name}`.trim() || roleUser.username }}</strong>
          </p>
          <div class="form-field">
            <label class="form-label">Role *</label>
            <select v-model="selectedRole" class="form-input">
              <option v-for="opt in creatableRoles" :key="opt.value" :value="opt.value">
                {{ opt.label }}
              </option>
            </select>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" @click="showRoleModal = false" :disabled="updatingRole">Cancel</button>
          <button class="btn btn-primary" @click="submitRole" :disabled="updatingRole">
            {{ updatingRole ? 'Saving...' : 'Save Role' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.users-view {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* Stats */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 14px;
  background-color: #ffffff;
  padding: 18px;
  border-radius: 16px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon svg {
  width: 22px;
  height: 22px;
}

.stat-icon-primary { background-color: #fff1ea; color: #e15b2b; }
.stat-icon-green { background-color: #d1fae5; color: #059669; }
.stat-icon-blue { background-color: #dbeafe; color: #2563eb; }
.stat-icon-purple { background-color: #ede9fe; color: #7c3aed; }

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-label {
  font-size: 0.75rem;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: #0f172a;
}

/* Control Bar */
.control-bar {
  display: flex;
  gap: 16px;
  background-color: #ffffff;
  padding: 16px;
  border-radius: 16px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  flex-wrap: wrap;
  align-items: center;
}

.search-box {
  flex-grow: 1;
  min-width: 250px;
}

.search-input {
  width: 100%;
  padding: 10px 14px;
  border-radius: 10px;
  border: 1px solid #cbd5e1;
  outline: none;
  font-size: 0.875rem;
  color: #334155;
  box-sizing: border-box;
}

.search-input:focus {
  border-color: #e15b2b;
  box-shadow: 0 0 0 2px rgba(225, 91, 43, 0.1);
}

.filter-select {
  padding: 10px 14px;
  border-radius: 10px;
  border: 1px solid #cbd5e1;
  outline: none;
  font-size: 0.875rem;
  color: #334155;
  background-color: #ffffff;
  cursor: pointer;
}

.filter-select:focus {
  border-color: #e15b2b;
}

/* Table */
.table-card {
  background-color: #ffffff;
  border-radius: 16px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
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
  border-bottom: 1px solid #e2e8f0;
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.data-table td {
  padding: 16px 20px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 0.875rem;
  color: #334155;
  vertical-align: middle;
}

.font-mono { font-family: monospace; }

.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #e15b2b, #f17b50);
  color: #ffffff;
  font-weight: 700;
  font-size: 0.8rem;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.user-meta {
  display: flex;
  flex-direction: column;
}

.user-fullname {
  font-weight: 600;
  color: #0f172a;
}

.user-id {
  font-size: 0.72rem;
  color: #94a3b8;
}

.badge {
  padding: 4px 10px;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  display: inline-block;
}

.badge-admin { background-color: #fff1ea; color: #c2410c; }
.badge-manager { background-color: #dbeafe; color: #1d4ed8; }
.badge-user { background-color: #d1fae5; color: #065f46; }
.badge-reception { background-color: #ede9fe; color: #6d28d9; }
.badge-finance { background-color: #ccfbf1; color: #0f766e; }
.badge-inventory { background-color: #fef3c7; color: #b45309; }

.action-buttons {
  display: flex;
  gap: 8px;
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
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-sm {
  padding: 4px 8px;
  font-size: 0.75rem;
  border-radius: 6px;
}

.btn-outline {
  border-color: #cbd5e1;
  color: #334155;
  background: #ffffff;
}
.btn-outline:hover {
  background-color: #f8fafc;
}

.btn-danger-outline {
  border-color: #ef4444;
  color: #ef4444;
  background: transparent;
}
.btn-danger-outline:hover {
  background-color: #fee2e2;
}

.no-data {
  text-align: center;
  color: #64748b;
  padding: 32px !important;
}

/* Button variants */
.btn-primary {
  background-color: #e15b2b;
  color: #ffffff;
}
.btn-primary:hover { background-color: #c84e20; }

.btn-primary-outline {
  border-color: #e15b2b;
  color: #e15b2b;
  background: transparent;
}
.btn-primary-outline:hover {
  background-color: #fff1ea;
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(15, 23, 42, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  z-index: 50;
}

.modal-card {
  background-color: #ffffff;
  border-radius: 16px;
  width: 100%;
  max-width: 560px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
  overflow: hidden;
}

.modal-card-sm {
  max-width: 420px;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 24px;
  border-bottom: 1px solid #e2e8f0;
}

.modal-title {
  font-size: 1.1rem;
  font-weight: 700;
  color: #0f172a;
  margin: 0;
}

.modal-close {
  border: none;
  background: transparent;
  font-size: 1.5rem;
  line-height: 1;
  color: #94a3b8;
  cursor: pointer;
  padding: 0 4px;
}
.modal-close:hover { color: #334155; }

.modal-body {
  padding: 24px;
}

.modal-desc {
  color: #475569;
  font-size: 0.875rem;
  margin: 0 0 16px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #e2e8f0;
}

/* Form */
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-field-full {
  grid-column: span 2;
}

.form-label {
  font-size: 0.8rem;
  font-weight: 600;
  color: #334155;
}

.form-input {
  padding: 10px 14px;
  border-radius: 10px;
  border: 1px solid #cbd5e1;
  outline: none;
  font-size: 0.875rem;
  color: #334155;
  background-color: #ffffff;
  box-sizing: border-box;
  width: 100%;
}

.form-input:focus {
  border-color: #e15b2b;
  box-shadow: 0 0 0 2px rgba(225, 91, 43, 0.1);
}
</style>
