package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRoleMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("missing role in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(RoleMiddleware("admin"))
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

	t.Run("user has forbidden role", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		// Create a middleware to inject a non-allowed role first
		r.Use(func(ctx *gin.Context) {
			ctx.Set("role", "user")
			ctx.Next()
		})
		r.Use(RoleMiddleware("admin"))
		r.GET("/test", func(ctx *gin.Context) {
			ctx.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		c.Request = req
		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("expected status 403, got %d", w.Code)
		}
	})

	t.Run("user has allowed role", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		// Create a middleware to inject an allowed role first
		r.Use(func(ctx *gin.Context) {
			ctx.Set("role", "admin")
			ctx.Next()
		})
		r.Use(RoleMiddleware("admin", "superadmin"))
		r.GET("/test", func(ctx *gin.Context) {
			ctx.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		c.Request = req
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}
