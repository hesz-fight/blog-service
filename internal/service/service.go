package service

import (
	"context"

	"github.com/go-programming-tour/blog-service/global"
	"github.com/go-programming-tour/blog-service/internal/dao"
)

// service 对象
type Service struct {
	ctx context.Context
	dao *dao.Dao
}

// service 对象的构造函数
func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = dao.New(global.DBEngine)

	return svc
}
