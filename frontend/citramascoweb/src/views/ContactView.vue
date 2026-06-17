<script setup lang="ts">
import { ref } from 'vue'

const name = ref('')
const email = ref('')
const message = ref('')
const submitted = ref(false)

const handleSubmit = () => {
  submitted.value = true
}
</script>

<template>
  <div class="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
    <div class="mb-10 text-center md:text-left">
      <h1 class="text-3xl font-extrabold tracking-tight">Contact Us</h1>
      <p class="mt-2 text-muted-foreground">Have any questions? We would love to hear from you.</p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-12">
      <!-- Contact Info -->
      <div class="space-y-8">
        <div>
          <h2 class="text-xl font-bold">Get in Touch</h2>
          <p class="text-muted-foreground mt-2">Feel free to reach out to us using the form or through our direct contact information below.</p>
        </div>

        <div class="space-y-4">
          <div class="flex items-start">
            <div class="flex-shrink-0 text-primary">
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </div>
            <div class="ml-3">
              <p class="font-medium text-foreground">Address</p>
              <p class="text-muted-foreground text-sm">Jl. Raya Citramasco No. 123, Indonesia</p>
            </div>
          </div>

          <div class="flex items-start">
            <div class="flex-shrink-0 text-primary">
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.94.725l.548 2.2a1 1 0 01-.321.988l-1.305.98a10.582 10.582 0 004.872 4.872l.98-1.305a1 1 0 01.988-.321l2.2.548a1 1 0 01.725.94V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
              </svg>
            </div>
            <div class="ml-3">
              <p class="font-medium text-foreground">Phone</p>
              <p class="text-muted-foreground text-sm">+62 21 1234 5678</p>
            </div>
          </div>

          <div class="flex items-start">
            <div class="flex-shrink-0 text-primary">
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
            </div>
            <div class="ml-3">
              <p class="font-medium text-foreground">Email</p>
              <p class="text-muted-foreground text-sm">info@citramascohotel.com</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Contact Form -->
      <div class="bg-card text-card-foreground border border-border rounded-lg p-8 shadow-sm">
        <div v-if="submitted" class="text-center py-8">
          <div class="inline-flex items-center justify-center h-12 w-12 rounded-full bg-green-100 text-green-600 mb-4">
            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h3 class="text-xl font-bold">Message Sent!</h3>
          <p class="text-muted-foreground mt-2">Thank you, {{ name }}. We will get back to you shortly.</p>
          <button @click="submitted = false" class="mt-6 inline-flex items-center justify-center rounded-md bg-primary text-primary-foreground hover:bg-primary/90 px-4 py-2 text-sm font-medium">
            Send Another Message
          </button>
        </div>

        <form v-else @submit.prevent="handleSubmit" class="space-y-6">
          <div>
            <label for="name" class="block text-sm font-medium text-foreground">Name</label>
            <input type="text" id="name" v-model="name" required placeholder="John Doe" class="mt-1 block w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>
          <div>
            <label for="email" class="block text-sm font-medium text-foreground">Email</label>
            <input type="email" id="email" v-model="email" required placeholder="john@example.com" class="mt-1 block w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>
          <div>
            <label for="message" class="block text-sm font-medium text-foreground">Message</label>
            <textarea id="message" v-model="message" rows="4" required placeholder="Your message here..." class="mt-1 block w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"></textarea>
          </div>
          <button type="submit" class="w-full inline-flex items-center justify-center rounded-md bg-primary text-primary-foreground hover:bg-primary/90 py-3 text-base font-semibold shadow-sm">
            Send Message
          </button>
        </form>
      </div>
    </div>
  </div>
</template>
