package model

import (
	"github.com/kms9/gvblog/libs/logs"
	"github.com/spf13/cast"
)

// Tag 标签

type Tag struct {
	Id    int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	Name  string `gorm:"column:name;type:varchar(255);comment:标签名" json:"name"`
	Intro string `gorm:"column:intro;type:varchar(255);comment:描述" json:"intro"`
}

func (m *Tag) TableName() string {
	return "tag"
}

// TagGet 单条标签
// int	==>	id
// str	==>	name
func TagGet(id interface{}) (*Tag, error) {
	mod := &Tag{}
	switch val := id.(type) {
	case int:
		if err := db.Model(&Tag{}).Where(&Tag{Id: cast.ToInt(val)}).First(mod).Error; err != nil {
			return nil, err
		}
		return mod, nil
	case string:
		if err := db.Model(&Tag{}).Where(&Tag{Name: val}).First(mod).Error; err != nil {
			return nil, err
		}
	}
	return mod, nil
}

// TagAll 所有标签信息
func TagAll() ([]Tag, error) {
	mods := make([]Tag, 0, 8)
	err := db.Model(&Tag{}).Find(&mods).Error
	return mods, err
}

// TagPage 标签分页
func TagPage(pi int, ps int, cols ...string) ([]Tag, error) {
	mods := make([]Tag, 0, ps)
	if err := db.Model(&Tag{}).Order("Id desc").Limit(ps).Offset((pi - 1) * ps).Find(&mods).Error; err != nil {
		return nil, err
	}
	return mods, nil
}

// TagCount 标签分页总数
func TagCount() int {
	var count int64
	if err := db.Model(&Tag{}).Count(&count).Error; err != nil {
		logs.Error(err)
		return 0
	}
	return int(count)
}

// TagIds 通过id集合返回标签
//func TagIds(ids []int) map[int]*Tag {
//	mods := make([]Tag, 0, len(ids))
//	db.In("id", ids).Find(&mods)
//	mapSet := make(map[int]*Tag, len(mods))
//	for idx := range mods {
//		mapSet[mods[idx].Id] = &mods[idx]
//	}
//	return mapSet
//}

// TagAdd 添加标签
func TagAdd(mod *Tag) error {
	if err := db.Model(&Tag{}).Save(mod).Error; err != nil {
		return err
	}
	return nil
}

// TagEdit 编辑标签
func TagEdit(mod *Tag, cols ...string) error {

	if err := db.Model(&Tag{}).Where(&Tag{Id: mod.Id}).Updates(mod).Error; err != nil {
		return err
	}
	return nil
}

// TagDrop 删除单条标签
func TagDrop(id int) error {
	if err := db.Model(&Tag{}).Where(&Tag{Id: id}).Delete(&Tag{}).Error; err != nil {
		return err
	}
	return nil
}

// ------------------------------------------------------ 前台使用 ------------------------------------------------------

// TagState 统计
type TagState struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Intro string `json:"intro"`
}

// TagStateAll 所有标签统计 当前标签下有文章才显示
func TagStateAll() ([]TagState, error) {
	mods := make([]TagState, 0, 8)
	err := db.Raw("SELECT `name`,intro,count(tag_id) as count FROM post_tag ,tag WHERE tag.id=tag_id GROUP BY tag_id HAVING count>0").Scan(&mods).Error
	return mods, err
}
