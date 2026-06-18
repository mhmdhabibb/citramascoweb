package types

import (
	"citramascoweb-backend/internal/dto"

	"github.com/gin-gonic/gin"
)

type typeHandler struct {
	typeService *typeService
}

func NewTypeHandler(typeService *typeService) *typeHandler {
	return &typeHandler{
		typeService: typeService,
	}
}

func (h *typeHandler) GetAll(c *gin.Context) {

	types, err := h.typeService.GetAll()

	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": "Error getting types data!"})
		return
	}

	c.JSON(200, gin.H{"success": true, "message": "Types data retrieved successfully!", "data": types})

}

func (h *typeHandler) GetById(c *gin.Context) {

	id := c.Param("id")

	types, err := h.typeService.GetById(id)
	if err != nil {
		c.JSON(404, gin.H{"success": false, "message": "Types not found!"})
		return
	}

	c.JSON(200, gin.H{"success": true, "message": "Detail of Types by Id ", "data": types})
}

func (h *typeHandler) Store(c *gin.Context) {

	var req dto.CreateTypesRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": "Bad Request"})
		return
	}

	err = h.typeService.Create(&req)

	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": "Error when creating types data!"})
		return
	}

	c.JSON(200, gin.H{"success": true, "message": "Types data created successfully!"})

}

func (h *typeHandler) Update(c *gin.Context) {

	id := c.Param("id")

	var req dto.UpdateTypesRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": "Bad Request"})
		return
	}

	err = h.typeService.Update(id, &req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": "Error when updating types data!"})
		return
	}

	c.JSON(200, gin.H{"success": true, "message": "Types data updated successfully!"})

}

func (h *typeHandler) Delete(c *gin.Context) {

	id := c.Param("id")

	err := h.typeService.Delete(id)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": "Error when deleting types data!"})
		return
	}

	c.JSON(200, gin.H{"success": true, "message": "Types data deleted successfully!"})

}

func (h *typeHandler) Search(c *gin.Context) {

	query := c.Query("query")

	types, err := h.typeService.Search(query)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": "Error when searching types data!"})
		return
	}

	c.JSON(200, gin.H{"success": true, "message": "Types data searched successfully!", "data": types})

}
