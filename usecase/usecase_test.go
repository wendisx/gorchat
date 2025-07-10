package usecase

import (
	"fmt"
	"os"
	"testing"

	"github.com/wendisx/gorchat/internal/constant"
	"github.com/wendisx/gorchat/internal/log"
	"github.com/wendisx/gorchat/model"
	r "github.com/wendisx/gorchat/repository"
)

const (
	connect_url = "gochat:gochat@tcp(100.86.230.59:3306)/im?charset=utf8mb4&timeout=30s&loc=Local"
)

var (
	logger log.Logger
	repo   r.UserRepository
	usc    UserUsecase
)

func TestMain(m *testing.M) {
	db := r.NewMysqlDB(connect_url)
	logger = log.NewLogger(constant.DEBUG).Sugar()
	repo = r.NewUserRepository(db, logger)
	usc = NewUserUsecase(repo)
	defer db.Close()

	code := m.Run()
	os.Exit(code)
}

func TestSignup(t *testing.T) {
	t.Skip()
	t.Parallel()
	user := []model.User{
		{
			UserName: "tom",
			Password: "tomlovegolang",
			Email:    "tom@gamil.com",
		},
		{
			UserName: "mike",
			Password: "mikelovevim",
			Email:    "mike@gamil.com",
		},
		{
			UserName: "lily",
			Password: "lilylovevim",
			Email:    "lily@gamil.com",
		},
	}
	userId1, err := usc.Signup(user[0])
	if err != nil {
		fmt.Printf("test -signup- fail: %s\n", err.Error())
	}
	userId2, err := usc.Signup(user[1])
	if err != nil {
		fmt.Printf("test -signup- fail: %s\n", err.Error())
	}
	userId3, err := usc.Signup(user[2])
	if err != nil {
		fmt.Printf("test -signup- fail: %s\n", err.Error())
	}
	fmt.Printf("test -signup- success userIds: [%d %d %d]\n", userId1, userId2, userId3)
}

func TestLogin(t *testing.T) {
	t.Skip()
	t.Parallel()
	sUser := model.User{
		UserId:   2,
		UserName: "tom",
		Password: "tomlovegolang",
		Email:    "tom@gamil.com",
	}
	fUser := model.User{
		UserId:   0,
		UserName: "tom",
		Password: "tomjust",
		Email:    "tom@gamil.com",
	}
	// login success
	sUser, err := usc.Login(sUser)
	if err != nil {
		fmt.Printf("test -login- fail: %s\n", err.Error())
	} else {
		fmt.Printf("test -login- success: userId: %d userName: %s\n", sUser.UserId, sUser.UserName)
	}
	// login fail
	fUser, err = usc.Login(fUser)
	if err != nil {
		fmt.Printf("test -login- fail: %s\n", err.Error())
	}
}

func TestUpdateInfo(t *testing.T) {
	t.Skip()
	t.Parallel()
	user := model.User{
		UserId:   2,
		UserName: "tomchange",
		Password: "tomlovegolang",
		Email:    "tom123@gamil.com",
	}
	nUser, err := usc.UpdateInfo(user)
	if err != nil {
		fmt.Printf("test -updateInfo- fail: %s", err.Error())
	} else {
		fmt.Printf("test -updateInfo- success: userId: %d userName: %s userEmail: %s", nUser.UserId, nUser.UserName, nUser.Email)
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	suserId := 5
	fuserId := 1
	// delete success
	err := usc.Delete(int64(suserId))
	if err != nil {
		fmt.Printf("test -delete- success: userId: %d", suserId)
	}
	// delete fail
	err = usc.Delete(int64(fuserId))
	if err != nil {
		fmt.Printf("test -delete- fail: userId: %d", fuserId)
	}
}
