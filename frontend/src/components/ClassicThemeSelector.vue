<script setup lang="ts">
import { Image, Palette } from '@lucide/vue'
import { useClassicTheme } from '../composables/useClassicTheme'

const { classicThemes, activeThemeId, setClassicTheme } = useClassicTheme()
</script>

<template>
  <section class="relative z-10 border-b border-slate-200 bg-white/90 backdrop-blur dark:border-white/10 dark:bg-[#171b22]/90" aria-label="经典版主题">
    <div class="nav-scroll mx-auto flex min-h-14 max-w-[1200px] items-center gap-2 overflow-x-auto px-4 py-2">
      <span class="mr-1 hidden shrink-0 items-center gap-2 text-xs font-semibold text-slate-500 sm:inline-flex dark:text-slate-400">
        <Palette :size="15" />经典主题
      </span>
      <button
        v-for="theme in classicThemes"
        :key="theme.id"
        type="button"
        class="flex h-10 shrink-0 items-center gap-2 border px-2.5 text-left transition-colors"
        :class="activeThemeId === theme.id
          ? 'border-brand-400 bg-brand-50 text-brand-700 dark:border-brand-600 dark:bg-brand-900/20 dark:text-brand-300'
          : 'border-slate-200 bg-white text-slate-500 hover:border-slate-300 dark:border-white/10 dark:bg-[#11151b] dark:text-slate-400'"
        :title="`${theme.name}：${theme.description}`"
        @click="setClassicTheme(theme.id)"
      >
        <img v-if="theme.background" :src="theme.background" alt="" class="h-7 w-9 object-cover" :style="{ objectPosition: theme.position }" />
        <span v-else class="flex h-7 w-9 items-center justify-center bg-slate-100 dark:bg-white/5"><Image :size="15" /></span>
        <span>
          <span class="block text-xs font-semibold">{{ theme.name }}</span>
          <span class="hidden text-[10px] opacity-70 md:block">{{ theme.description }}</span>
        </span>
      </button>
    </div>
  </section>
</template>
