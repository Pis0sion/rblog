package v1

import (
	"context"
	"github.com/Pis0sion/rblog/internal/rblog/dto"
	metav1 "github.com/Pis0sion/rblogrus/store/rblog/v1"
)

type ArticleSrv interface {
	Create(ctx context.Context, article *metav1.Article) error
	GetArticleList(ctx context.Context, page, pageSize int) (*metav1.ArticleList, error)
}

type ArticleService struct {
	dto dto.Factory
}

func newArticles(dto dto.Factory) ArticleSrv {
	return &ArticleService{dto: dto}
}

func (s *ArticleService) Create(ctx context.Context, article *metav1.Article) error {
	return s.dto.Articles().Create(ctx, article)
}

func (s *ArticleService) GetArticleList(ctx context.Context, page, pageSize int) (*metav1.ArticleList, error) {
	return s.dto.Articles().GetArticleList(ctx, page, pageSize)
}
