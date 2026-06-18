package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"citramascoweb-backend/internal/dto"

	"github.com/gin-gonic/gin"
)

func setupTestRouter(handler *typeHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/types", handler.GetAll)
	r.GET("/types/search", handler.Search)
	r.GET("/type/:id", handler.GetById)
	r.POST("/type", handler.Store)
	r.PATCH("/type/:id", handler.Update)
	r.DELETE("/type/:id", handler.Delete)
	return r
}

func TestTypeHandler_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockTypeRepository{
			GetAllFn: func() ([]Types, error) {
				return []Types{{Id: "1", Name: "Type A"}}, nil
			},
		}
		service := NewTypeService(mockRepo)
		handler := NewTypeHandler(service)
		router := setupTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/types", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := &mockTypeRepository{
			GetAllFn: func() ([]Types, error) {
				return nil, errors.New("db error")
			},
		}
		service := NewTypeService(mockRepo)
		handler := NewTypeHandler(service)
		router := setupTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/types", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestTypeHandler_GetById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockTypeRepository{
			GetByIdFn: func(id string) (*Types, error) {
				return &Types{Id: id, Name: "Type A"}, nil
			},
		}
		service := NewTypeService(mockRepo)
		handler := NewTypeHandler(service)
		router := setupTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/type/123", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := &mockTypeRepository{
			GetByIdFn: func(id string) (*Types, error) {
				return nil, errors.New("not found")
			},
		}
		service := NewTypeService(mockRepo)
		handler := NewTypeHandler(service)
		router := setupTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/type/123", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestTypeHandler_Store(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockTypeRepository{
			CreateFn: func(types *Types) error {
				return nil
			},
		}
		service := NewTypeService(mockRepo)
		handler := NewTypeHandler(service)
		router := setupTestRouter(handler)

		body := dto.CreateTypesRequest{Name: "New Type"}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/type", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("bad request - invalid json", func(t *testing.T) {
		service := NewTypeService(&mockTypeRepository{})
		handler := NewTypeHandler(service)
		router := setupTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/type", bytes.NewBufferString("{invalid"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		mockRepo := &mockTypeRepository{
			CreateFn: func(types *Types) error {
				return errors.New("db error")
			},
		}
		service := NewTypeService(mockRepo)
		handler := NewTypeHandler(service)
		router := setupTestRouter(handler)

		body := dto.CreateTypesRequest{Name: "New Type"}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/type", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestTypeHandler_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockTypeRepository{
			GetByIdFn: func(id string) (*Types, error) {
				return &Types{Id: id, Name: "Old Name"}, nil
			},
			UpdateFn: func(id string, types *Types) error {
				return nil
			},
		}
		service := NewTypeService(mockRepo)
		handler := NewTypeHandler(service)
		router := setupTestRouter(handler)

		body := dto.UpdateTypesRequest{Name: "New Name"}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/type/123", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d %s", w.Code, w.Body.String())
		}
	})

	t.Run("bad request - invalid json", func(t *testing.T) {
		service := NewTypeService(&mockTypeRepository{})
		handler := NewTypeHandler(service)
		router := setupTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/type/123", bytes.NewBufferString("{invalid"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		mockRepo := &mockTypeRepository{
			GetByIdFn: func(id string) (*Types, error) {
				return nil, errors.New("not found")
			},
		}
		service := NewTypeService(mockRepo)
		handler := NewTypeHandler(service)
		router := setupTestRouter(handler)

		body := dto.UpdateTypesRequest{Name: "New Name"}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/type/123", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestTypeHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockTypeRepository{
			DeleteFn: func(id string) error {
				return nil
			},
		}
		service := NewTypeService(mockRepo)
		handler := NewTypeHandler(service)
		router := setupTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/type/123", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		mockRepo := &mockTypeRepository{
			DeleteFn: func(id string) error {
				return errors.New("db error")
			},
		}
		service := NewTypeService(mockRepo)
		handler := NewTypeHandler(service)
		router := setupTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/type/123", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestTypeHandler_Search(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockTypeRepository{
			SearchFn: func(query string) ([]Types, error) {
				return []Types{{Id: "1", Name: "Type A"}}, nil
			},
		}
		service := NewTypeService(mockRepo)
		handler := NewTypeHandler(service)
		router := setupTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/types/search?query=Type", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		mockRepo := &mockTypeRepository{
			SearchFn: func(query string) ([]Types, error) {
				return nil, errors.New("db error")
			},
		}
		service := NewTypeService(mockRepo)
		handler := NewTypeHandler(service)
		router := setupTestRouter(handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/types/search?query=Type", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}
