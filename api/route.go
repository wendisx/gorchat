package api

import (
	"log"

	"github.com/wendisx/gorchat/handler"
	"github.com/wendisx/gorchat/model"
	"github.com/wendisx/gorchat/repository"
	"github.com/wendisx/gorchat/usecase"
)

const (
	GROUP_USER   = "/user"
	GROUP_SINGLE = "/single"
)

func SetupRoute(dependency *model.Dependency) {
	defer log.Printf("[init] -- (api/route) status: success\n")
	registerUserRoute(dependency)
	registerSingleRoute(dependency)
}

func registerUserRoute(dep *model.Dependency) {
	defer log.Printf("[init] -- (api/route/user) status: success\n")
	g := dep.Echo.Group(GROUP_USER)

	userRepo := repository.NewUserRepository(dep.Database, dep.Logger)
	userCase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userCase, dep.Response)

	g.POST("/signup", userHandler.Signup, dep.MiddleWare.ValidatorMiddleware(&model.SignupReq{}))
	g.GET("/login", userHandler.Login, dep.MiddleWare.ValidatorMiddleware(&model.LoginReq{}), dep.MiddleWare.SessionCheckMiddleware(true))
	g.PUT("/update", userHandler.UpdateInfo, dep.MiddleWare.ValidatorMiddleware(&model.UpdateInfoReq{}), dep.MiddleWare.SessionCheckMiddleware(false))
	g.DELETE("/delete", userHandler.Delete, dep.MiddleWare.SessionCheckMiddleware(false))
	g.GET("/detail", userHandler.GetUserdetail, dep.MiddleWare.SessionCheckMiddleware(false))
	g.GET("/search", userHandler.SearchUser, dep.MiddleWare.ValidatorMiddleware(&model.SearchUserReq{}), dep.MiddleWare.SessionCheckMiddleware(false))
}

func registerSingleRoute(dep *model.Dependency) {
	defer log.Printf("[init] -- (api/route/single) status: success")
	g := dep.Echo.Group(GROUP_SINGLE)

	singleRepo := repository.NewSingleRepository(dep.Database, dep.Logger)
	singleCase := usecase.NewSingleUsercase(singleRepo)
	singleHandler := handler.NewSingleHandler(singleCase, dep.Response)

	g.Use(dep.MiddleWare.SessionCheckMiddleware(false))

	g.POST("/invite", singleHandler.Invite, dep.MiddleWare.ValidatorMiddleware(&model.InviteReq{}))
	g.PATCH("/accept", singleHandler.Accept, dep.MiddleWare.ValidatorMiddleware(&model.AcceptReq{}))
	g.PATCH("/setNickname", singleHandler.UpdateNickname, dep.MiddleWare.ValidatorMiddleware(&model.UpdateNicknameReq{}))
	g.PATCH("/setDisturb", singleHandler.UpdateDisturb, dep.MiddleWare.ValidatorMiddleware(&model.UpdateDisturbReq{}))
	g.GET("/detail", singleHandler.GetDetail, dep.MiddleWare.ValidatorMiddleware(&model.GetDetailReq{}))
	g.DELETE("/delete", singleHandler.Delete, dep.MiddleWare.ValidatorMiddleware(&model.DeleteReq{}))
}
