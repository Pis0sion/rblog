package dto

import "context"

type ArticleDto interface {
	Create(ctx context.Context, articleEntity *ArticleEntity) error
	GetArticleList(ctx context.Context) ([]*ArticleEntity, int64, error)
}

type ArticleEntity struct {
	Title    string `json:"title"`
	State    int    `json:"state"`
	CreateAt string `json:"createdAt,omitempty"`
}
