package guest

import (
	"time"

	"gorm.io/gorm"
)

type IdentityType struct {
	Id        string         `gorm:"type:varchar(191);primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null;unique" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (IdentityType) TableName() string {
	return "identity_types"
}

type Guest struct {
	Id             string         `gorm:"type:varchar(191);primaryKey" json:"id"`
	Title          string         `gorm:"type:varchar(20)" json:"title"` // Mr, Mrs, Ms, Dr
	FirstName      string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName       string         `gorm:"type:varchar(100)" json:"last_name"`
	IdentityTypeId *string        `gorm:"type:varchar(191);index" json:"identity_type_id"`
	IdentityType   *IdentityType  `gorm:"foreignKey:IdentityTypeId;references:Id" json:"identity_type,omitempty"`
	IdentityNumber string         `gorm:"type:varchar(100);index" json:"identity_number"`
	IdentityImage  string         `gorm:"type:text" json:"identity_image"`
	Email          string         `gorm:"type:varchar(191);index" json:"email"`
	Phone          string         `gorm:"type:varchar(50);index" json:"phone"`
	Nationality    string         `gorm:"type:varchar(100);default:'Indonesia'" json:"nationality"`
	Address        string         `gorm:"type:text" json:"address"`
	GuestType      string         `gorm:"type:varchar(50);default:'regular'" json:"guest_type"` // regular, vip, corporate
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Guest) TableName() string {
	return "guests"
}
