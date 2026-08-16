# WarehouseStock（多仓库库存管理系统）

面向中小型仓储企业的多仓库/多库位精细化管理平台，覆盖入库、出库、盘点、调拨全流程。

## 快速启动（Docker Compose 一键部署）

```bash
cp .env.example .env
docker compose up -d --build
```

启动完成后访问：

- 前端：http://localhost:18705
- 后端 API：http://localhost:19205
- 后端健康检查：http://localhost:19205/healthz
- MySQL：localhost:3306

预置账号（密码见 database/init.sql 与 backend/internal/service/seed.go）：

| 手机号 | 密码 | 角色 |
| --- | --- | --- |
| 13800000001 | Admin@123 | 管理员 |
| 13800000002 | User@123 | 仓库经理 |
| 13800000003 | User@123 | 仓管员 |
| 13800000004 | User@123 | 观察员 |

## 本地开发

后端：

```bash
cd backend && go mod tidy && go run ./cmd/server
```

构建：`cd backend && go build ./...`

前端：

```bash
cd frontend && npm install && npm run dev
```

前端开发服务器通过 Vite 代理将 `/api` 转发到 `http://localhost:19205`。

## 技术栈

| 层 | 技术 |
| --- | --- |
| 前端 | Vue 3 + TypeScript + Element Plus + ECharts + Pinia + Vite |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | MySQL 8.0 |
| 认证 | JWT（github.com/golang-jwt/jwt/v5）+ RBAC |
| 其他依赖 | gin-contrib/cors、golang.org/x/crypto/bcrypt、go-playground/validator/v10 |

## 项目目录结构

```
wje-135/
├── docker-compose.yml
├── .env.example
├── database/init.sql              # MySQL 初始化脚本（建表 + 种子数据）
├── backend/
│   ├── cmd/server/main.go
│   └── internal/
│       ├── config/
│       ├── model/                 # user/warehouse/shelf/category/product/stock_record/stock_order/stock_order_item/operation_log
│       ├── repository/            # 按实体分文件
│       ├── service/               # 业务逻辑（入库/出库/调拨/盘点）+ dashboard + 种子数据 + 单元测试
│       ├── handler/               # 按实体分文件（含 upload_handler、operation_log_handler）
│       ├── router/                # router.go + 按实体分文件
│       ├── middleware/            # auth/rbac/audit_log/error_handler/rate_limiter/cors/request_logger
│       ├── dto/
│       ├── constants/             # 枚举、错误码、日志模板、文案
│       └── util/                  # jwt/logger/formatters/app_error/file_upload
└── frontend/
    └── src/
        ├── api/                   # user/warehouse/product/stockRecord/stockOrder/dashboard/operationLog
        ├── stores/                # authStore/userStore/warehouseStore/productStore/stockStore/orderStore
        ├── types/
        ├── components/common/     # StatusBadge/StockAlert/OccupancyBar/EmptyState/RoleGuard/ErrorBoundary
        ├── hooks/                 # useAuth/usePagination/useStockQuery
        ├── pages/                 # Dashboard/WarehouseManage/ProductManage/OrderManage/InventoryCheck/Profile/OperationLogs/Login
        ├── router/                # index.ts + guards.ts
        ├── utils/                 # formatStockNumber/orderNo/request
        └── constants/             # order/user/errorCodes
```

## 环境变量

| 变量 | 说明 | 默认值 |
| --- | --- | --- |
| COMPOSE_PROJECT_NAME | Docker Compose 项目名/容器前缀 | warehouse-stock |
| DB_NAME | 数据库名 | warehouse_db |
| DB_USER | 数据库用户 | warehouse_user |
| DB_PASSWORD | 数据库密码 | warehouse_pwd |
| DB_ROOT_PASSWORD | 数据库 root 密码 | warehouse_root |
| JWT_SECRET | JWT 签名密钥 | change_me_to_a_long_random_string |
| JWT_EXPIRE_HOURS | JWT 过期小时数 | 72 |
| APP_CORS_ORIGINS | 允许的跨域来源（逗号分隔） | http://localhost:18705 |
| FRONTEND_PORT | 前端端口 | 18705 |
| BACKEND_PORT | 后端端口 | 19205 |
| DB_PORT | 数据库端口 | 3306 |

## Docker 部署说明

- 端口映射：前端 18705:80、后端 19205:8080、MySQL 3306:3306。
- 数据卷：`db-data` 持久化 MySQL 数据；`upload-data` 持久化上传图片。
- 服务依赖：backend `depends_on` db（service_healthy），frontend `depends_on` backend（service_healthy）。
- 前端 Nginx 将 `/api/` 反代到 `http://backend:8080/`，支持 SPA 路由 `try_files`。
- 常见问题：
  - 端口冲突：修改 `.env` 中 `FRONTEND_PORT/BACKEND_PORT/DB_PORT` 后重新 `docker compose up -d`。
  - 数据库重置：`docker compose down -v` 后重新启动。

