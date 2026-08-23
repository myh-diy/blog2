<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { BookOpen, ChevronRight, Clock3, Hash, Layers3, Sparkles } from '@lucide/vue'
import api from '../../utils/api'
import type { Post } from '../../stores/posts'
import PostCard from '../../components/PostCard.vue'
import { useSiteTitle } from '../../composables/useSiteTitle'
import { useSiteAvatar } from '../../composables/useSiteAvatar'

interface Poem {
  content: string
  origin: string
  author: string
}

const { siteTitle } = useSiteTitle()
const { siteAvatar } = useSiteAvatar()
const posts = ref<Post[]>([])
const tags = ref<{ name: string; count: number }[]>([])
const poem = ref<Poem | null>(null)
const loading = ref(true)
const totalTags = computed(() => tags.value.length)

onMounted(async () => {
  const [postsResult, tagsResult, poemResult] = await Promise.allSettled([
    api.get('/posts', { params: { page: 1, per_page: 6 } }),
    api.get('/tags'),
    fetch('https://v1.jinrishici.com/all.json', { cache: 'no-store' }).then(async response => {
      if (!response.ok) throw new Error(`Poetry API returned ${response.status}`)
      return response.json() as Promise<Poem>
    }),
  ])

  if (postsResult.status === 'fulfilled') posts.value = postsResult.value.data.posts ?? []
  if (tagsResult.status === 'fulfilled') tags.value = tagsResult.value.data.tags ?? []
  if (poemResult.status === 'fulfilled') poem.value = poemResult.value
  loading.value = false
})
</script>

<template>
  <div class="grid items-start gap-4 lg:grid-cols-[minmax(0,1fr)_280px]">
    <section class="surface min-w-0 overflow-hidden">
      <div class="flex h-12 items-center border-b border-slate-100 px-5 dark:border-white/10">
        <span class="border-r border-slate-200 pr-4 text-sm font-medium text-brand-600 dark:border-white/10">推荐</span>
        <router-link to="/posts" class="px-4 text-sm text-slate-500 hover:text-brand-600 dark:text-slate-400">最新</router-link>
        <router-link to="/timeline" class="px-4 text-sm text-slate-500 hover:text-brand-600 dark:text-slate-400">归档</router-link>
      </div>

      <div v-if="loading" class="divide-y divide-slate-100 dark:divide-white/10">
        <div v-for="item in 5" :key="item" class="h-36 animate-pulse bg-slate-50 dark:bg-white/[0.02]"></div>
      </div>
      <div v-else-if="posts.length">
        <PostCard v-for="post in posts" :key="post.id" :post="post" compact />
        <router-link to="/posts" class="flex h-12 items-center justify-center gap-1 border-t border-slate-100 text-sm text-slate-500 hover:text-brand-600 dark:border-white/10 dark:text-slate-400">
          查看全部文章 <ChevronRight :size="15" />
        </router-link>
      </div>
      <div v-else class="flex min-h-72 flex-col items-center justify-center text-slate-400">
        <BookOpen :size="32" class="mb-3" />
        <p class="text-sm">还没有发布文章</p>
      </div>
    </section>

    <aside class="space-y-4 lg:sticky lg:top-[116px]">
      <section class="surface p-5">
        <div class="flex items-center gap-3">
          <img :src="siteAvatar" alt="站点头像" class="h-12 w-12 bg-brand-50 object-contain p-1 dark:bg-brand-900/20" />
          <div class="min-w-0">
            <h1 class="truncate text-base font-semibold text-slate-900 dark:text-white">{{ siteTitle }}</h1>
            <p class="mt-1 text-xs text-slate-400">持续学习，认真记录</p>
          </div>
        </div>
        <div class="mt-5 grid grid-cols-2 border-t border-slate-100 pt-4 text-center dark:border-white/10">
          <div><strong class="block text-lg text-slate-900 dark:text-white">{{ posts.length }}</strong><span class="text-xs text-slate-400">最近文章</span></div>
          <div class="border-l border-slate-100 dark:border-white/10"><strong class="block text-lg text-slate-900 dark:text-white">{{ totalTags }}</strong><span class="text-xs text-slate-400">技术标签</span></div>
        </div>
      </section>

      <section v-if="poem" class="surface p-5">
        <div class="mb-3 flex items-center gap-2 text-sm font-semibold text-slate-800 dark:text-slate-100"><Sparkles :size="16" class="text-amber-500" /> 今日诗句</div>
        <p class="text-sm leading-7 text-slate-600 dark:text-slate-300">{{ poem.content }}</p>
        <p class="mt-3 text-xs text-slate-400">{{ poem.author }} · {{ poem.origin }}</p>
      </section>

      <section class="surface p-5">
        <h2 class="mb-3 flex items-center gap-2 text-sm font-semibold text-slate-800 dark:text-slate-100"><Hash :size="16" class="text-brand-500" /> 热门标签</h2>
        <div class="flex flex-wrap gap-2">
          <router-link v-for="tag in tags.slice(0, 10)" :key="tag.name" :to="`/posts?tag=${tag.name}`" class="tag-pill">{{ tag.name }} {{ tag.count }}</router-link>
          <span v-if="!tags.length" class="text-xs text-slate-400">暂无标签</span>
        </div>
      </section>

      <div class="flex items-center justify-center gap-4 text-xs text-slate-400">
        <span class="inline-flex items-center gap-1"><Layers3 :size="13" /> Vue + Go</span>
        <span class="inline-flex items-center gap-1"><Clock3 :size="13" /> 持续更新</span>
      </div>
    </aside>
  </div>
</template>
