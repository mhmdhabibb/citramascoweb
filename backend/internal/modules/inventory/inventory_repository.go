package inventory

import (
	"time"

	"gorm.io/gorm"
)

type InventoryRepositoryInterface interface {
	// Items
	GetItems() ([]InventoryItem, error)
	GetItemById(id string) (*InventoryItem, error)
	CreateItem(item *InventoryItem) error
	UpdateItem(item *InventoryItem, id string) error
	DeleteItem(id string) error
	UpdateItemStock(id string, newStock int) error

	// Transactions
	CreateTransaction(tx *InventoryTransaction) error
	GetTransactions(from time.Time, to time.Time) ([]InventoryTransaction, error)

// Stock Take
	CreateStockTake(st *InventoryStockTake) error
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepositoryInterface {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) GetItems() ([]InventoryItem, error) {
	var items []InventoryItem
	err := r.db.Order("name asc").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *inventoryRepository) GetItemById(id string) (*InventoryItem, error) {
	var item InventoryItem
	err := r.db.Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *inventoryRepository) CreateItem(item *InventoryItem) error {
	return r.db.Create(item).Error
}

func (r *inventoryRepository) UpdateItem(item *InventoryItem, id string) error {
	return r.db.Model(&InventoryItem{}).Where("id = ?", id).Updates(item).Error
}

func (r *inventoryRepository) DeleteItem(id string) error {
	return r.db.Where("id = ?", id).Delete(&InventoryItem{}).Error
}

func (r *inventoryRepository) UpdateItemStock(id string, newStock int) error {
	return r.db.Model(&InventoryItem{}).Where("id = ?", id).Update("current_stock", newStock).Error
}

func (r *inventoryRepository) CreateTransaction(tx *InventoryTransaction) error {
	return r.db.Create(tx).Error
}

func (r *inventoryRepository) GetTransactions(from time.Time, to time.Time) ([]InventoryTransaction, error) {
	var txs []InventoryTransaction
	q := r.db.Preload("Item")
	if !from.IsZero() {
		q = q.Where("created_at >= ?", from)
	}
	if !to.IsZero() {
		q = q.Where("created_at <= ?", to)
	}
	err := q.Order("created_at desc").Find(&txs).Error
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (r *inventoryRepository) CreateStockTake(stt *InventoryStockTake) error {
	return r.db.Create(stt).Error
}