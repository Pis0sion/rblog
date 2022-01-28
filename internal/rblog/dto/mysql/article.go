package mysql

import (
	"context"
	"github.com/Pis0sion/rblog/internal/rblog/dto"
	metav1 "github.com/Pis0sion/rblogrus/store/rblog/v1"
	"gorm.io/gorm"
)

type articles struct {
	dbIns *gorm.DB
}

func newArticles(db *gorm.DB) dto.ArticleDto {
	return &articles{dbIns: db}
}

func (a *articles) Create(ctx context.Context, article *metav1.Article) error {
	return a.dbIns.Create(article).Error
}

func (a *articles) GetArticleList(ctx context.Context, page, pageSize int) (*metav1.ArticleList, error) {
	articleList := metav1.ArticleList{}

	if err := a.dbIns.Offset((page - 1) * pageSize).Limit(pageSize).Find(&articleList.Items).Offset(-1).Limit(-1).Count(&articleList.TotalCount).Error; err != nil {
		return nil, err
	}

	return &articleList, nil
}
