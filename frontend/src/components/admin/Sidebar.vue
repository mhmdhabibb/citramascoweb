<script setup>
import { authService } from '@/services/authService'
import { useToastStore } from '@/stores/toastStore'
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

// State
const isCollapsed = ref(false)
const searchQuery = ref('')
const activeItem = ref('Dashboard')
const toastStore = useToastStore()

const route = useRoute()
const router = useRouter()

// Menu Definition
const menuGroups = ref([
  {
    groupName: 'Master',
    items: [
      {
        name: 'Dashboard',
        icon: 'LayoutDashboard',
        route: '/admin/dashboard'
      },
      {
        name: 'Reservation',
        icon: 'CalendarCheck',
        route: '/admin/reservations'
      },
      {
        name: 'Manage Rooms',
        icon: 'Bed',
        isOpen: false,
        children: [
          { name: 'Room List', route: '/admin/rooms' },
          { name: 'Room Types', route: '/admin/room-types' },
          { name: 'Room Categories', route: '/admin/room-categories' },
        ]
      },
    
      // {
      //   name: 'Manage Staff',
      //   icon: 'UserCog',
      //   isOpen: false,
      //   children: [
      //     { name: 'Staff List', route: '/admin/staff' },
      //     { name: 'Roles & Permissions', route: '/admin/staff/roles' }
      //   ]
      // },
      {
        name: 'Promotions',
        icon: 'Percent',
        route: '/admin/promotions'
      },
     
    ]
  },
  
  
])

// Bottom Items
const bottomItems = [
  {
    name: 'Help & Center',
    icon: 'HelpCircle',
    route: '/admin/help'
  }
]

// Inline SVG Icon Paths
const icons = {
  LayoutDashboard: '<rect x="3" y="3" width="7" height="9" rx="1"/><rect x="14" y="3" width="7" height="5" rx="1"/><rect x="14" y="12" width="7" height="9" rx="1"/><rect x="3" y="16" width="7" height="5" rx="1"/>',
  CalendarCheck: '<rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/><path d="m9 16 2 2 4-4"/>',
  Bed: '<path d="M2 4v16"/><path d="M2 8h18a2 2 0 0 1 2 2v10"/><path d="M2 14h20"/><path d="M6 8v6"/>',
  Users: '<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/>',
  UserCog: '<circle cx="9" cy="7" r="4"/><path d="M10 15H6a4 4 0 0 0-4 4v2"/><circle cx="19" cy="16" r="2"/><path d="M19 13v1"/><path d="M19 18v1"/><path d="M18 15h-1"/><path d="M21 17h-1"/><path d="M18 17l-.7-.7"/><path d="M20.7 15.3l-.7-.7"/><path d="M18 15.3l-.7.7"/><path d="M20.7 17l-.7.7"/>',
  Percent: '<line x1="19" y1="5" x2="5" y2="19"/><circle cx="6.5" cy="6.5" r="2.5"/><circle cx="17.5" cy="17.5" r="2.5"/>',
  MessageSquare: '<path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>',
  BarChart3: '<path d="M3 3v18h18"/><path d="M18 17V9"/><path d="M13 17V5"/><path d="M8 17v-3"/>',
  Wrench: '<path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>',
  Calendar: '<rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/>',
  Monitor: '<rect x="2" y="3" width="20" height="14" rx="2" ry="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/>',
  Sparkles: '<path d="m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z"/><path d="m5 3 1 2.5L8.5 6 6 7 5 9.5 4 7 1.5 6 4 5.5Z"/><path d="M19 17l1 2.5 2.5.5-2.5 1-1 2.5-1-2.5-2.5-1 2.5-1Z"/>',
  Settings: '<circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>',
  HelpCircle: '<circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/><line x1="12" y1="17" x2="12.01" y2="17"/>',
  LogOut: '<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/>',
  ChevronLeft: '<polyline points="15 18 9 12 15 6"/>',
  ChevronRight: '<polyline points="9 18 15 12 9 6"/>',
  ChevronDown: '<polyline points="6 9 12 15 18 9"/>',
  Search: '<circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>'
}

// Toggle Sidebar Collapse
const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value
}

// Toggle Dropdown (only works if sidebar is expanded)
const toggleDropdown = (item) => {
  if (isCollapsed.value) {
    isCollapsed.value = false
    item.isOpen = true
    return
  }
  item.isOpen = !item.isOpen
}

// Check route and activate matching item
watch(
  () => route?.path,
  (newPath) => {
    if (!newPath) return
    for (const group of menuGroups.value) {
      for (const item of group.items) {
        if (item.route && newPath.startsWith(item.route)) {
          activeItem.value = item.name
          return
        }
        if (item.children) {
          for (const subItem of item.children) {
            if (subItem.route && newPath === subItem.route) {
              activeItem.value = item.name
              item.isOpen = true
              return
            }
          }
        }
      }
    }
  },
  { immediate: true }
)

