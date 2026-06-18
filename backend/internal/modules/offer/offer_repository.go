package offer

import "gorm.io/gorm"

type OfferRepositoryInterface interface {
	GetAll() ([]Offer, error)
	GetAlls() ([]Offer, error)
	GetById(id string) (*Offer, error)
	Create(offer *Offer) error
	Update(id string, offer *Offer) error
	Delete(id string) error
}

type offerRepository struct {
	db *gorm.DB
}

func NewOfferRepository(db *gorm.DB) OfferRepositoryInterface {
	return &offerRepository{db: db}
}

func (r *offerRepository) GetAll() ([]Offer, error) {
	var offers []Offer

	err := r.db.Where("status = ?", "active").Find(&offers).Error
	if err != nil {
		return nil, err
	}

	return offers, nil
}

func (r *offerRepository) GetAlls() ([]Offer, error) {
	var offers []Offer

	err := r.db.Find(&offers).Error
	if err != nil {
		return nil, err
	}

	return offers, nil
}

func (r *offerRepository) GetById(id string) (*Offer, error) {
	var offer Offer
	err := r.db.Where("id = ?", id).First(&offer).Error

	if err != nil {
		return nil, err
	}

	return &offer, nil
}

func (r *offerRepository) Create(offer *Offer) error {
	return r.db.Create(offer).Error
}

func (r *offerRepository) Update(id string, offer *Offer) error {
	err := r.db.Where("id = ?", id).Updates(offer).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *offerRepository) Delete(id string) error {
	err := r.db.Where("id = ?").Delete(&Offer{}).Error
	if err != nil {
		return err
	}
	return nil
}
