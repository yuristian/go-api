package infrastructure

import (
	user "github.com/yuristian/go-api/internal/modules/user/domain"
	"gorm.io/gorm"
)

type UserGormRepository struct {
	db *gorm.DB
}

func NewUserGormRepository(db *gorm.DB) user.Repository {
	return &UserGormRepository{db: db}
}

func (r *UserGormRepository) Create(u *user.User) error {
	return r.db.Create(u).Error
}

func (r *UserGormRepository) GetByID(id uint) (*user.User, error) {
	var u user.User
	err := r.db.First(&u, id).Error
	return &u, err
}

func (r *UserGormRepository) GetByEmail(email string) (*user.User, error) {
	var u user.User
	err := r.db.Where("email = ?", email).First(&u).Error
	return &u, err
}

func (r *UserGormRepository) Update(u *user.User) error {
	return r.db.Save(u).Error
}

func (r *UserGormRepository) Delete(id uint) error {
	return r.db.Delete(&user.User{}, id).Error
}
