package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kms9/gvblog/conf"
	"github.com/kms9/gvblog/control/appctl"
	"github.com/kms9/gvblog/control/base"
	"github.com/kms9/gvblog/middleware"
	"log"
	"net/http"
)

// RunApp 入口
func RunApp() {
	r := gin.Default()
	r.SetFuncMap(funcMap)
	r.HTMLRender = createMyRender()
	r.Use(middleware.Log())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	// 初始渲染引擎
	//r.Use(midRecover, midLogger)                 // 恢复 日志记录
	//r.Use(middleware.CORSWithConfig(crosConfig)) // 跨域设置
	//r.HideBanner = true                          // 不显示横幅
	//r.HTTPErrorHandler = HTTPErrorHandler // 自定义错误处理
	//r.Debug = conf.App.IsDev()                     // 运行模式 - echo框架好像没怎么使用这个
	//RegDocs(r)                                     // 注册文档
	r.Static(`/dist`, "dist")                   // 静态目录 - 后端专用
	r.Static(`/static`, "static")               // 静态目录
	r.StaticFile(`/favicon.ico`, "favicon.ico") // ico
	// r.Get("/dashboard", "dist/index.html")      // 前后端分离页面
	r.GET("/login.html", func(c *gin.Context) {
		c.Redirect(302, "/dashboard")
	})
	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dist/index.html", nil)
	})

	//--- 页面 -- start
	r.GET("/home", appctl.ViewIndex)        // 首页
	r.GET("/archives", appctl.ViewArchives) // 归档
	//r.GET("/archives.json", appctl.ArchivesJson) // 归档 json
	r.GET("/tags", appctl.ViewTags) // 标签
	//r.GET("/tags.json", appctl.TagsJson)         // 标签 json
	r.GET("/tag/:tag", appctl.ViewTagPost)    // 具体某个标签
	r.GET("/cate/:cate", appctl.ViewCatePost) // 分类
	r.GET("/about", appctl.ViewAbout)         // 关于
	r.GET("/links", appctl.ViewLinks)         // 友链
	r.GET("/post/:id", appctl.ViewPage)       // 具体某个文章
	r.GET("/page/:id", appctl.ViewPost)       // 具体某个页面
	//--- 页面 -- end

	api := r.Group("/api")                      // api/
	apiRouter(api)                              // 注册分组路由
	adm := r.Group("/adm", JWTAuthMiddleware()) // adm/ 需要登陆才能访问
	admRouter(adm)                              // 注册分组路由

	// 未知路由处理
	r.NoRoute(func(c *gin.Context) {
		t := c.Request
		fmt.Println(t)
		if c.Request.RequestURI == "/" {
			c.Request.URL.Path = "/home" //把请求的URL修改
			r.HandleContext(c)           //继续后续处理
		} else {
			c.HTML(404, "404.html", base.NewR().Add(nil))
		}
	})

	err := r.Run(conf.App.Addr)
	if err != nil {
		log.Fatalln("run error :", err)
	}
}
