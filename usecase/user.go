package usecase

import (
	"context"
	"strconv"
	"time"

	"github.com/wendisx/gorchat/internal/constant"
	"github.com/wendisx/gorchat/internal/log"
	"github.com/wendisx/gorchat/model"
	"github.com/wendisx/gorchat/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	GetLogger() log.Logger
	Signup(userName string, userPassword string) (int64, error)
	Login(userId string, userPassword string) (*model.User, error)
	UpdateInfo(user *model.User) (*model.User, error)
	Delete(userId int64) error
	GetUserDetail(userId int64) (*model.User, error)
	SearchUsers(userId int64, userName string, page *model.Page[model.UserBasic]) error
}

type userUsecase struct {
	repo   repository.UserRepository
	Logger log.Logger
	c      context.Context
	t      time.Duration
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{
		repo:   repo,
		Logger: repo.GetLogger(),
		c:      context.TODO(),
		t:      10 * time.Second,
	}
}

func (u *userUsecase) GetLogger() log.Logger {
	return u.Logger
}

func (u *userUsecase) Signup(userName string, userPassword string) (int64, error) {
	ctx, cancel := context.WithTimeout(u.c, u.t)
	defer cancel()
	// 由于账号才是唯一标识，注册时并不允许直接传递账号，账号尝试自动生成
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		return -1, &model.DError{
			Code:    constant.ErrSignupFail,
			Message: constant.MsgSignupFail,
		}
	}
	var user *model.User
	userPassword = string(hashPassword)
	user, err = u.repo.InsertOne(ctx, &model.User{
		UserName:     userName,
		UserPassword: userPassword,
	})
	if err != nil && user == nil {
		return -1, &model.DError{
			Code:    constant.ErrSignupFail,
			Message: constant.MsgSignupFail,
		}
	}
	return user.UserId, nil
}

func (u *userUsecase) Login(userId string, userPassword string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(u.c, u.t)
	defer cancel()
	// 转换账号类型
	userid, err := strconv.Atoi(userId)
	if err != nil {
		return nil, &model.DError{
			Code:    constant.ErrArgument,
			Message: constant.MsgArgumentErr,
		}
	}
	var user *model.User
	// 带着 userId 登录意味着刚刚注册存在 signup 返回的 userId
	// 带着 userId 的请求忽略登录时的账号或者邮箱字段,因为 userId 唯一
	user, err = u.repo.FindOneById(ctx, int64(userid))
	if err != nil && user == nil {
		return nil, &model.DError{
			Code:    constant.ErrUserNotExist,
			Message: constant.MsgUserNotExist,
		}
	}
	// 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(userPassword))
	// 密码校验失败
	if err != nil {
		return nil, &model.DError{
			Code:    constant.ErrPasswordAuthFail,
			Message: constant.MsgPasswordAuthFail,
		}
	}
	return user, nil
}

func (u *userUsecase) UpdateInfo(user *model.User) (*model.User, error) {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	tuser, err := u.repo.FindOneById(ctx, user.UserId)
	if err != nil || tuser == nil {
		return user, &model.DError{
			Code:    constant.ErrUserNotExist,
			Message: constant.MsgUserNotExist,
		}
	}
	user, err = u.repo.UpdateOneById(ctx, user)
	if err != nil || user == nil {
		return user, &model.DError{
			Code:    constant.ErrUserUpdateFail,
			Message: constant.MsgUserUpdateFail,
		}
	}
	return user, nil
}

func (u *userUsecase) Delete(userId int64) error {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	tuser, err := u.repo.FindOneById(ctx, userId)
	if err != nil && tuser == nil {
		return &model.DError{
			Code:    constant.ErrUserNotExist,
			Message: constant.MsgUserNotExist,
		}
	}
	err = u.repo.DeleteOneById(ctx, userId)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrUserDeleteFail,
			Message: constant.MsgUserDeleteFail,
		}
	}
	return nil
}

func (u *userUsecase) GetUserDetail(userId int64) (*model.User, error) {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	user, err := u.repo.FindOneById(ctx, userId)
	if err != nil && user == nil {
		return nil, &model.DError{
			Code:    constant.ErrGetUserDetail,
			Message: constant.MsgGetUserDetailFail,
		}
	}
	return user, nil
}

func (u *userUsecase) SearchUsers(userId int64, userName string, page *model.Page[model.UserBasic]) error {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	userBasic := model.UserBasic{
		UserId:   userId,
		UserName: userName,
	}
	err := u.repo.FindBasicLists(ctx, userBasic, page)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrSearchUser,
			Message: constant.MsgSearchUserFail,
		}
	}
	return nil
}
