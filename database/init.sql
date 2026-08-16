-- 多仓库库存管理系统 (warehouse-stock) 初始化脚本：首次启动容器时自动执行
SET NAMES utf8mb4;
CREATE DATABASE IF NOT EXISTS warehouse_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE warehouse_db;

CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  phone VARCHAR(20) NOT NULL,
  password_hash VARCHAR(100) NOT NULL,
  name VARCHAR(50) NOT NULL DEFAULT '',
  avatar VARCHAR(255) NOT NULL DEFAULT '',
  role VARCHAR(30) NOT NULL DEFAULT 'viewer',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id),
  UNIQUE KEY uk_users_phone (phone)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS warehouses (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  code VARCHAR(50) NOT NULL,
  address VARCHAR(255) NOT NULL DEFAULT '',
  area_sqm DECIMAL(10,2) NOT NULL DEFAULT 0,
  manager_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  status VARCHAR(30) NOT NULL DEFAULT 'active',
  phone VARCHAR(20) NOT NULL DEFAULT '',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id),
  UNIQUE KEY uk_warehouses_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS shelves (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  warehouse_id BIGINT UNSIGNED NOT NULL,
  shelf_no VARCHAR(50) NOT NULL,
  layer_count INT NOT NULL DEFAULT 1,
  column_count INT NOT NULL DEFAULT 1,
  capacity INT NOT NULL DEFAULT 100,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id),
  KEY idx_shelves_warehouse (warehouse_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS categories (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  sort_order INT NOT NULL DEFAULT 0,
  icon VARCHAR(100) NOT NULL DEFAULT '',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS products (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(200) NOT NULL,
  sku VARCHAR(50) NOT NULL,
  category_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  spec VARCHAR(100) NOT NULL DEFAULT '',
  unit VARCHAR(20) NOT NULL DEFAULT '件',
  weight DECIMAL(10,2) NOT NULL DEFAULT 0,
  volume DECIMAL(10,2) NOT NULL DEFAULT 0,
  barcode VARCHAR(100) NOT NULL DEFAULT '',
  min_stock INT NOT NULL DEFAULT 0,
  max_stock INT NOT NULL DEFAULT 0,
  shelf_life_days INT NOT NULL DEFAULT 0,
  image_url VARCHAR(255) NOT NULL DEFAULT '',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id),
  UNIQUE KEY uk_products_sku (sku)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS stock_records (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  product_id BIGINT UNSIGNED NOT NULL,
  warehouse_id BIGINT UNSIGNED NOT NULL,
  shelf_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  batch_no VARCHAR(50) NOT NULL DEFAULT '',
  quantity INT NOT NULL DEFAULT 0,
  inbound_date DATE,
  expire_date DATE,
  last_op_type VARCHAR(30) NOT NULL DEFAULT 'inbound',
  last_op_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id),
  KEY idx_records_product (product_id),
  KEY idx_records_warehouse (warehouse_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS stock_orders (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  order_no VARCHAR(50) NOT NULL,
  order_type VARCHAR(30) NOT NULL DEFAULT 'inbound',
  source_warehouse_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  target_warehouse_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  status VARCHAR(30) NOT NULL DEFAULT 'draft',
  creator_id BIGINT UNSIGNED NOT NULL,
  approver_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id),
  UNIQUE KEY uk_orders_order_no (order_no)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS stock_order_items (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  order_id BIGINT UNSIGNED NOT NULL,
  product_id BIGINT UNSIGNED NOT NULL,
  shelf_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  quantity INT NOT NULL DEFAULT 0,
  actual_quantity INT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  KEY idx_order_items_order (order_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS operation_logs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  operator_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  operator_name VARCHAR(50) NOT NULL DEFAULT '',
  action VARCHAR(50) NOT NULL,
  entity_type VARCHAR(50) NOT NULL,
  entity_id VARCHAR(50) NOT NULL DEFAULT '',
  detail TEXT,
  ip VARCHAR(50) NOT NULL DEFAULT '',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 预置种子数据（密码：admin/Admin@123，其余/User@123）
INSERT INTO users (id, phone, password_hash, name, avatar, role, created_at) VALUES
(1, '13800000001', '$2a$10$bFfMuQAuKWflKxpuDYdFpeGJPVgD83q/.278LHYLL5S0DDmEfChX2', '系统管理员', '', 'admin', NOW(3)),
(2, '13800000002', '$2a$10$TMTpnDbEwRbtcbF9VJxAxe5IswQjmo7pboKI9zVtU.BYnhzdJpX9a', '仓库经理', '', 'warehouse_manager', NOW(3)),
(3, '13800000003', '$2a$10$TMTpnDbEwRbtcbF9VJxAxe5IswQjmo7pboKI9zVtU.BYnhzdJpX9a', '仓管员', '', 'operator', NOW(3)),
(4, '13800000004', '$2a$10$TMTpnDbEwRbtcbF9VJxAxe5IswQjmo7pboKI9zVtU.BYnhzdJpX9a', '观察员', '', 'viewer', NOW(3));

INSERT INTO warehouses (id, name, code, address, area_sqm, manager_id, status, phone, created_at) VALUES
(1, '华东一号仓', 'WH-001', '上海市青浦区物流园 1 号', 12000.00, 2, 'active', '021-88880001', NOW(3)),
(2, '华南二号仓', 'WH-002', '广州市黄埔区仓储大道 8 号', 8000.00, 2, 'active', '020-88880002', NOW(3));

INSERT INTO shelves (id, warehouse_id, shelf_no, layer_count, column_count, capacity, created_at) VALUES
(1, 1, 'A-01', 3, 4, 120, NOW(3)),
(2, 1, 'A-02', 3, 4, 120, NOW(3)),
(3, 2, 'B-01', 2, 4, 80, NOW(3));

INSERT INTO categories (id, name, parent_id, sort_order, icon, created_at) VALUES
(1, '电子元器件', 0, 1, '', NOW(3)),
(2, '包装材料', 0, 2, '', NOW(3)),
(3, '成品', 0, 3, '', NOW(3));

INSERT INTO products (id, name, sku, category_id, spec, unit, weight, volume, barcode, min_stock, max_stock, shelf_life_days, image_url, created_at) VALUES
(1, '电阻 10KΩ', 'ELEC-R10K', 1, '0805', '件', 0.01, 0.001, '6900000000011', 100, 5000, 0, '', NOW(3)),
(2, '电容 100uF', 'ELEC-C100', 1, '1206', '件', 0.02, 0.002, '6900000000028', 200, 8000, 0, '', NOW(3)),
(3, '气泡膜', 'PKG-BUBBLE', 2, '50cm*100m', '卷', 1.50, 0.05, '6900000000035', 10, 200, 365, '', NOW(3));

INSERT INTO stock_records (id, product_id, warehouse_id, shelf_id, batch_no, quantity, inbound_date, expire_date, last_op_type, last_op_at) VALUES
(1, 1, 1, 1, 'B20260801', 2000, '2026-08-01', NULL, 'inbound', NOW(3)),
(2, 2, 1, 2, 'B20260802', 3500, '2026-08-02', NULL, 'inbound', NOW(3)),
(3, 3, 2, 3, 'B20260803', 60, '2026-08-03', '2027-08-01', 'inbound', NOW(3));

INSERT INTO stock_orders (id, order_no, order_type, source_warehouse_id, target_warehouse_id, status, creator_id, approver_id, remark, created_at) VALUES
(1, 'IN202608160001', 'inbound', 0, 1, 'completed', 3, 2, '电阻入库', NOW(3)),
(2, 'OUT202608160001', 'outbound', 1, 0, 'submitted', 3, 0, '电容出库', NOW(3));

INSERT INTO stock_order_items (id, order_id, product_id, shelf_id, quantity, actual_quantity) VALUES
(1, 1, 1, 1, 500, 500),
(2, 2, 2, 2, 200, 0);

INSERT INTO operation_logs (id, operator_id, operator_name, action, entity_type, entity_id, detail, ip, created_at) VALUES
(1, 1, 'admin', 'seed', 'system', '', 'init', '127.0.0.1', NOW(3));
