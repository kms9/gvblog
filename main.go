package main

import (
	"github.com/kms9/gvblog/conf"
	"github.com/kms9/gvblog/model"
	"github.com/kms9/gvblog/router"
	"os"
	"os/signal"
	"syscall"

	"github.com/kms9/gvblog/libs/logs"
)

// @Title Blog’s Api文档
// @Version 1.0
// @Description token传递方式包括 [get/post]token 、[header] Authorization=Bearer xxxx
// @Description 数据传递方式包括 json、formData 推荐使用 json
// @Description /api/* 公共访问
// @Description /adm/* 必须传入 token
// @Host 127.0.0.1:88
// @BasePath /
func main() {
	logs.Info("app initializing")
	defer logs.Flush()

	conf.Init()
	model.Init()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	logs.Info("app running")
	go router.RunApp()
	<-quit
	logs.Info("app quitted")
	logs.Flush()
}
