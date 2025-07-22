package usecase

import (
	"context"
	"time"

	"github.com/wendisx/gorchat/internal/constant"
	"github.com/wendisx/gorchat/internal/log"
	"github.com/wendisx/gorchat/model"
	"github.com/wendisx/gorchat/repository"
)

const (
	ROLE_OWNER  = "owner"
	ROLE_ADMIN  = "admin"
	ROLE_JOINER = "joiner"
)

type GroupUsecase interface {
	GetLogger() log.Logger
	GroupCreate(groupBasic *model.GroupBasic) error
	GroupJoin(groupToUser *model.GroupToUser) error
	GroupUpdate(group *model.Group) error
	GroupUpdateUser(groupToUser *model.GroupToUser) error
	GroupDelete(groupId int64) error
	GroupDeleteUser(groupId, userId int64) error
	GroupUserDetail(groupToUser *model.GroupToUser) error
	GroupSearchUser(groupUser *model.GroupUser, page *model.Page[*model.GroupToUserItem]) error
	GroupSearch(groupItem *model.GroupItem, page *model.Page[*model.GroupItem]) error
	GroupAllUsers(groupId int64, page *model.Page[*model.GroupToUserItem]) error
}

type groupUsecase struct {
	repo   repository.GroupRepository
	logger log.Logger
	c      context.Context
	t      time.Duration
}

func NewGroupUsecase(repo repository.GroupRepository) GroupUsecase {
	return &groupUsecase{
		repo:   repo,
		logger: repo.GetLogger(),
		c:      context.Background(),
		t:      5 * time.Second,
	}
}

func (u *groupUsecase) GetLogger() log.Logger {
	return u.logger
}

func (u *groupUsecase) GroupCreate(groupBasic *model.GroupBasic) error {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	group := &model.Group{
		GroupName:     groupBasic.GroupName,
		GroupPassword: groupBasic.GroupPassword,
		GroupMaxSize:  groupBasic.GroupMaxSize,
	}
	groupId, err := u.repo.InsertOneGroup(ctx, group)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrGroupCreateFail,
			Message: constant.MsgGroupCreateFail,
		}
	}
	userRoleId, err := u.repo.FindUserRoleId(ctx, ROLE_OWNER)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrGroupCreateFail,
			Message: constant.MsgGroupCreateFail,
		}
	}
	groupToUser := &model.GroupToUser{
		GroupId:       groupId,
		GroupNickname: groupBasic.GroupNickname,
		UserId:        groupBasic.UserId,
		UserNickname:  groupBasic.UserNickname,
		UserRoleId:    userRoleId,
		UserDisturb:   1,
	}
	err = u.repo.InsertUserInGroup(ctx, groupToUser)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrGroupCreateFail,
			Message: constant.MsgGroupCreateFail,
		}
	}
	groupBasic.GroupId = groupId
	err = u.repo.FindGroupBasic(ctx, groupBasic)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrGroupCreateFail,
			Message: constant.MsgGroupCreateFail,
		}
	}
	return nil
}

func (u *groupUsecase) GroupJoin(groupToUser *model.GroupToUser) error {
	ctx, cancel := context.WithTimeout(u.c, u.t)
	defer cancel()
	err := u.repo.InsertUserInGroup(ctx, groupToUser)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrGroupJoinFail,
			Message: constant.MsgGroupJoinFail,
		}
	}
	return nil
}

func (u *groupUsecase) GroupUpdate(group *model.Group) error {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	err := u.repo.UpdateGroup(ctx, group)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrGroupUpdateFail,
			Message: constant.MsgGroupUpdateFail,
		}
	}
	return nil
}

func (u *groupUsecase) GroupUpdateUser(groupToUser *model.GroupToUser) error {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	tmpGroupToUser := &model.GroupToUser{
		GroupId: groupToUser.GroupId,
		UserId:  groupToUser.UserId,
	}
	err := u.repo.FindGroupToUser(ctx, tmpGroupToUser)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrGroupUpdateUserFail,
			Message: constant.MsgGroupUpdateUserFail,
		}
	}
	if groupToUser.IsSetDisturb {
		tmpGroupToUser.UserDisturb = groupToUser.UserDisturb
	} else if groupToUser.IsSetRole {
		tmpGroupToUser.UserRoleId = groupToUser.UserRoleId
		tmpGroupToUser.UserRoleNickname = groupToUser.UserRoleNickname
	} else if groupToUser.IsSetGroupNickname {
		tmpGroupToUser.GroupNickname = groupToUser.GroupNickname
	} else if groupToUser.IsSetUserNickname {
		tmpGroupToUser.UserNickname = groupToUser.UserNickname
	}
	if groupToUser.IsSetDisturb {
	}
	err = u.repo.UpdateGroupToUser(ctx, tmpGroupToUser)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrGroupUpdateUserFail,
			Message: constant.MsgGroupUpdateUserFail,
		}
	}
	err = u.repo.FindGroupToUser(ctx, groupToUser)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrGroupUpdateUserFail,
			Message: constant.MsgGroupUpdateUserFail,
		}
	}
	return nil
}

func (u *groupUsecase) GroupDelete(groupId int64) error {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	err := u.repo.DeleteGroup(ctx, groupId)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrGroupDeleteFail,
			Message: constant.MsgGroupDeleteFail,
		}
	}
	return nil
}

func (u *groupUsecase) GroupDeleteUser(groupId, userId int64) error {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	err := u.repo.DeleteGroupToUser(ctx, groupId, userId)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrGroupDeleteUserFail,
			Message: constant.MsgGroupDeleteUserFail,
		}
	}
	return nil
}

func (u *groupUsecase) GroupUserDetail(groupToUser *model.GroupToUser) error {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	err := u.repo.FindGroupToUser(ctx, groupToUser)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrGroupGetDetailFail,
			Message: constant.MsgGroupGetDetailFail,
		}
	}
	return nil
}

func (u *groupUsecase) GroupSearchUser(groupUser *model.GroupUser, page *model.Page[*model.GroupToUserItem]) error {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	err := u.repo.FindGroupUsers(ctx, groupUser, page)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrGroupSearchUserFail,
			Message: constant.MsgGroupSearchUserFail,
		}
	}
	return nil
}

func (u *groupUsecase) GroupSearch(groupItem *model.GroupItem, page *model.Page[*model.GroupItem]) error {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	err := u.repo.FindGroups(ctx, groupItem, page)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrGroupSearchFail,
			Message: constant.MsgGroupSearchFail,
		}
	}
	return nil
}

func (u *groupUsecase) GroupAllUsers(groupId int64, page *model.Page[*model.GroupToUserItem]) error {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	err := u.repo.FindGroupAllUsers(ctx, groupId, page)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrGroupGetAllUsersFail,
			Message: constant.MsgGroupGetAllUsersFail,
		}
	}
	return nil
}
