package auth

import (
	"citramascoweb-backend/internal/dto"
	"citramascoweb-backend/pkg/utils"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	authRepo AuthRepositoryInterface
}

func NewAuthService(authRepo AuthRepositoryInterface) *authService {
	return &authService{authRepo: authRepo}
}

func (s *authService) Register(request *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// Validate password length
	if len(request.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}
	// Validate password has no blank spaces
	if strings.Contains(request.Password, " ") {
		return nil, errors.New("password must not contain spaces")
	}

	// Check if username already exists
	existingUser, err := s.authRepo.GetByUsername(request.Username)
	if err == nil && existingUser != nil {
		return nil, errors.New("Username is already taken!")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &User{
		Id:        uuid.New().String(),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Username:  request.Username,
		Password:  string(hash),
		Phone:     request.Phone,
		Email:     request.Email,
		Address:   request.Address,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.authRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Phone:     user.Phone,
		Email:     user.Email,
		Address:   user.Address,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *authService) Login(request *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.authRepo.GetByUsername(request.Username)
	if err != nil {
		return nil, errors.New("Invalid Credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(request.Password),
	)

	if err != nil {
		return nil, errors.New("Invalid Credentials")
	}

	token, err := utils.GenerateToken(user.Id, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken: token,
	}, nil
}

func (s *authService) Profile(id string) (*User, error) {
	user, err := s.authRepo.GetById(id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) UpdateProfile(id string, request *dto.UpdateProfileRequest) error {
	user, err := s.authRepo.GetById(id)
	if err != nil {
		return err
	}

	if request.FirstName != "" {
		user.FirstName = request.FirstName
	}
	if request.LastName != "" {
		user.LastName = request.LastName
	}
	if request.Phone != "" {
		user.Phone = request.Phone
	}
	if request.Email != "" {
		user.Email = request.Email
	}
	if request.Address != "" {
		user.Address = request.Address
	}
	user.UpdatedAt = time.Now()

	err = s.authRepo.Update(id, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) Logout() error {
	// Stateless logout is handled client-side by clearing the token.
	return nil
}
