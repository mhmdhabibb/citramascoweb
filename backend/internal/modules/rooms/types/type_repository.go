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
