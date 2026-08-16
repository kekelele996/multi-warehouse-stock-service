import request from '@/utils/request'

export function listOperationLogs(params: { page?: number; page_size?: number }) {
  return request.get('/operation-logs', { params })
}
