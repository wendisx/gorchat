package main

import (
	lg "log"

	"github.com/wendisx/gorchat/api"
)

func main() {
	// 项目启动显示
	startup()
	// 初始化依赖
	e, _, dep, addr := setup()
	// 初始化路由
	api.SetupRoute(dep)
	// 启动server
	lg.Printf("[start] -- (cmd/main) status: success loc: http://%s", addr)
	e.Start(addr)
	// 释放资源，合理退出
	defer func() {
		// 日志写回
		dep.Logger.Sync()
	}()
}
