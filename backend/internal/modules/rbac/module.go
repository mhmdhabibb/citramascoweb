package rbac

import (
	"citramascoweb-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	Handler *rbacHandler
}

func InitModule(db *gorm.DB) *Module {
	repo := NewRbacRepository(db)
	service := NewRbacService(repo)
	handler := NewRbacHandler(service)

	return &Module{
		Handler: handler,
	}
}

func (m *Module) RbacRoutes(router *gin.RouterGroup) {
	// Roles routes
	roles := router.Group("/roles")
	{
		roles.GET("", m.Handler.GetAllRoles)
		roles.GET("/:id", m.Handler.GetRoleById)
		roles.POST("", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.CreateRole)
		roles.PUT("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.UpdateRole)
		roles.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.DeleteRole)
		roles.POST("/:id/permissions", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.SyncRolePermissions)
	}

	// Permissions routes
	permissions := router.Group("/permissions")
	{
		permissions.GET("", m.Handler.GetAllPermissions)
		permissions.GET("/modules", m.Handler.GetModules)
		permissions.GET("/:id", m.Handler.GetPermissionById)
		permissions.POST("", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.CreatePermission)
		permissions.PUT("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.UpdatePermission)
		permissions.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.DeletePermission)
	}

	// User Roles routes (mapping external user_id to role_id)
	userRoles := router.Group("/user-roles", middlewares.AuthMiddleware())
	{
		userRoles.GET("/user/:userId", m.Handler.GetUserRoles)
		userRoles.GET("/user/:userId/permissions", m.Handler.GetUserPermissions)
		userRoles.GET("/role/:roleId", m.Handler.GetUsersByRole)
		userRoles.POST("", middlewares.RoleMiddleware("admin", "manager"), m.Handler.AssignUserRole)
		userRoles.DELETE("/:userId/:roleId", middlewares.RoleMiddleware("admin", "manager"), m.Handler.RevokeUserRole)
	}
}

