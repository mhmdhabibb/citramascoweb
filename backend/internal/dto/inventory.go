package dto

import "time"

type CreateItemRequest struct {
	Name         string `json:"name" binding:"required"`
	Category     string `json:"category" binding:"required"`
	Unit         string `json:"unit" binding:"required"`
	CurrentStock int    `json:"current_stock"`
	ReorderLevel int    `json:"reorder_level"`
}

type UpdateItemRequest struct {
	Name         string `json:"name"`
	Category     string `json:"category"`
	Unit         string `json:"unit"`
	ReorderLevel int    `json:"reorder_level"`
	Status       string `json:"status"`
}

type StockInRequest struct {
	Quantity  int    `json:"quantity" binding:"required"`
	Reference string `json:"reference"`
	Note      string `json:"note"`
}

type UsageRequest struct {
	Quantity  int    `json:"quantity" binding:"required"`
	Reference string `json:"reference"`
	Note      string `json:"note"`
}

type StockTakeRequest struct {
	ItemId      string `json:"item_id" binding:"required"`
	ActualStock int    `json:"actual_stock" binding:"required"`
	Note        string `json:"note"`
}

type InventoryReportItem struct {
	Item          InventoryItemReport `json:"item"`
	OpeningStock  int                 `json:"opening_stock"`
	StockIn       int                 `json:"stock_in"`
	Usage         int                 `json:"usage"`
	ClosingStock  int                 `json:"closing_stock"`
	IsLowStock    bool                `json:"is_low_stock"`
}

type InventoryItemReport struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	Unit         string `json:"unit"`
	ReorderLevel int    `json:"reorder_level"`
}

type ReportResponse struct {
	From       time.Time              `json:"from"`
	To         time.Time              `json:"to"`
	Summary    []InventoryReportItem  `json:"summary"`
	LowStock   []InventoryReportItem  `json:"low_stock"`
}
