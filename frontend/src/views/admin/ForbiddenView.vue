<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'
import { firstAllowedPage } from '@/config/access'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const roleLabel = computed(() => (authStore.role ? authStore.role.toUpperCase() : 'Unknown'))

const goBack = () => {
  const target = firstAllowedPage(authStore.role, 'admin-forbidden')
  if (target === route.name) {
    router.push('/admin/forbidden')
    return
  }
  router.push({ name: target })
}
</script>

<template>
  <div class="forbidden-view">
    <div class="forbidden-card">
      <div class="forbidden-icon">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10" />
          <path d="M15 9l-6 6" />
          <path d="M9 9l6 6" />
        </svg>
      </div>
      <h2>403 — Access Denied</h2>
      <p>
        Your role ({{ roleLabel }}) does not have permission to open this page.
      </p>
      <button class="btn btn-primary" @click="goBack">Go to your dashboard</button>
    </div>
  </div>
</template>

<style scoped>
.forbidden-view {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
}

.forbidden-card {
  text-align: center;
  background-color: #ffffff;
  border: 1px solid rgba(228, 228, 231, 0.8);
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  padding: 40px 48px;
  max-width: 420px;
}

.forbidden-icon {
  width: 56px;
  height: 56px;
  margin: 0 auto 16px;
  border-radius: 50%;
  background-color: #fee2e2;
  color: #dc2626;
  display: flex;
  align-items: center;
  justify-content: center;
}

.forbidden-icon svg {
  width: 28px;
  height: 28px;
}

.forbidden-card h2 {
  margin: 0 0 8px;
  color: #0f172a;
  font-size: 1.25rem;
}

.forbidden-card p {
  color: #64748b;
  margin: 0 0 24px;
}

.btn {
  padding: 10px 18px;
  border-radius: 10px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s;
}

.btn-primary {
  background-color: #e15b2b;
  color: #ffffff;
}
.btn-primary:hover { background-color: #c84e20; }
</style>