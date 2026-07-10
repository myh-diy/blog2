<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import api from '../utils/api'
import { usePostsStore } from '../stores/posts'
import PostCard from '../components/PostCard.vue'
import GradientButton from '../components/GradientButton.vue'
import EmptyState from '../components/EmptyState.vue'
import { useSiteAvatar } from '../composables/useSiteAvatar'
import { useSiteTitle } from '../composables/useSiteTitle'

const postStore = usePostsStore()
const { siteAvatar } = useSiteAvatar()
const { siteTitle } = useSiteTitle()

interface P { id: number; text: string; x: number; y: number; c: string; s: number; vx: number; vy: number; o: number }
const particles = ref<P[]>([])
const ready = ref(false)
let raf = 0
let quotePool: string[] = []

const DANMAKU_COLORS = [
  '#ef4444', '#f97316', '#f59e0b', '#84cc16', '#10b981',
  '#06b6d4', '#3b82f6', '#8b5cf6', '#d946ef', '#f43f5e'
]

function randomColor() {
  return DANMAKU_COLORS[Math.floor(Math.random() * DANMAKU_COLORS.length)]
}

onMounted(async () => {
  postStore.fetchPosts(1, '')

  try { const r = await api.get('/quotes'); quotePool = r.data.quotes } catch {}
  if (!quotePool.length) quotePool = ['Stay curious', 'Keep learning', 'Code is poetry', 'Build awesome things', 'Think different', 'Simplicity wins']

  const items: P[] = []
  for (let i = 0; i < 12; i++) {
    items.push({
      id: i,
      text: quotePool[i % quotePool.length],
      x: -20 + Math.random() * 130,
      y: Math.random() * 80 + 10,
      c: randomColor(),
      s: 14 + Math.floor(Math.random() * 16),
      vx: 0.08 + Math.random() * 0.12,
      vy: 0,
      o: 0.6 + Math.random() * 0.35,
    })
  }
  particles.value = items
  ready.value = true
  animate()
})

onUnmounted(() => cancelAnimationFrame(raf))

function resetParticle(p: P): P {
  return {
    ...p,
    text: quotePool[Math.floor(Math.random() * quotePool.length)],
    x: -20 - Math.random() * 30,
    y: Math.random() * 80 + 10,
    c: randomColor(),
    s: 14 + Math.floor(Math.random() * 16),
    vx: 0.08 + Math.random() * 0.12,
    o: 0.6 + Math.random() * 0.35,
  }
}

function animate() {
  particles.value = particles.value.map(p => {
    let x = p.x + p.vx
    if (x > 120) {
      return resetParticle(p)
    }
    return { ...p, x }
  })
  raf = requestAnimationFrame(animate)
}
</script>

<template>
  <div class="space-y-16">
    <!-- Hero -->
    <section class="relative overflow-hidden rounded-3xl bg-white/80 dark:bg-slate-900/80 backdrop-blur-sm border border-gray-100 dark:border-white/5 p-8 md:p-12 min-h-[420px] flex items-center">
      <!-- Floating quotes -->
      <div v-if="ready" class="absolute inset-0 overflow-hidden pointer-events-none">
        <div v-for="p in particles" :key="p.id"
          class="absolute font-medium whitespace-nowrap select-none transition-none"
          :style="{ left: p.x + '%', top: p.y + '%', fontSize: p.s + 'px', color: p.c, opacity: p.o }">
          {{ p.text }}
        </div>
      </div>

      <div class="relative z-10 flex flex-col md:flex-row items-center gap-8 w-full">
        <div class="flex-1 text-center md:text-left">
          <h1 class="text-4xl md:text-6xl font-black text-brand-600 dark:text-brand-400 leading-tight">
            {{ siteTitle }}
          </h1>
          <p class="mt-4 text-lg text-slate-500 dark:text-slate-400">Learning in public, one post at a time.</p>
          <div class="mt-8 flex flex-wrap justify-center md:justify-start gap-3">
            <router-link to="/posts">
              <GradientButton>Explore posts</GradientButton>
            </router-link>
            <router-link to="/search"
              class="px-5 py-2.5 rounded-xl font-semibold text-slate-600 dark:text-slate-300
                     bg-white dark:bg-slate-800 border border-gray-200 dark:border-white/10
                     hover:border-brand-300 dark:hover:border-brand-700 hover:shadow-md transition-all">
              Search
            </router-link>
          </div>
        </div>
        <div class="w-56 h-56 md:w-80 md:h-80 shrink-0">
          <div class="w-full h-full rounded-full overflow-hidden border-8 border-white/60 dark:border-slate-800/60 shadow-2xl bg-gray-100 dark:bg-slate-800">
            <img :src="siteAvatar" alt="site avatar" class="w-full h-full object-cover" />
          </div>
        </div>
      </div>
    </section>

    <!-- Latest posts -->
    <section>
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold text-slate-800 dark:text-slate-100">Latest Posts</h2>
        <router-link to="/posts" class="text-sm font-semibold text-brand-600 dark:text-brand-400 hover:underline">View all →</router-link>
      </div>

      <div v-if="postStore.loading" class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="i in 6" :key="i" class="h-48 rounded-2xl bg-gray-200 dark:bg-white/5 animate-pulse"></div>
      </div>

      <div v-else-if="postStore.posts.length" class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        <PostCard v-for="post in postStore.posts.slice(0, 6)" :key="post.id" :post="post" />
      </div>

      <EmptyState v-else icon="sad" title="No posts yet" description="Upload your first markdown post in admin.">
        <router-link to="/admin" class="mt-4">
          <GradientButton>Go to admin</GradientButton>
        </router-link>
      </EmptyState>
    </section>
  </div>
</template>
