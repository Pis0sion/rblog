package dto

import "context"

type ArticleDto interface {
	GetArticleList(ctx context.Context) ([]*ArticleEntity, error)
}

type ArticleEntity struct {
	Title    string `json:"title"`
	State    int    `json:"state"`
	CreateAt string `json:"createdAt,omitempty"`
}