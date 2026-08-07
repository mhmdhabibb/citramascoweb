package inventory

import (
	"citramascoweb-backend/internal/dto"
	"errors"
	"time"

	"github.com/google/uuid"
)

type inventoryService struct {
	repo InventoryRepositoryInterface
}

func NewInventoryService(repo InventoryRepositoryInterface) *inventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) ListItems() ([]InventoryItem, error) {
	return s.repo.GetItems()
}

func (s *inventoryService) CreateItem(req *dto.CreateItemRequest) (*InventoryItem, error) {
	now := time.Now()
	item := &InventoryItem{
		Id:           uuid.New().String(),
		Name:         req.Name,
		Category:     req.Category,
		Unit:         req.Unit,
		CurrentStock: req.CurrentStock,
		ReorderLevel: req.ReorderLevel,
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.CreateItem(item); err != nil {
		return nil, err
	}

	// Record initial stock as an opening stock-in transaction
	if item.CurrentStock > 0 {
		openingTx := &InventoryTransaction{
			Id:           uuid.New().String(),
			ItemId:       item.Id,
			Type:         TransactionTypeStockIn,
			Quantity:     item.CurrentStock,
			BalanceAfter: item.CurrentStock,
			Reference:    "OPENING",
			Note:         "Opening stock",
			CreatedAt:    now,
		}
		_ = s.repo.CreateTransaction(openingTx)
	}

	return item, nil
}

func (s *inventoryService) UpdateItem(req *dto.UpdateItemRequest, id string) error {
	item, err := s.repo.GetItemById(id)
	if err != nil {
		return err
	}
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.Unit != "" {
		item.Unit = req.Unit
	}
	if req.ReorderLevel != 0 {
		item.ReorderLevel = req.ReorderLevel
	}
	if req.Status != "" {
		item.Status = req.Status
	}
	item.UpdatedAt = time.Now()
	return s.repo.UpdateItem(item, id)
}

func (s *inventoryService) DeleteItem(id string) error {
	return s.repo.DeleteItem(id)
}

func (s *inventoryService) StockIn(id string, req *dto.StockInRequest, actor string) error {
	item, err := s.repo.GetItemById(id)
	if err != nil {
		return err
	}
	if req.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	newStock := item.CurrentStock + req.Quantity
	_ = s.repo.UpdateItemStock(id, newStock)

	tx := &InventoryTransaction{
		Id:           uuid.New().String(),
		ItemId:       item.Id,
		Type:         TransactionTypeStockIn,
		Quantity:     req.Quantity,
		BalanceAfter: newStock,
		Reference:    req.Reference,
		Note:         req.Note,
		CreatedBy:    actor,
		CreatedAt:    time.Now(),
	}
	return s.repo.CreateTransaction(tx)
}

func (s *inventoryService) Usage(id string, req *dto.UsageRequest, actor string) error {
	item, err := s.repo.GetItemById(id)
	if err != nil {
		return err
	}
	if req.Quantity <= 0 {
		return errors.New("quantity should be greater than zero")
	}
	if item.CurrentStock < req.Quantity {
		return errors.New("insufficient stock for this item")
	}

	newStock := item.CurrentStock - req.Quantity
	_ = s.repo.UpdateItemStock(id, newStock)

	tx := &InventoryTransaction{
		Id:           uuid.New().String(),
		ItemId:       item.Id,
		Type:         TransactionTypeUsage,
		Quantity:     req.Quantity,
		BalanceAfter: newStock,
		Reference:    req.Reference,
		Note:         req.Note,
		CreatedBy:    actor,
		CreatedAt:    time.Now(),
	}
	return s.repo.CreateTransaction(tx)
}

func (s *inventoryService) StockTake(req *dto.StockTakeRequest, actor string) error {
	item, err := s.repo.GetItemById(req.ItemId)
	if err != nil {
		return err
	}
	if req.ActualStock < 0 {
		return errors.New("actual stock cannot be negative")
	}

	systemStock := item.CurrentStock
	variance := req.ActualStock - systemStock

	st := &InventoryStockTake{
		Id:          uuid.New().String(),
		ItemId:      item.Id,
		SystemStock: systemStock,
		ActualStock: req.ActualStock,
		Variance:    variance,
		Note:        req.Note,
		CheckedBy:   actor,
		CreatedAt:   time.Now(),
	}
	if err := s.repo.CreateStockTake(st); err != nil {
		return err
	}

	_ = s.repo.UpdateItemStock(item.Id, req.ActualStock)

	// Record the difference as an adjustment transaction
	if variance != 0 {
		txType := TransactionTypeAdjust
		tx := &InventoryTransaction{
			Id:           uuid.New().String(),
			ItemId:       item.Id,
			Type:         txType,
			Quantity:     variance,
			BalanceAfter: req.ActualStock,
			Reference:    "STOCK-TAKE",
			Note:         req.Note,
			CreatedBy:    actor,
			CreatedAt:    time.Now(),
		}
		_ = s.repo.CreateTransaction(tx)
	}

	return nil
}

func (s *inventoryService) ListTransactions(from time.Time, to time.Time) ([]InventoryTransaction, error) {
	return s.repo.GetTransactions(from, to)
}

func (s *inventoryService) Report(from time.Time, to time.Time) (*dto.ReportResponse, error) {
	items, err := s.repo.GetItems()
	if err != nil {
		return nil, err
	}
	txs, err := s.repo.GetTransactions(from, to)
	if err != nil {
		return nil, err
	}

	stockInByItem := map[string]int{}
	usageByItem := map[string]int{}
	firstTxAt := map[string]time.Time{}

	for _, tx := range txs {
		key := tx.ItemId
		switch tx.Type {
		case TransactionTypeStockIn:
			stockInByItem[key] += tx.Quantity
		case TransactionTypeUsage:
			usageByItem[key] += tx.Quantity
		}
		if t, ok := firstTxAt[key]; !ok || tx.CreatedAt.Before(t) {
			firstTxAt[key] = tx.CreatedAt
		}
	}

	summary := make([]dto.InventoryReportItem, 0, len(items))
	lowStock := make([]dto.InventoryReportItem, 0)

	for _, item := range items {
		opening := item.CurrentStock - stockInByItem[item.Id] + usageByItem[item.Id]
		if opening < 0 {
			opening = 0
		}
		row := dto.InventoryReportItem{
			Item: dto.InventoryItemReport{
				Id:           item.Id,
				Name:         item.Name,
				Category:     item.Category,
				Unit:         item.Unit,
				ReorderLevel: item.ReorderLevel,
			},
			OpeningStock: opening,
			StockIn:      stockInByItem[item.Id],
			Usage:        usageByItem[item.Id],
			ClosingStock: item.CurrentStock,
		}

		row.IsLowStock = item.ReorderLevel > 0 && item.CurrentStock <= item.ReorderLevel
		summary = append(summary, row)
		if row.IsLowStock {
			lowStock = append(lowStock, row)
		}
	}

	return &dto.ReportResponse{
		From:     from,
		To:       to,
		Summary:  summary,
		LowStock: lowStock,
	}, nil
}