package constant

// 错误代码
const (
	// 服务器内部特殊错误
	ErrServerInternal  int16 = iota // 服务器内部错误
	ErrUnknown                      // 未知错误
	ErrValidate                     // 校验错误
	ErrArgument                     // 参数错误
	ErrOperationFail                // 操作失败
	ErrNotAuthenticate              // 鉴权失败
	ErrBadRequest                   // 错误请求
	ErrBindObjectErr                // 绑定对象错误

	// 数据库操作错误
	ErrTransactionBegin // 事务启动错误
	ErrTransactionFail  // 事务失败
	ErrSqlExcution      // sql执行异常
	ErrSqlInsertFail    // 插入错误
	ErrSqlUpdateFail    // 更新错误
	ErrSqlDeleteFail    // 删除错误
	ErrSqlSelectFail    // 查询错误

	// 逻辑错误
	ErrSignupFail       // 注册失败
	ErrLoginFail        // 登录失败
	ErrUserExist        // 用户已存在
	ErrUserNotExist     // 用户不存在
	ErrPasswordAuthFail // 密码验证失败
	ErrUserUpdateFail   // 更新失败
	ErrUserDeleteFail   // 删除失败
	ErrGetUserDetail    // 获取用户详细信息错误
	ErrSearchUser       // 搜索错误
)

// 错误信息
const (
	// 服务器内部错误
	MsgServerInternalErr = "服务器内部错误"
	MsgUnknownErr        = "未知错误"
	MsgValidateFail      = "校验失败"
	MsgArgumentErr       = "参数错误"
	MsgOperationFail     = "操作失败"
	MsgNotAuthenticate   = "鉴权失败"
	MsgBadRequest        = "错误请求"
	MsgBindObjectErr     = "绑定对象错误"
	// 数据库错误
	MsgTransactionBegin = "事务启动错误"
	MsgTransactionFail  = "事务失败"
	MsgSqlExcution      = "sql执行异常"
	MsgSqlInsertFail    = "sql插入错误"
	MsgSqlUpdateFail    = "sql更新错误"
	MsgSqlDeleteFail    = "sql删除错误"
	MsgSqlSelectFail    = "sql查询错误"
	// 逻辑错误
	MsgSignupFail        = "注册失败"
	MsgLoginFail         = "登录失败"
	MsgUserExist         = "用户已存在"
	MsgUserNotExist      = "用户不存在"
	MsgPasswordAuthFail  = "密码验证失败"
	MsgUserUpdateFail    = "用户更新失败"
	MsgUserDeleteFail    = "用户删除失败"
	MsgGetUserDetailFail = "获取用户详细信息失败"
	MsgSearchUserFail    = "搜索用户失败"
)

// 一般提示信息
const (
	MsgValidateSuccess = "校验成功"

	MsgUserSignupSuccess    = "用户注册成功"
	MsgUserLoginSuccess     = "用户登录成功"
	MsgUserUpdateSuccess    = "用户信息更新成功"
	MsgUserDeleteSuccess    = "用户删除成功"
	MsgGetUserDetailSuccess = "获取用户详细信息成功"
	MsgSearchUserSuccess    = "搜索用户成功"
)
