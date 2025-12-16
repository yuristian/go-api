package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/yuristian/go-api/internal/modules/{{.ModuleName}}/presentation"
	"github.com/yuristian/go-api/internal/modules/{{.ModuleName}}/usecase"
)

func RegisterRoutes(rg *gin.RouterGroup, uc *usecase.{{.EntityName}}Usecase) {
	handler := presentation.New{{.EntityName}}Handler(uc)

	r := rg.Group("/{{.ModuleName}}s")
	r.GET("/:id", handler.GetByID)
}
