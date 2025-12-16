package usecase

import "github.com/yuristian/go-api/internal/modules/{{.ModuleName}}/domain"

type {{.EntityName}}Usecase struct {
	repo domain.{{.EntityName}}Repository
}

func New{{.EntityName}}Usecase(repo domain.{{.EntityName}}Repository) *{{.EntityName}}Usecase {
	return &{{.EntityName}}Usecase{repo: repo}
}

func (u *{{.EntityName}}Usecase) GetByID(id uint) (*domain.{{.EntityName}}, error) {
	return u.repo.FindByID(id)
}
