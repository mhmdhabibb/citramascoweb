<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Sidebar from '@/components/admin/Sidebar.vue'
import { useAuthStore } from '@/stores/authStore'
import { notificationService } from '@/services/admin/notificationService'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const userProfile = computed(() => authStore.user)

const notifications = ref([])
const unreadCount = ref(0)
const showNotifications = ref(false)
let unreadTimer = null

const loadUnreadCount = async () => {
  try {
    unreadCount.value = await notificationService.getUnreadCount()
  } catch {
    /* header bell should not break navigation if the API is down */
  }
}

const toggleNotifications = async () => {
  showNotifications.value = !showNotifications.value
  if (showNotifications.value) {
    try {
      notifications.value = await notificationService.getAll()
      await Promise.all(
        notifications.value
          .filter((n) => !n.is_read)
          .map((n) => notificationService.markRead(n.id)),
      )
      unreadCount.value = 0
    } catch {
      /* ignore */
    }
  }
}

const closeNotifications = () => {
  showNotifications.value = false
}

const timeAgo = (iso) => {
  const diff = Date.now() - new Date(iso).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'just now'
  if (mins < 60) return `${mins}m ago`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  return `${days}d ago`
}

onMounted(() => {
  loadUnreadCount()
  unreadTimer = setInterval(loadUnreadCount, 30000)
})

onUnmounted(() => {
  if (unreadTimer) clearInterval(unreadTimer)
})

// Computes the current page title based on the active route
const pageTitle = computed(() => {
  if (route && route.meta && route.meta.title) {
    return route.meta.title
  }

  const path = route ? route.path : ''
  if (path.includes('dashboard')) return 'Dashboard'
  if (path.includes('reservations')) return 'Reservations'
  if (path.includes('rooms')) return 'Rooms Management'
  if (path.includes('guests')) return 'Guests Management'
  if (path.includes('staff')) return 'Staff Management'
  if (path.includes('promotions')) return 'Promotions'
  if (path.includes('reviews')) return 'Reviews'
  if (path.includes('reports')) return 'Reports'
  if (path.includes('maintenance')) return 'Maintenance'
  if (path.includes('calendar')) return 'Calendar'
  if (path.includes('platform')) return 'Platform Settings'
  if (path.includes('upgrade')) return 'Upgrade Plan'
  if (path.includes('settings')) return 'Settings'

  return 'Admin Portal'
})

const fetchProfile = async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  const profile = await authStore.fetchProfile(true)
  if (!profile) {
    localStorage.removeItem('token')
    router.push('/login')
  }
}

onMounted(() => {
  fetchProfile()
})
</script>

<template>
  <div class="admin-layout">
    <!-- Sidebar component -->
    <Sidebar />

    <!-- Main Content Container -->
    <main class="admin-main">
      <!-- Admin Top Header -->
      <header class="admin-header-modern">
        <div class="header-left">
          <h1 class="page-title">{{ pageTitle }}</h1>
        </div>

        <div class="header-right">
          <div class="notif-wrap">
            <button class="notif-btn" @click="toggleNotifications" aria-label="Notifications">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9" /><path d="M10.3 21a1.94 1.94 0 0 0 3.4 0" /></svg>
              <span v-if="unreadCount > 0" class="notif-badge">{{ unreadCount > 9 ? '9+' : unreadCount }}</span>
            </button>

            <transition name="notif-drop">
              <div v-if="showNotifications" class="notif-dropdown">
                <div class="notif-dropdown-head">
                  <span class="notif-title">Notifications</span>
                  <button class="notif-close" @click="closeNotifications">✕</button>
                </div>
                <div class="notif-list">
                  <div v-for="n in notifications" :key="n.id" class="notif-item" :class="{ unread: !n.is_read }">
                    <div class="notif-icon">💰</div>
                    <div class="notif-body">
                      <span class="notif-item-title">{{ n.title }}</span>
                      <p class="notif-msg">{{ n.message }}</p>
                      <span class="notif-time">{{ timeAgo(n.created_at) }}</span>
                    </div>
                  </div>
                  <div v-if="notifications.length === 0" class="notif-empty">
                    No notifications yet.
                  </div>
                </div>
              </div>
            </transition>
          </div>

          <div class="profile-pill">
            <div class="avatar-wrapper">
              <div class="profile-avatar">
                {{
                  (userProfile?.first_name?.[0] || userProfile?.username?.[0] || 'A').toUpperCase()
                }}
              </div>
              <span class="status-indicator online"></span>
            </div>
            <div class="profile-meta">
              <span class="user-name">
                {{
                  userProfile ? `${userProfile.first_name} ${userProfile.last_name}` : 'Loading...'
                }}
              </span>
              <span class="user-role">
                {{
                  userProfile?.role
                    ? userProfile.role.charAt(0).toUpperCase() + userProfile.role.slice(1)
                    : 'Please wait...'
                }}
              </span>
            </div>
          </div>
        </div>
      </header>

      <!-- Nested route content view -->
      <div class="admin-content-pane">
        <RouterView v-slot="{ Component }">
          <Transition name="page" mode="out-in">
            <component :is="Component" v-if="Component" />
          </Transition>
        </RouterView>
      </div>
    </main>
  </div>
