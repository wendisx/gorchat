package handler

import (
	"net/http"
	"regexp"
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
	GetUserdetail(c echo.Context) error
	SearchUser(c echo.Context) error
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
	signupReq := c.Get("body").(*model.SignupReq)
	userId, err := h.ucase.Signup(signupReq.UserName, signupReq.UserPassword)
	if err != nil {
		// 可以考虑更高层次的错误封装
		if derr, ok := err.(*model.DError); ok {
			return h.res.Fail(c, http.StatusBadRequest, int(derr.Code), derr.Message)
		} else {
			return err
		}
	}
	return h.res.Success(c, http.StatusOK, constant.MsgUserSignupSuccess, model.SignupRes{
		UserId: strconv.FormatInt(userId, 10),
	})
}

func (h *userHandler) Login(c echo.Context) error {
	loginReq := c.Get("body").(*model.LoginReq)
	user, err := h.ucase.Login(loginReq.UserId, loginReq.UserPassword)
	if err != nil && user == nil {
		if derr, ok := err.(*model.DError); ok {
			return h.res.Fail(c, http.StatusNotFound, int(derr.Code), derr.Message)
		} else {
			return err
		}
	}
	loginRes := model.LoginRes{
		UserId:       strconv.FormatInt(user.UserId, 10),
		UserName:     user.UserName,
		UserEmail:    user.UserEmail,
		UserPhone:    user.UserPhone,
		UserGender:   user.UserGender,
		UserAge:      user.UserAge,
		UserAddress:  user.UserAddress,
		UserLocation: user.UserLocation,
		UserAvatar:   user.UserAvatar,
	}
	return h.res.Success(c, http.StatusOK, constant.MsgUserLoginSuccess, loginRes)
}

func (h *userHandler) UpdateInfo(c echo.Context) error {
	updateInfoReq := c.Get("body").(*model.UpdateInfoReq)
	userid, err := strconv.Atoi(updateInfoReq.UserId)
	if err != nil {
		return h.res.Fail(c, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	updateUser := &model.User{
		UserId:       int64(userid),
		UserName:     updateInfoReq.UserName,
		UserEmail:    updateInfoReq.UserEmail,
		UserPhone:    updateInfoReq.UserPhone,
		UserGender:   updateInfoReq.UserGender,
		UserAge:      updateInfoReq.UserAge,
		UserAddress:  updateInfoReq.UserAddress,
		UserLocation: updateInfoReq.UserLocation,
		UserAvatar:   updateInfoReq.UserAvatar,
	}
	user, err := h.ucase.UpdateInfo(updateUser)
	if err != nil || user == nil {
		if derr, ok := err.(*model.DError); ok {
			return h.res.Fail(c, http.StatusInternalServerError, int(derr.Code), derr.Message)
		}
	}
	updateInfoRes := model.UpdateInfoRes{
		UserName:     user.UserName,
		UserEmail:    user.UserEmail,
		UserPhone:    user.UserPhone,
		UserGender:   user.UserGender,
		UserAge:      user.UserAge,
		UserAddress:  user.UserAddress,
		UserLocation: user.UserLocation,
		UserAvatar:   user.UserAvatar,
	}
	return h.res.Success(c, http.StatusOK, constant.MsgUserUpdateSuccess, updateInfoRes)
}

func (h *userHandler) Delete(c echo.Context) error {
	userIdRegex := regexp.MustCompile(`^\d+$`)
	userIdStr := c.QueryParam("userId")
	ok := userIdRegex.MatchString(userIdStr)
	userId, err := strconv.Atoi(userIdStr)
	if userIdStr == "" || !ok || err != nil {
		return h.res.Fail(c, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	err = h.ucase.Delete(int64(userId))
	if err != nil {
		if derr, ok := err.(*model.DError); ok {
			return h.res.Fail(c, http.StatusInternalServerError, int(derr.Code), derr.Message)
		}
	}
	return h.res.Success(c, http.StatusOK, constant.MsgUserDeleteSuccess, nil)
}

func (h *userHandler) GetUserdetail(c echo.Context) error {
	userIdRegex := regexp.MustCompile(`^\d+$`)
	userIdStr := c.QueryParam("userId")
	ok := userIdRegex.MatchString(userIdStr)
	userId, err := strconv.Atoi(userIdStr)
	if userIdStr == "" || !ok || err != nil {
		return h.res.Fail(c, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	user, err := h.ucase.GetUserDetail(int64(userId))
	if err != nil {
		return err
	}
	getUserdetailRes := model.GetUserdetailRes{
		UserId:       strconv.FormatInt(user.UserId, 10),
		UserName:     user.UserName,
		UserEmail:    user.UserEmail,
		UserPhone:    user.UserPhone,
		UserGender:   user.UserGender,
		UserAge:      user.UserAge,
		UserAddress:  user.UserAddress,
		UserLocation: user.UserLocation,
		UserAvatar:   user.UserAvatar,
	}
	return h.res.Success(c, http.StatusOK, constant.MsgGetUserDetailSuccess, getUserdetailRes)
}

func (h *userHandler) SearchUser(c echo.Context) error {
	searchUserReq := c.Get("body").(*model.SearchUserReq)
	var userBasic model.UserBasic
	userBasic.UserName = searchUserReq.UserName
	userIdRegex := regexp.MustCompile(`^\d+$`)
	ok := userIdRegex.MatchString(searchUserReq.UserId)
	if !ok {
		userBasic.UserId = 0
	} else {
		userId, err := strconv.Atoi(searchUserReq.UserId)
		if err != nil {
			return h.res.Fail(c, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
		}
		userBasic.UserId = int64(userId)
	}
	page := &model.Page[model.UserBasic]{
		CurrentPage: searchUserReq.CurrentPage,
		PageSize:    searchUserReq.PageSize,
	}
	err := h.ucase.SearchUsers(userBasic.UserId, userBasic.UserName, page)
	if err != nil {
		return err
	}
	searchUserRes := model.SearchUserRes{
		CurrentPage: page.CurrentPage,
		PageSize:    page.PageSize,
		Total:       page.Total,
		Items:       page.Items,
	}
	return h.res.Success(c, http.StatusOK, constant.MsgSearchUserSuccess, searchUserRes)
}
