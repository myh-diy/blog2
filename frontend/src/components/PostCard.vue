<script setup lang="ts">
import type { Post } from '../stores/posts'
defineProps<{ post: Post }>()
</script>

<template>
  <article class="group flex flex-col h-full bg-white dark:bg-slate-900 rounded-2xl border border-gray-100 dark:border-white/5
                   shadow-sm hover:shadow-xl hover:shadow-brand-500/10
                   hover:-translate-y-1 transition-all duration-300 overflow-hidden">
    <!-- Cover placeholder -->
    <div class="h-36 bg-gradient-to-br from-brand-100 to-accent-100 dark:from-brand-900/30 dark:to-accent-900/30 relative overflow-hidden">
      <div class="absolute inset-0 opacity-30 bg-[radial-gradient(circle_at_1px_1px,rgba(255,255,255,0.6)_1px,transparent_0)] bg-[length:20px_20px]"></div>
    </div>

    <div class="flex flex-col flex-1 p-5">
      <div class="flex items-center gap-2 text-xs text-slate-400 dark:text-slate-500 mb-3">
        <time>{{ new Date(post.created_at).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' }) }}</time>
        <span aria-hidden="true">·</span>
        <span>{{ post.tags.length }} tags</span>
      </div>

      <router-link :to="`/post/${post.slug}`" class="block flex-1">
        <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 group-hover:text-brand-600 dark:group-hover:text-brand-400 transition-colors line-clamp-2 leading-snug">
          {{ post.title }}
        </h2>
      </router-link>

      <div v-if="post.tags.length" class="flex gap-2 flex-wrap mt-4">
        <router-link v-for="tag in post.tags" :key="tag.id" :to="`/posts?tag=${tag.name}`" class="tag-pill">
          #{{ tag.name }}
        </router-link>
      </div>
    </div>
  </article>
</template>
