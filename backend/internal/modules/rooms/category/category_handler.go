package category

import (
	"citramascoweb-backend/internal/dto"

	"github.com/gin-gonic/gin"
)

type categoryHandler struct {
	categoryService *categoryService
}

func NewCategoryHandler(categoryService *categoryService) *categoryHandler {
	return &categoryHandler{categoryService: categoryService}
}

func (h *categoryHandler) GetById(c *gin.Context) {

	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{
			"error": "ID is required",
		})
		return
	}

	category, err := h.categoryService.FindById(id)
	if err != nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "Category Not Found!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"success": true, "message": "Detail of Category By Id", "data": category})

}

func (h *categoryHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	if slug == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "slug is required!",
		})
		return
	}

	category, err := h.categoryService.FindBySlug(slug)
	if err != nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "Category Not Found!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"success": true, "message": "Detail of Category By Slug", "data": category})
}

func (h *categoryHandler) GetAll(c *gin.Context) {
	categories, err := h.categoryService.FindAll()

	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Error when getting category data!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "List all category!",
		"data":    categories,
	})
}

func (h *categoryHandler) Create(c *gin.Context) {

	var req dto.CreateCategoryRequest

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Validation failed!",
			"error":   err.Error(),
		})
		return
	}

	err = h.categoryService.Create(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Error when creating category!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{"success": true, "message": "Category created successfully!"})

}

func (h *categoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "ID is required!",
		})
		return
	}

	var req dto.UpdateCategoryRequest

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Validation failed!",
			"error":   err.Error(),
		})
		return
	}

	err = h.categoryService.Update(id, &req)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Error when updating category!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"success": true, "message": "Category updated successfully!"})
}

func (h *categoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "ID is required!",
		})
		return
	}

	err := h.categoryService.Delete(id)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Error when deleting category!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"success": true, "message": "Category deleted successfully!"})

}
