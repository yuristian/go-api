package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/yuristian/go-api/internal/auth"
	"github.com/yuristian/go-api/internal/config"
	"github.com/yuristian/go-api/internal/infrastructure/db"
	"github.com/yuristian/go-api/internal/middleware"
	"github.com/yuristian/go-api/internal/modules"
)

func main() {
	// 1. Load configuration
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatal("Failed to load configuration")
	}

	// 2. Setup Gin engine
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 3. Initialize database (using the correct function)
	gormDB := db.NewGormDB(cfg)

	// 4. Setup JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpiresIn)

	// 5. Register all modules
	api := r.Group("/api")
	modules.RegisterAllModules(api, gormDB, jwtManager)

	// 6. Protected route group
	protected := api.Group("/protected")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	protected.GET("/profile", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"user_id": c.GetUint("user_id"),
			"email":   c.GetString("email"),
			"role":    c.GetString("role"),
		})
	})

	// 7. Admin route group
	admin := api.Group("/admin")
	admin.Use(
		middleware.AuthMiddleware(jwtManager),
		middleware.RequireRole("admin"),
	)
	admin.GET("/stats", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Admin access granted"})
	})

	// 8. Healthcheck
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 9. Start server
	log.Printf("Server running on port %d", cfg.Server.Port)
	r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
