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
    <div class="w-8 h-8 border-2 border-brand-500 border-t-transparent rounded-full animate-spin"></div>
  </div>

  <div v-else class="flex flex-col lg:flex-row gap-8">
    <article class="flex-1 min-w-0">
      <router-link to="/posts" class="inline-flex items-center gap-1 text-sm text-slate-400 dark:text-slate-500 hover:text-brand-600 dark:hover:text-brand-400 mb-6 transition-colors">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/></svg>
        Back to posts
      </router-link>

      <div class="bg-white dark:bg-slate-900 rounded-2xl border border-gray-100 dark:border-white/5 p-6 md:p-10 shadow-sm mb-8">
        <header class="mb-8">
          <h1 class="text-3xl md:text-4xl font-black text-slate-800 dark:text-slate-100 leading-tight mb-4">{{ post.title }}</h1>
          <time class="text-sm text-slate-400 dark:text-slate-500 mb-4 block">
            {{ new Date(post.created_at).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' }) }}
          </time>
          <div v-if="post.tags.length" class="flex gap-2 flex-wrap">
            <router-link v-for="tag in post.tags" :key="tag.id" :to="`/posts?tag=${tag.name}`" class="tag-pill">#{{ tag.name }}</router-link>
          </div>
        </header>

        <div class="prose prose-lg max-w-none
          prose-headings:text-slate-800 dark:prose-headings:text-slate-100
          prose-a:text-brand-600 dark:prose-a:text-brand-400 prose-a:no-underline hover:prose-a:underline
          prose-code:text-brand-700 dark:prose-code:text-brand-300
          prose-pre:bg-slate-900 dark:prose-pre:bg-[#0B0F19] prose-pre:rounded-xl
          prose-blockquote:border-brand-500 dark:prose-blockquote:border-brand-400
          prose-img:rounded-xl prose-img:shadow-md prose-img:mx-auto prose-img:max-w-full">
          <MarkdownRenderer :html="post.content_html" />
        </div>
      </div>
    </article>

    <aside class="hidden lg:block lg:w-64 shrink-0">
      <div class="sticky top-24 bg-white dark:bg-slate-900 rounded-2xl border border-gray-100 dark:border-white/5 p-5 shadow-sm">
        <TOCSidebar :toc-json="post.toc" />
      </div>
    </aside>
  </div>
</template>
