<script setup lang="ts">
import { onMounted, watch, ref } from 'vue'
import hljs from 'highlight.js/lib/core'
import 'highlight.js/styles/github-dark.css'

const props = defineProps<{ html: string }>()
const el = ref<HTMLElement>()

function highlight() {
  if (el.value) {
    el.value.querySelectorAll('pre code').forEach((block) => {
      hljs.highlightElement(block as HTMLElement)
    })
  }
}

onMounted(highlight)
watch(() => props.html, () => {
  setTimeout(highlight, 0)
})
</script>

<template>
  <div ref="el" class="prose dark:prose-invert max-w-none" v-html="html" />
</template>
