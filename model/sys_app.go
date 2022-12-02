package model

import (
	"errors"

	"github.com/kms9/gvblog/libs/logs"
)

// Goinfo go information
type Goinfo struct {
	ARCH    string `json:"arch"`
	OS      string `json:"os"`
	Version string `json:"version"`
	NumCPU  int    `json:"num_cpu"`
}

// ------------------------------------------------------ Global 全局配置 ------------------------------------------------------
// 结构体配置

// Global 全局配置
type Global struct {
	Id          int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT;comment:主键" json:"id"`
	SiteUrl     string `gorm:"column:site_url;type:varchar(255);comment:网站地址" json:"site_url"`
	LogoUrl     string `gorm:"column:logo_url;type:varchar(255);comment:Logo地址" json:"logo_url"`
	Title       string `gorm:"column:title;type:varchar(255);comment:网站标题" json:"title"`
	Keywords    string `gorm:"column:keywords;type:varchar(255);comment:网站关键词" json:"keywords"`
	Description string `gorm:"column:description;type:varchar(255);comment:网站描述" json:"description"`
	FaviconUrl  string `gorm:"column:favicon_url;type:varchar(255);comment:Favicon地址" json:"favicon_url"`
	BeianMiit   string `gorm:"column:beian_miit;type:varchar(255);comment:ICP备案" json:"beian_miit"`
	BeianNism   string `gorm:"column:beian_nism;type:varchar(255);comment:公安备案" json:"beian_nism"`
	Copyright   string `gorm:"column:copyright;type:varchar(255);comment:版权" json:"copyright"`
	SiteJs      string `gorm:"column:site_js;type:varchar(512);comment:全局js" json:"site_js"`
	SiteCss     string `gorm:"column:site_css;type:varchar(512);comment:全局css" json:"site_css"`
	PageSize    int    `gorm:"column:page_size;type:int(11);default:6;comment:分页大小" json:"page_size"`
	Comment     string `gorm:"column:comment;type:varchar(1024);comment:评论脚本" json:"comment"`
	GithubUrl   string `gorm:"column:github_url;type:varchar(255);comment:githu地址" json:"github_url"`
	WeiboUrl    string `gorm:"column:weibo_url;type:varchar(255);comment:微博地址" json:"weibo_url"`
	Analytic    string `gorm:"column:analytic;type:varchar(1024);comment:统计脚本" json:"analytic"`
	Author      string `gorm:"column:author;type:varchar(255);comment:网站作者" json:"author"`
}

func (m *Global) TableName() string {
	return "sys_global"
}

const globalId = 1

var globalCache Global

func initGlobal() error {
	mod := Global{}
	//has, _ := db.ID(1).Get(&mod)
	err := db.Model(&Global{}).Where(&Global{Id: globalId}).First(&mod).Error
	if err != nil {
		return errors.New("no")
	}
	globalCache = mod
	logs.Debug("global cache")
	return nil
}

func GlobalGet() (*Global, error) {
	// mod := &Global{}
	// has, _ := db.ID(globalId).Get(mod)
	return &globalCache, nil
}
func Gcfg() Global {
	cache := globalCache
	return cache
}

// GlobalEdit 编辑global信息
func GlobalEdit(mod *Global, cols ...string) error {
	if err := db.Model(&Global{}).Where(&Global{Id: 1}).Updates(mod).Error; err != nil {
		return err
	}
	initGlobal()
	return nil
}
