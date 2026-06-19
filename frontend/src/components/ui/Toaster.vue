<script setup lang="ts">
import { useToastStore } from '@/stores/toastStore'

const toastStore = useToastStore()
</script>

<template>
  <div class="toaster-container">
    <TransitionGroup name="toast" tag="div" class="toast-list">
      <div
        v-for="toast in toastStore.toasts"
        :key="toast.id"
        class="toast-item"
        :class="`toast-${toast.type}`"
      >
        <!-- Icon -->
        <span class="toast-icon">
          <!-- Success: Check Circle -->
          <svg
            v-if="toast.type === 'success'"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="icon"
          >
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
            <polyline points="22 4 12 14.01 9 11.01" />
          </svg>

          <!-- Error: Alert Circle -->
          <svg
            v-else-if="toast.type === 'error'"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="icon"
          >
            <circle cx="12" cy="12" r="10" />
            <line x1="12" y1="8" x2="12" y2="12" />
            <line x1="12" y1="16" x2="12.01" y2="16" />
          </svg>

          <!-- Warning: Alert Triangle -->
          <svg
            v-else-if="toast.type === 'warning'"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="icon"
          >
            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z" />
            <line x1="12" y1="9" x2="12" y2="13" />
            <line x1="12" y1="17" x2="12.01" y2="17" />
          </svg>

          <!-- Info: Info Circle -->
          <svg
            v-else
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="icon"
          >
            <circle cx="12" cy="12" r="10" />
            <line x1="12" y1="16" x2="12" y2="12" />
            <line x1="12" y1="8" x2="12.01" y2="8" />
          </svg>
        </span>

        <!-- Message -->
        <span class="toast-message">{{ toast.message }}</span>

        <!-- Close Button -->
        <button class="toast-close" @click="toastStore.remove(toast.id)">&times;</button>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toaster-container {
  position: fixed;
  top: 24px;
  right: 24px;
  z-index: 9999;
  pointer-events: none;
  max-width: 380px;
  width: 100%;
}

.toast-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.toast-item {
  pointer-events: auto;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 12px;
  background-color: #ffffff;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.05), 0 4px 6px -2px rgba(0, 0, 0, 0.02);
  color: #1e293b;
  font-family: 'Inter', sans-serif;
  font-size: 0.875rem;
  font-weight: 500;
  position: relative;
  overflow: hidden;
}

/* Glassmorphism styling based on types */
.toast-success {
  border-left: 4px solid #10b981;
  background-color: rgba(240, 253, 250, 0.95);
  color: #0f5132;
}
.toast-success .toast-icon {
  color: #10b981;
}

.toast-error {
  border-left: 4px solid #ef4444;
  background-color: rgba(254, 242, 242, 0.95);
  color: #842029;
}
.toast-error .toast-icon {
  color: #ef4444;
}

.toast-warning {
  border-left: 4px solid #f59e0b;
  background-color: rgba(255, 251, 235, 0.95);
  color: #664d03;
}
.toast-warning .toast-icon {
  color: #f59e0b;
}

.toast-info {
  border-left: 4px solid #3b82f6;
  background-color: rgba(239, 246, 255, 0.95);
  color: #084298;
}
.toast-info .toast-icon {
  color: #3b82f6;
}

.toast-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.icon {
  width: 20px;
  height: 20px;
}

.toast-message {
  flex-grow: 1;
  line-height: 1.4;
}

.toast-close {
  background: none;
  border: none;
  font-size: 1.25rem;
  color: #94a3b8;
  cursor: pointer;
  padding: 0 4px;
  line-height: 1;
  flex-shrink: 0;
  transition: color 0.15s;
}

.toast-close:hover {
  color: #1e293b;
}

/* Slide & Fade Transitions */
.toast-enter-from {
  opacity: 0;
  transform: translateX(100px) scale(0.9);
}
.toast-enter-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
.toast-leave-active {
  transition: all 0.2s ease;
  position: absolute;
}
.toast-leave-to {
  opacity: 0;
  transform: translateX(50px) scale(0.95);
}
.toast-move {
  transition: transform 0.3s ease;
}
</style>
