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
  try { return JSON.parse(props.tocJson) || [] }
  catch { return [] }
})
</script>

<template>
  <nav v-if="items.length" class="sticky top-4">
    <h3 class="font-semibold mb-2 text-sm text-gray-500">目录</h3>
    <ul class="space-y-1 text-sm">
      <li v-for="item in items" :key="item.id" :style="{ paddingLeft: (item.level - 1) * 12 + 'px' }">
        <a :href="`#${item.id}`" class="hover:text-blue-500 text-gray-600 dark:text-gray-400">
          {{ item.text }}
        </a>
        <template v-if="item.children?.length">
          <li v-for="child in item.children"
            :key="child.id"
            :style="{ paddingLeft: (child.level - 1) * 12 + 'px' }">
            <a :href="`#${child.id}`" class="hover:text-blue-500 text-gray-500 dark:text-gray-500">
              {{ child.text }}
            </a>
          </li>
        </template>
      </li>
    </ul>
  </nav>
</template>
