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
	checkinDate, err := time.Parse("2006-01-02", req.CheckInDate)
	if err != nil {
		checkinDate, err = time.Parse(time.RFC3339, req.CheckInDate)
		if err != nil {
			return fmt.Errorf("invalid check-in date: %v", err)
		}
	}

	checkoutDate, err := time.Parse("2006-01-02", req.CheckOutDate)
	if err != nil {
		checkoutDate, err = time.Parse(time.RFC3339, req.CheckOutDate)
		if err != nil {
			return fmt.Errorf("invalid check-out date: %v", err)
		}
	}

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
	totalPrice := roomData.Price * totalNight

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
		CheckinDate:   checkinDate,
		CheckoutDate:  checkoutDate,
		Price:         roomData.Price,
		TotalNight:    totalNight,
		TotalPrice:    totalPrice,
		Status:        ReservationStatusPending,
		NumberOfGuest: req.NumberOfGuest,
	}

	err = s.reservationRepo.Create(newReservation)
	if err != nil {
		return err
	}

	return nil
}

func (s *reservationService) Update(id string, req *dto.UpdateReservationRequest) error {
	reservation, err := s.reservationRepo.GetById(id)
	if err != nil {
		return err
	}

	if req.RoomId != "" {
		reservation.RoomId = req.RoomId
		// Recalculate price if room changes
		roomData, err := s.roomRepo.GetById(req.RoomId)
		if err != nil {
			return fmt.Errorf("room not found: %v", err)
		}
		reservation.Price = roomData.Price
	}

	// Track if dates changed to recalculate total price/nights
	datesChanged := false

	if req.CheckInDate != "" {
		checkinDate, err := time.Parse("2006-01-02", req.CheckInDate)
		if err != nil {
			checkinDate, err = time.Parse(time.RFC3339, req.CheckInDate)
			if err != nil {
				return fmt.Errorf("invalid check-in date: %v", err)
			}
		}
		reservation.CheckinDate = checkinDate
		datesChanged = true
	}

	if req.CheckOutDate != "" {
		checkoutDate, err := time.Parse("2006-01-02", req.CheckOutDate)
		if err != nil {
			checkoutDate, err = time.Parse(time.RFC3339, req.CheckOutDate)
			if err != nil {
				return fmt.Errorf("invalid check-out date: %v", err)
			}
		}
		reservation.CheckoutDate = checkoutDate
		datesChanged = true
	}

	if datesChanged {
		duration := reservation.CheckoutDate.Sub(reservation.CheckinDate)
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

	return nil
}

func (s *reservationService) CheckAvailability(req *dto.CheckAvailabilityRequest) (interface{}, error) {
	// Parse check-in and check-out dates
	checkinDate, err := time.Parse("2006-01-02", req.CheckInDate)
	if err != nil {
		checkinDate, err = time.Parse(time.RFC3339, req.CheckInDate)
		if err != nil {
			return nil, fmt.Errorf("invalid check-in date: %v", err)
		}
	}

	checkoutDate, err := time.Parse("2006-01-02", req.CheckOutDate)
	if err != nil {
		checkoutDate, err = time.Parse(time.RFC3339, req.CheckOutDate)
		if err != nil {
			return nil, fmt.Errorf("invalid check-out date: %v", err)
		}
	}

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
