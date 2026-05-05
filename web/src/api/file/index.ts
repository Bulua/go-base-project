import { request, unwrap } from '@/api/request'
import type { FileListQuery, FileListResult, FileRecord } from '@/types/file'

export function listFiles(params: FileListQuery): Promise<FileListResult> {
  return unwrap<FileListResult>(request.get('/api/v1/files', { params }))
}

export function uploadFile(
  file: File,
  onProgress?: (pct: number) => void,
): Promise<FileRecord> {
  const form = new FormData()
  form.append('file', file)
  return unwrap<FileRecord>(
    request.post('/api/v1/files', form, {
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (e) => {
        onProgress?.(Math.round((e.loaded / (e.total ?? e.loaded)) * 100))
      },
    }),
  )
}

export function deleteFile(id: number): Promise<{ success: boolean }> {
  return unwrap<{ success: boolean }>(request.delete(`/api/v1/files/${id}`))
}

export function downloadFileBlob(id: number): Promise<Blob> {
  return request
    .get<Blob>(`/api/v1/files/${id}/raw`, { responseType: 'blob' })
    .then((r) => r.data)
}
