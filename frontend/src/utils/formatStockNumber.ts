export function formatStockNumber(quantity?: number | null, unit = ''): string {
  const q = quantity ?? 0
  const n = Number(q).toLocaleString('zh-CN')
  return unit ? `${n} ${unit}` : n
}
