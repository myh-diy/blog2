<script setup lang="ts">
import { computed } from 'vue'
interface TOCItem { id: string; text: string; level: number; children: TOCItem[] }
const props = defineProps<{ tocJson: string }>()
const items = computed<TOCItem[]>(() => { try { return JSON.parse(props.tocJson) || [] } catch { return [] } })
</script>

<template>
  <nav v-if="items.length">
    <h3 class="text-xs font-semibold text-warm-400 dark:text-warm-500 uppercase tracking-wider mb-4">目录</h3>
    <ul class="space-y-0.5 text-sm border-l-2 border-warm-200 dark:border-white/5">
      <li v-for="item in items" :key="item.id">
        <a :href="`#${item.id}`"
          :style="{ paddingLeft: item.level * 12 + 'px' }"
          class="block py-1.5 border-l-2 -ml-0.5 border-transparent text-warm-500 dark:text-warm-400 hover:text-brand-600 dark:hover:text-pop-400 hover:border-brand-500 dark:hover:border-pop-500 transition-colors">
          {{ item.text }}
        </a>
        <template v-if="item.children?.length">
          <a v-for="child in item.children" :key="child.id" :href="`#${child.id}`"
            :style="{ paddingLeft: child.level * 12 + 'px' }"
            class="block py-1 text-xs border-l-2 -ml-0.5 border-transparent text-warm-400 dark:text-warm-500 hover:text-brand-600 dark:hover:text-pop-400 hover:border-brand-500 dark:hover:border-pop-500 transition-colors">
            {{ child.text }}
          </a>
        </template>
      </li>
    </ul>
  </nav>
</template>
