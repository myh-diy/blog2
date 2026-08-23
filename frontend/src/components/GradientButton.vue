<script setup lang="ts">
withDefaults(defineProps<{
  type?: 'button' | 'submit' | 'reset'
  disabled?: boolean
}>(), {
  type: 'button',
  disabled: false,
})
</script>

<template>
  <button :type="type" :disabled="disabled"
    class="inline-flex items-center justify-center gap-2 bg-brand-500 px-5 py-2.5 text-sm font-semibold text-white
           transition-colors hover:bg-brand-600 active:bg-brand-700 disabled:cursor-not-allowed disabled:opacity-60">
    <slot />
  </button>
</template>
