package channel

import (
	"citramascoweb-backend/internal/dto"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type channelService struct {
	channelRepo ChannelRepository
}

func NewChannelService(repo ChannelRepository) *channelService {
	return &channelService{channelRepo: repo}
}

func (s *channelService) FindAll() ([]Channel, error) {
	channels, err := s.channelRepo.FindAll()
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (s *channelService) Find(page int) ([]Channel, int, error) {
	channels, total, err := s.channelRepo.Find(page)
	if err != nil {
		return nil, 0, err
	}
	return channels, total, nil
}

func (s *channelService) FindById(id string) (*Channel, error) {
	channel, err := s.channelRepo.FindById(id)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (s *channelService) FindByName(name string) (*Channel, error) {
	channel, err := s.channelRepo.FindByName(name)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (s *channelService) Create(req *dto.CreateChannelRequest) error {
	now := time.Now()

	lastChannel, err := s.channelRepo.GetLastChannel()
	if err != nil {
		return err
	}

	code := "C-1"
	if lastChannel != nil {
		parts := strings.Split(lastChannel.Code, "-")
		if len(parts) > 1 {
			lastCodePart := parts[len(parts)-1]
			lastIndex, err := strconv.Atoi(lastCodePart)
			if err == nil {
				code = fmt.Sprintf("C-%d", lastIndex+1)
			}
		} else {
			code = "C-2"
		}
	}

	newChannel := &Channel{
		Id:        uuid.New().String(),
		Name:      req.Name,
		Code:      code,
		CreatedAt: now,
	}

	err = s.channelRepo.Create(newChannel)
	if err != nil {
		return err
	}

	return nil
}

func (s *channelService) Update(id string, req *dto.UpdateChannelRequest) error {
	channel, err := s.channelRepo.FindById(id)
	if err != nil {
		return err
	}

	channel.Name = req.Name
	channel.UpdatedAt = time.Now()

	err = s.channelRepo.Update(channel, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *channelService) Delete(id string) error {
	err := s.channelRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *channelService) GetLastChannel() (*Channel, error) {
	channel, err := s.channelRepo.GetLastChannel()
	if err != nil {
		return nil, err
	}
	return channel, nil
}
