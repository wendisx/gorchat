package model

import "time"

// entity for single_chat table
type Single struct {
	SingleId        int64     `json:"singleId"`        // 单聊id
	InviterId       int64     `json:"inviterId"`       // 邀请人id
	InviteeId       int64     `json:"inviteeId"`       // 被邀请人id
	InviterNickname string    `json:"inviterNickname"` // 邀请人别称
	InviteeNickname string    `json:"inviteeNickname"` // 被邀请人别称
	InviterDisturb  int       `json:"inviterDisturb"`  // 邀请人打扰模式
	InviteeDisturb  int       `json:"inviteeDisturb"`  // 被邀请人打扰模式
	CreateTime      time.Time `json:"createTime"`      // 单聊创建时间
	Deleted         int       `json:"deleted"`         // 逻辑删除
}

type SingleInvite struct {
	SingleId        int64  `json:"singleId"`
	InviterId       int64  `json:"inviterId"`
	InviteeId       int64  `json:"inviteeId"`
	InviteeNickname string `json:"inviteeNickname"`
	InviterDisturb  int    `json:"inviterDisturb"`
	Deleted         int    `json:"deleted"`
}

type SingleAccept struct {
	SingleId        int64  `json:"singleId"`
	InviteeId       int64  `json:"inviteeId"`
	InviterNickname string `json:"inviterNickname"`
	InviteeDisturb  int    `json:"inviteeDisturb"`
	Deleted         int    `json:"deleted"`
}

type SingleInviter struct {
	SingleId        int64  `json:"singleId"`
	InviterId       int64  `json:"inviterId"`
	InviteeId       int64  `json:"inviteeId"`
	InviteeNickname string `json:"inviterNickname"`
	InviteeName     string `json:"inviterName"`
	InviterDisturb  int    `json:"inviteeDisturb"`
}

type SingleInvitee struct {
	SingleId        int64  `json:"singleId"`
	InviteeId       int64  `json:"inviteeId"`
	InviterId       int64  `json:"inviterId"`
	InviterNickname string `json:"inviterNickname"`
	InviterName     string `json:"inviterName"`
	InviteeDisturb  int    `json:"inviteeDisturb"`
}

type SingleDelete struct {
	SingleId  int64 `json:"singleId"`
	InviterId int64 `json:"inviterId"`
	InviteeId int64 `json:"inviteeId"`
}

type InviteReq struct {
	InviterId       int64  `json:"inviterId"`
	InviteeId       int64  `json:"inviteeId"`
	InviteeNickname string `json:"inviteeNickname"`
	InviterDisturb  int    `json:"inviterDisturb"`
}

type InviteRes struct {
	SingleId        int64  `json:"singleId"`
	InviteeId       int64  `json:"inviteeId"`
	InviteeName     string `json:"inviteeName"`
	InviteeNickname string `json:"inviteeNickname"`
	InviterDisturb  int    `json:"inviterDisturb"`
}

type AcceptReq struct {
	SingleId        int64  `json:"singleId"`
	InviteeId       int64  `json:"inviteeId"`
	InviterNickname string `json:"inviterNickname"`
	InviteeDisturb  int    `json:"inviteeDisturb"`
}

type AcceptRes struct {
	SingleId        int64  `json:"singleId"`
	InviterId       int64  `json:"inviterId"`
	InviterName     string `json:"inviterName"`
	InviterNickname string `json:"inviterNickname"`
	InviteeDisturb  int    `json:"inviteeDisturb"`
}

type UpdateNicknameReq struct {
	SingleId    int64  `json:"singleId"`
	IsInviter   bool   `json:"isInviter"`
	UserId      int64  `json:"userId"`
	SetNickname string `json:"setNickname"`
	UserDisturb int    `json:"userDisturb"`
}

type UpdateNicknameRes struct {
	SingleId     int64  `json:"singleId"`
	UserId       int64  `json:"userId"`
	UserName     string `json:"userName"`
	UserNickname string `json:"userNickname"`
}

type UpdateDisturbReq struct {
	SingleId     int64  `json:"singleId"`
	IsInviter    bool   `json:"isInviter"`
	UserId       int64  `json:"userId"`
	SetDisturb   int    `json:"userDisturb"`
	UserNickname string `json:"userNickname"`
}

type UpdateDisturbRes struct {
	SingleId    int64 `json:"singleId"`
	UserDisturb int   `json:"userDisturb"`
}

type GetDetailReq struct {
	SingleId  int64 `json:"singleId"`
	IsInviter bool  `json:"isInviter"`
	UserId    int64 `json:"userId"`
}

type GetDetailRes struct {
	SingleId     int64  `json:"singleId"`
	UserId       int64  `json:"userId"`
	UserName     string `json:"userName"`
	UserNickname string `json:"userNickname"`
	UserDisturb  int    `json:"userDisturb"`
}

type DeleteReq struct {
	SingleId  int64 `json:"singleId"`
	InviterId int64 `json:"inviterId"`
	InviteeId int64 `json:"inviteeId"`
}
