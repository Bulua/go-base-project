import { request, unwrap } from '@/api/request'
import type {
  AssignIDsPayload,
  AssignedIDs,
  CreateRolePayload,
  RoleItem,
  RoleListQuery,
  RoleListResult,
  RoleResources,
  UpdateRolePayload,
} from '@/types/role'

export function listRoles(query: RoleListQuery = {}): Promise<RoleListResult> {
  return unwrap<RoleListResult>(
    request.get('/api/v1/roles', {
      params: {
        page: query.page ?? 1,
        page_size: query.page_size ?? 20,
        keyword: query.keyword || undefined,
        role_status: query.role_status || undefined,
      },
    }),
  )
}

export function getRoleTree(): Promise<RoleItem[]> {
  return unwrap<RoleItem[]>(request.get('/api/v1/roles/tree'))
}

export function getRoleResources(): Promise<RoleResources> {
  return unwrap<RoleResources>(request.get('/api/v1/roles/resources'))
}

export function getRole(id: number): Promise<RoleItem> {
  return unwrap<RoleItem>(request.get(`/api/v1/roles/${id}`))
}

export function createRole(payload: CreateRolePayload): Promise<RoleItem> {
  return unwrap<RoleItem>(request.post('/api/v1/roles', payload))
}

export function updateRole(id: number, payload: UpdateRolePayload): Promise<RoleItem> {
  return unwrap<RoleItem>(request.put(`/api/v1/roles/${id}`, payload))
}

export function deleteRole(id: number): Promise<{ success: boolean }> {
  return unwrap<{ success: boolean }>(request.delete(`/api/v1/roles/${id}`))
}

export function getRoleMenuIDs(id: number): Promise<AssignedIDs> {
  return unwrap<AssignedIDs>(request.get(`/api/v1/roles/${id}/menus`))
}

export function assignRoleMenus(id: number, payload: AssignIDsPayload): Promise<{ success: boolean }> {
  return unwrap<{ success: boolean }>(request.put(`/api/v1/roles/${id}/menus`, payload))
}

export function getRoleActionIDs(id: number): Promise<AssignedIDs> {
  return unwrap<AssignedIDs>(request.get(`/api/v1/roles/${id}/actions`))
}

export function assignRoleActions(id: number, payload: AssignIDsPayload): Promise<{ success: boolean }> {
  return unwrap<{ success: boolean }>(request.put(`/api/v1/roles/${id}/actions`, payload))
}

export function getRoleApiIDs(id: number): Promise<AssignedIDs> {
  return unwrap<AssignedIDs>(request.get(`/api/v1/roles/${id}/apis`))
}

export function assignRoleApis(id: number, payload: AssignIDsPayload): Promise<{ success: boolean }> {
  return unwrap<{ success: boolean }>(request.put(`/api/v1/roles/${id}/apis`, payload))
}

export function getRoleDataScopeIDs(id: number): Promise<AssignedIDs> {
  return unwrap<AssignedIDs>(request.get(`/api/v1/roles/${id}/data-scopes`))
}

export function assignRoleDataScopes(id: number, payload: AssignIDsPayload): Promise<{ success: boolean }> {
  return unwrap<{ success: boolean }>(request.put(`/api/v1/roles/${id}/data-scopes`, payload))
}
