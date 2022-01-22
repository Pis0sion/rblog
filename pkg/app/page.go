package app

import (
	"github.com/Pis0sion/rblog/internal/rblog/cfg"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetPage(ctx *gin.Context) int {

	i, _ := strconv.Atoi(ctx.Query("page"))

	if i <= 0 {
		return 1
	}

	return i
}

func GetPageSize(ctx *gin.Context) int {

	i, _ := strconv.Atoi(ctx.Query("pageSize"))

	if i <= 0 {
		return cfg.ApplicationParameters.DefaultPageSize
	}

	if i > cfg.ApplicationParameters.MaxPageSize {
		return cfg.ApplicationParameters.MaxPageSize
	}

	return i
}
