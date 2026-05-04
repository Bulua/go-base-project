export interface Dictionary {
  id: number
  dict_name: string
  dict_code: string
  dict_status: number
  parent_id: number
  remark?: string | null
  item_count: number
  created_at: string
  updated_at: string
}

export interface DictItem {
  id: number
  dict_id: number
  item_label: string
  item_value: string
  item_extra?: string | null
  item_status: number
  sort_no: number
  parent_id: number
  tree_level: number
  tree_path?: string | null
  created_at: string
  updated_at: string
}

export interface DictListQuery {
  page?: number
  page_size?: number
  keyword?: string
  dict_status?: number
}

export interface DictListResult {
  total: number
  items: Dictionary[]
  page: number
  page_size: number
}

export interface SaveDictPayload {
  dict_name: string
  dict_code: string
  dict_status: number
  parent_id: number
  remark?: string | null
}

export interface SaveItemPayload {
  item_label: string
  item_value: string
  item_extra?: string | null
  item_status: number
  sort_no: number
  parent_id: number
}
