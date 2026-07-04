<script setup lang="ts">
import type { Post } from '../stores/posts'

defineProps<{ post: Post }>()
</script>

<template>
  <article class="py-6 border-b border-gray-200 dark:border-gray-700">
    <time class="text-sm text-gray-500 dark:text-gray-400">
      {{ new Date(post.created_at).toLocaleDateString('zh-CN') }}
    </time>
    <router-link :to="`/post/${post.slug}`">
      <h2 class="text-xl font-semibold mt-1 hover:text-blue-500">{{ post.title }}</h2>
    </router-link>
    <div class="flex gap-2 mt-2">
      <router-link
        v-for="tag in post.tags"
        :key="tag.id"
        :to="`/?tag=${tag.name}`"
        class="text-xs px-2 py-0.5 bg-gray-100 dark:bg-gray-800 rounded hover:bg-blue-100"
      >
        {{ tag.name }}
      </router-link>
    </div>
  </article>
</template>
