package guest

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type GuestService interface {
	GetAllIdentityTypes() ([]IdentityType, error)
	GetIdentityTypeById(id string) (*IdentityType, error)
	CreateIdentityType(req *CreateIdentityTypeRequest) (*IdentityType, error)
	UpdateIdentityType(id string, req *UpdateIdentityTypeRequest) (*IdentityType, error)
	DeleteIdentityType(id string) error

	GetAllGuests(search string, page, limit int) ([]Guest, int64, error)
	GetGuestById(id string) (*Guest, error)
	CreateGuest(req *CreateGuestRequest) (*Guest, error)
	UpdateGuest(id string, req *UpdateGuestRequest) (*Guest, error)
	DeleteGuest(id string) error
}

type guestService struct {
	repo GuestRepository
}

func NewGuestService(repo GuestRepository) GuestService {
	return &guestService{repo: repo}
}

func (s *guestService) GetAllIdentityTypes() ([]IdentityType, error) {
	return s.repo.FindIdentityTypes()
}

func (s *guestService) GetIdentityTypeById(id string) (*IdentityType, error) {
	return s.repo.FindIdentityTypeById(id)
}

func (s *guestService) CreateIdentityType(req *CreateIdentityTypeRequest) (*IdentityType, error) {
	entity := &IdentityType{
		Id:        uuid.New().String(),
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.CreateIdentityType(entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *guestService) UpdateIdentityType(id string, req *UpdateIdentityTypeRequest) (*IdentityType, error) {
	entity, err := s.repo.FindIdentityTypeById(id)
	if err != nil {
		return nil, err
	}
	entity.Name = req.Name
	entity.UpdatedAt = time.Now()
	if err := s.repo.UpdateIdentityType(entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *guestService) DeleteIdentityType(id string) error {
	return s.repo.DeleteIdentityType(id)
}

func (s *guestService) GetAllGuests(search string, page, limit int) ([]Guest, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit
	return s.repo.FindGuests(search, limit, offset)
}

func (s *guestService) GetGuestById(id string) (*Guest, error) {
	return s.repo.FindGuestById(id)
}

func (s *guestService) CreateGuest(req *CreateGuestRequest) (*Guest, error) {
	if req.FirstName == "" {
		return nil, errors.New("first_name is required")
	}

	entity := &Guest{
		Id:             uuid.New().String(),
		Title:          req.Title,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		IdentityTypeId: req.IdentityTypeId,
		IdentityNumber: req.IdentityNumber,
		IdentityImage:  req.IdentityImage,
		Email:          req.Email,
		Phone:          req.Phone,
		Nationality:    req.Nationality,
		Address:        req.Address,
		GuestType:      req.GuestType,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if entity.Nationality == "" {
		entity.Nationality = "Indonesia"
	}
	if entity.GuestType == "" {
		entity.GuestType = "regular"
	}

	if err := s.repo.CreateGuest(entity); err != nil {
		return nil, err
	}
	return s.repo.FindGuestById(entity.Id)
}

func (s *guestService) UpdateGuest(id string, req *UpdateGuestRequest) (*Guest, error) {
	entity, err := s.repo.FindGuestById(id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != "" {
		entity.FirstName = req.FirstName
	}
	entity.Title = req.Title
	entity.LastName = req.LastName
	if req.IdentityTypeId != nil {
		entity.IdentityTypeId = req.IdentityTypeId
	}
	if req.IdentityNumber != "" {
		entity.IdentityNumber = req.IdentityNumber
	}
	if req.IdentityImage != "" {
		entity.IdentityImage = req.IdentityImage
	}
	if req.Email != "" {
		entity.Email = req.Email
	}
	if req.Phone != "" {
		entity.Phone = req.Phone
	}
	if req.Nationality != "" {
		entity.Nationality = req.Nationality
	}
	if req.Address != "" {
		entity.Address = req.Address
	}
	if req.GuestType != "" {
		entity.GuestType = req.GuestType
	}
	entity.UpdatedAt = time.Now()

	if err := s.repo.UpdateGuest(entity); err != nil {
		return nil, err
	}
	return s.repo.FindGuestById(id)
}

func (s *guestService) DeleteGuest(id string) error {
	return s.repo.DeleteGuest(id)
}
