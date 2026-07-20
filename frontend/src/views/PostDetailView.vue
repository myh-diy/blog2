<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { usePostsStore, type Post } from '../stores/posts'
import { useAuthStore } from '../stores/auth'
import api from '../utils/api'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'
import TOCSidebar from '../components/TOCSidebar.vue'

const route = useRoute()
const store = usePostsStore()
const auth = useAuthStore()
const post = ref<Post | null>(null)
const isEditing = ref(false)
const editorContent = ref('')
const editorLoading = ref(false)
const saving = ref(false)
const exporting = ref(false)
const actionError = ref('')

onMounted(async () => {
  post.value = await store.fetchPost(route.params.slug as string)
})

async function startEditing() {
  if (!post.value || !auth.isAuthenticated) return
  editorLoading.value = true
  actionError.value = ''
  try {
    const response = await api.get(`/admin/posts/${post.value.id}/source`)
    editorContent.value = response.data.content
    isEditing.value = true
  } catch {
    actionError.value = '无法加载 Markdown 原文。'
  } finally {
    editorLoading.value = false
  }
}

async function saveContent() {
  if (!post.value || !editorContent.value.trim()) return
  saving.value = true
  actionError.value = ''
  try {
    const response = await api.put(`/admin/posts/${post.value.id}/content`, {
      content: editorContent.value,
    })
    post.value = response.data.post
    isEditing.value = false
  } catch (error: any) {
    actionError.value = error.response?.data?.error || '保存失败，请稍后重试。'
  } finally {
    saving.value = false
  }
}

function cancelEditing() {
  isEditing.value = false
  editorContent.value = ''
  actionError.value = ''
}

async function exportMarkdown() {
  if (!post.value) return
  exporting.value = true
  actionError.value = ''
  try {
    const response = await api.get(`/posts/${post.value.slug}/export`, { responseType: 'blob' })
    const url = URL.createObjectURL(response.data)
    const link = document.createElement('a')
    const safeTitle = post.value.title.replace(/[\\/:*?"<>|]/g, '-').trim() || 'post'
    link.href = url
    link.download = `${safeTitle}.md`
    document.body.appendChild(link)
    link.click()
    link.remove()
    URL.revokeObjectURL(url)
  } catch {
    actionError.value = '导出失败，请稍后重试。'
  } finally {
    exporting.value = false
  }
}
</script>

<template>
  <div v-if="!post" class="flex justify-center py-20">
    <div class="w-8 h-8 border-2 border-brand-500 border-t-transparent rounded-full animate-spin"></div>
  </div>

  <div v-else class="relative max-w-6xl mx-auto">
    <article class="max-w-5xl mx-auto xl:mx-0">
      <div class="mb-6 flex flex-wrap items-center justify-between gap-3">
        <router-link to="/posts" class="inline-flex items-center gap-1 text-sm text-slate-400 dark:text-slate-500 hover:text-brand-600 dark:hover:text-brand-400 transition-colors">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/></svg>
          Back to posts
        </router-link>

        <div class="flex items-center gap-2">
          <button
            v-if="auth.isAuthenticated && !isEditing"
            type="button"
            :disabled="editorLoading"
            title="在线编辑"
            class="inline-flex h-9 items-center gap-2 border border-gray-200 dark:border-white/10 bg-white dark:bg-slate-900 px-3 text-sm font-medium text-slate-600 dark:text-slate-300 hover:border-brand-300 hover:text-brand-600 disabled:opacity-50 transition-colors"
            @click="startEditing"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16.862 3.487a2.25 2.25 0 113.182 3.182L8.25 18.463 3 20l1.537-5.25L16.862 3.487z"/></svg>
            {{ editorLoading ? '加载中' : '编辑' }}
          </button>
          <button
            type="button"
            :disabled="exporting"
            title="导出 Markdown"
            class="inline-flex h-9 items-center gap-2 bg-brand-500 px-3 text-sm font-semibold text-white hover:bg-brand-600 disabled:opacity-50 transition-colors"
            @click="exportMarkdown"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v12m0 0l-4-4m4 4l4-4M5 21h14"/></svg>
            {{ exporting ? '导出中' : '导出' }}
          </button>
        </div>
      </div>

      <p v-if="actionError" class="mb-4 border border-red-200 dark:border-red-900/50 bg-red-50 dark:bg-red-950/30 p-3 text-sm text-red-600 dark:text-red-400">
        {{ actionError }}
      </p>

      <div class="bg-white dark:bg-slate-900 rounded-2xl border border-gray-100 dark:border-white/5 p-6 md:p-10 shadow-sm mb-8">
        <header class="mb-8">
          <img v-if="post.cover_image" :src="post.cover_image" :alt="post.title" class="w-full h-auto rounded-xl mb-6 object-cover max-h-80" />
          <h1 class="text-3xl md:text-4xl font-black text-slate-800 dark:text-slate-100 leading-tight mb-4">{{ post.title }}</h1>
          <time class="text-sm text-slate-400 dark:text-slate-500 mb-4 block">
            {{ new Date(post.created_at).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' }) }}
          </time>
          <div v-if="post.tags.length" class="flex gap-2 flex-wrap">
            <router-link v-for="tag in post.tags" :key="tag.id" :to="`/posts?tag=${tag.name}`" class="tag-pill">#{{ tag.name }}</router-link>
          </div>
        </header>

        <div v-if="isEditing" class="border-t border-gray-100 dark:border-white/10 pt-6">
          <label for="post-markdown-editor" class="mb-2 block text-sm font-semibold text-slate-700 dark:text-slate-200">Markdown 原文</label>
          <textarea
            id="post-markdown-editor"
            v-model="editorContent"
            spellcheck="false"
            class="min-h-[60vh] w-full resize-y border border-gray-200 dark:border-white/10 bg-gray-50 dark:bg-slate-950 p-4 font-mono text-sm leading-6 text-slate-800 dark:text-slate-200 outline-none focus:border-brand-500 focus:ring-2 focus:ring-brand-500/20"
          ></textarea>
          <div class="mt-4 flex justify-end gap-2">
            <button type="button" :disabled="saving" class="h-10 px-4 text-sm font-medium text-slate-500 hover:bg-gray-100 dark:hover:bg-white/5 disabled:opacity-50 transition-colors" @click="cancelEditing">
              取消
            </button>
            <button type="button" :disabled="saving || !editorContent.trim()" class="h-10 bg-brand-500 px-5 text-sm font-semibold text-white hover:bg-brand-600 disabled:cursor-not-allowed disabled:opacity-50 transition-colors" @click="saveContent">
              {{ saving ? '保存中...' : '保存文章' }}
            </button>
          </div>
        </div>

        <div v-else class="prose prose-lg max-w-none
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

    <!-- Floating TOC beside the article (desktop) -->
    <aside class="hidden xl:block fixed top-24 w-56 z-40" style="left: calc(50% + 22rem);">
      <div class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-sm rounded-2xl border border-gray-100 dark:border-white/5 p-5 shadow-lg max-h-[calc(100vh-8rem)] overflow-y-auto scrollbar-thin">
        <TOCSidebar :toc-json="post.toc" />
      </div>
    </aside>
  </div>
</template>
