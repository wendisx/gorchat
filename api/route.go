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
	GROUP_GROUP  = "/group"
)

func SetupRoute(dependency *model.Dependency) {
	defer log.Printf("[init] -- (api/route) status: success\n")
	registerUserRoute(dependency)
	registerSingleRoute(dependency)
	registerGroupRoute(dependency)
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

func registerGroupRoute(dep *model.Dependency) {
	defer log.Printf("[init] -- (api/route/group) status: success")
	g := dep.Echo.Group(GROUP_GROUP)

	groupRepo := repository.NewGroupRepository(dep.Database, dep.Logger)
	groupUcase := usecase.NewGroupUsecase(groupRepo)
	groupHandler := handler.NewGroupHandler(groupUcase, dep.Response)

	g.Use(dep.MiddleWare.SessionCheckMiddleware(false))

	g.POST("/create", groupHandler.CreateGroup, dep.MiddleWare.ValidatorMiddleware(&model.CreateGroupReq{}))
	g.POST("/join", groupHandler.JoinGroup, dep.MiddleWare.ValidatorMiddleware(&model.JoinGroupReq{}))
	g.PATCH("/setRole", groupHandler.UpdateGroupUser, dep.MiddleWare.ValidatorMiddleware(&model.UpdateGroupUserReq{}))
	g.PATCH("/setDisturb", groupHandler.UpdateGroupUser, dep.MiddleWare.ValidatorMiddleware(&model.UpdateGroupUserReq{}))
	g.PATCH("/setUserNickname", groupHandler.UpdateGroupUser, dep.MiddleWare.ValidatorMiddleware(&model.UpdateGroupUserReq{}))
	g.PATCH("/setGroupNickname", groupHandler.UpdateGroupUser, dep.MiddleWare.ValidatorMiddleware(&model.UpdateGroupUserReq{}))
	g.PUT("/update", groupHandler.UpdateGroup, dep.MiddleWare.ValidatorMiddleware(&model.UpdateGroupReq{}))
	g.GET("/getAll", groupHandler.GetGroupUsers, dep.MiddleWare.ValidatorMiddleware(&model.GetGroupUsersReq{}))
	g.GET("/search", groupHandler.SearchGroup, dep.MiddleWare.ValidatorMiddleware(&model.SearchGroupReq{}))
	g.GET("/searchUsers", groupHandler.SearchGroupUsers, dep.MiddleWare.ValidatorMiddleware(&model.SearchGroupUsersReq{}))
	g.DELETE("/delete", groupHandler.DeleteGroup)
	g.DELETE("/deleteUser", groupHandler.DeleteGroupUser, dep.MiddleWare.ValidatorMiddleware(&model.DeleteGroupUserReq{}))
}
