package main

import (
	"time"

	"citramascoweb-backend/config"
	"citramascoweb-backend/internal/modules/audit"
	"citramascoweb-backend/internal/modules/auth"
	"citramascoweb-backend/internal/modules/channel"
	"citramascoweb-backend/internal/modules/dashboard"
	"citramascoweb-backend/internal/modules/facility"
	"citramascoweb-backend/internal/modules/finance"
	"citramascoweb-backend/internal/modules/guest"
	"citramascoweb-backend/internal/modules/inventory"
	"citramascoweb-backend/internal/modules/notification"
	"citramascoweb-backend/internal/modules/offer"
	"citramascoweb-backend/internal/modules/pricing"
	"citramascoweb-backend/internal/modules/rbac"
	"citramascoweb-backend/internal/modules/rooms"
	"citramascoweb-backend/internal/modules/rooms/category"
	"citramascoweb-backend/internal/modules/rooms/floor"
	"citramascoweb-backend/internal/modules/rooms/reservation"
	"citramascoweb-backend/internal/modules/rooms/types"
	"citramascoweb-backend/internal/modules/rooms/unit"
	"citramascoweb-backend/internal/modules/service_request"
	"citramascoweb-backend/internal/modules/sla"
	"citramascoweb-backend/internal/modules/smartlock"
	"citramascoweb-backend/pkg/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	app := gin.Default()

	db := config.ConnectDB()

	jwtSecret := config.Config("JWT_SECRET")
	if jwtSecret != "" {
		utils.SetJWTSecret([]byte(jwtSecret))
	}

	// Auto migrate all entities according to target ERD
	db.AutoMigrate(
		&rbac.Role{}, &rbac.Permission{}, &rbac.RolePermission{}, &rbac.UserRole{},
		&audit.StaffAuditLog{},
		&category.Category{},
		&types.Types{}, &types.BedType{},
		&facility.Facility{}, &facility.RoomTypeFacility{}, &facility.RoomTypeImage{},
		&floor.Floor{},
		&rooms.Room{},
		&unit.RoomStatus{}, &unit.RoomUnit{},
		&guest.IdentityType{}, &guest.Guest{},
		&channel.Channel{},
		&offer.Offer{}, &offer.OfferUsage{},
		&reservation.Reservation{}, &reservation.ReservationItem{}, &reservation.ReservationLog{},
		&smartlock.RoomAccessCredential{}, &smartlock.RoomAccessLog{},
		&sla.SLAPolicy{}, &sla.OperationalTask{}, &sla.StaffSLAPerformanceDaily{},
		&pricing.CalendarSeason{}, &pricing.DailyRoomRate{}, &pricing.DailyMarketSnapshot{}, &pricing.CompetitorRate{}, &pricing.MarketEvent{}, &pricing.MonthlyKPISummary{},
		&inventory.InventoryItem{}, &inventory.InventoryTransaction{}, &inventory.InventoryStockTake{},
		&service_request.ServiceRequest{},
	)

	// Seed default master data
	rbac.SeedDefaultRoles(db)
	guest.SeedDefaultIdentityTypes(db)
	facility.SeedDefaultFacilities(db)
	sla.SeedDefaultSLAPolicies(db)
	pricing.SeedDefaultSeasons(db)
	unit.SeedDefaultRoomStatuses(db)
	types.SeedDefaultBedTypes(db)

	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}

	app.Use(cors.New(corsConfig))

	port := config.Config("PORT")

	//  init modules
	categoryModule := category.InitModule(db)
	typeModule := types.InitModule(db)
	floorModule := floor.InitModule(db)
	unitModule := unit.InitModule(db)
	roomModule := rooms.InitModule(db)
	reservationModule := reservation.InitModule(db)
	offerModule := offer.InitModule(db)
	dashboardModule := dashboard.InitModule(db)
	inventoryModule := inventory.InitModule(db)
	notificationModule := notification.InitModule(db)
	financeModule := finance.InitModule(db)
	channelModule := channel.InitModule(db)
	serviceRequestModule := service_request.InitModule(db)
	rbacModule := rbac.InitModule(db)
	authModule := auth.InitModule(db)

	// New modules from ERD
	guestModule := guest.NewGuestModule(db)
	facilityModule := facility.NewFacilityModule(db)
	smartLockModule := smartlock.NewSmartLockModule(db)
	slaModule := sla.NewSLAModule(db)
	pricingModule := pricing.NewPricingModule(db)
	auditModule := audit.NewAuditModule(db)

	api := app.Group("/api")

	authModule.AuthRoutes(api)
	categoryModule.CategoryRoutes(api)
	typeModule.TypeRoutes(api)
	floorModule.FloorRoutes(api)
	unitModule.RoomUnitRoutes(api)
	roomModule.RoomRoutes(api)
	reservationModule.ReservationRoutes(api)
	offerModule.OfferRoutes(api)
	dashboardModule.DashboardRoutes(api)
	inventoryModule.InventoryRoutes(api)
	notificationModule.NotificationRoutes(api)
	financeModule.FinanceRoutes(api)
	channelModule.ChannelRoutes(api)
	serviceRequestModule.ServiceRequestRoutes(api)
	rbacModule.RbacRoutes(api)

	// Register routes for new modules
	guestModule.RegisterRoutes(api)
	facilityModule.RegisterRoutes(api)
	smartLockModule.RegisterRoutes(api)
	slaModule.RegisterRoutes(api)
	pricingModule.RegisterRoutes(api)
	auditModule.RegisterRoutes(api)

	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	app.Run(":" + port)
}
