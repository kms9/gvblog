package router

import (
	"github.com/kms9/gvblog/conf"
	"github.com/kms9/gvblog/internal/token"
	"github.com/kms9/gvblog/utils"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

var funcMap template.FuncMap

func init() {
	funcMap = template.FuncMap{"str2html": Str2html, "str2js": Str2js, "date": Date, "md5": Md5}
}

//func midRecover(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(ctx *gin.Context)  {
//		defer func() {
//			if r := recover(); r != nil {
//				err, ok := r.(error)
//				if !ok {
//					err = fmt.Errorf("%v", r)
//				}
//				stack := make([]byte, 1<<10)
//				length := runtime.Stack(stack, false)
//				// stdlog.Println(string(stack[:length]))
//				os.Stdout.Write(stack[:length])
//				ctx.Error(err)
//			}
//		}()
//		return next(ctx)
//	}
//}

// HTTPErrorHandler 全局错误捕捉
//func HTTPErrorHandler(err error, ctx *gin.Context) {
//	if !ctx.Response().Committed {
//		if he, ok := err.(*echo.HTTPError); ok {
//			if he.Code == 404 {
//				if strings.HasPrefix(ctx.Request().URL.Path, "/static") || strings.HasPrefix(ctx.Request().URL.Path, "/dist") {
//					ctx.NoContent(404)
//				} else if strings.HasPrefix(ctx.Request().URL.Path, "/api") || strings.HasPrefix(ctx.Request().URL.Path, "/adm") {
//					ctx.JSON(utils.NewErrSvr("系统错误", he.Message))
//				} else {
//					ctx.HTML(404, html404)
//				}
//			} else {
//				ctx.JSON(utils.NewErrSvr("系统错误", he.Message))
//			}
//		} else {
//			ctx.JSON(utils.NewErrSvr("系统错误", err.Error()))
//		}
//	}
//}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 先调用c.Next()执行后面的中间件
		// 所有中间件及router处理完毕后从这里开始执行
		// 检查c.Errors中是否有错误
		for _, e := range c.Errors {
			err := e.Err
			// 若是自定义的错误则将code、msg返回
			// 若非自定义错误则返回详细错误信息err.Error()
			// 比如save session出错时设置的err
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器异常",
				"data": err.Error(),
			})
			return // 检查一个错误就行
		}
	}
}

// 跨越配置
//var crosConfig = middleware.CORSConfig{
//	AllowOrigins: []string{"*"},
//	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
//}

// Str2html Convert string to template.HTML type.
func Str2html(raw string) template.HTML {
	return template.HTML(raw)
}

// Str2js Convert string to template.JS type.
func Str2js(raw string) template.JS {
	return template.JS(raw)
}

// Date Date
func Date(t time.Time, format string) string {
	return t.Format(format) //"2006-01-02 15:04:05"
}

// Md5 Md5
func Md5(str string) string {
	ctx := md5.New()
	ctx.Write([]byte(str))
	return hex.EncodeToString(ctx.Sum(nil))
}

func createMyRender() multitemplate.Renderer {
	r := loadTemplates("./views")
	r.AddFromFilesFuncs("dist/index.html", funcMap, "dist/index.html")
	return r
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	pages, err := filepath.Glob(templatesDir + "/page/*.html")
	if err != nil {
		panic(err.Error())
	}

	tpls, err := filepath.Glob(templatesDir + "/tpl/*.html")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, page := range pages {
		tplsCopy := make([]string, 0)
		files := append(tplsCopy, page)
		files = append(files, tpls...)
		r.AddFromFilesFuncs(filepath.Base(page), funcMap, files...)
	}

	return r
}

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		tokenRaw := ctx.Query("token") // query/form 查找 token
		if tokenRaw == "" {
			tokenRaw = ctx.Request.Header.Get("Authorization") // header 查找token
			if tokenRaw == "" {
				ctx.JSON(utils.ErrJwt(`请重新登陆`, `未发现jwt`))
				ctx.Abort()
				return
			}
		}
		auth := token.Auth{}
		err := auth.Decode(tokenRaw, conf.App.TokenSecret)
		if err != nil {
			ctx.JSON(utils.ErrJwt(`请重新登陆`, err.Error()))
			ctx.Abort()
			return
		}
		// 验证通过，保存信息
		ctx.Set("auth", auth)
		ctx.Set("uid", auth.Id)
		ctx.Set("rid", auth.RoleId)
		// 后续流程
		ctx.Next()
	}
}
