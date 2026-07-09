import { ref, onMounted } from 'vue'

const STORAGE_KEY = 'blog-bg-image'

const backgroundImage = ref('')

function apply() {
  document.documentElement.style.setProperty(
    '--bg-image',
    backgroundImage.value ? `url("${backgroundImage.value}")` : ''
  )
}

export function useBackgroundImage() {
  onMounted(() => {
    backgroundImage.value = localStorage.getItem(STORAGE_KEY) || ''
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

  return {
    backgroundImage,
    setBackground,
    clearBackground,
  }
}
