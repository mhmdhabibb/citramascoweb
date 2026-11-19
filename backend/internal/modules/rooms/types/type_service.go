package types

import (
	"citramascoweb-backend/internal/dto"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type typeService struct {
	typeRepository TypeRepositoryInterface
}

func Slugify(text string) string {
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, " ", "-")
	return text
}

func NewTypeService(typeRepository TypeRepositoryInterface) *typeService {
	return &typeService{typeRepository: typeRepository}
}

func (s *typeService) GetAll() ([]Types, error) {
	types, err := s.typeRepository.GetAll()
	if err != nil {
		return nil, err
	}

	return types, nil
}

func (s *typeService) GetById(id string) (*Types, error) {
	types, err := s.typeRepository.GetById(id)
	if err != nil {
		return nil, err
	}

	return types, nil
}

func (s *typeService) Create(req *dto.CreateTypesRequest) error {
	if req.Name == "" {
		return errors.New("Name is required!")
	}

	now := time.Now()
	types := &Types{
		Id:        uuid.New().String(),
		Name:      req.Name,
		Slug:      Slugify(req.Name),
		CreatedAt: now,
	}

	err := s.typeRepository.Create(types)

	if err != nil {
		return err
	}

	return nil
}

func (s *typeService) Update(id string, req *dto.UpdateTypesRequest) error {

	types, err := s.typeRepository.GetById(id)
	if err != nil {
		return err
	}

	types.Name = req.Name
	types.Slug = Slugify(req.Name)

	err = s.typeRepository.Update(id, types)
	if err != nil {
		return err
	}

	return nil

}

func (s *typeService) Delete(id string) error {
	err := s.typeRepository.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (s *typeService) Search(query string) ([]Types, error) {

	types, err := s.typeRepository.Search(query)

	if err != nil {
		return nil, err
	}

	return types, nil

}

// ---------------------------------------------------------------------------
// BedType Service
// ---------------------------------------------------------------------------

func (s *typeService) GetAllBedTypes() ([]BedType, error) {
	return s.typeRepository.GetAllBedTypes()
}

func (s *typeService) GetBedTypeById(id string) (*BedType, error) {
	return s.typeRepository.GetBedTypeById(id)
}

func (s *typeService) CreateBedType(req *dto.CreateBedTypeRequest) error {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return errors.New("name is required")
	}

	slug := req.Slug
	if slug == "" {
		slug = Slugify(name)
	}

	existing, _ := s.typeRepository.GetBedTypeBySlug(slug)
	if existing != nil && existing.Id != "" {
		return errors.New("bed type with this slug already exists")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	bedType := &BedType{
		Id:        uuid.New().String(),
		Name:      name,
		Slug:      slug,
		IsActive:  isActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.typeRepository.CreateBedType(bedType)
}

func (s *typeService) UpdateBedType(id string, req *dto.UpdateBedTypeRequest) error {
	existing, err := s.typeRepository.GetBedTypeById(id)
	if err != nil {
		return errors.New("bed type not found")
	}

	if req.Name != "" {
		existing.Name = strings.TrimSpace(req.Name)
		if req.Slug == "" {
			existing.Slug = Slugify(existing.Name)
		}
	}

	if req.Slug != "" {
		existing.Slug = Slugify(req.Slug)
	}

	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	existing.UpdatedAt = time.Now()

	return s.typeRepository.UpdateBedType(id, existing)
}

func (s *typeService) DeleteBedType(id string) error {
	_, err := s.typeRepository.GetBedTypeById(id)
	if err != nil {
		return errors.New("bed type not found")
	}
	return s.typeRepository.DeleteBedType(id)
}

func SeedDefaultBedTypes(db *gorm.DB) {
	defaults := []struct {
		Name string
		Slug string
	}{
		{Name: "Single Bed", Slug: "single-bed"},
		{Name: "Twin Bed", Slug: "twin-bed"},
		{Name: "Double Bed", Slug: "double-bed"},
		{Name: "Queen Bed", Slug: "queen-bed"},
		{Name: "King Bed", Slug: "king-bed"},
		{Name: "Super King Bed", Slug: "super-king-bed"},
	}

	for _, d := range defaults {
		var count int64
		db.Model(&BedType{}).Where("slug = ?", d.Slug).Count(&count)
		if count == 0 {
			db.Create(&BedType{
				Id:        uuid.New().String(),
				Name:      d.Name,
				Slug:      d.Slug,
				IsActive:  true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	}
}
