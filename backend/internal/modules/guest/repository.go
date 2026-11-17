package guest

import (
	"errors"

	"gorm.io/gorm"
)

type GuestRepository interface {
	// Identity Types
	FindIdentityTypes() ([]IdentityType, error)
	FindIdentityTypeById(id string) (*IdentityType, error)
	CreateIdentityType(entity *IdentityType) error
	UpdateIdentityType(entity *IdentityType) error
	DeleteIdentityType(id string) error

	// Guests
	FindGuests(search string, limit, offset int) ([]Guest, int64, error)
	FindGuestById(id string) (*Guest, error)
	CreateGuest(entity *Guest) error
	UpdateGuest(entity *Guest) error
	DeleteGuest(id string) error
}

type guestRepository struct {
	db *gorm.DB
}

func NewGuestRepository(db *gorm.DB) GuestRepository {
	return &guestRepository{db: db}
}

func (r *guestRepository) FindIdentityTypes() ([]IdentityType, error) {
	var list []IdentityType
	err := r.db.Order("name ASC").Find(&list).Error
	return list, err
}

func (r *guestRepository) FindIdentityTypeById(id string) (*IdentityType, error) {
	var item IdentityType
	err := r.db.Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *guestRepository) CreateIdentityType(entity *IdentityType) error {
	return r.db.Create(entity).Error
}

func (r *guestRepository) UpdateIdentityType(entity *IdentityType) error {
	return r.db.Save(entity).Error
}

func (r *guestRepository) DeleteIdentityType(id string) error {
	return r.db.Where("id = ?", id).Delete(&IdentityType{}).Error
}

func (r *guestRepository) FindGuests(search string, limit, offset int) ([]Guest, int64, error) {
	var guests []Guest
	var total int64

	q := r.db.Model(&Guest{})
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR phone ILIKE ? OR identity_number ILIKE ?", like, like, like, like, like)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit <= 0 {
		limit = 20
	}

	err := q.Preload("IdentityType").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&guests).Error

	return guests, total, err
}

func (r *guestRepository) FindGuestById(id string) (*Guest, error) {
	var guest Guest
	err := r.db.Preload("IdentityType").Where("id = ?", id).First(&guest).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("guest not found")
		}
		return nil, err
	}
	return &guest, nil
}

func (r *guestRepository) CreateGuest(entity *Guest) error {
	return r.db.Create(entity).Error
}

func (r *guestRepository) UpdateGuest(entity *Guest) error {
	return r.db.Save(entity).Error
}

func (r *guestRepository) DeleteGuest(id string) error {
	return r.db.Where("id = ?", id).Delete(&Guest{}).Error
}
