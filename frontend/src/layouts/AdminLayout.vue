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
      email: 'admin@citramas.com'
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
      <header class="admin-header">
        <div class="header-left">
          <h1 class="page-title">{{ pageTitle }}</h1>
        </div>
        <div class="header-right">
          <!-- Quick profile block -->
          <div class="user-profile">
            <div class="profile-avatar">
              {{ (userProfile?.first_name?.[0] || userProfile?.username?.[0] || 'A').toUpperCase() }}
            </div>
            <div class="profile-details">
              <span class="profile-name">
                {{ userProfile ? `${userProfile.first_name} ${userProfile.last_name}` : 'Loading...' }}
              </span>
              <span class="profile-role">
                {{ userProfile?.role ? userProfile.role.charAt(0).toUpperCase() + userProfile.role.slice(1) : 'Please wait...' }}
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

.admin-header {
  height: 72px;
  background-color: #ffffff;
  border-bottom: 1px solid rgba(228, 228, 231, 0.8); /* zinc-200/80 */
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32px;
  flex-shrink: 0;
  box-sizing: border-box;
}

.page-title {
  font-size: 1.25rem; /* text-xl */
  font-weight: 700;
  color: #0f172a; /* slate-900 */
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 10px;
}

.profile-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: #e15b2b; /* primary stayhub orange */
  color: #ffffff;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.875rem;
}

.profile-details {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}

.profile-name {
  font-size: 0.875rem;
  font-weight: 600;
  color: #0f172a;
}

.profile-role {
  font-size: 0.75rem;
  color: #64748b; /* slate-500 */
}

.admin-content-pane {
  flex-grow: 1;
  padding: 32px;
  overflow-y: auto;
  box-sizing: border-box;
}

/* Page transitions inside admin layout */
.page-enter-active {
  transition: opacity 0.3s ease, transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.page-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
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
