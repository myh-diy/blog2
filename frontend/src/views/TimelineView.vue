<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Archive, CalendarDays, ChevronRight } from '@lucide/vue'
import api from '../utils/api'

interface TimelinePost { id: number; title: string; slug: string; created_at: string }
interface TimelineEntry { year: number; month: number; posts: TimelinePost[] }

const timeline = ref<TimelineEntry[]>([])
onMounted(async () => { const r = await api.get('/timeline'); timeline.value = r.data.timeline ?? [] })
</script>

<template>
  <div class="mx-auto max-w-4xl">
    <section class="surface overflow-hidden">
      <header class="surface-header flex items-center gap-3"><span class="flex h-10 w-10 items-center justify-center bg-brand-50 text-brand-600 dark:bg-brand-900/20"><Archive :size="20" /></span><div><h1 class="page-title">文章归档</h1><p class="mt-1 text-xs text-slate-400">沿时间回看每一次记录</p></div></header>
    <div v-if="timeline.length" class="divide-y divide-slate-100 dark:divide-white/10">
      <section v-for="entry in timeline" :key="`${entry.year}-${entry.month}`" class="grid gap-4 px-5 py-6 sm:grid-cols-[120px_1fr]">
        <div><p class="flex items-center gap-2 text-sm font-semibold text-slate-900 dark:text-white"><CalendarDays :size="16" class="text-brand-500" />{{ entry.year }} 年</p><p class="mt-2 pl-6 text-sm text-slate-400">{{ entry.month }} 月</p></div>
        <div class="divide-y divide-slate-100 dark:divide-white/10">
          <router-link v-for="post in entry.posts" :key="post.id" :to="`/post/${post.slug}`" class="group flex items-center justify-between gap-4 py-3 first:pt-0 last:pb-0">
            <div class="min-w-0"><h2 class="truncate text-sm font-medium text-slate-700 group-hover:text-brand-600 dark:text-slate-200">{{ post.title }}</h2><time class="mt-1 block text-xs text-slate-400">{{ new Date(post.created_at).toLocaleDateString('zh-CN') }}</time></div><ChevronRight :size="16" class="shrink-0 text-slate-300 group-hover:text-brand-500" />
          </router-link>
        </div>
      </section>
    </div>
    <p v-else class="py-20 text-center text-sm text-slate-400">还没有文章。</p>
    </section>
  </div>
</template>
