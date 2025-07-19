package model

// entity for user and user_detail table
type User struct {
	UserId       int64  `json:"userId"`       // 用户账号
	UserName     string `json:"userName"`     // 用户名
	UserPassword string `json:"userPassword"` // 用户登录密码
	UserEmail    string `json:"userEmail"`    // 用户邮箱
	UserPhone    string `json:"userPhone"`    // 用户手机
	UserGender   string `json:"userGender"`   // 用户性别
	UserAge      int    `json:"userAge"`      // 用户年龄
	UserAddress  string `json:"userAddress"`  // 用户地址
	UserLocation string `json:"userLocation"` // 用户所在地
	UserAvatar   string `json:"userAvatar"`   // 用户头像
	Deleted      int64  `json:"deleted"`      // 用户注销软删除
}

type UserBasic struct {
	UserId   int64  `json:"userId"`
	UserName string `json:"userName"`
}

type SignupReq struct {
	UserName     string `json:"userName" valid:"required,min=1,max=16"`
	UserPassword string `json:"userPassword" valid:"required,min=8,max=20"`
}

type SignupRes struct {
	UserId int64 `json:"userId"`
}

type LoginReq struct {
	UserId       int64  `json:"userId" valid:"required,min=100000"`
	UserPassword string `json:"userPassword" valid:"required,min=8,max=20"`
}

type LoginRes struct {
	UserId       int64  `json:"userId"`
	UserName     string `json:"userName"`
	UserEmail    string `json:"userEmail"`
	UserPhone    string `json:"userPhone"`
	UserGender   string `json:"userGender"`
	UserAge      int    `json:"userAge"`
	UserAddress  string `json:"userAddress"`
	UserLocation string `json:"userLocation"`
	UserAvatar   string `json:"userAvatar"`
}

type UpdateInfoReq struct {
	UserId       int64  `json:"userId" valid:"required,min=100000"`
	UserName     string `json:"userName"`
	UserEmail    string `json:"userEmail"`
	UserPhone    string `json:"userPhone"`
	UserGender   string `json:"userGender"`
	UserAge      int    `json:"userAge"`
	UserAddress  string `json:"userAddress"`
	UserLocation string `json:"userLocation"`
	UserAvatar   string `json:"userAvatar"`
}

type UpdateInfoRes struct {
	UserName     string `json:"userName"`
	UserEmail    string `json:"userEmail"`
	UserPhone    string `json:"userPhone"`
	UserGender   string `json:"userGender"`
	UserAge      int    `json:"userAge"`
	UserAddress  string `json:"userAddress"`
	UserLocation string `json:"userLocation"`
	UserAvatar   string `json:"userAvatar"`
}

type GetUserdetailRes struct {
	UserId       int64  `json:"userId"`
	UserName     string `json:"userName"`
	UserEmail    string `json:"userEmail"`
	UserPhone    string `json:"userPhone"`
	UserGender   string `json:"userGender"`
	UserAge      int    `json:"userAge"`
	UserAddress  string `json:"userAddress"`
	UserLocation string `json:"userLocation"`
	UserAvatar   string `json:"userAvatar"`
}

type SearchUserReq struct {
	CurrentPage int    `json:"currentPage" valid:"required,min=1"`
	PageSize    int    `json:"pageSize" valid:"required,min=1,max=8"`
	UserId      int64  `json:"userId" valid:"required,min=100000"`
	UserName    string `json:"userName"`
}

type SearchUserRes struct {
	CurrentPage int         `json:"currentPage"`
	PageSize    int         `json:"pageSize" `
	Total       int         `json:"total"`
	Items       []UserBasic `json:"items"`
}
