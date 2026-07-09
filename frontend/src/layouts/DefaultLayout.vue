<script setup lang="ts">
import { useRoute } from 'vue-router'
import ThemeToggle from '../components/ThemeToggle.vue'
import ThemeCustomizer from '../components/ThemeCustomizer.vue'

const route = useRoute()

const navLinks = [
  { to: '/posts', label: 'Posts' },
  { to: '/timeline', label: 'Timeline' },
  { to: '/search', label: 'Search' },
  { to: '/admin', label: 'Admin' },
]
</script>

<template>
  <div class="min-h-screen flex flex-col bg-gray-50 dark:bg-slate-950">
    <!-- Navbar -->
    <header class="sticky top-0 z-50 backdrop-blur-sm bg-white/80 dark:bg-slate-900/80 border-b border-gray-200 dark:border-white/5">
      <nav class="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between">
        <!-- Logo -->
        <router-link to="/" class="group flex items-center gap-2.5">
          <div class="w-9 h-9 rounded-xl bg-gradient-to-br from-brand-400 to-accent-500 flex items-center justify-center text-white font-black text-base shadow-md group-hover:shadow-lg group-hover:shadow-brand-500/25 transition-all">
            B
          </div>
          <span class="text-lg font-bold text-slate-800 dark:text-slate-100 group-hover:text-brand-500 transition-colors">
            My Blog
          </span>
        </router-link>

        <!-- Nav -->
        <div class="flex items-center gap-1">
          <router-link v-for="link in navLinks" :key="link.to" :to="link.to"
            class="px-3 py-2 rounded-lg text-sm font-medium transition-all"
            :class="route.path === link.to
              ? 'text-brand-600 dark:text-brand-400 bg-brand-50 dark:bg-brand-900/20'
              : 'text-slate-500 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200 hover:bg-gray-100 dark:hover:bg-white/5'"
          >
            {{ link.label }}
          </router-link>
          <ThemeToggle />
        </div>
      </nav>
    </header>

    <main class="flex-1 max-w-6xl w-full mx-auto px-4 py-10">
      <slot />
    </main>

    <footer class="border-t border-gray-200 dark:border-white/5">
      <div class="max-w-6xl mx-auto px-4 py-8">
        <div class="flex flex-col md:flex-row items-center justify-between gap-6">
          <div class="text-center md:text-left">
            <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-brand-400 to-accent-500 flex items-center justify-center text-white font-black text-sm mx-auto md:mx-0 mb-3">B</div>
            <p class="text-sm text-slate-400 dark:text-slate-500">Built with Vue + Go · Learning in public</p>
          </div>
          <ThemeCustomizer />
        </div>
      </div>
    </footer>
  </div>
</template>
