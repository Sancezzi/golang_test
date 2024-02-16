package model

import (
	"github.com/guregu/null"
)

//go:generate reform

// reform:News
type Article struct {
	Id      int64       `reform:"Id,pk"`
	Title   null.String `reform:"Title"`
	Content null.String `reform:"Content"`
}

// reform:NewsCategories
type NewsCategory struct {
	NewsId     int64 `reform:"NewsId"`
	CategoryId int64 `reform:"CategoryId"`
}

func (a Article) UpdatedFields() []string {
	fields := make([]string, 0, 2)
	if a.Title.Valid {
		fields = append(fields, "Title")
	}
	if a.Content.Valid {
		fields = append(fields, "Content")
	}

	return fields
}
