//go:build !prod
// +build !prod

package router

import (
	"github.com/kms9/gvblog/libs/logs"
	// docs
	//_ "github.com/kms9/gvblog/docs"
	"bytes"
	"sync"
)

var pool *sync.Pool

func init() {
	logs.SetLevel(logs.DEBUG)
	logs.SetCallInfo(true)
	logs.SetConsole(true)
	pool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 512))
		},
	}
}

// midLogger 中间件-日志记录
//func midLogger(next echo.HandlerFunc) gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		start := time.Now()
//
//		stop := time.Now()
//		buf := pool.Get().(*bytes.Buffer)
//		buf.Reset()
//		defer pool.Put(buf)
//		buf.WriteString("\tip：" + ctx.Request.RemoteAddr)
//		buf.WriteString("\tmethod：" + ctx.Request.Method)
//		buf.WriteString("\tpath：" + ctx.Request.RequestURI)
//		buf.WriteString("\tspan：" + stop.Sub(start).String())
//		// 开发模式直接输出
//		// 生产模式中间层会记录
//		// os.Stdout.Write(buf.Bytes())
//		logs.Debug(buf.String())
//		return
//	}
//}
