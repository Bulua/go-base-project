import { request, unwrap } from '@/api/request'
import type {
  AssignRolesPayload,
  CreateUserPayload,
  ResetPasswordPayload,
  UpdateStatusPayload,
  UpdateUserPayload,
  UserItem,
  UserListQuery,
  UserListResult,
  UserRole,
} from '@/types/user'

export function listUsers(query: UserListQuery = {}): Promise<UserListResult> {
  return unwrap<UserListResult>(
    request.get('/api/v1/users', {
      params: {
        page: query.page ?? 1,
        page_size: query.page_size ?? 20,
        keyword: query.keyword || undefined,
        user_status: query.user_status || undefined,
        role_id: query.role_id || undefined,
      },
    }),
  )
}

export function getUser(id: number): Promise<UserItem> {
  return unwrap<UserItem>(request.get(`/api/v1/users/${id}`))
}

export function listRoleOptions(): Promise<UserRole[]> {
  return unwrap<UserRole[]>(request.get('/api/v1/users/role-options'))
}

export function createUser(payload: CreateUserPayload): Promise<UserItem> {
  return unwrap<UserItem>(request.post('/api/v1/users', payload))
}

export function updateUser(id: number, payload: UpdateUserPayload): Promise<UserItem> {
  return unwrap<UserItem>(request.put(`/api/v1/users/${id}`, payload))
}

export function deleteUser(id: number): Promise<{ success: boolean }> {
  return unwrap<{ success: boolean }>(request.delete(`/api/v1/users/${id}`))
}

export function updateUserStatus(id: number, payload: UpdateStatusPayload): Promise<{ success: boolean }> {
  return unwrap<{ success: boolean }>(request.put(`/api/v1/users/${id}/status`, payload))
}

export function resetUserPassword(id: number, payload: ResetPasswordPayload): Promise<{ success: boolean }> {
  return unwrap<{ success: boolean }>(request.put(`/api/v1/users/${id}/password`, payload))
}

export function assignUserRoles(id: number, payload: AssignRolesPayload): Promise<{ success: boolean }> {
  return unwrap<{ success: boolean }>(request.put(`/api/v1/users/${id}/roles`, payload))
}
