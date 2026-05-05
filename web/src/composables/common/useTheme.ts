import { ref, watchEffect } from 'vue'

const HUE_KEY = 'bp-hue'
const DARK_KEY = 'bp-dark'

const DEFAULT_HUE = 250

function readStoredHue(): number {
  try {
    const v = localStorage.getItem(HUE_KEY)
    if (v) {
      const n = parseFloat(v)
      if (!Number.isNaN(n)) return n
    }
  } catch {}
  return DEFAULT_HUE
}

function readStoredDark(): boolean {
  try {
    return localStorage.getItem(DARK_KEY) === '1'
  } catch {
    return false
  }
}

const hue = ref<number>(readStoredHue())
const isDark = ref<boolean>(readStoredDark())

watchEffect(() => {
  document.documentElement.style.setProperty('--hue', String(hue.value))
  try {
    localStorage.setItem(HUE_KEY, String(hue.value))
  } catch {}
})

watchEffect(() => {
  document.documentElement.setAttribute('data-theme', isDark.value ? 'dark' : 'light')
  document.documentElement.classList.toggle('dark', isDark.value)
  try {
    localStorage.setItem(DARK_KEY, isDark.value ? '1' : '0')
  } catch {}
})

export function useTheme() {
  return {
    hue,
    isDark,
    toggleDark: () => {
      isDark.value = !isDark.value
    },
    setHue: (v: number) => {
      hue.value = v
    },
  }
}
