import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../utils/api'

export interface Post {
  id: number
  title: string
  slug: string
  content_html: string
  toc: string
  created_at: string
  updated_at: string
  tags: { id: number; name: string }[]
}

export const usePostsStore = defineStore('posts', () => {
  const posts = ref<Post[]>([])
  const total = ref(0)
  const loading = ref(false)

  async function fetchPosts(page = 1, tag = '') {
    loading.value = true
    try {
      const res = await api.get('/posts', { params: { page, per_page: 10, tag } })
      posts.value = res.data.posts
      total.value = res.data.total
    } finally {
      loading.value = false
    }
  }

  async function fetchPost(slug: string): Promise<Post | null> {
    try {
      const res = await api.get(`/posts/${slug}`)
      return res.data.post
    } catch {
      return null
    }
  }

  return { posts, total, loading, fetchPosts, fetchPost }
})
