package facility

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type FacilityRepository interface {
	FindAllFacilities() ([]Facility, error)
	FindFacilityById(id string) (*Facility, error)
	CreateFacility(entity *Facility) error
	UpdateFacility(entity *Facility) error
	DeleteFacility(id string) error

	// Room Type Facilities
	GetFacilitiesByRoomTypeId(roomTypeId string) ([]Facility, error)
	SyncRoomTypeFacilities(roomTypeId string, facilityIds []string) error

	// Room Type Images
	GetImagesByRoomTypeId(roomTypeId string) ([]RoomTypeImage, error)
	CreateRoomTypeImage(img *RoomTypeImage) error
	DeleteRoomTypeImage(id string) error
	SetPrimaryImage(roomTypeId, imageId string) error
}

type facilityRepository struct {
	db *gorm.DB
}

func NewFacilityRepository(db *gorm.DB) FacilityRepository {
	return &facilityRepository{db: db}
}

func (r *facilityRepository) FindAllFacilities() ([]Facility, error) {
	var list []Facility
	err := r.db.Order("name ASC").Find(&list).Error
	return list, err
}

func (r *facilityRepository) FindFacilityById(id string) (*Facility, error) {
	var item Facility
	err := r.db.Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *facilityRepository) CreateFacility(entity *Facility) error {
	return r.db.Create(entity).Error
}

func (r *facilityRepository) UpdateFacility(entity *Facility) error {
	return r.db.Save(entity).Error
}

func (r *facilityRepository) DeleteFacility(id string) error {
	return r.db.Where("id = ?", id).Delete(&Facility{}).Error
}

func (r *facilityRepository) GetFacilitiesByRoomTypeId(roomTypeId string) ([]Facility, error) {
	var facilities []Facility
	err := r.db.Table("facilities").
		Joins("JOIN room_type_facilities ON room_type_facilities.facility_id = facilities.id").
		Where("room_type_facilities.room_type_id = ? AND facilities.deleted_at IS NULL", roomTypeId).
		Order("facilities.name ASC").
		Find(&facilities).Error
	return facilities, err
}

func (r *facilityRepository) SyncRoomTypeFacilities(roomTypeId string, facilityIds []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("room_type_id = ?", roomTypeId).Delete(&RoomTypeFacility{}).Error; err != nil {
			return err
		}

		for _, fid := range facilityIds {
			rf := RoomTypeFacility{
				RoomTypeId: roomTypeId,
				FacilityId: fid,
				CreatedAt:  time.Now(),
			}
			if err := tx.Create(&rf).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *facilityRepository) GetImagesByRoomTypeId(roomTypeId string) ([]RoomTypeImage, error) {
	var images []RoomTypeImage
	err := r.db.Where("room_type_id = ?", roomTypeId).Order("is_primary DESC, sort_order ASC").Find(&images).Error
	return images, err
}

func (r *facilityRepository) CreateRoomTypeImage(img *RoomTypeImage) error {
	return r.db.Create(img).Error
}

func (r *facilityRepository) DeleteRoomTypeImage(id string) error {
	return r.db.Where("id = ?", id).Delete(&RoomTypeImage{}).Error
}

func (r *facilityRepository) SetPrimaryImage(roomTypeId, imageId string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&RoomTypeImage{}).Where("room_type_id = ?", roomTypeId).Update("is_primary", false).Error; err != nil {
			return err
		}
		res := tx.Model(&RoomTypeImage{}).Where("id = ? AND room_type_id = ?", imageId, roomTypeId).Update("is_primary", true)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("image not found")
		}
		return nil
	})
}
