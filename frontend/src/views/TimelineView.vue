<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../utils/api'

interface TimelinePost { id: number; title: string; slug: string; created_at: string }
interface TimelineMonth { month: string; posts: TimelinePost[] }
interface TimelineEntry { year: string; months: TimelineMonth[] }

const timeline = ref<TimelineEntry[]>([])
onMounted(async () => { const r = await api.get('/timeline'); timeline.value = r.data.timeline })
</script>

<template>
  <div>
    <h1 class="text-3xl font-black text-slate-800 dark:text-slate-100 mb-8">Timeline</h1>

    <div v-if="timeline.length" class="relative">
      <!-- Vertical line -->
      <div class="absolute left-[23px] top-4 bottom-4 w-0.5 bg-gradient-to-b from-brand-300 to-accent-300 dark:from-brand-700 dark:to-accent-700 rounded-full"></div>

      <div v-for="entry in timeline" :key="entry.year" class="mb-12 relative">
        <!-- Year node -->
        <div class="flex items-center gap-4 mb-6">
          <div class="w-12 h-12 rounded-full bg-brand-500 text-white flex items-center justify-center text-sm font-bold shadow-md z-10">
            {{ entry.year.slice(2) }}
          </div>
          <h2 class="text-xl font-bold text-slate-700 dark:text-slate-200">{{ entry.year }}</h2>
        </div>

        <div v-for="month in entry.months" :key="month.month" class="ml-16 mb-8">
          <h3 class="text-sm font-bold text-brand-600 dark:text-brand-400 mb-3">{{ month.month }}月</h3>
          <div class="space-y-3">
            <router-link v-for="post in month.posts" :key="post.id" :to="`/post/${post.slug}`"
              class="group block bg-white dark:bg-slate-900 rounded-xl border border-gray-100 dark:border-white/5 p-4
                     hover:border-brand-300 dark:hover:border-brand-700 hover:shadow-md hover:shadow-brand-500/5 hover:-translate-y-0.5 transition-all">
              <div class="flex items-baseline gap-3">
                <span class="text-xs font-mono text-slate-400 dark:text-slate-500 whitespace-nowrap">{{ post.created_at.split('T')[0].slice(5) }}</span>
                <span class="font-medium text-slate-700 dark:text-slate-200 group-hover:text-brand-600 dark:group-hover:text-brand-400 transition-colors">{{ post.title }}</span>
              </div>
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <p v-else class="text-center py-20 text-slate-400">No posts yet.</p>
  </div>
</template>
