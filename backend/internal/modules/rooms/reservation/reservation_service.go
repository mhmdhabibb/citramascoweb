package reservation

import (
	"citramascoweb-backend/internal/dto"
	"citramascoweb-backend/internal/modules/rooms"
	"citramascoweb-backend/pkg/utils"
	"errors"
	"fmt"
	"strings"
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

	// 🚨 ABSOLUTE LOCK: Prevent booking if the physical room is under maintenance
	if strings.ToLower(strings.TrimSpace(string(roomData.Status))) == "maintenance" {
		return errors.New("cannot process reservation: this room is currently under maintenance")
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
			return errors.New("invalid offer code")
		}

		// Check expiration (ValidStart & ValidEnd)
		now := time.Now().UTC()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		if offerData.ValidStart != nil && today.Before(time.Time(*offerData.ValidStart)) {
			return errors.New("offer code is not active yet")
		}
		if offerData.ValidEnd != nil && today.After(time.Time(*offerData.ValidEnd)) {
			return errors.New("offer code has expired")
		}

		// Validate stay check-in date is within offer validity range
		if offerData.ValidStart != nil && checkinDate.Before(time.Time(*offerData.ValidStart)) {
			return errors.New("offer code is not active ")
		}
		if offerData.ValidEnd != nil && checkinDate.After(time.Time(*offerData.ValidEnd)) {
			return errors.New("offer code has expired ")
		}

		// Check reservation quota (MaxQuota)
		if offerData.MaxQuota != nil && *offerData.MaxQuota <= 0 {
			return errors.New("offer code has reached its usage limit")
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
			// Always validate checkin date against offer validity period
			if reservation.CheckinDate != nil {
				checkinDate := time.Time(*reservation.CheckinDate)
				if offerData.ValidStart != nil && checkinDate.Before(time.Time(*offerData.ValidStart)) {
					return errors.New("offer code is not active for the selected check-in date")
				}
				if offerData.ValidEnd != nil && checkinDate.After(time.Time(*offerData.ValidEnd)) {
					return errors.New("offer code has expired for the selected check-in date")
				}
			}

			// If code is different from before, we validate & decrement quota
			if reservation.OfferCode == nil || *reservation.OfferCode != offerCode {
				now := time.Now().UTC()
				today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
				if offerData.ValidStart != nil && today.Before(time.Time(*offerData.ValidStart)) {
					return errors.New("offer code is not active yet")
				}
				if offerData.ValidEnd != nil && today.After(time.Time(*offerData.ValidEnd)) {
					return errors.New("offer code has expired")
				}
				if offerData.MaxQuota != nil && *offerData.MaxQuota <= 0 {
					return errors.New("offer code has reached its usage limit")
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
			return errors.New("invalid offer code")
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

	// 1. Ambil kamar yang bebas dari overlap reservasi dan bebas dari status maintenance
	availableRooms, err := s.reservationRepo.CheckAvailability(checkinDate, checkoutDate)
	if err != nil {
		return nil, err
	}

	// Jika request meminta pengecekan spesifik satu ID Kamar
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

	// 2. Ambil seluruh master data kamar untuk pemetaan status global
	allRooms, err := s.roomRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Pindahkan list kamar tersedia ke dalam map pencarian cepat
	availableMap := make(map[string]bool)
	for _, r := range availableRooms {
		availableMap[r.Id] = true
	}

	// 3. Sinkronisasi status gabungan antara keterisian (Availability) dan fisik (Maintenance)
	for i := range allRooms {
		currentStatus := strings.ToLower(string(allRooms[i].Status))

		// KONDISI KHUSUS: Jika kamar memang sedang di-maintenance oleh teknisi
		if currentStatus == "maintenance" {
			allRooms[i].Status = "maintenance"
			allRooms[i].AvailabilityStatus = "occupied" // Kamar dikunci agar tidak bisa dipilih di sistem booking
			continue
		}

		// Kondisi Normal: Petakan berdasarkan data overlap reservasi dari repository
		if availableMap[allRooms[i].Id] {
			allRooms[i].Status = "active" // Kembalikan ke Enum fisik yang valid untuk frontend
			allRooms[i].AvailabilityStatus = "available"
		} else {
			allRooms[i].Status = "active"
			allRooms[i].AvailabilityStatus = "occupied"
		}
	}

	return allRooms, nil
}
func (s *reservationService) CheckIn(id string) error {
	reservation, err := s.reservationRepo.GetById(id)
	if err != nil {
		return errors.New("reservation tidak ditemukan")
	}

	// Validasi: Hanya reservasi yang sudah di-approve oleh admin yang bisa check-in
	if reservation.Status != ReservationStatusApproved {
		return fmt.Errorf("gagal check-in, status reservasi saat ini masih '%s' (harus approved)", reservation.Status)
	}

	return s.reservationRepo.CheckIn(id, reservation.RoomId)
}

func (s *reservationService) CheckOut(id string) error {
	reservation, err := s.reservationRepo.GetById(id)
	if err != nil {
		return errors.New("reservation tidak ditemukan")
	}

	// Validasi: Hanya tamu yang sedang menginap (checked-in) yang bisa check-out
	if reservation.Status != ReservationStatusCheckedIn {
		return errors.New("gagal check-out, tamu belum berstatus checked-in")
	}

	return s.reservationRepo.CheckOut(id, reservation.RoomId)
}
