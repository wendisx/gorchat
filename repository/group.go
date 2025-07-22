package repository

import (
	"context"

	"github.com/wendisx/gorchat/internal/constant"
	"github.com/wendisx/gorchat/internal/log"
	"github.com/wendisx/gorchat/model"
)

type GroupRepository interface {
	GetLogger() log.Logger
	InsertOneGroup(ctx context.Context, group *model.Group) (int64, error)
	InsertUserInGroup(ctx context.Context, groupToUser *model.GroupToUser) error
	FindGroupToUser(ctx context.Context, groupToUser *model.GroupToUser) error
	FindGroup(ctx context.Context, group *model.Group) error
	FindGroupPassword(ctx context.Context, groupId int64, groupPassword *string) error
	FindUserRoleId(ctx context.Context, userRole string) (int, error)
	FindGroupBasic(ctx context.Context, groupBasic *model.GroupBasic) error
	FindGroupUsers(ctx context.Context, groupUser *model.GroupUser, page *model.Page[*model.GroupToUserItem]) error
	FindGroups(ctx context.Context, groupItem *model.GroupItem, page *model.Page[*model.GroupItem]) error
	FindGroupAllUsers(ctx context.Context, groupId int64, page *model.Page[*model.GroupToUserItem]) error
	UpdateGroup(ctx context.Context, group *model.Group) error
	UpdateGroupToUser(ctx context.Context, groupToUser *model.GroupToUser) error
	DeleteGroup(ctx context.Context, groupId int64) error
	DeleteGroupToUser(ctx context.Context, groupId, userId int64) error
}

type groupRepository struct {
	db     DBTX
	logger log.Logger
}

func NewGroupRepository(db DBTX, logger log.Logger) GroupRepository {
	return &groupRepository{
		db:     db,
		logger: logger,
	}
}

func (r *groupRepository) GetLogger() log.Logger {
	return r.logger
}

