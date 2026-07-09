<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import api from '../utils/api'
import UploadZone from '../components/UploadZone.vue'
import GradientButton from '../components/GradientButton.vue'
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
const editingPost = ref<Post | null>(null)
const editTags = ref('')

function startEdit(post: Post) {
  editingPost.value = post
  editTags.value = post.tags.map(t => t.name).join(', ')
}
async function saveEdit() {
  if (!editingPost.value) return
  const tags = editTags.value.split(',').map(t => t.trim()).filter(Boolean)
  await api.put(`/admin/posts/${editingPost.value.id}`, { tags })
  editingPost.value = null
  await loadPosts()
}
function cancelEdit() { editingPost.value = null }

function logout() { auth.logout(); router.push('/login') }
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-3xl font-black text-slate-800 dark:text-slate-100 mb-1">Admin</h1>
        <p class="text-sm text-slate-400 dark:text-slate-500">Manage your blog</p>
      </div>
      <button @click="logout" class="text-sm font-medium text-red-500 hover:text-red-600 transition-colors">Logout</button>
    </div>

    <!-- Upload -->
    <div class="bg-white dark:bg-slate-900 rounded-2xl border border-gray-100 dark:border-white/5 p-6 shadow-sm mb-10">
      <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-4">Upload Post</h2>
      <UploadZone @uploaded="loadPosts" />
    </div>

    <!-- Quotes -->
    <div class="bg-white dark:bg-slate-900 rounded-2xl border border-gray-100 dark:border-white/5 p-6 shadow-sm mb-10">
      <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-1">Floating Quotes <span class="ml-2 text-sm font-normal text-slate-400">({{ quotes.length }})</span></h2>
      <p class="text-sm text-slate-400 dark:text-slate-500 mb-4">Homepage floating sentences. Short 10-80 char phrases work best.</p>
      <form @submit.prevent="addQuote" class="flex gap-3 mb-4">
        <input v-model="newQuote" placeholder="Enter a quote..." maxlength="200"
          class="flex-1 px-4 py-2.5 bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-white/10 rounded-xl text-sm text-slate-800 dark:text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500/30 focus:border-brand-500 transition-all" />
        <GradientButton type="submit">Add</GradientButton>
      </form>
      <div v-if="quotes.length" class="rounded-xl border border-gray-100 dark:border-white/5 overflow-hidden">
        <div v-for="q in quotes" :key="q.id" class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100 dark:border-white/5 last:border-0 hover:bg-gray-50 dark:hover:bg-white/5 transition-colors">
          <span class="text-sm text-slate-700 dark:text-slate-300 flex-1 mr-4">{{ q.text }}</span>
          <button @click="deleteQuote(q.id)" class="text-xs font-medium text-red-500 hover:text-red-600 shrink-0 transition-colors">Delete</button>
        </div>
      </div>
      <p v-else class="text-sm text-slate-400 dark:text-slate-500 italic">No quotes yet. Falling back to sentences from blog posts.</p>
    </div>

    <!-- Posts -->
    <div class="bg-white dark:bg-slate-900 rounded-2xl border border-gray-100 dark:border-white/5 p-6 shadow-sm">
      <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-4">All Posts <span class="ml-2 text-sm font-normal text-slate-400">({{ posts.length }})</span></h2>
      <div v-if="loading" class="text-center py-12 text-slate-400">Loading...</div>
      <div v-else class="rounded-xl border border-gray-100 dark:border-white/5 overflow-hidden">
        <table class="w-full text-sm">
          <thead><tr class="border-b border-gray-100 dark:border-white/5 bg-gray-50 dark:bg-white/5">
            <th class="text-left py-3 px-4 font-semibold text-slate-500 dark:text-slate-400">Title</th>
            <th class="text-left py-3 px-4 font-semibold text-slate-500 dark:text-slate-400 w-28">Date</th>
            <th class="text-right py-3 px-4 font-semibold text-slate-500 dark:text-slate-400 w-32">Actions</th>
          </tr></thead>
          <tbody class="divide-y divide-gray-100 dark:divide-white/5">
            <template v-for="post in posts" :key="post.id">
              <tr class="hover:bg-gray-50 dark:hover:bg-white/5 transition-colors">
                <td class="py-3 px-4"><router-link :to="`/post/${post.slug}`" class="font-medium text-slate-800 dark:text-slate-200 hover:text-brand-600 dark:hover:text-brand-400 transition-colors">{{ post.title }}</router-link></td>
                <td class="py-3 px-4 text-slate-400 dark:text-slate-500 font-mono text-xs">{{ post.created_at.split('T')[0] }}</td>
                <td class="py-3 px-4 text-right space-x-2">
                  <button @click="startEdit(post)" class="text-xs font-medium text-brand-600 dark:text-brand-400 hover:underline transition-colors">Edit</button>
                  <button @click="deletePost(post.id)" class="text-xs font-medium text-red-500 hover:text-red-600 transition-colors">Delete</button>
                </td>
              </tr>
              <!-- Edit row -->
              <tr v-if="editingPost?.id === post.id">
                <td colspan="3" class="px-4 py-3 bg-gray-50 dark:bg-white/5">
                  <div class="flex flex-col gap-2">
                    <label class="text-xs font-semibold text-slate-500">Tags (comma separated)</label>
                    <input v-model="editTags" class="w-full px-3 py-2 bg-white dark:bg-slate-800 border border-gray-200 dark:border-white/10 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/30" />
                    <div class="flex gap-2">
                      <button @click="saveEdit" class="px-3 py-1.5 text-xs font-medium bg-gradient-to-r from-brand-400 to-accent-500 text-white rounded-lg hover:from-brand-500 hover:to-accent-600 transition-all">Save</button>
                      <button @click="cancelEdit" class="px-3 py-1.5 text-xs font-medium text-slate-500 hover:bg-gray-100 dark:hover:bg-white/10 rounded-lg transition-colors">Cancel</button>
                    </div>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
