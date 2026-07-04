<script setup lang="ts">
import { useTheme } from '../utils/theme'
const { isDark, toggle } = useTheme()
</script>

<template>
  <button @click="toggle"
    class="w-9 h-9 rounded-lg flex items-center justify-center transition-all hover:bg-warm-200/50 dark:hover:bg-white/10"
    :title="isDark ? 'Light mode' : 'Dark mode'">
    <svg v-if="isDark" class="w-5 h-5 text-brand-400" fill="currentColor" viewBox="0 0 24 24">
      <path d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"/>
    </svg>
    <svg v-else class="w-5 h-5 text-warm-500" fill="currentColor" viewBox="0 0 24 24">
      <path d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"/>
    </svg>
  </button>
</template>
