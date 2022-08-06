package model

import "github.com/go-programming-tour/blog-service/pkg/app"

// 文章结构体
type Article struct {
	*Common
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	State         uint8  `json:"state"`
}

// swagger结构体
type ArticleSwagger struct {
	List  []*Article
	Paper *app.Pager
}

func (a *Article) TableName() string {
	return "blog_article"
}
