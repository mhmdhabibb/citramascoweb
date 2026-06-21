package rooms

import (
	"citramascoweb-backend/internal/dto"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type roomHandler struct {
	roomService *roomService
	supabaseURL string
	anonKey     string
}

func NewRoomHandler(roomService *roomService) *roomHandler {
	return &roomHandler{
		roomService: roomService,
		supabaseURL: os.Getenv("SUPABASE_URL"),
		anonKey:     os.Getenv("SUPABASE_ANON_KEY"),
	}
}

func (h *roomHandler) uploadToSupabase(c *gin.Context, fileKey string) (string, error) {
	file, header, err := c.Request.FormFile(fileKey)
	if err != nil {
		return "", err // Mengembalikan error jika file tidak ditemukan/tidak di-upload
	}
	defer file.Close()

	// Membuat nama unik file di storage agar tidak bentrok
	uniqueFileName := fmt.Sprintf("room_%d_%s", time.Now().UnixNano(), header.Filename)
	uploadURL := fmt.Sprintf("%s/storage/v1/object/room-images/%s", h.supabaseURL, uniqueFileName)

	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, uploadURL, file)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+h.anonKey)
	req.Header.Set("apiKey", h.anonKey)
	req.Header.Set("Content-Type", header.Header.Get("Content-Type"))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("supabase storage error dengan status: %d", resp.StatusCode)
	}

	// Mengembalikan URL publik gambar
	return fmt.Sprintf("%s/storage/v1/object/public/room-images/%s", h.supabaseURL, uniqueFileName), nil
}

func (h *roomHandler) GetAll(c *gin.Context) {
	rooms, err := h.roomService.GetAll()
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Rooms retrieved successfully", "data": rooms})
}

func (h *roomHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	room, err := h.roomService.GetById(id)
	if err != nil {
		c.JSON(404, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Room retrieved successfully", "data": room})
}

func (h *roomHandler) Store(c *gin.Context) {
	var req dto.CreateRoomRequest

	// Gunakan ShouldBind untuk membaca data Multipart Form-Data
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Proses upload file gambar ke Supabase
	imageURL, err := h.uploadToSupabase(c, "image")
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": "Failed to upload image: " + err.Error()})
		return
	}
	req.Image = imageURL // Suntikkan URL gambar ke DTO

	err = h.roomService.Store(&req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Room created successfully!"})
}

func (h *roomHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateRoomRequest
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Proses upload file gambar ke Supabase jika ada file baru yang diupload
	imageURL, err := h.uploadToSupabase(c, "image")
	if err == nil {
		req.Image = imageURL
	} else if err != http.ErrMissingFile {
		c.JSON(400, gin.H{"success": false, "message": "Failed to upload image: " + err.Error()})
		return
	}

	err = h.roomService.Update(id, &req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Room updated successfully!"})
}

func (h *roomHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.roomService.Delete(id)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Room deleted successfully!"})
}

func (h *roomHandler) FilerByCategory(c *gin.Context) {
	id := c.Param("id")
	rooms, err := h.roomService.FilerByCategory(id)
	if err != nil {
		c.JSON(404, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Room filtered by category successfully", "data": rooms})
}

func (h *roomHandler) FilterByType(c *gin.Context) {
	id := c.Param("id")
	rooms, err := h.roomService.FilterByType(id)
	if err != nil {
		c.JSON(404, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Room filtered by type successfully", "data": rooms})
}

func (h *roomHandler) Filter(c *gin.Context) {
	status := c.Query("status")
	availabilityStatus := c.Query("availability_status")
	checkinStr := c.Query("check_in_date")
	checkoutStr := c.Query("check_out_date")

	var checkinDate, checkoutDate time.Time
	var err error

	if checkinStr != "" {
		checkinDate, err = time.Parse("02-01-2006", checkinStr)
		if err != nil {
			checkinDate, err = time.Parse(time.RFC3339, checkinStr)
		}
	}
	if checkoutStr != "" {
		checkoutDate, err = time.Parse("02-01-2006", checkoutStr)
		if err != nil {
			checkoutDate, err = time.Parse(time.RFC3339, checkoutStr)
		}
	}

	rooms, err := h.roomService.Filter(status, availabilityStatus, checkinDate, checkoutDate)

	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": "Failed to Filter Rooms", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": true, "message": "Rooms filtered successfully", "data": rooms})
}

// Buat struct DTO lokal ringkas untuk menerima payload status
type UpdateRoomStatusPayload struct {
	Status string `json:"status" binding:"required"`
}

func (h *roomHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	log.Printf("[DEBUG][ROOM-HANDLER] Incoming request HTTP PATCH /room/status/%s", id)

	var payload UpdateRoomStatusPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("[WARN][ROOM-HANDLER] Validation failed for incoming JSON payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload. The 'status' field is required.",
		})
		return
	}

	log.Printf("[DEBUG][ROOM-HANDLER] Payload extracted successfully. Target status: '%s'", payload.Status)

	// Execute service mutation
	if err := h.roomService.UpdateStatus(id, payload.Status); err != nil {
		log.Printf("[ERROR][ROOM-HANDLER] Service failed to update room status: %v", err)

		// Professional Error Handling Replacement: Avoiding misleading 404s
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Room operational status has been updated successfully to '%s'.", payload.Status),
	})
}
