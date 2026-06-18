package offer

import (
	"citramascoweb-backend/internal/dto"
	"errors"
	"time"

	"github.com/google/uuid"
)

type offerService struct {
	repo OfferRepositoryInterface
}

func NewOfferService(repo OfferRepositoryInterface) *offerService {
	return &offerService{repo: repo}
}

func (s *offerService) GetAll() ([]Offer, error) {

	offers, err := s.repo.GetAll()

	if err != nil {
		return nil, err
	}

	return offers, nil

}
func (s *offerService) GetAlls() ([]Offer, error) {

	offers, err := s.repo.GetAlls()

	if err != nil {
		return nil, err
	}

	return offers, nil

}

func (s *offerService) GetById(id string) (*Offer, error) {
	offer, err := s.repo.GetById(id)
	if err != nil {
		return nil, errors.New("Offer Not Found!")
	}

	return offer, nil
}

func (s *offerService) Create(req *dto.CreateOfferRequest) error {

	if req.Title == "" {
		return errors.New("Offer is required!")
	}

	offer := &Offer{
		Id:          uuid.New().String(),
		Title:       req.Title,
		Price:       req.Price,
		Discount:    req.Discount,
		Description: req.Description,
		Image:       *req.Image,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repo.Create(offer)

	if err != nil {
		return errors.New("Error when inserting the offer data")
	}

	return nil

}

func (s *offerService) Update(id string, req *dto.UpdateOfferRequest) error {

	offer, err := s.repo.GetById(id)
	if err != nil {
		return errors.New("Offer Not Found!")
	}

	offer.Title = req.Title
	offer.Discount = req.Discount
	offer.Price = req.Price
	offer.Description = req.Description
	if req.Image != nil {
		offer.Image = *req.Image
	}
	if req.Status != nil {
		offer.Status = OfferStatus(*req.Status)
	}
	offer.UpdatedAt = time.Now()

	return s.repo.Update(id, offer)
}

func (s *offerService) Delete(id string) error {

	_, err := s.repo.GetById(id)
	if err != nil {
		return errors.New("Offer Not Found!")
	}

	err = s.repo.Delete(id)
	if err != nil {
		return errors.New("Error when deleting the offer data")
	}

	return nil
}
