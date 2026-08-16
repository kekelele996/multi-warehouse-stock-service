export interface User {
  id: number
  phone: string
  name: string
  avatar: string
  role: string
  created_at: string
}

export interface Warehouse {
  id: number
  name: string
  code: string
  address: string
  area_sqm: number
  manager_id: number
  status: string
  phone: string
  created_at: string
}

export interface Shelf {
  id: number
  warehouse_id: number
  shelf_no: string
  layer_count: number
  column_count: number
  capacity: number
}

export interface Category {
  id: number
  name: string
  parent_id: number
  sort_order: number
  icon: string
}

export interface Product {
  id: number
  name: string
  sku: string
  category_id: number
  spec: string
  unit: string
  weight: number
  volume: number
  barcode: string
  min_stock: number
  max_stock: number
  shelf_life_days: number
  image_url: string
  created_at: string
}

export interface StockRecord {
  id: number
  product_id: number
  warehouse_id: number
  shelf_id: number
  batch_no: string
  quantity: number
  inbound_date: string | null
  expire_date: string | null
  last_op_type: string
  last_op_at: string
}

export interface StockOrder {
  id: number
  order_no: string
  order_type: string
  source_warehouse_id: number
  target_warehouse_id: number
  status: string
  creator_id: number
  approver_id: number
  remark: string
  created_at: string
}

export interface StockOrderItem {
  id: number
  order_id: number
  product_id: number
  shelf_id: number
  quantity: number
  actual_quantity: number
}

export interface OperationLog {
  id: number
  operator_id: number
  operator_name: string
  action: string
  entity_type: string
  entity_id: string
  detail: string
  ip: string
  created_at: string
}

export interface PageResult<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}
