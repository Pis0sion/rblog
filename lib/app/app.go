package app

import (
	"github.com/Pis0sion/rblog/lib/errs"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	ctx *gin.Context
}

type Paginate struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalCount int64 `json:"totalCount"`
}

type unifiedRespBody struct {
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	RespData interface{} `json:"respData,omitempty"`
}

func NewResponse(ctx *gin.Context) *Response {

	return &Response{ctx: ctx}
}

func newUnifiedRespBody(errs *errs.Errors, respData interface{}) *unifiedRespBody {

	return &unifiedRespBody{
		Code:     errs.Code(),
		Message:  errs.Msg(),
		RespData: respData,
	}
}

func (r Response) ToResponse(respData interface{}) {

	r.ctx.JSON(http.StatusOK, newUnifiedRespBody(errs.Success, respData))
}

func (r Response) ToResponseList(respData interface{}, totalCount int64) {

	r.ctx.JSON(http.StatusOK, newUnifiedRespBody(errs.Success, gin.H{
		"list": respData,
		"page": Paginate{
			Page:       GetPage(r.ctx),
			PageSize:   GetPageSize(r.ctx),
			TotalCount: totalCount,
		},
	}))
}

func (r Response) ToErrorResponse(errors *errs.Errors) {

	r.ctx.JSON(errors.StatusCode(), newUnifiedRespBody(errors, nil))
}
