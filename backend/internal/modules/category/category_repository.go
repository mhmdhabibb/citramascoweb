package category

import (
	"gorm.io/gorm"
)

type CategoryRepositoryInterface interface {
	GetAll() ([]Category, error)
	Store(category *Category) error
	Update(id string, category *Category) error
	Delete(ud string) error
	GetById(id string) (*Category, error)
	GetBySlug(slug string) (*Category, error)
}

type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepositoryInterface {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) GetAll() ([]Category, error) {
	var categories []Category

	err := r.db.Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepo) Store(category *Category) error {

	return r.db.Create(&category).Error

}

func (r *categoryRepo) Update(id string, category *Category) error {

	return r.db.Model(Category{}).Where("id = ?", id).Updates(&category).Error

}

func (r *categoryRepo) Delete(id string) error {

	return r.db.Where("id = ?", id).Delete(&Category{}).Error

}

// func (r *categoryRepo) SoftDelete(id string) error {
// 	return r.db.Where("id = ?", id).Updates(map[string]interface{}{"is_deleted": true, "deleted_at": time.Now()}).Error
// }

func (r *categoryRepo) GetById(id string) (*Category, error) {
	var category *Category

	err := r.db.Where("id = ?", id).First(&category).Error

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *categoryRepo) GetBySlug(slug string) (*Category, error) {
	var category *Category

	err := r.db.Where("slug = ?", slug).First(&category).Error

	if err != nil {
		return nil, err
	}

	return category, nil
}
