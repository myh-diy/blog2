<script setup lang="ts">
import { useRoute } from 'vue-router'
import ThemeToggle from '../components/ThemeToggle.vue'

const route = useRoute()

const navLinks = [
  { to: '/posts', label: 'Posts' },
  { to: '/timeline', label: 'Timeline' },
  { to: '/search', label: 'Search' },
  { to: '/admin', label: 'Admin' },
]
</script>

<template>
  <div class="min-h-screen bg-warm-50 dark:bg-[#0d0714]">
    <!-- Glass navbar -->
    <header class="sticky top-0 z-50 backdrop-blur-xl bg-warm-50/70 dark:bg-[#0d0714]/70 border-b border-warm-200/40 dark:border-white/5">
      <nav class="max-w-4xl mx-auto px-4 h-16 flex items-center justify-between">
        <!-- Logo -->
        <router-link to="/" class="group flex items-center gap-2.5">
          <div class="w-9 h-9 rounded-xl bg-gradient-to-br from-brand-400 to-brand-600 dark:from-pop-400 dark:to-pop-600 flex items-center justify-center text-white font-black text-base shadow-md group-hover:shadow-lg transition-shadow">
            B
          </div>
          <span class="text-lg font-bold text-warm-800 dark:text-warm-100 group-hover:text-brand-600 dark:group-hover:text-pop-400 transition-colors">
            My Blog
          </span>
        </router-link>

        <!-- Nav -->
        <div class="flex items-center gap-0.5">
          <router-link v-for="link in navLinks" :key="link.to" :to="link.to"
            class="px-3 py-2 rounded-lg text-sm font-medium transition-all"
            :class="route.path === link.to
              ? 'text-brand-600 dark:text-pop-400 bg-brand-50 dark:bg-pop-900/20'
              : 'text-warm-500 dark:text-warm-400 hover:text-warm-800 dark:hover:text-warm-200 hover:bg-warm-100 dark:hover:bg-white/5'"
          >
            {{ link.label }}
          </router-link>
          <ThemeToggle />
        </div>
      </nav>
    </header>

    <main class="max-w-4xl mx-auto px-4 py-10">
      <slot />
    </main>

    <footer class="border-t border-warm-200/50 dark:border-white/5 mt-20">
      <div class="max-w-4xl mx-auto px-4 py-8 text-center text-sm text-warm-400 dark:text-warm-500">
        Built with Vue + Go · Learning in public
      </div>
    </footer>
  </div>
</template>
