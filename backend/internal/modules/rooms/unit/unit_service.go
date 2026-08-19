package unit

import (
	"citramascoweb-backend/internal/dto"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RoomUnitServiceInterface interface {
	GetAll() ([]RoomUnit, error)
	GetById(id string) (*RoomUnit, error)
	GetByRoomId(roomId string) ([]RoomUnit, error)
	GetByFloorId(floorId string) ([]RoomUnit, error)
	GetNextAvailableUnit(roomId string) (*RoomUnit, error)
	Create(req *dto.CreateRoomUnitRequest) error
	Update(id string, req *dto.UpdateRoomUnitRequest) error
	UpdateStatus(id string, status string) error
	Delete(id string) error
}

type roomUnitService struct {
	repo RoomUnitRepositoryInterface
}

func NewRoomUnitService(repo RoomUnitRepositoryInterface) RoomUnitServiceInterface {
	return &roomUnitService{repo: repo}
}

func (s *roomUnitService) GetAll() ([]RoomUnit, error) {
	return s.repo.GetAll()
}

func (s *roomUnitService) GetById(id string) (*RoomUnit, error) {
	return s.repo.GetById(id)
}

func (s *roomUnitService) GetByRoomId(roomId string) ([]RoomUnit, error) {
	return s.repo.GetByRoomId(roomId)
}

func (s *roomUnitService) GetByFloorId(floorId string) ([]RoomUnit, error) {
	return s.repo.GetByFloorId(floorId)
}

func (s *roomUnitService) GetNextAvailableUnit(roomId string) (*RoomUnit, error) {
	return s.repo.GetNextAvailableUnit(roomId)
}

func (s *roomUnitService) Create(req *dto.CreateRoomUnitRequest) error {
	if req.RoomNumber == "" {
		return errors.New("nomor kamar wajib diisi")
	}
	if req.RoomId == "" {
		return errors.New("tipe kamar wajib dipilih")
	}
	if req.FloorId == "" {
		return errors.New("lantai wajib dipilih")
	}

	// Cek nomor kamar sudah ada
	existing, _ := s.repo.GetByRoomNumber(req.RoomNumber)
	if existing != nil && existing.Id != "" {
		return errors.New("nomor kamar sudah digunakan")
	}

	status := "available"
	if req.Status != "" {
		status = strings.ToLower(req.Status)
	}

	newUnit := &RoomUnit{
		Id:         uuid.New().String(),
		RoomNumber: req.RoomNumber,
		RoomId:     req.RoomId,
		FloorId:    req.FloorId,
		Status:     status,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return s.repo.Create(newUnit)
}

func (s *roomUnitService) Update(id string, req *dto.UpdateRoomUnitRequest) error {
	existing, err := s.repo.GetById(id)
	if err != nil {
		return errors.New("unit kamar tidak ditemukan")
	}

	if req.RoomNumber != "" && req.RoomNumber != existing.RoomNumber {
		// Cek keunikan nomor kamar baru
		check, _ := s.repo.GetByRoomNumber(req.RoomNumber)
		if check != nil && check.Id != id {
			return errors.New("nomor kamar sudah digunakan")
		}
		existing.RoomNumber = req.RoomNumber
	}

	if req.RoomId != "" {
		existing.RoomId = req.RoomId
	}
	if req.FloorId != "" {
		existing.FloorId = req.FloorId
	}
	if req.Status != "" {
		existing.Status = strings.ToLower(req.Status)
	}
	existing.UpdatedAt = time.Now()

	return s.repo.Update(existing, id)
}

func (s *roomUnitService) UpdateStatus(id string, status string) error {
	status = strings.ToLower(strings.TrimSpace(status))
	validStatuses := map[string]bool{
		"available":   true,
		"occupied":    true,
		"dirty":       true,
		"maintenance": true,
	}
	if !validStatuses[status] {
		return errors.New("status tidak valid (harus: available, occupied, dirty, atau maintenance)")
	}

	_, err := s.repo.GetById(id)
	if err != nil {
		return errors.New("unit kamar tidak ditemukan")
	}

	return s.repo.UpdateStatus(id, status)
}

func (s *roomUnitService) Delete(id string) error {
	_, err := s.repo.GetById(id)
	if err != nil {
		return errors.New("unit kamar tidak ditemukan")
	}
	return s.repo.Delete(id)
}
