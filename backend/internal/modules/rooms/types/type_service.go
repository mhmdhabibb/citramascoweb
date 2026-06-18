package types

import (
	"citramascoweb-backend/internal/dto"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
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
