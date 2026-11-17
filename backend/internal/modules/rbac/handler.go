package rbac

import (
	"citramascoweb-backend/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type rbacHandler struct {
	service RbacService
}

func NewRbacHandler(service RbacService) *rbacHandler {
	return &rbacHandler{service: service}
}

// -------------------------------------------------------------
// ROLES
// -------------------------------------------------------------

func (h *rbacHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.service.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": roles})
}

func (h *rbacHandler) GetRoleById(c *gin.Context) {
	id := c.Param("id")
	role, err := h.service.GetRoleById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": role})
}

func (h *rbacHandler) CreateRole(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	role, err := h.service.CreateRole(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Role created successfully", "data": role})
}

func (h *rbacHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	role, err := h.service.UpdateRole(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Role updated successfully", "data": role})
}

func (h *rbacHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteRole(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Role deleted successfully"})
}

func (h *rbacHandler) SyncRolePermissions(c *gin.Context) {
	id := c.Param("id")
	var req dto.SyncRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.service.SyncRolePermissions(id, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Role permissions synchronized successfully"})
}

// -------------------------------------------------------------
// PERMISSIONS
// -------------------------------------------------------------

func (h *rbacHandler) GetAllPermissions(c *gin.Context) {
	module := c.Query("module")
	permissions, err := h.service.GetAllPermissions(module)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": permissions})
}

func (h *rbacHandler) GetPermissionById(c *gin.Context) {
	id := c.Param("id")
	permission, err := h.service.GetPermissionById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": permission})
}

func (h *rbacHandler) CreatePermission(c *gin.Context) {
	var req dto.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	permission, err := h.service.CreatePermission(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Permission created successfully", "data": permission})
}

func (h *rbacHandler) UpdatePermission(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	permission, err := h.service.UpdatePermission(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Permission updated successfully", "data": permission})
}

func (h *rbacHandler) DeletePermission(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeletePermission(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Permission deleted successfully"})
}

func (h *rbacHandler) GetModules(c *gin.Context) {
	modules, err := h.service.GetModules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": modules})
}

// -------------------------------------------------------------
// USER ROLES
// -------------------------------------------------------------

func (h *rbacHandler) GetUserRoles(c *gin.Context) {
	userId := c.Param("userId")
	roles, err := h.service.GetUserRoles(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": roles})
}

func (h *rbacHandler) GetUsersByRole(c *gin.Context) {
	roleId := c.Param("roleId")
	users, err := h.service.GetUsersByRole(roleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
}

func (h *rbacHandler) AssignUserRole(c *gin.Context) {
	var req dto.AssignUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	assignedBy := c.GetString("user_id") // Optional from JWT middleware if available

	if err := h.service.AssignUserRole(&req, assignedBy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Role assigned to user successfully"})
}

func (h *rbacHandler) RevokeUserRole(c *gin.Context) {
	userId := c.Param("userId")
	roleId := c.Param("roleId")

	req := dto.RevokeUserRoleRequest{
		UserId: userId,
		RoleId: roleId,
	}

	if err := h.service.RevokeUserRole(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Role revoked from user successfully"})
}

func (h *rbacHandler) GetUserPermissions(c *gin.Context) {
	userId := c.Param("userId")
	permissions, err := h.service.GetUserPermissions(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": permissions})
}

