package app

import (
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour/blog-service/global"
	"github.com/go-programming-tour/blog-service/pkg/convert"
)

// 分页处理
// 获取分页编号
func GetPage(c *gin.Context) int {
	s := convert.StrTo(c.Query("page"))
	page := s.MustInt()
	if page <= 0 {
		return 1
	}

	return page
}

// 获取页大小
func GetPageSize(c *gin.Context) int {
	s := convert.StrTo(c.Query("page_size"))
	pageSize := s.MustInt()
	if pageSize <= 0 {
		return global.AppSetting.DefaultPageSize
	}

	if pageSize > global.AppSetting.MaxPageSize {
		return global.AppSetting.MaxPageSize
	}

	return pageSize
}

// 获取页偏移
func GetPageOffset(page, pageSize int) int {
	ret := 0
	if page > 0 {
		ret = (page - 1) * pageSize
	}

	return ret
}
