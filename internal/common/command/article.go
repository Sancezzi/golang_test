package command

import "github.com/guregu/null"

type UpdateArticle struct {
	Id         int64       `json:"Id"`
	Title      null.String `json:"Title" validate:"max=255"`
	Content    null.String `json:"Content"`
	Categories []int64     `json:"Categories"`
}
