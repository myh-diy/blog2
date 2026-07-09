<script setup lang="ts">
import KawaiiIcon from './KawaiiIcon.vue'

defineProps<{
  icon: 'search' | 'sad' | 'happy'
  title: string
  description?: string
}>()
</script>

<template>
  <div class="flex flex-col items-center justify-center py-16 text-center">
    <div class="w-28 h-28 mb-5">
      <KawaiiIcon :name="icon" />
    </div>
    <h3 class="text-lg font-bold text-slate-700 dark:text-slate-200">{{ title }}</h3>
    <p v-if="description" class="mt-1.5 text-sm text-slate-400 dark:text-slate-500 max-w-xs">{{ description }}</p>
    <slot />
  </div>
</template>
