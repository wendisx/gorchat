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

type GroupHandler interface {
	CreateGroup(e echo.Context) error
	JoinGroup(e echo.Context) error
	UpdateGroup(e echo.Context) error
	UpdateGroupUser(e echo.Context) error
	SearchGroup(e echo.Context) error
	SearchGroupUsers(e echo.Context) error
	GetGroupUsers(e echo.Context) error
	DeleteGroup(e echo.Context) error
	DeleteGroupUser(e echo.Context) error
}

type groupHandler struct {
	ucase  usecase.GroupUsecase
	logger log.Logger
	res    model.Response
}

func NewGroupHandler(ucase usecase.GroupUsecase, res model.Response) GroupHandler {
	return &groupHandler{
		ucase:  ucase,
		logger: ucase.GetLogger(),
		res:    res,
	}
}

func (h *groupHandler) CreateGroup(e echo.Context) error {
	createGroupReq, ok := e.Get("body").(*model.CreateGroupReq)
	if !ok {
		return h.res.Fail(e, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	groupBasic := &model.GroupBasic{
		GroupName:     createGroupReq.GroupName,
		GroupNickname: createGroupReq.GroupName,
		GroupPassword: createGroupReq.GroupPassword,
		GroupMaxSize:  createGroupReq.GroupMaxSize,
		UserId:        createGroupReq.UserId,
		UserNickname:  createGroupReq.UserNickname,
	}
	err := h.ucase.GroupCreate(groupBasic)
	if err != nil {
		return err
	}
	createGroupRes := &model.CreateGroupRes{
		GroupId:       groupBasic.GroupId,
		GroupName:     groupBasic.GroupName,
		GroupNickname: groupBasic.GroupNickname,
		GroupAvatar:   groupBasic.GroupAvatar,
	}
	return h.res.Success(e, http.StatusOK, constant.MsgGroupCreateSuccess, createGroupRes)
}

func (h *groupHandler) JoinGroup(e echo.Context) error {
	joinGroupReq, ok := e.Get("body").(*model.JoinGroupReq)
	if !ok {
		return h.res.Fail(e, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	groupToUser := &model.GroupToUser{
		GroupId:      joinGroupReq.GroupId,
		UserId:       joinGroupReq.UserId,
		UserNickname: joinGroupReq.UserNickname,
		UserDisturb:  joinGroupReq.UserDisturb,
		UserRoleId:   3,
	}
	err := h.ucase.GroupJoin(groupToUser)
	if err != nil {
		return err
	}
	joinGroupRes := &model.JoinGroupRes{
		GroupId:      groupToUser.GroupId,
		UserNickname: groupToUser.UserNickname,
		UserDisturb:  groupToUser.UserDisturb,
	}
	return h.res.Success(e, http.StatusOK, constant.MsgGroupJoinSuccess, joinGroupRes)
}

func (h *groupHandler) UpdateGroup(e echo.Context) error {
	updateGroupReq, ok := e.Get("body").(*model.UpdateGroupReq)
	if !ok {
		return h.res.Fail(e, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	group := &model.Group{
		GroupId:       updateGroupReq.GroupId,
		GroupName:     updateGroupReq.GroupName,
		GroupPassword: updateGroupReq.GroupPassword,
		GroupMaxSize:  updateGroupReq.GroupMaxSize,
		GroupAvatar:   updateGroupReq.GroupAvatar,
	}
	err := h.ucase.GroupUpdate(group)
	if err != nil {
		return err
	}
	return h.res.Success(e, http.StatusOK, constant.MsgGroupUpdateSuccess, nil)
}

func (h *groupHandler) UpdateGroupUser(e echo.Context) error {
	updateGroupUserReq, ok := e.Get("body").(*model.UpdateGroupUserReq)
	if !ok {
		return h.res.Fail(e, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	groupToUser := &model.GroupToUser{
		GroupId:            updateGroupUserReq.GroupId,
		GroupNickname:      updateGroupUserReq.SetGroupNickname,
		UserId:             updateGroupUserReq.SetUserId,
		UserNickname:       updateGroupUserReq.SetUserNickname,
		UserRoleId:         updateGroupUserReq.SetUserRole,
		UserRoleNickname:   updateGroupUserReq.SetUserRoleNickname,
		UserDisturb:        updateGroupUserReq.SetUserDisturb,
		IsSetRole:          updateGroupUserReq.IsSetRole,
		IsSetDisturb:       updateGroupUserReq.IsSetDisturb,
		IsSetUserNickname:  updateGroupUserReq.IsSetUserNickname,
		IsSetGroupNickname: updateGroupUserReq.IsSetGroupNickname,
	}
	err := h.ucase.GroupUpdateUser(groupToUser)
	if err != nil {
		return err
	}
	return h.res.Success(e, http.StatusOK, constant.MsgGroupUpdateUserSuccess, nil)
}

func (h *groupHandler) SearchGroup(e echo.Context) error {
	searchGroupReq, ok := e.Get("body").(*model.SearchGroupReq)
	if !ok {
		return h.res.Fail(e, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	groupItem := &model.GroupItem{
		GroupId:   searchGroupReq.GroupId,
		GroupName: searchGroupReq.GroupName,
	}
	page := &model.Page[*model.GroupItem]{
		CurrentPage: searchGroupReq.CurrentPage,
		PageSize:    searchGroupReq.PageSize,
	}
	err := h.ucase.GroupSearch(groupItem, page)
	if err != nil {
		return err
	}
	return h.res.Success(e, http.StatusOK, constant.MsgGroupSearchSuccess, page)
}

func (h *groupHandler) SearchGroupUsers(e echo.Context) error {
	searchGroupUserReq, ok := e.Get("body").(*model.SearchGroupUsersReq)
	if !ok {
		return h.res.Fail(e, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	groupUser := &model.GroupUser{
		GroupId:      searchGroupUserReq.GroupId,
		UserId:       searchGroupUserReq.UserId,
		UserName:     searchGroupUserReq.UserName,
		UserNickname: searchGroupUserReq.UserNickname,
	}
	page := &model.Page[*model.GroupToUserItem]{
		CurrentPage: searchGroupUserReq.CurrentPage,
		PageSize:    searchGroupUserReq.PageSize,
	}
	err := h.ucase.GroupSearchUser(groupUser, page)
	if err != nil {
		return err
	}
	return h.res.Success(e, http.StatusOK, constant.MsgGroupSearchUserSuccess, page)
}

func (h *groupHandler) GetGroupUsers(e echo.Context) error {
	getGroupUsers, ok := e.Get("body").(*model.GetGroupUsersReq)
	if !ok {
		return h.res.Fail(e, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	page := &model.Page[*model.GroupToUserItem]{
		CurrentPage: getGroupUsers.CurrentPage,
		PageSize:    getGroupUsers.PageSize,
	}
	err := h.ucase.GroupAllUsers(getGroupUsers.GroupId, page)
	if err != nil {
		return err
	}
	return h.res.Success(e, http.StatusOK, constant.MsgGroupGetAllUsersSuccess, page)
}

func (h *groupHandler) DeleteGroup(e echo.Context) error {
	groupIdStr := e.QueryParam("groupId")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		return h.res.Fail(e, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	err = h.ucase.GroupDelete(int64(groupId))
	if err != nil {
		return err
	}
	return h.res.Success(e, http.StatusOK, constant.MsgGroupDeleteSuccess, nil)
}

func (h *groupHandler) DeleteGroupUser(e echo.Context) error {
	deleteGroupUserReq, ok := e.Get("body").(*model.DeleteGroupUserReq)
	if !ok {
		return h.res.Fail(e, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	err := h.ucase.GroupDeleteUser(deleteGroupUserReq.GroupId, deleteGroupUserReq.UserId)
	if err != nil {
		return err
	}
	return h.res.Success(e, http.StatusOK, constant.MsgGroupDeleteUserSuccess, nil)
}
