import { createRouter, createWebHistory } from 'vue-router'
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import HomeView from '@/views/HomeView.vue'
import RoomsView from '@/views/RoomsView.vue'
import RoomDetailView from '@/views/RoomDetailView.vue'
import BookingView from '@/views/BookingView.vue'
import ContactView from '@/views/ContactView.vue'
import OffersView from '@/views/OffersView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    if (to.hash) {
      return { el: to.hash, behavior: 'smooth' }
    }
    return { top: 0 }
  },
  routes: [
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
        },
        {
          path: 'offers',
          name: 'offers',
          component: OffersView
        }
      ]
    }
  ]
})

export default router

