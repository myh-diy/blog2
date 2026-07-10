import { ref, onMounted } from 'vue'

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

export function useBackgroundImage() {
  onMounted(() => {
    backgroundImage.value = localStorage.getItem(STORAGE_KEY) || ''
    const savedOpacity = parseFloat(localStorage.getItem(OPACITY_KEY) || '0.75')
    bgOpacity.value = Number.isNaN(savedOpacity) ? 0.75 : savedOpacity
    apply()
  })

  function setBackground(url: string) {
    backgroundImage.value = url
    localStorage.setItem(STORAGE_KEY, url)
    apply()
  }

  function clearBackground() {
    backgroundImage.value = ''
    localStorage.removeItem(STORAGE_KEY)
    apply()
  }

  function setOpacity(value: number) {
    bgOpacity.value = value
    localStorage.setItem(OPACITY_KEY, String(value))
  }

  return {
    backgroundImage,
    bgOpacity,
    setBackground,
    clearBackground,
    setOpacity,
  }
}
