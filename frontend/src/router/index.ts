import { createRouter, createWebHistory } from 'vue-router'
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import AdminLayout from '@/layouts/AdminLayout.vue'
import HomeView from '@/views/HomeView.vue'
import RoomsView from '@/views/RoomsView.vue'
import RoomDetailView from '@/views/RoomDetailView.vue'
import BookingView from '@/views/BookingView.vue'
import ContactView from '@/views/ContactView.vue'
import LoginView from '@/views/LoginView.vue'

// Admin Views
import DashboardView from '@/views/admin/DashboardView.vue'
import ReservationsView from '@/views/admin/ReservationsView.vue'
import AdminRoomsView from '@/views/admin/RoomsView.vue'
import RoomTypesView from '@/views/admin/RoomTypesView.vue'
import RoomCategoriesView from '@/views/admin/RoomCategoriesView.vue'
import StaffView from '@/views/admin/StaffView.vue'
import PromotionsView from '@/views/admin/PromotionsView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { title: 'Login' }
    },
    {
      path: '/',
      component: DefaultLayout,
      children: [
        {
          path: '',
          name: 'home',
          component: HomeView
        },
        {
          path: 'rooms',
          name: 'rooms',
          component: RoomsView
        },
        {
          path: 'rooms/:id',
          name: 'room-detail',
          component: RoomDetailView
        },
        {
          path: 'booking',
          name: 'booking',
          component: BookingView
        },
        {
          path: 'contact',
          name: 'contact',
          component: ContactView
        }
      ]
    },
    {
      path: '/admin',
      component: AdminLayout,
      children: [
        {
          path: '',
          redirect: '/admin/dashboard'
        },
        {
          path: 'dashboard',
          name: 'admin-dashboard',
          component: DashboardView,
          meta: { title: 'Dashboard' }
        },
        {
          path: 'reservations',
          name: 'admin-reservations',
          component: ReservationsView,
          meta: { title: 'Reservations' }
        },
        {
          path: 'rooms',
          name: 'admin-rooms',
          component: AdminRoomsView,
          meta: { title: 'Manage Rooms' }
        },
        {
          path: 'room-types',
          name: 'admin-room-types',
          component: RoomTypesView,
          meta: { title: 'Room Types' }
        },
        {
          path: 'room-categories',
          name: 'admin-room-categories',
          component: RoomCategoriesView,
          meta: { title: 'Room Categories' }
        },
        {
          path: 'staff',
          name: 'admin-staff',
          component: StaffView,
          meta: { title: 'Manage Staff' }
        },
        {
          path: 'promotions',
          name: 'admin-promotions',
          component: PromotionsView,
          meta: { title: 'Promotions' }
        }
      ]
    }
  ]
})

export default router

