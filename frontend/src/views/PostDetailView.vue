<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { usePostsStore, type Post } from '../stores/posts'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'
import TOCSidebar from '../components/TOCSidebar.vue'

const route = useRoute()
const store = usePostsStore()
const post = ref<Post | null>(null)

onMounted(async () => {
  post.value = await store.fetchPost(route.params.slug as string)
})
</script>

<template>
  <div v-if="!post" class="text-center py-12 text-gray-500">Loading...</div>
  <div v-else class="lg:grid lg:grid-cols-[1fr_200px] gap-8">
    <article>
      <h1 class="text-3xl font-bold mb-2">{{ post.title }}</h1>
      <p class="text-gray-500 text-sm mb-4">
        {{ new Date(post.created_at).toLocaleDateString('zh-CN') }}
      </p>
      <div class="flex gap-2 mb-8">
        <router-link v-for="tag in post.tags" :key="tag.id" :to="`/?tag=${tag.name}`"
          class="text-xs px-2 py-0.5 bg-gray-100 dark:bg-gray-800 rounded hover:bg-blue-100">
          {{ tag.name }}
        </router-link>
      </div>
      <MarkdownRenderer :html="post.content_html" />
    </article>
    <aside class="hidden lg:block">
      <TOCSidebar :toc-json="post.toc" />
    </aside>
  </div>
</template>
