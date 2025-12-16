package presentation

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yuristian/go-api/internal/modules/{{.ModuleName}}/usecase"
)

type {{.EntityName}}Handler struct {
	uc *usecase.{{.EntityName}}Usecase
}

func New{{.EntityName}}Handler(uc *usecase.{{.EntityName}}Usecase) *{{.EntityName}}Handler {
	return &{{.EntityName}}Handler{uc: uc}
}

func (h *{{.EntityName}}Handler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.uc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
