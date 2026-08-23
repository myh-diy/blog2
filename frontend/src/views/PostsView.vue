<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { usePostsStore } from '../stores/posts'
import api from '../utils/api'
import PostCard from '../components/PostCard.vue'
import EmptyState from '../components/EmptyState.vue'
import { ChevronLeft, ChevronRight, Hash } from '@lucide/vue'

const route = useRoute()
const store = usePostsStore()
const allTags = ref<{ name: string; count: number }[]>([])
const currentPage = computed(() => Number(route.query.page) || 1)
const totalPages = computed(() => Math.max(1, Math.ceil(store.total / 10)))

onMounted(async () => {
  store.fetchPosts(1, (route.query.tag as string) || '')
  try {
    const response = await api.get('/tags')
    allTags.value = response.data.tags ?? []
  } catch {}
})

watch(() => route.query, () => {
  store.fetchPosts(Number(route.query.page) || 1, (route.query.tag as string) || '')
})
</script>

<template>
  <div class="grid items-start gap-4 lg:grid-cols-[minmax(0,1fr)_260px]">
    <section class="surface min-w-0 overflow-hidden">
      <header class="surface-header flex items-center justify-between gap-4">
        <div>
          <h1 class="page-title">{{ route.query.tag ? `# ${route.query.tag}` : '全部文章' }}</h1>
          <p class="mt-1 text-xs text-slate-400">共 {{ store.total }} 篇内容</p>
        </div>
        <router-link to="/search" class="text-sm text-brand-600 hover:text-brand-700">搜索文章</router-link>
      </header>

      <div v-if="store.loading" class="divide-y divide-slate-100 dark:divide-white/10">
        <div v-for="i in 5" :key="i" class="h-36 animate-pulse bg-slate-50 dark:bg-white/[0.02]"></div>
      </div>
      <EmptyState v-else-if="!store.posts.length" icon="sad" title="没有找到文章" description="换个标签试试。" />
      <div v-else><PostCard v-for="post in store.posts" :key="post.id" :post="post" compact /></div>

      <div v-if="totalPages > 1" class="flex items-center justify-center gap-1 border-t border-slate-100 p-5 dark:border-white/10">
        <button class="icon-button" :disabled="currentPage <= 1" @click="$router.push({ query: { ...route.query, page: currentPage - 1 } })"><ChevronLeft :size="17" /></button>
        <button v-for="page in totalPages" :key="page" class="h-9 min-w-9 px-2 text-sm"
          :class="page === currentPage ? 'bg-brand-500 text-white' : 'text-slate-500 hover:bg-slate-100 dark:text-slate-400 dark:hover:bg-white/5'"
          @click="$router.push({ query: { ...route.query, page } })">{{ page }}</button>
        <button class="icon-button" :disabled="currentPage >= totalPages" @click="$router.push({ query: { ...route.query, page: currentPage + 1 } })"><ChevronRight :size="17" /></button>
      </div>
    </section>

    <aside class="surface p-5 lg:sticky lg:top-[116px]">
      <h2 class="mb-4 flex items-center gap-2 text-sm font-semibold text-slate-900 dark:text-white"><Hash :size="16" class="text-brand-500" /> 内容分类</h2>
      <div class="space-y-1">
        <button class="flex w-full items-center justify-between px-3 py-2 text-left text-sm" :class="!route.query.tag ? 'bg-brand-50 text-brand-600 dark:bg-brand-900/20' : 'text-slate-500 hover:bg-slate-50 dark:text-slate-400 dark:hover:bg-white/5'" @click="$router.push('/posts')"><span>全部</span><span>{{ store.total }}</span></button>
        <button v-for="tag in allTags" :key="tag.name" class="flex w-full items-center justify-between px-3 py-2 text-left text-sm" :class="route.query.tag === tag.name ? 'bg-brand-50 text-brand-600 dark:bg-brand-900/20' : 'text-slate-500 hover:bg-slate-50 dark:text-slate-400 dark:hover:bg-white/5'" @click="$router.push({ path: '/posts', query: { tag: tag.name } })"><span class="truncate">{{ tag.name }}</span><span>{{ tag.count }}</span></button>
      </div>
    </aside>
  </div>
</template>
