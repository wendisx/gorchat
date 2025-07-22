package repository

import (
	"context"

	"github.com/wendisx/gorchat/internal/constant"
	"github.com/wendisx/gorchat/internal/log"
	"github.com/wendisx/gorchat/model"
)

type SingleRepository interface {
	GetLogger() log.Logger
	InsertUnAccepted(ctx context.Context, singleInvite *model.SingleInvite) error
	FindByInviter(ctx context.Context, singleInviter *model.SingleInviter) error
	FindByInvitee(ctx context.Context, singleInvitee *model.SingleInvitee) error
	UpdateByInviter(ctx context.Context, singleInviter *model.SingleInviter) error
	UpdateByInvitee(ctx context.Context, singleInvitee *model.SingleInvitee) error
	UpdateByAccept(ctx context.Context, singleAccept *model.SingleAccept) error
	Update(ctx context.Context, single *model.Single) error
	Delete(ctx context.Context, single *model.SingleDelete) error
}

type singleRepository struct {
	db     DBTX
	logger log.Logger
}

func NewSingleRepository(db DBTX, logger log.Logger) SingleRepository {
	return &singleRepository{
		db:     db,
		logger: logger,
	}
}

func (r *singleRepository) GetLogger() log.Logger {
	return r.logger
}

// 插入新的 single row,默认为逻辑删除
func (r *singleRepository) InsertUnAccepted(ctx context.Context, singleInvite *model.SingleInvite) error {
	insertSql := `
		insert into im_single_chat(single_id,inviter_id,invitee_id,invitee_nickname,inviter_disturb,deleted)
		values
		(?,?,?,?,?,?)
	`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrTransactionBegin,
			Message: constant.MsgOperationFail,
		}
	}
	result, err := tx.ExecContext(
		ctx,
		insertSql,
		singleInvite.SingleId,
		singleInvite.InviterId,
		singleInvite.InviteeId,
		singleInvite.InviteeNickname,
		singleInvite.InviterDisturb,
		singleInvite.Deleted,
	)
	if err != nil {
		tx.Rollback()
		log.Error(
			r.logger,
			insertSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlInsertFail,
			Message: constant.MsgOperationFail,
		}
	}
	rowChange, err := result.RowsAffected()
	if err != nil || rowChange > 1 {
		tx.Rollback()
		log.Error(
			r.logger,
			"can't insert",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlInsertFail,
			Message: constant.MsgOperationFail,
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Error(
			r.logger,
			"transaction commit fail",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrTransactionFail,
			Message: constant.MsgOperationFail,
		}
	}
	return nil
}

func (r *singleRepository) FindByInviter(ctx context.Context, singleInviter *model.SingleInviter) error {
	selectSql := `
		select isc.invitee_id,isc.invitee_nickname,iu.user_name,isc.inviter_disturb
		from im_single_chat isc
		left join im_users iu
		on isc.invitee_id = iu.user_id
		where
		isc.single_id = ? and isc.deleted = ? and isc.inviter_id = ?
	`
	err := r.db.QueryRowContext(
		ctx,
		selectSql,
		singleInviter.SingleId,
		0,
		singleInviter.InviterId,
	).Scan(
		&singleInviter.InviteeId,
		&singleInviter.InviteeNickname,
		&singleInviter.InviteeName,
		&singleInviter.InviterDisturb,
	)
	if err != nil {
		log.Error(
			r.logger,
			selectSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlSelectFail,
			Message: constant.MsgOperationFail,
		}
	}
	return nil
}

func (r *singleRepository) FindByInvitee(ctx context.Context, singleInvitee *model.SingleInvitee) error {
	selectSql := `
		select isc.inviter_id,isc.inviter_nickname,iu.user_name,isc.invitee_disturb
		from im_single_chat isc
		left join im_users iu
		on isc.inviter_id = iu.user_id
		where
		isc.single_id = ? and isc.deleted = ? and isc.invitee_id = ?
	`
	err := r.db.QueryRowContext(
		ctx,
		selectSql,
		singleInvitee.SingleId,
		0,
		singleInvitee.InviteeId,
	).Scan(
		&singleInvitee.InviterId,
		&singleInvitee.InviterNickname,
		&singleInvitee.InviterName,
		&singleInvitee.InviteeDisturb,
	)
	if err != nil {
		log.Error(
			r.logger,
			selectSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlSelectFail,
			Message: constant.MsgOperationFail,
		}
	}
	return nil
}

func (r *singleRepository) UpdateByInviter(ctx context.Context, singleInviter *model.SingleInviter) error {
	updateSql := `
		update im_single_chat
		set
			invitee_nickname = ?,
			inviter_disturb = ?
		where
		single_id = ? and deleted = ? and inviter_id = ?
	`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrTransactionBegin,
			Message: constant.MsgOperationFail,
		}
	}
	result, err := tx.ExecContext(
		ctx,
		updateSql,
		singleInviter.InviteeNickname,
		singleInviter.InviterDisturb,
		singleInviter.SingleId,
		0,
		singleInviter.InviterId,
	)
	if err != nil {
		tx.Rollback()
		log.Error(
			r.logger,
			updateSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlUpdateFail,
			Message: constant.MsgOperationFail,
		}
	}
	rowChange, err := result.RowsAffected()
	if err != nil || rowChange > 1 {
		tx.Rollback()
		log.Error(
			r.logger,
			"update fail",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlUpdateFail,
			Message: constant.MsgOperationFail,
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Error(
			r.logger,
			"transaction commit fail",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrTransactionFail,
			Message: constant.MsgOperationFail,
		}
	}
	err = r.FindByInviter(ctx, singleInviter)
	if err != nil {
		return err
	}
	return nil
}

