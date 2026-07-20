<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { usePostsStore } from '../stores/posts'
import api from '../utils/api'
import PostCard from '../components/PostCard.vue'
import EmptyState from '../components/EmptyState.vue'

const route = useRoute()
const store = usePostsStore()
const allTags = ref<{ name: string; count: number }[]>([])
const currentPage = computed(() => Number(route.query.page) || 1)
const totalPages = computed(() => Math.max(1, Math.ceil(store.total / 10)))

onMounted(async () => {
  store.fetchPosts(1, (route.query.tag as string) || '')
  try {
    const response = await api.get('/tags')
    allTags.value = response.data.tags
  } catch {}
})

watch(() => route.query, () => {
  store.fetchPosts(Number(route.query.page) || 1, (route.query.tag as string) || '')
})
</script>

<template>
  <div class="mx-auto max-w-5xl">
    <div class="mb-6 flex items-end justify-between gap-4">
      <div>
        <p class="mb-1 text-sm font-medium text-slate-400 dark:text-slate-500">文章归档</p>
        <h1 class="text-3xl font-black text-slate-800 dark:text-slate-100">{{ route.query.tag ? `#${route.query.tag}` : '全部文章' }}</h1>
      </div>
      <span class="text-sm text-slate-400 dark:text-slate-500">{{ store.total }} 篇</span>
    </div>

    <div class="mb-6 flex items-center gap-2 overflow-x-auto pb-2 scrollbar-thin">
      <button
        class="shrink-0 px-3 py-1.5 text-sm font-medium transition-colors"
        :class="!route.query.tag ? 'bg-brand-500 text-white' : 'bg-white text-slate-500 hover:text-brand-600 dark:bg-slate-900 dark:text-slate-400'"
        @click="$router.push('/posts')"
      >
        全部
      </button>
      <button
        v-for="tag in allTags"
        :key="tag.name"
        class="shrink-0 px-3 py-1.5 text-sm font-medium transition-colors"
        :class="route.query.tag === tag.name ? 'bg-brand-500 text-white' : 'bg-white text-slate-500 hover:text-brand-600 dark:bg-slate-900 dark:text-slate-400'"
        @click="$router.push({ path: '/posts', query: { tag: tag.name } })"
      >
        {{ tag.name }} · {{ tag.count }}
      </button>
    </div>

    <div v-if="store.loading" class="divide-y divide-gray-200 border-y border-gray-200 dark:divide-white/10 dark:border-white/10">
      <div v-for="i in 5" :key="i" class="h-32 animate-pulse bg-gray-100/70 dark:bg-white/5"></div>
    </div>

    <EmptyState v-else-if="!store.posts.length" icon="sad" title="没有找到文章" description="换个标签试试。" />

    <template v-else>
      <div>
        <PostCard v-for="post in store.posts" :key="post.id" :post="post" compact />
      </div>

      <div v-if="totalPages > 1" class="mt-10 flex justify-center gap-2">
        <button
          v-for="page in totalPages"
          :key="page"
          class="h-10 w-10 text-sm font-semibold transition-colors"
          :class="page === currentPage ? 'bg-brand-500 text-white' : 'bg-white text-slate-600 hover:text-brand-600 dark:bg-slate-900 dark:text-slate-300'"
          @click="$router.push({ query: { ...route.query, page } })"
        >
          {{ page }}
        </button>
      </div>
    </template>
  </div>
</template>
