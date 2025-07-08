package repository

import (
	"context"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/wendisx/gorchat/internal/constant"
	"github.com/wendisx/gorchat/internal/log"
	"github.com/wendisx/gorchat/model"
)

var (
	ctx context.Context

	logger log.Logger
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	logger = log.NewLogger(constant.DEBUG).Sugar()
	code := m.Run()
	os.Exit(code)
}

func TestInsertOneUser(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatalf("setup sqlmock fail")
	}
	defer db.Close()

	userRepo := NewUserRepository(db, logger)

	user := model.User{
		UserName:   "tom",
		Password:   "tom123",
		Email:      "tom@gmail.com",
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta("insert into `user`(user_name,user_password,user_email,create_time,update_time) values (?,?,?,?,?)")).
		WithArgs("tom", "tom123", "tom@gmail.com", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = userRepo.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(
			userRepo.GetLogger(),
			"test -- InsertOne",
			map[string]any{
				"error": err.Error(),
			},
		)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		userRepo.GetLogger().Fatalf("unexpected test error: %s", err)
	}
}

func TestFindOneUser(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatalf("setup sqlmock fail")
	}
	defer db.Close()

	userRepo := NewUserRepository(db, logger)

	mock.ExpectQuery(regexp.QuoteMeta("select user_id,user_name,user_email from `user` where user_id = ?")).
		WithArgs(1).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id", "user_name", "user_email"}).AddRow(1, "tom", "tom@gmail.com"),
		)

	mock.ExpectQuery(regexp.QuoteMeta("select user_id,user_name,user_email from `user` where user_email = ?")).
		WithArgs("tom@gmail.com").
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id", "user_name", "user_email"}).AddRow(1, "tom", "tom@gmail.com"),
		)

	var userId int64 = 1
	var userEmail string = "tom@gmail.com"
	_, err = userRepo.FindOneById(ctx, userId)
	if err != nil {
		log.Fatal(
			userRepo.GetLogger(),
			"test -- FindOneById",
			map[string]any{
				"error": err.Error(),
			},
		)
	}
	_, err = userRepo.FindOneByEmail(ctx, userEmail)
	if err != nil {
		log.Fatal(
			userRepo.GetLogger(),
			"test -- FindOneByEmail",
			map[string]any{
				"error": err.Error(),
			},
		)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		userRepo.GetLogger().Fatalf("unexpected test error: %s", err)
	}
}

func TestUpdateOneUser(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatalf("setup sqlmock fail")
	}
	defer db.Close()

	userRepo := NewUserRepository(db, logger)

	mock.ExpectExec(regexp.QuoteMeta("update `user` set user_name = ?, user_email = ?, update_time = ? where user_id = ?")).
		WithArgs("ttt", "ttt@firefox.com", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(regexp.QuoteMeta("select user_id,user_name,user_email from `user` where user_id = ?")).
		WithArgs(1).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id", "user_name", "user_email"}).AddRow(1, "ttt", "ttt@firefox.com"),
		)

	user := model.User{
		UserId:   1,
		UserName: "ttt",
		Email:    "ttt@firefox.com",
	}

	_, err = userRepo.UpdateOneById(ctx, user)
	if err != nil {
		log.Fatal(
			userRepo.GetLogger(),
			"test -- UpdateOneById",
			map[string]any{
				"error": err.Error(),
			},
		)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		userRepo.GetLogger().Fatalf("unexpected test error: %s", err)
	}
}

func TestDeleteOneUserById(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatalf("setup sqlmock fail")
	}
	defer db.Close()

	userRepo := NewUserRepository(db, logger)

	var userId int64 = 1

	mock.ExpectExec(regexp.QuoteMeta("delete from `user` where user_id = ?")).
		WithArgs(userId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = userRepo.DeleteById(ctx, userId)
	if err != nil {
		log.Fatal(
			userRepo.GetLogger(),
			"test -- DeleteOneUserById",
			map[string]any{
				"error": err.Error(),
			},
		)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		userRepo.GetLogger().Fatalf("unexpected test error: %s", err)
	}
}
