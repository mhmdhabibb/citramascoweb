package floor

import (
	"citramascoweb-backend/internal/dto"
	"errors"
	"time"

	"github.com/google/uuid"
)

type FloorServiceInterface interface {
	GetAll() ([]Floor, error)
	GetById(id string) (*Floor, error)
	Create(req *dto.CreateFloorRequest) error
	Update(id string, req *dto.UpdateFloorRequest) error
	Delete(id string) error
}

type floorService struct {
	repo FloorRepositoryInterface
}

func NewFloorService(repo FloorRepositoryInterface) FloorServiceInterface {
	return &floorService{repo: repo}
}

func (s *floorService) GetAll() ([]Floor, error) {
	return s.repo.GetAll()
}

func (s *floorService) GetById(id string) (*Floor, error) {
	return s.repo.GetById(id)
}

func (s *floorService) Create(req *dto.CreateFloorRequest) error {
	if req.Name == "" {
		return errors.New("nama lantai wajib diisi")
	}

	newFloor := &Floor{
		Id:          uuid.New().String(),
		Name:        req.Name,
		FloorNumber: req.FloorNumber,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.repo.Create(newFloor)
}

func (s *floorService) Update(id string, req *dto.UpdateFloorRequest) error {
	existing, err := s.repo.GetById(id)
	if err != nil {
		return errors.New("lantai tidak ditemukan")
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.FloorNumber != 0 {
		existing.FloorNumber = req.FloorNumber
	}
	existing.Description = req.Description
	existing.UpdatedAt = time.Now()

	return s.repo.Update(existing, id)
}

func (s *floorService) Delete(id string) error {
	_, err := s.repo.GetById(id)
	if err != nil {
		return errors.New("lantai tidak ditemukan")
	}
	return s.repo.Delete(id)
}
