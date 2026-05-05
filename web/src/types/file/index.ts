export interface FileRecord {
  id: number
  original_name: string
  storage_key: string
  file_size: number
  mime_type: string
  uploader_id?: number
  created_at: string
}

export interface FileListQuery {
  page: number
  page_size: number
  keyword?: string
}

export interface FileListResult {
  total: number
  items: FileRecord[]
  page: number
  page_size: number
}
