<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import api from '../utils/api'
import UploadZone from '../components/UploadZone.vue'
import type { Post } from '../stores/posts'

const router = useRouter()
const auth = useAuthStore()
const posts = ref<Post[]>([])
const loading = ref(true)
const quotes = ref<{ id: number; text: string; created_at: string }[]>([])
const newQuote = ref('')

onMounted(() => { loadPosts(); loadQuotes() })

async function loadPosts() {
  loading.value = true
  try { const r = await api.get('/posts', { params: { per_page: 100 } }); posts.value = r.data.posts }
  finally { loading.value = false }
}
async function loadQuotes() {
  try { const r = await api.get('/admin/quotes'); quotes.value = r.data.quotes }
  catch { quotes.value = [] }
}
async function addQuote() {
  if (!newQuote.value.trim()) return
  await api.post('/admin/quotes', { text: newQuote.value.trim() })
  newQuote.value = ''; await loadQuotes()
}
async function deleteQuote(id: number) {
  await api.delete(`/admin/quotes/${id}`)
  quotes.value = quotes.value.filter(q => q.id !== id)
}
async function deletePost(id: number) {
  if (!confirm('Delete this post?')) return
  await api.delete(`/admin/posts/${id}`)
  posts.value = posts.value.filter(p => p.id !== id)
}
function logout() { auth.logout(); router.push('/login') }
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-4xl font-bold text-warm-800 dark:text-warm-100 mb-1">Admin</h1>
        <p class="text-sm text-warm-400 dark:text-warm-500">Manage your blog</p>
      </div>
      <button @click="logout" class="text-sm font-medium text-red-500 hover:text-red-600 transition-colors">Logout</button>
    </div>

    <UploadZone @uploaded="loadPosts" />

    <!-- Quotes -->
    <div class="mt-10">
      <h2 class="text-xl font-bold text-warm-800 dark:text-warm-100 mb-1">Floating Quotes <span class="ml-2 text-sm font-normal text-warm-400">({{ quotes.length }})</span></h2>
      <p class="text-sm text-warm-400 dark:text-warm-500 mb-4">Homepage floating sentences. Short 10-80 char phrases work best.</p>
      <form @submit.prevent="addQuote" class="flex gap-3 mb-4">
        <input v-model="newQuote" placeholder="Enter a quote..." maxlength="200"
          class="flex-1 px-4 py-2.5 bg-white dark:bg-white/5 border border-warm-200 dark:border-white/10 rounded-xl text-sm text-warm-800 dark:text-warm-100 placeholder-warm-400 focus:outline-none focus:ring-2 focus:ring-brand-500/30 dark:focus:ring-pop-500/30 focus:border-brand-500 dark:focus:border-pop-500 transition-all" />
        <button type="submit" class="btn-primary text-sm">Add</button>
      </form>
      <div v-if="quotes.length" class="bg-white dark:bg-white/5 rounded-xl border border-warm-200 dark:border-white/5 overflow-hidden">
        <div v-for="q in quotes" :key="q.id" class="flex items-center justify-between px-4 py-2.5 border-b border-warm-100 dark:border-white/5 last:border-0 hover:bg-warm-50 dark:hover:bg-white/5 transition-colors">
          <span class="text-sm text-warm-700 dark:text-warm-300 flex-1 mr-4">{{ q.text }}</span>
          <button @click="deleteQuote(q.id)" class="text-xs font-medium text-red-500 hover:text-red-600 shrink-0 transition-colors">Delete</button>
        </div>
      </div>
      <p v-else class="text-sm text-warm-400 dark:text-warm-500 italic">No quotes yet. Falling back to sentences from blog posts.</p>
    </div>

    <!-- Posts -->
    <div class="mt-10">
      <h2 class="text-xl font-bold text-warm-800 dark:text-warm-100 mb-4">All Posts <span class="ml-2 text-sm font-normal text-warm-400">({{ posts.length }})</span></h2>
      <div v-if="loading" class="text-center py-12 text-warm-400">Loading...</div>
      <div v-else class="bg-white dark:bg-white/5 rounded-xl border border-warm-200 dark:border-white/5 overflow-hidden">
        <table class="w-full text-sm">
          <thead><tr class="border-b border-warm-200 dark:border-white/5 bg-warm-50 dark:bg-white/5">
            <th class="text-left py-3 px-4 font-semibold text-warm-500 dark:text-warm-400">Title</th>
            <th class="text-left py-3 px-4 font-semibold text-warm-500 dark:text-warm-400 w-28">Date</th>
            <th class="text-right py-3 px-4 font-semibold text-warm-500 dark:text-warm-400 w-20">Actions</th>
          </tr></thead>
          <tbody class="divide-y divide-warm-100 dark:divide-white/5">
            <tr v-for="post in posts" :key="post.id" class="hover:bg-warm-50 dark:hover:bg-white/5 transition-colors">
              <td class="py-3 px-4"><router-link :to="`/post/${post.slug}`" class="font-medium text-warm-800 dark:text-warm-200 hover:text-brand-600 dark:hover:text-pop-400 transition-colors">{{ post.title }}</router-link></td>
              <td class="py-3 px-4 text-warm-400 dark:text-warm-500 font-mono text-xs">{{ post.created_at.split('T')[0] }}</td>
              <td class="py-3 px-4 text-right"><button @click="deletePost(post.id)" class="text-xs font-medium text-red-500 hover:text-red-600 transition-colors">Delete</button></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
