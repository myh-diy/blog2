<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../utils/api'
const tags = ref<{ name: string; count: number }[]>([])
onMounted(async () => { const r = await api.get('/tags'); tags.value = r.data.tags })
</script>

<template>
  <div>
    <h1 class="text-4xl font-bold text-warm-800 dark:text-warm-100 mb-1">Tags</h1>
    <p class="text-warm-400 dark:text-warm-500 mb-8 text-sm">{{ tags.length }} topics</p>
    <div v-if="tags.length" class="flex flex-wrap gap-3">
      <router-link v-for="t in tags" :key="t.name" :to="`/?tag=${t.name}`"
        class="group inline-flex items-center gap-2 px-5 py-3 rounded-2xl
               bg-white dark:bg-white/5 border border-warm-200 dark:border-white/5
               hover:border-brand-300 dark:hover:border-pop-600
               hover:shadow-lg hover:shadow-brand-500/5 dark:hover:shadow-pop-500/5 transition-all">
        <span class="font-medium text-warm-700 dark:text-warm-200 group-hover:text-brand-600 dark:group-hover:text-pop-400 transition-colors">{{ t.name }}</span>
        <span class="text-xs px-2 py-0.5 rounded-full bg-warm-100 dark:bg-white/10 text-warm-500 dark:text-warm-400">{{ t.count }}</span>
      </router-link>
    </div>
    <p v-else class="text-center py-20 text-warm-400">No tags yet.</p>
  </div>
</template>
