package model

import (
	"github.com/kms9/gvblog/conf"
	"github.com/kms9/gvblog/libs/logs"
	"time"
)

type Post struct {
	Id       int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	CateId   int       `gorm:"column:cate_id;type:int(11);comment:分类Id" json:"cate_id"`
	Kind     int       `gorm:"column:kind;type:int(11);comment:类型1-文章，2-页面" json:"kind"`
	Status   int       `gorm:"column:status;type:int(11);comment:状态1-草稿，2-已发布" json:"status"`
	Title    string    `gorm:"column:title;type:varchar(255);comment:标题" json:"title"`
	Path     string    `gorm:"column:path;type:varchar(255);comment:访问路径" json:"path"`
	Summary  string    `gorm:"column:summary;type:text;comment:摘要" json:"summary"`
	Markdown string    `gorm:"column:markdown;type:mediumtext;comment:markdown内容" json:"markdown"`
	Richtext string    `gorm:"column:richtext;type:mediumtext;comment:富文本内容" json:"richtext"`
	Allow    bool      `gorm:"column:allow;type:tinyint(4);default:1;comment:允许评论" json:"allow"`
	Created  time.Time `gorm:"column:created;type:datetime;comment:创建时间" json:"created"`
	Updated  time.Time `gorm:"column:updated;type:datetime;comment:修改时间" json:"updated"`
	Tags     []Tag     `gorm:"-" json:"tags"` //标签
	Cate     *Cate     `gorm:"-" json:"cate"`
}

func (m *Post) TableName() string {
	return "post"
}

const (
	PostKindPost = 1 //文章
	PostKindPage = 2 //页面
)
const (
	PostStatusDraft  = 1 //草稿
	PostStatusFinish = 2 //完成
)

// PostGet 单条文章/页面
func PostGet(id int) (*Post, error) {
	mod := &Post{}
	// has, _ := db.ID(id).Get(mod)
	err := db.Model(&Post{}).Where(&Post{Id: id}).First(mod).Error
	if err != nil {
		return nil, err
	}
	if err == nil && mod.Kind == PostKindPost {
		tags := make([]Tag, 0, 4)
		if err := db.Raw("SELECT * FROM tag WHERE id IN (SELECT tag_id FROM post_tag WHERE post_id = ?)", mod.Id).Scan(&tags).Error; err != nil {
			logs.Error("find post tag err:", err.Error())
		}
		mod.Tags = tags
	}

	return mod, err
}

// PostAll 所有文章/页面
//func PostAll(cateId int, kind int, cols ...string) ([]Post, error) {
//	mods := make([]Post, 0, 4)
//	sess := db.Model(&Post{})
//	if cateId > 0 {
//		sess.Where("cate_id = ?", cateId)
//	}
//	if kind > 0 {
//		sess.Where("kind = ?", kind)
//	}
//	if len(cols) > 0 {
//		sess.Cols(cols...)
//	}
//	err := db.Desc("created").Find(&mods)
//	return mods, err
//}

//PostExist 判断是否存在
func PostExist(ptah string) bool {
	t := &Post{}
	if err := db.Model(&Post{}).Where(&Post{Path: ptah}).First(&t).Error; err != nil {
		logs.Error(err)
		return false
	}
	return true
}

// PostPage 文章/页面分页
func PostPage(cateId int, kind int, pi int, ps int, cols ...string) ([]Post, error) {
	mods := make([]Post, 0, ps)
	tx := db.Model(&Post{})
	if cateId > 0 {
		tx = tx.Where("cate_id = ?", cateId)
	}
	if kind > 0 {
		tx = tx.Where("kind = ?", kind)
	}
	err := tx.Order("created desc").Limit(ps).Offset((pi - 1) * ps).Find(&mods).Error
	return mods, err
}

// PostCount 返回总数
func PostCount(cateId int, kind int) int {
	var count int64
	tx := db.Model(&Post{})
	if cateId > 0 {
		tx = tx.Where("cate_id = ?", cateId)
	}
	if kind > 0 {
		tx = tx.Where("kind = ?", kind)
	}
	if err := tx.Count(&count).Error; err != nil {
		logs.Error(err)
		return 0
	}
	return int(count)
}

