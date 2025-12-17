package usecase

import "github.com/yuristian/go-api/internal/modules/product/domain"

type ProductUsecase struct {
	repo domain.ProductRepository
}

func NewProductUsecase(repo domain.ProductRepository) *ProductUsecase {
	return &ProductUsecase{repo: repo}
}

func (u *ProductUsecase) GetByID(id uint) (*domain.Product, error) {
	return u.repo.FindByID(id)
}
