package modules

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yuristian/go-api/internal/auth"

	productInfra "github.com/yuristian/go-api/internal/modules/product/infrastructure"
	productUsecase "github.com/yuristian/go-api/internal/modules/product/usecase"
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
	productRepo := productInfra.NewProductGormRepository(gormDB)
	productUC := productUsecase.NewProductUsecase(productRepo)
	productInfra.RegisterRoutes(rg, productUC)

}
