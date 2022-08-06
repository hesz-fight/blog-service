package model

type ArticleTag struct {
	*Common
	TageID    uint32 `json:"tag_id"`
	ArticleID uint32 `json:"article_id"`
}

func (a *ArticleTag) TableName() string {
	return "blog_article_tag"
}
