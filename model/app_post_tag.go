package model

import "github.com/kms9/gvblog/libs/logs"

//PostTag 文章标签
type PostTag struct {
	Id     int   `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	PostId int   `gorm:"column:post_id;type:int(11)" json:"post_id"`
	TagId  int   `gorm:"column:tag_id;type:int(11)" json:"tag_id"`
	Post   *Post `gorm:"-" swaggerignore:"true" json:"post"` //文章
	Tag    *Tag  `gorm:"-" swaggerignore:"true" json:"tag"`  //标签
}

func (m *PostTag) TableName() string {
	return "post_tag"
}

// TagPostCount 通过标签查询文章分页总数
func TagPostCount(tagId int) int {
	var count int
	if err := db.Raw(`SELECT count(post_id) as count FROM post_tag WHERE tag_id=?`, tagId).Scan(&count).Error; err != nil {
		logs.Error(err)
		return 0
	}
	return count
}

// TagPostPage 通过标签查询文章分页
func TagPostPage(tagId, pi, ps int) ([]PostTag, error) {
	mods := make([]PostTag, 0, ps)
	err := db.Model(&PostTag{}).Raw(`SELECT id,post_id,tag_id FROM post_tag WHERE tag_id=? ORDER BY post_id DESC LIMIT ?,? `, tagId, (pi-1)*ps, ps).Scan(&mods).Error
	if len(mods) > 0 {
		ids := make([]int, 0, len(mods))
		for idx := range mods {
			if !inOf(mods[idx].PostId, ids) {
				ids = append(ids, mods[idx].PostId)
			}
		}
		mapSet := PostIds(ids)
		for idx := range mods {
			mods[idx].Post = mapSet[mods[idx].PostId]
		}
	}
	return mods, err
}

// TagPostAdds 添加文章标签[]
func TagPostAdds(mods *[]PostTag) error {

	if err := db.Model(&PostTag{}).Save(mods).Error; err != nil {
		return err
	}
	return nil
}

// TagPostDrop 删除标签时候
// 删除对应的标签_文章
func TagPostDrop(tagId int) error {
	if err := db.Model(&PostTag{}).Where(&PostTag{TagId: tagId}).Delete(&PostTag{}).Error; err != nil {
		return err
	}
	return nil
}

// PostTagDrops 修改文章时候
// 删除对应的标签_文章
func PostTagDrops(postId int, tagIds []int) error {
	if len(tagIds) < 1 {
		return nil
	}

	if err := db.Model(&PostTag{}).Where(&PostTag{PostId: postId}).Where("tag_id IN ? ", tagIds).Delete(&PostTag{}).Error; err != nil {
		return err
	}
	return nil

}

// PostTagDrop 删除文章时候
// 删除对应的标签_文章
func PostTagDrop(postId int) error {

	if err := db.Model(&PostTag{}).Where(&PostTag{PostId: postId}).Delete(&PostTag{}).Error; err != nil {
		return err
	}
	return nil
}
