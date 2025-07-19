package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wendisx/gorchat/internal/constant"
	"github.com/wendisx/gorchat/internal/log"
	"github.com/wendisx/gorchat/model"
	"github.com/wendisx/gorchat/usecase"
)

type SingleHandler interface {
	Invite(e echo.Context) error
	Accept(e echo.Context) error
	UpdateNickname(e echo.Context) error
	UpdateDisturb(e echo.Context) error
	GetDetail(e echo.Context) error
	Delete(e echo.Context) error
}

type singleHandler struct {
	ucase  usecase.SingleUsecase
	logger log.Logger
	res    model.Response
}

func NewSingleHandler(ucase usecase.SingleUsecase, res model.Response) SingleHandler {
	return &singleHandler{
		ucase:  ucase,
		logger: ucase.GetLogger(),
		res:    res,
	}
}

func (h *singleHandler) Invite(e echo.Context) error {
	inviteReq, ok := e.Get("body").(*model.InviteReq)
	if !ok {
		return h.res.Fail(e, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	singleInvite := &model.SingleInvite{
		InviterId:       inviteReq.InviterId,
		InviteeId:       inviteReq.InviteeId,
		InviteeNickname: inviteReq.InviteeNickname,
		InviterDisturb:  inviteReq.InviterDisturb,
	}
	singleInviter, err := h.ucase.InviteSingle(singleInvite)
	if err != nil {
		return err
	}
	inviteRes := model.InviteRes{
		SingleId:        singleInvite.SingleId,
		InviteeId:       singleInvite.InviteeId,
		InviteeName:     singleInviter.InviteeName,
		InviteeNickname: singleInvite.InviteeNickname,
		InviterDisturb:  singleInvite.InviterDisturb,
	}
	return h.res.Success(e, http.StatusOK, constant.MsgSingleInviteSuccess, inviteRes)
}

func (h *singleHandler) Accept(e echo.Context) error {
	acceptReq, ok := e.Get("body").(*model.AcceptReq)
	if !ok {
		return h.res.Fail(e, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	singleAccept := &model.SingleAccept{
		SingleId:        acceptReq.SingleId,
		InviteeId:       acceptReq.InviteeId,
		InviterNickname: acceptReq.InviterNickname,
		InviteeDisturb:  acceptReq.InviteeDisturb,
	}
	singleInvitee, err := h.ucase.AcceptSingle(singleAccept)
	if err != nil {
		return err
	}
	acceptRes := &model.AcceptRes{
		SingleId:        singleInvitee.SingleId,
		InviterId:       singleInvitee.InviterId,
		InviterName:     singleInvitee.InviterName,
		InviterNickname: singleInvitee.InviterNickname,
		InviteeDisturb:  singleInvitee.InviteeDisturb,
	}
	return h.res.Success(e, http.StatusOK, constant.MsgSingleAcceptSuccess, acceptRes)
}

func (h *singleHandler) UpdateNickname(e echo.Context) error {
	updateNicknameReq, ok := e.Get("body").(*model.UpdateNicknameReq)
	if !ok {
		return h.res.Fail(e, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	var err error
	if updateNicknameReq.IsInviter {
		inviter := &model.SingleInviter{
			SingleId:        updateNicknameReq.SingleId,
			InviterId:       updateNicknameReq.UserId,
			InviteeNickname: updateNicknameReq.SetNickname,
			InviterDisturb:  updateNicknameReq.UserDisturb,
		}
		err = h.ucase.UpdateByInviter(inviter)
		if err != nil {
			return err
		}
		updateNicknameRes := &model.UpdateNicknameRes{
			SingleId:     inviter.SingleId,
			UserId:       inviter.InviteeId,
			UserName:     inviter.InviteeName,
			UserNickname: inviter.InviteeNickname,
		}
		return h.res.Success(e, http.StatusOK, constant.MsgSingleUpdateSuccess, updateNicknameRes)
	}
	invitee := &model.SingleInvitee{
		SingleId:        updateNicknameReq.SingleId,
		InviteeId:       updateNicknameReq.UserId,
		InviterNickname: updateNicknameReq.SetNickname,
		InviteeDisturb:  updateNicknameReq.UserDisturb,
	}
	err = h.ucase.UpdateByInvitee(invitee)
	if err != nil {
		return err
	}
	updateNicknameRes := &model.UpdateNicknameRes{
		SingleId:     invitee.SingleId,
		UserId:       invitee.InviterId,
		UserName:     invitee.InviterName,
		UserNickname: invitee.InviterNickname,
	}
	return h.res.Success(e, http.StatusOK, constant.MsgSingleUpdateSuccess, updateNicknameRes)
}

func (h *singleHandler) UpdateDisturb(e echo.Context) error {
	updateDisturbReq, ok := e.Get("body").(*model.UpdateDisturbReq)
	if !ok {
		return h.res.Fail(e, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	var err error
	if updateDisturbReq.IsInviter {
		inviter := &model.SingleInviter{
			SingleId:        updateDisturbReq.SingleId,
			InviterId:       updateDisturbReq.UserId,
			InviterDisturb:  updateDisturbReq.SetDisturb,
			InviteeNickname: updateDisturbReq.UserNickname,
		}
		err = h.ucase.UpdateByInviter(inviter)
		if err != nil {
			return err
		}
		updateDisturbRes := &model.UpdateDisturbRes{
			SingleId:    inviter.SingleId,
			UserDisturb: inviter.InviterDisturb,
		}
		return h.res.Success(e, http.StatusOK, constant.MsgSingleUpdateSuccess, updateDisturbRes)
	}
	invitee := &model.SingleInvitee{
		SingleId:        updateDisturbReq.SingleId,
		InviteeId:       updateDisturbReq.UserId,
		InviteeDisturb:  updateDisturbReq.SetDisturb,
		InviterNickname: updateDisturbReq.UserNickname,
	}
	err = h.ucase.UpdateByInvitee(invitee)
	if err != nil {
		return err
	}
	updateDisturbRes := &model.UpdateDisturbRes{
		SingleId:    invitee.SingleId,
		UserDisturb: invitee.InviteeDisturb,
	}
	return h.res.Success(e, http.StatusOK, constant.MsgSingleUpdateSuccess, updateDisturbRes)

}

func (h *singleHandler) GetDetail(e echo.Context) error {
	getDetailReq, ok := e.Get("body").(*model.GetDetailReq)
	if !ok {
		return h.res.Fail(e, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	var err error
	if getDetailReq.IsInviter {
		inviter := &model.SingleInviter{
			SingleId:  getDetailReq.SingleId,
			InviterId: getDetailReq.UserId,
		}
		err = h.ucase.GetDetailForInviter(inviter)
		if err != nil {
			return err
		}
		getDetailRes := &model.GetDetailRes{
			SingleId:     inviter.SingleId,
			UserId:       inviter.InviteeId,
			UserName:     inviter.InviteeName,
			UserNickname: inviter.InviteeNickname,
			UserDisturb:  inviter.InviterDisturb,
		}
		return h.res.Success(e, http.StatusOK, constant.MsgSingleGetDetailSuccess, getDetailRes)
	}
	invitee := &model.SingleInvitee{
		SingleId:  getDetailReq.SingleId,
		InviteeId: getDetailReq.UserId,
	}
	err = h.ucase.GetDetailForInvitee(invitee)
	if err != nil {
		return err
	}
	getDetailRes := &model.GetDetailRes{
		SingleId:     invitee.SingleId,
		UserId:       invitee.InviterId,
		UserName:     invitee.InviterName,
		UserNickname: invitee.InviterNickname,
		UserDisturb:  invitee.InviteeDisturb,
	}
	return h.res.Success(e, http.StatusOK, constant.MsgSingleGetDetailSuccess, getDetailRes)
}

func (h *singleHandler) Delete(e echo.Context) error {
	deleteReq, ok := e.Get("body").(*model.DeleteReq)
	if !ok {
		return h.res.Fail(e, http.StatusBadRequest, int(constant.ErrBadRequest), constant.MsgBadRequest)
	}
	singleDelete := &model.SingleDelete{
		SingleId:  deleteReq.SingleId,
		InviterId: deleteReq.InviterId,
		InviteeId: deleteReq.InviteeId,
	}
	err := h.ucase.Delete(singleDelete)
	if err != nil {
		return err
	}
	return h.res.Success(e, http.StatusOK, constant.MsgSingleDeleteSuccess, nil)
}
