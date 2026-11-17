package smartlock

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
)

type SmartLockService interface {
	GetCredentialsByUnit(unitId string) ([]RoomAccessCredential, error)
	GetCredentialsByReservation(resId string) ([]RoomAccessCredential, error)
	GenerateCredential(req *CreateCredentialRequest) (*RoomAccessCredential, error)
	RevokeCredential(id string) error

	GetAccessLogs(unitId string, page, limit int) ([]RoomAccessLog, int64, error)
	RecordLog(req *RecordAccessLogRequest) (*RoomAccessLog, error)
}

type smartLockService struct {
	repo SmartLockRepository
}

func NewSmartLockService(repo SmartLockRepository) SmartLockService {
	return &smartLockService{repo: repo}
}

func (s *smartLockService) GetCredentialsByUnit(unitId string) ([]RoomAccessCredential, error) {
	return s.repo.FindCredentialsByUnit(unitId)
}

func (s *smartLockService) GetCredentialsByReservation(resId string) ([]RoomAccessCredential, error) {
	return s.repo.FindCredentialsByReservation(resId)
}

func (s *smartLockService) GenerateCredential(req *CreateCredentialRequest) (*RoomAccessCredential, error) {
	val := req.CredentialValue
	if val == "" {
		// generate random 6-digit PIN
		n, _ := rand.Int(rand.Reader, big.NewInt(900000))
		val = fmt.Sprintf("%06d", n.Int64()+100000)
	}

	cred := &RoomAccessCredential{
		Id:              uuid.New().String(),
		ReservationId:   req.ReservationId,
		RoomUnitId:      req.RoomUnitId,
		GuestId:         req.GuestId,
		CredentialType:  req.CredentialType,
		CredentialValue: val,
		ValidFrom:       req.ValidFrom,
		ValidUntil:      req.ValidUntil,
		Status:          "active",
		SyncStatus:      "synced",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.CreateCredential(cred); err != nil {
		return nil, err
	}
	return cred, nil
}

func (s *smartLockService) RevokeCredential(id string) error {
	cred, err := s.repo.FindCredentialById(id)
	if err != nil {
		return err
	}
	cred.Status = "revoked"
	cred.UpdatedAt = time.Now()
	return s.repo.UpdateCredential(cred)
}

func (s *smartLockService) GetAccessLogs(unitId string, page, limit int) ([]RoomAccessLog, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 50
	}
	offset := (page - 1) * limit
	return s.repo.FindLogs(unitId, limit, offset)
}

func (s *smartLockService) RecordLog(req *RecordAccessLogRequest) (*RoomAccessLog, error) {
	log := &RoomAccessLog{
		Id:             uuid.New().String(),
		RoomUnitId:     req.RoomUnitId,
		ReservationId:  req.ReservationId,
		CredentialType: req.CredentialType,
		UnlockedByName: req.UnlockedByName,
		UnlockedAt:     time.Now(),
		Success:        req.Success,
		RawPayload:     req.RawPayload,
		CreatedAt:      time.Now(),
	}
	if err := s.repo.CreateLog(log); err != nil {
		return nil, err
	}
	return log, nil
}
