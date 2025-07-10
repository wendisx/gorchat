package constant

// 错误码
const (
	ErrServerInternal int16 = iota
	ErrUnknown
	ErrValidator

	ErrSqlExcution
	ErrInsert
	ErrUpdate
	ErrDelete
	ErrFind

	ErrSignupFail
	ErrUserExist
	ErrUserNotExist
	ErrPasswordAuth
	ErrUpdateFail
	ErrDeleteFail
)

// 错误信息
const (
	MsgServerInternal   = "服务器内部错误"
	MsgOperationFail    = "操作失败"
	MsgUnAuthorized     = "身份验证失败"
	MsgValidatorFail    = "参数校验失败"
	MsgValidatorSuccess = "参数校验成功"
	MsgSignupFail       = "注册失败"
	MsgUserExist        = "用户存在"
	MsgUserNotExist     = "用户不存在"
	MsgPasswordAuth     = "密码校验失败"
	MsgUpdateFail       = "更新失败"
	MsgDeleteFail       = "删除失败"
)
