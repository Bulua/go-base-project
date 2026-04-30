import { describe, expect, it } from 'vitest'
import { AxiosError } from 'axios'
import { unwrap } from '@/api/request'

describe('api request unwrap', () => {
  it('uses backend message from non-2xx response bodies', async () => {
    const error = new AxiosError('Request failed with status code 409')
    error.response = {
      data: {
        code: 409201,
        message: '登录账号已存在',
        trace_id: 'trace-001',
      },
      status: 409,
      statusText: 'Conflict',
      headers: {},
      config: {} as never,
    }

    await expect(unwrap(Promise.reject(error))).rejects.toThrow('登录账号已存在')
  })

  it('uses backend message from successful transport responses with non-zero business code', async () => {
    await expect(
      unwrap(
        Promise.resolve({
          data: {
            code: 400001,
            message: '请求参数格式不正确',
            data: null,
          },
        }),
      ),
    ).rejects.toThrow('请求参数格式不正确')
  })
})
