<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import api from '../utils/api'
import PostCard from '../components/PostCard.vue'
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
    results.value = r.data.posts; total.value = r.data.total
  } finally { loading.value = false }
}
</script>

<template>
  <div>
    <h1 class="text-4xl font-bold text-warm-800 dark:text-warm-100 mb-8">Search</h1>
    <form @submit.prevent="search" class="flex gap-3 mb-10">
      <div class="relative flex-1">
        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-warm-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
        <input v-model="q" placeholder="Search posts..."
          class="w-full pl-10 pr-4 py-3 bg-white dark:bg-white/5 border border-warm-200 dark:border-white/10 rounded-xl text-warm-800 dark:text-warm-100 placeholder-warm-400 focus:outline-none focus:ring-2 focus:ring-brand-500/30 dark:focus:ring-pop-500/30 focus:border-brand-500 dark:focus:border-pop-500 transition-all" />
      </div>
      <button type="submit" class="btn-primary px-6">Search</button>
    </form>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="w-8 h-8 border-2 border-brand-500 dark:border-pop-500 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <template v-else-if="results.length">
      <p class="text-sm text-warm-400 dark:text-warm-500 mb-6">{{ total }} result(s)</p>
      <PostCard v-for="post in results" :key="post.id" :post="post" />
    </template>

    <p v-else-if="q" class="text-center py-20 text-warm-400">No results found.</p>
  </div>
</template>