func (r *groupRepository) InsertOneGroup(ctx context.Context, group *model.Group) (int64, error) {
	insertSql := `
		insert into im_groups(group_name,group_password)
		values
		(?,?)
	`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, &model.DError{
			Code:    constant.ErrTransactionBegin,
			Message: constant.MsgTransactionBegin,
		}
	}
	result, err := tx.ExecContext(
		ctx,
		insertSql,
		group.GroupName,
		group.GroupPassword,
	)
	if err != nil {
		log.Error(
			r.logger,
			insertSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		tx.Rollback()
		return -1, &model.DError{
			Code:    constant.ErrSqlInsertFail,
			Message: constant.MsgSqlInsertFail,
		}
	}
	group.GroupId, err = result.LastInsertId()
	if err != nil || group.GroupId < 1000000 {
		log.Error(
			r.logger,
			insertSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		return -1, &model.DError{
			Code:    constant.ErrSqlInsertFail,
			Message: constant.MsgSqlInsertFail,
		}
	}
	insertSql = `
		insert into im_groups_detail(group_id,group_avatar,max_size)
		values
		(?,?,?)
	`
	result, err = tx.ExecContext(
		ctx,
		insertSql,
		group.GroupId,
		group.GroupAvatar,
		group.GroupMaxSize,
	)
	if err != nil {
		log.Error(
			r.logger,
			insertSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		tx.Rollback()
		return -1, &model.DError{
			Code:    constant.ErrSqlInsertFail,
			Message: constant.MsgSqlInsertFail,
		}
	}
	rowChange, err := result.RowsAffected()
	if err != nil || rowChange > 1 {
		log.Error(
			r.logger,
			insertSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		tx.Rollback()
		return -1, &model.DError{
			Code:    constant.ErrSqlInsertFail,
			Message: constant.MsgSqlInsertFail,
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
		return -1, &model.DError{
			Code:    constant.ErrTransactionFail,
			Message: constant.MsgTransactionFail,
		}
	}
	return group.GroupId, nil
}

func (r *groupRepository) InsertUserInGroup(ctx context.Context, groupToUser *model.GroupToUser) error {
	insertSql := `
		insert into im_groups_users(group_id,user_id,group_nickname,user_nickname,user_role,user_role_nickname,disturb)
		values
		(?,?,?,?,?,?,?)
	`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrTransactionBegin,
			Message: constant.MsgTransactionBegin,
		}
	}
	result, err := tx.ExecContext(
		ctx,
		insertSql,
		groupToUser.GroupId,
		groupToUser.UserId,
		groupToUser.GroupNickname,
		groupToUser.UserNickname,
		groupToUser.UserRoleId,
		groupToUser.UserRoleNickname,
		groupToUser.UserDisturb,
	)
	if err != nil {
		log.Error(
			r.logger,
			insertSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		tx.Rollback()
		return &model.DError{
			Code:    constant.ErrSqlInsertFail,
			Message: constant.MsgSqlInsertFail,
		}
	}
	rowChange, err := result.RowsAffected()
	if err != nil || rowChange > 1 {
		log.Error(
			r.logger,
			insertSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		tx.Rollback()
		return &model.DError{
			Code:    constant.ErrSqlInsertFail,
			Message: constant.MsgSqlInsertFail,
		}
	}
	updateSql := `
		update im_groups_detail
		set
			current_size = current_size+1
		where
			group_id = ?
	`
	_, err = tx.ExecContext(
		ctx,
		updateSql,
		groupToUser.GroupId,
	)
	if err != nil {
		log.Error(
			r.logger,
			updateSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		tx.Rollback()
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
			Message: constant.MsgTransactionFail,
		}
	}
	return nil
}

func (r *groupRepository) FindUserRoleId(ctx context.Context, userRole string) (int, error) {
	selectSql := `
		select role_id from im_users_role
		where 
			role_name = ? and deleted = ?
	`
	userRoleId := new(int)
	err := r.db.QueryRowContext(
		ctx,
		selectSql,
		userRole,
		0,
	).Scan(
		userRoleId,
	)
	if err != nil {
		log.Error(
			r.logger,
			selectSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		return -1, &model.DError{
			Code:    constant.ErrSqlSelectFail,
			Message: constant.MsgOperationFail,
		}
	}
	return *userRoleId, nil
}

func (r *groupRepository) FindGroupToUser(ctx context.Context, groupToUser *model.GroupToUser) error {
	selectSql := `
		select ig.group_name,igd.group_avatar,igd.max_size,igd.current_size
		from im_groups ig
		left join im_groups_detail igd
		on ig.group_id = igd.group_id
		where
			ig.group_id = ? and ig.deleted = ?
	`
	err := r.db.QueryRowContext(
		ctx,
		selectSql,
		groupToUser.GroupId,
		0,
	).Scan(
		&groupToUser.GroupName,
		&groupToUser.GroupAvatar,
		&groupToUser.GroupMaxSize,
		&groupToUser.GroupCurrentSize,
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
			Message: constant.MsgSqlSelectFail,
		}
	}
	selectSql = `
		select group_nickname,user_nickname,user_role,user_role_nickname,disturb
		from im_groups_users
		where 
			group_id = ? and deleted = ? and user_id = ?
	`
	err = r.db.QueryRowContext(
		ctx,
		selectSql,
		groupToUser.GroupId,
		0,
		groupToUser.UserId,
	).Scan(
		&groupToUser.GroupNickname,
		&groupToUser.UserNickname,
		&groupToUser.UserRoleId,
		&groupToUser.UserRoleNickname,
		&groupToUser.UserDisturb,
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
			Message: constant.MsgSqlSelectFail,
		}
	}
	return nil
}

func (r *groupRepository) FindGroup(ctx context.Context, group *model.Group) error {
	selectSql := `
		select ig.group_name,igd.group_avatar,igd.max_size,igd.current_size
		from im_groups ig
		left join im_groups_detail igd
		on ig.group_id = igd.group_id
		where
			ig.group_id = ? and ig.deleted = ? and ig.group_name like ?
	`
	err := r.db.QueryRowContext(
		ctx,
		selectSql,
		group.GroupId,
		0,
		group.GroupName,
	).Scan(
		&group.GroupName,
		&group.GroupAvatar,
		&group.GroupMaxSize,
		&group.GroupCurrentSize,
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
			Message: constant.MsgSqlSelectFail,
		}
	}
	return nil
}

func (r *groupRepository) FindGroupPassword(ctx context.Context, groupId int64, groupPassword *string) error {
	selectSql := `
		select group_password 
		from im_groups
		where
			group_id = ? and deleted = ?
	`
	err := r.db.QueryRowContext(
		ctx,
		selectSql,
		groupId,
		0,
	).Scan(groupPassword)
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
			Message: constant.MsgSqlSelectFail,
		}
	}
	return nil
}

