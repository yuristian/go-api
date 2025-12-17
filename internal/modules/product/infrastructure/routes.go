package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/yuristian/go-api/internal/modules/product/presentation"
	"github.com/yuristian/go-api/internal/modules/product/usecase"
)

func RegisterRoutes(rg *gin.RouterGroup, uc *usecase.ProductUsecase) {
	handler := presentation.NewProductHandler(uc)

	r := rg.Group("/products")
	r.GET("/:id", handler.GetByID)
}
