package rooms

import (
	"citramascoweb-backend/internal/dto"
	"citramascoweb-backend/pkg/utils"
	"time"

	"github.com/google/uuid"
)

type roomService struct {
	roomRepo RoomRepositoryInterface
}

func NewRoomService(roomRepo RoomRepositoryInterface) *roomService {
	return &roomService{roomRepo: roomRepo}
}

func (s *roomService) GetAll() ([]Room, error) {
	rooms, err := s.roomRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *roomService) GetById(id string) (*Room, error) {
	room, err := s.roomRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *roomService) Store(req *dto.CreateRoomRequest) error {
	now := time.Now()

	code, err := utils.GenerateRandomNumberCode(6)
	if err != nil {
		return err
	}

	newRoom := &Room{
		Id:          uuid.New().String(),
		Code:        code,
		Name:        req.Name,
		Slug:        utils.Slugify(req.Name),
		CategoryId:  req.CategoryId,
		Price:       req.Price,
		Capacity:    req.Capacity,
		Size:        req.Size,
		TypeId:      req.TypeId,
		Description: req.Description,
		Image:       req.Image, // Menggunakan string URL publik Supabase yang sudah disiapkan handler
		CreatedAt:   now,
	}

	err = s.roomRepo.Create(newRoom)
	if err != nil {
		return err
	}

	return nil
}

func (s *roomService) Update(id string, req *dto.UpdateRoomRequest) error {
	room, err := s.roomRepo.GetById(id)
	if err != nil {
		return err
	}

	room.Name = req.Name
	room.Slug = utils.Slugify(req.Name)
	room.CategoryId = req.CategoryId
	room.Price = req.Price
	room.Capacity = req.Capacity
	room.Size = req.Size
	room.TypeId = req.TypeId
	room.Description = req.Description

	// PENGAMAN: Hanya ganti nilai gambar jika payload DTO membawa link URL baru dari handler.
	if req.Image != "" {
		room.Image = req.Image
	}

	err = s.roomRepo.Update(room, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *roomService) Delete(id string) error {
	err := s.roomRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *roomService) FilerByCategory(categoryId string) ([]Room, error) {
	rooms, err := s.roomRepo.FilerByCategory(categoryId)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *roomService) FilterByType(typeId string) ([]Room, error) {
	rooms, err := s.roomRepo.FilterByType(typeId)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
