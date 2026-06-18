package user

import "citramascoweb-backend/internal/modules/auth"

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
