<script setup lang="ts">
import { computed } from 'vue'
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
</script>

<template>
  <article v-if="compact" class="group border-b border-gray-200 py-6 first:border-t dark:border-white/10">
    <div class="flex min-w-0 items-start gap-5">
      <div class="min-w-0 flex-1">
        <router-link :to="`/post/${post.slug}`" class="block">
          <h2 class="text-xl font-bold leading-snug text-slate-800 transition-colors group-hover:text-brand-600 dark:text-slate-100 dark:group-hover:text-brand-400">{{ post.title }}</h2>
          <p v-if="excerpt" class="mt-2 line-clamp-2 text-sm leading-6 text-slate-500 dark:text-slate-400">{{ excerpt }}</p>
        </router-link>
        <div class="mt-3 flex flex-wrap items-center gap-x-4 gap-y-2 text-sm text-slate-400 dark:text-slate-500">
          <time>{{ new Date(post.created_at).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' }) }}</time>
          <div v-if="post.tags.length" class="flex flex-wrap gap-2">
            <router-link v-for="tag in post.tags" :key="tag.id" :to="`/posts?tag=${tag.name}`" class="bg-gray-100 px-2 py-1 text-xs text-slate-500 transition-colors hover:text-brand-600 dark:bg-white/5 dark:text-slate-400">{{ tag.name }}</router-link>
          </div>
        </div>
      </div>
      <router-link v-if="post.cover_image" :to="`/post/${post.slug}`" class="h-20 w-28 shrink-0 overflow-hidden bg-gray-100 sm:h-24 sm:w-40 dark:bg-slate-800">
        <img :src="post.cover_image" :alt="post.title" class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105" />
      </router-link>
    </div>
  </article>

  <article v-else class="group flex h-full flex-col overflow-hidden rounded-2xl border border-gray-100 bg-white shadow-sm transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:shadow-brand-500/10 dark:border-white/5 dark:bg-slate-900">
    <div class="relative flex h-36 items-center justify-center overflow-hidden bg-brand-100 dark:bg-brand-900/30">
      <img v-if="post.cover_image" :src="post.cover_image" :alt="post.title" class="absolute inset-0 h-full w-full object-cover" />
      <div v-if="!post.cover_image" class="absolute inset-0 opacity-30 bg-[radial-gradient(circle_at_1px_1px,rgba(255,255,255,0.6)_1px,transparent_0)] bg-[length:20px_20px]"></div>
      <div v-if="!post.cover_image" class="relative z-10 flex flex-col items-center gap-2">
        <div v-if="post.tags.length" class="rounded-full bg-white/80 px-4 py-1.5 text-sm font-bold text-brand-600 shadow-sm backdrop-blur-sm dark:bg-slate-900/60 dark:text-brand-300">#{{ post.tags[0].name }}</div>
        <div v-else class="flex h-14 w-14 items-center justify-center rounded-2xl bg-white/80 text-2xl font-black text-brand-500 shadow-sm backdrop-blur-sm dark:bg-slate-900/60">{{ post.title.charAt(0).toUpperCase() }}</div>
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