// Search filtering logic
const filteredMenuGroups = computed(() => {
  if (!searchQuery.value) return menuGroups.value

  const query = searchQuery.value.toLowerCase()
  return menuGroups.value
    .map((group) => {
      const filteredItems = group.items.filter((item) => {
        // Match main item name
        const matchName = item.name.toLowerCase().includes(query)
        // Match sub-items name
        const matchChildren = item.children?.some((sub) =>
          sub.name.toLowerCase().includes(query)
        )
        return matchName || matchChildren
      })

      return {
        ...group,
        items: filteredItems.map((item) => {
          // If query matches sub-items, automatically open this dropdown
          const hasMatchingChildren = item.children?.some((sub) =>
            sub.name.toLowerCase().includes(query)
          )
          return {
            ...item,
            isOpen: hasMatchingChildren ? true : item.isOpen
          }
        })
      }
    })
    .filter((group) => group.items.length > 0)
})

// Handle Item Click (both local fallback & router navigation)
const handleItemClick = (item) => {
  if (item.children) {
    toggleDropdown(item)
  } else {
    activeItem.value = item.name
    if (item.route && router) {
      router.push(item.route)
    }
  }
}

// Handle Sub-item Click
const handleSubItemClick = (parentName, subItem) => {
  activeItem.value = parentName
  if (subItem.route && router) {
    router.push(subItem.route)
  }
}

// Handle Logout
const handleLogout = async () => {
  await authService.logout()
  localStorage.removeItem('token')
  toastStore.success('Logout success')
  router.push('/login')
}
</script>

<template>
  <aside class="sidebar" :class="{ 'collapsed': isCollapsed }">
    <!-- Header / Brand Logo -->
    <div class="header" :class="{ 'header-collapsed': isCollapsed }">
      <div class="brand">
        <!-- Orange Icon Block -->
        <div class="logo-box">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="3"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="logo-icon"
          >
            <line x1="18" y1="19" x2="18" y2="10" />
            <line x1="12" y1="19" x2="12" y2="4" />
            <line x1="6" y1="19" x2="6" y2="14" />
          </svg>
        </div>
        <!-- Brand Name -->
        <span v-if="!isCollapsed" class="brand-name">CitraMas Co </span>
      </div>

      <!-- Collapse Toggle Arrow -->
      <button
        @click="toggleCollapse"
        class="toggle-btn"
        :class="{ 'toggle-btn-collapsed': isCollapsed }"
        aria-label="Toggle Sidebar"
      >
        <svg
          v-if="!isCollapsed"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="arrow-icon"
          v-html="icons.ChevronLeft"
        ></svg>
        <svg
          v-else
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="arrow-icon"
          v-html="icons.ChevronRight"
        ></svg>
      </button>
    </div>

   
    <!-- Scrollable Menu Items -->
    <div class="menu-container sidebar-scroll">
      <div v-for="group in filteredMenuGroups" :key="group.groupName" class="menu-group">
        <!-- Group Header -->
        <h3 v-if="!isCollapsed" class="group-title">{{ group.groupName }}</h3>
        <!-- Group Divider (when collapsed) -->
        <div v-else class="group-divider"></div>

        <!-- Menu List -->
        <div class="menu-list">
          <div v-for="item in group.items" :key="item.name" class="menu-item-wrapper">
            <!-- Active Indicator Pill on the very left -->
            <div v-if="activeItem === item.name" class="active-indicator"></div>

            <!-- Main Button Link -->
            <button
              @click="handleItemClick(item)"
              class="menu-item"
              :class="{ 'active': activeItem === item.name }"
            >
              <!-- Icon -->
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="item-icon"
                :class="{ 'icon-active': activeItem === item.name }"
                v-html="icons[item.icon]"
              ></svg>

              <!-- Label -->
              <span v-if="!isCollapsed" class="item-label">{{ item.name }}</span>

              <!-- Dropdown Chevron -->
              <svg
                v-if="item.children && !isCollapsed"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="chevron-icon"
                :class="{ 'chevron-open': item.isOpen }"
                v-html="icons.ChevronDown"
              ></svg>
            </button>

            <!-- Dropdown Sub-menu -->
            <div
              v-if="item.children && !isCollapsed"
              class="submenu"
              :style="{ maxHeight: item.isOpen ? `${item.children.length * 36 + 8}px` : '0px' }"
            >
              <div class="submenu-inner">
                <button
                  v-for="sub in item.children"
                  :key="sub.name"
                  @click="handleSubItemClick(item.name, sub)"
                  class="submenu-item"
                >
                  {{ sub.name }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Section -->
    <div class="bottom-section">
      <!-- Help and Center -->
      <div v-for="item in bottomItems" :key="item.name" class="menu-item-wrapper">
        <div v-if="activeItem === item.name" class="active-indicator"></div>

        <button
          @click="handleItemClick(item)"
          class="menu-item"
          :class="{ 'active': activeItem === item.name }"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="item-icon"
            :class="{ 'icon-active': activeItem === item.name }"
            v-html="icons[item.icon]"
          ></svg>
          <span v-if="!isCollapsed" class="item-label">{{ item.name }}</span>
        </button>
      </div>

      <!-- Logout -->
      <button @click="handleLogout" class="menu-item logout-btn">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="item-icon logout-icon"
          v-html="icons.LogOut"
        ></svg>
        <span v-if="!isCollapsed" class="item-label">Logout</span>
      </button>
    </div>
  </aside>
</template>

<style scoped>
/* Color Palette Variables */
.sidebar {
  --primary-orange: #E15B2B;
  --primary-orange-gradient-end: #F17B50;
  --bg-orange-light: rgba(225, 91, 43, 0.05);
  --border-orange-light: rgba(225, 91, 43, 0.1);
  --bg-sidebar: #ffffff;
  --text-main: #27272a; /* zinc-800 */
  --text-muted: #71717a; /* zinc-500 */
  --text-light: #a1a1aa; /* zinc-400 */
  --bg-hover: #f4f4f5; /* zinc-100 */
  --border-color: rgba(228, 228, 231, 0.8); /* zinc-200/80 */
  --border-light: #f4f4f5; /* zinc-100 */
  --shadow-sidebar: 4px 0 24px -10px rgba(0, 0, 0, 0.03);
  --shadow-logo: 0 4px 6px -1px rgba(225, 91, 43, 0.1), 0 2px 4px -2px rgba(225, 91, 43, 0.1);
  --font-family: 'Inter', sans-serif;
  --transition-speed: 0.3s;
}

/* Sidebar Container */
.sidebar {
  position: relative;
  height: 100vh;
  width: 256px; /* w-64 */
  background-color: var(--bg-sidebar);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  transition: width var(--transition-speed) cubic-bezier(0.4, 0, 0.2, 1);
  user-select: none;
  box-shadow: var(--shadow-sidebar);
  font-family: var(--font-family);
  box-sizing: border-box;
}

.sidebar.collapsed {
  width: 80px; /* w-20 */
}

/* Header Section */
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 16px;
  border-bottom: 1px solid var(--border-light);
  box-sizing: border-box;
}

