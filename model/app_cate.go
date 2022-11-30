package model

import (
	"github.com/silenceper/log"
	"github.com/spf13/cast"
)

// Cate 分类
type Cate struct {
	Id    int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	Name  string `gorm:"column:name;type:varchar(255);comment:分类名" json:"name"`
	Intro string `gorm:"column:intro;type:varchar(255);comment:描述" json:"intro"`
}

func (m *Cate) TableName() string {
	return "cate"
}

// CateGet 单条分类
// int	==>	id
// str	==>	name
func CateGet(id interface{}) (*Cate, error) {
	mod := &Cate{}
	switch val := id.(type) {
	case int:
		err := db.Where(val).Find(mod).Error
		return mod, err
	case string:
		if err := db.Where(&Cate{Name: cast.ToString(id)}).Find(mod).Error; err != nil {
			log.Error(err.Error())
		}
		return mod, nil
	default:
		return mod, nil
	}
}

// CateAll 所有分类
func CateAll() ([]Cate, error) {
	mods := make([]Cate, 0, 8)
	err := db.Find(&mods).Error
	return mods, err
}

// CatePage 分类分页
func CatePage(pi int, ps int, cols ...string) ([]Cate, error) {
	mods := make([]Cate, 0, ps)
	err := db.Model(&Cate{}).Order("id desc").Limit(ps).Offset((pi - 1) * ps).Find(&mods).Error
	return mods, err
}

// CateCount 分类分页总数
func CateCount() int {
	var count int64
	if err := db.Model(&Cate{}).Count(&count).Error; err != nil {
		log.Error(err.Error())
		return 0
	}
	return int(count)
}

// CateIds 通过id集合返回分类
func CateIds(ids []int) map[int]*Cate {
	mods := make([]Cate, 0, len(ids))
	db.Model(&Cate{}).Where("in IN ?", ids).Find(&mods)
	mapSet := make(map[int]*Cate, len(mods))
	for idx := range mods {
		mapSet[mods[idx].Id] = &mods[idx]
	}
	return mapSet
}

// CateAdd 添加分类
func CateAdd(mod *Cate) error {
	if err := db.Model(&Cate{}).Save(mod).Error; err != nil {
		return err
	}
	return nil
}

// CateEdit 编辑分类
func CateEdit(mod *Cate, cols ...string) error {

	if err := db.Model(&Cate{}).Where(&Cate{Id: mod.Id}).Updates(mod).Error; err != nil {
		return err
	}
	return nil
}

// CateDrop 删除单条分类
func CateDrop(id int) error {
	if err := db.Model(&Cate{}).Where(&Cate{Id: id}).Delete(&Cate{}).Error; err != nil {
		return err
	}
	return nil
}
