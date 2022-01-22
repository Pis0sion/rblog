package v1

import (
	"github.com/Pis0sion/rblog/internal/rblog/dto"
)

type Service interface {
	Article() ArticleSrv
}

type srv struct {
	dto dto.Factory
}

func NewSrv() *srv {
	return &srv{dto: dto.Client()}
}

func (s *srv) Article() ArticleSrv {
	return newArticles(s.dto)
}
