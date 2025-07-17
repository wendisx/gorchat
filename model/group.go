package model

// entity for group and group_detail table
type Group struct {
	GroupId          int64  `json:"groupId"`          // 群id
	GroupName        string `json:"groupName"`        // 群名称
	GroupPassword    string `json:"groupPassword"`    // 群密码
	GroupAvatar      string `json:"groupAvatar"`      // 群头像
	GroupMaxSize     int    `json:"groupMaxSize"`     // 群最大容量
	GroupCurrentSize int    `json:"groupCurrentSize"` // 群当前容量
	GroupOnlineSize  int    `json:"groupOnlineSize"`  // 群在线人数
	Deleted          int    `json:"deleted"`          // 群逻辑删除
	CreaterId        int64  `json:"createId"`         // 群创建者id
	CreaterName      string `json:"createrName"`      // 群创建者名称
}

type GroupToUser struct {
	GroupId          int64  `json:"groupId"`          // 群号
	UserId           int64  `json:"userId"`           // 用户账号
	GroupNickname    string `json:"groupNickname"`    // 群别称
	UserNickname     string `json:"userNickname"`     // 用户别称
	UserRole         string `json:"userRole"`         // 用户职责
	UserRoleNickname string `json:"userRoleNickname"` // 用户职责别称
	UserDisturb      int    `json:"userDisturb"`      // 用户打扰模式
}
