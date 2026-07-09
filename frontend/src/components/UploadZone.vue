<script setup lang="ts">
import { ref, computed } from 'vue'
import api from '../utils/api'
import KawaiiIcon from './KawaiiIcon.vue'

const emit = defineEmits<{ uploaded: [] }>()
const dragging = ref(false)
const uploading = ref(false)
const message = ref('')
const status = ref<'idle' | 'uploading' | 'success' | 'error'>('idle')
const fileNames = ref<string[]>([])

const hasMd = computed(() => fileNames.value.some(f => f.endsWith('.md')))

function onDragOver(e: DragEvent) { e.preventDefault(); dragging.value = true }
function onDragLeave() { dragging.value = false }

async function onDrop(e: DragEvent) {
  e.preventDefault(); dragging.value = false
  const items = e.dataTransfer?.items
  if (!items) return
  const files: File[] = []

  async function collect(entry: FileSystemEntry) {
    if (entry.isFile) {
      const file = await getFile(entry as FileSystemFileEntry)
      if (file) files.push(file)
    } else if (entry.isDirectory) {
      const reader = (entry as FileSystemDirectoryEntry).createReader()
      const entries = await readAllEntries(reader)
      for (const e of entries) await collect(e)
    }
  }

  const entries: FileSystemEntry[] = []
  for (let i = 0; i < items.length; i++) {
    const entry = items[i].webkitGetAsEntry?.()
    if (entry) entries.push(entry)
  }

  if (entries.length === 0) {
    const dropped = e.dataTransfer?.files
    if (dropped) for (let i = 0; i < dropped.length; i++) files.push(dropped[i])
  } else {
    for (const entry of entries) await collect(entry)
  }

  if (!files.length) return
  fileNames.value = files.map(f => (f as any).webkitRelativePath || f.name)

  if (!hasMd.value) {
    message.value = 'No .md file found in dropped files'
    status.value = 'error'
    return
  }

  await doUpload(files)
}

function getFile(entry: FileSystemFileEntry): Promise<File | null> {
  return new Promise(resolve => entry.file(resolve, () => resolve(null)))
}

function readAllEntries(reader: FileSystemDirectoryReader): Promise<FileSystemEntry[]> {
  return new Promise(resolve => {
    const all: FileSystemEntry[] = []
    const read = () => reader.readEntries(entries => {
      if (entries.length === 0) resolve(all)
      else { all.push(...entries); read() }
    })
    read()
  })
}

async function doUpload(files: File[]) {
  uploading.value = true; status.value = 'uploading'; message.value = ''
  try {
    const fd = new FormData()
    for (const f of files) {
      const relPath = (f as any).webkitRelativePath || f.name
      fd.append('file', f, relPath)
    }
    await api.post('/admin/upload', fd)
    const mdName = files.find(f => f.name.endsWith('.md'))?.name || ''
    const imgCount = files.filter(f => !f.name.endsWith('.md')).length
    message.value = imgCount > 0 ? `${mdName} + ${imgCount} image(s)` : mdName
    status.value = 'success'
    emit('uploaded')
  } catch (err: any) {
    message.value = err.response?.data?.error || 'Upload failed'
    status.value = 'error'
  } finally { uploading.value = false }
}

function retry() {
  status.value = 'idle'
  message.value = ''
  fileNames.value = []
}
</script>

<template>
  <div
    :class="[
      'min-h-[200px] rounded-2xl border-2 border-dashed p-8 text-center flex items-center justify-center transition-all duration-200 cursor-pointer',
      dragging
        ? 'border-brand-400 bg-brand-50 dark:bg-brand-900/20 scale-[1.01]'
        : 'border-gray-300 dark:border-white/10 bg-white dark:bg-slate-900 hover:border-brand-300 dark:hover:border-brand-700 hover:bg-gray-50 dark:hover:bg-white/5'
    ]"
    @dragover="onDragOver" @dragleave="onDragLeave" @drop="onDrop">

    <!-- Idle -->
    <div v-if="status === 'idle'" class="flex flex-col items-center gap-3">
      <div class="w-16 h-16">
        <KawaiiIcon name="upload" />
      </div>
      <p class="text-slate-600 dark:text-slate-300 font-medium">Drop <code class="px-1.5 py-0.5 rounded bg-gray-100 dark:bg-white/10 text-brand-600 dark:text-brand-400 text-sm">.md</code> + images here</p>
      <p class="text-xs text-slate-400">Supports folders — directory structure is preserved</p>
    </div>

    <!-- Uploading -->
    <div v-else-if="status === 'uploading'" class="flex flex-col items-center gap-3">
      <div class="w-8 h-8 border-2 border-brand-500 border-t-transparent rounded-full animate-spin"></div>
      <p class="text-slate-600 dark:text-slate-300 font-medium">Uploading &amp; parsing {{ fileNames.length }} file(s)...</p>
    </div>

    <!-- Success -->
    <div v-else-if="status === 'success'" class="flex flex-col items-center gap-3">
      <div class="w-16 h-16">
        <KawaiiIcon name="happy" />
      </div>
      <p class="text-green-600 dark:text-green-400 font-medium">Uploaded: {{ message }}</p>
      <button @click.stop="retry" class="text-sm text-slate-500 hover:text-brand-600 dark:hover:text-brand-400 underline">Upload more</button>
    </div>

    <!-- Error -->
    <div v-else class="flex flex-col items-center gap-3">
      <div class="w-16 h-16">
        <KawaiiIcon name="sad" />
      </div>
      <p class="text-red-500 font-medium">{{ message }}</p>
      <button @click.stop="retry" class="text-sm text-slate-500 hover:text-brand-600 dark:hover:text-brand-400 underline">Try again</button>
    </div>
  </div>
</template>
