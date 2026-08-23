<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ArrowLeft, CalendarDays, Clock3, Download, FileUp, PencilLine } from '@lucide/vue'
import { usePostsStore, type Post } from '../stores/posts'
import { useAuthStore } from '../stores/auth'
import api from '../utils/api'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'
import TOCSidebar from '../components/TOCSidebar.vue'
import { useSiteTitle } from '../composables/useSiteTitle'

const route = useRoute()
const store = usePostsStore()
const auth = useAuthStore()
const { siteTitle } = useSiteTitle()
const post = ref<Post | null>(null)
const loaded = ref(false)
const isEditing = ref(false)
const editTitle = ref('')
const editorContent = ref('')
const editorLoading = ref(false)
const saving = ref(false)
const exporting = ref(false)
const actionError = ref('')
const markdownInput = ref<HTMLInputElement | null>(null)
const readingProgress = ref(0)
const readingMinutes = computed(() => {
  if (!post.value) return 1
  const text = post.value.content_html.replace(/<[^>]+>/g, '').replace(/\s+/g, '')
  return Math.max(1, Math.ceil(text.length / 500))
})

onMounted(async () => {
	window.addEventListener('scroll', updateReadingProgress, { passive: true })
  try {
    post.value = await store.fetchPost(route.params.slug as string)
    if (post.value) {
      document.title = `${post.value.title} | ${siteTitle.value}`
      setMetaDescription(post.value.content_html)
    }
  } finally {
    loaded.value = true
  }
})

onUnmounted(() => {
	window.removeEventListener('scroll', updateReadingProgress)
  document.title = siteTitle.value
})

function updateReadingProgress() {
  const scrollable = document.documentElement.scrollHeight - window.innerHeight
  readingProgress.value = scrollable > 0 ? Math.min(100, Math.max(0, window.scrollY / scrollable * 100)) : 0
}

function setMetaDescription(html: string) {
  const text = new DOMParser().parseFromString(html, 'text/html').body.textContent?.replace(/\s+/g, ' ').trim() || ''
  let meta = document.querySelector<HTMLMetaElement>('meta[name="description"]')
  if (!meta) {
    meta = document.createElement('meta')
    meta.name = 'description'
    document.head.appendChild(meta)
  }
  meta.content = text.slice(0, 160)
}

async function startEditing() {
  if (!post.value || !auth.isAuthenticated) return
  editorLoading.value = true
  actionError.value = ''
  try {
    const response = await api.get(`/admin/posts/${post.value.id}/source`)
    editTitle.value = post.value.title
    editorContent.value = response.data.content
    isEditing.value = true
  } catch {
    actionError.value = '无法加载 Markdown 原文。'
  } finally {
    editorLoading.value = false
  }
}

function chooseMarkdownFile() {
  markdownInput.value?.click()
}

async function overwriteFromMarkdown(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file || !post.value) return

  try {
    const content = await file.text()
    if (!content.trim()) throw new Error('empty markdown')
    editorContent.value = content
    editTitle.value = extractMarkdownTitle(content) || post.value.title
    actionError.value = ''
    isEditing.value = true
  } catch {
    actionError.value = '无法读取该 Markdown 文件。'
  } finally {
    input.value = ''
  }
}

