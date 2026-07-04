<script setup lang="ts">
import ThemeToggle from '../components/ThemeToggle.vue'
</script>

<template>
  <div class="min-h-screen bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100">
    <header class="border-b border-gray-200 dark:border-gray-700">
      <nav class="max-w-4xl mx-auto px-4 h-14 flex items-center justify-between">
        <router-link to="/" class="text-xl font-bold">My Blog</router-link>
        <div class="flex items-center gap-4 text-sm">
          <router-link to="/timeline" class="hover:text-blue-500">Timeline</router-link>
          <router-link to="/tags" class="hover:text-blue-500">Tags</router-link>
          <router-link to="/search" class="hover:text-blue-500">Search</router-link>
          <router-link to="/admin" class="hover:text-blue-500">Admin</router-link>
          <ThemeToggle />
        </div>
      </nav>
    </header>
    <main class="max-w-4xl mx-auto px-4 py-8">
      <slot />
    </main>
  </div>
</template>
