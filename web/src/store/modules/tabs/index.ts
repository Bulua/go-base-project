import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

export interface TabItem {
  path: string   // full route path, used as unique key
  name: string   // route name, used for <keep-alive :include>
  title: string
}

export const useTabsStore = defineStore('tabs', () => {
  const list = ref<TabItem[]>([])

  const cachedNames = computed(() => list.value.map((t) => t.name).filter(Boolean))

  function add(tab: TabItem) {
    if (list.value.some((t) => t.path === tab.path)) return
    list.value.push(tab)
  }

  function close(path: string): string | null {
    const idx = list.value.findIndex((t) => t.path === path)
    if (idx === -1) return null
    list.value.splice(idx, 1)
    if (list.value.length === 0) return null
    return (list.value[idx] ?? list.value[idx - 1])?.path ?? null
  }

  function closeOthers(keepPath: string) {
    list.value = list.value.filter((t) => t.path === keepPath)
  }

  function closeAll() {
    list.value = []
  }

  function reset() {
    list.value = []
  }

  return { list, cachedNames, add, close, closeOthers, closeAll, reset }
})
