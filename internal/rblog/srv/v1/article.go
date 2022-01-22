package v1

import (
	"context"
	"github.com/Pis0sion/rblog/internal/rblog/dto"
)

type ArticleSrv interface {
	Create() error
	GetArticleList(ctx context.Context) []*dto.ArticleEntity
}

type ArticleService struct {
	dto dto.Factory
}

func newArticles(dto dto.Factory) *ArticleService {
	return &ArticleService{dto: dto}
}

func (s *ArticleService) Create() error {
	return nil
}

func (s *ArticleService) GetArticleList(ctx context.Context) []*dto.ArticleEntity {

	if articles, err := s.dto.Articles().GetArticleList(ctx); err == nil {
		return articles
	}

	return nil
}
