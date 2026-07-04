<script setup lang="ts">
import { ref, computed } from 'vue'
import api from '../utils/api'

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

  // Recursively collect files from directories
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
    // Fallback: use files directly
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

// Retry on error
function retry() {
  status.value = 'idle'
  message.value = ''
  fileNames.value = []
}
</script>

<template>
  <div :class="['upload-zone', dragging && 'upload-drag']"
    @dragover="onDragOver" @dragleave="onDragLeave" @drop="onDrop">
    <!-- Idle -->
    <div v-if="status === 'idle'" class="upload-content">
      <div class="upload-icon">📂</div>
      <p class="upload-text">Drop <code>.md</code> + images here</p>
      <p class="upload-hint">Supports folders — directory structure is preserved</p>
    </div>

    <!-- Uploading -->
    <div v-else-if="status === 'uploading'" class="upload-content">
      <div class="upload-spin"></div>
      <p class="upload-text">Uploading &amp; parsing {{ fileNames.length }} file(s)...</p>
    </div>

    <!-- Success -->
    <div v-else-if="status === 'success'" class="upload-content">
      <div class="upload-icon">✅</div>
      <p class="text-green-600 dark:text-green-400 font-medium">{{ message }}</p>
      <button @click="retry" class="upload-retry">Upload more</button>
    </div>

    <!-- Error -->
    <div v-else class="upload-content">
      <div class="upload-icon">❌</div>
      <p class="text-red-500 font-medium">{{ message }}</p>
      <button @click="retry" class="upload-retry">Try again</button>
    </div>
  </div>
</template>

<style>
.upload-zone {
  border: 2px dashed #d6d3d1; border-radius: 1rem; padding: 3rem 1rem;
  text-align: center; transition: all .2s; background: #fff; cursor: pointer;
  min-height: 180px; display: flex; align-items: center; justify-content: center;
}
.dark .upload-zone { background: rgba(255,255,255,.03); border-color: rgba(255,255,255,.1); }
.upload-drag { border-color: #f97316; background: #fff7ed; transform: scale(1.01); border-style: solid; }
.dark .upload-drag { border-color: #f43f5e; background: rgba(244,63,94,.06); }
.upload-content { display: flex; flex-direction: column; align-items: center; gap: .5rem; }
.upload-icon { font-size: 3rem; }
.upload-text { color: #78716c; font-weight: 500; }
.upload-text code { background: #f5f5f4; padding: .15rem .4rem; border-radius: .25rem; font-size: .85em; }
.dark .upload-text { color: #a8a29e; }
.dark .upload-text code { background: rgba(255,255,255,.08); }
.upload-hint { font-size: .75rem; color: #a8a29e; }
.upload-spin { width: 2rem; height: 2rem; border: 2px solid #f97316; border-top-color: transparent; border-radius: 50%; animation: spin .6s linear infinite; }
.dark .upload-spin { border-color: #f43f5e; border-top-color: transparent; }
@keyframes spin { to { transform: rotate(360deg); } }
.upload-retry { font-size: .8rem; color: #78716c; text-decoration: underline; background: none; border: none; cursor: pointer; margin-top: .5rem; }
</style>
