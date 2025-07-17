package main

import (
	"context"
	"fmt"
	lg "log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wendisx/gorchat/config"
	"github.com/wendisx/gorchat/config/middleware"
	"github.com/wendisx/gorchat/config/redis"
	"github.com/wendisx/gorchat/internal/constant"
	"github.com/wendisx/gorchat/internal/log"
	"github.com/wendisx/gorchat/internal/redistore"
	"github.com/wendisx/gorchat/internal/validator"
	"github.com/wendisx/gorchat/model"
	"github.com/wendisx/gorchat/repository"
)

func startup() {
	startUp := `
			┌─┐┌─┐┬─┐┌─┐┬ ┬┌─┐┌┬┐
			│ ┬│ │├┬┘│  ├─┤├─┤ │ 
			└─┘└─┘┴└─└─┘┴ ┴┴ ┴ ┴ 

   	   A clean architecture go backend for rchat.
	`
	fmt.Println(startUp)
}

func globalErrorHandler(e error, c echo.Context) {
	var err error
	if derr, ok := e.(*model.DError); ok {
		err = c.JSON(
			http.StatusInternalServerError,
			map[string]any{
				"code":    int(derr.Code),
				"message": derr.Message,
				"data":    nil,
			},
		)
		lg.Printf("[error] -- (globalErrorHandler) type: *DError\n")
	} else if derr, ok := e.(validator.ValidatorErrors); ok {
		err = c.JSON(
			http.StatusInternalServerError,
			map[string]any{
				"code":    constant.ErrValidate,
				"message": constant.MsgValidateFail,
				"data":    nil,
			},
		)
		lg.Printf("[error] -- (globalErrorHandler) type: ValidatorErrors\n")
		lg.Printf("[error] -- (globalErrorHandler) %v\n", derr.Error())
	} else if derr, ok := e.(*echo.HTTPError); ok {
		err = c.JSON(
			derr.Code,
			map[string]any{
				"code":    constant.FAIL_CODE,
				"message": derr.Message,
				"data":    nil,
			},
		)
		lg.Printf("[error] -- (globalErrorHandler) type: unknown\n")
	}
	if err != nil {
		lg.Printf("[error] -- (globalErrorHandler) %v\n", err)
	}
}

func setup() (*echo.Echo, config.Env, *model.Dependency, string) {
	// echo -- 路由配置(初始化)
	e := echo.New()
	e.HTTPErrorHandler = globalErrorHandler
	e.HideBanner = true
	// env -- 加载环境变量
	env := config.NewEnv(constant.DEV_ENV_FILE)
	// logger -- 全局日志器
	logger := log.NewLogger(constant.DEBUG)
	sugar := logger.Sugar()
	// redis -- 加载redis
	rdb := redis.NewRedisClient(env)
	rstore := redistore.NewRedistore(context.Background(), rdb)
	// mysql database -- 加载数据库
	db := repository.NewMysqlDB(env[constant.MYSQL_URL])
	// response -- 响应器初始化
	res := model.NewResponser()
	// middleware -- 中间件初始化
	validator := validator.NewValidator()
	md := middleware.NewMiddleware(validator, rstore)

	dep := &model.Dependency{
		Echo:        e,
		Database:    db,
		Logger:      sugar,
		RedisClient: rdb,
		Response:    res,
		MiddleWare:  md,
	}

	// echo -- 服务监听地址
	addr := fmt.Sprintf("%s:%s", env[constant.SERVER_IP], env[constant.SERVER_PORT])
	// 多值初始化返回
	return e, env, dep, addr
}
