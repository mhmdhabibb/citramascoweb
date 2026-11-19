package types

import (
	"gorm.io/gorm"
)

type TypeRepositoryInterface interface {
	GetAll() ([]Types, error)
	GetById(id string) (*Types, error)
	Create(types *Types) error
	Update(id string, types *Types) error
	Delete(id string) error
	Search(query string) ([]Types, error)

	// BedType operations
	GetAllBedTypes() ([]BedType, error)
	GetBedTypeById(id string) (*BedType, error)
	GetBedTypeBySlug(slug string) (*BedType, error)
	CreateBedType(bedType *BedType) error
	UpdateBedType(id string, bedType *BedType) error
	DeleteBedType(id string) error
}

type typeRepository struct {
	db *gorm.DB
}

func NewTypeRepository(db *gorm.DB) TypeRepositoryInterface {
	return &typeRepository{db: db}
}

func (r *typeRepository) GetAll() ([]Types, error) {
	var types []Types

	err := r.db.Preload("Rooms").Find(&types).Error

	if err != nil {
		return nil, err
	}

	return types, nil
}

func (r *typeRepository) GetById(id string) (*Types, error) {
	var types Types

	err := r.db.Preload("Rooms").First(&types, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &types, nil
}

func (r *typeRepository) Create(types *Types) error {
	return r.db.Create(&types).Error
}

func (r *typeRepository) Update(id string, types *Types) error {
	return r.db.Model(Types{}).Where("id = ?", id).Updates(&types).Error
}

func (r *typeRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&Types{}).Error
}

func (r *typeRepository) Search(query string) ([]Types, error) {
	var types []Types
	err := r.db.Preload("Rooms").Where("name LIKE ?", "%"+query+"%").Find(&types).Error
	if err != nil {
		return nil, err
	}
	return types, nil
}

// ---------------------------------------------------------------------------
// BedType CRUD
// ---------------------------------------------------------------------------

func (r *typeRepository) GetAllBedTypes() ([]BedType, error) {
	var bedTypes []BedType
	err := r.db.Order("name asc").Find(&bedTypes).Error
	if err != nil {
		return nil, err
	}
	return bedTypes, nil
}

func (r *typeRepository) GetBedTypeById(id string) (*BedType, error) {
	var bt BedType
	err := r.db.First(&bt, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &bt, nil
}

func (r *typeRepository) GetBedTypeBySlug(slug string) (*BedType, error) {
	var bt BedType
	err := r.db.First(&bt, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &bt, nil
}

func (r *typeRepository) CreateBedType(bedType *BedType) error {
	return r.db.Create(bedType).Error
}

func (r *typeRepository) UpdateBedType(id string, bedType *BedType) error {
	return r.db.Model(&BedType{}).Where("id = ?", id).Updates(bedType).Error
}

func (r *typeRepository) DeleteBedType(id string) error {
	return r.db.Where("id = ?", id).Delete(&BedType{}).Error
}
