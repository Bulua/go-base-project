export interface UserRole {
  id: number
  role_code: string
  role_name: string
}

export interface UserItem {
  id: number
  user_uuid: string
  login_name: string
  display_name: string
  avatar_url?: string | null
  primary_role_id?: number | null
  phone_number?: string | null
  email_address?: string | null
  user_status: number
  must_change_password: boolean
  last_login_at?: string | null
  remark?: string | null
  roles?: UserRole[]
  created_at: string
  updated_at: string
}

export interface UserListQuery {
  page?: number
  page_size?: number
  keyword?: string
  user_status?: number
  role_id?: number
}

export interface UserListResult {
  total: number
  items: UserItem[]
  page: number
  page_size: number
}

export interface CreateUserPayload {
  login_name: string
  password: string
  display_name: string
  avatar_url?: string | null
  primary_role_id?: number | null
  phone_number?: string | null
  email_address?: string | null
  user_status: number
  remark?: string | null
  role_ids?: number[]
}

export interface UpdateUserPayload {
  display_name: string
  avatar_url?: string | null
  primary_role_id?: number | null
  phone_number?: string | null
  email_address?: string | null
  remark?: string | null
}

export interface ResetPasswordPayload {
  password: string
  must_change_password?: boolean
}

export interface AssignRolesPayload {
  role_ids: number[]
}

export interface UpdateStatusPayload {
  user_status: number
}
