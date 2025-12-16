package infrastructure

import (
	"github.com/yuristian/go-api/internal/modules/{{.ModuleName}}/domain"
	"gorm.io/gorm"
)

type {{.EntityName}}GormRepository struct {
	db *gorm.DB
}

func New{{.EntityName}}GormRepository(db *gorm.DB) domain.{{.EntityName}}Repository {
	return &{{.EntityName}}GormRepository{db: db}
}

func (r *{{.EntityName}}GormRepository) FindByID(id uint) (*domain.{{.EntityName}}, error) {
	var entity domain.{{.EntityName}}
	if err := r.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}
