import request from '@/utils/request'

export function queryStock(params: Record<string, unknown>) {
  return request.get('/stock-records', { params })
}

export function inbound(data: Record<string, unknown>) {
  return request.post('/stock-records/inbound', data)
}

export function outbound(data: Record<string, unknown>) {
  return request.post('/stock-records/outbound', data)
}

export function adjustStock(data: Record<string, unknown>) {
  return request.post('/stock-records/adjust', data)
}
