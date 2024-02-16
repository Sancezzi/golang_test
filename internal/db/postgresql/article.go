package postgresql

import (
	"context"
	"fmt"
	"test-api/internal/common/errs"
	"test-api/internal/db"
	"test-api/internal/db/postgresql/model"

	"gopkg.in/reform.v1"
)

type articleRepository struct {
	db *db.DbManager
}

func NewArticleRepository(db *db.DbManager) *articleRepository {
	return &articleRepository{
		db: db,
	}
}

func (r *articleRepository) FetchList(ctx context.Context) ([]model.Article, error) {
	db := r.db.GetConnection()
	rows, err := db.WithContext(ctx).SelectRows(model.ArticleTable, "")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]model.Article, 0)
	for {
		var article model.Article
		if err = db.NextRow(&article, rows); err != nil {
			break
		}
		res = append(res, article)
	}

	if err != reform.ErrNoRows {
		return nil, err
	}

	return res, nil
}

func (r *articleRepository) Update(ctx context.Context, article model.Article) (model.Article, error) {
	updatedFileds := article.UpdatedFields()
	err := r.db.GetConnection().WithContext(ctx).UpdateColumns(&article, updatedFileds...)
	if err != nil {
		if err == reform.ErrNoRows {
			return article, errs.ErrNotFound
		}
		return article, err
	}

	return article, nil
}

func (r *articleRepository) FetchArticleCategories(ctx context.Context, id int64) ([]int64, error) {
	db := r.db.GetConnection()

	tail := fmt.Sprintf("WHERE \"NewsId\" = %s", db.Placeholder(1))

	rows, err := db.WithContext(ctx).SelectRows(model.NewsCategoryView, tail, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]int64, 0)
	for {
		var cat model.NewsCategory
		if err = db.NextRow(&cat, rows); err != nil {
			break
		}
		res = append(res, cat.CategoryId)
	}

	if err != reform.ErrNoRows {
		return nil, err
	}

	return res, nil
}

func (r *articleRepository) UpdateCategories(ctx context.Context, id int64, cats []int64) error {
	db := r.db.GetConnection()

	tail := fmt.Sprintf("WHERE \"NewsId\" = %s", db.Placeholder(1))

	_, err := db.WithContext(ctx).DeleteFrom(model.NewsCategoryView, tail, id)
	if err != nil {
		return err
	}

	rows := make([]reform.Struct, 0, len(cats))
	for _, v := range cats {
		rows = append(rows, &model.NewsCategory{
			NewsId:     id,
			CategoryId: v,
		})
	}

	err = db.WithContext(ctx).InsertMulti(rows...)
	if err != nil {
		return err
	}

	return nil
}
