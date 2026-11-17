package rbac

import (
	"citramascoweb-backend/internal/dto"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RbacService interface {
	// Roles
	GetAllRoles() ([]Role, error)
	GetRoleById(id string) (*Role, error)
	CreateRole(req *dto.CreateRoleRequest) (*Role, error)
	UpdateRole(id string, req *dto.UpdateRoleRequest) (*Role, error)
	DeleteRole(id string) error
	SyncRolePermissions(roleId string, req *dto.SyncRolePermissionsRequest) error

	// Permissions
	GetAllPermissions(module string) ([]Permission, error)
	GetPermissionById(id string) (*Permission, error)
	CreatePermission(req *dto.CreatePermissionRequest) (*Permission, error)
	UpdatePermission(id string, req *dto.UpdatePermissionRequest) (*Permission, error)
	DeletePermission(id string) error
	GetModules() ([]string, error)

	// User Roles
	GetUserRoles(userId string) ([]UserRole, error)
	GetUsersByRole(roleId string) ([]UserRole, error)
	AssignUserRole(req *dto.AssignUserRoleRequest, assignedBy string) error
	RevokeUserRole(req *dto.RevokeUserRoleRequest) error
	GetUserPermissions(userId string) ([]Permission, error)
}

type rbacService struct {
	repo RbacRepository
}

func NewRbacService(repo RbacRepository) RbacService {
	return &rbacService{repo: repo}
}

// -------------------------------------------------------------
// ROLES
// -------------------------------------------------------------

func (s *rbacService) GetAllRoles() ([]Role, error) {
	return s.repo.FindAllRoles()
}

func (s *rbacService) GetRoleById(id string) (*Role, error) {
	role, err := s.repo.FindRoleById(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}
	return role, nil
}

func (s *rbacService) CreateRole(req *dto.CreateRoleRequest) (*Role, error) {
	slug := strings.TrimSpace(strings.ToLower(req.Slug))
	existing, err := s.repo.FindRoleBySlug(slug)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("role with this slug already exists")
	}

	role := &Role{
		Id:          uuid.New().String(),
		Name:        strings.TrimSpace(req.Name),
		Slug:        slug,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateRole(role); err != nil {
		return nil, err
	}

	if len(req.PermissionIds) > 0 {
		if err := s.repo.SyncRolePermissions(role.Id, req.PermissionIds); err != nil {
			return nil, err
		}
	}

	return s.repo.FindRoleById(role.Id)
}

func (s *rbacService) UpdateRole(id string, req *dto.UpdateRoleRequest) (*Role, error) {
	role, err := s.repo.FindRoleById(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	updated := &Role{
		Name:        strings.TrimSpace(req.Name),
		Description: req.Description,
	}

	if req.Slug != "" {
		newSlug := strings.TrimSpace(strings.ToLower(req.Slug))
		if newSlug != role.Slug {
			existing, err := s.repo.FindRoleBySlug(newSlug)
			if err != nil {
				return nil, err
			}
			if existing != nil && existing.Id != id {
				return nil, errors.New("role with this slug already exists")
			}
			updated.Slug = newSlug
		}
	}

	if err := s.repo.UpdateRole(id, updated); err != nil {
		return nil, err
	}

	if req.PermissionIds != nil {
		if err := s.repo.SyncRolePermissions(id, req.PermissionIds); err != nil {
			return nil, err
		}
	}

	return s.repo.FindRoleById(id)
}

func (s *rbacService) DeleteRole(id string) error {
	role, err := s.repo.FindRoleById(id)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}
	return s.repo.DeleteRole(id)
}

func (s *rbacService) SyncRolePermissions(roleId string, req *dto.SyncRolePermissionsRequest) error {
	role, err := s.repo.FindRoleById(roleId)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}
	return s.repo.SyncRolePermissions(roleId, req.PermissionIds)
}

// -------------------------------------------------------------
// PERMISSIONS
// -------------------------------------------------------------

func (s *rbacService) GetAllPermissions(module string) ([]Permission, error) {
	return s.repo.FindAllPermissions(module)
}

func (s *rbacService) GetPermissionById(id string) (*Permission, error) {
	p, err := s.repo.FindPermissionById(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("permission not found")
	}
	return p, nil
}

func (s *rbacService) CreatePermission(req *dto.CreatePermissionRequest) (*Permission, error) {
	permission := &Permission{
		Id:          uuid.New().String(),
		Name:        strings.TrimSpace(req.Name),
		Module:      strings.TrimSpace(strings.ToLower(req.Module)),
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreatePermission(permission); err != nil {
		return nil, err
	}
	return permission, nil
}

func (s *rbacService) UpdatePermission(id string, req *dto.UpdatePermissionRequest) (*Permission, error) {
	p, err := s.repo.FindPermissionById(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("permission not found")
	}

	updated := &Permission{
		Name:        strings.TrimSpace(req.Name),
		Module:      strings.TrimSpace(strings.ToLower(req.Module)),
		Description: req.Description,
	}

	if err := s.repo.UpdatePermission(id, updated); err != nil {
		return nil, err
	}
	return s.repo.FindPermissionById(id)
}

func (s *rbacService) DeletePermission(id string) error {
	p, err := s.repo.FindPermissionById(id)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("permission not found")
	}
	return s.repo.DeletePermission(id)
}

func (s *rbacService) GetModules() ([]string, error) {
	return s.repo.GetDistinctModules()
}

// -------------------------------------------------------------
// USER ROLES
// -------------------------------------------------------------

func (s *rbacService) GetUserRoles(userId string) ([]UserRole, error) {
	return s.repo.FindUserRolesByUserId(userId)
}

func (s *rbacService) GetUsersByRole(roleId string) ([]UserRole, error) {
	return s.repo.FindUsersByRoleId(roleId)
}

func (s *rbacService) AssignUserRole(req *dto.AssignUserRoleRequest, assignedBy string) error {
	role, err := s.repo.FindRoleById(req.RoleId)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	assigner := assignedBy
	if req.AssignedBy != nil && *req.AssignedBy != "" {
		assigner = *req.AssignedBy
	}

	var assignerPtr *string
	if assigner != "" {
		assignerPtr = &assigner
	}

	userRole := &UserRole{
		UserId:     req.UserId,
		RoleId:     req.RoleId,
		AssignedBy: assignerPtr,
		AssignedAt: time.Now(),
	}

	return s.repo.AssignUserRole(userRole)
}

func (s *rbacService) RevokeUserRole(req *dto.RevokeUserRoleRequest) error {
	return s.repo.RevokeUserRole(req.UserId, req.RoleId)
}

func (s *rbacService) GetUserPermissions(userId string) ([]Permission, error) {
	return s.repo.GetUserPermissions(userId)
}

