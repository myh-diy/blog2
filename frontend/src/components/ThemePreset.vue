<script setup lang="ts">
import { useCustomTheme } from '../composables/useCustomTheme'

const { presets, activePreset, setPreset } = useCustomTheme()
</script>

<template>
  <div class="bg-white dark:bg-slate-900 rounded-2xl border border-gray-100 dark:border-white/5 p-4 shadow-sm">
    <h3 class="text-sm font-bold text-slate-700 dark:text-slate-200 mb-3">Theme</h3>
    <div class="flex gap-2">
      <button v-for="p in presets" :key="p.name" @click="setPreset(p.name)"
        class="relative w-8 h-8 rounded-full border-2 transition-all"
        :class="activePreset === p.name
          ? 'border-slate-400 dark:border-slate-300 scale-110 shadow-md'
          : 'border-transparent hover:scale-105'"
        :style="{ backgroundColor: p.brand }"
        :title="p.name">
        <span class="absolute bottom-0 right-0 w-3 h-3 rounded-full border border-white dark:border-slate-900" :style="{ backgroundColor: p.accent }"></span>
      </button>
    </div>
    <p class="text-xs text-slate-400 mt-2">{{ activePreset }}</p>
  </div>
</template>
