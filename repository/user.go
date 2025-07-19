package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/wendisx/gorchat/internal/constant"
	"github.com/wendisx/gorchat/internal/log"
	"github.com/wendisx/gorchat/model"
)

// repository -> usecase
type UserRepository interface {
	InsertOne(ctx context.Context, user *model.User) (*model.User, error)
	FindOneById(ctx context.Context, userId int64) (*model.User, error)
	FindOneByName(ctx context.Context, userName string) ([]model.User, error)
	FindBasicLists(ctx context.Context, userSearch model.UserBasic, page *model.Page[model.UserBasic]) error
	UpdateOneById(ctx context.Context, user *model.User) (*model.User, error)
	DeleteOneById(ctx context.Context, userId int64) error
	GetLogger() log.Logger
}

// internal
type userRepository struct {
	db     DBTX
	logger log.Logger
}

func NewUserRepository(db *sql.DB, logger log.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

// The repo layer does not do error handling, it is only responsible for executing SQL.
// The real error handling is judged in the usecase layer.

func (r *userRepository) GetLogger() log.Logger {
	return r.logger
}

// 插入一个用户
func (r *userRepository) InsertOne(ctx context.Context, user *model.User) (*model.User, error) {
	// 启动事务
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error(
			r.logger,
			"Set transaction",
			map[string]any{
				"error": err,
			},
		)
		return nil, &model.DError{
			Code:    constant.ErrTransactionBegin,
			Message: constant.MsgOperationFail,
		}
	}
	// 插入 im_users 基本数据
	insertSql := `
		insert into im_users(user_name,user_password)
		values
		(?,?)
	`
	result, err := tx.Exec(
		insertSql,
		user.UserName,
		user.UserPassword,
	)
	// 插入失败
	if err != nil {
		log.Error(
			r.logger,
			insertSql,
			map[string]any{
				"error": err,
			},
		)
		tx.Rollback()
		return nil, &model.DError{
			Code:    constant.ErrSqlInsertFail,
			Message: constant.MsgOperationFail,
		}
	}
	// 尝试获取插入主键 -- userId
	userId, err := result.LastInsertId()
	if err != nil {
		log.Error(
			r.logger,
			"Get user_id",
			map[string]any{
				"error": err,
			},
		)
		tx.Rollback()
		return nil, &model.DError{
			Code:    constant.ErrOperationFail,
			Message: constant.MsgOperationFail,
		}
	}
	// 尝试创建关联依赖表
	insertSql = `
		insert into im_users_detail(user_id) values (?)
	`
	_, err = tx.Exec(
		insertSql,
		userId,
	)
	// 插入失败
	if err != nil {
		log.Error(
			r.logger,
			insertSql,
			map[string]any{
				"error": err,
			},
		)
		return nil, &model.DError{
			Code:    constant.ErrSqlInsertFail,
			Message: constant.MsgOperationFail,
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Error(
			r.logger,
			"transaction rollback",
			map[string]any{
				"error": err,
			},
		)
		return nil, &model.DError{
			Code:    constant.ErrTransactionFail,
			Message: constant.MsgOperationFail,
		}
	}
	// 插入成功并创建依赖表
	user.UserId = userId
	return user, nil
}

func (r *userRepository) FindOneById(ctx context.Context, userId int64) (*model.User, error) {
	var user model.User
	selectSql := `
		select iu.user_id,iu.user_name,iu.user_password,iud.email,iud.phone,iud.gender,iud.age,iud.address,iud.location,iud.avatar 
		from im_users iu
		left join im_users_detail iud
		on iu.user_id = iud.user_id
		where 
		iu.deleted = ?
		and iu.user_id = ?
	 `
	// 查找用户详细信息
	err := r.db.QueryRowContext(
		ctx,
		selectSql,
		0,
		userId,
	).Scan(
		&user.UserId,
		&user.UserName,
		&user.UserPassword,
		&user.UserEmail,
		&user.UserPhone,
		&user.UserGender,
		&user.UserAge,
		&user.UserAddress,
		&user.UserLocation,
		&user.UserAvatar,
	)
	// 查找失败
	if err != nil {
		log.Error(
			r.logger,
			selectSql,
			map[string]any{
				"error": err,
			},
		)
		return nil, &model.DError{
			Code:    constant.ErrSqlSelectFail,
			Message: constant.MsgOperationFail,
		}
	}
	// 找到用户
	return &user, nil
}

func (r *userRepository) FindOneByName(ctx context.Context, userName string) ([]model.User, error) {
	var users []model.User
	var user model.User
	selectSql := `
		select iu.user_id,iu.user_name,iud.email,iud.phone,iud.gender,iud.age,iud.address,iud.location,iud.avatar 
		from im_users iu
		left join im_users_detail iud
		on iu.user_id = iud.user_id
		where 
		iu.deleted = ?
		and iu.user_id like ?
	 `
	// 查找用户详细信息
	rows, err := r.db.QueryContext(
		ctx,
		selectSql,
		0,
		fmt.Sprintf("%%%s%%", userName),
	)
	// 查找失败
	if err != nil {
		log.Error(
			r.logger,
			selectSql,
			map[string]any{
				"error": err,
			},
		)
		return users, &model.DError{
			Code:    constant.ErrSqlSelectFail,
			Message: constant.MsgOperationFail,
		}
	}
	defer rows.Close()
	// 查找成功
	for rows.Next() {
		err = rows.Scan(
			&user.UserId,
			&user.UserName,
			&user.UserEmail,
			&user.UserPhone,
			&user.UserGender,
			&user.UserAge,
			&user.UserAddress,
			&user.UserLocation,
			&user.UserAvatar,
		)
		users = append(users, user)
		if err != nil {
			log.Error(
				r.logger,
				selectSql,
				map[string]any{
					"error": err,
				},
			)
			return users, &model.DError{
				Code:    constant.ErrSqlSelectFail,
				Message: constant.MsgOperationFail,
			}
		}
	}
	return users, nil
}

func (r *userRepository) FindBasicLists(ctx context.Context, userSearch model.UserBasic, page *model.Page[model.UserBasic]) error {
	var userBasic model.UserBasic
	selectSql := `
		select user_id,user_name
		from im_users 
		where deleted = ? and (user_name like ? or user_id = ?)
		order by user_id
		limit ? offset ?
	`
	rows, err := r.db.QueryContext(
		ctx,
		selectSql,
		0,
		fmt.Sprintf("%%%s%%", userSearch.UserName),
		userSearch.UserId,
		page.PageSize,
		(page.CurrentPage-1)*page.PageSize,
	)
	if err != nil {
		log.Error(
			r.logger,
			selectSql,
			map[string]any{
				"error": err,
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
			&userBasic.UserId,
			&userBasic.UserName,
		)
		page.Items = append(page.Items, userBasic)
		if err != nil {
			log.Error(
				r.logger,
				selectSql,
				map[string]any{
					"error": err,
				},
			)
			return &model.DError{
				Code:    constant.ErrSqlSelectFail,
				Message: constant.MsgOperationFail,
			}
		}
	}
	page.Total = len(page.Items)
	return nil
}

func (r *userRepository) UpdateOneById(ctx context.Context, user *model.User) (*model.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error(
			r.logger,
			"Set transaction",
			map[string]any{
				"error": err,
			},
		)
		return nil, &model.DError{
			Code:    constant.ErrTransactionBegin,
			Message: constant.MsgOperationFail,
		}
	}
	updateSql := `
		update im_users as iu join im_users_detail as iud
		on iu.user_id = iud.user_id
		set
			iu.user_name=?,
			iud.email=?,
			iud.phone=?,
			iud.gender=?,
			iud.age=?,
			iud.address=?,	
			iud.location=?,
			iud.avatar=?
		where 
			iu.user_id = ?
	`
	result, err := tx.Exec(
		updateSql,
		user.UserName,
		user.UserEmail,
		user.UserPhone,
		user.UserGender,
		user.UserAge,
		user.UserAddress,
		user.UserLocation,
		user.UserAvatar,
		user.UserId,
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
		return nil, &model.DError{
			Code:    constant.ErrSqlUpdateFail,
			Message: constant.MsgOperationFail,
		}
	}
	rowChange, err := result.RowsAffected()
	if err != nil || rowChange > 2 {
		log.Error(
			r.logger,
			"not found to update",
			map[string]any{
				"error": err.Error(),
			},
		)
		tx.Rollback()
		return nil, &model.DError{
			Code:    constant.ErrSqlUpdateFail,
			Message: constant.MsgOperationFail,
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Error(
			r.logger,
			"transaction rollback",
			map[string]any{
				"error": err,
			},
		)
		return nil, &model.DError{
			Code:    constant.ErrTransactionFail,
			Message: constant.MsgOperationFail,
		}
	}
	user, err = r.FindOneById(ctx, user.UserId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) DeleteOneById(ctx context.Context, userId int64) error {
	deleteSql := `
		update im_users set deleted = ? where user_id = ? and deleted = ?
	`
	result, err := r.db.ExecContext(
		ctx,
		deleteSql,
		1,
		userId,
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
		return &model.DError{
			Code:    constant.ErrSqlDeleteFail,
			Message: constant.MsgOperationFail,
		}
	}
	rowChange, err := result.RowsAffected()
	if err != nil || rowChange != 1 {
		log.Error(
			r.logger,
			"not found to delete",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlDeleteFail,
			Message: constant.MsgOperationFail,
		}
	}
	return nil
}
