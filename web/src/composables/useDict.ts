import { ref } from 'vue'
import { getDictItems } from '@/api/dict'
import type { DictItem } from '@/types/dict'

const cache = new Map<string, DictItem[]>()
const pending = new Map<string, Promise<DictItem[]>>()

export function useDict(dictCode: string) {
  const items = ref<DictItem[]>(cache.get(dictCode) ?? [])
  const loading = ref(false)

  function load() {
    if (cache.has(dictCode)) {
      items.value = cache.get(dictCode)!
      return
    }
    if (!pending.has(dictCode)) {
      const p = getDictItems(dictCode)
        .then((data) => {
          cache.set(dictCode, data)
          return data
        })
        .finally(() => pending.delete(dictCode))
      pending.set(dictCode, p)
    }
    loading.value = true
    pending.get(dictCode)!
      .then((data) => { items.value = data })
      .finally(() => { loading.value = false })
  }

  function refresh() {
    cache.delete(dictCode)
    load()
  }

  // Label lookup helper
  function labelOf(value: string | number): string {
    const v = String(value)
    return items.value.find((i) => i.item_value === v)?.item_label ?? v
  }

  load()

  return { items, loading, refresh, labelOf }
}
