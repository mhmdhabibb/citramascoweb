package rbac

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	Id          string         `gorm:"type:varchar(255);primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"slug"`
	Description *string        `gorm:"type:text" json:"description"`
	Permissions []Permission   `gorm:"many2many:role_permissions;foreignKey:Id;joinForeignKey:RoleId;references:Id;joinReferences:PermissionId" json:"permissions,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}

type Permission struct {
	Id          string         `gorm:"type:varchar(255);primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Module      string         `gorm:"type:varchar(100);not null;index" json:"module"`
	Description *string        `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Permission) TableName() string {
	return "permissions"
}

type RolePermission struct {
	RoleId       string `gorm:"type:varchar(255);primaryKey" json:"role_id"`
	PermissionId string `gorm:"type:varchar(255);primaryKey" json:"permission_id"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

type UserRole struct {
	UserId     string    `gorm:"type:varchar(255);primaryKey" json:"user_id"`
	RoleId     string    `gorm:"type:varchar(255);primaryKey" json:"role_id"`
	Role       *Role     `gorm:"foreignKey:RoleId;references:Id" json:"role,omitempty"`
	AssignedBy *string   `gorm:"type:varchar(255)" json:"assigned_by,omitempty"`
	AssignedAt time.Time `json:"assigned_at"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

