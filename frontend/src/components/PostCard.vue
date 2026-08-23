<script setup lang="ts">
import { computed } from 'vue'
import { Clock3, Eye, MessageSquare } from '@lucide/vue'
import type { Post } from '../stores/posts'

const props = withDefaults(defineProps<{ post: Post; compact?: boolean }>(), { compact: false })

const excerpt = computed(() => props.post.content_html
  .replace(/<style[\s\S]*?<\/style>/gi, '')
  .replace(/<script[\s\S]*?<\/script>/gi, '')
  .replace(/<[^>]+>/g, ' ')
  .replace(/&nbsp;/g, ' ')
  .replace(/&amp;/g, '&')
  .replace(/&lt;/g, '<')
  .replace(/&gt;/g, '>')
  .replace(/\s+/g, ' ')
  .trim()
  .slice(0, 150))

const readingMinutes = computed(() => Math.max(1, Math.ceil(excerpt.value.length / 120)))
</script>

<template>
  <article v-if="compact" class="group border-b border-slate-100 px-5 py-5 transition-colors hover:bg-slate-50/70 dark:border-white/10 dark:hover:bg-white/[0.025]">
    <div class="flex min-w-0 items-center gap-4">
      <div class="min-w-0 flex-1">
        <div class="mb-2 flex flex-wrap items-center gap-2 text-xs text-slate-400 dark:text-slate-500">
          <span>博主</span><span class="text-slate-200 dark:text-white/10">|</span>
          <time>{{ new Date(post.created_at).toLocaleDateString('zh-CN') }}</time>
          <template v-if="post.tags.length"><span class="text-slate-200 dark:text-white/10">|</span><span>{{ post.tags.map(tag => tag.name).slice(0, 2).join(' · ') }}</span></template>
        </div>
        <router-link :to="`/post/${post.slug}`" class="block">
          <h2 class="line-clamp-1 text-base font-semibold leading-7 text-slate-900 transition-colors group-hover:text-brand-600 dark:text-slate-100 dark:group-hover:text-brand-300">{{ post.title }}</h2>
          <p v-if="excerpt" class="mt-1 line-clamp-2 text-sm leading-6 text-slate-500 dark:text-slate-400">{{ excerpt }}</p>
        </router-link>
        <div class="mt-3 flex items-center gap-5 text-xs text-slate-400 dark:text-slate-500">
          <span class="inline-flex items-center gap-1"><Eye :size="14" /> 阅读</span>
          <span class="inline-flex items-center gap-1"><Clock3 :size="14" /> {{ readingMinutes }} 分钟</span>
          <span class="inline-flex items-center gap-1"><MessageSquare :size="14" /> 0</span>
        </div>
      </div>
      <router-link v-if="post.cover_image" :to="`/post/${post.slug}`" class="h-20 w-28 shrink-0 overflow-hidden bg-slate-100 sm:h-24 sm:w-36 dark:bg-slate-800">
        <img :src="post.cover_image" :alt="post.title" class="h-full w-full object-cover" />
      </router-link>
    </div>
  </article>

  <article v-else class="surface group flex h-full flex-col overflow-hidden transition-colors hover:border-brand-300 dark:hover:border-brand-700">
    <div class="relative flex h-36 items-center justify-center overflow-hidden bg-slate-100 dark:bg-slate-800">
      <img v-if="post.cover_image" :src="post.cover_image" :alt="post.title" class="absolute inset-0 h-full w-full object-cover" />
      <div v-if="!post.cover_image" class="absolute inset-0 opacity-30 bg-[radial-gradient(circle_at_1px_1px,rgba(255,255,255,0.6)_1px,transparent_0)] bg-[length:20px_20px]"></div>
      <div v-if="!post.cover_image" class="relative z-10 flex flex-col items-center gap-2">
        <div v-if="post.tags.length" class="bg-white/90 px-4 py-1.5 text-sm font-bold text-brand-600 dark:bg-slate-900/80 dark:text-brand-300">#{{ post.tags[0].name }}</div>
        <div v-else class="flex h-12 w-12 items-center justify-center bg-white text-xl font-bold text-brand-500 dark:bg-slate-900">{{ post.title.charAt(0).toUpperCase() }}</div>
      </div>
    </div>
    <div class="flex flex-1 flex-col p-5">
      <div class="mb-3 flex items-center gap-2 text-xs text-slate-400 dark:text-slate-500">
        <time>{{ new Date(post.created_at).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' }) }}</time>
        <span aria-hidden="true">·</span>
        <span>{{ post.tags.length }} tags</span>
      </div>
      <router-link :to="`/post/${post.slug}`" class="block flex-1">
        <h2 class="line-clamp-2 text-lg font-bold leading-snug text-slate-800 transition-colors group-hover:text-brand-600 dark:text-slate-100 dark:group-hover:text-brand-400">{{ post.title }}</h2>
      </router-link>
      <div v-if="post.tags.length" class="mt-4 flex flex-wrap gap-2">
        <router-link v-for="tag in post.tags" :key="tag.id" :to="`/posts?tag=${tag.name}`" class="tag-pill">#{{ tag.name }}</router-link>
      </div>
    </div>
  </article>
</template>
