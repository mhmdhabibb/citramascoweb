package types

import (
	"citramascoweb-backend/internal/dto"
	"errors"
	"testing"
	"time"
)

type mockTypeRepository struct {
	GetAllFn  func() ([]Types, error)
	GetByIdFn func(id string) (*Types, error)
	CreateFn  func(types *Types) error
	UpdateFn  func(id string, types *Types) error
	DeleteFn  func(id string) error
	SearchFn  func(query string) ([]Types, error)
}

func (m *mockTypeRepository) GetAll() ([]Types, error) {
	return m.GetAllFn()
}

func (m *mockTypeRepository) GetById(id string) (*Types, error) {
	return m.GetByIdFn(id)
}

func (m *mockTypeRepository) Create(types *Types) error {
	return m.CreateFn(types)
}

func (m *mockTypeRepository) Update(id string, types *Types) error {
	return m.UpdateFn(id, types)
}

func (m *mockTypeRepository) Delete(id string) error {
	return m.DeleteFn(id)
}

func (m *mockTypeRepository) Search(query string) ([]Types, error) {
	return m.SearchFn(query)
}

func TestTypeService_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedTypes := []Types{
			{Id: "1", Name: "Type 1", Slug: "type-1"},
			{Id: "2", Name: "Type 2", Slug: "type-2"},
		}

		mockRepo := &mockTypeRepository{
			GetAllFn: func() ([]Types, error) {
				return expectedTypes, nil
			},
		}

		service := NewTypeService(mockRepo)
		results, err := service.GetAll()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(results) != len(expectedTypes) {
			t.Errorf("expected %d types, got %d", len(expectedTypes), len(results))
		}
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := &mockTypeRepository{
			GetAllFn: func() ([]Types, error) {
				return nil, errors.New("database error")
			},
		}

		service := NewTypeService(mockRepo)
		_, err := service.GetAll()

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestTypeService_GetById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedType := &Types{
			Id:   "test-id",
			Name: "Test Type",
			Slug: "test-type",
		}

		mockRepo := &mockTypeRepository{
			GetByIdFn: func(id string) (*Types, error) {
				if id != "test-id" {
					t.Errorf("expected id 'test-id', got '%s'", id)
				}
				return expectedType, nil
			},
		}

		service := NewTypeService(mockRepo)
		result, err := service.GetById("test-id")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Id != expectedType.Id || result.Name != expectedType.Name {
			t.Errorf("expected type %+v, got %+v", expectedType, result)
		}
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := &mockTypeRepository{
			GetByIdFn: func(id string) (*Types, error) {
				return nil, errors.New("not found")
			},
		}

		service := NewTypeService(mockRepo)
		result, err := service.GetById("test-id")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if result != nil {
			t.Errorf("expected nil result, got %+v", result)
		}
	})
}

func TestTypeService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := &dto.CreateTypesRequest{
			Name: "New Type",
		}

		var storedType *Types
		mockRepo := &mockTypeRepository{
			CreateFn: func(types *Types) error {
				storedType = types
				return nil
			},
		}

		service := NewTypeService(mockRepo)
		err := service.Create(req)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if storedType == nil {
			t.Fatal("expected type to be stored, but create function was not called or stored nil")
		}
		if storedType.Name != req.Name {
			t.Errorf("expected type name '%s', got '%s'", req.Name, storedType.Name)
		}
		if storedType.Slug != "new-type" {
			t.Errorf("expected slug 'new-type', got '%s'", storedType.Slug)
		}
		if storedType.Id == "" {
			t.Error("expected generated type ID to be non-empty")
		}
	})

	t.Run("validation error - empty name", func(t *testing.T) {
		req := &dto.CreateTypesRequest{
			Name: "",
		}

		mockRepo := &mockTypeRepository{}
		service := NewTypeService(mockRepo)
		err := service.Create(req)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "Name is required!" {
			t.Errorf("expected error 'Name is required!', got '%s'", err.Error())
		}
	})
}

func TestTypeService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		existingType := &Types{
			Id:        "test-id",
			Name:      "Old Name",
			Slug:      "old-name",
			CreatedAt: time.Now(),
		}

		req := &dto.UpdateTypesRequest{
			Name: "Updated Name",
		}

		var updatedId string
		var updatedType *Types

		mockRepo := &mockTypeRepository{
			GetByIdFn: func(id string) (*Types, error) {
				return existingType, nil
			},
			UpdateFn: func(id string, types *Types) error {
				updatedId = id
				updatedType = types
				return nil
			},
		}

		service := NewTypeService(mockRepo)
		err := service.Update("test-id", req)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updatedId != "test-id" {
			t.Errorf("expected updated id 'test-id', got '%s'", updatedId)
		}
		if updatedType == nil {
			t.Fatal("expected updated type to be set")
		}
		if updatedType.Name != "Updated Name" {
			t.Errorf("expected name 'Updated Name', got '%s'", updatedType.Name)
		}
		if updatedType.Slug != "updated-name" {
			t.Errorf("expected slug 'updated-name', got '%s'", updatedType.Slug)
		}
	})
}

func TestTypeService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var deletedId string
		mockRepo := &mockTypeRepository{
			DeleteFn: func(id string) error {
				deletedId = id
				return nil
			},
		}

		service := NewTypeService(mockRepo)
		err := service.Delete("test-id")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if deletedId != "test-id" {
			t.Errorf("expected deleted id 'test-id', got '%s'", deletedId)
		}
	})
}

func TestTypeService_Search(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedTypes := []Types{
			{Id: "1", Name: "Type 1", Slug: "type-1"},
		}

		mockRepo := &mockTypeRepository{
			SearchFn: func(query string) ([]Types, error) {
				if query != "Type" {
					t.Errorf("expected query 'Type', got '%s'", query)
				}
				return expectedTypes, nil
			},
		}

		service := NewTypeService(mockRepo)
		results, err := service.Search("Type")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(results) != len(expectedTypes) {
			t.Errorf("expected %d results, got %d", len(expectedTypes), len(results))
		}
	})
}
