package route

import (
	article2 "github.com/Pis0sion/rblog/internal/rblog/route/api/v1/article"
	v1 "github.com/Pis0sion/rblog/internal/rblog/route/api/v1/article-tags"
	"github.com/gin-gonic/gin"
)

// InitializeRouters
// load routers
func InitializeRouters(engine *gin.Engine) {

	article := article2.NewArticle()
	articleTag := v1.NewArticleTag()

	apiV1 := engine.Group("/api/v1")

	{
		apiV1.POST("/articles", article.Create)
		apiV1.DELETE("/article/:id", article.Delete)
		apiV1.PUT("/article/:id", article.Update)
		apiV1.GET("/articles", article.List)
		apiV1.GET("/article/:id", article.Get)

		apiV1.POST("/article-tags", articleTag.Create)
		apiV1.DELETE("/article-tag/:id", articleTag.Delete)
		apiV1.PUT("/article-tag/:id", articleTag.Update)
		apiV1.GET("/article-tags", articleTag.List)
		apiV1.GET("/article-tag/:id", articleTag.Get)
	}

	return
}
