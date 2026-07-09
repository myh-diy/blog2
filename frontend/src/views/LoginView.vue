<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import GradientButton from '../components/GradientButton.vue'

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
        <div class="w-16 h-16 rounded-2xl bg-brand-500 flex items-center justify-center text-white font-black text-2xl shadow-lg mx-auto mb-4">B</div>
        <h1 class="text-2xl font-black text-slate-800 dark:text-slate-100">Admin Login</h1>
        <p class="text-sm text-slate-400 dark:text-slate-500 mt-1">Welcome back, creator!</p>
      </div>

      <form @submit.prevent="handleLogin" class="space-y-4 bg-white dark:bg-slate-900 p-6 rounded-2xl border border-gray-100 dark:border-white/5 shadow-sm">
        <div>
          <label class="block text-sm font-semibold text-slate-600 dark:text-slate-300 mb-1.5">Username</label>
          <input v-model="username" autocomplete="username"
            class="w-full px-4 py-2.5 bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-white/10 rounded-xl text-slate-800 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500/30 focus:border-brand-500 transition-all" />
        </div>
        <div>
          <label class="block text-sm font-semibold text-slate-600 dark:text-slate-300 mb-1.5">Password</label>
          <input v-model="password" type="password" autocomplete="current-password"
            class="w-full px-4 py-2.5 bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-white/10 rounded-xl text-slate-800 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500/30 focus:border-brand-500 transition-all" />
        </div>
        <p v-if="error" class="text-sm text-red-500 font-medium">{{ error }}</p>
        <GradientButton type="submit" class="w-full">Login</GradientButton>
      </form>
    </div>
  </div>
</template>
