package inventory

import (
	"time"

	"citramascoweb-backend/internal/dto"

	"github.com/gin-gonic/gin"
)

type inventoryHandler struct {
	service *inventoryService
}

func NewInventoryHandler(service *inventoryService) *inventoryHandler {
	return &inventoryHandler{service: service}
}

func (h *inventoryHandler) ListItems(c *gin.Context) {
	items, err := h.service.ListItems()
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": "Error getting inventory items!", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "List of inventory items", "data": items})
}

func (h *inventoryHandler) CreateItem(c *gin.Context) {
	var req dto.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	item, err := h.service.CreateItem(&req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Item created successfully!", "data": item})
}

func (h *inventoryHandler) UpdateItem(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := h.service.UpdateItem(&req, id); err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Item updated successfully!"})
}

func (h *inventoryHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteItem(id); err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Item deleted successfully!"})
}

func (h *inventoryHandler) StockIn(c *gin.Context) {
	id := c.Param("id")
	actor := actorFrom(c)
	var req dto.StockInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := h.service.StockIn(id, &req, actor); err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Stock added successfully!"})
}

func (h *inventoryHandler) Usage(c *gin.Context) {
	id := c.Param("id")
	actor := actorFrom(c)
	var req dto.UsageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := h.service.Usage(id, &req, actor); err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Usage recorded successfully!"})
}

func (h *inventoryHandler) StockTake(c *gin.Context) {
	actor := actorFrom(c)
	var req dto.StockTakeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := h.service.StockTake(&req, actor); err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Stock take recorded and reconciled!"})
}

func (h *inventoryHandler) Transactions(c *gin.Context) {
	from := parseDate(c.Query("from"))
	to := parseDateEndOfDay(c.Query("to"))
	txs, err := h.service.ListTransactions(from, to)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Inventory transactions", "data": txs})
}

func (h *inventoryHandler) Report(c *gin.Context) {
	from := parseDate(c.Query("from"))
	to := parseDateEndOfDay(c.Query("to"))
	report, err := h.service.Report(from, to)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Inventory report", "data": report})
}

func actorFrom(c *gin.Context) string {
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func parseDate(raw string) time.Time {
	if raw == "" {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Time{}
	}
	return t
}

func parseDateEndOfDay(raw string) time.Time {
	if raw == "" {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Time{}
	}
	return t.Add(24*time.Hour - time.Second)
}