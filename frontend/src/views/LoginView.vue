<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from '@/services/authService'
import { useAuthStore } from '@/stores/authStore'

const router = useRouter()

const username = ref('')
const password = ref('')
const isPasswordVisible = ref(false)
const isLoading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const rememberMe = ref(false)
const loginVideoSrc = ref('/CM_LP2.mp4')

const videoA = ref<HTMLVideoElement | null>(null)
const videoB = ref<HTMLVideoElement | null>(null)
const activePlayer = ref<'A' | 'B'>('A')

const CROSSFADE_DURATION = 0.8 // seconds before end to crossfade
let isCrossfading = false

const handleTimeUpdate = (player: 'A' | 'B') => {
  const currentVideo = player === 'A' ? videoA.value : videoB.value
  const nextVideo = player === 'A' ? videoB.value : videoA.value
  if (!currentVideo || !nextVideo || !currentVideo.duration) return

  const remaining = currentVideo.duration - currentVideo.currentTime
  if (remaining <= CROSSFADE_DURATION && !isCrossfading && activePlayer.value === player) {
    isCrossfading = true
    nextVideo.currentTime = 0
    nextVideo.muted = true
    nextVideo.play().catch(() => {})
    activePlayer.value = player === 'A' ? 'B' : 'A'

    setTimeout(() => {
      isCrossfading = false
      currentVideo.pause()
      currentVideo.currentTime = 0
    }, CROSSFADE_DURATION * 1000)
  }
}

onMounted(async () => {
  if (videoA.value) {
    videoA.value.muted = true
    videoA.value.defaultMuted = true
    try {
      await videoA.value.play()
    } catch (e) {
      console.warn('Login Video A autoplay fallback:', e)
    }
  }
  if (videoB.value) {
    videoB.value.muted = true
    videoB.value.defaultMuted = true
  }
  if (localStorage.getItem('rememberMe') === 'true') {
    rememberMe.value = true
    username.value = localStorage.getItem('savedUsername') || ''
    password.value = localStorage.getItem('savedPassword') || ''
  }
})

const handleLogin = async () => {
  try {
    isLoading.value = true
    errorMessage.value = ''
    successMessage.value = ''

    const data = await authService.login(username.value, password.value)
    
    // Store token
    localStorage.setItem('token', data.access_token)
    
    if (rememberMe.value) {
      localStorage.setItem('savedUsername', username.value)
      localStorage.setItem('savedPassword', password.value)
      localStorage.setItem('rememberMe', 'true')
    } else {
      localStorage.removeItem('savedUsername')
      localStorage.removeItem('savedPassword')
      localStorage.removeItem('rememberMe')
    }

    successMessage.value = 'Sign in successful! Redirecting...'
    
    // reset cached profile so the guard picks up the new user's role
    useAuthStore().reset()

    // Brief delay to let the user see the success message
    setTimeout(() => {
      router.push('/admin')
    }, 1000)
  } catch (error: any) {
    console.error('Login error:', error)
    
    // Developer Fallback: If network connection is refused (backend not running)
    const isNetworkError = error.message === 'Network Error' || 
                           error.code === '500' || 
                           error.message?.includes('Network Error') ||
                           error.message?.includes('NetworkError')
    console.log(isNetworkError)

    if (isNetworkError) {
      errorMessage.value = 'Cannot connect to the server. Please check your internet connection.'
      return
    }
                           
    
    
    errorMessage.value =  'Invalid credentials. Please try again.'
  } finally {
    isLoading.value = false
  }
}



</script>

