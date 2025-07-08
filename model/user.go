package model

import "time"

type User struct {
	UserId     int64     `json:"userId"`
	UserName   string    `json:"userName"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	Deleted    int64     `json:"deleted"`
}

type UserDTO struct {
	UserId       int64  `json:"userId"`
	UserName     string `json:"userName" valid:"min=6"`
	UserPassword string `json:"userPassword" valid:"required,min=8,max=20"`
	UserEmail    string `json:"userEmail" valid:"required,email"`
}

type UserVO struct {
	UserId    int64  `json:"userId"`
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
}
