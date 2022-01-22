package v1

import (
	v1 "github.com/Pis0sion/rblog/internal/rblog/srv/v1"
	"github.com/Pis0sion/rblog/pkg/app"
	"github.com/gin-gonic/gin"
)

// Article 文章结构体
type Article struct{}

func NewArticle() *Article {
	return &Article{}
}

func (a Article) Get(ctx *gin.Context) {}

func (a Article) List(ctx *gin.Context) {

	srvv1 := v1.NewSrv()
	articleList := srvv1.Article().GetArticleList(ctx)
	app.NewResponse(ctx).ToResponseList(articleList, len(articleList))

	return
}

func (a Article) Create(ctx *gin.Context) {}

func (a Article) Update(ctx *gin.Context) {}

func (a Article) Delete(ctx *gin.Context) {}