.header.header-collapsed {
  justify-content: center;
}

/* Brand Logo */
.brand {
  display: flex;
  align-items: center;
  gap: 12px;
  overflow: hidden;
  white-space: nowrap;
}

.logo-box {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, var(--primary-orange), var(--primary-orange-gradient-end));
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-logo);
  flex-shrink: 0;
  transition: transform 0.2s ease;
}

.logo-box:hover {
  transform: scale(1.05);
}

.logo-icon {
  width: 20px;
  height: 20px;
  stroke: #ffffff;
}

.brand-name {
  font-weight: 800;
  font-size: 1.125rem;
  color: var(--text-main);
  letter-spacing: -0.025em;
  animation: fadeIn 0.3s ease;
}

/* Collapse Toggle Button */
.toggle-btn {
  width: 28px;
  height: 28px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-light);
  background-color: var(--bg-sidebar);
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  padding: 0;
}

.toggle-btn:hover {
  color: var(--text-main);
  background-color: var(--bg-hover);
}

.toggle-btn.toggle-btn-collapsed {
  position: absolute;
  right: -14px;
  top: 24px;
  z-index: 10;
}

.arrow-icon {
  width: 16px;
  height: 16px;
}

/* Search Section */
.search-section {
  padding: 12px 16px;
  box-sizing: border-box;
}

.search-collapsed {
  display: flex;
  justify-content: center;
}

.search-box {
  position: relative;
  display: flex;
  align-items: center;
  background-color: rgba(244, 244, 245, 0.6);
  border: 1px solid rgba(228, 228, 231, 0.6);
  border-radius: 12px;
  padding: 8px 12px;
  transition: all 0.2s ease;
}

.search-box:hover {
  background-color: #f4f4f5;
}

.search-box:focus-within {
  background-color: #ffffff;
  border-color: rgba(225, 91, 43, 0.5);
  box-shadow: 0 0 0 2px rgba(225, 91, 43, 0.15);
}

.search-icon {
  width: 16px;
  height: 16px;
  color: var(--text-light);
  flex-shrink: 0;
  margin-right: 8px;
}

.search-input {
  width: 100%;
  background: transparent;
  border: none;
  outline: none;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-main);
}

