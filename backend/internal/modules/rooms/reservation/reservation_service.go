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

	// Send confirmation email
	newReservation.Room = *roomData
	s.sendReservationEmail(newReservation, "Booking Confirmation - Citramasco Hotel", "Your booking has been successfully created and is currently awaiting approval.")

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

	statusChanged := false
	var oldStatus ReservationStatus
	if req.Status != "" {
		oldStatus = reservation.Status
		if ReservationStatus(req.Status) != oldStatus {
			reservation.Status = ReservationStatus(req.Status)
			statusChanged = true
		}
	}

	if req.NumberOfGuest != 0 {
		reservation.NumberOfGuest = req.NumberOfGuest
	}

	err = s.reservationRepo.Update(reservation, id)
	if err != nil {
		return err
	}

	if statusChanged {
		var subject string
		var message string

		switch reservation.Status {
		case ReservationStatusApproved:
			subject = "Booking Approved - Citramasco Hotel"
			message = "Congratulations! Your booking has been approved. We look forward to your arrival."
		case ReservationStatusCancel:
			subject = "Booking Cancelled - Citramasco Hotel"
			message = "Your booking has been successfully cancelled."
		case ReservationStatusRejected:
			subject = "Booking Rejected - Citramasco Hotel"
			message = "We regret to inform you that your booking has been rejected."
		}

		if message != "" {
			s.sendReservationEmail(reservation, subject, message)
		}
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

	s.sendReservationEmail(reservation, "Booking Approved - Citramasco Hotel", "Congratulations! Your booking has been approved. We look forward to your arrival.")

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

	s.sendReservationEmail(reservation, "Booking Cancelled - Citramasco Hotel", "Your booking has been successfully cancelled.")

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

	s.sendReservationEmail(reservation, "Booking Rejected - Citramasco Hotel", "We regret to inform you that your booking has been rejected.")

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

func (s *reservationService) sendReservationEmail(r *Reservation, subject, message string) {
	go func() {
		var roomName string = "Kamar"
		if r.Room.Name != "" {
			roomName = r.Room.Name
		} else {
			room, err := s.roomRepo.GetById(r.RoomId)
			if err == nil && room != nil {
				roomName = room.Name
			}
		}

		htmlContent := fmt.Sprintf(`
			<table cellpadding="0" cellspacing="0" width="100%%" style="font-family: 'Plus Jakarta Sans', 'Inter', 'Helvetica Neue', Helvetica, Arial, sans-serif; background-color: #f7f5f0; padding: 40px 20px;">
				<tr>
					<td align="center">
						<table cellpadding="0" cellspacing="0" width="100%%" style="max-width: 600px; background-color: #ffffff; border-radius: 16px; box-shadow: 0 8px 30px rgba(28, 22, 18, 0.06); overflow: hidden; border: 1px solid #ebdccb;">
							<!-- Header Banner -->
							<tr>
								<td align="center" style="background: linear-gradient(135deg, #1e1713 0%%, #3e332c 100%%); padding: 40px 30px; border-bottom: 3px solid #bfa280;">
									<h1 style="color: #ffffff; font-family: 'Playfair Display', Georgia, serif; font-size: 28px; font-weight: 600; letter-spacing: 2px; margin: 0 0 5px 0; text-transform: uppercase;">Citramasco Hotel</h1>
									<p style="color: #bfa280; font-size: 13px; font-weight: 500; letter-spacing: 3px; margin: 0; text-transform: uppercase;">Luxury & Comfort Await</p>
								</td>
							</tr>

							<!-- Content Area -->
							<tr>
								<td style="padding: 40px 35px 30px 35px;">
									<h2 style="color: #1e1713; font-size: 20px; font-weight: 600; margin-top: 0; margin-bottom: 15px;">Hello %s,</h2>
									<p style="color: #5c544e; font-size: 15px; line-height: 1.6; margin-bottom: 30px;">%s</p>

									<!-- Info Card Container -->
									<table cellpadding="0" cellspacing="0" width="100%%" style="background-color: #faf8f5; border: 1px solid #ebdccb; border-radius: 12px; padding: 25px; margin-bottom: 30px;">
										<!-- Booking Code Row -->
										<tr>
											<td align="center" style="border-bottom: 1px dashed #ebdccb; padding-bottom: 15px; margin-bottom: 15px;">
												<span style="color: #8c7355; font-size: 12px; font-weight: 600; text-transform: uppercase; letter-spacing: 1.5px; display: block; margin-bottom: 5px;">Booking Code</span>
												<strong style="color: #1e1713; font-size: 24px; font-weight: 700; letter-spacing: 1px;">%s</strong>
											</td>
										</tr>
										
										<!-- Spacing -->
										<tr><td height="15"></td></tr>

										<!-- Details Table -->
										<tr>
											<td>
												<table cellpadding="0" cellspacing="0" width="100%%">
													<!-- Guest Name -->
													<tr>
														<td align="left" style="color: #8c837a; font-size: 14px; font-weight: 500; padding: 6px 0;">Guest Name</td>
														<td align="right" style="color: #1e1713; font-size: 14px; font-weight: 600; padding: 6px 0;">%s</td>
													</tr>
													<!-- Guest Email -->
													<tr>
														<td align="left" style="color: #8c837a; font-size: 14px; font-weight: 500; padding: 6px 0;">Guest Email</td>
														<td align="right" style="color: #1e1713; font-size: 14px; font-weight: 600; padding: 6px 0;">%s</td>
													</tr>
													<!-- Room Type -->
													<tr>
														<td align="left" style="color: #8c837a; font-size: 14px; font-weight: 500; padding: 6px 0;">Room Type</td>
														<td align="right" style="color: #1e1713; font-size: 14px; font-weight: 600; padding: 6px 0;">%s</td>
													</tr>
													<!-- Check-In -->
													<tr>
														<td align="left" style="color: #8c837a; font-size: 14px; font-weight: 500; padding: 6px 0;">Check-In Date</td>
														<td align="right" style="color: #1e1713; font-size: 14px; font-weight: 600; padding: 6px 0;">%s</td>
													</tr>
													<!-- Check-Out -->
													<tr>
														<td align="left" style="color: #8c837a; font-size: 14px; font-weight: 500; padding: 6px 0;">Check-Out Date</td>
														<td align="right" style="color: #1e1713; font-size: 14px; font-weight: 600; padding: 6px 0;">%s</td>
													</tr>
													<!-- Number of Guests -->
													<tr>
														<td align="left" style="color: #8c837a; font-size: 14px; font-weight: 500; padding: 6px 0;">Guests</td>
														<td align="right" style="color: #1e1713; font-size: 14px; font-weight: 600; padding: 6px 0;">%d Guest(s)</td>
													</tr>
													<!-- Duration of Stay -->
													<tr>
														<td align="left" style="color: #8c837a; font-size: 14px; font-weight: 500; padding: 6px 0;">Duration of Stay</td>
														<td align="right" style="color: #1e1713; font-size: 14px; font-weight: 600; padding: 6px 0;">%d Night(s)</td>
													</tr>
												</table>
											</td>
										</tr>

										<!-- Spacing -->
										<tr><td height="15"></td></tr>

										<!-- Total Price -->
										<tr>
											<td style="border-top: 1px dashed #ebdccb; padding-top: 15px;">
												<table cellpadding="0" cellspacing="0" width="100%%">
													<tr>
														<td align="left" style="color: #1e1713; font-size: 16px; font-weight: 700;">Total Price</td>
														<td align="right" style="color: #bfa280; font-size: 20px; font-weight: 700;">Rp %s</td>
													</tr>
												</table>
											</td>
										</tr>
									</table>

									<!-- Decorative Quote -->
									<table cellpadding="0" cellspacing="0" width="100%%" style="text-align: center; margin-top: 35px;">
										<tr>
											<td style="color: #8c837a; font-size: 13px; line-height: 1.5; font-style: italic;">
												"A true sanctuary of comfort and elegance, designed for your ultimate peace of mind."
											</td>
										</tr>
									</table>
								</td>
							</tr>

							<!-- Footer -->
							<tr>
								<td align="center" style="background-color: #1e1713; padding: 25px 30px; border-top: 1px solid #2e2621;">
									<p style="color: #a89f97; font-size: 11px; margin: 0 0 8px 0;">© 2026 Citramasco Hotel. All rights reserved.</p>
									<p style="color: #8c837a; font-size: 10px; margin: 0;">If you have any questions, please contact our customer support.</p>
								</td>
							</tr>
						</table>
					</td>
				</tr>
			</table>
		`, r.FullName, message, r.Code, r.FullName, r.Email, roomName,
			time.Time(*r.CheckinDate).Format("02-01-2006"),
			time.Time(*r.CheckoutDate).Format("02-01-2006"),
			r.NumberOfGuest, r.TotalNight, utils.FormatDotSeparator(r.TotalPrice))

		_ = utils.SendEmail(r.Email, subject, htmlContent)
	}()
}
