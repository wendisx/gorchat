package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/wendisx/gorchat/internal/constant"
	"github.com/wendisx/gorchat/internal/log"
	"github.com/wendisx/gorchat/model"
	"github.com/wendisx/gorchat/usecase"
)

type UserHandler interface {
	Signup(c echo.Context) error
	Login(c echo.Context) error
	UpdateInfo(c echo.Context) error
	Delete(c echo.Context) error
}

type userHandler struct {
	ucase  usecase.UserUsecase
	logger log.Logger
	res    model.Response
}

func NewUserHandler(ucase usecase.UserUsecase, res model.Response) UserHandler {
	return &userHandler{
		ucase:  ucase,
		logger: ucase.GetLogger(),
		res:    res,
	}
}

func (h *userHandler) Signup(c echo.Context) error {
	// 断言一定成功，逻辑没问题，预期也差不多，但是不建议
	signupInfo := c.Get("body").(*model.UserDTO)
	signupUser := model.User{
		UserName: signupInfo.UserName,
		Password: signupInfo.UserPassword,
		Email:    signupInfo.UserEmail,
	}
	userId, err := h.ucase.Signup(signupUser)
	if err != nil {
		if derr, ok := err.(*model.DError); ok {
			return h.res.Fail(c, http.StatusBadRequest, int(derr.Code), derr.Message)
		} else {
			return err
		}
	}
	return h.res.Success(c, http.StatusOK, constant.MsgSignupSuccess, model.UserVO{
		UserId:    userId,
		UserName:  signupInfo.UserName,
		UserEmail: signupInfo.UserEmail,
	})
}

func (h *userHandler) Login(c echo.Context) error {
	loginInfo := c.Get("body").(*model.UserDTO)
	loginUser := model.User{
		UserId:   loginInfo.UserId,
		UserName: loginInfo.UserName,
		Password: loginInfo.UserPassword,
		Email:    loginInfo.UserEmail,
	}
	user, err := h.ucase.Login(loginUser)
	if err != nil {
		if derr, ok := err.(*model.DError); ok {
			return h.res.Fail(c, http.StatusInternalServerError, int(derr.Code), derr.Message)
		}
	}
	return h.res.Success(c, http.StatusOK, constant.MsgLoginSuccess, model.UserVO{
		UserId:    user.UserId,
		UserName:  user.UserName,
		UserEmail: user.Email,
	})
}

func (h *userHandler) UpdateInfo(c echo.Context) error {
	newInfo := c.Get("body").(*model.UserDTO)
	newUserInfo := model.User{
		UserId:   newInfo.UserId,
		UserName: newInfo.UserName,
		Password: newInfo.UserPassword,
		Email:    newInfo.UserEmail,
	}
	newUser, err := h.ucase.UpdateInfo(newUserInfo)
	if err != nil {
		if derr, ok := err.(*model.DError); ok {
			return h.res.Fail(c, http.StatusInternalServerError, int(derr.Code), derr.Message)
		}
	}
	return h.res.Success(c, http.StatusOK, constant.MsgUpdateInfoSuccess, model.UserVO{
		UserId:    newUser.UserId,
		UserName:  newUser.UserName,
		UserEmail: newUser.Email,
	})
}

func (h *userHandler) Delete(c echo.Context) error {
	userIdStr := c.QueryParam("userId")
	if userIdStr == "" {
		return h.res.Fail(c, http.StatusBadRequest, constant.FAIL_CODE, constant.MsgBadRequest)
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return err
	}
	err = h.ucase.Delete(int64(userId))
	if err != nil {
		if derr, ok := err.(*model.DError); ok {
			return h.res.Fail(c, http.StatusInternalServerError, int(derr.Code), derr.Message)
		}
	}
	return h.res.Success(c, http.StatusOK, constant.MsgDeleteSuccess, nil)
}
