package errs

var (
	Success                   = NewErrors(0, "成功")
	ServerErrors              = NewErrors(100000, "服务器内部错误")
	InvalidParams             = NewErrors(100001, "入参错误")
	NotFound                  = NewErrors(100002, "找不到")
	UnauthorizedAuthNotExist  = NewErrors(100003, "鉴权失败")
	UnauthorizedTokenErrors   = NewErrors(100004, "鉴权失败")
	UnauthorizedTokenTimeout  = NewErrors(100005, "鉴权失败")
	UnauthorizedTokenGenerate = NewErrors(100006, "鉴权失败")
	TooManyRequests           = NewErrors(100007, "请求过多")

	ArticleNotFound = NewErrors(100008, "article 不存在...")
)
