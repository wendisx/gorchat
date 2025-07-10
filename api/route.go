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

	var userDTO model.UserDTO

	userRepo := repository.NewUserRepository(dep.Database, dep.Logger)
	userCase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userCase, dep.Response)

	g.POST("/signup", userHandler.Signup, dep.MiddleWare.ValidatorMiddleware(&userDTO))
	g.POST("/login", userHandler.Login, dep.MiddleWare.ValidatorMiddleware(&userDTO), dep.MiddleWare.SessionCheckMiddleware(true))
	g.PUT("/update", userHandler.UpdateInfo, dep.MiddleWare.ValidatorMiddleware(&userDTO), dep.MiddleWare.SessionCheckMiddleware(false))
	g.DELETE("/delete", userHandler.Delete, dep.MiddleWare.SessionCheckMiddleware(false))
}
