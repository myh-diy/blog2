<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import GradientButton from '../components/GradientButton.vue'
import { LockKeyhole, UserRound } from '@lucide/vue'
import { useSiteAvatar } from '../composables/useSiteAvatar'

const router = useRouter()
const auth = useAuthStore()
const username = ref('')
const password = ref('')
const error = ref('')
const { siteAvatar } = useSiteAvatar()

async function handleLogin() {
  error.value = ''
  try { await auth.login(username.value, password.value); router.push('/admin') }
  catch { error.value = 'Invalid credentials' }
}
</script>

<template>
  <div class="flex min-h-[65vh] items-center justify-center">
    <div class="w-full max-w-sm">
      <div class="mb-6 text-center">
        <img :src="siteAvatar" alt="站点头像" class="mx-auto mb-3 h-14 w-14 bg-brand-50 object-contain p-1 dark:bg-brand-900/20" />
        <h1 class="text-xl font-bold text-slate-900 dark:text-white">创作者登录</h1>
        <p class="mt-1 text-sm text-slate-400">管理文章、主题与系统状态</p>
      </div>

      <form @submit.prevent="handleLogin" class="surface space-y-4 p-6">
        <div>
          <label class="mb-1.5 block text-sm font-medium text-slate-600 dark:text-slate-300">用户名</label>
          <div class="relative"><UserRound class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" :size="17" /><input v-model="username" autocomplete="username" class="text-input pl-9" /></div>
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-slate-600 dark:text-slate-300">密码</label>
          <div class="relative"><LockKeyhole class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" :size="17" /><input v-model="password" type="password" autocomplete="current-password" class="text-input pl-9" /></div>
        </div>
        <p v-if="error" class="text-sm text-red-500 font-medium">{{ error }}</p>
        <GradientButton type="submit" class="w-full">登录</GradientButton>
      </form>
    </div>
  </div>
</template>
