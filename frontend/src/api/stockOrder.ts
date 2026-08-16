import request from '@/utils/request'

export function listOrders(params: Record<string, unknown>) {
  return request.get('/orders', { params })
}

export function getOrder(id: number) {
  return request.get(`/orders/${id}`)
}

export function createOrder(data: Record<string, unknown>) {
  return request.post('/orders', data)
}

export function submitOrder(id: number) {
  return request.post(`/orders/${id}/submit`)
}

export function approveOrder(id: number) {
  return request.post(`/orders/${id}/approve`)
}

export function executeOrder(id: number) {
  return request.post(`/orders/${id}/execute`)
}

export function cancelOrder(id: number) {
  return request.post(`/orders/${id}/cancel`)
}
