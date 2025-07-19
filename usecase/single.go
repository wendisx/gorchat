package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"time"

	"github.com/wendisx/gorchat/internal/constant"
	"github.com/wendisx/gorchat/internal/log"
	"github.com/wendisx/gorchat/model"
	"github.com/wendisx/gorchat/repository"
)

type SingleUsecase interface {
	GetLogger() log.Logger
	InviteSingle(singleInvite *model.SingleInvite) (*model.SingleInviter, error)
	AcceptSingle(singleAccept *model.SingleAccept) (*model.SingleInvitee, error)
	UpdateByInviter(singleInviter *model.SingleInviter) error
	UpdateByInvitee(singleInvitee *model.SingleInvitee) error
	GetDetailForInviter(singleInviter *model.SingleInviter) error
	GetDetailForInvitee(singleInvitee *model.SingleInvitee) error
	Delete(singleDelete *model.SingleDelete) error
}

type singleUsecase struct {
	repo   repository.SingleRepository
	logger log.Logger
	c      context.Context
	t      time.Duration
}

func NewSingleUsercase(repo repository.SingleRepository) SingleUsecase {
	return &singleUsecase{
		repo:   repo,
		logger: repo.GetLogger(),
		c:      context.Background(),
		t:      5 * time.Second,
	}
}

func (u *singleUsecase) generateSingleId(inviterId, inviteeId int64) int64 {
	var low, high int64
	if inviterId > inviteeId {
		high, low = inviterId, inviteeId
	} else {
		high, low = inviteeId, inviterId
	}
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[:8], uint64(high))
	binary.BigEndian.PutUint64(buf[8:], uint64(low))
	hashCode := sha256.Sum256(buf)
	singleId := int64(binary.BigEndian.Uint64(hashCode[:8]))
	return singleId & 0x7FFFFFFFFFFFFFFF
}

func (u *singleUsecase) GetLogger() log.Logger {
	return u.logger
}

func (u *singleUsecase) InviteSingle(singleInvite *model.SingleInvite) (*model.SingleInviter, error) {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	// 逻辑未建立
	singleInvite.Deleted = 1
	// 生成singleId
	singleInvite.SingleId = u.generateSingleId(singleInvite.InviterId, singleInvite.InviteeId)
	err := u.repo.InsertUnAccepted(ctx, singleInvite)
	if err != nil {
		return nil, &model.DError{
			Code:    constant.ErrSingleInviteFail,
			Message: constant.MsgSingleInviteFail,
		}
	}
	singleInviter := &model.SingleInviter{
		SingleId:  singleInvite.SingleId,
		InviterId: singleInvite.InviterId,
	}
	err = u.repo.FindByInviter(ctx, singleInviter)
	if err != nil {
		return nil, &model.DError{
			Code:    constant.ErrSingleInviteFail,
			Message: constant.MsgSingleInviteFail,
		}
	}
	return singleInviter, nil
}

func (u *singleUsecase) AcceptSingle(singleAccept *model.SingleAccept) (*model.SingleInvitee, error) {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	singleAccept.Deleted = 0
	err := u.repo.UpdateByAccept(ctx, singleAccept)
	if err != nil {
		return nil, &model.DError{
			Code:    constant.ErrSingleAcceptFail,
			Message: constant.MsgSingleAcceptFail,
		}
	}
	singleInvitee := &model.SingleInvitee{
		SingleId:  singleAccept.SingleId,
		InviteeId: singleAccept.InviteeId,
	}
	err = u.repo.FindByInvitee(ctx, singleInvitee)
	if err != nil {
		return nil, &model.DError{
			Code:    constant.ErrSingleAcceptFail,
			Message: constant.MsgSingleAcceptFail,
		}
	}
	return singleInvitee, nil
}

func (u *singleUsecase) UpdateByInviter(singleInviter *model.SingleInviter) error {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	err := u.repo.UpdateByInviter(ctx, singleInviter)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrSingleUpdateFail,
			Message: constant.MsgSingleUpdateFail,
		}
	}
	err = u.repo.FindByInviter(ctx, singleInviter)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrSingleUpdateFail,
			Message: constant.MsgSingleUpdateFail,
		}
	}
	return nil
}

func (u *singleUsecase) UpdateByInvitee(singleInvitee *model.SingleInvitee) error {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	err := u.repo.UpdateByInvitee(ctx, singleInvitee)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrSingleUpdateFail,
			Message: constant.MsgSingleUpdateFail,
		}
	}
	err = u.repo.FindByInvitee(ctx, singleInvitee)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrSingleUpdateFail,
			Message: constant.MsgSingleUpdateFail,
		}
	}
	return nil
}

func (u *singleUsecase) GetDetailForInviter(singleInviter *model.SingleInviter) error {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	err := u.repo.FindByInviter(ctx, singleInviter)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrSingleGetDetailFail,
			Message: constant.MsgSingleGetDetailFail,
		}
	}
	return nil
}

func (u *singleUsecase) GetDetailForInvitee(singleInvitee *model.SingleInvitee) error {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	err := u.repo.FindByInvitee(ctx, singleInvitee)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrSingleGetDetailFail,
			Message: constant.MsgSingleGetDetailFail,
		}
	}
	return nil
}

func (u *singleUsecase) Delete(singleDelete *model.SingleDelete) error {
	ctx, cancle := context.WithTimeout(u.c, u.t)
	defer cancle()
	err := u.repo.Delete(ctx, singleDelete)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrSingleDeleteFail,
			Message: constant.MsgSingleDeleteFail,
		}
	}
	return nil
}
