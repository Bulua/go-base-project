export interface RoleItem {
  id: number
  role_code: string
  role_name: string
  parent_role_id: number
  default_route: string
  sort_no: number
  role_status: number
  remark?: string | null
  user_count: number
  created_at: string
  updated_at: string
  children?: RoleItem[]
}

export interface RoleListQuery {
  page?: number
  page_size?: number
  keyword?: string
  role_status?: number
}

export interface RoleListResult {
  total: number
  items: RoleItem[]
  page: number
  page_size: number
}

export interface CreateRolePayload {
  role_code: string
  role_name: string
  parent_role_id: number
  default_route: string
  sort_no: number
  role_status: number
  remark?: string | null
}

export interface UpdateRolePayload {
  role_code: string
  role_name: string
  parent_role_id: number
  default_route: string
  sort_no: number
  role_status: number
  remark?: string | null
}

export interface AssignIDsPayload {
  ids: number[]
}

export interface AssignedIDs {
  ids: number[]
}

export interface MenuOption {
  id: number
  parent_id: number
  menu_title: string
  menu_type: number
  sort_no: number
  children?: MenuOption[]
}

export interface ActionOption {
  id: number
  menu_id: number
  menu_title: string
  action_code: string
  action_name: string
  sort_no: number
}

export interface ApiOption {
  id: number
  api_path: string
  api_method: string
  api_group: string
  api_desc: string
  api_status: number
}

export interface RoleResources {
  menus: MenuOption[]
  actions: ActionOption[]
  apis: ApiOption[]
  roles: RoleItem[]
}
