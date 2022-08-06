package model

import (
	"github.com/go-programming-tour/blog-service/pkg/app"
	"github.com/jinzhu/gorm"
)

type Tag struct {
	*Common        // 匿名结构体
	Name    string `json:"name"`
	State   uint8  `json:"state"`
}

// swagger结构体
type TagSwagger struct {
	List  []*Tag
	Paper *app.Pager
}

// 用于 gorm 框架返回表名
func (t *Tag) TableName() string {
	return "blog_tag"
}

// 返回指定条件的标签数量
func (t *Tag) Count(db *gorm.DB) (int, error) {
	var count int
	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	db = db.Where("state = ?", t.State)

	if err := db.Model(&t).Where("is_del = ?", 0).Count(&count).Error; err != nil {
		return 0, nil
	}

	return count, nil
}

// 返回分页查询的标签集合
func (t *Tag) List(db *gorm.DB, pageOffset, pageSize int) ([]*Tag, error) {
	var tags []*Tag
	// 对 SQL 语句的参数进行校验
	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}

	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	db = db.Where("state = ?", t.State)
	if err := db.Where("is_del = ?", 0).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (t *Tag) Create(db *gorm.DB) error {
	return db.Create(&t).Error
}

func (t *Tag) Update(db *gorm.DB, values interface{}) error {
	if err := db.Model(&Tag{}).Where("id = ? AND is_del = ?", t.Common.ID, 0).Update(values).Error; err != nil {
		return err
	}
	return nil
}

func (t *Tag) Delete(db *gorm.DB) error {
	return db.Where("id = ? AND is_del = ?", t.Common.ID, 0).Delete(t).Error
}
