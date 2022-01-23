package article

import (
	v1 "github.com/Pis0sion/rblog/internal/rblog/srv/v1"
	"github.com/Pis0sion/rblog/pkg/app"
	"github.com/gin-gonic/gin"
)

// Article article struct
type Article struct{}

func NewArticle() *Article {
	return &Article{}
}

type ArtRequestEntity struct {
	Title string `json:"title" form:"title" binding:"required"`
}

func (a Article) Get(ctx *gin.Context) {}

func (a Article) List(ctx *gin.Context) {

	srvv1 := v1.NewSrv()
	articleList, totalCount := srvv1.Article().GetArticleList(ctx)
	app.NewResponse(ctx).ToResponseList(articleList, totalCount)

	return
}

func (a Article) Update(ctx *gin.Context) {}

func (a Article) Delete(ctx *gin.Context) {}
