package reservation

import (
	"citramascoweb-backend/internal/dto"
	"citramascoweb-backend/internal/modules/rooms"
	"citramascoweb-backend/pkg/utils"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type reservationService struct {
	reservationRepo ReservationRepositoryInterface
	roomRepo        rooms.RoomRepositoryInterface
}

func NewReservationService(reservationRepo ReservationRepositoryInterface, roomRepo rooms.RoomRepositoryInterface) *reservationService {
	return &reservationService{
		reservationRepo: reservationRepo,
		roomRepo:        roomRepo,
	}
}

func (s *reservationService) GetAll() ([]Reservation, error) {
	reservations, err := s.reservationRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return reservations, nil
}

func (s *reservationService) GetById(id string) (*Reservation, error) {
	reservation, err := s.reservationRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *reservationService) Store(req *dto.CreateReservationRequest) error {
	// 1. Fetch Room details to get the price
	roomData, err := s.roomRepo.GetById(req.RoomId)
	if err != nil {
		return fmt.Errorf("room not found: %v", err)
	}

	// 2. Parse check-in and check-out dates
	if req.CheckInDate == "" {
		return errors.New("check-in date is required")
	}
	checkinCustom, err := dto.ParseCustomDate(req.CheckInDate)
	if err != nil || checkinCustom == nil {
		return errors.New("invalid check-in date format")
	}
	checkinDate := time.Time(*checkinCustom)

	if req.CheckOutDate == "" {
		return errors.New("check-out date is required")
	}
	checkoutCustom, err := dto.ParseCustomDate(req.CheckOutDate)
	if err != nil || checkoutCustom == nil {
		return errors.New("invalid check-out date format")
	}
	checkoutDate := time.Time(*checkoutCustom)

	// 3. Calculate total nights
	duration := checkoutDate.Sub(checkinDate)
	totalNight := int(duration.Hours() / 24)
	if totalNight <= 0 {
		return errors.New("check-out date must be after check-in date")
	}

	// Check room availability before creating reservation
	availableRooms, err := s.reservationRepo.CheckAvailability(checkinDate, checkoutDate)
	if err != nil {
		return fmt.Errorf("failed to check room availability: %v", err)
	}

	isAvailable := false
	for _, room := range availableRooms {
		if room.Id == req.RoomId {
			isAvailable = true
			break
		}
	}
	if !isAvailable {
		return errors.New("room is not available for the selected dates")
	}

	// 4. Calculate total price
	pricePerNight := roomData.Price
	totalPriceBeforeDiscount := pricePerNight * totalNight
	totalPrice := totalPriceBeforeDiscount

	var offerApplied *string
	var isOfferVal bool = false

	if req.IsOffer != nil && *req.IsOffer && req.OfferCode != nil && *req.OfferCode != "" {
		offerData, err := s.roomRepo.GetOfferByCode(*req.OfferCode)
		if err != nil {
			return errors.New("invalid promo code")
		}

		// Check expiration (ValidStart & ValidEnd)
		now := time.Now()
		if offerData.ValidStart != nil && now.Before(time.Time(*offerData.ValidStart)) {
			return errors.New("the discount offer associated with this promo code is not active yet")
		}
		if offerData.ValidEnd != nil && now.After(time.Time(*offerData.ValidEnd)) {
			return errors.New("the discount offer associated with this promo code has expired")
		}

		// Check reservation quota (MaxQuota)
		if offerData.MaxQuota != nil && *offerData.MaxQuota <= 0 {
			return errors.New("the discount offer for this promo code has reached its usage limit")
		}

		discountVal := 0
		if offerData.Discount != nil {
			discountVal = *offerData.Discount
		}

		totalDiscount := totalPriceBeforeDiscount * discountVal / 100
		totalPrice = totalPriceBeforeDiscount - totalDiscount

		offerApplied = req.OfferCode
		isOfferVal = true
	}

	// 5. Generate random reservation code
	code, err := utils.GenerateRandomNumberCode(6)
	if err != nil {
		return err
	}

	newReservation := &Reservation{
		Id:            uuid.New().String(),
		Code:          code,
		RoomId:        req.RoomId,
		FullName:      req.FullName,
		Email:         req.Email,
		CheckinDate:   checkinCustom,
		CheckoutDate:  checkoutCustom,
		Price:         pricePerNight,
		TotalNight:    totalNight,
		TotalPrice:    totalPrice,
		Status:        ReservationStatusPending,
		NumberOfGuest: req.NumberOfGuest,
		IsOffer:       &isOfferVal,
		OfferCode:     offerApplied,
	}

	err = s.reservationRepo.Create(newReservation)
	if err != nil {
		return err
	}

	// Decrement offer quota if applied
	if offerApplied != nil {
		_ = s.roomRepo.DecrementOfferQuota(*offerApplied)
	}

	return nil
}

func (s *reservationService) Update(id string, req *dto.UpdateReservationRequest) error {
	reservation, err := s.reservationRepo.GetById(id)
	if err != nil {
		return err
	}

	if req.RoomId != "" && req.RoomId != reservation.RoomId {
		// Recalculate price if room changes
		roomData, err := s.roomRepo.GetById(req.RoomId)
		if err != nil {
			return fmt.Errorf("room not found: %v", err)
		}

		reservation.Price = roomData.Price
		reservation.RoomId = req.RoomId
	}

	// Track if dates changed to recalculate total price/nights
	datesChanged := false

	if req.CheckInDate != "" {
		checkinCustom, err := dto.ParseCustomDate(req.CheckInDate)
		if err != nil {
			return errors.New("invalid check-in date format")
		}
		reservation.CheckinDate = checkinCustom
		datesChanged = true
	}

	if req.CheckOutDate != "" {
		checkoutCustom, err := dto.ParseCustomDate(req.CheckOutDate)
		if err != nil {
			return errors.New("invalid check-out date format")
		}
		reservation.CheckoutDate = checkoutCustom
		datesChanged = true
	}

	if datesChanged {
		if reservation.CheckinDate == nil || reservation.CheckoutDate == nil {
			return errors.New("check-in and check-out dates are required")
		}
		duration := time.Time(*reservation.CheckoutDate).Sub(time.Time(*reservation.CheckinDate))
		totalNight := int(duration.Hours() / 24)
		if totalNight <= 0 {
			return errors.New("check-out date must be after check-in date")
		}
		reservation.TotalNight = totalNight
	}

	// Calculate total price based on promo codes
	totalPriceBeforeDiscount := reservation.Price * reservation.TotalNight
	totalPrice := totalPriceBeforeDiscount

	// Determine if reservation uses offer
	isOffer := false
	if reservation.IsOffer != nil {
		isOffer = *reservation.IsOffer
	}
	if req.IsOffer != nil {
		isOffer = *req.IsOffer
	}

	offerCode := ""
	if reservation.OfferCode != nil {
		offerCode = *reservation.OfferCode
	}
	if req.OfferCode != nil {
		offerCode = *req.OfferCode
	}

	// If the offer code changed or was removed, increment old quota back
	if reservation.OfferCode != nil && *reservation.OfferCode != "" && (offerCode != *reservation.OfferCode || !isOffer) {
		_ = s.roomRepo.IncrementOfferQuota(*reservation.OfferCode)
		reservation.OfferCode = nil
		isFalse := false
		reservation.IsOffer = &isFalse
	}

	// Apply new offer if set
	if isOffer && offerCode != "" {
		offerData, err := s.roomRepo.GetOfferByCode(offerCode)
		if err == nil {
			// If code is different from before, we validate & decrement quota
			if reservation.OfferCode == nil || *reservation.OfferCode != offerCode {
				now := time.Now()
				if offerData.ValidStart != nil && now.Before(time.Time(*offerData.ValidStart)) {
					return errors.New("the discount offer associated with this promo code is not active yet")
				}
				if offerData.ValidEnd != nil && now.After(time.Time(*offerData.ValidEnd)) {
					return errors.New("the discount offer associated with this promo code has expired")
				}
				if offerData.MaxQuota != nil && *offerData.MaxQuota <= 0 {
					return errors.New("the discount offer for this promo code has reached its usage limit")
				}
				_ = s.roomRepo.DecrementOfferQuota(offerCode)
			}

			discountVal := 0
			if offerData.Discount != nil {
				discountVal = *offerData.Discount
			}
			totalDiscount := totalPriceBeforeDiscount * discountVal / 100
			totalPrice = totalPriceBeforeDiscount - totalDiscount

			isTrue := true
			reservation.IsOffer = &isTrue
			reservation.OfferCode = &offerCode
		} else {
			return errors.New("invalid promo code")
		}
	}

	reservation.TotalPrice = totalPrice

	if req.Status != "" {
		reservation.Status = ReservationStatus(req.Status)
	}

	if req.NumberOfGuest != 0 {
		reservation.NumberOfGuest = req.NumberOfGuest
	}

	err = s.reservationRepo.Update(reservation, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *reservationService) Delete(id string) error {
	err := s.reservationRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *reservationService) ApproveReservation(id string) error {

	reservation, err := s.reservationRepo.GetById(id)
	if err != nil {
		return errors.New("Reservation Not Found!")
	}

	if reservation.Status == ReservationStatusApproved {
		return errors.New("reservation is already approved")
	}
	if reservation.Status == ReservationStatusRejected {
		return errors.New("reservation is already rejected")
	}
	if reservation.Status == ReservationStatusCancel {
		return errors.New("reservation is already cancelled")
	}

	err = s.reservationRepo.ApproveReservation(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *reservationService) CancelReservation(id string) error {

	reservation, err := s.reservationRepo.GetById(id)
	if err != nil {
		return errors.New("Reservation Not Found!")
	}

	if reservation.Status == ReservationStatusApproved {
		return errors.New("reservation is already approved")
	}
	if reservation.Status == ReservationStatusRejected {
		return errors.New("reservation is already rejected")
	}
	if reservation.Status == ReservationStatusCancel {
		return errors.New("reservation is already cancelled")
	}

	err = s.reservationRepo.CancelReservation(id)
	if err != nil {
		return err
	}

	// Increment offer quota back
	if reservation.IsOffer != nil && *reservation.IsOffer && reservation.OfferCode != nil && *reservation.OfferCode != "" {
		_ = s.roomRepo.IncrementOfferQuota(*reservation.OfferCode)
	}

	return nil
}

func (s *reservationService) RejectReservation(id string) error {
	reservation, err := s.reservationRepo.GetById(id)
	if err != nil {
		return errors.New("Reservation Not Found!")
	}

	if reservation.Status == ReservationStatusApproved {
		return errors.New("reservation is already approved")
	}
	if reservation.Status == ReservationStatusRejected {
		return errors.New("reservation is already rejected")
	}
	if reservation.Status == ReservationStatusCancel {
		return errors.New("reservation is already cancelled")
	}
	err = s.reservationRepo.RejectReservation(id)
	if err != nil {
		return err
	}

	// Increment offer quota back
	if reservation.IsOffer != nil && *reservation.IsOffer && reservation.OfferCode != nil && *reservation.OfferCode != "" {
		_ = s.roomRepo.IncrementOfferQuota(*reservation.OfferCode)
	}

	return nil
}

func (s *reservationService) CheckAvailability(req *dto.CheckAvailabilityRequest) (interface{}, error) {
	if req.CheckInDate == "" {
		return nil, errors.New("check-in date is required")
	}
	checkinCustom, err := dto.ParseCustomDate(req.CheckInDate)
	if err != nil || checkinCustom == nil {
		return nil, errors.New("invalid check-in date format")
	}
	checkinDate := time.Time(*checkinCustom)

	if req.CheckOutDate == "" {
		return nil, errors.New("check-out date is required")
	}
	checkoutCustom, err := dto.ParseCustomDate(req.CheckOutDate)
	if err != nil || checkoutCustom == nil {
		return nil, errors.New("invalid check-out date format")
	}
	checkoutDate := time.Time(*checkoutCustom)

	if checkoutDate.Before(checkinDate) || checkoutDate.Equal(checkinDate) {
		return nil, errors.New("check-out date must be after check-in date")
	}

	availableRooms, err := s.reservationRepo.CheckAvailability(checkinDate, checkoutDate)
	if err != nil {
		return nil, err
	}

	if req.RoomId != "" {
		isAvailable := false
		for _, room := range availableRooms {
			if room.Id == req.RoomId {
				isAvailable = true
				break
			}
		}
		status := "tidak tersedia"
		if isAvailable {
			status = "tersedia"
		}
		return map[string]interface{}{
			"available": isAvailable,
			"status":    status,
		}, nil
	}

	allRooms, err := s.roomRepo.GetAll()
	if err != nil {
		return nil, err
	}

	availableMap := make(map[string]bool)
	for _, r := range availableRooms {
		availableMap[r.Id] = true
	}

	for i := range allRooms {
		if availableMap[allRooms[i].Id] {
			allRooms[i].Status = "tersedia"
			allRooms[i].AvailabilityStatus = "tersedia"
		} else {
			allRooms[i].Status = "tidak tersedia"
			allRooms[i].AvailabilityStatus = "tidak tersedia"
		}
	}

	return allRooms, nil
}
