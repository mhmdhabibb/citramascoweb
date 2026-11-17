package rbac

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type RbacRepository interface {
	// Roles
	FindAllRoles() ([]Role, error)
	FindRoleById(id string) (*Role, error)
	FindRoleBySlug(slug string) (*Role, error)
	CreateRole(role *Role) error
	UpdateRole(id string, role *Role) error
	DeleteRole(id string) error
	SyncRolePermissions(roleId string, permissionIds []string) error

	// Permissions
	FindAllPermissions(module string) ([]Permission, error)
	FindPermissionById(id string) (*Permission, error)
	CreatePermission(permission *Permission) error
	UpdatePermission(id string, permission *Permission) error
	DeletePermission(id string) error
	GetDistinctModules() ([]string, error)

	// User Roles
	FindUserRolesByUserId(userId string) ([]UserRole, error)
	FindUsersByRoleId(roleId string) ([]UserRole, error)
	AssignUserRole(userRole *UserRole) error
	RevokeUserRole(userId string, roleId string) error
	GetUserPermissions(userId string) ([]Permission, error)
}

type rbacRepo struct {
	db *gorm.DB
}

func NewRbacRepository(db *gorm.DB) RbacRepository {
	return &rbacRepo{db: db}
}

// -------------------------------------------------------------
// ROLES
// -------------------------------------------------------------

func (r *rbacRepo) FindAllRoles() ([]Role, error) {
	var roles []Role
	err := r.db.Preload("Permissions").Order("created_at asc").Find(&roles).Error
	return roles, err
}

func (r *rbacRepo) FindRoleById(id string) (*Role, error) {
	var role Role
	err := r.db.Preload("Permissions").First(&role, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *rbacRepo) FindRoleBySlug(slug string) (*Role, error) {
	var role Role
	err := r.db.Preload("Permissions").First(&role, "slug = ?", slug).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *rbacRepo) CreateRole(role *Role) error {
	return r.db.Create(role).Error
}

func (r *rbacRepo) UpdateRole(id string, role *Role) error {
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}
	if role.Name != "" {
		updates["name"] = role.Name
	}
	if role.Slug != "" {
		updates["slug"] = role.Slug
	}
	if role.Description != nil {
		updates["description"] = role.Description
	}

	return r.db.Model(&Role{}).Where("id = ?", id).Updates(updates).Error
}

func (r *rbacRepo) DeleteRole(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Clean pivot records
		if err := tx.Where("role_id = ?", id).Delete(&RolePermission{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", id).Delete(&UserRole{}).Error; err != nil {
			return err
		}
		return tx.Delete(&Role{}, "id = ?", id).Error
	})
}

func (r *rbacRepo) SyncRolePermissions(roleId string, permissionIds []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing associations
		if err := tx.Where("role_id = ?", roleId).Delete(&RolePermission{}).Error; err != nil {
			return err
		}

		if len(permissionIds) == 0 {
			return nil
		}

		var items []RolePermission
		for _, pid := range permissionIds {
			items = append(items, RolePermission{
				RoleId:       roleId,
				PermissionId: pid,
			})
		}
		return tx.Create(&items).Error
	})
}

// -------------------------------------------------------------
// PERMISSIONS
// -------------------------------------------------------------

func (r *rbacRepo) FindAllPermissions(module string) ([]Permission, error) {
	var permissions []Permission
	query := r.db.Order("module asc, name asc")
	if module != "" && module != "all" {
		query = query.Where("module = ?", module)
	}
	err := query.Find(&permissions).Error
	return permissions, err
}

func (r *rbacRepo) FindPermissionById(id string) (*Permission, error) {
	var permission Permission
	err := r.db.First(&permission, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

func (r *rbacRepo) CreatePermission(permission *Permission) error {
	return r.db.Create(permission).Error
}

func (r *rbacRepo) UpdatePermission(id string, permission *Permission) error {
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}
	if permission.Name != "" {
		updates["name"] = permission.Name
	}
	if permission.Module != "" {
		updates["module"] = permission.Module
	}
	if permission.Description != nil {
		updates["description"] = permission.Description
	}

	return r.db.Model(&Permission{}).Where("id = ?", id).Updates(updates).Error
}

func (r *rbacRepo) DeletePermission(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("permission_id = ?", id).Delete(&RolePermission{}).Error; err != nil {
			return err
		}
		return tx.Delete(&Permission{}, "id = ?", id).Error
	})
}

func (r *rbacRepo) GetDistinctModules() ([]string, error) {
	var modules []string
	err := r.db.Model(&Permission{}).Distinct().Order("module asc").Pluck("module", &modules).Error
	return modules, err
}

// -------------------------------------------------------------
// USER ROLES
// -------------------------------------------------------------

func (r *rbacRepo) FindUserRolesByUserId(userId string) ([]UserRole, error) {
	var userRoles []UserRole
	err := r.db.Preload("Role").Preload("Role.Permissions").Where("user_id = ?", userId).Find(&userRoles).Error
	return userRoles, err
}

func (r *rbacRepo) FindUsersByRoleId(roleId string) ([]UserRole, error) {
	var userRoles []UserRole
	err := r.db.Preload("Role").Where("role_id = ?", roleId).Find(&userRoles).Error
	return userRoles, err
}

func (r *rbacRepo) AssignUserRole(userRole *UserRole) error {
	var existing UserRole
	err := r.db.Where("user_id = ? AND role_id = ?", userRole.UserId, userRole.RoleId).First(&existing).Error
	if err == nil {
		// Already assigned, update assigned_by / assigned_at
		return r.db.Model(&UserRole{}).
			Where("user_id = ? AND role_id = ?", userRole.UserId, userRole.RoleId).
			Updates(map[string]interface{}{
				"assigned_by": userRole.AssignedBy,
				"assigned_at": userRole.AssignedAt,
			}).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return r.db.Create(userRole).Error
}

func (r *rbacRepo) RevokeUserRole(userId string, roleId string) error {
	return r.db.Where("user_id = ? AND role_id = ?", userId, roleId).Delete(&UserRole{}).Error
}

func (r *rbacRepo) GetUserPermissions(userId string) ([]Permission, error) {
	var permissions []Permission
	err := r.db.Distinct().
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.deleted_at IS NULL", userId).
		Order("permissions.module asc, permissions.name asc").
		Find(&permissions).Error
	return permissions, err
}

