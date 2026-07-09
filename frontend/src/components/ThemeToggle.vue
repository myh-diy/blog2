<script setup lang="ts">
import { useTheme } from '../utils/theme'
const { isDark, toggle } = useTheme()
</script>

<template>
  <button @click="toggle"
    class="w-9 h-9 rounded-lg flex items-center justify-center transition-all hover:bg-brand-100 dark:hover:bg-white/10"
    :title="isDark ? 'Light mode' : 'Dark mode'">
    <!-- Sun -->
    <svg v-if="isDark" class="w-5 h-5 text-brand-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
      <circle cx="12" cy="12" r="5"/>
      <path stroke-linecap="round" d="M12 1v2m0 18v2M4.22 4.22l1.42 1.42m12.72 12.72l1.42 1.42M1 12h2m18 0h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
    </svg>
    <!-- Moon -->
    <svg v-else class="w-5 h-5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"/>
    </svg>
  </button>
</template>
