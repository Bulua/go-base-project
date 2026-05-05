<script setup lang="ts">
import { ref } from 'vue'
import { Check } from '@element-plus/icons-vue'
import { useTheme } from '@/composables/common/useTheme'
import { useWatermark } from '@/composables/common/useWatermark'

defineModel<boolean>({ required: true })

const { hue, setHue } = useTheme()
const { isWatermark, watermarkText, toggleWatermark } = useWatermark()

// ── Preset color swatches ──────────────────────────────────────────────────

interface Preset {
  label: string
  hue: number
}

const PRESETS: Preset[] = [
  { label: '蓝色', hue: 210 },
  { label: '紫色', hue: 270 },
  { label: '青色', hue: 185 },
  { label: '绿色', hue: 145 },
  { label: '橙色', hue: 42 },
  { label: '玫红', hue: 330 },
  { label: '红色', hue: 18 },
  { label: '黄绿', hue: 100 },
]

function swatchBg(h: number): string {
  return `oklch(0.52 0.18 ${h})`
}

function isActivePreset(h: number): boolean {
  return Math.abs(((hue.value - h) + 360) % 360) < 15
}

// ── Custom color picker ────────────────────────────────────────────────────

function hslChannel(p: number, q: number, t: number): number {
  t = ((t % 1) + 1) % 1
  if (t < 1 / 6) return p + (q - p) * 6 * t
  if (t < 1 / 2) return q
  if (t < 2 / 3) return p + (q - p) * (2 / 3 - t) * 6
  return p
}

// Approximate oklch(0.52 0.18 h) → hex via HSL for the color picker seed value.
function hueToHex(h: number): string {
  const s = 0.65
  const l = 0.45
  const q = l < 0.5 ? l * (1 + s) : l + s - l * s
  const p = 2 * l - q
  const r = hslChannel(p, q, h / 360 + 1 / 3)
  const g = hslChannel(p, q, h / 360)
  const b = hslChannel(p, q, h / 360 - 1 / 3)
  const hex = (v: number) => Math.round(v * 255).toString(16).padStart(2, '0')
  return `#${hex(r)}${hex(g)}${hex(b)}`
}

// Extract hue angle from a hex color string.
function hexToHue(hex: string): number {
  if (!/^#[0-9a-f]{6}$/i.test(hex)) return hue.value
  const r = parseInt(hex.slice(1, 3), 16) / 255
  const g = parseInt(hex.slice(3, 5), 16) / 255
  const b = parseInt(hex.slice(5, 7), 16) / 255
  const max = Math.max(r, g, b)
  const min = Math.min(r, g, b)
  const d = max - min
  if (d === 0) return 0
  let h: number
  if (max === r) h = ((g - b) / d % 6 + 6) % 6
  else if (max === g) h = (b - r) / d + 2
  else h = (r - g) / d + 4
  return Math.round(h * 60)
}

const pickerColor = ref(hueToHex(hue.value))

function applyPreset(presetHue: number) {
  setHue(presetHue)
  pickerColor.value = hueToHex(presetHue)
}

function onPickerChange(val: string | null) {
  if (val) setHue(hexToHue(val))
}
</script>

<template>
  <el-drawer
    :model-value="modelValue"
    title="外观设置"
    direction="rtl"
    size="300px"
    :append-to-body="true"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <!-- Theme Color -->
    <div class="sd-section">
      <div class="sd-section-title">主题颜色</div>

      <div class="sd-label">预设颜色</div>
      <div class="sd-preset-grid">
        <button
          v-for="p in PRESETS"
          :key="p.hue"
          class="sd-swatch"
          :style="{ background: swatchBg(p.hue) }"
          :title="p.label"
          @click="applyPreset(p.hue)"
        >
          <el-icon v-if="isActivePreset(p.hue)" class="sd-swatch-check"><Check /></el-icon>
        </button>
      </div>

      <div class="sd-label">自定义颜色</div>
      <div class="sd-custom-color-row">
        <el-color-picker
          v-model="pickerColor"
          size="default"
          @change="onPickerChange"
        />
        <span class="sd-hex-label">{{ pickerColor }}</span>
      </div>
    </div>

    <el-divider />

    <!-- Watermark -->
    <div class="sd-section">
      <div class="sd-section-title">水印</div>

      <div class="sd-row">
        <span class="sd-row-label">启用水印</span>
        <el-switch :model-value="isWatermark" @change="toggleWatermark" />
      </div>

      <template v-if="isWatermark">
        <div class="sd-label" style="margin-top: 12px">水印文字</div>
        <el-input
          v-model="watermarkText"
          placeholder="默认使用登录用户名"
          clearable
        />
      </template>
    </div>
  </el-drawer>
</template>

<style scoped>
.sd-section {
  margin-bottom: 8px;
}

.sd-section-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin-bottom: 16px;
}

.sd-label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-bottom: 8px;
}

.sd-preset-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
  margin-bottom: 20px;
}

.sd-swatch {
  width: 100%;
  aspect-ratio: 1;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.15s, box-shadow 0.15s;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
}

.sd-swatch:hover {
  transform: scale(1.1);
  box-shadow: 0 3px 8px rgba(0, 0, 0, 0.3);
}

.sd-swatch-check {
  color: #fff;
  font-size: 16px;
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.4));
}

.sd-custom-color-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.sd-hex-label {
  font-size: 13px;
  font-family: var(--font-mono, monospace);
  color: var(--el-text-color-regular);
}

.sd-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.sd-row-label {
  font-size: 13px;
  color: var(--el-text-color-primary);
}
</style>
