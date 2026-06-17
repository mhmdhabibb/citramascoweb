package middlewares

import (
	"citramascoweb-backend/pkg/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	utils.SetJWTSecret([]byte("test-secret-key"))

	t.Run("missing authorization header", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(AuthMiddleware())
		r.GET("/test", func(ctx *gin.Context) {
			ctx.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		c.Request = req
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("invalid token format", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(AuthMiddleware())
		r.GET("/test", func(ctx *gin.Context) {
			ctx.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "InvalidFormat tokenhere")
		c.Request = req
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("valid token", func(t *testing.T) {
		token, err := utils.GenerateToken("user-123", "admin")
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(AuthMiddleware())
		r.GET("/test", func(ctx *gin.Context) {
			userID, _ := ctx.Get("user_id")
			role, _ := ctx.Get("role")
			ctx.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"role":    role,
			})
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		c.Request = req
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("options method bypass", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(AuthMiddleware())
		r.OPTIONS("/test", func(ctx *gin.Context) {
			ctx.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodOptions, "/test", nil)
		c.Request = req
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}
