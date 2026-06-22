<script setup>
import { computed, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Sidebar from '@/components/admin/Sidebar.vue'
import { authService } from '@/services/authService'

const route = useRoute()
const router = useRouter()
const userProfile = ref(null)

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

  // Developer mock session check
  if (token === 'mock-developer-token-citramas') {
    userProfile.value = {
      first_name: 'Developer',
      last_name: 'Mock',
      username: 'admin',
      role: 'admin',
      email: 'admin@citramas.com',
    }
    return
  }

  try {
    const data = await authService.getProfile()
    userProfile.value = data
  } catch (error) {
    console.error('Failed to load user profile:', error)
    // Clear invalid session token and force redirection
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
