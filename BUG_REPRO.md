# BUG_REPRO

## Bug 是什么
错误包装用 %v 丢掉哨兵；用户按手机号查询未命中返回新错误；单据查询未命中返回新错误；登录时丢掉了未命中转 401 的分支。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestErrorChainSentinelsAndLogin 失败：新手机号注册报错、未知手机号登录未返回 invalid credentials、单据未命中不再匹配 ErrNotFound。