// PostEdit 编辑文章
func PostEdit(mod *Post, cols ...string) error {

	if err := db.Model(&Post{}).Where(&Post{Id: mod.Id}).Updates(mod).Error; err != nil {
		return err
	}
	return nil
}

// PostAdd 添加文章/页面
func PostAdd(mod *Post) error {

	if err := db.Model(&Post{}).Save(mod).Error; err != nil {
		return err
	}
	return nil
}

// PostDrop 删除单条文章
func PostDrop(id int) error {
	if err := db.Model(&Post{}).Where(&Post{Id: id}).Delete(&Post{}).Error; err != nil {
		return err
	}
	return nil
}

// PostIds 通过id集合返回文章
func PostIds(ids []int) map[int]*Post {
	mods := make([]Post, 0, len(ids))
	mapSet := make(map[int]*Post, 0)
	if err := db.Model(&Post{}).Where("id in ? ", ids).Find(&mods); err != nil {
		return mapSet
	}
	for idx := range mods {
		mapSet[mods[idx].Id] = &mods[idx]
	}
	return mapSet
}

// ------------------------------------------------------ 页面使用 ------------------------------------------------------
//PostSingle 单页面 page
func PostSingle(path string) (*Post, error) {
	mod := &Post{}
	if err := db.Model(&Post{}).Where(&Post{Path: path}).First(mod).Error; err != nil {
		return nil, err
	}
	return mod, nil
}

// ------------------------------------------------------ 归档使用 ------------------------------------------------------
// Archive 归档
type Archive struct {
	Time  time.Time `json:"time"`  // 日期
	Posts []Post    `json:"posts"` //文章
}

// PostArchive 归档
func PostArchive() ([]Archive, error) {
	posts := make([]Post, 0, 8)
	err := db.Model(&Post{}).Where("kind = 1  and status = 2 ").Order("created desc").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	mods := make([]Archive, 0, 8)
	for _, v := range posts {
		v.Markdown = ""
		v.Richtext = ""
		v.Summary = ""
		if idx := archInOf(v.Created, mods); idx == -1 { //没有
			mods = append(mods, Archive{
				Time:  v.Created,
				Posts: []Post{v},
			})
		} else { //有
			mods[idx].Posts = append(mods[idx].Posts, v)
		}
	}
	return mods, nil
}

func archInOf(time time.Time, mods []Archive) int {
	for idx := 0; idx < len(mods); idx++ {
		if time.Year() == mods[idx].Time.Year() && time.Month() == mods[idx].Time.Month() {
			return idx
		}
	}
	return -1
}

// PostPath 一条post
func PostPath(path string) (*Post, *Naver, error) {
	mod := &Post{
		Path: path,
		Kind: 1,
	}
	err := db.Model(&Post{}).Where(mod).First(&mod).Error
	if err == nil {
		mod.Cate, _ = CateGet(mod.CateId)
		if mod.Kind == PostKindPost {
			tags := make([]Tag, 0, 4)
			db.Model(&Post{}).Raw("SELECT * FROM tag WHERE id IN (SELECT tag_id FROM post_tag WHERE post_id = ?)", mod.Id).Scan(&tags)
			mod.Tags = tags
		}
		naver := &Naver{}
		p := Post{}
		err = db.Where("kind = 1 and status = 2 and created >?", mod.Created.Format(conf.StdDateTime)).Order("created Asc").First(&p).Error
		if err == nil {
			// <a href="{{.Naver.Prev}}" class="prev">&laquo; 上一页</a>
			naver.Prev = `<a href="/post/` + p.Path + `.html" class="prev">&laquo; ` + p.Title + `</a>`
		}
		n := Post{}
		err = db.Where("kind = 1  and status = 2 and created <?", mod.Created.Format(conf.StdDateTime)).Order("created Desc").First(&n).Error
		if err == nil {
			//<a href="{{.Naver.Next}}" class="next">下一页 &raquo;</a>
			naver.Next = `<a href="/post/` + n.Path + `.html" class="next"> ` + n.Title + ` &raquo;</a>`
		}
		return mod, naver, nil
	}
	return nil, nil, err
}
