import { initializeApp } from 'firebase/app'
import { getMessaging, getToken, onMessage, type Messaging } from 'firebase/messaging'
import { notificationService } from '@/services/admin/notificationService'
import { useToastStore } from '@/stores/toastStore'

// Firebase Web Configuration
export const firebaseConfig = {
  apiKey: import.meta.env.VITE_FIREBASE_API_KEY || '',
  authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN || 'cmliving-1bf26.firebaseapp.com',
  projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID || 'cmliving-1bf26',
  storageBucket: import.meta.env.VITE_FIREBASE_STORAGE_BUCKET || 'cmliving-1bf26.appspot.com',
  messagingSenderId: import.meta.env.VITE_FIREBASE_MESSAGING_SENDER_ID || '',
  appId: import.meta.env.VITE_FIREBASE_APP_ID || '',
}

let messaging: Messaging | null = null

export const initFirebaseMessaging = async () => {
  try {
    if (typeof window === 'undefined' || !('Notification' in window) || !('serviceWorker' in navigator)) {
      return null
    }

    const app = initializeApp(firebaseConfig)
    messaging = getMessaging(app)

    // Listen for foreground push notifications
    onMessage(messaging, (payload) => {
      const toastStore = useToastStore()
      const title = payload.notification?.title || (payload.data && payload.data.title) || 'New Notification'
      const body = payload.notification?.body || (payload.data && payload.data.body) || ''
      toastStore.info(`${title}: ${body}`)
    })

    return messaging
  } catch (error) {
    console.warn('[FCM] Firebase messaging initialization note:', error)
    return null
  }
}

export const requestFCMToken = async (): Promise<string | null> => {
  try {
    if (typeof window === 'undefined' || !('Notification' in window)) return null

    const permission = await Notification.requestPermission()
    if (permission !== 'granted') {
      return null
    }

    const msg = messaging || (await initFirebaseMessaging())
    if (!msg) return null

    const vapidKey = import.meta.env.VITE_FIREBASE_VAPID_KEY || undefined
    const currentToken = await getToken(msg, { vapidKey })
    if (currentToken) {
      // Send token to backend database
      await notificationService.saveDeviceToken(currentToken, 'web')
      return currentToken
    }
    return null
  } catch (error) {
    console.warn('[FCM] Token retrieval note:', error)
    return null
  }
}
