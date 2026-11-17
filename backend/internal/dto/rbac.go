package dto

type CreateRoleRequest struct {
	Name          string   `json:"name" binding:"required"`
	Slug          string   `json:"slug" binding:"required"`
	Description   *string  `json:"description"`
	PermissionIds []string `json:"permission_ids"`
}

type UpdateRoleRequest struct {
	Name          string   `json:"name"`
	Slug          string   `json:"slug"`
	Description   *string  `json:"description"`
	PermissionIds []string `json:"permission_ids"`
}

type SyncRolePermissionsRequest struct {
	PermissionIds []string `json:"permission_ids" binding:"required"`
}

type CreatePermissionRequest struct {
	Name        string  `json:"name" binding:"required"`
	Module      string  `json:"module" binding:"required"`
	Description *string `json:"description"`
}

type UpdatePermissionRequest struct {
	Name        string  `json:"name"`
	Module      string  `json:"module"`
	Description *string `json:"description"`
}

type AssignUserRoleRequest struct {
	UserId     string  `json:"user_id" binding:"required"`
	RoleId     string  `json:"role_id" binding:"required"`
	AssignedBy *string `json:"assigned_by"`
}

type RevokeUserRoleRequest struct {
	UserId string `json:"user_id" binding:"required"`
	RoleId string `json:"role_id" binding:"required"`
}

