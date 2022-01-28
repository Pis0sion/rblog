package article

import (
	"github.com/gin-gonic/gin"
)

// Article article struct
type Article struct{}

func NewArticle() *Article {
	return &Article{}
}

func (a Article) Get(ctx *gin.Context) {}

func (a Article) Update(ctx *gin.Context) {}

func (a Article) Delete(ctx *gin.Context) {}