</template>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background-color: #f8fafc; /* slate-50 background for content */
  color: #1e293b; /* slate-800 text */
  font-family: 'Inter', sans-serif;
  overflow: hidden;
}

.admin-main {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

.admin-header-modern {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 32px;
  background-color: #ffffff;
  border-bottom: 1px solid #e2e8f0;
}
.page-title {
  font-size: 1.4rem;
  font-weight: 700;
  color: #0f172a;
}

/* Notification bell */
.header-right {
  display: flex;
  align-items: center;
  gap: 14px;
}
.notif-wrap {
  position: relative;
}
.notif-btn {
  position: relative;
  width: 42px;
  height: 42px;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  background: #ffffff;
  color: #475569;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}
.notif-btn:hover {
  background: #f8fafc;
  color: #0f172a;
}
.notif-btn svg {
  width: 20px;
  height: 20px;
}
.notif-badge {
  position: absolute;
  top: -6px;
  right: -6px;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  border-radius: 9999px;
  background: #ef4444;
  color: #ffffff;
  font-size: 0.7rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid #ffffff;
}
.notif-dropdown {
  position: absolute;
  right: 0;
  top: 52px;
  width: 360px;
  max-width: calc(100vw - 40px);
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(15, 23, 42, 0.12);
  z-index: 50;
  overflow: hidden;
}
.notif-dropdown-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 16px;
  border-bottom: 1px solid #f1f5f9;
}
.notif-title {
  font-size: 0.95rem;
  font-weight: 700;
  color: #0f172a;
}
.notif-close {
  background: #f1f5f9;
  border: none;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  color: #64748b;
  cursor: pointer;
}
.notif-list {
  max-height: 360px;
  overflow-y: auto;
}
.notif-item {
  display: flex;
  gap: 12px;
  padding: 14px 16px;
  border-bottom: 1px solid #f8fafc;
}
.notif-item.unread {
  background-color: #fff7f2;
}
.notif-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: #d1fae5;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  flex-shrink: 0;
}
.notif-body {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.notif-item-title {
  font-size: 0.85rem;
  font-weight: 700;
  color: #1e293b;
}
.notif-msg {
  font-size: 0.8rem;
  color: #64748b;
  margin: 0;
  line-height: 1.4;
}
.notif-time {
  font-size: 0.7rem;
  color: #94a3b8;
}
.notif-empty {
  text-align: center;
  color: #94a3b8;
  padding: 40px;
  font-size: 0.85rem;
}

.notif-drop-enter-active,
.notif-drop-leave-active {
  transition: all 0.2s ease;
}
.notif-drop-enter-from,
.notif-drop-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

.profile-pill {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 6px 14px 6px 6px;
  border-radius: 9999px;
  background-color: #f8fafc;
  border: 1px solid #f1f5f9;
}
.avatar-wrapper {
  position: relative;
}
.profile-avatar {
  width: 38px;
  height: 38px;
  background: orangered;
  color: #ffffff;
  font-weight: 700;
  font-size: 0.95rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  box-shadow: 0 2px 4px rgba(79, 70, 229, 0.15);
}
.status-indicator {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  border: 2px solid #f8fafc;
}
.status-indicator.online {
  background-color: #10b981;
}
.profile-meta {
  display: flex;
  flex-direction: column;
}
.user-name {
  font-size: 0.875rem;
  font-weight: 600;
  color: #1e293b;
  line-height: 1.25;
}
.user-role {
  font-size: 0.75rem;
  font-weight: 500;
  color: #64748b;
}

.admin-content-pane {
  flex-grow: 1;
  padding: 32px;
  overflow-y: auto;
  box-sizing: border-box;
}

/* Page transitions inside admin layout */
.page-enter-active {
  transition:
    opacity 0.3s ease,
    transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.page-leave-active {
  transition:
    opacity 0.2s ease,
    transform 0.2s ease;
}

.page-enter-from {
  opacity: 0;
  transform: translateY(8px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
