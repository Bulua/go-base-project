import { request, unwrap } from '@/api/request'
import type {
  CreateParamPayload,
  MenuItem,
  MenuListQuery,
  MenuListResult,
  RouteParam,
  SaveMenuPayload,
} from '@/types/menu'

export function listMenus(query: MenuListQuery = {}): Promise<MenuListResult> {
  return unwrap<MenuListResult>(
    request.get('/api/v1/menus', {
      params: {
        page: query.page ?? 1,
        page_size: query.page_size ?? 20,
        keyword: query.keyword || undefined,
        menu_status: query.menu_status || undefined,
        menu_type: query.menu_type || undefined,
      },
    }),
  )
}

export function getMenuTree(): Promise<MenuItem[]> {
  return unwrap<MenuItem[]>(request.get('/api/v1/menus/tree'))
}

export function getMenu(id: number): Promise<MenuItem> {
  return unwrap<MenuItem>(request.get(`/api/v1/menus/${id}`))
}

export interface CreateMenuResult extends MenuItem {
  code_generated: boolean
  code_path?: string
}

export function createMenu(payload: SaveMenuPayload): Promise<CreateMenuResult> {
  return unwrap<CreateMenuResult>(request.post('/api/v1/menus', payload))
}

export function updateMenu(id: number, payload: SaveMenuPayload): Promise<MenuItem> {
  return unwrap<MenuItem>(request.put(`/api/v1/menus/${id}`, payload))
}

export function deleteMenu(id: number): Promise<{ success: boolean }> {
  return unwrap<{ success: boolean }>(request.delete(`/api/v1/menus/${id}`))
}

export function listMenuParams(menuId: number): Promise<RouteParam[]> {
  return unwrap<RouteParam[]>(request.get(`/api/v1/menus/${menuId}/params`))
}

export function createMenuParam(menuId: number, payload: CreateParamPayload): Promise<RouteParam> {
  return unwrap<RouteParam>(request.post(`/api/v1/menus/${menuId}/params`, payload))
}

export function deleteMenuParam(paramId: number): Promise<{ success: boolean }> {
  return unwrap<{ success: boolean }>(request.delete(`/api/v1/menus/params/${paramId}`))
}
