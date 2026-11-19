package unit

import (
	"citramascoweb-backend/internal/dto"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

	// RoomStatus operations
	GetAllStatuses() ([]RoomStatus, error)
	GetStatusById(id string) (*RoomStatus, error)
	CreateStatus(req *dto.CreateRoomStatusRequest) error
	UpdateStatusModel(id string, req *dto.UpdateRoomStatusRequest) error
	DeleteStatus(id string) error
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
		return errors.New("room number is already in use")
	}

	status := "available"
	if req.Status != "" {
		status = strings.ToLower(req.Status)
	}

	newUnit := &RoomUnit{
		Id:           uuid.New().String(),
		RoomNumber:   req.RoomNumber,
		RoomId:       req.RoomId,
		FloorId:      req.FloorId,
		RoomStatusId: req.RoomStatusId,
		Status:       status,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.repo.Create(newUnit)
}

func (s *roomUnitService) Update(id string, req *dto.UpdateRoomUnitRequest) error {
	existing, err := s.repo.GetById(id)
	if err != nil {
		return errors.New("room unit not found")
	}

	if req.RoomNumber != "" && req.RoomNumber != existing.RoomNumber {
		// Cek keunikan nomor kamar baru
		check, _ := s.repo.GetByRoomNumber(req.RoomNumber)
		if check != nil && check.Id != id {
			return errors.New("room number is already in use")
		}
		existing.RoomNumber = req.RoomNumber
	}

	if req.RoomId != "" {
		existing.RoomId = req.RoomId
	}
	if req.FloorId != "" {
		existing.FloorId = req.FloorId
	}
	if req.RoomStatusId != nil {
		existing.RoomStatusId = req.RoomStatusId
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
		"vc":          true,
		"vd":          true,
		"oc":          true,
		"od":          true,
		"ooo":         true,
		"oos":         true,
		"dnd":         true,
		"so":          true,
	}
	if !validStatuses[status] {
		return errors.New("status tidak valid (harus: available, occupied, dirty, maintenance, atau kode status VC, VD, OC, OD, OOO, OOS)")
	}

	_, err := s.repo.GetById(id)
	if err != nil {
		return errors.New("room unit not found")
	}

	// Sinkronisasi status kode ke RoomStatus jika cocok
	codeMap := map[string]string{
		"available":   "VC",
		"occupied":    "OC",
		"dirty":       "VD",
		"maintenance": "OOO",
	}
	targetCode := strings.ToUpper(status)
	if mapped, exists := codeMap[status]; exists {
		targetCode = mapped
	}
	if roomStat, err := s.repo.GetStatusByCode(targetCode); err == nil && roomStat != nil {
		_ = s.repo.UpdateRoomStatusId(id, &roomStat.Id)
	}

	return s.repo.UpdateStatus(id, status)
}

func (s *roomUnitService) Delete(id string) error {
	_, err := s.repo.GetById(id)
	if err != nil {
		return errors.New("room unit not found")
	}
	return s.repo.Delete(id)
}

// ---------------------------------------------------------------------------
// RoomStatus Service
// ---------------------------------------------------------------------------

func (s *roomUnitService) GetAllStatuses() ([]RoomStatus, error) {
	return s.repo.GetAllStatuses()
}

func (s *roomUnitService) GetStatusById(id string) (*RoomStatus, error) {
	return s.repo.GetStatusById(id)
}

func (s *roomUnitService) CreateStatus(req *dto.CreateRoomStatusRequest) error {
	code := strings.ToUpper(strings.TrimSpace(req.Code))
	if code == "" {
		return errors.New("kode status wajib diisi (misal: VC, VD, OOO)")
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return errors.New("nama status wajib diisi (misal: Vacant Clean)")
	}

	existing, _ := s.repo.GetStatusByCode(code)
	if existing != nil && existing.Id != "" {
		return errors.New("kode status sudah terdaftar")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	status := &RoomStatus{
		Id:          uuid.New().String(),
		Code:        code,
		Name:        name,
		Description: req.Description,
		IsActive:    isActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.repo.CreateStatus(status)
}

func (s *roomUnitService) UpdateStatusModel(id string, req *dto.UpdateRoomStatusRequest) error {
	existing, err := s.repo.GetStatusById(id)
	if err != nil {
		return errors.New("room status tidak ditemukan")
	}

	if req.Code != "" {
		code := strings.ToUpper(strings.TrimSpace(req.Code))
		if code != existing.Code {
			check, _ := s.repo.GetStatusByCode(code)
			if check != nil && check.Id != id {
				return errors.New("kode status sudah terdaftar")
			}
			existing.Code = code
		}
	}

	if req.Name != "" {
		existing.Name = strings.TrimSpace(req.Name)
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	existing.UpdatedAt = time.Now()

	return s.repo.UpdateStatusModel(existing, id)
}

func (s *roomUnitService) DeleteStatus(id string) error {
	_, err := s.repo.GetStatusById(id)
	if err != nil {
		return errors.New("room status tidak ditemukan")
	}
	return s.repo.DeleteStatus(id)
}

func SeedDefaultRoomStatuses(db *gorm.DB) {
	defaults := []struct {
		Code        string
		Name        string
		Description string
	}{
		{Code: "VC", Name: "Vacant Clean", Description: "Kamar kosong, bersih, dan siap untuk check-in"},
		{Code: "VD", Name: "Vacant Dirty", Description: "Kamar kosong namun kotor, perlu dibersihkan oleh housekeeping"},
		{Code: "OC", Name: "Occupied Clean", Description: "Kamar terisi tamu dan telah dibersihkan"},
		{Code: "OD", Name: "Occupied Dirty", Description: "Kamar terisi tamu dan membutuhkan pembersihan harian"},
		{Code: "OOO", Name: "Out of Order", Description: "Kamar rusak berat atau sedang dalam renovasi / perbaikan"},
		{Code: "OOS", Name: "Out of Service", Description: "Kamar tidak dijual sementara untuk perbaikan ringan"},
		{Code: "DND", Name: "Do Not Disturb", Description: "Tamu memasang tanda Jangan Diganggu"},
		{Code: "SO", Name: "Sleep Out", Description: "Tamu terdaftar menginap tetapi tidak tidur di kamar"},
	}

	for _, item := range defaults {
		var count int64
		db.Model(&RoomStatus{}).Where("code = ?", item.Code).Count(&count)
		if count == 0 {
			db.Create(&RoomStatus{
				Id:          uuid.New().String(),
				Code:        item.Code,
				Name:        item.Name,
				Description: item.Description,
				IsActive:    true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			})
		}
	}
}
