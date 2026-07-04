<script setup lang="ts">
import { useTheme } from '../utils/theme'
const { isDark, toggle } = useTheme()
</script>

<template>
  <button @click="toggle" class="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700" :title="isDark ? 'Light mode' : 'Dark mode'">
    <span v-if="isDark">☀️</span>
    <span v-else>🌙</span>
  </button>
</template>
