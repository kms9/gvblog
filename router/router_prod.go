//go:build prod
// +build prod

package router

import (
	"github.com/kms9/gvblog/libs/logs"
)

const AppJsUrl = "/static/js/app.min.js"
const AppCssUrl = "/static/css/app.min.css"

func init() {
	logs.SetLevel(logs.WARN)
}

/*  正式模式 编译 取消文档
 *  生成文档 swag init
 *  go build -tags=prod -o blog.exe .\main.go
 *
 *  开发模式 编译 添加文档
 *  go build -o blogdev.exe .\main.go
 */
