import { ref, watch, onMounted } from 'vue'
import api from '../utils/api'
import { generateCSSVars, hexToRgb } from '../utils/color'

export interface ThemePreset {
  name: string
  brand: string
  accent: string
}

export const presets: ThemePreset[] = [
  { name: 'Sakura', brand: '#fb7299', accent: '#f472b6' },
  { name: 'Spring', brand: '#22c55e', accent: '#14b8a6' },
  { name: 'Sky',    brand: '#3b82f6', accent: '#06b6d4' },
  { name: 'Lavender', brand: '#8b5cf6', accent: '#a855f7' },
]

const STORAGE_KEY = 'theme-preset'
const DEFAULT_PRESET = 'Sakura'

function findPreset(name: string): ThemePreset {
  return presets.find(p => p.name === name) || presets.find(p => p.name === DEFAULT_PRESET)!
}

function loadPresetName(): string {
  try {
    return localStorage.getItem(STORAGE_KEY) || DEFAULT_PRESET
  } catch {}
  return DEFAULT_PRESET
}

const activePreset = ref<string>(loadPresetName())

async function savePreset(name: string) {
  try {
    await api.put('/admin/settings', { theme_preset: name })
  } catch {
    // Offline or not admin; localStorage is already updated as cache
  }
}

export function useCustomTheme() {
  function apply() {
    const preset = findPreset(activePreset.value)
    const vars = generateCSSVars(preset.brand, preset.accent)
    const root = document.documentElement
    for (const [key, value] of Object.entries(vars)) {
      root.style.setProperty(key, value)
    }
    const brandRgb = hexToRgb(preset.brand)
    const accentRgb = hexToRgb(preset.accent)
    root.style.setProperty('--brand-hex', preset.brand)
    root.style.setProperty('--accent-hex', preset.accent)
    root.style.setProperty('--brand-rgb', `${brandRgb.r} ${brandRgb.g} ${brandRgb.b}`)
    root.style.setProperty('--accent-rgb', `${accentRgb.r} ${accentRgb.g} ${accentRgb.b}`)
  }

  function setPreset(name: string) {
    activePreset.value = name
  }

  onMounted(async () => {
    try {
      const r = await api.get('/settings')
      const s = r.data.settings as Record<string, string>
      if (s.theme_preset !== undefined && presets.some(p => p.name === s.theme_preset)) {
        activePreset.value = s.theme_preset
        localStorage.setItem(STORAGE_KEY, s.theme_preset)
      }
    } catch {
      // Use localStorage default already loaded
    }
    apply()
  })

  watch(activePreset, () => {
    apply()
    localStorage.setItem(STORAGE_KEY, activePreset.value)
    savePreset(activePreset.value)
  })

  return {
    presets,
    activePreset,
    apply,
    setPreset,
  }
}
