# wje-135 多仓库库存管理系统 执行记录

- 项目编号/名称：wje-135 多仓库库存管理系统（WarehouseStock）
- 日期：2026-08-16
- 短名：warehouse-stock
- 端口：前端 18705 / 后端 19205 / MySQL 33308（本机 3306 被 OrbStack 占用，.env 覆盖 DB_PORT=33308；compose 默认 3306）
- 技术栈：Vue 3 + TypeScript + Element Plus + ECharts + Pinia + Vite；Go 1.22 + Gin + GORM；MySQL 8.0；JWT + RBAC

## Docker Compose 结果

| 容器 | 状态 | 端口 |
| --- | --- | --- |
| warehouse-stock-db | Up (healthy) | 33308:3306 |
| warehouse-stock-backend | Up (healthy) | 19205:8080 |
| warehouse-stock-frontend | Up | 18705:80 |

`docker compose config --quiet` 通过；`docker compose up -d --build` 一键启动成功。

## 关键 API 冒烟结果（47/47 通过）

| 接口 | 方法 | 状态码 | 结果摘要 |
| --- | --- | --- | --- |
| /healthz | GET | 200 | ok |
| /api/healthz（Nginx 反代） | GET | 200 | ok |
| /api/v1/healthz（Nginx 反代） | GET | 200 | ok |
| /api/v1/auth/register | POST | 200 | 注册成功 |
| /api/v1/auth/login（仓库经理） | POST | 200 | 返回 JWT + 用户 |
| /api/v1/users/me | GET | 200 | 当前用户 |
| /api/v1/users/me（未授权） | GET | 401 | 未登录 |
| /api/v1/warehouses（观察员创建） | POST | 403 | 越权拦截 |
| /api/v1/dashboard/stats | GET | 200 | 仓库汇总/预警/今日单数 |
| /api/v1/warehouses?page_size=5 | GET | 200 | 仓库列表 |
| /api/v1/warehouses/1 | GET | 200 | 仓库详情+库位 |
| /api/v1/warehouses | POST | 200 | 创建仓库 |
| /api/v1/warehouses/:id/status | POST | 200 | 启用/停用 |
| /api/v1/warehouses/shelves | POST | 200 | 创建库位 |
| /api/v1/categories | GET | 200 | 分类列表 |
| /api/v1/categories | POST | 200 | 创建分类 |
| /api/v1/products?page_size=5 | GET | 200 | 商品列表 |
| /api/v1/products | POST | 200 | 创建商品 |
| /api/v1/products/:id | PUT | 200 | 更新商品 |
| /api/v1/stock-records?page_size=5 | GET | 200 | 库存查询 |
| /api/v1/stock-records/inbound | POST | 200 | 入库 |
| /api/v1/stock-records/outbound | POST | 200 | 出库 |
| /api/v1/stock-records/outbound（超量） | POST | 409 | 库存不足 |
| /api/v1/stock-records/adjust | POST | 200 | 盘点校准 |
| /api/v1/orders?page_size=5 | GET | 200 | 单据列表 |
| /api/v1/orders | POST | 200 | 创建入库单 |
| /api/v1/orders/:id | GET | 200 | 单据详情+明细 |
| /api/v1/orders/:id/submit | POST | 200 | 提交 |
| /api/v1/orders/:id/submit（重复） | POST | 409 | 状态冲突 |
| /api/v1/orders/:id/approve | POST | 200 | 审批 |
| /api/v1/orders/:id/execute | POST | 200 | 执行入库 |
| /api/v1/orders/999999/execute | POST | 404 | 资源不存在 |
| /api/v1/orders（调拨单） | POST | 200 | 创建调拨单 |
| /api/v1/orders/:id/submit（调拨） | POST | 200 | 提交 |
| /api/v1/orders/:id/approve（调拨） | POST | 200 | 审批 |
| /api/v1/orders/:id/execute（调拨） | POST | 200 | 执行调拨 |
| /api/v1/orders（盘点单） | POST | 200 | 创建盘点单 |
| /api/v1/orders/:id/execute（盘点） | POST | 200 | 执行盘点 |
| /api/v1/orders/:id/cancel | POST | 200 | 取消草稿单 |
| /api/v1/orders/:id/cancel（重复） | POST | 409 | 状态冲突 |
| /api/v1/auth/login（admin） | POST | 200 | 管理员登录 |
| /api/v1/operation-logs（admin） | GET | 200 | 操作日志 |
| /api/v1/operation-logs（经理） | GET | 403 | 越权拦截 |
| /api/v1/upload/image | POST | 200 | 图片上传返回 url |

异常路径覆盖：401 未授权、403 越权（仓库创建/操作日志）、409 库存不足与单据状态冲突、404 资源不存在。

## 浏览器验证结论（内置 Playwright，无外部 Chrome）

- /login 登录页打开正常，输入 13800000002/User@123 登录后进入仪表盘，头部显示「仓库经理」与全部菜单（库存总览/仓库管理/商品管理/出入库管理/盘点管理/操作日志/个人中心）。
- /dashboard 渲染统计卡片（仓库数量/今日出入库单/库存预警商品）与「库存预警 TOP10」表格，柱状图区域正常。
- /warehouses 渲染真实仓库卡片（华东一号仓 WH-001、华南二号仓 WH-002）与「查看库位」。
- /orders 渲染真实单据（OUT202608160001 出库、IN 入库已完成、冒烟入库）与入库/出库 Tab、创建单据按钮。
- /products 渲染真实商品（电阻 10KΩ ELEC-R10K、电容 100uF、气泡膜）。
- 截图：output/wje135_dashboard.png、output/wje135_warehouses.png。

## README 检查项

- Docker Compose 一键启动命令在最前；本地开发命令；技术栈表格（后端 Go 1.22 + Gin + GORM）；目录结构；环境变量；部署说明；License。
- 枚举出现位置清单：OrderType、OrderStatus、UserRole 前后端出现位置已列出。

## 其他质量项

- 后端 `go build ./...` 通过；`go vet` 通过。
- 单元测试：internal/service 与 internal/util 表驱动测试通过（go test ./... ok）。
- 前端 `npm run build` 零错误（vue-tsc + vite build 通过）。
- database/init.sql 含建表与种子数据（4 用户 + 2 仓库 + 3 库位 + 3 分类 + 3 商品 + 3 库存记录 + 2 单据 + 2 明细 + 操作日志），容器首次启动自动执行。

- 提交记录：init commit（见 git log）；本文件独立提交 docs commit。
