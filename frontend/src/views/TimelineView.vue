<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../utils/api'
const timeline = ref<any[]>([])
onMounted(async () => { const r = await api.get('/timeline'); timeline.value = r.data.timeline })
</script>

<template>
  <div>
    <h1 class="text-4xl font-bold text-warm-800 dark:text-warm-100 mb-8">Timeline</h1>
    <div v-if="timeline.length" class="relative">
      <div class="absolute left-[19px] top-3 bottom-3 w-px bg-warm-200 dark:bg-white/5"></div>
      <div v-for="entry in timeline" :key="entry.year" class="mb-10">
        <div class="flex items-center gap-3 mb-6">
          <div class="w-10 h-10 rounded-full bg-brand-500 dark:bg-pop-500 text-white flex items-center justify-center text-sm font-bold shadow-md z-10">{{ entry.year }}</div>
        </div>
        <div v-for="month in entry.months" :key="month.month" class="ml-14 mb-6">
          <h3 class="text-xs font-semibold text-warm-400 dark:text-warm-500 mb-3 uppercase tracking-wide">{{ month.month }}月</h3>
          <ul class="space-y-2">
            <li v-for="post in month.posts" :key="post.id" class="group">
              <router-link :to="`/post/${post.slug}`" class="block p-3 -ml-2 rounded-xl hover:bg-white dark:hover:bg-white/5 transition-colors">
                <div class="flex items-baseline gap-3">
                  <span class="text-xs text-warm-400 dark:text-warm-600 font-mono whitespace-nowrap">{{ post.created_at.split('T')[0].slice(5) }}</span>
                  <span class="font-medium text-warm-700 dark:text-warm-200 group-hover:text-brand-600 dark:group-hover:text-pop-400 transition-colors">{{ post.title }}</span>
                </div>
              </router-link>
            </li>
          </ul>
        </div>
      </div>
    </div>
    <p v-else class="text-center py-20 text-warm-400">No posts yet.</p>
  </div>
</template>