## 枚举出现位置清单

### OrderType（inbound/outbound/transfer/inventory_check）
- 后端：`backend/internal/constants/order.go`、`backend/internal/model/stock_order.go`、`backend/internal/service/stock_order_service.go`、`backend/internal/util/formatters.go`、`backend/internal/constants/log_templates.go`、`backend/internal/constants/error_codes.go`、`backend/internal/dto/dto_order.go`、`database/init.sql`
- 前端：`frontend/src/constants/order.ts`、`frontend/src/pages/OrderManage.vue`、`frontend/src/utils/orderNo.ts`

### OrderStatus（draft/submitted/processing/completed/cancelled）
- 后端：`backend/internal/constants/order.go`、`backend/internal/model/stock_order.go`、`backend/internal/service/stock_order_service.go`、`backend/internal/util/formatters.go`、`backend/internal/constants/log_templates.go`、`backend/internal/constants/error_codes.go`、`database/init.sql`
- 前端：`frontend/src/constants/order.ts`、`frontend/src/components/common/StatusBadge.vue`、`frontend/src/pages/OrderManage.vue`

### UserRole（admin/warehouse_manager/operator/viewer）
- 后端：`backend/internal/constants/user.go`、`backend/internal/model/user.go`、`backend/internal/middleware/rbac.go`、`backend/internal/router/*.go`、`backend/internal/util/formatters.go`、`database/init.sql`
- 前端：`frontend/src/constants/user.ts`、`frontend/src/stores/authStore.ts`、`frontend/src/components/common/RoleGuard.vue`、`frontend/src/router/guards.ts`、`frontend/src/pages/Login.vue`

## API 接口清单

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /healthz | 服务健康检查 |
| GET | /api/healthz | Nginx 反代健康检查 |
| GET | /api/v1/healthz | API 版本健康检查 |
| POST | /api/v1/auth/register | 用户注册 |
| POST | /api/v1/auth/login | 用户登录，返回 JWT |
| GET | /api/v1/users/me | 当前登录用户信息 |
| PUT | /api/v1/users/me | 修改个人资料 |
| GET | /api/v1/users | 用户列表（仅管理员） |
| GET | /api/v1/dashboard/stats | 库存总览统计 |
| GET | /api/v1/warehouses | 仓库分页列表 |
| GET | /api/v1/warehouses/:id | 仓库详情与库位 |
| POST | /api/v1/warehouses | 创建仓库 |
| PUT | /api/v1/warehouses/:id | 更新仓库 |
| POST | /api/v1/warehouses/:id/status | 启用/停用仓库 |
| POST | /api/v1/warehouses/shelves | 创建库位 |
| PUT | /api/v1/warehouses/shelves/:id | 更新库位 |
| GET | /api/v1/categories | 商品分类列表 |
| POST | /api/v1/categories | 创建商品分类 |
| GET | /api/v1/products | 商品分页列表 |
| POST | /api/v1/products | 创建商品 |
| GET | /api/v1/products/:id | 商品详情 |
| PUT | /api/v1/products/:id | 更新商品 |
| GET | /api/v1/stock-records | 组合条件查询库存 |
| POST | /api/v1/stock-records/inbound | 入库 |
| POST | /api/v1/stock-records/outbound | 出库 |
| POST | /api/v1/stock-records/adjust | 盘点校准 |
| GET | /api/v1/orders | 单据分页列表 |
| POST | /api/v1/orders | 创建单据（入库/出库/调拨/盘点） |
| GET | /api/v1/orders/:id | 单据详情与明细 |
| POST | /api/v1/orders/:id/submit | 提交单据 |
| POST | /api/v1/orders/:id/approve | 审批单据 |
| POST | /api/v1/orders/:id/execute | 执行单据 |
| POST | /api/v1/orders/:id/cancel | 取消草稿/已提交单据 |
| GET | /api/v1/operation-logs | 操作日志（仅管理员） |
| POST | /api/v1/upload/image | 图片上传 |

## 主要功能

- 库存总览：各仓库库存量柱状图、库存预警商品 TOP10、今日出入库单数。
- 仓库管理：仓库卡片（面积/库存量）、库位地图（网格占用率）、启用/停用。
- 商品管理：树形分类、搜索筛选、各仓库库存分布表、库存阈值预警。
- 出入库管理：创建入库/出库/调拨/盘点单，提交→审批→执行状态流转，动态明细行。
- 盘点管理：按仓库/商品查询账面库存、录入实盘数量、自动计算差异并保存。
- 操作日志：写操作自动记录（管理员查看）。
- 角色权限：JWT + RBAC（admin/warehouse_manager/operator/viewer）。

## License

MIT License
