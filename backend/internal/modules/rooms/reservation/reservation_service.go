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
	if req.CheckInDate == nil {
		return errors.New("check-in date is required")
	}
	checkinDate := time.Time(*req.CheckInDate)

	if req.CheckOutDate == nil {
		return errors.New("check-out date is required")
	}
	checkoutDate := time.Time(*req.CheckOutDate)

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
	var offerApplied *string
	if roomData.IsOffer != nil && *roomData.IsOffer && roomData.OfferCode != nil && *roomData.OfferCode != "" {
		offerData, err := s.roomRepo.GetOfferByCode(*roomData.OfferCode)
		if err == nil {
			// Check expiration (ValidStart & ValidEnd)
			now := time.Now()
			if offerData.ValidStart != nil && now.Before(time.Time(*offerData.ValidStart)) {
				return errors.New("the discount offer associated with this room is not active yet")
			}
			if offerData.ValidEnd != nil && now.After(time.Time(*offerData.ValidEnd)) {
				return errors.New("the discount offer associated with this room has expired")
			}

			// Check reservation quota (MaxQuota)
			if offerData.MaxQuota != nil && *offerData.MaxQuota <= 0 {
				return errors.New("the discount offer for this room has reached its usage limit")
			}

			if roomData.PriceAfterDiscount != nil {
				pricePerNight = *roomData.PriceAfterDiscount
			}
			offerApplied = roomData.OfferCode
		}
	}
	totalPrice := pricePerNight * totalNight

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
		CheckinDate:   req.CheckInDate,
		CheckoutDate:  req.CheckOutDate,
		Price:         pricePerNight,
		TotalNight:    totalNight,
		TotalPrice:    totalPrice,
		Status:        ReservationStatusPending,
		NumberOfGuest: req.NumberOfGuest,
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
		// Increment old room offer quota
		oldRoom, err := s.roomRepo.GetById(reservation.RoomId)
		if err == nil && oldRoom.IsOffer != nil && *oldRoom.IsOffer && oldRoom.OfferCode != nil && *oldRoom.OfferCode != "" {
			_ = s.roomRepo.IncrementOfferQuota(*oldRoom.OfferCode)
		}

		// Recalculate price if room changes
		roomData, err := s.roomRepo.GetById(req.RoomId)
		if err != nil {
			return fmt.Errorf("room not found: %v", err)
		}
		
		pricePerNight := roomData.Price
		if roomData.IsOffer != nil && *roomData.IsOffer && roomData.OfferCode != nil && *roomData.OfferCode != "" {
			offerData, err := s.roomRepo.GetOfferByCode(*roomData.OfferCode)
			if err == nil {
				// Check expiration (ValidStart & ValidEnd)
				now := time.Now()
				if offerData.ValidStart != nil && now.Before(time.Time(*offerData.ValidStart)) {
					return errors.New("the discount offer associated with the new room is not active yet")
				}
				if offerData.ValidEnd != nil && now.After(time.Time(*offerData.ValidEnd)) {
					return errors.New("the discount offer associated with the new room has expired")
				}

				// Check reservation quota (MaxQuota)
				if offerData.MaxQuota != nil && *offerData.MaxQuota <= 0 {
					return errors.New("the discount offer for the new room has reached its usage limit")
				}

				if roomData.PriceAfterDiscount != nil {
					pricePerNight = *roomData.PriceAfterDiscount
				}

				// Decrement new room offer quota
				_ = s.roomRepo.DecrementOfferQuota(*roomData.OfferCode)
			}
		}
		reservation.Price = pricePerNight
		reservation.RoomId = req.RoomId
	}

	// Track if dates changed to recalculate total price/nights
	datesChanged := false

	if req.CheckInDate != nil {
		reservation.CheckinDate = req.CheckInDate
		datesChanged = true
	}

	if req.CheckOutDate != nil {
		reservation.CheckoutDate = req.CheckOutDate
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

	// Always update total price based on price and total nights
	reservation.TotalPrice = reservation.Price * reservation.TotalNight

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
	if reservation.Room.IsOffer != nil && *reservation.Room.IsOffer && reservation.Room.OfferCode != nil && *reservation.Room.OfferCode != "" {
		_ = s.roomRepo.IncrementOfferQuota(*reservation.Room.OfferCode)
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
	if reservation.Room.IsOffer != nil && *reservation.Room.IsOffer && reservation.Room.OfferCode != nil && *reservation.Room.OfferCode != "" {
		_ = s.roomRepo.IncrementOfferQuota(*reservation.Room.OfferCode)
	}

	return nil
}

func (s *reservationService) CheckAvailability(req *dto.CheckAvailabilityRequest) (interface{}, error) {
	if req.CheckInDate == nil {
		return nil, errors.New("check-in date is required")
	}
	checkinDate := time.Time(*req.CheckInDate)

	if req.CheckOutDate == nil {
		return nil, errors.New("check-out date is required")
	}
	checkoutDate := time.Time(*req.CheckOutDate)

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
