package category

import (
	"citramascoweb-backend/internal/dto"
	"strings"
	"time"

	"github.com/google/uuid"
)

type categoryService struct {
	categoryRepo CategoryRepositoryInterface
}

func Slugify(text string) string {
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, " ", "-")
	return text
}

func NewCategoryService(categoryRepo CategoryRepositoryInterface) *categoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) FindById(id string) (*Category, error) {
	category, err := s.categoryRepo.GetById(id)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) FindAll() ([]Category, error) {
	categories, err := s.categoryRepo.GetAll()

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *categoryService) FindBySlug(slug string) (*Category, error) {
	category, err := s.categoryRepo.GetBySlug(slug)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Create(req *dto.CreateCategoryRequest) error {

	now := time.Now()
	category := &Category{
		Id:        uuid.New().String(),
		Name:      req.Name,
		Slug:      Slugify(req.Name),
		CreatedAt: now,
		IsDeleted: false,
	}

	err := s.categoryRepo.Store(category)
	if err != nil {
		return err
	}

	return nil
}

func (s *categoryService) Update(id string, req *dto.UpdateCategoryRequest) error {

	category, err := s.categoryRepo.GetById(id)
	if err != nil {
		return err
	}

	category.Name = req.Name
	category.Slug = Slugify(req.Name)

	err = s.categoryRepo.Update(id, category)
	if err != nil {
		return err
	}

	return nil

}

func (s *categoryService) Delete(id string) error {
	err := s.categoryRepo.Delete(id)

	if err != nil {
		return err
	}

	return nil
}