func (r *singleRepository) UpdateByInvitee(ctx context.Context, singleInvitee *model.SingleInvitee) error {
	updateSql := `
		update im_single_chat
		set
			inviter_nickname = ?,
			invitee_disturb = ?
		where
		single_id = ? and deleted = ? and invitee_id = ?
	`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrTransactionBegin,
			Message: constant.MsgOperationFail,
		}
	}
	result, err := tx.ExecContext(
		ctx,
		updateSql,
		singleInvitee.InviterNickname,
		singleInvitee.InviteeDisturb,
		singleInvitee.SingleId,
		0,
		singleInvitee.InviteeId,
	)
	if err != nil {
		tx.Rollback()
		log.Error(
			r.logger,
			updateSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlUpdateFail,
			Message: constant.MsgOperationFail,
		}
	}
	rowChange, err := result.RowsAffected()
	if err != nil || rowChange > 1 {
		tx.Rollback()
		log.Error(
			r.logger,
			"update fail",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlUpdateFail,
			Message: constant.MsgOperationFail,
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Error(
			r.logger,
			"transaction commit fail",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrTransactionFail,
			Message: constant.MsgOperationFail,
		}
	}
	err = r.FindByInvitee(ctx, singleInvitee)
	if err != nil {
		return err
	}
	return nil

}

func (r *singleRepository) UpdateByAccept(ctx context.Context, singleAccept *model.SingleAccept) error {
	updateSql := `
		update im_single_chat
		set
			inviter_nickname = ?,
			invitee_disturb = ?,
			deleted = ?
		where
		single_id = ? and deleted = ?
	`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrTransactionBegin,
			Message: constant.MsgOperationFail,
		}
	}
	result, err := tx.ExecContext(
		ctx,
		updateSql,
		singleAccept.InviterNickname,
		singleAccept.InviteeDisturb,
		singleAccept.Deleted,
		singleAccept.SingleId,
		singleAccept.Deleted+1,
	)
	if err != nil {
		tx.Rollback()
		log.Error(
			r.logger,
			updateSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlUpdateFail,
			Message: constant.MsgOperationFail,
		}
	}
	rowChange, err := result.RowsAffected()
	if err != nil || rowChange > 1 {
		tx.Rollback()
		log.Error(
			r.logger,
			"update fail",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlUpdateFail,
			Message: constant.MsgOperationFail,
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Error(
			r.logger,
			"transaction commit fail",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrTransactionFail,
			Message: constant.MsgOperationFail,
		}
	}
	return nil
}

func (r *singleRepository) Update(ctx context.Context, single *model.Single) error {
	updateSql := `
		update im_single_chat
		set
		inviter_nickname = ?,
		inviter_disturb = ?,
		invitee_nickname = ?,
		invitee_disturb = ?
		where
		single_id = ? and deleted = ? and inviter_id = ? and invitee_id = ? 
	`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrTransactionBegin,
			Message: constant.MsgOperationFail,
		}
	}
	result, err := tx.ExecContext(
		ctx,
		updateSql,
		single.InviterNickname,
		single.InviterDisturb,
		single.InviteeNickname,
		single.InviteeDisturb,
		single.SingleId,
		single.Deleted,
		single.InviterId,
		single.InviteeId,
	)
	if err != nil {
		tx.Rollback()
		log.Error(
			r.logger,
			updateSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlUpdateFail,
			Message: constant.MsgOperationFail,
		}
	}
	rowChange, err := result.RowsAffected()
	if err != nil || rowChange > 1 {
		tx.Rollback()
		log.Error(
			r.logger,
			"nothing to update",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlUpdateFail,
			Message: constant.MsgOperationFail,
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Error(
			r.logger,
			"transaction commit fail",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrTransactionFail,
			Message: constant.MsgOperationFail,
		}
	}
	return nil
}

func (r *singleRepository) Delete(ctx context.Context, single *model.SingleDelete) error {
	deleteSql := `
		update im_single_chat 
		set
		deleted = ?
		where
		single_id = ? and deleted = ? and inviter_id = ? and invitee_id = ?
	`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrTransactionBegin,
			Message: constant.MsgOperationFail,
		}
	}
	result, err := tx.ExecContext(
		ctx,
		deleteSql,
		1,
		single.SingleId,
		0,
		single.InviterId,
		single.InviteeId,
	)
	if err != nil {
		tx.Rollback()
		log.Error(
			r.logger,
			deleteSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlDeleteFail,
			Message: constant.MsgOperationFail,
		}
	}
	rowChange, err := result.RowsAffected()
	if err != nil || rowChange > 1 {
		tx.Rollback()
		log.Error(
			r.logger,
			"nothing to delete",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlDeleteFail,
			Message: constant.MsgOperationFail,
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Error(
			r.logger,
			"transaction commit fail",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrTransactionFail,
			Message: constant.MsgOperationFail,
		}
	}
	return nil
}
