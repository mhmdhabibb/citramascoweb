package floor

import (
	"gorm.io/gorm"
)

type FloorRepositoryInterface interface {
	GetAll() ([]Floor, error)
	GetById(id string) (*Floor, error)
	Create(floor *Floor) error
	Update(floor *Floor, id string) error
	Delete(id string) error
}

type floorRepository struct {
	db *gorm.DB
}

func NewFloorRepository(db *gorm.DB) FloorRepositoryInterface {
	return &floorRepository{db: db}
}

func (r *floorRepository) GetAll() ([]Floor, error) {
	var floors []Floor
	err := r.db.Order("floor_number asc").Find(&floors).Error
	if err != nil {
		return nil, err
	}
	return floors, nil
}

func (r *floorRepository) GetById(id string) (*Floor, error) {
	var floor Floor
	err := r.db.Where("id = ?", id).First(&floor).Error
	if err != nil {
		return nil, err
	}
	return &floor, nil
}

func (r *floorRepository) Create(floor *Floor) error {
	return r.db.Create(floor).Error
}

func (r *floorRepository) Update(floor *Floor, id string) error {
	return r.db.Model(&Floor{}).Where("id = ?", id).Updates(floor).Error
}

func (r *floorRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&Floor{}).Error
}
