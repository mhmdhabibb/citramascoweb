package inventory

import "time"

type TransactionType string

const (
	TransactionTypeStockIn  TransactionType = "stock-in"
	TransactionTypeUsage    TransactionType = "usage"
	TransactionTypeAdjust   TransactionType = "adjustment"
)

type InventoryItem struct {
	Id             string    `gorm:"type:varchar(191);primaryKey" json:"id"`
	Name           string    `json:"name"`
	Category       string    `gorm:"type:varchar(191)" json:"category"`
	Unit           string    `gorm:"type:varchar(50)" json:"unit"`
	CurrentStock   int       `json:"current_stock"`
	ReorderLevel   int       `json:"reorder_level"`
	Status         string    `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type InventoryTransaction struct {
	Id             string    `gorm:"type:varchar(191);primaryKey" json:"id"`
	ItemId         string    `gorm:"type:varchar(191)" json:"item_id"`
	Item           InventoryItem `gorm:"foreignKey:ItemId;references:Id" json:"item"`
	Type           TransactionType `gorm:"type:varchar(50)" json:"type"`
	Quantity       int       `json:"quantity"`
	BalanceAfter   int       `json:"balance_after"`
	Reference      string    `gorm:"type:varchar(191)" json:"reference"`
	Note           string    `gorm:"type:text" json:"note"`
	CreatedBy      string    `gorm:"type:varchar(191)" json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
}

type InventoryStockTake struct {
	Id          string `gorm:"type:varchar(191);primaryKey" json:"id"`
	ItemId      string `gorm:"type:varchar(191)" json:"item_id"`
	Item        InventoryItem `gorm:"foreignKey:ItemId;references:Id" json:"item"`
	SystemStock int    `json:"system_stock"`
	ActualStock int    `json:"actual_stock"`
	Variance    int    `json:"variance"`
	Note        string `gorm:"type:text" json:"note"`
	CheckedBy   string `gorm:"type:varchar(191)" json:"checked_by"`
	CreatedAt   time.Time `json:"created_at"`
}