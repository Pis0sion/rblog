package v1

import (
	"context"
	"github.com/Pis0sion/rblog/internal/rblog/dto"
)

type ArticleSrv interface {
	Create(ctx context.Context, title string) error
	GetArticleList(ctx context.Context) ([]*dto.ArticleEntity, int64)
}

type ArticleService struct {
	dto dto.Factory
}

func newArticles(dto dto.Factory) *ArticleService {
	return &ArticleService{dto: dto}
}

func (s *ArticleService) Create(ctx context.Context, title string) error {

	article := &dto.ArticleEntity{
		Title: title,
	}
	return s.dto.Articles().Create(ctx, article)
}

func (s *ArticleService) GetArticleList(ctx context.Context) ([]*dto.ArticleEntity, int64) {

	if articles, count, err := s.dto.Articles().GetArticleList(ctx); err == nil {
		return articles, count
	}

	return nil, 0
}
