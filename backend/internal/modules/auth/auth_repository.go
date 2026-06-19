package auth

import "gorm.io/gorm"

type AuthRepositoryInterface interface {
	GetById(id string) (*User, error)

	GetByUsername(username string) (*User, error)
	Create(user *User) error
	Update(id string, user *User) error
	Delete(id string) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepositoryInterface {
	return &authRepository{db: db}
}

func (r *authRepository) GetById(id string) (*User, error) {
	var user User

	err := r.db.First(&user, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepository) GetByUsername(username string) (*User, error) {
	var user User

	err := r.db.First(&user, "username = ? OR email = ?", username, username).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepository) Create(user *User) error {

	return r.db.Create(&user).Error
}

func (r *authRepository) Update(id string, user *User) error {
	return r.db.Model(&User{}).Where("id = ?", id).Updates(&user).Error
}

func (r *authRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&User{}).Error
}
