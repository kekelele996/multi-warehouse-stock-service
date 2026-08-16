# BUG_REPRO

## Bug 是什么
盘点校准的差异计算方向写反；批次查询按 id 倒序导致补差补到新批次；库存数量默认值写错；盘点单执行时取实际数量当盘点数。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestAdjustSetsStockToActual 失败：盘点后总数不是盘点数、老批次数量不对，盘点单执行后总数也不对。