<template>
  <div class="login-container">
    <!-- Dual-Player Seamless Crossfade Video Engine -->
    <video
      ref="videoA"
      class="login-bg-video"
      :class="{ 'video-active': activePlayer === 'A' }"
      :src="loginVideoSrc"
      muted
      playsinline
      preload="auto"
      disablePictureInPicture
      @timeupdate="handleTimeUpdate('A')"
    ></video>
    <video
      ref="videoB"
      class="login-bg-video"
      :class="{ 'video-active': activePlayer === 'B' }"
      :src="loginVideoSrc"
      muted
      playsinline
      preload="auto"
      disablePictureInPicture
      @timeupdate="handleTimeUpdate('B')"
    ></video>
    
    <!-- Transparent Glass Overlay -->
    <div class="login-video-overlay"></div>

    <!-- Centered Glassmorphic Login Container -->
    <div class="login-centered-layout">
      <div class="login-card">
        <!-- Logo Emblem -->
        <div class="brand-logo-box">
          <img src="/logo-light.png" alt="CM Living" class="login-logo-img" />
        </div>

        <!-- Header -->
        <div class="login-header">
          <p class="login-portal-subtitle">Staff & Management Portal</p>
        </div>

          <!-- Feedback Messages -->
          <Transition name="fade">
            <div v-if="errorMessage" class="feedback-alert error-alert">
              <svg class="alert-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="8" x2="12" y2="12"></line>
                <line x1="12" y1="16" x2="12.01" y2="16"></line>
              </svg>
              <span>{{ errorMessage }}</span>
            </div>
          </Transition>

          <Transition name="fade">
            <div v-if="successMessage" class="feedback-alert success-alert">
              <svg class="alert-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                <polyline points="22 4 12 14.01 9 11.01"></polyline>
              </svg>
              <span>{{ successMessage }}</span>
            </div>
          </Transition>

          <!-- Form -->
          <form @submit.prevent="handleLogin" class="login-form">
            <!-- Username / Email Input -->
            <div class="input-group">
              <span class="input-icon">
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path>
                  <polyline points="22,6 12,13 2,6"></polyline>
                </svg>
              </span>
              <input 
                v-model="username" 
                type="text" 
                placeholder="Username" 
                required
                :disabled="isLoading"
                class="form-input"
                autocomplete="username"
              />
            </div>

            <!-- Password Input -->
            <div class="input-group password-group">
              <span class="input-icon">
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
                </svg>
              </span>
              <input 
                v-model="password" 
                :type="isPasswordVisible ? 'text' : 'password'" 
                placeholder="Password" 
                required
                :disabled="isLoading"
                class="form-input"
                autocomplete="current-password"
              />
              <button 
                type="button" 
                class="toggle-password" 
                @click="isPasswordVisible = !isPasswordVisible"
                tabindex="-1"
              >
                <!-- Eye Off (Hidden) -->
                <svg v-if="!isPasswordVisible" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path>
                  <line x1="1" y1="1" x2="23" y2="23"></line>
                </svg>
                <!-- Eye (Visible) -->
                <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                  <circle cx="12" cy="12" r="3"></circle>
                </svg>
              </button>
            </div>

            <!-- Remember Me -->
            <div class="remember-me-container">
              <label class="remember-me-label">
                <input type="checkbox" v-model="rememberMe" class="remember-me-checkbox" />
                Ingat Saya
              </label>
            </div>

            <!-- Submit Button -->
            <button type="submit" class="submit-btn" :disabled="isLoading">
              <span v-if="isLoading" class="btn-spinner"></span>
              <span>{{ isLoading ? 'Signing in...' : 'Sign In' }}</span>
            </button>
          </form>
        </div>
      </div>
    </div>
  </template>

<style scoped>
/* Google Fonts Import for high-fidelity typography */
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=Outfit:wght@500;600;700&display=swap');

.login-container {
  height: 100vh;
  width: 100vw;
  display: flex;
  font-family: 'Inter', sans-serif;
  background-color: #0f172a;
  overflow: hidden;
  box-sizing: border-box;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.login-container::-webkit-scrollbar {
  display: none;
  width: 0;
  height: 0;
}

.login-bg-video {
  position: absolute;
  top: 50%;
  left: 50%;
  min-width: 100%;
  min-height: 100%;
  width: auto;
  height: auto;
  transform: translate3d(-50%, -50%, 0);
  object-fit: cover;
  z-index: 1;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.8s ease-in-out;
}

.login-bg-video.video-active {
  opacity: 1;
}

.login-video-overlay {
  position: absolute;
  inset: 0;
  background: radial-gradient(
    circle at center,
    rgba(15, 23, 42, 0.45) 0%,
    rgba(15, 23, 42, 0.75) 100%
  );
  z-index: 2;
}

/* Centered Glassmorphic Container */
.login-centered-layout {
  position: relative;
  z-index: 10;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  box-sizing: border-box;
  overflow: hidden;
}

/* Glassmorphism login card design */
.login-card {
  width: 440px;
  max-width: 100%;
  background: rgba(28, 22, 18, 0.68);
  backdrop-filter: blur(20px) saturate(190%);
  -webkit-backdrop-filter: blur(20px) saturate(190%);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 28px;
  padding: 44px 38px 40px;
  box-shadow: 
    0 24px 60px rgba(0, 0, 0, 0.45),
    0 0 0 1px rgba(255, 255, 255, 0.08) inset;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  animation: cardFadeIn 0.8s cubic-bezier(0.16, 1, 0.3, 1);
}

.brand-logo-box {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 20px;
}

.login-logo-img {
  height: 96px;
  width: auto;
  object-fit: contain;
  filter: drop-shadow(0 6px 16px rgba(0, 0, 0, 0.45));
  margin-bottom: 6px;
  transition: transform 0.3s ease;
}

.login-logo-img:hover {
  transform: scale(1.03);
}

/* Subtitle styles */
.login-portal-subtitle {
  font-size: 0.95rem;
  font-weight: 500;
  letter-spacing: 0.04em;
  color: rgba(255, 255, 255, 0.85);
  margin: 0 0 24px 0;
  text-align: center;
}

.login-form {
  display: flex;
  flex-direction: column;
  width: 100%;
}

/* Feedback message styling */
.feedback-alert {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: 14px;
  font-size: 0.85rem;
  margin-bottom: 20px;
  text-align: left;
  line-height: 1.4;
  backdrop-filter: blur(10px);
}

.error-alert {
  background: rgba(239, 68, 68, 0.2);
  border: 1px solid rgba(239, 68, 68, 0.4);
  color: #fca5a5;
}

.success-alert {
  background: rgba(34, 197, 94, 0.2);
  border: 1px solid rgba(34, 197, 94, 0.4);
  color: #86efac;
}

.alert-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

/* Input group layouts */
.input-group {
  position: relative;
  margin-bottom: 18px;
  width: 100%;
  box-sizing: border-box;
}

.input-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: #ffffff;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
  transition: color 0.2s ease;
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.3));
}

