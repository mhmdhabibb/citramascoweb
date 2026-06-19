package rooms

import (
	"citramascoweb-backend/internal/modules/offer"

	"gorm.io/gorm"
)

type RoomRepositoryInterface interface {
	GetAll() ([]Room, error)
	GetById(id string) (*Room, error)
	Create(room *Room) error
	Update(room *Room, id string) error
	Delete(id string) error
	GetLastRoom() (*Room, error)
	FilerByCategory(categoryId string) ([]Room, error)
	FilterByType(typeId string) ([]Room, error)
	GetOfferByCode(code string) (*offer.Offer, error)
	CountRoomsWithOfferCode(code string, excludeRoomId string) (int64, error)
	DecrementOfferQuota(code string) error
	IncrementOfferQuota(code string) error
}

type roomRepo struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepositoryInterface {
	return &roomRepo{db: db}
}

func (r *roomRepo) GetAll() ([]Room, error) {
	var rooms []Room

	err := r.db.Preload("Category").Preload("Type").Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomRepo) GetById(id string) (*Room, error) {
	var room Room
	err := r.db.Preload("Category").Preload("Type").Where("id = ?", id).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepo) Create(room *Room) error {
	return r.db.Omit("Category", "Type").Create(room).Error
}

func (r *roomRepo) Update(room *Room, id string) error {
	return r.db.Model(&Room{}).Where("id = ?", id).Omit("Category", "Type").Updates(room).Error
}

func (r *roomRepo) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&Room{}).Error
}

func (r *roomRepo) GetLastRoom() (*Room, error) {
	var room Room
	err := r.db.Order("created_at desc").First(&room).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

func (r *roomRepo) FilerByCategory(categoryId string) ([]Room, error) {

	var rooms []Room

	err := r.db.Preload("Type").Preload("Category").Where("category_id = ?", categoryId).Find(&rooms).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rooms, nil

}

func (r *roomRepo) FilterByType(typeId string) ([]Room, error) {

	var rooms []Room

	err := r.db.Preload("Type").Preload("Category").Where("type_id = ?", typeId).Find(&rooms).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rooms, nil

}

func (r *roomRepo) GetOfferByCode(code string) (*offer.Offer, error) {
	var o offer.Offer
	err := r.db.Where("code = ?", code).First(&o).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *roomRepo) CountRoomsWithOfferCode(code string, excludeRoomId string) (int64, error) {
	var count int64
	query := r.db.Model(&Room{}).Where("offer_code = ?", code)
	if excludeRoomId != "" {
		query = query.Where("id != ?", excludeRoomId)
	}
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *roomRepo) DecrementOfferQuota(code string) error {
	return r.db.Model(&offer.Offer{}).Where("code = ? AND max_quota IS NOT NULL AND max_quota > 0", code).UpdateColumn("max_quota", gorm.Expr("max_quota - ?", 1)).Error
}

func (r *roomRepo) IncrementOfferQuota(code string) error {
	return r.db.Model(&offer.Offer{}).Where("code = ? AND max_quota IS NOT NULL", code).UpdateColumn("max_quota", gorm.Expr("max_quota + ?", 1)).Error
}
