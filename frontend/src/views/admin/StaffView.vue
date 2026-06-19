<script setup>
import { ref, computed } from 'vue'

// Dummy Data
const staffMembers = ref([
  { id: 'STF-001', name: 'John Miller', role: 'Super Admin', email: 'john.m@citramasco.com', status: 'Active', dateJoined: '2025-01-15' },
  { id: 'STF-002', name: 'Sarah Connor', role: 'Receptionist', email: 'sarah.c@citramasco.com', status: 'Active', dateJoined: '2025-04-10' },
  { id: 'STF-003', name: 'Kyle Reese', role: 'Housekeeping', email: 'kyle.r@citramasco.com', status: 'Inactive', dateJoined: '2025-06-01' },
  { id: 'STF-004', name: 'Marcus Wright', role: 'Maintenance', email: 'marcus.w@citramasco.com', status: 'Active', dateJoined: '2025-09-20' }
])

const searchQuery = ref('')
const roleFilter = ref('All')

// Computed list
const filteredStaff = computed(() => {
  return staffMembers.value.filter(staff => {
    const matchesSearch = staff.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                          staff.email.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                          staff.id.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesRole = roleFilter.value === 'All' || staff.role === roleFilter.value
    return matchesSearch && matchesRole
  })
})

// Actions
const toggleStatus = (id) => {
  const staff = staffMembers.value.find(s => s.id === id)
  if (staff) {
    staff.status = staff.status === 'Active' ? 'Inactive' : 'Active'
  }
}

const addStaff = () => {
  const name = prompt("Enter staff name:")
  if (!name) return
  const role = prompt("Enter role (e.g. Admin, Receptionist, Housekeeping, Maintenance):", "Receptionist")
  const email = prompt("Enter email address:")
  if (!email) return

  const nextIdNum = staffMembers.value.length + 1
  const id = `STF-00${nextIdNum}`

  staffMembers.value.push({
    id,
    name,
    role,
    email,
    status: 'Active',
    dateJoined: new Date().toISOString().split('T')[0]
  })
}

const deleteStaff = (id) => {
  if (confirm(`Are you sure you want to delete staff member ${id}?`)) {
    staffMembers.value = staffMembers.value.filter(s => s.id !== id)
  }
}
</script>

<template>
  <div class="staff-view">
    <!-- Control Bar -->
    <div class="control-bar">
      <div class="search-box">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Search by ID, name or email..." 
          class="search-input"
        />
      </div>
      <div class="filter-box">
        <select v-model="roleFilter" class="filter-select">
          <option value="All">All Roles</option>
          <option value="Super Admin">Super Admin</option>
          <option value="Receptionist">Receptionist</option>
          <option value="Housekeeping">Housekeeping</option>
          <option value="Maintenance">Maintenance</option>
        </select>
      </div>
      <button @click="addStaff" class="btn btn-primary add-btn">
        + Add Staff
      </button>
    </div>

    <!-- Data Table Card -->
    <div class="table-card">
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>Staff ID</th>
              <th>Name</th>
              <th>Role</th>
              <th>Email</th>
              <th>Join Date</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="staff in filteredStaff" :key="staff.id">
              <td class="font-mono font-bold">{{ staff.id }}</td>
              <td>{{ staff.name }}</td>
              <td>{{ staff.role }}</td>
              <td>{{ staff.email }}</td>
              <td>{{ staff.dateJoined }}</td>
              <td>
                <span class="badge" :class="{
                  'badge-success': staff.status === 'Active',
                  'badge-danger': staff.status === 'Inactive'
                }">{{ staff.status }}</span>
              </td>
              <td>
                <div class="action-buttons">
                  <button 
                    @click="toggleStatus(staff.id)" 
                    class="btn btn-sm"
                    :class="staff.status === 'Active' ? 'btn-warning-outline' : 'btn-success-outline'"
                  >
                    {{ staff.status === 'Active' ? 'Deactivate' : 'Activate' }}
                  </button>
                  <button 
                    @click="deleteStaff(staff.id)" 
                    class="btn btn-sm btn-danger-outline"
                  >Delete</button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredStaff.length === 0">
              <td colspan="7" class="no-data">No staff members found.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.staff-view {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

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
.font-bold { font-weight: 600; }

.badge {
  padding: 4px 8px;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  display: inline-block;
}

.badge-success { background-color: #d1fae5; color: #065f46; }
.badge-danger { background-color: #fee2e2; color: #991b1b; }

.action-buttons {
  display: flex;
  gap: 8px;
}

.btn {
  padding: 8px 16px;
  border-radius: 10px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s;
}

.btn-sm {
  padding: 4px 8px;
  font-size: 0.75rem;
  border-radius: 6px;
}

.btn-primary {
  background-color: #e15b2b;
  color: #ffffff;
}
.btn-primary:hover { background-color: #c84e20; }

.btn-success-outline {
  border-color: #10b981;
  color: #10b981;
  background: transparent;
}
.btn-success-outline:hover {
  background-color: #e6fcf5;
}

.btn-warning-outline {
  border-color: #f59e0b;
  color: #f59e0b;
  background: transparent;
}
.btn-warning-outline:hover {
  background-color: #fffbeb;
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
</style>
