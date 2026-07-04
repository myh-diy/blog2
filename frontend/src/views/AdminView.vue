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

onMounted(loadPosts)

async function loadPosts() {
  loading.value = true
  try {
    const res = await api.get('/posts', { params: { per_page: 100 } })
    posts.value = res.data.posts
  } finally {
    loading.value = false
  }
}

async function deletePost(id: number) {
  if (!confirm('确定删除？')) return
  await api.delete(`/admin/posts/${id}`)
  posts.value = posts.value.filter((p) => p.id !== id)
}

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-8">
      <h1 class="text-3xl font-bold">Admin</h1>
      <button @click="logout" class="text-sm text-red-500 hover:underline">Logout</button>
    </div>

    <UploadZone @uploaded="loadPosts" />

    <h2 class="text-xl font-semibold mt-8 mb-4">All Posts ({{ posts.length }})</h2>
    <div v-if="loading" class="text-gray-500">Loading...</div>
    <table v-else class="w-full text-sm">
      <thead>
        <tr class="text-left border-b dark:border-gray-700">
          <th class="py-2">Title</th>
          <th class="py-2">Date</th>
          <th class="py-2">Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="post in posts" :key="post.id" class="border-b dark:border-gray-700/50">
          <td class="py-2">
            <router-link :to="`/post/${post.slug}`" class="hover:text-blue-500">{{ post.title }}</router-link>
          </td>
          <td class="py-2 text-gray-500">{{ post.created_at.split('T')[0] }}</td>
          <td class="py-2">
            <button @click="deletePost(post.id)" class="text-red-500 hover:underline">Delete</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
