<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import api from '../utils/api'
import { useAuthStore } from '../stores/auth'
import { Activity, Cpu, HardDrive, LogIn, RefreshCw } from '@lucide/vue'

interface SystemMetrics {
  cpu_usage_ratio: number
  memory_total_bytes: number
  memory_available_bytes: number
  memory_used_bytes: number
  memory_usage_ratio: number
}

const auth = useAuthStore()
const metrics = ref<SystemMetrics | null>(null)
const loading = ref(false)
const error = ref('')
const updatedAt = ref<Date | null>(null)
const isAuthorized = computed(() => auth.isAuthenticated)
let refreshTimer: ReturnType<typeof setInterval> | undefined

async function loadMetrics() {
  if (!isAuthorized.value || loading.value) return
  loading.value = true
  error.value = ''
  try {
    const response = await api.get('/admin/system/metrics')
    metrics.value = response.data.metrics
    updatedAt.value = new Date()
  } catch {
    error.value = '暂时无法获取监控数据，请检查 exporter 是否正常运行。'
  } finally {
    loading.value = false
  }
}

function formatPercent(value: number) {
  return `${(value * 100).toFixed(1)}%`
}

function formatBytes(value: number) {
  return `${(value / 1024 / 1024 / 1024).toFixed(1)} GiB`
}

onMounted(() => {
  if (!isAuthorized.value) return
  loadMetrics()
  refreshTimer = setInterval(loadMetrics, 5000)
})

onUnmounted(() => {
  if (refreshTimer) clearInterval(refreshTimer)
})
</script>

<template>
  <section class="mx-auto max-w-5xl">
    <div class="surface-header surface mb-4 flex items-center justify-between gap-4">
      <div>
        <p class="flex items-center gap-2 text-xs font-medium text-brand-600"><Activity :size="15" /> MINI EXPORTER</p>
        <h1 class="mt-1 page-title">系统监控</h1>
      </div>
      <div v-if="isAuthorized && updatedAt" class="text-right text-xs text-slate-400 dark:text-slate-500">
        <span class="inline-flex items-center gap-1"><RefreshCw :size="13" /> 每 5 秒更新</span><br>{{ updatedAt.toLocaleTimeString() }}
      </div>
    </div>

    <div v-if="!isAuthorized" class="surface p-10 text-center">
      <LogIn :size="32" class="mx-auto mb-3 text-slate-300" />
      <h2 class="text-xl font-bold text-slate-900 dark:text-white">无权限</h2>
      <p class="mt-2 text-sm text-slate-500 dark:text-slate-400">该页面仅限管理员查看，请先登录管理员账号。</p>
      <router-link to="/login" class="btn-primary mt-5">
        管理员登录
      </router-link>
    </div>

    <div v-else>
      <div v-if="metrics" class="grid gap-4 md:grid-cols-2">
        <article class="surface p-6">
          <div class="flex items-center justify-between gap-4">
            <h2 class="flex items-center gap-2 font-semibold text-slate-600 dark:text-slate-300"><Cpu :size="18" class="text-brand-500" /> CPU 使用率</h2>
            <span class="text-3xl font-black text-slate-900 dark:text-white">{{ formatPercent(metrics.cpu_usage_ratio) }}</span>
          </div>
          <div class="mt-6 h-2 overflow-hidden bg-slate-100 dark:bg-white/10">
            <div class="h-full bg-brand-500 transition-all duration-500" :style="{ width: formatPercent(metrics.cpu_usage_ratio) }"></div>
          </div>
        </article>

        <article class="surface p-6">
          <div class="flex items-center justify-between gap-4">
            <h2 class="flex items-center gap-2 font-semibold text-slate-600 dark:text-slate-300"><HardDrive :size="18" class="text-emerald-500" /> 内存使用率</h2>
            <span class="text-3xl font-black text-slate-900 dark:text-white">{{ formatPercent(metrics.memory_usage_ratio) }}</span>
          </div>
          <div class="mt-6 h-2 overflow-hidden bg-slate-100 dark:bg-white/10">
            <div class="h-full bg-emerald-500 transition-all duration-500" :style="{ width: formatPercent(metrics.memory_usage_ratio) }"></div>
          </div>
          <div class="mt-4 flex justify-between gap-4 text-sm text-slate-500 dark:text-slate-400">
            <span>已用 {{ formatBytes(metrics.memory_used_bytes) }}</span>
            <span>总计 {{ formatBytes(metrics.memory_total_bytes) }}</span>
          </div>
        </article>
      </div>

      <div v-else-if="loading" class="border border-gray-200 dark:border-white/10 bg-white/80 dark:bg-slate-900/80 p-8 text-center text-slate-500">
        正在获取监控数据...
      </div>

      <div v-if="error" class="mt-5 border border-red-200 dark:border-red-900/50 bg-red-50 dark:bg-red-950/30 p-4 text-sm text-red-600 dark:text-red-400">
        {{ error }}
      </div>
    </div>
  </section>
</template>
