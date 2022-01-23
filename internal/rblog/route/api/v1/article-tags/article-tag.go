package article_tags

import "github.com/gin-gonic/gin"

// ArticleTag 文章标签
type ArticleTag struct {}

func NewArticleTag() *ArticleTag {
	return &ArticleTag{}
}

func (a ArticleTag) Get(ctx *gin.Context) {}

func (a ArticleTag) List(ctx *gin.Context) {}

func (a ArticleTag) Create(ctx *gin.Context) {}

func (a ArticleTag) Update(ctx *gin.Context) {}

func (a ArticleTag) Delete(ctx *gin.Context) {}