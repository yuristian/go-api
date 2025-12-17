package domain

type ProductRepository interface {
	FindByID(id uint) (*Product, error)
}
