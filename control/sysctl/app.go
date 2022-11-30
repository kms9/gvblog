package sysctl

import (
	"github.com/gin-gonic/gin"
	"github.com/kms9/gvblog/libs/logs"
	"github.com/kms9/gvblog/model"

	"github.com/kms9/gvblog/utils"
)

// ------------------------------------------------------ 配置中心 ------------------------------------------------------

// GlobalGet doc
// @Tags global-全局配置
// @Summary 获取global信息
// @Success 200 {object} model.Reply{data=model.Global} "返回数据"
// @Router /api/global/get [get]
func GlobalGet(ctx *gin.Context) {
	mod, err := model.GlobalGet()
	if err != nil {
		logs.Error(err)
		ctx.JSON(utils.ErrOpt("未查询到global信息"))
		return
	}
	ctx.JSON(utils.Succ("succ", mod))
	return
}

// GlobalEdit doc
// @Tags global-全局配置
// @Summary 编辑global信息
// @Param token query string true "token"
// @Param body body model.Global true "请求数据"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/global/edit [post]
func GlobalEdit(ctx *gin.Context) {
	ipt := &model.Global{}
	err := ctx.ShouldBindQuery(ipt)
	if err != nil {
		ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
		return
	}
	err = model.GlobalEdit(ipt, "site_url", "logo_url", "title", "author", "keywords", "description", "favicon_url", "beian_miit", "beian_nism", "copyright", "site_js", "site_css", "page_size", "analytic", "comment", "github_url", "weibo_url")
	if err != nil {
		ctx.JSON(utils.Fail("修改失败", err.Error()))
		return
	}
	ctx.JSON(utils.Succ("succ"))
	return
}
