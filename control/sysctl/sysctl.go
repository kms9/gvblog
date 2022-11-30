package sysctl

import (
	"github.com/gin-gonic/gin"
	"github.com/kms9/gvblog/conf"
	"github.com/kms9/gvblog/model"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/kms9/gvblog/utils"
	"github.com/nfnt/resize"
	_ "golang.org/x/image/bmp"
)

// StatusGoinfo doc
// @Tags status-状态信息
// @Summary 获取服务器go信息
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=model.Goinfo} "返回数据"
// @Router /adm/status/goinfo [get]
func StatusGoinfo(ctx *gin.Context) {
	mod := model.Goinfo{
		ARCH:    runtime.GOARCH,
		OS:      runtime.GOOS,
		Version: runtime.Version(),
		NumCPU:  runtime.NumCPU(),
	}
	ctx.JSON(utils.Succ("系统信息", mod))
	return
}

// StatusApp doc
// @Tags status-状态信息
// @Summary 获取服务器go信息
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=model.State} "返回数据"
// @Router /adm/status/app [get]
func StatusAppinfo(ctx *gin.Context) {
	if mod, err := model.Collect(); err == nil {
		ctx.JSON(utils.Succ(`统计信息`, mod))
		return
	} else {
		ctx.JSON(utils.Fail(`未查询到统计信息`))
		return
	}
}

// UploadFile doc
// @Tags ctrl-系统相关
// @Summary 上传文件
// @Accept  mpfd
// @Param token query string true "token"
// @Param file formData file true "file"
// @Router /adm/upload/file [post]
func UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(utils.Fail("未发现文件", err.Error()))
		return
	}
	src, err := file.Open()
	if err != nil {
		ctx.JSON(utils.Fail("文件打开失败", err.Error()))
		return

	}
	dir := time.Now().Format("200601/02")
	os.MkdirAll("./static/upload/"+dir[:6], 0755)
	name := "static/upload/" + dir + utils.RandStr(10) + path.Ext(file.Filename)
	dst, err := os.Create(name)
	if err != nil {
		ctx.JSON(utils.Fail("文件打创建文件失败", err.Error()))
		return

	}
	_, err = io.Copy(dst, src)
	if err != nil {
		ctx.JSON(utils.Fail("文件保存失败", err.Error()))
		return
	}
	src.Close()
	dst.Close()
	ctx.JSON(utils.Succ("上传成功", "/"+name))
	return
}

// UploadImage doc
// @Tags ctrl-系统相关
// @Summary 上传图片并裁剪
// @Accept  mpfd
// @Param token query string true "token"
// @Param file formData file true "file"
// @Router /adm/upload/image [post]
func UploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(utils.ErrIpt(`未发现文件,请重试`, err.Error()))
		return
	}
	if !strings.Contains(file.Header.Get("Content-Type"), "image") {
		ctx.JSON(utils.ErrIpt("请选择图片文件", file.Header.Get("Content-Type")))
		return
	}
	src, err := file.Open()
	if err != nil {
		ctx.JSON(utils.ErrIpt(`文件打开失败,请重试`, err.Error()))
		return
	}
	defer src.Close()
	attr := &struct {
		Cut bool `json:"cut" query:"cut" form:"cut"`
		Wd  int  `json:"wd" query:"wd" form:"wd"`
		Hd  int  `json:"hd" query:"hd" form:"hd"`
	}{}
	ctx.ShouldBindQuery(attr)
	dir := time.Now().Format("200601/02")
	os.MkdirAll("./static/upload/"+dir[:6], 0755)
	ext := path.Ext(file.Filename)
	if conf.App.ImageCut && attr.Cut {
		ext = ".jpg"
	}
	name := "static/upload/" + dir + utils.RandStr(10) + ext
	dst, err := os.Create(name)
	if err != nil {
		ctx.JSON(utils.ErrIpt("目标文件创建失败,请重试", err.Error()))
		return
	}
	defer dst.Close()
	if conf.App.ImageCut && attr.Cut { //图片裁剪
		imgSrc, _, err := image.Decode(src)
		// 图片解码
		if err != nil {
			ctx.JSON(utils.ErrIpt("读取图片失败,请重试", err.Error()))
			return
		}
		if attr.Wd <= 0 {
			attr.Wd = conf.App.ImageWidth
		}
		if attr.Hd <= 0 {
			attr.Hd = conf.App.ImageHeight
		}
		// 宽度>指定宽度 防止负调整
		dx := imgSrc.Bounds().Dx()
		if dx > attr.Wd {
			imgSrc = resize.Resize(uint(attr.Wd), 0, imgSrc, resize.Lanczos3)
		}
		//高度>指定高度 防止负调整
		dy := imgSrc.Bounds().Dy()
		if dy > attr.Hd {
			imgSrc = resize.Resize(0, uint(attr.Hd), imgSrc, resize.Lanczos3)
		}
		err = jpeg.Encode(dst, imgSrc, nil)
		if err != nil {
			ctx.JSON(utils.ErrIpt("文件写入失败,请重试", err.Error()))
			return
		}
	} else {
		_, err = io.Copy(dst, src)
		if err != nil {
			ctx.JSON(utils.ErrIpt("文件写入失败,请重试", err.Error()))
			return
		}
	}
	ctx.JSON(utils.Succ("上传成功", "/"+name))
	return
}
