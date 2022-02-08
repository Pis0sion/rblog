package article

import (
	srvv1 "github.com/Pis0sion/rblog/internal/rblog/srv/v1"
	"github.com/Pis0sion/rblog/pkg/app"
	"github.com/Pis0sion/rblog/pkg/errs"
	"github.com/gin-gonic/gin"
)

func (a Article) List(ctx *gin.Context) {

	articleList, err := srvv1.NewSrv().Article().List(ctx, app.GetPage(ctx), app.GetPageSize(ctx))

	if err != nil {
		app.NewResponse(ctx).ToErrorResponse(errs.NotFound)
		return
	}

	app.NewResponse(ctx).ToResponseList(articleList.Items, articleList.TotalCount)
	return
}
