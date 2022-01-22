package mysql

import (
	"context"
	"github.com/Pis0sion/rblog/internal/rblog/dto"
	"github.com/Pis0sion/rblog/internal/rblog/mod"
	"gorm.io/gorm"
)

type articles struct {
	dbIns *gorm.DB
}

func newArticles(db *gorm.DB) dto.ArticleDto {
	return &articles{dbIns: db}
}

func (a *articles) GetArticleList(ctx context.Context) (list []*dto.ArticleEntity, err error) {

	article := mod.Article{
		Title: "",
		State: 0,
	}

	articleList, _ := article.List(a.dbIns, 0, 10)

	for _, item := range articleList {
		list = append(
			list, &dto.ArticleEntity{
				Title:    item.Title,
				State:    item.State,
				CreateAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
			},
		)
	}

	return list, nil
}
