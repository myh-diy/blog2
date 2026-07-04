<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../utils/api'

const timeline = ref<any[]>([])

onMounted(async () => {
  const res = await api.get('/timeline')
  timeline.value = res.data.timeline
})
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold mb-8">Timeline</h1>
    <div v-for="entry in timeline" :key="entry.year" class="mb-8">
      <h2 class="text-2xl font-bold text-blue-500 mb-4">{{ entry.year }}</h2>
      <div v-for="month in entry.months" :key="month.month" class="ml-4 mb-4">
        <h3 class="text-lg font-semibold text-gray-500 mb-2">{{ month.month }}月</h3>
        <ul class="space-y-2">
          <li v-for="post in month.posts" :key="post.id">
            <router-link :to="`/post/${post.slug}`" class="hover:text-blue-500">
              <span class="text-sm text-gray-400 mr-2">{{ post.created_at.split('T')[0] }}</span>
              {{ post.title }}
            </router-link>
          </li>
        </ul>
      </div>
    </div>
    <p v-if="timeline.length === 0" class="text-gray-500 text-center py-12">No posts yet.</p>
  </div>
</template>
