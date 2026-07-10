<script setup lang="ts">
import { useRoute } from 'vue-router'
import ThemeToggle from '../components/ThemeToggle.vue'
import PageDecorations from '../components/PageDecorations.vue'
import SillyTavernLink from '../components/SillyTavernLink.vue'
import { useBackgroundImage } from '../composables/useBackgroundImage'
import { useSiteTitle } from '../composables/useSiteTitle'

const route = useRoute()
const { backgroundImage, bgOpacity } = useBackgroundImage()
const { siteTitle } = useSiteTitle()

const navLinks = [
  { to: '/posts', label: 'Posts' },
  { to: '/timeline', label: 'Timeline' },
  { to: '/search', label: 'Search' },
  { to: '/admin', label: 'Admin' },
]
</script>

<template>
  <div class="relative min-h-screen flex flex-col">
    <!-- Global background image layer -->
    <div
      class="fixed inset-0 z-0 bg-cover bg-center bg-fixed bg-no-repeat transition-all duration-700"
      :style="{ backgroundImage: backgroundImage ? `url('${backgroundImage}')` : 'none' }"
      aria-hidden="true"
    >
      <div
        class="absolute inset-0 transition-colors"
        :class="backgroundImage ? 'bg-white dark:bg-slate-950' : 'bg-gray-50 dark:bg-slate-950'"
        :style="backgroundImage ? { opacity: bgOpacity } : undefined"
      ></div>
    </div>
    <!-- Navbar -->
    <header class="relative z-10 sticky top-0 z-50 backdrop-blur-md bg-white/60 dark:bg-slate-900/60 border-b border-gray-200/50 dark:border-white/5">
      <nav class="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between">
        <!-- Logo -->
        <router-link to="/" class="px-3 py-2 rounded-lg text-sm font-medium transition-all text-slate-500 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200 hover:bg-gray-100 dark:hover:bg-white/5">
          {{ siteTitle }}
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

    <main class="relative z-10 flex-1 max-w-6xl w-full mx-auto px-4 py-10">
      <PageDecorations />
      <slot />
    </main>

    <SillyTavernLink />

    <footer class="relative z-10 border-t border-gray-200/50 dark:border-white/5">
      <div class="max-w-6xl mx-auto px-4 py-6 text-center text-xs text-slate-400 dark:text-slate-500">
        {{ siteTitle }}
      </div>
    </footer>
  </div>
</template>
