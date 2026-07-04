<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { usePostsStore } from '../stores/posts'
import PostCard from '../components/PostCard.vue'

const route = useRoute()
const store = usePostsStore()

onMounted(() => load())
watch(() => route.query, load)

function load() {
  const page = Number(route.query.page) || 1
  const tag = (route.query.tag as string) || ''
  store.fetchPosts(page, tag)
}
</script>

<template>
  <div>
    <h1 v-if="!route.query.tag" class="text-3xl font-bold mb-8">Blog Posts</h1>
    <h1 v-else class="text-3xl font-bold mb-8">
      Tag: {{ route.query.tag }}
    </h1>

    <div v-if="store.loading" class="text-center py-12 text-gray-500">Loading...</div>

    <div v-else-if="store.posts.length === 0" class="text-center py-12 text-gray-500">
      No posts yet.
    </div>

    <template v-else>
      <PostCard v-for="post in store.posts" :key="post.id" :post="post" />

      <div class="flex justify-center gap-2 mt-8">
        <button
          v-for="p in Math.ceil(store.total / 10)"
          :key="p"
          :class="['px-3 py-1 rounded', p === (Number(route.query.page) || 1) ? 'bg-blue-500 text-white' : 'bg-gray-100 dark:bg-gray-800']"
          @click="$router.push({ query: { ...$route.query, page: p } })"
        >
          {{ p }}
        </button>
      </div>
    </template>
  </div>
</template>
