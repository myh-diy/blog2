<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import api from '../utils/api'
import PostCard from '../components/PostCard.vue'
import EmptyState from '../components/EmptyState.vue'
import GradientButton from '../components/GradientButton.vue'
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
    results.value = r.data.posts
    total.value = r.data.total
  } finally { loading.value = false }
}
</script>

<template>
  <div class="max-w-3xl mx-auto">
    <div class="text-center mb-10">
      <h1 class="text-3xl font-black text-slate-800 dark:text-slate-100 mb-2">Search</h1>
      <p class="text-slate-400 dark:text-slate-500 text-sm">Find posts by title, content, or tags</p>
    </div>

    <form @submit.prevent="search" class="flex gap-3 mb-10">
      <div class="relative flex-1">
        <svg class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
        <input v-model="q" placeholder="Search posts..."
          class="w-full pl-11 pr-4 py-3.5 bg-white dark:bg-slate-900 border border-gray-200 dark:border-white/10 rounded-2xl text-slate-800 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500/30 focus:border-brand-500 transition-all" />
      </div>
      <GradientButton type="submit">Search</GradientButton>
    </form>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="w-8 h-8 border-2 border-brand-500 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <template v-else-if="results.length">
      <p class="text-sm text-slate-400 dark:text-slate-500 mb-6">{{ total }} result(s)</p>
      <div class="grid md:grid-cols-2 gap-6">
        <PostCard v-for="post in results" :key="post.id" :post="post" />
      </div>
    </template>

    <EmptyState v-else-if="q" icon="search" title="No results found" :description="`We couldn't find anything for '${q}'. Try a different keyword.`" />
    <EmptyState v-else icon="search" title="Start typing" description="Enter a keyword to search your posts." />
  </div>
</template>
