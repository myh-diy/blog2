<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { usePostsStore } from '../stores/posts'
import api from '../utils/api'
import PostCard from '../components/PostCard.vue'

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
  <div class="posts-root">
    <div class="posts-layout">
      <!-- Left: posts -->
      <div>
        <h1 v-if="route.query.tag" class="posts-title">#{{ route.query.tag }}</h1>
        <h1 v-else class="posts-title">All Posts</h1>

        <div v-if="store.loading" class="space-y-4">
          <div v-for="i in 3" :key="i" class="animate-pulse space-y-2">
            <div class="h-4 bg-warm-200 dark:bg-white/10 rounded w-24"></div>
            <div class="h-5 bg-warm-200 dark:bg-white/10 rounded w-3/4"></div>
          </div>
        </div>

        <div v-else-if="!store.posts.length" class="py-20 text-center">
          <div class="text-4xl mb-2">📝</div>
          <p class="text-warm-400 dark:text-warm-500">No posts yet.</p>
        </div>

        <template v-else>
          <PostCard v-for="post in store.posts" :key="post.id" :post="post" />
          <div v-if="totalPages > 1" class="flex justify-center gap-2 mt-10">
            <button v-for="p in totalPages" :key="p" @click="$router.push({query:{...route.query,page:p}})"
              :style="{background:p===currentPage?'var(--c-brand)':'var(--c-btn)',color:p===currentPage?'#fff':'var(--c-text)'}"
              class="page-btn">{{ p }}</button>
          </div>
        </template>
      </div>

      <!-- Right: tag sidebar -->
      <aside class="sidebar">
        <nav class="sidebar-nav">
          <h3 class="sidebar-h">Tags</h3>
          <button @click="$router.push('/posts')"
            :style="{background:!route.query.tag?'var(--c-tag)':'',color:!route.query.tag?'var(--c-tag-text)':'var(--c-side)'}"
            class="sidebar-link">All posts</button>
          <button v-for="t in allTags" :key="t.name" @click="$router.push({path:'/posts',query:{tag:t.name}})"
            :style="{background:route.query.tag===t.name?'var(--c-tag)':'',color:route.query.tag===t.name?'var(--c-tag-text)':'var(--c-side)'}"
            class="sidebar-link sb-row">
            <span>#{{ t.name }}</span><span class="sb-num">{{ t.count }}</span>
          </button>
        </nav>
      </aside>
    </div>
  </div>
</template>

<style>
:root {
  --c-brand: #f97316; --c-btn: #fff; --c-text: #57534e;
  --c-side: #78716c; --c-tag: #fff7ed; --c-tag-text: #c2410c;
}
.dark {
  --c-brand: #f43f5e; --c-btn: rgba(255,255,255,.06); --c-text: #a8a29e;
  --c-side: #a8a29e; --c-tag: rgba(244,63,94,.12); --c-tag-text: #fb7185;
}
.posts-root { min-height: 100vh; }
.posts-layout { max-width: 60rem; margin: 0 auto; padding: 2.5rem 1rem 5rem; display: grid; grid-template-columns: 1fr 220px; gap: 2.5rem; }
.posts-title { font-size: 2rem; font-weight: 800; color: var(--c-text); margin-bottom: 1.5rem; }
.dark .posts-title { color: #f5f5f4; }
.sidebar-nav { position: sticky; top: 5rem; }
.sidebar-h { font-size: .7rem; font-weight: 600; color: #a8a29e; text-transform: uppercase; letter-spacing: .08em; margin-bottom: .5rem; }
.sidebar-link { display: flex; width: 100%; text-align: left; padding: .45rem .7rem; border-radius: .5rem; font-size: .875rem; border: none; cursor: pointer; background: transparent; color: var(--c-side); transition: all .1s; }
.sidebar-link:hover { background: var(--c-btn); }
.sb-row { justify-content: space-between; }
.sb-num { font-size: .7rem; opacity: .5; }
.page-btn { width: 2.25rem; height: 2.25rem; border-radius: .5rem; font-weight: 500; font-size: .875rem; border: none; cursor: pointer; }
@media (max-width: 768px) {
  .posts-layout { grid-template-columns: 1fr; }
  .sidebar { display: none; }
}
</style>
