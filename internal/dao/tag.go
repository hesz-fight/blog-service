package dao

import (
	"github.com/go-programming-tour/blog-service/internal/model"
	"github.com/go-programming-tour/blog-service/pkg/app"
)

// 返回 Tag 的数量
func (d *Dao) CountTag(name string, state uint8) (int, error) {
	tag := model.Tag{Name: name, State: state}

	return tag.Count(d.engine)
}

// 返回 Tag 标签的分页数据
func (d *Dao) GetTagList(name string, state uint8, page int, pageSize int) ([]*model.Tag, error) {
	tag := model.Tag{Name: name, State: state}

	pageOffset := app.GetPageOffset(page, pageSize)

	return tag.List(d.engine, pageOffset, pageSize)
}

// 创建新的 Tag
func (d *Dao) CreateTag(name string, state uint8, cratedBy string) error {
	tag := model.Tag{
		Name:   name,
		State:  state,
		Common: &model.Common{CreatedBy: cratedBy},
	}

	return tag.Create(d.engine)
}

// 修改 Tag 的信息
func (d *Dao) UpdateTag(id uint32, name string, state uint8, modifiedBy string) error {
	tag := model.Tag{
		Common: &model.Common{ID: id},
	}

	values := map[string]interface{}{
		"state":      state,
		"modifiedBy": modifiedBy,
	}
	if name != "" {
		values["name"] = name
	}

	return tag.Update(d.engine, values)
}

// 删除某个 id 的 Tag
func (d *Dao) DeleteTag(id uint32) error {
	tag := model.Tag{Common: &model.Common{ID: id}}

	return tag.Delete(d.engine)
}
