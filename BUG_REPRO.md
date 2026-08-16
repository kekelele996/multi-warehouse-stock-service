# BUG_REPRO

## Bug 是什么
出库单和仓库查询在未命中时返回 nil,nil，上层拿到 nil 后直接解引用，导致空指针 panic。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestMissingOrderAndWarehouseReturnNotFound 失败：提交不存在单据时发生 invalid memory address or nil pointer dereference。
