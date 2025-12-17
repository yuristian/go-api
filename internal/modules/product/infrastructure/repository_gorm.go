package infrastructure

import (
	"github.com/yuristian/go-api/internal/modules/product/domain"
	"gorm.io/gorm"
)

type ProductGormRepository struct {
	db *gorm.DB
}

func NewProductGormRepository(db *gorm.DB) domain.ProductRepository {
	return &ProductGormRepository{db: db}
}

func (r *ProductGormRepository) FindByID(id uint) (*domain.Product, error) {
	var entity domain.Product
	if err := r.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}
