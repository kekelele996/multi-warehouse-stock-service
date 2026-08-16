// 单号生成工具（与后端逻辑对应：IN/OUT/TR/CK + 年份 + 序号）
export function generateOrderNo(type: string): string {
  const prefixMap: Record<string, string> = {
    inbound: 'IN',
    outbound: 'OUT',
    transfer: 'TR',
    inventory_check: 'CK',
  }
  const prefix = prefixMap[type] || 'IN'
  const now = new Date()
  const rand = Math.floor(Math.random() * 1000000).toString().padStart(6, '0')
  return `${prefix}${now.getFullYear()}${rand}`
}

export function isValidOrderNo(no: string): boolean {
  return /^(IN|OUT|TR|CK)\d{10,}$/.test(no)
}
