import request from '@/utils/request'

export function listWarehouses(params: { page?: number; page_size?: number }) {
  return request.get('/warehouses', { params })
}

export function getWarehouse(id: number) {
  return request.get(`/warehouses/${id}`)
}

export function createWarehouse(data: Record<string, unknown>) {
  return request.post('/warehouses', data)
}

export function updateWarehouse(id: number, data: Record<string, unknown>) {
  return request.put(`/warehouses/${id}`, data)
}

export function changeWarehouseStatus(id: number, status: string) {
  return request.post(`/warehouses/${id}/status`, { status })
}

export function createShelf(data: Record<string, unknown>) {
  return request.post('/warehouses/shelves', data)
}

export function updateShelf(id: number, data: Record<string, unknown>) {
  return request.put(`/warehouses/shelves/${id}`, data)
}
