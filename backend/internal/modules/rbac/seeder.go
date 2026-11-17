package rbac

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedDefaultRoles(db *gorm.DB) {
	// 1. Seed Roles
	roles := []Role{
		{
			Name:        "Administrator",
			Slug:        "admin",
			Description: strPtr("Full system access and configurations"),
		},
		{
			Name:        "Hotel Manager",
			Slug:        "manager",
			Description: strPtr("Management, reports, and approvals"),
		},
		{
			Name:        "Receptionist",
			Slug:        "reception",
			Description: strPtr("Front desk, check-in, check-out, and guest bookings"),
		},
		{
			Name:        "Finance",
			Slug:        "finance",
			Description: strPtr("Billing, transactions, deposits, and financial reports"),
		},
		{
			Name:        "Housekeeping",
			Slug:        "housekeeping",
			Description: strPtr("Room cleaning, maintenance status, and service requests"),
		},
		{
			Name:        "Inventory",
			Slug:        "inventory",
			Description: strPtr("Item stock, transactions, and inventory takes"),
		},
	}

	roleMap := make(map[string]*Role)
	for _, r := range roles {
		var existing Role
		if err := db.Where("slug = ?", r.Slug).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				r.Id = uuid.New().String()
				r.CreatedAt = time.Now()
				r.UpdatedAt = time.Now()
				if err := db.Create(&r).Error; err != nil {
					log.Printf("Failed to seed role %s: %v", r.Slug, err)
				} else {
					roleMap[r.Slug] = &r
				}
			}
		} else {
			roleMap[r.Slug] = &existing
		}
	}

	// 2. Seed Master Permissions for Citramasco Modules
	defaultPermissions := []Permission{
		// Rooms
		{Name: "rooms.view", Module: "rooms", Description: strPtr("View room types, units, and status")},
		{Name: "rooms.create", Module: "rooms", Description: strPtr("Create room categories, types, and units")},
		{Name: "rooms.update", Module: "rooms", Description: strPtr("Update room details and pricing")},
		{Name: "rooms.delete", Module: "rooms", Description: strPtr("Delete rooms or units")},

		// Reservations
		{Name: "reservations.view", Module: "reservations", Description: strPtr("View guest reservations and history")},
		{Name: "reservations.create", Module: "reservations", Description: strPtr("Create new booking / reservation")},
		{Name: "reservations.checkin", Module: "reservations", Description: strPtr("Perform guest check-in")},
		{Name: "reservations.checkout", Module: "reservations", Description: strPtr("Perform guest check-out")},
		{Name: "reservations.cancel", Module: "reservations", Description: strPtr("Cancel or reject reservation")},

		// Service Requests & Housekeeping
		{Name: "service_requests.view", Module: "housekeeping", Description: strPtr("View room service requests and cleaning status")},
		{Name: "service_requests.create", Module: "housekeeping", Description: strPtr("Create room service request or report damage")},
		{Name: "service_requests.assign", Module: "housekeeping", Description: strPtr("Assign cleaning task to staff")},
		{Name: "service_requests.complete", Module: "housekeeping", Description: strPtr("Mark cleaning task completed")},

		// Finance
		{Name: "finance.view", Module: "finance", Description: strPtr("View finance dashboard and chart of accounts")},
		{Name: "finance.transaction_create", Module: "finance", Description: strPtr("Record cash/bank transaction and deposits")},
		{Name: "finance.reports", Module: "finance", Description: strPtr("View income statement and balance sheet")},

		// Inventory
		{Name: "inventory.view", Module: "inventory", Description: strPtr("View stock items and catalog")},
		{Name: "inventory.transaction", Module: "inventory", Description: strPtr("Create stock-in or stock-out transaction")},
		{Name: "inventory.stock_take", Module: "inventory", Description: strPtr("Conduct physical stock take")},

		// RBAC
		{Name: "roles.manage", Module: "rbac", Description: strPtr("Create, update, and delete roles")},
		{Name: "permissions.manage", Module: "rbac", Description: strPtr("Manage system permissions")},
		{Name: "user_roles.manage", Module: "rbac", Description: strPtr("Assign roles to users")},

		// Guests & Identity
		{Name: "guests.view", Module: "guests", Description: strPtr("View guest profiles and identification documents")},
		{Name: "guests.manage", Module: "guests", Description: strPtr("Create, update, or remove guest profiles")},

		// Facilities & Amenities
		{Name: "facilities.view", Module: "facilities", Description: strPtr("View facilities and amenities list")},
		{Name: "facilities.manage", Module: "facilities", Description: strPtr("Manage room type facilities and gallery images")},

		// Smart Lock (IoT Access)
		{Name: "smartlock.view", Module: "smartlock", Description: strPtr("View door lock credentials and unlock history")},
		{Name: "smartlock.manage", Module: "smartlock", Description: strPtr("Generate PINs, assign RFID cards, and revoke keys")},

		// SLA & Operations
		{Name: "sla.view", Module: "sla", Description: strPtr("View operational tasks and SLA compliance reports")},
		{Name: "sla.manage", Module: "sla", Description: strPtr("Configure SLA response/resolution policies and assign tasks")},

		// Revenue Management & Pricing
		{Name: "pricing.view", Module: "pricing", Description: strPtr("View dynamic rates, seasons, and market snapshots")},
		{Name: "pricing.manage", Module: "pricing", Description: strPtr("Set daily pricing, manage seasons, and record market events")},

		// Audit Log
		{Name: "audit.view", Module: "audit", Description: strPtr("View staff audit logs and system activity trails")},
	}

	var allPermIds []string
	permIdByName := make(map[string]string)

	for _, p := range defaultPermissions {
		var existing Permission
		if err := db.Where("name = ?", p.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				p.Id = uuid.New().String()
				p.CreatedAt = time.Now()
				p.UpdatedAt = time.Now()
				if err := db.Create(&p).Error; err == nil {
					allPermIds = append(allPermIds, p.Id)
					permIdByName[p.Name] = p.Id
				}
			}
		} else {
			allPermIds = append(allPermIds, existing.Id)
			permIdByName[p.Name] = existing.Id
		}
	}

	// 3. Link Permissions to Roles (Default mapping)
	// Admin gets ALL permissions
	if adminRole, ok := roleMap["admin"]; ok {
		syncRolePermissionsIfEmpty(db, adminRole.Id, allPermIds)
	}

	// Manager gets operational + reports
	if managerRole, ok := roleMap["manager"]; ok {
		var managerPerms []string
		for name, id := range permIdByName {
			if name != "roles.manage" && name != "permissions.manage" {
				managerPerms = append(managerPerms, id)
			}
		}
		syncRolePermissionsIfEmpty(db, managerRole.Id, managerPerms)
	}

	// Receptionist gets reservations & room view
	if recRole, ok := roleMap["reception"]; ok {
		recPermNames := []string{
			"rooms.view", "reservations.view", "reservations.create",
			"reservations.checkin", "reservations.checkout", "service_requests.create", "service_requests.view",
		}
		var recPerms []string
		for _, name := range recPermNames {
			if id, ok := permIdByName[name]; ok {
				recPerms = append(recPerms, id)
			}
		}
		syncRolePermissionsIfEmpty(db, recRole.Id, recPerms)
	}

	// Housekeeping gets cleaning tasks & room status
	if hkRole, ok := roleMap["housekeeping"]; ok {
		hkPermNames := []string{
			"rooms.view", "service_requests.view", "service_requests.complete", "service_requests.create",
		}
		var hkPerms []string
		for _, name := range hkPermNames {
			if id, ok := permIdByName[name]; ok {
				hkPerms = append(hkPerms, id)
			}
		}
		syncRolePermissionsIfEmpty(db, hkRole.Id, hkPerms)
	}
}

func syncRolePermissionsIfEmpty(db *gorm.DB, roleId string, permissionIds []string) {
	var count int64
	db.Model(&RolePermission{}).Where("role_id = ?", roleId).Count(&count)
	if count == 0 && len(permissionIds) > 0 {
		var items []RolePermission
		for _, pid := range permissionIds {
			items = append(items, RolePermission{
				RoleId:       roleId,
				PermissionId: pid,
			})
		}
		db.Create(&items)
	}
}

func strPtr(s string) *string {
	return &s
}
