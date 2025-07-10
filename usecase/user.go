package usecase

import (
	"context"
	"time"

	"github.com/wendisx/gorchat/internal/constant"
	"github.com/wendisx/gorchat/internal/log"
	"github.com/wendisx/gorchat/model"
	"github.com/wendisx/gorchat/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	GetLogger() log.Logger
	Signup(user model.User) (int64, error)
	Login(loginUser model.User) (model.User, error)
	UpdateInfo(newInfo model.User) (model.User, error)
	Delete(userId int64) error
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

func (u *userUsecase) Signup(user model.User) (int64, error) {
	ctx, cancel := context.WithTimeout(u.c, u.t)
	defer cancel()
	// block when exist user
	_, err := u.repo.FindOneByEmail(ctx, user.Email)
	if err == nil {
		return -1, &model.DError{
			Code:    constant.ErrUserExist,
			Message: constant.MsgUserExist,
		}
	}
	// bypass when not exist user
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return -1, &model.DError{
			Code:    constant.ErrSignupFail,
			Message: constant.MsgSignupFail,
		}
	}
	user.Password = string(hashPassword)
	userId, err := u.repo.InsertOne(ctx, user)
	if err != nil {
		return -1, &model.DError{
			Code:    constant.ErrSignupFail,
			Message: constant.MsgSignupFail,
		}
	}
	return userId, nil
}

func (u *userUsecase) Login(loginUser model.User) (model.User, error) {
	ctx, cancel := context.WithTimeout(u.c, u.t)
	defer cancel()
	var user model.User
	var err error
	// 带着 userId 登录意味着刚刚注册存在 signup 返回的 userId
	// 带着 userId 的请求忽略登录时的账号或者邮箱字段,因为 userId 唯一
	if loginUser.UserId > 0 {
		user, err = u.repo.FindOneById(ctx, loginUser.UserId)
	} else {
		user, err = u.repo.FindOneByEmail(ctx, loginUser.Email)
	}
	if err != nil {
		return user, &model.DError{
			Code:    constant.ErrUserNotExist,
			Message: constant.MsgUserNotExist,
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		return user, &model.DError{
			Code:    constant.ErrPasswordAuth,
			Message: constant.MsgPasswordAuth,
		}
	}
	return user, nil
}

func (u *userUsecase) UpdateInfo(newInfo model.User) (model.User, error) {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	var user model.User
	_, err := u.repo.FindOneById(ctx, newInfo.UserId)
	if err != nil {
		return user, &model.DError{
			Code:    constant.ErrUserNotExist,
			Message: constant.MsgUserNotExist,
		}
	}
	user, err = u.repo.UpdateOneById(ctx, newInfo)
	if err != nil {
		return user, &model.DError{
			Code:    constant.ErrUpdateFail,
			Message: constant.MsgUpdateFail,
		}
	}
	return user, nil
}

func (u *userUsecase) Delete(userId int64) error {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	_, err := u.repo.FindOneById(ctx, userId)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrUserNotExist,
			Message: constant.MsgUserNotExist,
		}
	}
	err = u.repo.DeleteById(ctx, userId)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrDeleteFail,
			Message: constant.MsgDeleteFail,
		}
	}
	return nil
}
