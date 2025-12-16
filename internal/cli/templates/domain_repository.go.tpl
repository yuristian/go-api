package domain

type {{.EntityName}}Repository interface {
	FindByID(id uint) (*{{.EntityName}}, error)
}
