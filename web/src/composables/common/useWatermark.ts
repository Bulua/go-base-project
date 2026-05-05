import { ref, watchEffect } from 'vue'
import type { Ref } from 'vue'

const WATERMARK_KEY = 'bp-watermark'
const WATERMARK_TEXT_KEY = 'bp-watermark-text'

function readStoredBool(key: string): boolean {
  try {
    return localStorage.getItem(key) === '1'
  } catch {
    return false
  }
}

function readStoredStr(key: string): string {
  try {
    return localStorage.getItem(key) ?? ''
  } catch {
    return ''
  }
}

const isWatermark = ref<boolean>(readStoredBool(WATERMARK_KEY))
const watermarkText = ref<string>(readStoredStr(WATERMARK_TEXT_KEY))

let styleEl: HTMLStyleElement | null = null

function renderWatermark(text: string): void {
  const canvas = document.createElement('canvas')
  canvas.width = 240
  canvas.height = 140
  const ctx = canvas.getContext('2d')
  if (!ctx) return
  ctx.save()
  ctx.translate(canvas.width / 2, canvas.height / 2)
  ctx.rotate(-Math.PI / 6)
  ctx.font = '13px sans-serif'
  ctx.fillStyle = 'rgba(0,0,0,0.08)'
  ctx.textAlign = 'center'
  ctx.textBaseline = 'middle'
  ctx.fillText(text, 0, 0)
  ctx.restore()
  const url = canvas.toDataURL()
  if (!styleEl) {
    styleEl = document.createElement('style')
    document.head.appendChild(styleEl)
  }
  styleEl.textContent = [
    'body::after{content:"";position:fixed;inset:0;',
    'pointer-events:none;z-index:9999;',
    `background-image:url(${url});`,
    'background-repeat:repeat;background-size:240px 140px;}',
  ].join('')
}

function clearWatermark(): void {
  if (styleEl) styleEl.textContent = ''
}

// defaultText: fallback label when watermarkText is empty (typically the username).
// Only pass it from one call site (the root layout) to avoid duplicate watchEffects.
export function useWatermark(defaultText?: Ref<string> | (() => string)) {
  if (defaultText !== undefined) {
    watchEffect(() => {
      const custom = watermarkText.value.trim()
      const label = custom || (typeof defaultText === 'function' ? defaultText() : defaultText.value)
      try {
        localStorage.setItem(WATERMARK_KEY, isWatermark.value ? '1' : '0')
        localStorage.setItem(WATERMARK_TEXT_KEY, watermarkText.value)
      } catch {}
      if (isWatermark.value && label) {
        renderWatermark(label)
      } else {
        clearWatermark()
      }
    })
  }

  return {
    isWatermark,
    watermarkText,
    toggleWatermark: () => {
      isWatermark.value = !isWatermark.value
    },
  }
}
