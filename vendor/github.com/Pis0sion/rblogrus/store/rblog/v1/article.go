package v1

import (
	"github.com/Pis0sion/components/validate"
	metav1 "github.com/Pis0sion/rblogrus/meta/types"
)

type unifiedValidator interface {
	Validate() error
}

type Article struct {
	metav1.ObjectMeta  `json:",inline"`
	AuthorID           uint64 `json:"authorID"     gorm:"column:authorID" validate:"omitempty"`
	ArticleTitle       string `json:"articleTitle" gorm:"column:articleTitle" validate:"omitempty"`
	ArticleContent     string `json:"articleContent"  gorm:"column:articleContent" validate:"omitempty"`
	ArticleDescription string `json:"articleDescription"  gorm:"column:articleDescription" validate:"omitempty"`
	Status             int    `json:"status"    gorm:"column:status" validate:"omitempty"`
}

type ArticleList struct {
	metav1.ListMeta `json:",inline"`
	Items           []*Article `json:"items"`
}

func (a *Article) TableName() string {
	return "article"
}

type Art2Create struct {
	AuthorID           uint64 `json:"authorID"  validate:"required"`
	ArticleTitle       string `json:"articleTitle"  validate:"required,min=1,max=100"`
	ArticleContent     string `json:"articleContent"  validate:"required,min=1"`
	ArticleDescription string `json:"articleDescription" validate:"omitempty"`
}

func (c *Art2Create) Validate() error {
	return validate.NewValidator(c).Validate()
}
