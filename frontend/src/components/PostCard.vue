<script setup lang="ts">
import type { Post } from '../stores/posts'
defineProps<{ post: Post }>()
</script>

<template>
  <article class="group py-7 border-b border-warm-200/60 dark:border-white/5 last:border-0">
    <div class="flex items-center gap-3 text-sm text-warm-400 dark:text-warm-500 mb-3">
      <time>{{ new Date(post.created_at).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' }) }}</time>
      <span aria-hidden="true">·</span>
      <span>{{ post.tags.length }} tags</span>
    </div>
    <router-link :to="`/post/${post.slug}`" class="block">
      <h2 class="text-2xl font-bold text-warm-800 dark:text-warm-100 group-hover:text-brand-600 dark:group-hover:text-pop-400 transition-colors mb-3 leading-tight">
        {{ post.title }}
      </h2>
    </router-link>
    <div v-if="post.tags.length" class="flex gap-2 flex-wrap">
      <router-link v-for="tag in post.tags" :key="tag.id" :to="`/?tag=${tag.name}`" class="tag-pill">
        #{{ tag.name }}
      </router-link>
    </div>
  </article>
</template>
