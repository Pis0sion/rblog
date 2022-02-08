package dto

import (
	"context"
	metav1 "github.com/Pis0sion/rblogrus/store/rblog/v1"
)

type ArticleDto interface {
	GetArticle(ctx context.Context, articleID int) (*metav1.Article, error)
	CreateArticle(ctx context.Context, article *metav1.Article) error
	GetArticleList(ctx context.Context, page, pageSize int) (*metav1.ArticleList, error)
}
