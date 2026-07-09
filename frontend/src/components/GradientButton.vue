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
    class="inline-flex items-center justify-center gap-2 px-5 py-2.5 rounded-xl font-semibold text-white
           bg-brand-500 hover:bg-brand-600 active:bg-brand-700
           hover:-translate-y-0.5 hover:shadow-lg
           active:scale-[0.98] active:translate-y-0
           transition-all duration-200 disabled:opacity-60 disabled:cursor-not-allowed disabled:hover:translate-y-0">
    <slot />
  </button>
</template>
