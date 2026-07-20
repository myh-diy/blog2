<script setup lang="ts">
import { computed } from 'vue'

interface TOCItem {
  id: string
  text: string
  level: number
  children: TOCItem[]
}

const props = defineProps<{ tocJson: string }>()

const items = computed<TOCItem[]>(() => {
  try {
    const tree = JSON.parse(props.tocJson) || []
    const result: TOCItem[] = []

    function walk(nodes: TOCItem[]) {
      for (const node of nodes) {
        if (node.id && node.text) result.push(node)
        if (node.children?.length) walk(node.children)
      }
    }

    walk(tree)
    return result
  } catch {
    return []
  }
})
</script>

<template>
  <nav v-if="items.length" aria-label="文章目录">
    <h3 class="mb-4 text-xs font-bold uppercase text-slate-400 dark:text-slate-500">目录</h3>
    <ul class="border-l border-gray-200 text-sm dark:border-white/10">
      <li v-for="item in items" :key="item.id">
        <a
          :href="`#${item.id}`"
          :style="{ paddingLeft: `${Math.max(12, (item.level - 1) * 12)}px` }"
          class="-ml-px block border-l-2 border-transparent py-1.5 leading-5 text-slate-500 transition-colors hover:border-brand-500 hover:text-brand-600 dark:text-slate-400 dark:hover:border-brand-400 dark:hover:text-brand-400"
        >
          {{ item.text }}
        </a>
      </li>
    </ul>
  </nav>
</template>