.form-input {
  width: 100%;
  padding: 15px 16px 15px 48px;
  box-sizing: border-box;
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 14px;
  font-size: 0.95rem;
  color: #ffffff;
  outline: none;
  font-family: inherit;
  transition: all 0.25s ease;
}

.form-input::placeholder {
  color: rgba(255, 255, 255, 0.6);
}

/* Input focus dynamics */
.form-input:focus {
  background: rgba(255, 255, 255, 0.2);
  border-color: #fb923c;
  box-shadow: 0 0 0 3px rgba(251, 146, 60, 0.35);
}

/* WebKit Browser Autofill High Contrast Fix */
.form-input:-webkit-autofill,
.form-input:-webkit-autofill:hover,
.form-input:-webkit-autofill:focus,
.form-input:-webkit-autofill:active {
  -webkit-box-shadow: 0 0 0 1000px rgba(32, 24, 18, 0.96) inset !important;
  -webkit-text-fill-color: #ffffff !important;
  caret-color: #ffffff !important;
  transition: background-color 5000s ease-in-out 0s;
  border-color: rgba(255, 255, 255, 0.4) !important;
}

.input-group:focus-within .input-icon {
  color: #fb923c;
}

/* Password toggler */
.password-group .form-input {
  padding-right: 48px;
}

.toggle-password {
  position: absolute;
  right: 16px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #ffffff;
  z-index: 10;
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.2s ease;
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.3));
}

.toggle-password:hover {
  color: #fb923c;
}

.toggle-password:hover {
  color: #ffffff;
}

/* Forgot password link container */
.forgot-password-container {
  text-align: right;
  margin-bottom: 24px;
  margin-top: -8px;
}

.forgot-link {
  font-size: 0.85rem;
  color: #475569;
  font-weight: 500;
  text-decoration: none;
  transition: color 0.2s ease;
}

.forgot-link:hover {
  color: #0f172a;
}

/* Remember me container */
.remember-me-container {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  margin-bottom: 24px;
  margin-top: 8px;
}

.remember-me-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.85rem;
  color: rgba(255, 255, 255, 0.85);
  cursor: pointer;
  user-select: none;
}

.remember-me-checkbox {
  appearance: none;
  width: 16px;
  height: 16px;
  border: 1.5px solid rgba(255, 255, 255, 0.4);
  border-radius: 4px;
  outline: none;
  cursor: pointer;
  position: relative;
  transition: all 0.2s ease;
  background: rgba(255, 255, 255, 0.08);
}

.remember-me-checkbox:checked {
  background: #fb923c;
  border-color: #fb923c;
}

.remember-me-checkbox:checked::after {
  content: '';
  position: absolute;
  top: 1px;
  left: 4px;
  width: 4px;
  height: 8px;
  border: solid white;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}

.remember-me-label:hover .remember-me-checkbox:not(:checked) {
  border-color: #fb923c;
}

/* Main action button */
.submit-btn {
  width: 100%;
  padding: 15px 24px;
  box-sizing: border-box;
  background: #f8f5f1;
  border: none;
  border-radius: 14px;
  color: #1c1612;
  font-size: 0.95rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
}

.submit-btn:hover:not(:disabled) {
  background: #ffffff;
  transform: translateY(-2px);
  box-shadow: 0 8px 28px rgba(255, 255, 255, 0.25);
}

.submit-btn:active:not(:disabled) {
  transform: translateY(0);
}

.submit-btn:disabled {
  opacity: 0.65;
  cursor: not-allowed;
}

/* Loading spinner */
.btn-spinner {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #ffffff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

/* Custom dashed divider */
.social-divider {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 28px 0;
  position: relative;
}

.social-divider::before {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  top: 50%;
  border-top: 1px dotted rgba(148, 163, 184, 0.55);
  z-index: 1;
}

.divider-text {
  font-size: 0.75rem;
  font-weight: 500;
  color: #94a3b8;
  background: #e9eff5; /* Soft background to overlay dotted line cleanly */
  padding: 4px 12px;
  border-radius: 20px;
  z-index: 2;
  text-transform: capitalize;
}

/* Social logins buttons grid */
.social-login-grid {
  display: flex;
  gap: 12px;
  width: 100%;
}

.social-btn {
  flex: 1;
  height: 52px;
  border-radius: 14px;
  border: 1px solid rgba(255, 255, 255, 0.6);
  background: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.01);
}

.social-btn:hover {
  background: #ffffff;
  border-color: rgba(15, 23, 42, 0.15);
  box-shadow: 
    0 10px 20px -8px rgba(15, 23, 42, 0.05),
    0 4px 12px rgba(15, 23, 42, 0.02);
  transform: translateY(-1px);
}

.social-btn:active {
  transform: translateY(0);
}

/* Animations */
@keyframes cardFadeIn {
  from {
    opacity: 0;
    transform: translateY(24px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Alert fade transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