function extractMarkdownTitle(content: string) {
  const frontmatter = content.match(/^---\s*\r?\n([\s\S]*?)\r?\n---/)
  const frontmatterTitle = frontmatter?.[1].match(/^title:\s*["']?(.*?)["']?\s*$/m)?.[1]
  if (frontmatterTitle?.trim()) return frontmatterTitle.trim()
  return content.match(/^#\s+(.+)$/m)?.[1]?.trim() || ''
}

async function saveContent() {
  if (!post.value || !editTitle.value.trim() || !editorContent.value.trim()) return
  saving.value = true
  actionError.value = ''
  try {
    const response = await api.put(`/admin/posts/${post.value.id}/content`, {
      title: editTitle.value,
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
  editTitle.value = ''
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
  <div class="fixed left-0 top-0 z-[60] h-0.5 bg-brand-500 transition-[width]" :style="{ width: `${readingProgress}%` }" aria-hidden="true"></div>
  <div v-if="!loaded" class="flex justify-center py-20">
    <div class="h-8 w-8 animate-spin rounded-full border-2 border-brand-500 border-t-transparent"></div>
  </div>

  <div v-else-if="!post" class="surface mx-auto max-w-3xl px-6 py-20 text-center">
    <h1 class="text-xl font-semibold text-slate-900 dark:text-white">文章不存在</h1>
    <p class="mt-2 text-sm text-slate-400">它可能已被删除，或者链接地址有误。</p>
    <router-link to="/posts" class="btn-primary mt-5"><ArrowLeft :size="16" />返回文章列表</router-link>
  </div>

  <div v-else>
    <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
      <router-link to="/posts" class="inline-flex items-center gap-1.5 text-sm text-slate-500 transition-colors hover:text-brand-600 dark:text-slate-400">
        <ArrowLeft :size="16" />
        返回文章列表
      </router-link>

      <div class="flex flex-wrap items-center gap-2">
        <input ref="markdownInput" type="file" accept=".md,text/markdown,text/plain" class="hidden" @change="overwriteFromMarkdown" />
        <button
          v-if="auth.isAuthenticated && !isEditing"
          type="button"
          title="上传 Markdown 覆盖当前文章"
          class="btn-secondary h-9 px-3"
          @click="chooseMarkdownFile"
        >
          <FileUp :size="16" />
          上传 MD
        </button>
        <button
          v-if="auth.isAuthenticated && !isEditing"
          type="button"
          :disabled="editorLoading"
          title="在线编辑"
          class="btn-secondary h-9 px-3"
          @click="startEditing"
        >
          <PencilLine :size="16" />
          {{ editorLoading ? '加载中' : '编辑' }}
        </button>
        <button
          type="button"
          :disabled="exporting"
          title="导出 Markdown"
          class="btn-primary h-9 px-3"
          @click="exportMarkdown"
        >
          <Download :size="16" />
          {{ exporting ? '导出中' : '导出' }}
        </button>
      </div>
    </div>

    <p v-if="actionError" class="mb-4 border border-red-200 bg-red-50 p-3 text-sm text-red-600 dark:border-red-900/50 dark:bg-red-950/30 dark:text-red-400">
      {{ actionError }}
    </p>

    <div class="grid items-start gap-4 lg:grid-cols-[minmax(0,1fr)_260px]">
      <article class="surface min-w-0 px-5 py-7 sm:px-8 md:px-12 md:py-10">
        <header class="mb-8 border-b border-slate-100 pb-8 dark:border-white/10">
          <img v-if="post.cover_image" :src="post.cover_image" :alt="post.title" class="mb-6 max-h-80 w-full object-cover" />
          <div v-if="isEditing">
            <label for="post-title-editor" class="mb-2 block text-sm font-semibold text-slate-700 dark:text-slate-200">文章标题</label>
            <input
              id="post-title-editor"
              v-model="editTitle"
              type="text"
              class="text-input py-3 text-2xl font-bold"
            />
          </div>
          <h1 v-else class="mb-4 text-3xl font-bold leading-tight text-slate-900 md:text-4xl dark:text-white">{{ post.title }}</h1>
          <div class="flex flex-wrap items-center gap-4 text-sm text-slate-400 dark:text-slate-500">
            <time class="inline-flex items-center gap-1.5"><CalendarDays :size="15" />{{ new Date(post.created_at).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' }) }}</time>
            <span class="inline-flex items-center gap-1.5"><Clock3 :size="15" />阅读约 {{ readingMinutes }} 分钟</span>
          </div>
          <div v-if="post.tags.length" class="mt-4 flex flex-wrap gap-2">
            <router-link v-for="tag in post.tags" :key="tag.id" :to="`/posts?tag=${tag.name}`" class="tag-pill">#{{ tag.name }}</router-link>
          </div>
        </header>

        <div v-if="isEditing">
          <div class="mb-3 flex flex-wrap items-center justify-between gap-2">
            <label for="post-markdown-editor" class="text-sm font-semibold text-slate-700 dark:text-slate-200">Markdown 原文</label>
            <button type="button" class="text-sm font-medium text-brand-600 hover:text-brand-700 dark:text-brand-400" @click="chooseMarkdownFile">重新选择 MD 文件</button>
          </div>
          <textarea
            id="post-markdown-editor"
            v-model="editorContent"
            spellcheck="false"
            class="min-h-[60vh] w-full resize-y border border-slate-200 bg-slate-50 p-4 font-mono text-sm leading-6 text-slate-800 outline-none focus:border-brand-500 dark:border-white/10 dark:bg-[#11151b] dark:text-slate-200"
          ></textarea>
          <div class="mt-4 flex justify-end gap-2">
            <button type="button" :disabled="saving" class="h-10 px-4 text-sm font-medium text-slate-500 transition-colors hover:bg-gray-100 disabled:opacity-50 dark:hover:bg-white/5" @click="cancelEditing">取消</button>
            <button type="button" :disabled="saving || !editTitle.trim() || !editorContent.trim()" class="h-10 bg-brand-500 px-5 text-sm font-semibold text-white transition-colors hover:bg-brand-600 disabled:cursor-not-allowed disabled:opacity-50" @click="saveContent">
              {{ saving ? '保存中...' : '保存文章' }}
            </button>
          </div>
        </div>

        <div v-else class="prose prose-lg max-w-none
          prose-headings:text-slate-800 dark:prose-headings:text-slate-100
          prose-a:text-brand-600 dark:prose-a:text-brand-400 prose-a:no-underline hover:prose-a:underline
          prose-code:text-brand-700 dark:prose-code:text-brand-300
          prose-pre:bg-slate-900 dark:prose-pre:bg-[#0B0F19]
          prose-blockquote:border-brand-500 dark:prose-blockquote:border-brand-400
          prose-img:mx-auto prose-img:max-w-full prose-img:shadow-md">
          <MarkdownRenderer :html="post.content_html" />
        </div>
      </article>

      <aside v-if="!isEditing" class="sticky top-[116px] hidden lg:block">
        <div class="surface max-h-[calc(100vh-8rem)] overflow-y-auto p-5">
          <TOCSidebar :toc-json="post.toc" />
        </div>
      </aside>
    </div>
  </div>
</template>
