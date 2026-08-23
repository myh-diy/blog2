<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../utils/api'
import EmptyState from '../components/EmptyState.vue'
import { Hash } from '@lucide/vue'

const tags = ref<{ name: string; count: number }[]>([])
onMounted(async () => { const r = await api.get('/tags'); tags.value = r.data.tags ?? [] })
</script>

<template>
  <div class="mx-auto max-w-5xl">
    <section class="surface overflow-hidden">
      <header class="surface-header"><h1 class="page-title">技术标签</h1><p class="mt-1 text-sm text-slate-400">{{ tags.length }} 个内容方向</p></header>
    <div v-if="tags.length" class="grid grid-cols-2 gap-px bg-slate-100 sm:grid-cols-3 lg:grid-cols-4 dark:bg-white/10">
      <router-link v-for="t in tags" :key="t.name" :to="`/posts?tag=${t.name}`"
        class="group flex min-h-32 flex-col justify-center bg-white p-5 transition-colors hover:bg-brand-50/60 dark:bg-[#171b22] dark:hover:bg-brand-900/10">
        <Hash :size="19" class="mb-3 text-brand-500" />
        <span class="font-semibold text-slate-700 transition-colors group-hover:text-brand-600 dark:text-slate-200">{{ t.name }}</span>
        <span class="mt-1 text-xs text-slate-400">{{ t.count }} 篇文章</span>
      </router-link>
    </div>
    <EmptyState v-else icon="sad" title="暂无标签" description="上传包含标签的文章后会显示在这里。" />
    </section>
  </div>
</template>
