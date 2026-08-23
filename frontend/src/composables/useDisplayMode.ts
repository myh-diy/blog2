import { ref } from 'vue'

export type DisplayMode = 'classic' | 'minimal'

const STORAGE_KEY = 'display-mode'
const DEFAULT_VERSION_KEY = 'display-mode-default-v2'
const DEFAULT_DISPLAY_MODE: DisplayMode = 'minimal'

function getInitialDisplayMode(): DisplayMode {
  if (localStorage.getItem(DEFAULT_VERSION_KEY) !== '1') {
    localStorage.setItem(DEFAULT_VERSION_KEY, '1')
    localStorage.setItem(STORAGE_KEY, DEFAULT_DISPLAY_MODE)
    return DEFAULT_DISPLAY_MODE
  }

  return localStorage.getItem(STORAGE_KEY) === 'classic' ? 'classic' : DEFAULT_DISPLAY_MODE
}

const displayMode = ref<DisplayMode>(getInitialDisplayMode())

export function useDisplayMode() {
  function setDisplayMode(mode: DisplayMode) {
    displayMode.value = mode
    localStorage.setItem(STORAGE_KEY, mode)
  }

  return { displayMode, setDisplayMode }
}
