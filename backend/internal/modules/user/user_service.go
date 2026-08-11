package user

import (
	"citramascoweb-backend/internal/dto"
	"citramascoweb-backend/internal/modules/auth"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo UserRepositoryInterface
}

func NewUserService(userRepository UserRepositoryInterface) *userService {
	return &userService{userRepo: userRepository}
}

func (s *userService) GetAllCustomer() ([]auth.User, error) {

	users, err := s.userRepo.GetAllCustomer()
	if err != nil {
		return nil, err
	}

	return users, nil

}

func (s *userService) GetAllByRole(role string) ([]auth.User, error) {
	users, err := s.userRepo.GetAllByRole(role)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *userService) Delete(id string) error {
	err := s.userRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) CreateUser(request *dto.CreateUserRequest) (*auth.User, error) {
	if len(request.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}
	if strings.Contains(request.Password, " ") {
		return nil, errors.New("password must not contain spaces")
	}

	exists, err := s.userRepo.CheckUsername(request.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username or email is already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &auth.User{
		Id:        uuid.New().String(),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Username:  request.Username,
		Password:  string(hash),
		Phone:     request.Phone,
		Email:     request.Email,
		Address:   request.Address,
		Role:      auth.Role(request.Role),
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateRole(id string, role auth.Role) error {
	err := s.userRepo.UpdateRole(id, role)
	if err != nil {
		return err
	}

	return nil
}
