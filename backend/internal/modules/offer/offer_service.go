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

	var image string
	if req.Image != nil {
		image = *req.Image
	}

	var discounted *int
	if req.Price != nil {
		priceVal := *req.Price
		discountVal := 0
		if req.Discount != nil {
			discountVal = *req.Discount
		}
		calculated := priceVal - (priceVal * discountVal / 100)
		discounted = &calculated
	}

	validStart, err := dto.ParseCustomDate(req.ValidStart)
	if err != nil {
		return errors.New("Invalid valid_start date format")
	}

	validEnd, err := dto.ParseCustomDate(req.ValidEnd)
	if err != nil {
		return errors.New("Invalid valid_end date format")
	}

	statusVal := OfferStatusDraft
	if req.Status != nil && *req.Status != "" {
		statusVal = OfferStatus(*req.Status)
	}

	var subtitle string
	if req.SubTitle != nil {
		subtitle = *req.SubTitle
	}

	offer := &Offer{
		Id:          uuid.New().String(),
		Title:       req.Title,
		SubTitle:    subtitle,
		Price:       req.Price,
		Discount:    req.Discount,
		Code:        req.Code,
		ValidStart:  validStart,
		ValidEnd:    validEnd,
		MaxQuota:    req.MaxQuota,
		Description: req.Description,
		Discounted:  discounted,
		Image:       image,
		Status:      statusVal,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.repo.Create(offer)

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

	validStart, err := dto.ParseCustomDate(req.ValidStart)
	if err != nil {
		return errors.New("Invalid valid_start date format")
	}

	validEnd, err := dto.ParseCustomDate(req.ValidEnd)
	if err != nil {
		return errors.New("Invalid valid_end date format")
	}

	offer.Title = req.Title
	offer.Discount = req.Discount
	if req.SubTitle != nil {
		offer.SubTitle = *req.SubTitle
	}
	offer.Price = req.Price
	offer.Description = req.Description
	offer.ValidStart = validStart
	offer.ValidEnd = validEnd
	offer.MaxQuota = req.MaxQuota
	offer.Code = req.Code

	if offer.Price != nil {
		priceVal := *offer.Price
		discountVal := 0
		if offer.Discount != nil {
			discountVal = *offer.Discount
		}
		calculated := priceVal - (priceVal * discountVal / 100)
		offer.Discounted = &calculated
	} else {
		offer.Discounted = nil
	}

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
