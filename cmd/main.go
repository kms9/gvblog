package main

import (
	"fmt"
	"github.com/kms9/gvblog/conf"
	"github.com/kms9/gvblog/model"
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

	logs.SetConsole(true)

	conf.Init()
	model.Init()

	fmt.Println("Dsn:", conf.App.Dsn())
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	logs.Info("app running")

	//test db
	fmt.Println("test db")

	model.AllGetTest()

	//结束应用
	<-quit
	logs.Info("app quitted")
	logs.Flush()
}
