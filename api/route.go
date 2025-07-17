package api

import (
	"log"

	"github.com/wendisx/gorchat/handler"
	"github.com/wendisx/gorchat/model"
	"github.com/wendisx/gorchat/repository"
	"github.com/wendisx/gorchat/usecase"
)

const (
	GROUP_USER = "/user"
)

func SetupRoute(dependency *model.Dependency) {
	defer log.Printf("[init] -- (api/route) status: success\n")
	registerUserRoute(dependency)
}

func registerUserRoute(dep *model.Dependency) {
	defer log.Printf("[init] -- (api/route) status: success\n")
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
