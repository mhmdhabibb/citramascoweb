package user

import (
	"citramascoweb-backend/internal/modules/auth"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetAllCustomer() ([]auth.User, error)
	GetAllByRole(role string) ([]auth.User, error)
	Delete(id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &userRepository{db: db}
}

func (r *userRepository) GetAllByRole(role string) ([]auth.User, error) {
	var users []auth.User

	err := r.db.Where("role = ?", role).Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetAllCustomer() ([]auth.User, error) {
	var users []auth.User

	err := r.db.Where("role = ?", "user").Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Delete(id string) error {
	err := r.db.Where("id = ?", id).Delete(&auth.User{}).Error

	if err != nil {
		return err
	}

	return nil
}
