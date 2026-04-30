export interface Role {
  id: number
  role_code: string
  role_name: string
  default_route: string
}

export interface CurrentUser {
  id: number
  user_uuid: string
  login_name: string
  display_name: string
  avatar_url?: string
  primary_role_id?: number
  phone_number?: string
  email_address?: string
  user_status: number
  must_change_password: boolean
  roles?: Role[]
}

export interface AuthSession {
  access_token: string
  refresh_token: string
  token_type: string
  expires_at: string
  refresh_expires_at: string
  user?: CurrentUser
}

export interface MenuRoute {
  id: number
  parent_id: number
  menu_type: number
  route_path?: string
  route_name?: string
  component_path?: string
  redirect_path?: string
  menu_title: string
  menu_icon?: string
  sort_no: number
  is_hidden: boolean
  is_keep_alive: boolean
  is_affix: boolean
  active_route?: string
  transition_name?: string
  external_url?: string
  children?: MenuRoute[]
}

export interface AuthAction {
  menu_id: number
  action_id: number
  action_code: string
  action_name: string
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
  trace_id?: string
}
