package facility

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type FacilityService interface {
	GetAllFacilities() ([]Facility, error)
	GetFacilityById(id string) (*Facility, error)
	CreateFacility(req *CreateFacilityRequest) (*Facility, error)
	UpdateFacility(id string, req *UpdateFacilityRequest) (*Facility, error)
	DeleteFacility(id string) error

	GetFacilitiesByRoomType(roomTypeId string) ([]Facility, error)
	SyncRoomTypeFacilities(roomTypeId string, req *AssignFacilitiesRequest) error

	GetImagesByRoomType(roomTypeId string) ([]RoomTypeImage, error)
	CreateRoomTypeImage(roomTypeId string, req *CreateRoomTypeImageRequest) (*RoomTypeImage, error)
	DeleteRoomTypeImage(id string) error
	SetPrimaryImage(roomTypeId, imageId string) error
}

type facilityService struct {
	repo FacilityRepository
}

func NewFacilityService(repo FacilityRepository) FacilityService {
	return &facilityService{repo: repo}
}

func (s *facilityService) GetAllFacilities() ([]Facility, error) {
	return s.repo.FindAllFacilities()
}

func (s *facilityService) GetFacilityById(id string) (*Facility, error) {
	return s.repo.FindFacilityById(id)
}

func (s *facilityService) CreateFacility(req *CreateFacilityRequest) (*Facility, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	entity := &Facility{
		Id:        uuid.New().String(),
		Name:      req.Name,
		Icon:      req.Icon,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.CreateFacility(entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *facilityService) UpdateFacility(id string, req *UpdateFacilityRequest) (*Facility, error) {
	entity, err := s.repo.FindFacilityById(id)
	if err != nil {
		return nil, err
	}
	entity.Name = req.Name
	entity.Icon = req.Icon
	entity.UpdatedAt = time.Now()
	if err := s.repo.UpdateFacility(entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *facilityService) DeleteFacility(id string) error {
	return s.repo.DeleteFacility(id)
}

func (s *facilityService) GetFacilitiesByRoomType(roomTypeId string) ([]Facility, error) {
	return s.repo.GetFacilitiesByRoomTypeId(roomTypeId)
}

func (s *facilityService) SyncRoomTypeFacilities(roomTypeId string, req *AssignFacilitiesRequest) error {
	return s.repo.SyncRoomTypeFacilities(roomTypeId, req.FacilityIds)
}

func (s *facilityService) GetImagesByRoomType(roomTypeId string) ([]RoomTypeImage, error) {
	return s.repo.GetImagesByRoomTypeId(roomTypeId)
}

func (s *facilityService) CreateRoomTypeImage(roomTypeId string, req *CreateRoomTypeImageRequest) (*RoomTypeImage, error) {
	if req.ImageUrl == "" {
		return nil, errors.New("image_url is required")
	}
	img := &RoomTypeImage{
		Id:         uuid.New().String(),
		RoomTypeId: roomTypeId,
		ImageUrl:   req.ImageUrl,
		IsPrimary:  req.IsPrimary,
		SortOrder:  req.SortOrder,
		CreatedAt:  time.Now(),
	}
	if err := s.repo.CreateRoomTypeImage(img); err != nil {
		return nil, err
	}
	return img, nil
}

func (s *facilityService) DeleteRoomTypeImage(id string) error {
	return s.repo.DeleteRoomTypeImage(id)
}

func (s *facilityService) SetPrimaryImage(roomTypeId, imageId string) error {
	return s.repo.SetPrimaryImage(roomTypeId, imageId)
}