func (r *groupRepository) FindGroupBasic(ctx context.Context, groupBasic *model.GroupBasic) error {
	selectSql := `
		select ig.group_name,igd.group_avatar,igd.max_size,igu.group_nickname
		from im_groups ig
		left join im_groups_detail igd on ig.group_id = igd.group_id
		left join im_groups_users igu on ig.group_id = igu.group_id
		where
			ig.group_id = ? and ig.deleted = ? and igu.deleted = ? and igu.user_id = ?
	`
	err := r.db.QueryRowContext(
		ctx,
		selectSql,
		groupBasic.GroupId,
		0,
		0,
		groupBasic.UserId,
	).Scan(
		&groupBasic.GroupName,
		&groupBasic.GroupAvatar,
		&groupBasic.GroupMaxSize,
		&groupBasic.GroupNickname,
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
			Message: constant.MsgSqlSelectFail,
		}
	}
	return nil
}

func (r *groupRepository) FindGroupUsers(ctx context.Context, groupUser *model.GroupUser, page *model.Page[*model.GroupToUserItem]) error {
	selectSql := `
		select iu.user_id,iu.user_name,igu.user_nickname,iur.role_name,igu.user_role_nickname
		from im_groups_users igu
		left join im_users iu on iu.user_id = igu.user_id
		left join im_users_role iur on igu.user_role = iur.role_id
		where
			igu.deleted = ? and iur.deleted = ? and iu.deleted = ? and igu.group_id = ?
			and (igu.user_id = ? or (igu.user_nickname like ? or iu.user_name like ?))
		order by user_id
		limit ? offset ?
	`
	rows, err := r.db.QueryContext(
		ctx,
		selectSql,
		0,
		0,
		0,
		groupUser.GroupId,
		groupUser.UserId,
		groupUser.UserNickname,
		groupUser.UserName,
		page.PageSize,
		(page.CurrentPage-1)*page.PageSize,
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
			Message: constant.MsgSqlSelectFail,
		}
	}
	defer rows.Close()
	for rows.Next() {
		var groupItem model.GroupToUserItem
		err := rows.Scan(
			&groupItem.UserId,
			&groupItem.UserName,
			&groupItem.UserNickname,
			&groupItem.UserRole,
			&groupItem.UserRoleNickname,
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
				Message: constant.MsgSqlSelectFail,
			}
		}
		page.Items = append(page.Items, &groupItem)
	}
	page.Total = len(page.Items)
	return nil
}

func (r *groupRepository) FindGroupAllUsers(ctx context.Context, groupId int64, page *model.Page[*model.GroupToUserItem]) error {
	selectSql := `
		select igu.user_id,iu.user_name,igu.user_nickname,iur.role_name,igu.user_role_nickname
		from im_groups_users igu
		left join im_users iu on iu.user_id = igu.user_id
		left join im_users_role iur on igu.user_role = iur.role_id
		where
			igu.deleted = ? and igu.group_id = ? and iur.deleted = ? and iu.deleted = ? 
		order by igu.user_id
		limit ? offset ?
	`
	rows, err := r.db.QueryContext(
		ctx,
		selectSql,
		0,
		groupId,
		0,
		0,
		page.PageSize,
		(page.CurrentPage-1)*page.PageSize,
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
			Message: constant.MsgSqlSelectFail,
		}
	}
	defer rows.Close()
	for rows.Next() {
		var groupItem model.GroupToUserItem
		err := rows.Scan(
			&groupItem.UserId,
			&groupItem.UserName,
			&groupItem.UserNickname,
			&groupItem.UserRole,
			&groupItem.UserRoleNickname,
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
				Message: constant.MsgSqlSelectFail,
			}
		}
		page.Items = append(page.Items, &groupItem)
	}
	page.Total = len(page.Items)
	return nil
}

func (r *groupRepository) FindGroups(ctx context.Context, groupItem *model.GroupItem, page *model.Page[*model.GroupItem]) error {
	selectSql := `
		select group_id,group_name
		from im_groups
		where
			deleted = ? and (group_id = ? or group_name like ?)
		order by group_id
		limit ? offset ?
	`
	rows, err := r.db.QueryContext(
		ctx,
		selectSql,
		0,
		groupItem.GroupId,
		groupItem.GroupName,
		page.PageSize,
		(page.CurrentPage-1)*page.PageSize,
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
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&groupItem.GroupId,
			&groupItem.GroupName,
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
		page.Items = append(page.Items, groupItem)
	}
	page.Total = len(page.Items)
	return nil
}

