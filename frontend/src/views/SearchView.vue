<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import api from '../utils/api'
import PostCard from '../components/PostCard.vue'
import EmptyState from '../components/EmptyState.vue'
import GradientButton from '../components/GradientButton.vue'
import { Search, Sparkles } from '@lucide/vue'
import type { Post } from '../stores/posts'

const route = useRoute()
const results = ref<Post[]>([])
const total = ref(0)
const loading = ref(false)
const q = ref((route.query.q as string) || '')

watch(() => route.query.q, (val) => { q.value = (val as string) || '' })

async function search() {
  if (!q.value) return
  loading.value = true
  try {
    const r = await api.get('/search', { params: { q: q.value } })
    results.value = r.data.posts ?? []
    total.value = r.data.total
  } finally { loading.value = false }
}
</script>

<template>
  <div class="mx-auto max-w-4xl">
    <section class="surface overflow-hidden">
      <header class="surface-header">
        <h1 class="page-title">搜索文章</h1>
        <p class="mt-1 text-sm text-slate-400">搜索标题与 Markdown 正文</p>
      </header>
    <form @submit.prevent="search" class="flex gap-3 border-b border-slate-100 p-5 dark:border-white/10">
      <div class="relative flex-1">
        <Search class="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" :size="18" />
        <input v-model="q" placeholder="输入关键词，例如 Go、Vue、Prometheus" class="text-input h-11 pl-10" />
      </div>
      <GradientButton type="submit">搜索</GradientButton>
    </form>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="w-8 h-8 border-2 border-brand-500 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <template v-else-if="results.length">
      <p class="border-b border-slate-100 px-5 py-3 text-sm text-slate-400 dark:border-white/10">找到 {{ total }} 篇结果</p>
      <div>
        <PostCard v-for="post in results" :key="post.id" :post="post" compact />
      </div>
    </template>

    <EmptyState v-else-if="q" icon="search" title="没有找到结果" :description="`没有匹配“${q}”的文章，请换个关键词。`" />
    <div v-else class="flex min-h-64 flex-col items-center justify-center text-slate-400"><Sparkles :size="32" class="mb-3" /><p class="text-sm">输入关键词开始探索</p></div>
    </section>
  </div>
</template>
