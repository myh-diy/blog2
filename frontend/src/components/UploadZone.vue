<script setup lang="ts">
import { ref } from 'vue'
import api from '../utils/api'

const emit = defineEmits<{ uploaded: [] }>()
const dragging = ref(false)
const uploading = ref(false)
const message = ref('')

function onDragOver(e: DragEvent) {
  e.preventDefault()
  dragging.value = true
}
function onDragLeave() { dragging.value = false }
async function onDrop(e: DragEvent) {
  e.preventDefault()
  dragging.value = false
  const file = e.dataTransfer?.files[0]
  if (!file) return
  if (!file.name.endsWith('.md')) {
    message.value = 'Please drop a .md file'
    return
  }
  uploading.value = true
  try {
    const form = new FormData()
    form.append('file', file)
    await api.post('/admin/upload', form)
    message.value = `"${file.name}" uploaded!`
    emit('uploaded')
  } catch (err: any) {
    message.value = err.response?.data?.error || 'Upload failed'
  } finally {
    uploading.value = false
  }
}
</script>

<template>
  <div
    :class="['border-2 border-dashed rounded-lg p-12 text-center transition',
      dragging ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20' : 'border-gray-300 dark:border-gray-600']"
    @dragover="onDragOver"
    @dragleave="onDragLeave"
    @drop="onDrop"
  >
    <p v-if="uploading" class="text-blue-500">Uploading...</p>
    <p v-else class="text-gray-500">Drop .md file here</p>
    <p v-if="message" class="mt-2 text-sm" :class="message.includes('!') ? 'text-green-500' : 'text-red-500'">
      {{ message }}
    </p>
  </div>
</template>
