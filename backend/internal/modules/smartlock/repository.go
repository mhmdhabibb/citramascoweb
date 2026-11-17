package smartlock

import (
	"gorm.io/gorm"
)

type SmartLockRepository interface {
	FindCredentialsByUnit(unitId string) ([]RoomAccessCredential, error)
	FindCredentialsByReservation(resId string) ([]RoomAccessCredential, error)
	CreateCredential(cred *RoomAccessCredential) error
	UpdateCredential(cred *RoomAccessCredential) error
	FindCredentialById(id string) (*RoomAccessCredential, error)

	FindLogs(unitId string, limit, offset int) ([]RoomAccessLog, int64, error)
	CreateLog(log *RoomAccessLog) error
}

type smartLockRepository struct {
	db *gorm.DB
}

func NewSmartLockRepository(db *gorm.DB) SmartLockRepository {
	return &smartLockRepository{db: db}
}

func (r *smartLockRepository) FindCredentialsByUnit(unitId string) ([]RoomAccessCredential, error) {
	var list []RoomAccessCredential
	err := r.db.Preload("Guest").Where("room_unit_id = ?", unitId).Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *smartLockRepository) FindCredentialsByReservation(resId string) ([]RoomAccessCredential, error) {
	var list []RoomAccessCredential
	err := r.db.Preload("RoomUnit").Where("reservation_id = ?", resId).Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *smartLockRepository) CreateCredential(cred *RoomAccessCredential) error {
	return r.db.Create(cred).Error
}

func (r *smartLockRepository) UpdateCredential(cred *RoomAccessCredential) error {
	return r.db.Save(cred).Error
}

func (r *smartLockRepository) FindCredentialById(id string) (*RoomAccessCredential, error) {
	var cred RoomAccessCredential
	err := r.db.Preload("RoomUnit").Preload("Guest").Where("id = ?", id).First(&cred).Error
	if err != nil {
		return nil, err
	}
	return &cred, nil
}

func (r *smartLockRepository) FindLogs(unitId string, limit, offset int) ([]RoomAccessLog, int64, error) {
	var logs []RoomAccessLog
	var total int64
	q := r.db.Model(&RoomAccessLog{})
	if unitId != "" {
		q = q.Where("room_unit_id = ?", unitId)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if limit <= 0 {
		limit = 50
	}
	err := q.Preload("RoomUnit").Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

func (r *smartLockRepository) CreateLog(log *RoomAccessLog) error {
	return r.db.Create(log).Error
}
