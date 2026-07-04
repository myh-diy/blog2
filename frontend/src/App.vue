<script setup lang="ts">
import { useTheme } from './utils/theme'
import DefaultLayout from './layouts/DefaultLayout.vue'

useTheme() // Initialize theme on app load
</script>

<template>
  <DefaultLayout>
    <router-view />
  </DefaultLayout>
</template>
