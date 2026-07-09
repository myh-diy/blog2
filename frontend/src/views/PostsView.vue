<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
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
  try { const r = await api.get('/tags'); allTags.value = r.data.tags } catch {}
})

watch(() => route.query, () => {
  store.fetchPosts(Number(route.query.page) || 1, (route.query.tag as string) || '')
})
</script>

<template>
  <div>
    <!-- Mobile tag filter -->
    <div class="lg:hidden mb-6">
      <div class="flex items-center gap-2 overflow-x-auto pb-2 scrollbar-thin">
        <button @click="$router.push('/posts')"
          class="shrink-0 px-4 py-2 rounded-full text-sm font-medium transition-all"
          :class="!route.query.tag
            ? 'bg-brand-500 text-white shadow-md'
            : 'bg-white dark:bg-slate-900 text-slate-600 dark:text-slate-300 border border-gray-200 dark:border-white/10'">
          All
        </button>
        <button v-for="t in allTags" :key="t.name" @click="$router.push({ path: '/posts', query: { tag: t.name } })"
          class="shrink-0 px-4 py-2 rounded-full text-sm font-medium transition-all"
          :class="route.query.tag === t.name
            ? 'bg-brand-500 text-white shadow-md'
            : 'bg-white dark:bg-slate-900 text-slate-600 dark:text-slate-300 border border-gray-200 dark:border-white/10'">
          #{{ t.name }} <span class="opacity-70 ml-1">{{ t.count }}</span>
        </button>
      </div>
    </div>

    <div class="flex flex-col lg:flex-row gap-10">
      <!-- Posts grid -->
      <div class="flex-1 min-w-0">
        <h1 class="text-3xl font-black text-slate-800 dark:text-slate-100 mb-6">
          {{ route.query.tag ? `#${route.query.tag}` : 'All Posts' }}
        </h1>

        <div v-if="store.loading" class="grid md:grid-cols-2 gap-6">
          <div v-for="i in 4" :key="i" class="h-52 rounded-2xl bg-gray-200 dark:bg-white/5 animate-pulse"></div>
        </div>

        <EmptyState v-else-if="!store.posts.length" icon="sad" title="No posts found" description="Try another tag or upload a new post." />

        <template v-else>
          <div class="grid md:grid-cols-2 gap-6">
            <PostCard v-for="post in store.posts" :key="post.id" :post="post" />
          </div>

          <div v-if="totalPages > 1" class="flex justify-center gap-2 mt-10">
            <button v-for="p in totalPages" :key="p" @click="$router.push({ query: { ...route.query, page: p } })"
              class="w-10 h-10 rounded-full text-sm font-semibold transition-all"
              :class="p === currentPage
                ? 'bg-brand-500 text-white shadow-md'
                : 'bg-white dark:bg-slate-900 text-slate-600 dark:text-slate-300 border border-gray-200 dark:border-white/10 hover:border-brand-300 dark:hover:border-brand-700 hover:shadow-sm'">
              {{ p }}
            </button>
          </div>
        </template>
      </div>

      <!-- Sidebar -->
      <aside class="hidden lg:block lg:w-64 shrink-0">
        <div class="sticky top-24 bg-white dark:bg-slate-900 rounded-2xl border border-gray-100 dark:border-white/5 p-5 shadow-sm">
          <h3 class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider mb-4">Tags</h3>
          <div class="space-y-1 max-h-[70vh] overflow-y-auto pr-1 scrollbar-thin">
            <button @click="$router.push('/posts')"
              class="w-full flex items-center justify-between px-3 py-2 rounded-xl text-sm transition-all"
              :class="!route.query.tag
                ? 'bg-brand-50 dark:bg-brand-900/20 text-brand-700 dark:text-brand-300 font-semibold'
                : 'text-slate-600 dark:text-slate-300 hover:bg-gray-50 dark:hover:bg-white/5'">
              <span>All posts</span>
            </button>
            <button v-for="t in allTags" :key="t.name" @click="$router.push({ path: '/posts', query: { tag: t.name } })"
              class="w-full flex items-center justify-between px-3 py-2 rounded-xl text-sm transition-all"
              :class="route.query.tag === t.name
                ? 'bg-brand-50 dark:bg-brand-900/20 text-brand-700 dark:text-brand-300 font-semibold'
                : 'text-slate-600 dark:text-slate-300 hover:bg-gray-50 dark:hover:bg-white/5'">
              <span>#{{ t.name }}</span>
              <span class="text-xs px-2 py-0.5 rounded-full bg-gray-100 dark:bg-white/10 text-slate-400">{{ t.count }}</span>
            </button>
          </div>
        </div>
      </aside>
    </div>
  </div>
</template>
