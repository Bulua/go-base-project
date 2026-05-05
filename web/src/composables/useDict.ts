import { ref, type Ref } from 'vue'
import { getDictItems } from '@/api/dict'
import type { DictItem } from '@/types/dict'

const sharedItems = new Map<string, Ref<DictItem[]>>()
const loaded = new Set<string>()
const pending = new Map<string, Promise<void>>()

function getSharedRef(dictCode: string): Ref<DictItem[]> {
  if (!sharedItems.has(dictCode)) {
    sharedItems.set(dictCode, ref<DictItem[]>([]))
  }
  return sharedItems.get(dictCode)!
}

function fetchAndUpdate(dictCode: string): Promise<void> {
  if (pending.has(dictCode)) return pending.get(dictCode)!
  const p = getDictItems(dictCode)
    .then((data) => {
      getSharedRef(dictCode).value = data
      loaded.add(dictCode)
    })
    .finally(() => pending.delete(dictCode))
  pending.set(dictCode, p)
  return p
}

// Call this after creating / updating / deleting a dict item.
// All consumers of useDict(dictCode) will reactively receive fresh data.
export function clearDictCache(dictCode: string) {
  loaded.delete(dictCode)
  pending.delete(dictCode)
  fetchAndUpdate(dictCode)
}

const PRESET_TAG_TYPES = new Set(['primary', 'success', 'warning', 'danger', 'info'])

export function useDict(dictCode: string) {
  const items = getSharedRef(dictCode)
  const loading = ref(false)

  function load() {
    if (loaded.has(dictCode)) return
    loading.value = true
    fetchAndUpdate(dictCode).finally(() => { loading.value = false })
  }

  function refresh() {
    loaded.delete(dictCode)
    pending.delete(dictCode)
    load()
  }

  function labelOf(value: string | number): string {
    const v = String(value)
    return items.value.find((i) => i.item_value === v)?.item_label ?? v
  }

  function typeOf(value: string | number): '' | 'primary' | 'success' | 'warning' | 'danger' | 'info' {
    const extra = items.value.find((i) => i.item_value === String(value))?.item_extra ?? ''
    return PRESET_TAG_TYPES.has(extra) ? (extra as '' | 'primary' | 'success' | 'warning' | 'danger' | 'info') : ''
  }

  // Returns the hex color string for use with el-tag's `color` prop.
  // Returns '' when item_extra is a preset type or empty.
  function colorOf(value: string | number): string {
    const extra = items.value.find((i) => i.item_value === String(value))?.item_extra ?? ''
    return extra.startsWith('#') ? extra : ''
  }

  load()

  return { items, loading, refresh, labelOf, typeOf, colorOf }
}
