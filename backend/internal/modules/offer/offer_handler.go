package offer

import (
	"citramascoweb-backend/internal/dto"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type offerHandler struct {
	offerService *offerService
	supabaseURL  string
	anonKey      string
}

func NewOfferHandler(offerService *offerService) *offerHandler {
	return &offerHandler{
		offerService: offerService,
		supabaseURL:  os.Getenv("SUPABASE_URL"),
		anonKey:      os.Getenv("SUPABASE_ANON_KEY"),
	}
}

func (h *offerHandler) uploadToSupabase(c *gin.Context, fileKey string) (string, error) {
	file, header, err := c.Request.FormFile(fileKey)
	if err != nil {
		return "", err // Mengembalikan error jika file tidak ditemukan/tidak di-upload
	}
	defer file.Close()

	// Membuat nama unik file di storage agar tidak bentrok
	uniqueFileName := fmt.Sprintf("room_%d_%s", time.Now().UnixNano(), header.Filename)
	uploadURL := fmt.Sprintf("%s/storage/v1/object/offer-image/%s", h.supabaseURL, uniqueFileName)

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
	return fmt.Sprintf("%s/storage/v1/object/public/offer-image/%s", h.supabaseURL, uniqueFileName), nil
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
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("=== BINDING ERROR IN CREATE OFFER ===")
		fmt.Println(err.Error())
		fmt.Println("=====================================")
		c.JSON(400, gin.H{"message": "Bad Request", "error": err.Error()})
		return
	}

	imageURL, err := h.uploadToSupabase(c, "image")
	if err == nil {
		req.Image = &imageURL
	} else if err != http.ErrMissingFile && err != http.ErrNotMultipart {
		c.JSON(400, gin.H{"success": false, "message": "Failed to upload image: " + err.Error()})
		return
	}

	err = h.offerService.Create(&req)
	if err != nil {
		c.JSON(500, gin.H{"message": "Error creating offer data", "error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Offer data created successfully!"})
}

func (h *offerHandler) Update(c *gin.Context) {

	id := c.Param("id")
	var req dto.UpdateOfferRequest
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("=== BINDING ERROR IN UPDATE OFFER ===")
		fmt.Println(err.Error())
		fmt.Println("=====================================")
		c.JSON(400, gin.H{"message": "Error binding JSON", "error": err.Error()})
		return
	}

	imageURL, err := h.uploadToSupabase(c, "image")
	if err == nil {
		req.Image = &imageURL
	} else if err != http.ErrMissingFile && err != http.ErrNotMultipart {
		c.JSON(400, gin.H{"success": false, "message": "Failed to upload image: " + err.Error()})
		return
	}

	err = h.offerService.Update(id, &req)
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
