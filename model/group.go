package model

// entity for group and group_detail table
type Group struct {
	GroupId          int64  `json:"groupId"`          // 群id
	GroupName        string `json:"groupName"`        // 群名称
	GroupPassword    string `json:"groupPassword"`    // 群密码
	GroupAvatar      string `json:"groupAvatar"`      // 群头像
	GroupMaxSize     int    `json:"groupMaxSize"`     // 群最大容量
	GroupCurrentSize int    `json:"groupCurrentSize"` // 群当前容量
	GroupOnlineSize  int    `json:"groupOnlineSize"`  // 群在线人数 -- 去除
	Deleted          int    `json:"deleted"`          // 群逻辑删除
}

type GroupToUser struct {
	GroupId            int64  `json:"groupId"`          // 群号
	GroupAvatar        string `json:"groupAvatar"`      // 群头像
	GroupName          string `json:"groupName"`        // 群名称
	GroupNickname      string `json:"groupNickname"`    // 群别称
	GroupMaxSize       int    `json:"groupMaxSize"`     // 群最大容量
	GroupCurrentSize   int    `json:"groupCurrentSize"` // 群当前容量
	UserId             int64  `json:"userId"`           // 用户账号
	UserName           string `json:"userName"`         // 用户名
	UserNickname       string `json:"userNickname"`     // 用户别称
	UserRoleId         int    `json:"userRoleId"`       // 用户职责id
	UserRole           string `json:"userRole"`         // 用户职责
	UserRoleNickname   string `json:"userRoleNickname"` // 用户职责别称
	UserDisturb        int    `json:"userDisturb"`      // 用户打扰模式
	IsSetRole          bool   `json:"isSetRole"`
	IsSetDisturb       bool   `json:"isSetDisturb"`
	IsSetUserNickname  bool   `json:"isSetUserNickname"`
	IsSetGroupNickname bool   `json:"isSetGroupNickname"`
}

type GroupToUserItem struct {
	UserId           int64  `json:"userId"`
	UserName         string `json:"userName"`
	UserNickname     string `json:"userNickname"`
	UserRole         string `json:"userRole"`
	UserRoleNickname string `json:"userRoleNickname"`
}

type GroupBasic struct {
	GroupId       int64  `json:"groupId"`
	GroupName     string `json:"groupName"`
	GroupNickname string `json:"groupNickname"`
	GroupPassword string `json:"groupPassword"`
	GroupMaxSize  int    `json:"groupMaxsize"`
	GroupAvatar   string `json:"groupAvatar"`
	UserId        int64  `json:"userId"`
	UserNickname  string `json:"userNickname"`
}

type GroupUser struct {
	GroupId      int64  `json:"groupId"`
	UserId       int64  `json:"userId"`
	UserName     string `json:"userName"`
	UserNickname string `json:"userNickname"`
}

type CreateGroupReq struct {
	GroupName     string `json:"groupName"`
	GroupPassword string `json:"groupPassword"`
	GroupMaxSize  int    `json:"groupMaxSize"`
	UserId        int64  `json:"userId"`
	UserNickname  string `json:"userNickname"`
}

type CreateGroupRes struct {
	GroupId       int64  `json:"groupId"`
	GroupName     string `json:"groupName"`
	GroupNickname string `json:"groupNickname"`
	GroupAvatar   string `json:"groupAvatar"`
}

type JoinGroupReq struct {
	GroupId      int64  `json:"groupId"`
	UserId       int64  `json:"userId"`
	UserNickname string `json:"userNickname"`
	UserDisturb  int    `json:"userDisturb"`
}

type JoinGroupRes struct {
	GroupId      int64  `json:"groupId"`
	UserNickname string `json:"userNickname"`
	UserDisturb  int    `json:"userDisturb"`
}

type UpdateGroupReq struct {
	GroupId       int64  `json:"groupId"`
	GroupName     string `json:"groupName"`
	GroupPassword string `json:"groupPassword"`
	GroupMaxSize  int    `json:"groupMaxSize"`
	GroupAvatar   string `json:"groupAvatar"`
}

type UpdateGroupUserReq struct {
	GroupId             int64  `json:"groupId"`
	SetGroupNickname    string `json:"setGroupNickname"`
	SetUserId           int64  `json:"setUserId"`
	SetUserNickname     string `json:"setUserNickname"`
	SetUserRole         int    `json:"setUserRole"`
	SetUserRoleNickname string `json:"setUserRoleNickname"`
	SetUserDisturb      int    `json:"setUserDisturb"`
	IsSetRole           bool   `json:"isSetRole"`
	IsSetDisturb        bool   `json:"isSetDisturb"`
	IsSetUserNickname   bool   `json:"isSetUserNickname"`
	IsSetGroupNickname  bool   `json:"isSetGroupNickname"`
}

type SearchGroupReq struct {
	CurrentPage int    `json:"currentPage"`
	PageSize    int    `json:"pageSize"`
	GroupId     int64  `json:"groupId"`
	GroupName   string `json:"groupName"`
}

type GroupItem struct {
	GroupId   int64  `json:"groupId"`
	GroupName string `json:"groupName"`
}

type SearchGroupUsersReq struct {
	GroupId      int64  `json:"groupId"`
	CurrentPage  int    `json:"currentPage"`
	PageSize     int    `json:"pageSize"`
	UserId       int64  `json:"userId"`
	UserName     string `json:"userName"`
	UserNickname string `json:"userNickname"`
}

type GetGroupUsersReq struct {
	GroupId     int64 `json:"groupId"`
	CurrentPage int   `json:"currentPage"`
	PageSize    int   `json:"pageSize"`
}

type DeleteGroupUserReq struct {
	GroupId int64 `json:"groupId"`
	UserId  int64 `json:"userId"`
}
