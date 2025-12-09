package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuristian/go-api/internal/infrastructure/auth"
	"github.com/yuristian/go-api/internal/infrastructure/config"
	"github.com/yuristian/go-api/internal/infrastructure/db"
	"github.com/yuristian/go-api/internal/infrastructure/middleware"
	"github.com/yuristian/go-api/internal/modules"
	// userDomain "github.com/yuristian/go-api/internal/modules/user/domain"
	// userInfra "github.com/yuristian/go-api/internal/modules/user/infrastructure"
	// userUsecase "github.com/yuristian/go-api/internal/modules/user/usecase"
	// // Module registry
)

func main() {
	cfg := config.LoadConfig("configs/config.yaml")

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gormDB := db.NewGormDB(cfg)
	// gormDB.AutoMigrate(&userDomain.User{}) // migrate User table

	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpiresIn)

	// userRepo := userInfra.NewUserGormRepository(gormDB)
	// userUC := userUsecase.NewUserUsecase(userRepo, jwtManager)

	api := r.Group("/api")
	modules.RegisterAllModules(api, gormDB, jwtManager)

	protected := api.Group("/protected")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	protected.GET("/profile", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"user_id": c.GetUint("user_id"),
			"email":   c.GetString("email"),
			"role":    c.GetString("role"),
		})
	})

	// Example admin-only route
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(jwtManager), middleware.RequireRole("admin"))
	admin.GET("/stats", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Admin access granted"})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
