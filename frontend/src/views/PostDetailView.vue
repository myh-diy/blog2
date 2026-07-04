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
  <div v-if="!post" class="flex justify-center py-20">
    <div class="w-8 h-8 border-2 border-brand-500 dark:border-pop-500 border-t-transparent rounded-full animate-spin"></div>
  </div>

  <div v-else class="lg:grid lg:grid-cols-[1fr_220px] gap-12">
    <article>
      <router-link to="/" class="inline-flex items-center gap-1 text-sm text-warm-400 dark:text-warm-500 hover:text-brand-600 dark:hover:text-pop-400 mb-6 transition-colors">
        &larr; Back to posts
      </router-link>

      <header class="mb-10">
        <h1 class="text-4xl font-bold text-warm-800 dark:text-warm-100 leading-tight mb-4">{{ post.title }}</h1>
        <time class="text-sm text-warm-400 dark:text-warm-500 mb-4 block">
          {{ new Date(post.created_at).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' }) }}
        </time>
        <div v-if="post.tags.length" class="flex gap-2 flex-wrap">
          <router-link v-for="tag in post.tags" :key="tag.id" :to="`/?tag=${tag.name}`" class="tag-pill">#{{ tag.name }}</router-link>
        </div>
      </header>

      <div class="prose prose-lg max-w-none
        prose-headings:text-warm-800 dark:prose-headings:text-warm-100
        prose-a:text-brand-600 dark:prose-a:text-pop-400 prose-a:no-underline hover:prose-a:underline
        prose-code:text-brand-700 dark:prose-code:text-pop-300
        prose-pre:bg-warm-900 dark:prose-pre:bg-[#0d0714] prose-pre:rounded-xl
        prose-blockquote:border-brand-500 dark:prose-blockquote:border-pop-500
        prose-img:rounded-xl prose-img:shadow-md prose-img:mx-auto prose-img:max-w-full">
        <MarkdownRenderer :html="post.content_html" />
      </div>
    </article>

    <aside class="hidden lg:block">
      <div class="sticky top-20">
        <TOCSidebar :toc-json="post.toc" />
      </div>
    </aside>
  </div>
</template>
