<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import api from '../utils/api'
import { usePostsStore } from '../stores/posts'
import PostCard from '../components/PostCard.vue'
import GradientButton from '../components/GradientButton.vue'
import KawaiiIcon from '../components/KawaiiIcon.vue'
import EmptyState from '../components/EmptyState.vue'

const postStore = usePostsStore()

interface P { id: number; text: string; x: number; y: number; c: string; s: number; vx: number; vy: number; o: number }
const particles = ref<P[]>([])
const ready = ref(false)
const C = ['#fb7299', '#ff5fb8', '#ff9dd4', '#a78bfa', '#8b5cf6', '#c4b5fd', '#ffc9e9']
let raf = 0

onMounted(async () => {
  postStore.fetchPosts(1, '')

  let quotes: string[] = []
  try { const r = await api.get('/quotes'); quotes = r.data.quotes } catch {}
  if (!quotes.length) quotes = ['Stay curious', 'Keep learning', 'Code is poetry', 'Build awesome things', 'Think different', 'Simplicity wins']

  const items: P[] = []
  for (let i = 0; i < 12; i++) {
    items.push({
      id: i,
      text: quotes[i % quotes.length],
      x: Math.random() * 90 + 5,
      y: Math.random() * 80 + 10,
      c: C[i % C.length],
      s: 12 + Math.floor(Math.random() * 16),
      vx: 0.02 + Math.random() * 0.06,
      vy: 0.01 + Math.random() * 0.04,
      o: 0.15 + Math.random() * 0.25,
    })
  }
  particles.value = items
  ready.value = true
  animate()
})

onUnmounted(() => cancelAnimationFrame(raf))

function animate() {
  particles.value = particles.value.map(p => {
    let x = p.x + p.vx
    let y = p.y + p.vy
    if (x < -10 || x > 110) p.vx = -p.vx
    if (y < -10 || y > 110) p.vy = -p.vy
    return { ...p, x, y }
  })
  raf = requestAnimationFrame(animate)
}
</script>

<template>
  <div class="space-y-16">
    <!-- Hero -->
    <section class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-brand-100 via-white to-accent-100 dark:from-brand-900/20 dark:via-slate-900 dark:to-accent-900/20 p-8 md:p-14 min-h-[360px] flex items-center">
      <!-- Floating quotes -->
      <div v-if="ready" class="absolute inset-0 overflow-hidden pointer-events-none">
        <div v-for="p in particles" :key="p.id"
          class="absolute font-medium whitespace-nowrap select-none transition-none"
          :style="{ left: p.x + '%', top: p.y + '%', fontSize: p.s + 'px', color: p.c, opacity: p.o }">
          {{ p.text }}
        </div>
      </div>

      <div class="relative z-10 flex flex-col md:flex-row items-center gap-10 w-full">
        <div class="flex-1 text-center md:text-left">
          <h1 class="text-4xl md:text-6xl font-black text-slate-800 dark:text-slate-100 leading-tight">
            My <span class="text-transparent bg-clip-text bg-gradient-to-r from-brand-400 to-accent-500">Blog</span>
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
        <div class="w-40 h-40 md:w-56 md:h-56 shrink-0">
          <KawaiiIcon name="wave" />
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
