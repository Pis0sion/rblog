package mysql

import (
	"github.com/Pis0sion/rblog/internal/rblog/dto"
	"gorm.io/gorm"
)

type articleTags struct {
	dbIns *gorm.DB
}

func newArticleTags(db *gorm.DB) dto.ArticleTagsDto {
	return &articleTags{dbIns: db}
}
