package model

import (
	"time"
)

// User 用户
type User struct {
	Id       int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT;comment:主键" json:"id"`
	Name     string    `gorm:"column:name;type:varchar(255);comment:姓名" json:"name"`
	Num      string    `gorm:"column:num;type:varchar(255);comment:账号" json:"num"`
	Passwd   string    `gorm:"column:passwd;type:varchar(255);comment:密码" json:"passwd"`
	Email    string    `gorm:"column:email;type:varchar(255);comment:邮箱" json:"email"`
	Phone    string    `gorm:"column:phone;type:varchar(255);comment:电话" json:"phone"`
	Ecount   int       `gorm:"column:ecount;type:int(11);default:0;comment:错误次数" json:"ecount"`
	Ltime    time.Time `gorm:"column:ltime;type:datetime;comment:上次登录时间" json:"ltime"`
	Ctime    time.Time `gorm:"column:ctime;type:datetime;comment:创建时间" json:"ctime"`
	OpenidQq string    `gorm:"column:openid_qq;type:varchar(64);comment:qq_openid" json:"openid_qq"`
}

func (m *User) TableName() string {
	return "sys_user"
}

//UserLogin 用户登录
func UserLogin(num string) (*User, error) {
	mod := &User{}
	if err = db.Model(&User{}).Where("num = ?", num).First(&mod).Error; err != nil {
		return mod, err
	}
	return mod, nil
}

// UserGet 单条用户信息
func UserGet(id int) (*User, error) {
	mod := &User{}
	if err := db.Model(&User{}).Where(&User{Id: id}).First(mod).Error; err != nil {
		return nil, err
	}
	return mod, nil
}

// UserEdit 编辑用户信息
func UserEdit(mod *User, cols ...string) error {

	if err := db.Model(&User{}).Where(&User{Id: mod.Id}).Updates(mod).Error; err != nil {
		return err
	}
	return nil
}
