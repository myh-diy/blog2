<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../utils/api'

const tags = ref<{ name: string; count: number }[]>([])

onMounted(async () => {
  const res = await api.get('/tags')
  tags.value = res.data.tags
})
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold mb-8">Tags</h1>
    <div class="flex flex-wrap gap-3">
      <router-link v-for="tag in tags" :key="tag.name" :to="`/?tag=${tag.name}`"
        class="px-4 py-2 bg-gray-100 dark:bg-gray-800 rounded-full hover:bg-blue-500 hover:text-white transition">
        {{ tag.name }}
        <span class="ml-1 text-sm opacity-60">{{ tag.count }}</span>
      </router-link>
    </div>
    <p v-if="tags.length === 0" class="text-gray-500 text-center py-12">No tags yet.</p>
  </div>
</template>
