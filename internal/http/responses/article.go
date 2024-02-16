package responses

import (
	"test-api/internal/service/model"

	"github.com/guregu/null"
)

type Article struct {
	Id         int64       `json:"Id"`
	Title      null.String `json:"Title"`
	Content    null.String `json:"Content"`
	Categories []int64     `json:"Categories"`
}

func NewArticleResponse(article model.Article) Article {
	return Article{
		Id:         article.Id,
		Title:      article.Title,
		Content:    article.Content,
		Categories: article.Categories,
	}
}

func ErrorResponse(err error) map[string]string {
	return map[string]string{
		"error": err.Error(),
	}
}
