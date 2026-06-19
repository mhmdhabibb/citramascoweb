<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from '@/services/authService'

const router = useRouter()

const username = ref('')
const password = ref('')
const isPasswordVisible = ref(false)
const isLoading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const handleLogin = async () => {
  try {
    isLoading.value = true
    errorMessage.value = ''
    successMessage.value = ''

    const data = await authService.login(username.value, password.value)
    
    // Store token
    localStorage.setItem('token', data.access_token)
    
    successMessage.value = 'Sign in successful! Redirecting...'
    
    // Brief delay to let the user see the success message
    setTimeout(() => {
      router.push('/admin')
    }, 1000)
  } catch (error: any) {
    console.error('Login error:', error)
    
    // Developer Fallback: If network connection is refused (backend not running)
    const isNetworkError = error.message === 'Network Error' || 
                           error.code === 'ERR_NETWORK' || 
                           error.message?.includes('Network Error') ||
                           error.message?.includes('NetworkError')
                           
    if (isNetworkError) {
      if ((username.value === 'admin' || username.value === 'admin@citramas.com') && password.value === 'admin123') {
        localStorage.setItem('token', 'mock-developer-token-citramas')
        successMessage.value = 'Offline/Mock Sign in successful! Redirecting...'
        setTimeout(() => {
          router.push('/admin')
        }, 1000)
        return
      } else {
        errorMessage.value = 'Database/Server offline. Try offline mock login: admin / admin123'
        return
      }
    }
    
    errorMessage.value = error.response?.data?.message || error.message || 'Invalid credentials. Please try again.'
  } finally {
    isLoading.value = false
  }
}



</script>

<template>
  <div class="login-container">
    <div class="login-card-wrapper">
      <div class="login-card">
        <!-- Squircle Brand Icon Box -->
        <div class="brand-icon-box">
          <svg 
            class="brand-icon" 
            viewBox="0 0 24 24" 
            width="24" 
            height="24" 
            stroke="currentColor" 
            stroke-width="2.5" 
            fill="none" 
            stroke-linecap="round" 
            stroke-linejoin="round"
          >
            <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4" />
            <polyline points="10 17 15 12 10 7" />
            <line x1="15" y1="12" x2="3" y2="12" />
          </svg>
        </div>

        <!-- Header -->
        <div class="login-header">
          <h2>Sign in </h2>
          <p>Access your Citra Mas account to manage bookings, rooms, and hotel operations.</p>
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
              placeholder="Email or Username" 
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

         

          <!-- Submit Button -->
          <button type="submit" class="submit-btn" :disabled="isLoading">
            <span v-if="isLoading" class="btn-spinner"></span>
            <span>{{ isLoading ? 'Signing in...' : 'Get Started' }}</span>
          </button>
        </form>

        <!-- Divider -->
       
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Google Fonts Import for high-fidelity typography */
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=Outfit:wght@500;600;700&display=swap');

.login-container {
  min-height: 100vh;
  width: 100vw;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'Inter', sans-serif;
  
  /* Deep sky/clouds background configuration matching reference image layout */
  background: radial-gradient(circle at 50% 50%, rgba(255, 255, 255, 0.4) 0%, rgba(255, 255, 255, 0) 100%),
              url('https://images.unsplash.com/photo-1513002749550-c59d786b8e6c?q=80&w=1920');
  background-size: cover;
  background-position: center bottom;
  background-repeat: no-repeat;
  overflow: hidden;
  box-sizing: border-box;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
}

.login-card-wrapper {
  animation: cardFadeIn 0.8s cubic-bezier(0.16, 1, 0.3, 1);
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  padding: 20px;
  box-sizing: border-box;
}

/* Glassmorphism login card design */
.login-card {
  width: 440px;
  max-width: 100%;
  background: rgba(255, 255, 255, 0.55);
  backdrop-filter: blur(28px) saturate(180%);
  -webkit-backdrop-filter: blur(28px) saturate(180%);
  border-radius: 28px;
  border: 1px solid rgba(255, 255, 255, 0.45);
  box-shadow: 
    0 24px 64px -12px rgba(15, 23, 42, 0.08),
    0 0 1px 0 rgba(255, 255, 255, 0.6) inset;
  padding: 48px 36px;
  box-sizing: border-box;
  text-align: center;
  display: flex;
  flex-direction: column;
}

/* Brand icon box - Squircle shape with shadow */
.brand-icon-box {
  width: 64px;
  height: 64px;
  background: #ffffff;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 24px auto;
  box-shadow: 
    0 10px 25px -5px rgba(0, 0, 0, 0.04),
    0 0 1px 0 rgba(0, 0, 0, 0.1) inset;
  color: #1e293b;
  border: 1px solid rgba(255, 255, 255, 0.8);
}

.brand-icon {
  width: 24px;
  height: 24px;
  stroke: #1e293b;
}

/* Header styles */
.login-header h2 {
  font-family: 'Outfit', sans-serif;
  font-size: 1.55rem;
  font-weight: 700;
  color: #0f172a;
  margin: 0 0 8px 0;
  letter-spacing: -0.02em;
}

.login-header p {
  font-size: 0.875rem;
  line-height: 1.5;
  color: #64748b;
  margin: 0 0 28px 0;
  padding: 0 8px;
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
  border-radius: 12px;
  font-size: 0.85rem;
  margin-bottom: 20px;
  text-align: left;
  line-height: 1.4;
  box-shadow: 0 4px 12px -2px rgba(0, 0, 0, 0.03);
}

.error-alert {
  background: rgba(254, 242, 242, 0.8);
  border: 1px solid rgba(239, 68, 68, 0.2);
  color: #b91c1c;
}

.success-alert {
  background: rgba(240, 253, 244, 0.8);
  border: 1px solid rgba(34, 197, 94, 0.2);
  color: #15803d;
}

.alert-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

/* Input group layouts */
.input-group {
  position: relative;
  margin-bottom: 16px;
  width: 100%;
  box-sizing: border-box;
}

.input-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: #94a3b8;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
  transition: color 0.2s ease;
}

.form-input {
  width: 100%;
  padding: 15px 16px 15px 48px;
  box-sizing: border-box;
  background: rgba(241, 243, 249, 0.7);
  border: 1px solid transparent;
  border-radius: 14px;
  font-size: 0.95rem;
  color: #0f172a;
  outline: none;
  font-family: inherit;
  transition: all 0.2s ease;
}

.form-input::placeholder {
  color: #94a3b8;
}

/* Input focus dynamics */
.form-input:focus {
  background: #ffffff;
  box-shadow: 
    0 10px 20px -10px rgba(15, 23, 42, 0.04),
    0 0 0 2px rgba(15, 23, 42, 0.08);
}

.form-input:focus + .input-icon,
.input-group:focus-within .input-icon {
  color: #0f172a;
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
  color: #94a3b8;
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.2s ease;
}

.toggle-password:hover {
  color: #0f172a;
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

/* Main action button */
.submit-btn {
  width: 100%;
  padding: 15px 24px;
  box-sizing: border-box;
  background: #18181b;
  border: none;
  border-radius: 14px;
  color: #ffffff;
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  box-shadow: 0 4px 12px rgba(24, 24, 27, 0.12);
}

.submit-btn:hover:not(:disabled) {
  background: #000000;
  transform: translateY(-1px);
  box-shadow: 
    0 12px 24px -10px rgba(24, 24, 27, 0.25),
    0 4px 12px rgba(24, 24, 27, 0.15);
}

.submit-btn:active:not(:disabled) {
  transform: translateY(0);
}

.submit-btn:disabled {
  opacity: 0.75;
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
