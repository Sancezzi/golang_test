package model

import (
	"test-api/internal/common/command"
	data "test-api/internal/db/postgresql/model"

	"github.com/guregu/null"
)

type Article struct {
	Id         int64
	Title      null.String
	Content    null.String
	Categories []int64
}

func FromUpdate(cmd command.UpdateArticle) Article {
	return Article{
		Id:         cmd.Id,
		Title:      cmd.Title,
		Content:    cmd.Content,
		Categories: cmd.Categories,
	}
}

func NewArticle(data data.Article) Article {
	return Article{
		Id:      data.Id,
		Title:   data.Title,
		Content: data.Content,
	}
}

func (a Article) ToData() data.Article {
	return data.Article{
		Id:      a.Id,
		Title:   a.Title,
		Content: a.Content,
	}
}
