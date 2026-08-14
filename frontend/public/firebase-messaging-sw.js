// Firebase Messaging Service Worker for Background Web Push Notifications
importScripts('https://www.gstatic.com/firebasejs/10.12.0/firebase-app-compat.js')
importScripts('https://www.gstatic.com/firebasejs/10.12.0/firebase-messaging-compat.js')

const firebaseConfig = {
  apiKey: 'AIzaSyDemoKeyCitramasHotelPlaceholder',
  authDomain: 'citramas-hotel.firebaseapp.com',
  projectId: 'citramas-hotel',
  storageBucket: 'citramas-hotel.appspot.com',
  messagingSenderId: '123456789012',
  appId: '1:123456789012:web:abcdef1234567890',
}

firebase.initializeApp(firebaseConfig)

const messaging = firebase.messaging()

messaging.onBackgroundMessage((payload) => {
  const notificationTitle = payload.notification ? payload.notification.title : 'Notifikasi Citramas Hotel'
  const notificationOptions = {
    body: payload.notification ? payload.notification.body : '',
    icon: '/logo.png',
    data: payload.data,
  }

  self.registration.showNotification(notificationTitle, notificationOptions)
})
