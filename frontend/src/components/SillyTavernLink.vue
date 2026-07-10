<script setup lang="ts">
const sillyUrl = `${window.location.protocol}//${window.location.hostname}:8000`
</script>

<template>
  <a :href="sillyUrl" target="_blank" rel="noopener noreferrer"
    class="fixed left-4 bottom-4 z-50 w-12 h-12 rounded-full bg-slate-800 dark:bg-slate-700 text-white shadow-lg hover:shadow-xl hover:scale-110 transition-all flex items-center justify-center"
    title="SillyTavern">
    <svg class="w-7 h-7" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <!-- Comedy/tragedy mask style icon -->
      <path d="M4 10c0-3 2.5-6 6-6s6 3 6 6v4c0 3-2.5 6-6 6s-6-3-6-6v-4z"/>
      <path d="M9 10h.01"/>
      <path d="M15 10h.01"/>
      <path d="M10 14c.5 1 1.5 1.5 2 1.5s1.5-.5 2-1.5"/>
      <path d="M12 2v2"/>
    </svg>
  </a>
</template>
