import { request, unwrap } from '@/api/request'
import type {
  Dictionary,
  DictItem,
  DictListQuery,
  DictListResult,
  SaveDictPayload,
  SaveItemPayload,
} from '@/types/dict'

export function listDicts(query: DictListQuery = {}): Promise<DictListResult> {
  return unwrap<DictListResult>(
    request.get('/api/v1/dictionaries', {
      params: {
        page: query.page ?? 1,
        page_size: query.page_size ?? 20,
        keyword: query.keyword || undefined,
        dict_status: query.dict_status || undefined,
      },
    }),
  )
}

export function createDict(payload: SaveDictPayload): Promise<Dictionary> {
  return unwrap<Dictionary>(request.post('/api/v1/dictionaries', payload))
}

export function updateDict(id: number, payload: SaveDictPayload): Promise<Dictionary> {
  return unwrap<Dictionary>(request.put(`/api/v1/dictionaries/${id}`, payload))
}

export function deleteDict(id: number): Promise<{ success: boolean }> {
  return unwrap<{ success: boolean }>(request.delete(`/api/v1/dictionaries/${id}`))
}

export function getDictItems(dictCode: string): Promise<DictItem[]> {
  return unwrap<DictItem[]>(request.get(`/api/v1/dictionaries/${dictCode}/items`))
}

export function createDictItem(dictId: number, payload: SaveItemPayload): Promise<DictItem> {
  return unwrap<DictItem>(request.post(`/api/v1/dictionaries/${dictId}/items`, payload))
}

export function updateDictItem(id: number, payload: SaveItemPayload): Promise<DictItem> {
  return unwrap<DictItem>(request.put(`/api/v1/dictionary-items/${id}`, payload))
}

export function deleteDictItem(id: number): Promise<{ success: boolean }> {
  return unwrap<{ success: boolean }>(request.delete(`/api/v1/dictionary-items/${id}`))
}
