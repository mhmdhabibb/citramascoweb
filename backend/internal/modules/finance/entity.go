package finance

import (
	"time"
)

type ChartOfAccount struct {
	Code      string    `json:"code" gorm:"primaryKey;type:varchar(50)"`
	Name      string    `json:"name" gorm:"type:varchar(191)"`
	Type      string    `json:"type" gorm:"type:varchar(50)"` // e.g., Asset, Liability, Equity, Revenue, Expense
	ParentCode *string   `json:"parent_code" gorm:"type:varchar(50)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CashTransaction struct {
	Id          string         `json:"id" gorm:"primaryKey;type:varchar(191)"`
	Date        time.Time      `json:"date"`
	Description string         `json:"description" gorm:"type:text"`
	Type        string         `json:"type" gorm:"type:varchar(50)"` // "in" or "out"
	Amount      int            `json:"amount"`
	AccountId   string         `json:"account_id" gorm:"type:varchar(50)"`
	Account     ChartOfAccount `json:"account" gorm:"foreignKey:AccountId;references:Code"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type GeneralJournal struct {
	Id            string         `json:"id" gorm:"primaryKey;type:varchar(191)"`
	Date          time.Time      `json:"date"`
	Description   string         `json:"description" gorm:"type:text"`
	AccountCode   string         `json:"account_code" gorm:"type:varchar(50)"`
	Account       ChartOfAccount `json:"account" gorm:"foreignKey:AccountCode;references:Code"`
	Debit         int            `json:"debit" gorm:"default:0"`
	Credit        int            `json:"credit" gorm:"default:0"`
	ReservationId *string        `json:"reservation_id" gorm:"type:varchar(191)"` // nullable, to link back to reservation if generated automatically
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type Invoice struct {
	Id            string    `json:"id" gorm:"primaryKey;type:varchar(191)"`
	InvoiceNumber string    `json:"invoice_number" gorm:"type:varchar(191);unique"`
	Type          string    `json:"type" gorm:"type:varchar(50)"` // "AP" (Hutang) or "AR" (Piutang)
	PartnerName   string    `json:"partner_name" gorm:"type:varchar(191)"`
	Date          time.Time `json:"date"`
	DueDate       time.Time `json:"due_date"`
	Amount        int       `json:"amount"`
	Status        string    `json:"status" gorm:"type:varchar(50);default:'Unpaid'"` // Unpaid, Paid
	Description   string    `json:"description" gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
