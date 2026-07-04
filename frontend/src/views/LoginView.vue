<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()
const username = ref('')
const password = ref('')
const error = ref('')

async function handleLogin() {
  error.value = ''
  try {
    await auth.login(username.value, password.value)
    router.push('/admin')
  } catch (e) {
    error.value = 'Invalid credentials'
  }
}
</script>

<template>
  <div class="max-w-sm mx-auto mt-20">
    <h1 class="text-2xl font-bold mb-6">Admin Login</h1>
    <form @submit.prevent="handleLogin" class="space-y-4">
      <div>
        <label class="block text-sm mb-1">Username</label>
        <input v-model="username" type="text"
          class="w-full px-3 py-2 border rounded dark:bg-gray-800 dark:border-gray-600" />
      </div>
      <div>
        <label class="block text-sm mb-1">Password</label>
        <input v-model="password" type="password"
          class="w-full px-3 py-2 border rounded dark:bg-gray-800 dark:border-gray-600" />
      </div>
      <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
      <button type="submit"
        class="w-full py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
        Login
      </button>
    </form>
  </div>
</template>
