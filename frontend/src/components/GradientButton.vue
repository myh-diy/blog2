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
           bg-gradient-to-r from-brand-400 to-accent-500
           hover:from-brand-500 hover:to-accent-600
           hover:-translate-y-0.5 hover:shadow-lg hover:shadow-brand-500/25
           active:scale-[0.98] active:translate-y-0
           transition-all duration-200 disabled:opacity-60 disabled:cursor-not-allowed disabled:hover:translate-y-0">
    <slot />
  </button>
</template>
