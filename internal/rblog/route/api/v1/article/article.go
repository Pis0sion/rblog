package article

import (
	srvv1 "github.com/Pis0sion/rblog/internal/rblog/srv/v1"
	"github.com/Pis0sion/rblog/pkg/app"
	"github.com/Pis0sion/rblog/pkg/errs"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Article article struct
type Article struct{}

func NewArticle() *Article {
	return &Article{}
}

func (a Article) Get(ctx *gin.Context) {

	artID, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		app.NewResponse(ctx).ToErrorResponse(errs.ServerErrors)
		return
	}

	if article, err := srvv1.NewSrv().Article().Get(ctx, artID); err == nil {
		a := struct {
			ID                 uint64 `json:"id"`
			AuthorID           uint64 `json:"authorID"`
			ArticleTitle       string `json:"articleTitle"`
			ArticleContent     string `json:"articleContent"`
			ArticleDescription string `json:"articleDescription"`
			CreatedAt          string `json:"createdAt"`
			UpdatedAt          string `json:"updatedAt"`
		}{
			article.ID,
			article.AuthorID,
			article.ArticleTitle,
			article.ArticleContent,
			article.ArticleDescription,
			article.CreatedAt.Format("2006-01-02 15:04:05"),
			article.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		app.NewResponse(ctx).ToResponse(&a)
		return
	} else {
		app.NewResponse(ctx).ToErrorResponse(errs.NotFound)
		return
	}
}

func (a Article) Update(ctx *gin.Context) {}

func (a Article) Delete(ctx *gin.Context) {}
