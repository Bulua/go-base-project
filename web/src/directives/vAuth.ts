import type { DirectiveBinding, ObjectDirective } from 'vue'
import { useAuthStore } from '@/store/modules/auth'

/**
 * v-auth directive — hides element when user lacks the action permission.
 *
 * Usage:
 *   v-auth="'add'"            hide if user has no 'add' action on any menu
 *   v-auth="'11:add'"         hide if user has no 'add' action on menu 11
 *   v-auth="{ id: 11, code: 'add' }"
 */
export const vAuth: ObjectDirective = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    if (!check(binding.value)) {
      el.style.display = 'none'
    }
  },
  updated(el: HTMLElement, binding: DirectiveBinding) {
    el.style.display = check(binding.value) ? '' : 'none'
  },
}

function check(value: unknown): boolean {
  const authStore = useAuthStore()
  if (typeof value === 'string') {
    if (value.includes(':')) {
      return authStore.actionKeys.has(value)
    }
    return authStore.actionCodes.has(value)
  }
  if (value && typeof value === 'object') {
    const { id, code } = value as { id: number; code: string }
    return authStore.actionKeys.has(`${id}:${code}`)
  }
  return false
}
