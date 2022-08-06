package dao

import "github.com/jinzhu/gorm"

// Dao 实例
type Dao struct {
	engine *gorm.DB
}

// 二次封装 gorm.DB
func New(engine *gorm.DB) *Dao {
	return &Dao{engine: engine}
}
