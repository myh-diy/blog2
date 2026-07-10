<script setup lang="ts">
import { useTheme } from './utils/theme'
import { useCustomTheme } from './composables/useCustomTheme'
import DefaultLayout from './layouts/DefaultLayout.vue'

useTheme() // Initialize dark/light theme on app load
useCustomTheme() // Initialize color preset on app load
</script>

<template>
  <DefaultLayout>
    <router-view />
  </DefaultLayout>
</template>
