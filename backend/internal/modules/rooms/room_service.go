package rooms

import (
	"citramascoweb-backend/internal/dto"
	"citramascoweb-backend/pkg/utils"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

type roomService struct {
	roomRepo RoomRepositoryInterface
}

func NewRoomService(roomRepo RoomRepositoryInterface) *roomService {
	return &roomService{roomRepo: roomRepo}
}

func (s *roomService) GetAll() ([]Room, error) {
	rooms, err := s.roomRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *roomService) GetById(id string) (*Room, error) {
	room, err := s.roomRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *roomService) Store(req *dto.CreateRoomRequest) error {
	now := time.Now()

	code, err := utils.GenerateRandomNumberCode(6)
	if err != nil {
		return err
	}

	newRoom := &Room{
		Id:         uuid.New().String(),
		Code:       code,
		Name:       req.Name,
		Slug:       utils.Slugify(req.Name),
		CategoryId: req.CategoryId,
		Price:      req.Price,
		Capacity:   req.Capacity,
		Size:       req.Size,

		TypeId:      req.TypeId,
		Description: req.Description,
		Image:       req.Image, // Menggunakan string URL publik Supabase yang sudah disiapkan handler
		CreatedAt:   now,
	}

	err = s.roomRepo.Create(newRoom)
	if err != nil {
		return err
	}

	return nil
}

func (s *roomService) Update(id string, req *dto.UpdateRoomRequest) error {
	room, err := s.roomRepo.GetById(id)
	if err != nil {
		return err
	}

	room.Name = req.Name
	room.Slug = utils.Slugify(req.Name)
	room.CategoryId = req.CategoryId
	room.Price = req.Price
	room.Capacity = req.Capacity
	room.Size = req.Size

	room.TypeId = req.TypeId
	room.Description = req.Description

	// PENGAMAN: Hanya ganti nilai gambar jika payload DTO membawa link URL baru dari handler.
	if req.Image != "" {
		room.Image = req.Image
	}

	err = s.roomRepo.Update(room, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *roomService) calculateOfferDiscount(isOffer *bool, offerCode *string, price int, excludeRoomId string) (*bool, *int, *string, error) {
	if isOffer == nil || !*isOffer {
		isOfferFalse := false
		return &isOfferFalse, nil, nil, nil
	}

	if offerCode == nil || *offerCode == "" {
		return nil, nil, nil, errors.New("Offer code is required when is_offer is active!")
	}

	offerData, err := s.roomRepo.GetOfferByCode(*offerCode)
	if err != nil {
		return nil, nil, nil, errors.New("Invalid offer code!")
	}

	// Check ValidStart and ValidEnd ranges
	now := time.Now()
	if offerData.ValidStart != nil && now.Before(time.Time(*offerData.ValidStart)) {
		return nil, nil, nil, errors.New("Offer code is not active yet!")
	}
	if offerData.ValidEnd != nil && now.After(time.Time(*offerData.ValidEnd)) {
		return nil, nil, nil, errors.New("Offer code has expired!")
	}

	// Calculate discounted price
	discountVal := 0
	if offerData.Discount != nil {
		discountVal = *offerData.Discount
	}
	priceAfterDiscount := price - (price * discountVal / 100)

	isOfferTrue := true
	return &isOfferTrue, &priceAfterDiscount, offerCode, nil
}

func (s *roomService) Delete(id string) error {
	err := s.roomRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *roomService) FilerByCategory(categoryId string) ([]Room, error) {
	rooms, err := s.roomRepo.FilerByCategory(categoryId)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *roomService) FilterByType(typeId string) ([]Room, error) {
	rooms, err := s.roomRepo.FilterByType(typeId)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *roomService) Filter(status, availabilityStatus string, checkinDate, checkoutDate time.Time) ([]Room, error) {
	rooms, err := s.roomRepo.Filter(status, availabilityStatus, checkinDate, checkoutDate)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *roomService) UpdateStatus(id string, statusStr string) error {
	log.Printf("[DEBUG][ROOM-SERVICE] Memproses request perubahan status kamar ID: %s, Payload: '%s'", id, statusStr)

	// Konversi string input menjadi tipe typed-enum RoomStatus
	var targetStatus RoomStatus
	switch strings.ToLower(statusStr) {
	case "active", "available":
		targetStatus = RoomStatusActive
	case "maintenance":
		targetStatus = RoomStatusMaintenance
	case "dirty":
		targetStatus = RoomStatusDirty
	default:
		log.Printf("[WARN][ROOM-SERVICE] Payload status tidak dikenal: '%s'", statusStr)
		return fmt.Errorf("status fisik kamar '%s' tidak valid", statusStr)
	}

	// Ambil data kamar terlebih dahulu untuk memastikan unitnya eksis
	room, err := s.roomRepo.GetById(id)
	if err != nil {
		log.Printf("[ERROR][ROOM-SERVICE] Kamar dengan ID %s tidak ditemukan di sistem", id)
		return fmt.Errorf("unit kamar tidak ditemukan: %v", err)
	}

	log.Printf("[DEBUG][ROOM-SERVICE] Kamar '%s' ditemukan. Mengubah status lama '%s' -> '%s'", room.Code, room.Status, targetStatus)
	return s.roomRepo.UpdateStatus(id, targetStatus)
}
