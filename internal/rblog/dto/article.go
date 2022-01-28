package dto

import (
	"context"
	metav1 "github.com/Pis0sion/rblogrus/store/rblog/v1"
)

type ArticleDto interface {
	Create(ctx context.Context, article *metav1.Article) error
	GetArticleList(ctx context.Context, page, pageSize int) (*metav1.ArticleList, error)
}

type ArticleEntity struct {
	Title    string `json:"title"`
	State    int    `json:"state"`
	CreateAt string `json:"createdAt,omitempty"`
}
