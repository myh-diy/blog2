<script setup lang="ts">
import { useDisplayMode, type DisplayMode } from '../composables/useDisplayMode'

const { displayMode, setDisplayMode } = useDisplayMode()

const options: { value: DisplayMode; label: string }[] = [
  { value: 'minimal', label: '资讯版' },
  { value: 'classic', label: '经典版' },
]
</script>

<template>
  <div
    class="mode-switch"
    role="group"
    aria-label="页面版本"
  >
    <button
      v-for="option in options"
      :key="option.value"
      type="button"
      class="mode-switch__button"
      :class="displayMode === option.value
        ? 'mode-switch__button--active'
        : ''"
      :aria-pressed="displayMode === option.value"
      @click="setDisplayMode(option.value)"
    >
      {{ option.label }}
    </button>
  </div>
</template>
