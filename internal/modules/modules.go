package modules

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yuristian/go-api/internal/infrastructure/auth"

	userInfra "github.com/yuristian/go-api/internal/modules/user/infrastructure"
	userUsecase "github.com/yuristian/go-api/internal/modules/user/usecase"
)

func RegisterAllModules(rg *gin.RouterGroup, gormDB *gorm.DB, jwtManager *auth.JWTManager) {
	// ----------------------
	// User Module
	// ----------------------
	userRepo := userInfra.NewUserGormRepository(gormDB)
	userUC := userUsecase.NewUserUsecase(userRepo, jwtManager)
	userInfra.RegisterRoutes(rg, userUC)

	// ----------------------
	// Module baru nanti ditambahkan di sini:
	// productRepo := productInfra.NewProductGormRepository(gormDB)
	// productUC := productUsecase.NewProductUsecase(productRepo)
	// productInfra.RegisterRoutes(rg, productUC)
}
