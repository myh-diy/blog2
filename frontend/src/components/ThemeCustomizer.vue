<script setup lang="ts">
import { useCustomTheme } from '../composables/useCustomTheme'

const { colors, setBrand, setAccent, reset } = useCustomTheme()
</script>

<template>
  <div class="bg-white dark:bg-slate-900 rounded-2xl border border-gray-100 dark:border-white/5 p-4 shadow-sm w-64">
    <h3 class="text-sm font-bold text-slate-700 dark:text-slate-200 mb-3">Theme Colors</h3>

    <div class="space-y-3">
      <div>
        <label class="block text-xs font-medium text-slate-500 mb-1">Primary</label>
        <div class="flex items-center gap-2">
          <input type="color" :value="colors.brand" @input="e => setBrand((e.target as HTMLInputElement).value)"
            class="w-8 h-8 rounded-lg border border-gray-200 dark:border-white/10 cursor-pointer bg-transparent" />
          <span class="text-xs font-mono text-slate-400">{{ colors.brand }}</span>
        </div>
      </div>

      <div>
        <label class="block text-xs font-medium text-slate-500 mb-1">Accent</label>
        <div class="flex items-center gap-2">
          <input type="color" :value="colors.accent" @input="e => setAccent((e.target as HTMLInputElement).value)"
            class="w-8 h-8 rounded-lg border border-gray-200 dark:border-white/10 cursor-pointer bg-transparent" />
          <span class="text-xs font-mono text-slate-400">{{ colors.accent }}</span>
        </div>
      </div>

      <button @click="reset"
        class="w-full text-xs font-medium text-slate-500 hover:text-brand-600 dark:hover:text-brand-400 py-1.5 rounded-lg hover:bg-gray-50 dark:hover:bg-white/5 transition-colors">
        Reset to default
      </button>
    </div>
  </div>
</template>
