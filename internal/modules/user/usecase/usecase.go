package usecase

import (
	"strconv"

	"github.com/yuristian/go-api/internal/auth"
	user "github.com/yuristian/go-api/internal/modules/user/domain"
	"golang.org/x/crypto/bcrypt"
)

type Usecase struct {
	repo       user.Repository
	jwtManager *auth.JWTManager
}

func NewUserUsecase(r user.Repository, jwtManager *auth.JWTManager) *Usecase {
	return &Usecase{repo: r, jwtManager: jwtManager}
}

func (u *Usecase) Register(name, email, password, role string) (*user.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	newUser := &user.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
		Role:     role,
	}
	err = u.repo.Create(newUser)
	return newUser, err
}

func (u *Usecase) Login(email, password string) (string, *user.User, error) {
	existingUser, err := u.repo.GetByEmail(email)
	if err != nil {
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password)); err != nil {
		return "", nil, err
	}

	token, err := u.jwtManager.GenerateToken(existingUser.ID, existingUser.Email, existingUser.Role)
	if err != nil {
		return "", nil, err
	}

	return token, existingUser, nil
}

func (u *Usecase) GetByID(idStr string) (*user.User, error) {
	// convert string to uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return nil, err
	}

	return u.repo.GetByID(uint(id))
}
