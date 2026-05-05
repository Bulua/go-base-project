import { defineComponent, h, type Component } from 'vue'
import type { MenuRoute } from '@/types/auth'
import type { TabItem } from '@/store/modules/tabs'
import router from './index'

type ViewLoader = () => Promise<Component>

const viewModules = import.meta.glob('@/views/**/*.vue') as Record<string, ViewLoader>

const REGISTERED_NAMES = new Set<string>()

interface FlatLeaf {
  fullPath: string
  componentPath: string
  routeName: string
  title: string
  keepAlive: boolean
  affix: boolean
}

function joinPath(parent: string, child: string): string {
  if (!child) return parent || '/'
  if (child.startsWith('/')) return child
  const sep = parent.endsWith('/') ? '' : '/'
  return `${parent || ''}${sep}${child}`.replace(/\/+/g, '/')
}

function flattenLeaves(menus: MenuRoute[], parentPath = ''): FlatLeaf[] {
  const out: FlatLeaf[] = []
  for (const m of menus) {
    const segment = m.route_path ?? ''
    const fullPath = joinPath(parentPath, segment)
    const isLeaf =
      m.menu_type === 2 && !!m.component_path && m.component_path !== 'layouts/default'
    if (isLeaf) {
      out.push({
        fullPath,
        componentPath: m.component_path!,
        routeName: m.route_name || `menu-${m.id}`,
        title: m.menu_title,
        keepAlive: m.is_keep_alive,
        affix: m.is_affix,
      })
    }
    if (m.children?.length) {
      out.push(...flattenLeaves(m.children, fullPath))
    }
  }
  return out
}

function resolveLoader(componentPath: string): ViewLoader | null {
  const normalized = componentPath.replace(/^\/+/, '').replace(/\.vue$/, '')
  const candidate = `/src/${normalized}.vue`
  const loader = viewModules[candidate]
  if (loader) return loader
  const lower = candidate.toLowerCase()
  const match = Object.keys(viewModules).find((k) => k.toLowerCase() === lower)
  return match ? viewModules[match] : null
}

export function registerDynamicRoutes(menus: MenuRoute[]): void {
  const leaves = flattenLeaves(menus)
  for (const leaf of leaves) {
    if (REGISTERED_NAMES.has(leaf.routeName)) continue
    const loader = resolveLoader(leaf.componentPath)
    if (!loader) {
      console.warn(
        `[dynamic-routes] component not found for path "${leaf.componentPath}" (menu: ${leaf.title})`,
      )
      continue
    }
    const routeName = leaf.routeName
    const wrappedLoader = async () => {
      const mod = await loader()
      const base = (mod as { default?: Component }).default ?? (mod as Component)
      return defineComponent({
        name: routeName,
        setup() {
          return () => h(base)
        },
      })
    }
    router.addRoute('Layout', {
      path: leaf.fullPath,
      name: leaf.routeName,
      component: wrappedLoader,
      meta: { title: leaf.title, keepAlive: leaf.keepAlive, affix: leaf.affix },
    })
    REGISTERED_NAMES.add(leaf.routeName)
  }
}

export function collectAffixTabs(menus: MenuRoute[]): TabItem[] {
  return flattenLeaves(menus)
    .filter((leaf) => leaf.affix)
    .map((leaf) => ({
      path: leaf.fullPath,
      name: leaf.routeName,
      title: leaf.title,
      affix: true,
      keepAlive: leaf.keepAlive,
    }))
}

export function clearDynamicRoutes(): void {
  for (const name of REGISTERED_NAMES) {
    if (router.hasRoute(name)) router.removeRoute(name)
  }
  REGISTERED_NAMES.clear()
}
