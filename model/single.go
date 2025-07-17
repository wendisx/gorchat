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
