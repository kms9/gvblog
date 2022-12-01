package model

import (
	"errors"
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/kms9/gvblog/conf"
	//"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strings"
	"time"

	"github.com/kms9/gvblog/libs/logs"
	// 数据库驱动
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
)

// db 数据库操作句柄
var db *gorm.DB
var err error

func InitSqlite() {

	fmt.Println("SqliteSrc:", conf.App.SqliteSrc())
	db, err = gorm.Open(sqlite.Open(conf.App.SqliteSrc()), &gorm.Config{
		// gorm日志模式：silent
		Logger: logger.Default.LogMode(logger.Info),
		// 外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
	})
	if err != nil {
		logs.Fatal("db sqlite 数据库 dsn:", err.Error())
	}

	sqlDB, _ := db.DB()
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	//初始化数据库
	// 仅在备份数据之后使用
	//err := db.AutoMigrate(
	//	&User{},
	//	&Cate{},
	//	&Post{},
	//	&PostTag{},
	//	&Tag{},
	//	&Global{},
	//	&Dict{},
	//)
	//if err != nil {
	//	logs.Fatal("sqlite db.AutoMigrate err:", err.Error())
	//	return
	//}

	logs.Info("model init by sqlite success")
}

func InitMysql() {
	// 初始化数据库操作的 gorm
	db, err = gorm.Open(mysql.Open(conf.App.Dsn()), &gorm.Config{
		// gorm日志模式：silent
		Logger: logger.Default.LogMode(logger.Info),
		// 外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
	})
	if err != nil {
		logs.Fatal("mysql 数据库 dsn:", err.Error())
	}

	sqlDB, _ := db.DB()
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	logs.Info("model init by mysql success")

}

func Init() {
	switch conf.App.DbKind {
	case "mysql":
		InitMysql()
	case "sqlite":
		InitSqlite()
	}
}

// Page 分页基本数据
type Page struct {
	Pi   int    `json:"pi" form:"pi" query:"pi"`       //分页页码
	Ps   int    `json:"ps" form:"ps" query:"ps"`       //分页大小
	Mult string `json:"mult" form:"mult" query:"mult"` //多条件信息
}

// Trim 去除空白字符
func (p *Page) Trim() string {
	p.Mult = strings.TrimSpace(p.Mult)
	return p.Mult
}

// Stat 检查状态
func (p *Page) Stat() error {
	if p.Ps < conf.App.PageMin {
		return errors.New("page size 过小")
	}
	if p.Ps > conf.App.PageMax {
		return errors.New("page size 过大")
	}
	return nil
}

type IptId struct {
	Id int `form:"id" binding:"required" query:"id" json:"id"` //仅包含Id的请求
}

// Naver 上下页
type Naver struct {
	Prev string
	Next string
}

// State 统计信息
type State struct {
	Post int `json:"post"`
	Page int `json:"page"`
	Cate int `json:"cate"`
	Tag  int `json:"tag"`
}

// Collect 统计信息
func Collect() (*State, error) {
	mod := &State{}
	if err := db.Raw(`SELECT * FROM(SELECT COUNT(id) as post FROM post WHERE kind=1)as a ,(SELECT COUNT(id) as page FROM post WHERE kind=2) as b, (SELECT COUNT(id) as cate FROM cate) as c, (SELECT COUNT(id) as tag FROM tag) as d`).Scan(&mod).Error; err != nil {
		logs.Error(err)
		return mod, err
	}

	return mod, nil
}
func inOf(goal int, arr []int) bool {
	for idx := range arr {
		if goal == arr[idx] {
			return true
		}
	}
	return false
}

// Reply 生成api文档使用
// 代码里未使用，也不要使用
type Reply struct {
	Code int    `json:"code" example:"200"`
	Msg  string `json:"msg" example:"提示信息"`
}
