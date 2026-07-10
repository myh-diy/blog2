import { ref, onMounted } from 'vue'
import api from '../utils/api'

const STORAGE_KEY = 'blog-bg-image'
const OPACITY_KEY = 'blog-bg-opacity'

const backgroundImage = ref('')
const bgOpacity = ref(0.75)

function apply() {
  document.documentElement.style.setProperty(
    '--bg-image',
    backgroundImage.value ? `url("${backgroundImage.value}")` : ''
  )
}

async function saveSettings() {
  try {
    await api.put('/admin/settings', {
      background_image: backgroundImage.value,
      background_opacity: bgOpacity.value,
    })
  } catch {
    // Offline or not admin; localStorage is already updated as cache
  }
}

export function useBackgroundImage() {
  onMounted(async () => {
    try {
      const r = await api.get('/settings')
      const s = r.data.settings as Record<string, string>
      if (s.background_image !== undefined) {
        backgroundImage.value = s.background_image
        localStorage.setItem(STORAGE_KEY, s.background_image)
      } else {
        backgroundImage.value = localStorage.getItem(STORAGE_KEY) || ''
      }
      if (s.background_opacity !== undefined) {
        const v = parseFloat(s.background_opacity)
        bgOpacity.value = Number.isNaN(v) ? 0.75 : v
        localStorage.setItem(OPACITY_KEY, String(bgOpacity.value))
      } else {
        const saved = parseFloat(localStorage.getItem(OPACITY_KEY) || '0.75')
        bgOpacity.value = Number.isNaN(saved) ? 0.75 : saved
      }
    } catch {
      backgroundImage.value = localStorage.getItem(STORAGE_KEY) || ''
      const saved = parseFloat(localStorage.getItem(OPACITY_KEY) || '0.75')
      bgOpacity.value = Number.isNaN(saved) ? 0.75 : saved
    }
    apply()
  })

  function setBackground(url: string) {
    backgroundImage.value = url
    localStorage.setItem(STORAGE_KEY, url)
    apply()
    saveSettings()
  }

  function clearBackground() {
    backgroundImage.value = ''
    localStorage.removeItem(STORAGE_KEY)
    apply()
    saveSettings()
  }

  function setOpacity(value: number) {
    bgOpacity.value = value
    localStorage.setItem(OPACITY_KEY, String(value))
    saveSettings()
  }

  return {
    backgroundImage,
    bgOpacity,
    setBackground,
    clearBackground,
    setOpacity,
  }
}
