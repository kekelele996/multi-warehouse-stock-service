import request from '@/utils/request'

export function listProducts(params: { page?: number; page_size?: number; category_id?: number; keyword?: string }) {
  return request.get('/products', { params })
}

export function getProduct(id: number) {
  return request.get(`/products/${id}`)
}

export function createProduct(data: Record<string, unknown>) {
  return request.post('/products', data)
}

export function updateProduct(id: number, data: Record<string, unknown>) {
  return request.put(`/products/${id}`, data)
}

export function listCategories() {
  return request.get('/categories')
}

export function createCategory(data: Record<string, unknown>) {
  return request.post('/categories', data)
}
