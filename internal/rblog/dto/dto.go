package dto

var agent Factory

type Factory interface {
	Articles() ArticleDto
	ArticleTags() ArticleTagsDto
}

func Client() Factory {
	return agent
}

func SetClient(factory Factory) {
	agent = factory
}
