package smartlock

import "time"

type CreateCredentialRequest struct {
	ReservationId   *string   `json:"reservation_id"`
	RoomUnitId      string    `json:"room_unit_id" binding:"required"`
	GuestId         *string   `json:"guest_id"`
	CredentialType  string    `json:"credential_type" binding:"required"`
	CredentialValue string    `json:"credential_value"` // if empty, auto-generate 6-digit PIN
	ValidFrom       time.Time `json:"valid_from" binding:"required"`
	ValidUntil      time.Time `json:"valid_until" binding:"required"`
}

type UpdateCredentialStatusRequest struct {
	Status string `json:"status" binding:"required"` // active, revoked
}

type RecordAccessLogRequest struct {
	RoomUnitId     string  `json:"room_unit_id" binding:"required"`
	ReservationId  *string `json:"reservation_id"`
	CredentialType string  `json:"credential_type"`
	UnlockedByName string  `json:"unlocked_by_name"`
	Success        bool    `json:"success"`
	RawPayload     string  `json:"raw_payload"`
}
