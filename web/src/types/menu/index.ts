export interface MenuItem {
  id: number
  parent_id: number
  menu_type: number
  route_path?: string | null
  route_name?: string | null
  component_path?: string | null
  redirect_path?: string | null
  menu_title: string
  menu_icon?: string | null
  sort_no: number
  is_hidden: boolean
  is_keep_alive: boolean
  is_affix: boolean
  active_route?: string | null
  transition_name?: string | null
  external_url?: string | null
  menu_status: number
  params?: RouteParam[]
  children?: MenuItem[]
  created_at: string
  updated_at: string
}

export interface RouteParam {
  id: number
  menu_id: number
  param_mode: string
  param_key: string
  param_value?: string | null
  created_at: string
}

export interface MenuListQuery {
  page?: number
  page_size?: number
  keyword?: string
  menu_status?: number
  menu_type?: number
}

export interface MenuListResult {
  total: number
  items: MenuItem[]
  page: number
  page_size: number
}

export interface SaveMenuPayload {
  parent_id: number
  menu_type: number
  route_path?: string | null
  route_name?: string | null
  component_path?: string | null
  redirect_path?: string | null
  menu_title: string
  menu_icon?: string | null
  sort_no: number
  is_hidden: boolean
  is_keep_alive: boolean
  is_affix: boolean
  active_route?: string | null
  transition_name?: string | null
  external_url?: string | null
  menu_status: number
}

export interface CreateParamPayload {
  param_mode: string
  param_key: string
  param_value?: string | null
}

export interface MenuAction {
  id: number
  menu_id: number
  action_code: string
  action_name: string
  action_desc?: string | null
  sort_no: number
  action_status: number
  created_at: string
  updated_at: string
}

export interface SaveActionPayload {
  action_code: string
  action_name: string
  action_desc?: string | null
  sort_no: number
  action_status: number
}
