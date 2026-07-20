<script setup lang="ts">
import { useSiteAvatar } from '../composables/useSiteAvatar'
const { siteAvatar, isDefaultAvatar } = useSiteAvatar()
</script>

<template>
  <div class="hidden xl:block fixed right-4 top-24 bottom-24 w-20 pointer-events-none z-0">
    <!-- Small floating site avatar -->
    <div class="absolute -top-1 -right-1 w-16 h-16 rounded-full overflow-hidden border-4 border-white dark:border-slate-800 shadow-lg bg-gray-100 dark:bg-slate-800 animate-bounce" style="animation-duration: 3.5s;">
      <img :src="siteAvatar" alt="site avatar" class="h-full w-full" :class="isDefaultAvatar ? 'object-contain p-1' : 'object-cover'" />
    </div>

    <!-- Stars -->
    <div class="absolute top-28 right-2 text-brand-300 animate-pulse">
      <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2l2.4 7.2h7.6l-6 4.8 2.4 7.2-6-4.8-6 4.8 2.4-7.2-6-4.8h7.6z"/></svg>
    </div>
    <div class="absolute top-44 right-8 text-accent-300 animate-pulse" style="animation-delay: 0.5s;">
      <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2l2.4 7.2h7.6l-6 4.8 2.4 7.2-6-4.8-6 4.8 2.4-7.2-6-4.8h7.6z"/></svg>
    </div>
    <div class="absolute top-60 right-1 text-brand-300 animate-pulse" style="animation-delay: 1s;">
      <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2l2.4 7.2h7.6l-6 4.8 2.4 7.2-6-4.8-6 4.8 2.4-7.2-6-4.8h7.6z"/></svg>
    </div>

    <!-- Hearts -->
    <div class="absolute top-80 right-6 text-accent-400 animate-bounce" style="animation-duration: 2.5s;">
      <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/></svg>
    </div>
    <div class="absolute top-96 right-2 text-brand-400 animate-bounce" style="animation-duration: 3.5s;">
      <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/></svg>
    </div>

    <!-- Small planet/bubble -->
    <div class="absolute bottom-40 right-4 w-10 h-10 rounded-full bg-brand-100 dark:bg-brand-900/30 flex items-center justify-center animate-pulse" style="animation-delay: 0.3s;">
      <div class="w-3 h-3 rounded-full bg-brand-300"></div>
    </div>
    <div class="absolute bottom-24 right-7 w-6 h-6 rounded-full bg-accent-100 dark:bg-accent-900/30 animate-pulse" style="animation-delay: 0.8s;"></div>

    <!-- Vertical dotted line -->
    <div class="absolute top-16 bottom-16 right-9 w-px border-r-2 border-dashed border-brand-200 dark:border-brand-800 opacity-50"></div>
  </div>
</template>