.search-input::placeholder {
  color: var(--text-light);
}

.shortcut-badge {
  display: flex;
  align-items: center;
  gap: 2px;
  font-family: monospace;
  font-size: 10px;
  color: var(--text-light);
  background-color: #ffffff;
  border: 1px solid var(--border-color);
  padding: 2px 6px;
  border-radius: 4px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  flex-shrink: 0;
  margin-left: 6px;
}

/* Collapsed Search Icon-only Button */
.search-btn-collapsed {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  border: 1px dashed var(--border-color);
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-light);
  cursor: pointer;
  transition: all 0.2s ease;
  padding: 0;
}

.search-btn-collapsed:hover {
  border-color: var(--text-muted);
  color: var(--text-muted);
}

/* Menu Container */
.menu-container {
  flex: 1;
  overflow-y: auto;
  padding: 8px 12px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  box-sizing: border-box;
}

.menu-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.group-title {
  padding: 0 12px;
  font-size: 10px;
  font-weight: 700;
  color: rgba(113, 113, 122, 0.8);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 6px;
  margin-top: 0;
}

.group-divider {
  height: 1px;
  background-color: var(--border-light);
  margin: 8px 8px;
}

.menu-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

/* Active indicator pill */
.menu-item-wrapper {
  position: relative;
}

.active-indicator {
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 28px;
  background-color: var(--primary-orange);
  border-radius: 0 4px 4px 0;
  z-index: 10;
  transition: all 0.2s ease;
}

/* Main Button Link */
.menu-item {
  width: 100%;
  display: flex;
  align-items: center;
  padding: 10px 14px;
  border-radius: 12px;
  border: 1px solid transparent;
  background-color: transparent;
  color: var(--text-muted);
  font-size: 0.875rem;
  font-weight: 500;
  text-align: left;
  cursor: pointer;
  position: relative;
  transition: all 0.2s ease;
  box-sizing: border-box;
}

.menu-item:hover {
  color: var(--text-main);
  background-color: var(--bg-hover);
}

.menu-item.active {
  color: var(--primary-orange);
  background-color: var(--bg-orange-light);
  border-color: var(--border-orange-light);
  font-weight: 600;
}

.item-icon {
  width: 20px;
  height: 20px;
  color: var(--text-light);
  flex-shrink: 0;
  transition: transform 0.2s ease, color 0.2s ease;
}

.menu-item:hover .item-icon {
  color: var(--text-muted);
  transform: scale(1.03);
}

.item-icon.icon-active {
  color: var(--primary-orange);
}

.item-label {
  margin-left: 12px;
  flex-grow: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  animation: fadeIn 0.2s ease;
}

.chevron-icon {
  width: 14px;
  height: 14px;
  color: rgba(113, 113, 122, 0.8);
  transition: transform 0.2s ease;
  flex-shrink: 0;
}

.chevron-open {
  transform: rotate(180deg);
  color: var(--primary-orange);
}

/* Submenu */
.submenu {
  overflow: hidden;
  transition: max-height 0.3s ease-in-out;
}

.submenu-inner {
  padding-left: 44px;
  padding-right: 8px;
  padding-top: 4px;
  padding-bottom: 4px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  border-left: 1px solid var(--border-light);
  margin-left: 24px;
  margin-top: 2px;
}

.submenu-item {
  display: block;
  width: 100%;
  padding: 6px 12px;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--text-muted);
  background-color: transparent;
  border: none;
  border-radius: 8px;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s ease;
  box-sizing: border-box;
}

.submenu-item:hover {
  color: var(--text-main);
  background-color: rgba(244, 244, 245, 0.8);
}

/* Bottom Section */
.bottom-section {
  padding: 12px;
  border-top: 1px solid var(--border-light);
  background-color: var(--bg-sidebar);
  display: flex;
  flex-direction: column;
  gap: 2px;
  box-sizing: border-box;
}

.logout-btn {
  color: #e11d48; /* rose-600 */
}

.logout-btn:hover {
  color: #be123c; /* rose-700 */
  background-color: rgba(254, 242, 242, 0.5); /* rose-50/50 */
}

.logout-icon {
  color: rgba(244, 63, 94, 0.8); /* rose-500/80 */
  transition: transform 0.2s ease;
}

.logout-btn:hover .logout-icon {
  color: #e11d48;
  transform: translateX(2px);
}

/* Scrollbar styling */
.sidebar-scroll::-webkit-scrollbar {
  width: 4px;
}

.sidebar-scroll::-webkit-scrollbar-track {
  background: transparent;
}

.sidebar-scroll::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.05);
  border-radius: 99px;
}

.sidebar-scroll:hover::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.12);
}

/* Animations */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateX(-4px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}
</style>
