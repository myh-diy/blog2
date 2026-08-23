<script setup lang="ts">
import { computed } from 'vue'
import { useDisplayMode } from '../composables/useDisplayMode'
import ClassicHome from './home/ClassicHome.vue'
import MinimalHome from './home/MinimalHome.vue'

const { displayMode } = useDisplayMode()
const activeHome = computed(() => displayMode.value === 'minimal' ? MinimalHome : ClassicHome)
</script>

<template>
  <component :is="activeHome" />
</template>
