package appctl

import (
	"github.com/gin-gonic/gin"
	"github.com/kms9/gvblog/control/base"
	"github.com/kms9/gvblog/libs/logs"
	"github.com/kms9/gvblog/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/kms9/gvblog/utils"
)

// ------------------------------------------------------ 主页面 ------------------------------------------------------

// ViewIndex 主页面
func ViewIndex(ctx *gin.Context) {
	pi, _ := strconv.Atoi(ctx.PostForm("page"))
	if pi == 0 {
		pi = 1
	}
	ps := model.Gcfg().PageSize
	mods, _ := model.PostPage(-1, model.PostKindPost, pi, ps, "id", "title", "path", "created", "summary")
	total := model.PostCount(-1, model.PostKindPost)
	naver := model.Naver{}
	if pi > 1 {
		naver.Prev = "/?page=" + strconv.Itoa(pi-1)
	}
	if total > (pi * ps) {
		naver.Next = "/?page=" + strconv.Itoa(pi+1)
	}
	ctx.HTML(http.StatusOK, "index.html", base.NewR().Add(map[string]interface{}{
		"Posts": mods,
		"Naver": naver,
	}))
}

// ------------------------------------------------------ 文章页面 ------------------------------------------------------
// ViewPost 文章页面
func ViewPost(ctx *gin.Context) {
	paths := strings.Split(ctx.Param("id"), ".")
	if len(paths) == 2 {
		mod, naver, err := model.PostPath(paths[0])
		if err != nil {
			logs.Error(err)
			ctx.Redirect(302, "/")
		}
		if paths[1] == "html" {
			mod.Richtext = regImg.ReplaceAllString(mod.Richtext, `<img class="lazy-load" src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==" data-src="$1" alt="$2">`)
			ctx.HTML(http.StatusOK, "post.html", base.NewR().Add(map[string]interface{}{
				"Post":  mod,
				"Show":  mod.Status == 2,
				"Naver": naver,
			}))
		}
		ctx.JSON(utils.Succ("", mod))
	}
	ctx.Redirect(302, "/404")
}

// ------------------------------------------------------ 单个页面 ------------------------------------------------------
// ViewAbout 关于页面
func ViewAbout(ctx *gin.Context) {
	RenderPage("about", ctx)
}

// ViewLinks 友链页面
func ViewLinks(ctx *gin.Context) {
	RenderPage("links", ctx)
}

// ViewPage 其它页面
func ViewPage(ctx *gin.Context) {
	paths := strings.Split(ctx.Param("id"), ".")
	if len(paths) == 2 {
		RenderPage(paths[0], ctx)
		return
	}
	ctx.Redirect(302, "/404")

}

// ------------------------------------------------------ 归档页面 ------------------------------------------------------
// ViewArchives 归档页面
func ViewArchives(ctx *gin.Context) {
	archives, err := model.PostArchive()
	if err != nil {
		ctx.Redirect(302, "/")
	}
	ctx.HTML(http.StatusOK, "archive.html", base.NewR().Add(map[string]interface{}{
		"Archives": archives,
	}))
}

// ------------------------------------------------------ 分类页面 ------------------------------------------------------
// ViewCatePost 分类文章列表
func ViewCatePost(ctx *gin.Context) {
	cate := ctx.Param("cate")
	if cate == "" {
		ctx.Redirect(302, "/")
	}
	mod, err := model.CateGet(cate)
	if err != nil {
		logs.Error(err)
		ctx.Redirect(302, "/")
	}
	pi, _ := strconv.Atoi(ctx.PostForm("page"))
	if pi == 0 {
		pi = 1
	}
	ps := model.Gcfg().PageSize
	mods, err := model.PostPage(mod.Id, model.PostKindPost, pi, ps, "id", "title", "path", "created", "summary", "updated", "status")
	if err != nil || len(mods) < 1 {
		ctx.Redirect(302, "/")
	}
	total := model.PostCount(mod.Id, model.PostKindPost)
	naver := model.Naver{}
	if pi > 1 {
		naver.Prev = "/cate/" + mod.Name + "?page=1"
	}
	if total > (pi * ps) {
		naver.Next = "/cate/" + mod.Name + "?page=" + strconv.Itoa(pi+1)
	}
	ctx.HTML(http.StatusOK, "post-cate.html", base.NewR().Add(map[string]interface{}{
		"Cate":      mod,
		"CatePosts": mods,
		"Naver":     naver,
	}))
}

// ------------------------------------------------------ 标签页面 ------------------------------------------------------
// ViewTags 标签页面
func ViewTags(ctx *gin.Context) {
	mods, err := model.TagStateAll()
	if err != nil {
		ctx.Redirect(302, "/")
	}
	ctx.HTML(http.StatusOK, "tags.html", base.NewR().Add(map[string]interface{}{
		"Tags": mods,
	}))
}

// ViewTagPost 标签下的文章列表
func ViewTagPost(ctx *gin.Context) {
	tag := ctx.Param("tag")
	if tag == "" {
		ctx.Redirect(302, "/tags")
	}
	mod, err := model.TagGet(tag)
	if err != nil {
		logs.Error(err)
		ctx.Redirect(302, "/tags")
	}
	pi, _ := strconv.Atoi(ctx.PostForm("page"))
	if pi == 0 {
		pi = 1
	}
	ps := model.Gcfg().PageSize
	mods, err := model.TagPostPage(mod.Id, pi, ps)
	if err != nil || len(mods) < 1 {
		ctx.Redirect(302, "/tags")
	}
	total := model.TagPostCount(mod.Id)
	naver := model.Naver{}
	if pi > 1 {
		naver.Prev = "/tag/" + mod.Name + "?page=1"
	}
	if total > (pi * ps) {
		naver.Next = "/tag/" + mod.Name + "?page=" + strconv.Itoa(pi+1)
	}
	ctx.HTML(http.StatusOK, "post-tag.html", base.NewR().Add(map[string]interface{}{
		"Tag":      mod,
		"TagPosts": mods,
		"Naver":    naver,
	}))
}