func (r *groupRepository) UpdateGroup(ctx context.Context, group *model.Group) error {
	updateSql := `
		update im_groups as ig join im_groups_detail as igd
		on ig.group_id = igd.group_id
		set
			ig.group_name = ?,
			ig.group_password = ?,
			igd.group_avatar = ?,
			igd.max_size = ?
		where
			ig.group_id = ? and ig.deleted = ?
	`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrTransactionBegin,
			Message: constant.MsgTransactionBegin,
		}
	}
	result, err := tx.ExecContext(
		ctx,
		updateSql,
		group.GroupName,
		group.GroupPassword,
		group.GroupAvatar,
		group.GroupMaxSize,
		group.GroupId,
		0,
	)
	if err != nil {
		log.Error(
			r.logger,
			updateSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		tx.Rollback()
		return &model.DError{
			Code:    constant.ErrSqlDeleteFail,
			Message: constant.MsgSqlDeleteFail,
		}
	}
	rowChange, err := result.RowsAffected()
	if err != nil || rowChange > 2 {
		log.Error(
			r.logger,
			"update too much",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlDeleteFail,
			Message: constant.MsgSqlDeleteFail,
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
			Message: constant.MsgTransactionFail,
		}
	}
	return nil
}

func (r *groupRepository) UpdateGroupToUser(ctx context.Context, groupToUser *model.GroupToUser) error {
	updateSql := `
		update im_groups_users
		set
			group_nickname = ?,
			user_nickname = ?,
			user_role = ?,
			user_role_nickname = ?,
			disturb = ?
		where
			group_id = ? and deleted = ? and user_id = ?
	`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrTransactionBegin,
			Message: constant.MsgTransactionBegin,
		}
	}
	result, err := tx.ExecContext(
		ctx,
		updateSql,
		groupToUser.GroupNickname,
		groupToUser.UserNickname,
		groupToUser.UserRoleId,
		groupToUser.UserRoleNickname,
		groupToUser.UserDisturb,
		groupToUser.GroupId,
		0,
		groupToUser.UserId,
	)
	if err != nil {
		log.Error(
			r.logger,
			updateSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		tx.Rollback()
		return &model.DError{
			Code:    constant.ErrSqlDeleteFail,
			Message: constant.MsgSqlDeleteFail,
		}
	}
	rowChange, err := result.RowsAffected()
	if err != nil || rowChange > 1 {
		log.Error(
			r.logger,
			"update wrong",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlDeleteFail,
			Message: constant.MsgSqlDeleteFail,
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
			Message: constant.MsgTransactionFail,
		}
	}
	return nil
}

// 理论上的删除群聊需要递归先将对应群聊下的所有用户删除确保无意外操作后执行删除群聊
func (r *groupRepository) DeleteGroup(ctx context.Context, groupId int64) error {
	deleteSql := `
		update im_groups
		set
			deleted = ?
		where
			group_id = ? and deleted = ? 
	`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrTransactionBegin,
			Message: constant.MsgTransactionBegin,
		}
	}
	result, err := tx.ExecContext(
		ctx,
		deleteSql,
		1,
		groupId,
		0,
	)
	if err != nil {
		log.Error(
			r.logger,
			deleteSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		tx.Rollback()
		return &model.DError{
			Code:    constant.ErrSqlDeleteFail,
			Message: constant.MsgSqlDeleteFail,
		}
	}
	rowChange, err := result.RowsAffected()
	if err != nil || rowChange != 1 {
		log.Error(
			r.logger,
			"nothing to delete",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlDeleteFail,
			Message: constant.MsgSqlDeleteFail,
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
			Message: constant.MsgTransactionFail,
		}
	}
	return nil
}

func (r *groupRepository) DeleteGroupToUser(ctx context.Context, groupId, userId int64) error {
	deleteSql := `
		update im_groups_users
		set
			deleted = ?
		where
			group_id = ? and deleted = ? and user_id = ?
	`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &model.DError{
			Code:    constant.ErrTransactionBegin,
			Message: constant.MsgTransactionBegin,
		}
	}
	result, err := tx.ExecContext(
		ctx,
		deleteSql,
		1,
		groupId,
		0,
		userId,
	)
	if err != nil {
		log.Error(
			r.logger,
			deleteSql,
			map[string]any{
				"error": err.Error(),
			},
		)
		tx.Rollback()
		return &model.DError{
			Code:    constant.ErrSqlDeleteFail,
			Message: constant.MsgSqlDeleteFail,
		}
	}
	rowChange, err := result.RowsAffected()
	if err != nil || rowChange != 1 {
		log.Error(
			r.logger,
			"nothing to delete",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlDeleteFail,
			Message: constant.MsgSqlDeleteFail,
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
			Message: constant.MsgTransactionFail,
		}
	}
	return nil
}
