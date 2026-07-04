<script setup lang="ts">
import { ref } from 'vue'
import api from '../utils/api'

const emit = defineEmits<{ uploaded: [] }>()
const dragging = ref(false); const uploading = ref(false)
const message = ref(''); const isSuccess = ref(false)

function onDragOver(e: DragEvent) { e.preventDefault(); dragging.value = true }
function onDragLeave() { dragging.value = false }
async function onDrop(e: DragEvent) {
  e.preventDefault(); dragging.value = false
  const file = e.dataTransfer?.files[0]
  if (!file) return
  if (!file.name.endsWith('.md')) { message.value = 'Please drop a .md file'; isSuccess.value = false; return }
  uploading.value = true; message.value = ''
  try {
    const fd = new FormData(); fd.append('file', file)
    await api.post('/admin/upload', fd)
    message.value = `"${file.name}" uploaded!`; isSuccess.value = true; emit('uploaded')
  } catch (err: any) {
    message.value = err.response?.data?.error || 'Upload failed'; isSuccess.value = false
  } finally { uploading.value = false }
}
</script>

<template>
  <div :class="['relative border-2 border-dashed rounded-2xl p-16 text-center transition-all',
    dragging ? 'border-brand-500 dark:border-pop-500 bg-brand-50 dark:bg-pop-900/10 scale-[1.01]' : 'border-warm-300 dark:border-white/10 hover:border-warm-400 dark:hover:border-white/20 bg-white dark:bg-white/5']"
    @dragover="onDragOver" @dragleave="onDragLeave" @drop="onDrop">
    <div v-if="!uploading && !message" class="space-y-3">
      <div class="text-5xl">📄</div>
      <p class="text-warm-500 dark:text-warm-400 font-medium">Drop your <code class="px-1.5 py-0.5 bg-warm-100 dark:bg-white/10 rounded text-sm font-mono">.md</code> file here</p>
    </div>
    <div v-else-if="uploading" class="space-y-3">
      <div class="w-10 h-10 border-2 border-brand-500 dark:border-pop-500 border-t-transparent rounded-full animate-spin mx-auto"></div>
      <p class="text-brand-600 dark:text-pop-400 font-medium">Uploading &amp; parsing...</p>
    </div>
    <div v-else class="space-y-3">
      <div class="text-5xl">{{ isSuccess ? '✅' : '❌' }}</div>
      <p :class="isSuccess ? 'text-brand-600 dark:text-pop-400 font-medium' : 'text-red-500 font-medium'">{{ message }}</p>
      <button v-if="isSuccess" @click="message = ''" class="text-sm text-warm-400 hover:text-warm-600 dark:hover:text-warm-300 underline">Upload another</button>
    </div>
  </div>
</template>
