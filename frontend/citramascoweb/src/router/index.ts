import { createRouter, createWebHistory } from 'vue-router'
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import HomeView from '@/views/HomeView.vue'
import RoomsView from '@/views/RoomsView.vue'
import RoomDetailView from '@/views/RoomDetailView.vue'
import BookingView from '@/views/BookingView.vue'
import ContactView from '@/views/ContactView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
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
        }
      ]
    }
  ]
})

export default router

