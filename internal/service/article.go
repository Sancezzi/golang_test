package service

import (
	"context"
	"test-api/internal/common/command"
	data "test-api/internal/db/postgresql/model"
	"test-api/internal/service/model"
)

type NewsRepository interface {
	FetchList(ctx context.Context) ([]data.Article, error)
	Update(ctx context.Context, article data.Article) (data.Article, error)
	FetchArticleCategories(ctx context.Context, id int64) ([]int64, error)
	UpdateCategories(ctx context.Context, id int64, cats []int64) error
}

type newsService struct {
	repo NewsRepository
}

func NewNewsService(repo NewsRepository) *newsService {
	return &newsService{
		repo: repo,
	}
}

func (s *newsService) List(ctx context.Context) ([]model.Article, error) {
	list, err := s.repo.FetchList(ctx)
	if err != nil {
		return nil, err
	}

	news := make([]model.Article, 0, len(list))
	for _, v := range list {
		article := model.NewArticle(v)
		article.Categories, err = s.repo.FetchArticleCategories(ctx, article.Id)
		if err != nil {
			return nil, err
		}

		news = append(news, article)
	}

	return news, nil
}

func (s *newsService) Update(ctx context.Context, cmd command.UpdateArticle, id int64) (model.Article, error) {
	article := model.FromUpdate(cmd)

	data := article.ToData()

	data, err := s.repo.Update(ctx, data)
	if err != nil {
		return article, err
	}

	err = s.repo.UpdateCategories(ctx, id, cmd.Categories)
	if err != nil {
		return article, err
	}

	article = model.NewArticle(data)
	article.Categories = cmd.Categories

	return article, nil
}
