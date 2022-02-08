package v1

import (
	"context"
	"github.com/Pis0sion/rblog/internal/rblog/dto"
	metav1 "github.com/Pis0sion/rblogrus/store/rblog/v1"
)

type ArticleSrv interface {
	Get(ctx context.Context, articleID int) (*metav1.Article, error)
	Create(ctx context.Context, article *metav1.Article) error
	List(ctx context.Context, page, pageSize int) (*metav1.ArticleList, error)
}

type ArticleService struct {
	dto dto.Factory
}

func newArticles(dto dto.Factory) ArticleSrv {
	return &ArticleService{dto: dto}
}

func (s *ArticleService) Get(ctx context.Context, articleID int) (*metav1.Article, error) {
	return s.dto.Articles().GetArticle(ctx, articleID)
}

func (s *ArticleService) Create(ctx context.Context, article *metav1.Article) error {
	return s.dto.Articles().CreateArticle(ctx, article)
}

func (s *ArticleService) List(ctx context.Context, page, pageSize int) (*metav1.ArticleList, error) {
	return s.dto.Articles().GetArticleList(ctx, page, pageSize)
}
