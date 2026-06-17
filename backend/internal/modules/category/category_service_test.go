package category

import (
	"citramascoweb-backend/internal/dto"
	"errors"
	"testing"
	"time"
)

type mockCategoryRepository struct {
	GetAllFn    func() ([]Category, error)
	StoreFn     func(category *Category) error
	UpdateFn    func(id string, category *Category) error
	DeleteFn    func(id string) error
	GetByIdFn   func(id string) (*Category, error)
	GetBySlugFn func(slug string) (*Category, error)
}

func (m *mockCategoryRepository) GetAll() ([]Category, error) {
	return m.GetAllFn()
}
func (m *mockCategoryRepository) Store(category *Category) error {
	return m.StoreFn(category)
}
func (m *mockCategoryRepository) Update(id string, category *Category) error {
	return m.UpdateFn(id, category)
}
func (m *mockCategoryRepository) Delete(id string) error {
	return m.DeleteFn(id)
}
func (m *mockCategoryRepository) GetById(id string) (*Category, error) {
	return m.GetByIdFn(id)
}
func (m *mockCategoryRepository) GetBySlug(slug string) (*Category, error) {
	return m.GetBySlugFn(slug)
}

func TestCategoryService_FindById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedCategory := &Category{
			Id:   "test-id",
			Name: "Test Category",
			Slug: "test-category",
		}

		mockRepo := &mockCategoryRepository{
			GetByIdFn: func(id string) (*Category, error) {
				if id != "test-id" {
					t.Errorf("expected id 'test-id', got '%s'", id)
				}
				return expectedCategory, nil
			},
		}

		service := NewCategoryService(mockRepo)
		result, err := service.FindById("test-id")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Id != expectedCategory.Id || result.Name != expectedCategory.Name {
			t.Errorf("expected category %+v, got %+v", expectedCategory, result)
		}
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			GetByIdFn: func(id string) (*Category, error) {
				return nil, errors.New("database error")
			},
		}

		service := NewCategoryService(mockRepo)
		result, err := service.FindById("test-id")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if result != nil {
			t.Errorf("expected nil result, got %+v", result)
		}
	})
}

func TestCategoryService_FindAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedCategories := []Category{
			{Id: "1", Name: "Cat 1", Slug: "cat-1"},
			{Id: "2", Name: "Cat 2", Slug: "cat-2"},
		}

		mockRepo := &mockCategoryRepository{
			GetAllFn: func() ([]Category, error) {
				return expectedCategories, nil
			},
		}

		service := NewCategoryService(mockRepo)
		results, err := service.FindAll()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(results) != len(expectedCategories) {
			t.Errorf("expected %d categories, got %d", len(expectedCategories), len(results))
		}
	})
}

func TestCategoryService_FindBySlug(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedCategory := &Category{
			Id:   "test-id",
			Name: "Test Category",
			Slug: "test-category",
		}

		mockRepo := &mockCategoryRepository{
			GetBySlugFn: func(slug string) (*Category, error) {
				if slug != "test-category" {
					t.Errorf("expected slug 'test-category', got '%s'", slug)
				}
				return expectedCategory, nil
			},
		}

		service := NewCategoryService(mockRepo)
		result, err := service.FindBySlug("test-category")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Slug != expectedCategory.Slug {
			t.Errorf("expected slug '%s', got '%s'", expectedCategory.Slug, result.Slug)
		}
	})
}

func TestCategoryService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := &dto.CreateCategoryRequest{
			Name: "New Category",
		}

		var storedCategory *Category
		mockRepo := &mockCategoryRepository{
			StoreFn: func(category *Category) error {
				storedCategory = category
				return nil
			},
		}

		service := NewCategoryService(mockRepo)
		err := service.Create(req)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if storedCategory == nil {
			t.Fatal("expected category to be stored, but store function was not called or stored nil")
		}
		if storedCategory.Name != req.Name {
			t.Errorf("expected category name '%s', got '%s'", req.Name, storedCategory.Name)
		}
		if storedCategory.Slug != "new-category" {
			t.Errorf("expected slug 'new-category', got '%s'", storedCategory.Slug)
		}
		if storedCategory.Id == "" {
			t.Error("expected generated category ID to be non-empty")
		}
	})
}

func TestCategoryService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		existingCategory := &Category{
			Id:        "test-id",
			Name:      "Old Name",
			Slug:      "old-name",
			CreatedAt: time.Now(),
		}

		req := &dto.UpdateCategoryRequest{
			Name: "Updated Name",
		}

		var updatedId string
		var updatedCategory *Category

		mockRepo := &mockCategoryRepository{
			GetByIdFn: func(id string) (*Category, error) {
				return existingCategory, nil
			},
			UpdateFn: func(id string, category *Category) error {
				updatedId = id
				updatedCategory = category
				return nil
			},
		}

		service := NewCategoryService(mockRepo)
		err := service.Update("test-id", req)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updatedId != "test-id" {
			t.Errorf("expected updated id 'test-id', got '%s'", updatedId)
		}
		if updatedCategory == nil {
			t.Fatal("expected updated category to be set")
		}
		if updatedCategory.Name != "Updated Name" {
			t.Errorf("expected name 'Updated Name', got '%s'", updatedCategory.Name)
		}
		if updatedCategory.Slug != "updated-name" {
			t.Errorf("expected slug 'updated-name', got '%s'", updatedCategory.Slug)
		}
	})
}

func TestCategoryService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var deletedId string
		mockRepo := &mockCategoryRepository{
			DeleteFn: func(id string) error {
				deletedId = id
				return nil
			},
		}

		service := NewCategoryService(mockRepo)
		err := service.Delete("test-id")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if deletedId != "test-id" {
			t.Errorf("expected deleted id 'test-id', got '%s'", deletedId)
		}
	})
}
