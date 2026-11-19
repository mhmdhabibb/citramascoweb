package unit

import (
	"gorm.io/gorm"
)

type RoomUnitRepositoryInterface interface {
	GetAll() ([]RoomUnit, error)
	GetById(id string) (*RoomUnit, error)
	GetByRoomNumber(roomNumber string) (*RoomUnit, error)
	GetByRoomId(roomId string) ([]RoomUnit, error)
	GetByFloorId(floorId string) ([]RoomUnit, error)
	GetNextAvailableUnit(roomId string) (*RoomUnit, error)
	Create(unit *RoomUnit) error
	Update(unit *RoomUnit, id string) error
	UpdateStatus(id string, status string) error
	UpdateRoomStatusId(id string, statusId *string) error
	Delete(id string) error

	// RoomStatus operations
	GetAllStatuses() ([]RoomStatus, error)
	GetStatusById(id string) (*RoomStatus, error)
	GetStatusByCode(code string) (*RoomStatus, error)
	CreateStatus(status *RoomStatus) error
	UpdateStatusModel(status *RoomStatus, id string) error
	DeleteStatus(id string) error
}

type roomUnitRepository struct {
	db *gorm.DB
}

func NewRoomUnitRepository(db *gorm.DB) RoomUnitRepositoryInterface {
	return &roomUnitRepository{db: db}
}

func (r *roomUnitRepository) GetAll() ([]RoomUnit, error) {
	var units []RoomUnit
	err := r.db.Preload("Room").Preload("Room.Category").Preload("Room.Type").Preload("Floor").Preload("RoomStatus").Order("room_number asc").Find(&units).Error
	if err != nil {
		return nil, err
	}
	return units, nil
}

func (r *roomUnitRepository) GetById(id string) (*RoomUnit, error) {
	var unit RoomUnit
	err := r.db.Preload("Room").Preload("Room.Category").Preload("Room.Type").Preload("Floor").Preload("RoomStatus").Where("id = ?", id).First(&unit).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *roomUnitRepository) GetByRoomNumber(roomNumber string) (*RoomUnit, error) {
	var unit RoomUnit
	err := r.db.Preload("Room").Preload("Floor").Preload("RoomStatus").Where("room_number = ?", roomNumber).First(&unit).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *roomUnitRepository) GetByRoomId(roomId string) ([]RoomUnit, error) {
	var units []RoomUnit
	err := r.db.Preload("Room").Preload("Floor").Preload("RoomStatus").Where("room_id = ?", roomId).Order("room_number asc").Find(&units).Error
	if err != nil {
		return nil, err
	}
	return units, nil
}

func (r *roomUnitRepository) GetByFloorId(floorId string) ([]RoomUnit, error) {
	var units []RoomUnit
	err := r.db.Preload("Room").Preload("Floor").Preload("RoomStatus").Where("floor_id = ?", floorId).Order("room_number asc").Find(&units).Error
	if err != nil {
		return nil, err
	}
	return units, nil
}

// GetNextAvailableUnit mencari kamar kosong pertama secara berurutan (Lantai terendah & No Kamar terkecil dulu)
func (r *roomUnitRepository) GetNextAvailableUnit(roomId string) (*RoomUnit, error) {
	var unit RoomUnit
	err := r.db.Joins("LEFT JOIN floors ON floors.id = room_units.floor_id").
		Preload("Room").Preload("Floor").Preload("RoomStatus").
		Where("room_units.room_id = ? AND (room_units.status = ? OR room_units.status = ?)", roomId, "available", "VC").
		Order("floors.floor_number ASC, room_units.room_number ASC").
		First(&unit).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *roomUnitRepository) Create(unit *RoomUnit) error {
	return r.db.Create(unit).Error
}

func (r *roomUnitRepository) Update(unit *RoomUnit, id string) error {
	return r.db.Model(&RoomUnit{}).Where("id = ?", id).Updates(unit).Error
}

func (r *roomUnitRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&RoomUnit{}).Where("id = ?", id).Update("status", status).Error
}

func (r *roomUnitRepository) UpdateRoomStatusId(id string, statusId *string) error {
	return r.db.Model(&RoomUnit{}).Where("id = ?", id).Update("room_status_id", statusId).Error
}

func (r *roomUnitRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&RoomUnit{}).Error
}

// ---------------------------------------------------------------------------
// RoomStatus CRUD
// ---------------------------------------------------------------------------

func (r *roomUnitRepository) GetAllStatuses() ([]RoomStatus, error) {
	var statuses []RoomStatus
	err := r.db.Order("code asc").Find(&statuses).Error
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

func (r *roomUnitRepository) GetStatusById(id string) (*RoomStatus, error) {
	var status RoomStatus
	err := r.db.Where("id = ?", id).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *roomUnitRepository) GetStatusByCode(code string) (*RoomStatus, error) {
	var status RoomStatus
	err := r.db.Where("code = ?", code).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *roomUnitRepository) CreateStatus(status *RoomStatus) error {
	return r.db.Create(status).Error
}

func (r *roomUnitRepository) UpdateStatusModel(status *RoomStatus, id string) error {
	return r.db.Model(&RoomStatus{}).Where("id = ?", id).Updates(status).Error
}

func (r *roomUnitRepository) DeleteStatus(id string) error {
	return r.db.Where("id = ?", id).Delete(&RoomStatus{}).Error
}
