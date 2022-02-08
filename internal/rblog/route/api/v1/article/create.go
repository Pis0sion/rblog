package article

import (
	srvv1 "github.com/Pis0sion/rblog/internal/rblog/srv/v1"
	"github.com/Pis0sion/rblog/pkg/app"
	"github.com/Pis0sion/rblog/pkg/errs"
	metav1 "github.com/Pis0sion/rblogrus/store/rblog/v1"
	"github.com/gin-gonic/gin"
)

func (a Article) Create(ctx *gin.Context) {

	art2Req := metav1.Art2Create{}
	if err := ctx.ShouldBindJSON(&art2Req); err != nil {
		app.NewResponse(ctx).ToErrorResponse(errs.InvalidParams)
		return
	}

	if err := art2Req.Validate(); err != nil {
		app.NewResponse(ctx).ToErrorResponse(errs.NewErrors(10050, err.Error()))
		return
	}

	if err := srvv1.NewSrv().Article().Create(ctx, &metav1.Article{
		AuthorID:           art2Req.AuthorID,
		ArticleTitle:       art2Req.ArticleTitle,
		ArticleContent:     art2Req.ArticleContent,
		ArticleDescription: art2Req.ArticleDescription,
	}); err != nil {
		app.NewResponse(ctx).ToErrorResponse(errs.ServerErrors)
		return
	}

	app.NewResponse(ctx).ToResponse(nil)
	return
}
