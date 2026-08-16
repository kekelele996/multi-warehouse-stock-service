# BUG_REPRO

## Bug 是什么
出库扣减把批次数量覆盖写而不是相减；可用量判断用 <= 把刚好够误判为不足；查询批次时按 id 倒序导致先扣新批次；单据执行时把实际数量当作出库数量。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestOutboundOrderFIFOAndExactStock 失败：出库单执行报数量必须大于 0 / 库存不足，或扣减后总量与老批次数量不对。
