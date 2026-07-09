<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../utils/api'
import EmptyState from '../components/EmptyState.vue'

const tags = ref<{ name: string; count: number }[]>([])
onMounted(async () => { const r = await api.get('/tags'); tags.value = r.data.tags })
</script>

<template>
  <div>
    <h1 class="text-3xl font-black text-slate-800 dark:text-slate-100 mb-1">Tags</h1>
    <p class="text-slate-400 dark:text-slate-500 mb-8 text-sm">{{ tags.length }} topics to explore</p>

    <div v-if="tags.length" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
      <router-link v-for="t in tags" :key="t.name" :to="`/posts?tag=${t.name}`"
        class="group flex flex-col items-center justify-center p-6 rounded-2xl
               bg-white dark:bg-slate-900 border border-gray-100 dark:border-white/5
               hover:border-brand-300 dark:hover:border-brand-700
               hover:shadow-lg hover:shadow-brand-500/10 hover:-translate-y-1
               transition-all duration-300">
        <div class="w-10 h-10 mb-3 rounded-full bg-gradient-to-br from-brand-100 to-accent-100 dark:from-brand-900/30 dark:to-accent-900/30 flex items-center justify-center text-lg">🏷️</div>
        <span class="font-bold text-slate-700 dark:text-slate-200 group-hover:text-brand-600 dark:group-hover:text-brand-400 transition-colors">{{ t.name }}</span>
        <span class="text-xs text-slate-400 mt-1">{{ t.count }} posts</span>
      </router-link>
    </div>

    <EmptyState v-else icon="sad" title="No tags yet" description="Upload posts with tags to see them here." />
  </div>
</template>
