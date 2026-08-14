package channel

import (
	"citramascoweb-backend/internal/dto"
	"strconv"

	"github.com/gin-gonic/gin"
)

type channelHandler struct {
	channelService *channelService
}

func NewChannelHandler(channelService *channelService) *channelHandler {
	return &channelHandler{channelService: channelService}
}

func (h *channelHandler) GetAll(c *gin.Context) {
	channels, err := h.channelService.FindAll()
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "List of Channel",
		"data":    channels,
	})
}

func (h *channelHandler) Get(c *gin.Context) {
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)

	if page < 1 {
		page = 1
	}

	channels, total, err := h.channelService.Find(page)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "List of Channel",
		"data":    channels,
		"meta": gin.H{
			"page":  page,
			"total": total,
			"limit": 10,
		},
	})
}

func (h *channelHandler) Store(c *gin.Context) {
	var req *dto.CreateChannelRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	channel, err := h.channelService.FindByName(req.Name)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if channel != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Channel name already exists",
		})
		return
	}

	if err := h.channelService.Create(req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Channel created successfully!",
	})
}

func (h *channelHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req *dto.UpdateChannelRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := h.channelService.Update(id, req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Channel updated successfully!",
	})
}

func (h *channelHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.channelService.Delete(id); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Channel deleted successfully!",
	})
}

func (h *channelHandler) GetLastChannel(c *gin.Context) {
	channel, err := h.channelService.GetLastChannel()
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Last Channel",
		"data":    channel,
	})
}
func (h *channelHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	channel, err := h.channelService.FindById(id)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Detail Channel",
		"data":    channel,
	})
}
func (h *channelHandler) GetByName(c *gin.Context) {
	id := c.Query("name")
	channel, err := h.channelService.FindByName(id)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Channel by Name",
		"data":    channel,
	})
}
