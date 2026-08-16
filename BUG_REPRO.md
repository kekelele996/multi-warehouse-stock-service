# BUG_REPRO

## Bug 是什么
调拨入库写成了入到源仓库；调拨执行时取实际数量当调拨数量；库存汇总把 SUM 写成 AVG；库存记录库位默认值写错。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestTransferMovesStockBetweenWarehouses 失败：调拨执行报数量必须大于 0，或目标仓库存为 0、源仓库存不对、汇总不是合计。
