package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/wendisx/gorchat/internal/constant"
	"github.com/wendisx/gorchat/internal/log"
	"github.com/wendisx/gorchat/model"
)

// repository -> usecase
type UserRepository interface {
	InsertOne(ctx context.Context, user model.User) (int64, error)
	FindOneById(ctx context.Context, userId int64) (model.User, error)
	FindOneByEmail(ctx context.Context, email string) (model.User, error)
	UpdateOneById(ctx context.Context, user model.User) (model.User, error)
	DeleteById(ctx context.Context, userId int64) error
	GetLogger() log.Logger
}

// internal
type userRepository struct {
	db     DBTX
	Logger log.Logger
}

func NewUserRepository(db *sql.DB, logger log.Logger) UserRepository {
	return &userRepository{
		db:     db,
		Logger: logger,
	}
}

// The repo layer does not do error handling, it is only responsible for executing SQL.
// The real error handling is judged in the usecase layer.

func (r *userRepository) GetLogger() log.Logger {
	return r.Logger
}

func (r *userRepository) InsertOne(ctx context.Context, user model.User) (int64, error) {
	result, err := r.db.ExecContext(ctx, "insert into `user`(user_name,user_password,user_email,create_time,update_time) values (?,?,?,?,?)", user.UserName, user.Password, user.Email, time.Now(), time.Now())
	if err != nil {
		log.Error(
			r.Logger,
			"insert a user",
			map[string]any{
				"error": err.Error(),
			},
		)
		return -1, &model.DError{
			Code:    constant.ErrSqlExcution,
			Message: constant.MsgOperationFail,
		}
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Error(
			r.Logger,
			"get inserted user_id",
			map[string]any{
				"error": err.Error(),
			},
		)
		return -1, &model.DError{
			Code:    constant.ErrInsert,
			Message: constant.MsgOperationFail,
		}
	}
	log.Info(
		r.Logger,
		"insert a user",
		map[string]any{
			"userId": id,
		},
	)
	return id, nil
}

func (r *userRepository) FindOneById(ctx context.Context, userId int64) (model.User, error) {
	var user model.User
	err := r.db.QueryRowContext(ctx, "select user_id,user_name,user_password,user_email from `user` where user_id = ?", userId).Scan(
		&user.UserId,
		&user.UserName,
		&user.Password,
		&user.Email,
	)
	if err != nil {
		log.Error(
			r.Logger,
			"find a user by user_id",
			map[string]any{
				"error": err.Error(),
			},
		)
		return user, &model.DError{
			Code:    constant.ErrFind,
			Message: constant.MsgOperationFail,
		}
	}
	log.Info(
		r.Logger,
		"find a user by user_id",
		map[string]any{
			"userId":    user.UserId,
			"userName":  user.UserName,
			"userEmail": user.Email,
		},
	)
	return user, nil
}

func (r *userRepository) FindOneByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	err := r.db.QueryRowContext(ctx, "select user_id,user_name,user_password,user_email from `user` where user_email = ?", email).Scan(
		&user.UserId,
		&user.UserName,
		&user.Password,
		&user.Email,
	)
	if err != nil {
		log.Error(
			r.Logger,
			"find a user by user_email",
			map[string]any{
				"error": err.Error(),
			},
		)
		return user, &model.DError{
			Code:    constant.ErrFind,
			Message: constant.MsgOperationFail,
		}
	}
	log.Info(
		r.Logger,
		"find a user by user_email",
		map[string]any{
			"userId":    user.UserId,
			"userName":  user.UserName,
			"userEmail": user.Email,
		},
	)
	return user, nil
}

func (r *userRepository) UpdateOneById(ctx context.Context, user model.User) (model.User, error) {
	var nUser model.User
	result, err := r.db.ExecContext(ctx, "update `user` set user_name = ?, user_email = ?, update_time = ? where user_id = ?", user.UserName, user.Email, time.Now(), user.UserId)
	if err != nil {
		log.Error(
			r.Logger,
			"update a user by user_id",
			map[string]any{
				"error": err.Error(),
			},
		)
		return nUser, &model.DError{
			Code:    constant.ErrUpdate,
			Message: constant.MsgOperationFail,
		}
	}
	rowUp, err := result.RowsAffected()
	if err != nil || rowUp != 1 {
		log.Error(
			r.Logger,
			"update a user by user_id",
			map[string]any{
				"error": err.Error(),
			},
		)
		return nUser, &model.DError{
			Code:    constant.ErrUpdate,
			Message: constant.MsgOperationFail,
		}
	}
	nUser, err = r.FindOneById(ctx, user.UserId)
	if err != nil {
		return nUser, err
	}
	log.Info(
		r.Logger,
		"update a user by user_id",
		map[string]any{
			"userId":    nUser.UserId,
			"userName":  nUser.UserName,
			"userEmail": nUser.Email,
		},
	)
	return nUser, nil
}

func (r *userRepository) DeleteById(ctx context.Context, userId int64) error {
	result, err := r.db.ExecContext(ctx, "delete from `user` where user_id = ?", userId)
	if err != nil {
		log.Error(
			r.Logger,
			"delete a user by user_id",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrSqlExcution,
			Message: constant.MsgOperationFail,
		}
	}
	rowNum, err := result.RowsAffected()
	if err != nil || rowNum != 1 {
		log.Error(
			r.Logger,
			"delete a user by user_id",
			map[string]any{
				"error": err.Error(),
			},
		)
		return &model.DError{
			Code:    constant.ErrDelete,
			Message: constant.MsgOperationFail,
		}
	}
	log.Info(
		r.Logger,
		"delete a user by user_id",
		map[string]any{
			"userId": userId,
		},
	)
	return nil
}
