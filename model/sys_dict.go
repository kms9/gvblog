package model

import "github.com/kms9/gvblog/libs/logs"

// Dict 配置
type Dict struct {
	Key   string `gorm:"column:key;type:varchar(64);primary_key;comment:唯一Key" json:"key"`
	Value string `gorm:"column:value;type:varchar(255);comment:值" json:"value"`
	Inner bool   `gorm:"column:inner;type:tinyint(4);comment:内部禁止删除" json:"inner"`
	Intro string `gorm:"column:intro;type:varchar(255);comment:说明" json:"intro"`
}

func (m *Dict) TableName() string {
	return "sys_dict"
}

// DictGet 单条字典
func DictGet(key string) (*Dict, error) {
	mod := &Dict{}
	if err := db.Model(&Dict{}).Where(&Dict{Key: key}).First(mod).Error; err != nil {
		return nil, err
	}
	return mod, nil
}

// DictPage 字典分页
func DictPage(pi int, ps int, cols ...string) ([]Dict, error) {
	mods := make([]Dict, 0, ps)

	if err := db.Model(&Dict{}).Limit(ps).Offset((pi - 1) * ps).Find(&mods).Error; err != nil {
		return nil, err
	}

	return mods, err
}

// DictCount 字典分页总数
func DictCount() int {

	var count int64
	if err := db.Model(&Dict{}).Count(&count).Error; err != nil {
		logs.Error(err)
		return 0
	}
	return int(count)
}

// DictAdd 添加字典
func DictAdd(mod *Dict) error {

	if err := db.Model(&Dict{}).Save(mod).Error; err != nil {
		return err
	}
	return nil
}

// DictEdit 编辑字典
func DictEdit(mod *Dict, cols ...string) error {

	if err := db.Model(&Dict{}).Where(&Dict{Key: mod.Key}).Updates(mod).Error; err != nil {
		return err
	}
	return nil
}

// DictDrop 删除单条字典
func DictDrop(key string) error {
	if err := db.Model(&Dict{}).Where(&Dict{Key: key}).Delete(&Dict{}).Error; err != nil {
		return err
	}
	return nil
}
