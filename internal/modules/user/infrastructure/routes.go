package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/yuristian/go-api/internal/modules/user/presentation"
	"github.com/yuristian/go-api/internal/modules/user/usecase"
)

func RegisterRoutes(rg *gin.RouterGroup, uc *usecase.Usecase) {
	handler := presentation.NewUserHandler(uc)

	users := rg.Group("/users")
	{
		users.POST("/register", handler.Register)
		users.POST("/login", handler.Login)
		users.GET("/:id", handler.GetByID)
	}
}
