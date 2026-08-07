import { createRouter, createWebHistory } from 'vue-router'
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import AdminLayout from '@/layouts/AdminLayout.vue'
import HomeView from '@/views/HomeView.vue'
import RoomsView from '@/views/RoomsView.vue'
import RoomDetailView from '@/views/RoomDetailView.vue'
import BookingView from '@/views/BookingView.vue'
import ContactView from '@/views/ContactView.vue'
import LoginView from '@/views/LoginView.vue'
import OffersView from '@/views/OffersView.vue'
import { useAuthStore } from '@/stores/authStore'
import { canAccess, firstAllowedPage } from '@/config/access'

// Admin Views
import DashboardView from '@/views/admin/DashboardView.vue'
import ReservationsView from '@/views/admin/ReservationsView.vue'
import AdminRoomsView from '@/views/admin/RoomsView.vue'
import RoomTypesView from '@/views/admin/RoomTypesView.vue'
import RoomCategoriesView from '@/views/admin/RoomCategoriesView.vue'
import StaffView from '@/views/admin/StaffView.vue'
import PromotionsView from '@/views/admin/PromotionsView.vue'
import UsersView from '@/views/admin/UsersView.vue'
import InventoryView from '@/views/admin/InventoryView.vue'
import ForbiddenView from '@/views/admin/ForbiddenView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior(to: { hash: any }, from: any, savedPosition: any) {
    if (savedPosition) {
      return savedPosition
    }
    if (to.hash) {
      return { el: to.hash, behavior: 'smooth' }
    }
    return { top: 0 }
  },
  routes: [    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { title: 'Login' },
    },
    {
      path: '/',
      component: DefaultLayout,
      children: [
        {
          path: '',
          name: 'home',
          component: HomeView,
        },
        {
          path: 'rooms',
          name: 'rooms',
          component: RoomsView,
        },
        {
          path: 'rooms/:id',
          name: 'room-detail',
          component: RoomDetailView,
        },
        {
          path: 'booking',
          name: 'booking',
          component: BookingView,
        },
        {
          path: 'contact',
          name: 'contact',
          component: ContactView,
        },
        {
          path: 'offers',
          name: 'offers',
          component: OffersView,
        },
      ],
    },
    {
      path: '/admin',
      component: AdminLayout,
      children: [
        {
          path: '',
          redirect: '/admin/dashboard',
        },
        {
          path: 'dashboard',
          name: 'admin-dashboard',
          component: DashboardView,
          meta: { title: 'Dashboard' },
        },
        {
          path: 'reservations',
          name: 'admin-reservations',
          component: ReservationsView,
          meta: { title: 'Reservations' },
        },
{
          path: 'users',
          name: 'admin-users',
          component: UsersView,
          meta: { title: 'Manage Users' }
        },
        {
          path: 'inventory',
          name: 'admin-inventory',
          component: InventoryView,
          meta: { title: 'Inventory' }
        },
        {
          path: 'rooms',
          name: 'admin-rooms',
          component: AdminRoomsView,
          meta: { title: 'Manage Rooms' },
        },
        {
          path: 'room-types',
          name: 'admin-room-types',
          component: RoomTypesView,
          meta: { title: 'Room Types' },
        },
        {
          path: 'room-categories',
          name: 'admin-room-categories',
          component: RoomCategoriesView,
          meta: { title: 'Room Categories' },
        },
        {
          path: 'staff',
          name: 'admin-staff',
          component: StaffView,
          meta: { title: 'Manage Staff' },
        },
        {
          path: 'promotions',
          name: 'admin-promotions',
          component: PromotionsView,
          meta: { title: 'Promotions' },
        },
        {
          path: 'forbidden',
          name: 'admin-forbidden',
          component: ForbiddenView,
          meta: { title: 'Access Denied' },
        },
      ],
    },
  ],
})

/**
 * Global role-based guard for the admin area.
 *  - Blocks access when not logged in (redirect to login).
 *  - Redirects to the role's first allowed page when a chosen page isn't permitted.
 */
router.beforeEach(async (to) => {
  if (!to.path.startsWith('/admin')) return true

  const token = localStorage.getItem('token')
  if (!token) return { name: 'login' }

  const authStore = useAuthStore()
  const profile = await authStore.fetchProfile()
  if (!profile) return { name: 'login' }

  // always allow the forbidden fallback page itself
  if (to.path === '/admin/forbidden') return true

  // default landing for /admin and empty paths
  if (to.path === '/admin' || to.path === '/admin/') {
    return { name: firstAllowedPage(profile.role) }
  }

  if (canAccess(profile.role, to.path)) return true

  // not permitted -> send to that role's first accessible page
  return { name: firstAllowedPage(profile.role) }
})

export default router
