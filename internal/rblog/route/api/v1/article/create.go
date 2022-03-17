package article

import (
	v1 "github.com/Pis0sion/rblog/internal/rblog/srv/v1"
	"github.com/Pis0sion/rblog/lib/app"
	"github.com/Pis0sion/rblog/lib/errs"
	"github.com/gin-gonic/gin"
)

func (a Article) Create(ctx *gin.Context) {
	art2Req := ArtRequestEntity{}

	if err := ctx.ShouldBindJSON(&art2Req); err != nil {
		app.NewResponse(ctx).ToErrorResponse(errs.InvalidParams)
		return
	}

	srvv1 := v1.NewSrv()
	if err := srvv1.Article().Create(ctx, art2Req.Title); err != nil {
		app.NewResponse(ctx).ToErrorResponse(errs.ServerErrors)
		return
	}

	app.NewResponse(ctx).ToResponse(nil)
	return
}
