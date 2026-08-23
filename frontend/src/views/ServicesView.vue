<script setup lang="ts">
import { ExternalLink, LockKeyhole, Server } from '@lucide/vue'

interface PrivateService {
  name: string
  description: string
  port: number
}

const services: PrivateService[] = [
  {
    name: 'SillyTavern',
    description: '私人 AI 会话与角色管理服务',
    port: 8000,
  },
]
const currentHost = window.location.hostname

function serviceUrl(port: number) {
  return `${window.location.protocol}//${currentHost}:${port}`
}
</script>

<template>
  <div class="mx-auto max-w-4xl">
    <section class="surface overflow-hidden">
      <header class="surface-header flex items-start gap-3">
        <span class="flex h-10 w-10 shrink-0 items-center justify-center bg-brand-50 text-brand-600 dark:bg-brand-900/20 dark:text-brand-300">
          <LockKeyhole :size="20" />
        </span>
        <div>
          <h1 class="page-title">私人服务</h1>
          <p class="mt-1 text-sm text-slate-400">仅管理员可访问的个人服务入口</p>
        </div>
      </header>

      <div class="divide-y divide-slate-100 dark:divide-white/10">
        <a
          v-for="service in services"
          :key="service.port"
          :href="serviceUrl(service.port)"
          target="_blank"
          rel="noopener noreferrer"
          class="group flex items-center gap-4 px-5 py-5 transition-colors hover:bg-slate-50 dark:hover:bg-white/[0.03]"
        >
          <span class="flex h-11 w-11 shrink-0 items-center justify-center bg-slate-100 text-slate-500 group-hover:bg-brand-50 group-hover:text-brand-600 dark:bg-white/5 dark:text-slate-400 dark:group-hover:bg-brand-900/20 dark:group-hover:text-brand-300">
            <Server :size="21" />
          </span>
          <span class="min-w-0 flex-1">
            <span class="block font-semibold text-slate-900 group-hover:text-brand-600 dark:text-white dark:group-hover:text-brand-300">{{ service.name }}</span>
            <span class="mt-1 block text-sm text-slate-400">{{ service.description }}</span>
            <code class="mt-2 block truncate text-xs text-slate-400">{{ currentHost }}:{{ service.port }}</code>
          </span>
          <span class="inline-flex shrink-0 items-center gap-1.5 text-sm font-medium text-brand-600 dark:text-brand-300">
            打开
            <ExternalLink :size="16" />
          </span>
        </a>
      </div>
    </section>
  </div>
</template>
