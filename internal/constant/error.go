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

	ErrSingleInviteFail    // 单聊邀请失败
	ErrSingleAcceptFail    // 单聊接受失败
	ErrSingleUpdateFail    // 单聊更新失败
	ErrSingleGetDetailFail // 单聊详细信息获取失败
	ErrSingleDeleteFail    // 单聊删除失败

	ErrGroupCreateFail      // 群聊创建失败
	ErrGroupJoinFail        // 群聊加入失败
	ErrGroupUpdateFail      // 群聊更新失败
	ErrGroupUpdateUserFail  // 群聊设置失败
	ErrGroupDeleteFail      // 群聊删除失败
	ErrGroupDeleteUserFail  // 删除用户失败
	ErrGroupGetDetailFail   // 用户群设置获取失败
	ErrGroupSearchUserFail  // 搜索用户失败
	ErrGroupSearchFail      // 群搜索失败
	ErrGroupGetAllUsersFail // 群获取用户失败
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

	MsgSingleInviteFail    = "单聊邀请失败"
	MsgSingleAcceptFail    = "单聊接受失败"
	MsgSingleUpdateFail    = "单聊更新失败"
	MsgSingleGetDetailFail = "单聊详细信息获取失败"
	MsgSingleDeleteFail    = "单聊删除失败"

	MsgGroupCreateFail      = "群聊创建失败"
	MsgGroupJoinFail        = "群聊加入失败"
	MsgGroupUpdateFail      = "群聊更新失败"
	MsgGroupUpdateUserFail  = "群聊设置失败"
	MsgGroupDeleteFail      = "群聊删除失败"
	MsgGroupDeleteUserFail  = "删除用户失败"
	MsgGroupGetDetailFail   = "用户群设置获取失败"
	MsgGroupSearchUserFail  = "搜索用户失败"
	MsgGroupSearchFail      = "群搜索失败"
	MsgGroupGetAllUsersFail = "群用户返回失败"
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

	MsgSingleInviteSuccess    = "单聊邀请成功"
	MsgSingleAcceptSuccess    = "单聊接受成功"
	MsgSingleUpdateSuccess    = "单聊更新成功"
	MsgSingleGetDetailSuccess = "单聊详细信息获取成功"
	MsgSingleDeleteSuccess    = "单聊删除成功"

	MsgGroupCreateSuccess      = "群聊创建成功"
	MsgGroupJoinSuccess        = "群聊加入成功"
	MsgGroupUpdateSuccess      = "群聊更新成功"
	MsgGroupUpdateUserSuccess  = "群聊设置成功"
	MsgGroupDeleteSuccess      = "群聊删除成功"
	MsgGroupDeleteUserSuccess  = "删除用户成功"
	MsgGroupGetDetailSuccess   = "用户群设置获取成功"
	MsgGroupSearchUserSuccess  = "搜索用户成功"
	MsgGroupSearchSuccess      = "群搜索成功"
	MsgGroupGetAllUsersSuccess = "群用户返回成功"
)
