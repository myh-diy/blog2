<script setup lang="ts">
import { useTheme } from '../utils/theme'
import { Moon, Sun } from '@lucide/vue'
const { isDark, toggle } = useTheme()
</script>

<template>
  <button class="icon-button" type="button" @click="toggle" :title="isDark ? '切换浅色模式' : '切换深色模式'">
    <Sun v-if="isDark" :size="18" />
    <Moon v-else :size="18" />
  </button>
</template>
