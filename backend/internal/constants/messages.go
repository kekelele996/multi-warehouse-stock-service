package constants

// messages.go 同时承载前端提示文案、后端返回文案与日志文案。
const (
	MsgOK                  = "ok"
	MsgUnauthorized        = "未登录或登录已过期"
	MsgForbidden           = "没有权限执行该操作"
	MsgNotFound            = "资源不存在"
	MsgTooManyRequests     = "请求过于频繁，请稍后再试"
	MsgInternalError       = "服务器内部错误"
	MsgPhoneExists         = "手机号已注册"
	MsgInvalidCredentials  = "手机号或密码错误"
	MsgOrderStatusConflict = "单据状态流转冲突"
	MsgInsufficientStock   = "库存不足"
	MsgLoginSuccess        = "登录成功"
	MsgOrderCreated        = "单据创建成功"
	MsgOrderSubmitted      = "单据已提交"
	MsgOrderApproved       = "单据已审批"
	MsgOrderExecuted       = "单据已执行"
	MsgOrderCancelled      = "单据已取消"
	MsgInventorySaved      = "盘点结果已保存"
)
