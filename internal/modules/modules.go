package modules

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yuristian/go-api/internal/auth"

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
}
