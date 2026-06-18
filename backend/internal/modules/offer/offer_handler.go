package offer

import (
	"citramascoweb-backend/internal/dto"

	"github.com/gin-gonic/gin"
)

type offerHandler struct {
	offerService *offerService
}

func NewOfferHandler(offerService *offerService) *offerHandler {
	return &offerHandler{offerService: offerService}
}

func (h *offerHandler) GetAll(c *gin.Context) {
	offers, err := h.offerService.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"message": "Error getting offers data", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "List of Offers data!", "data": offers})
}
func (h *offerHandler) GetAlls(c *gin.Context) {
	offers, err := h.offerService.GetAlls()
	if err != nil {
		c.JSON(500, gin.H{"message": "Error getting offers data", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "List of Offers data!", "data": offers})
}

func (h *offerHandler) GetById(c *gin.Context) {
	id := c.Param("id")

	offer, err := h.offerService.GetById(id)
	if err != nil {
		c.JSON(500, gin.H{"message": "Error getting offer data", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Detail of offer data", "data": offer})
}

func (h *offerHandler) Create(c *gin.Context) {
	var req dto.CreateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Bad Request", "error": err.Error()})
		return
	}

	err := h.offerService.Create(&req)
	if err != nil {
		c.JSON(500, gin.H{"message": "Error creating offer data", "error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Offer data created successfully!"})
}

func (h *offerHandler) Update(c *gin.Context) {
	var req dto.UpdateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Error binding JSON", "error": err.Error()})
		return
	}

	err := h.offerService.Update(c.Param("id"), &req)
	if err != nil {
		c.JSON(500, gin.H{"message": "Error updating offer data", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Offer data updated successfully!"})
}

func (h *offerHandler) Delete(c *gin.Context) {
	err := h.offerService.Delete(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"message": "Error deleting offer data", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Offer data deleted successfully!"})
}
