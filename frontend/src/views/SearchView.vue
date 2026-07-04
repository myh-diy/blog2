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

watch(() => route.query.q, (val) => {
  q.value = (val as string) || ''
})

async function search() {
  if (!q.value) return
  loading.value = true
  try {
    const res = await api.get('/search', { params: { q: q.value } })
    results.value = res.data.posts
    total.value = res.data.total
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold mb-8">Search</h1>
    <form @submit.prevent="search" class="flex gap-2 mb-8">
      <input v-model="q" type="text" placeholder="Search posts..."
        class="flex-1 px-3 py-2 border rounded dark:bg-gray-800 dark:border-gray-600" />
      <button type="submit"
        class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
        Search
      </button>
    </form>

    <div v-if="loading" class="text-center py-12 text-gray-500">Searching...</div>

    <template v-else-if="results.length">
      <p class="text-gray-500 mb-4">{{ total }} result(s)</p>
      <PostCard v-for="post in results" :key="post.id" :post="post" />
    </template>
  </div>
</template>
