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
  try { await auth.login(username.value, password.value); router.push('/admin') }
  catch { error.value = 'Invalid credentials' }
}
</script>

<template>
  <div class="min-h-[60vh] flex items-center justify-center">
    <div class="w-full max-w-sm">
      <div class="text-center mb-8">
        <div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-brand-400 to-brand-600 dark:from-pop-400 dark:to-pop-600 flex items-center justify-center text-white font-black text-xl shadow-lg mx-auto mb-4">B</div>
        <h1 class="text-2xl font-bold text-warm-800 dark:text-warm-100">Admin Login</h1>
      </div>
      <form @submit.prevent="handleLogin" class="space-y-4 bg-white dark:bg-white/5 p-6 rounded-2xl border border-warm-200 dark:border-white/5">
        <div>
          <label class="block text-sm font-medium text-warm-600 dark:text-warm-300 mb-1.5">Username</label>
          <input v-model="username" autocomplete="username"
            class="w-full px-4 py-2.5 bg-warm-50 dark:bg-white/5 border border-warm-200 dark:border-white/10 rounded-xl text-warm-800 dark:text-warm-100 placeholder-warm-400 focus:outline-none focus:ring-2 focus:ring-brand-500/30 dark:focus:ring-pop-500/30 focus:border-brand-500 dark:focus:border-pop-500 transition-all" />
        </div>
        <div>
          <label class="block text-sm font-medium text-warm-600 dark:text-warm-300 mb-1.5">Password</label>
          <input v-model="password" type="password" autocomplete="current-password"
            class="w-full px-4 py-2.5 bg-warm-50 dark:bg-white/5 border border-warm-200 dark:border-white/10 rounded-xl text-warm-800 dark:text-warm-100 placeholder-warm-400 focus:outline-none focus:ring-2 focus:ring-brand-500/30 dark:focus:ring-pop-500/30 focus:border-brand-500 dark:focus:border-pop-500 transition-all" />
        </div>
        <p v-if="error" class="text-sm text-red-500 font-medium">{{ error }}</p>
        <button type="submit" class="btn-primary w-full py-2.5">Login</button>
      </form>
    </div>
  </div>
</template>
