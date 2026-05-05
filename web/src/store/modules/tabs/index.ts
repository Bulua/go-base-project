import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

export interface TabItem {
  path: string
  name: string
  title: string
  affix?: boolean
  keepAlive?: boolean
}

export const useTabsStore = defineStore('tabs', () => {
  const list = ref<TabItem[]>([])

  const cachedNames = computed(() =>
    list.value.filter((t) => t.keepAlive).map((t) => t.name).filter(Boolean),
  )

  function initAffixTabs(tabs: TabItem[]) {
    for (const tab of tabs) {
      if (!list.value.some((t) => t.path === tab.path)) {
        list.value.unshift(tab)
      }
    }
  }

  function add(tab: TabItem) {
    if (list.value.some((t) => t.path === tab.path)) return
    list.value.push(tab)
  }

  function close(path: string): string | null {
    const idx = list.value.findIndex((t) => t.path === path)
    if (idx === -1) return null
    if (list.value[idx].affix) return path
    list.value.splice(idx, 1)
    if (list.value.length === 0) return null
    return (list.value[idx] ?? list.value[idx - 1])?.path ?? null
  }

  function closeOthers(keepPath: string) {
    list.value = list.value.filter((t) => t.path === keepPath || t.affix)
  }

  function closeAll() {
    list.value = list.value.filter((t) => t.affix)
  }

  function reset() {
    list.value = []
  }

  return { list, cachedNames, add, close, closeOthers, closeAll, initAffixTabs, reset }
})
