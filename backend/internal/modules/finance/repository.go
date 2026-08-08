package finance

import (
	"citramascoweb-backend/internal/dto"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FinanceRepositoryInterface interface {
	// COA
	GetAllCOA() ([]ChartOfAccount, error)
	GetCOAByCode(code string) (*ChartOfAccount, error)
	CreateCOA(coa *ChartOfAccount) error
	UpdateCOA(coa *ChartOfAccount, code string) error
	DeleteCOA(code string) error
	BulkUpsertCOA(coas []ChartOfAccount) error

	// Cash Transactions
	GetAllCashTransactions() ([]CashTransaction, error)
	GetCashTransactionById(id string) (*CashTransaction, error)
	CreateCashTransaction(transaction *CashTransaction) error
	DeleteCashTransaction(id string) error

	// General Journal
	GetAllGeneralJournals() ([]GeneralJournal, error)
	CreateGeneralJournal(journal *GeneralJournal) error
	CreateJournalEntries(journals []GeneralJournal) error // for multiple entries at once (e.g., debit and credit)

	// Verify Payment
	VerifyPayment(req *dto.VerifyPaymentRequest) error

	// Invoices (AP/AR)
	GetAllInvoices() ([]Invoice, error)
	GetInvoiceById(id string) (*Invoice, error)
	CreateInvoice(invoice *Invoice) error
	UpdateInvoice(invoice *Invoice) error
	DeleteInvoice(id string) error
}

type financeRepository struct {
	db *gorm.DB
}

func NewFinanceRepository(db *gorm.DB) FinanceRepositoryInterface {
	return &financeRepository{db: db}
}

// COA
func (r *financeRepository) GetAllCOA() ([]ChartOfAccount, error) {
	var coas []ChartOfAccount
	err := r.db.Find(&coas).Error
	return coas, err
}

func (r *financeRepository) GetCOAByCode(code string) (*ChartOfAccount, error) {
	var coa ChartOfAccount
	err := r.db.Where("code = ?", code).First(&coa).Error
	if err != nil {
		return nil, err
	}
	return &coa, nil
}

func (r *financeRepository) CreateCOA(coa *ChartOfAccount) error {
	return r.db.Create(coa).Error
}

func (r *financeRepository) UpdateCOA(coa *ChartOfAccount, code string) error {
	return r.db.Model(&ChartOfAccount{}).Where("code = ?", code).Updates(coa).Error
}

func (r *financeRepository) DeleteCOA(code string) error {
	return r.db.Where("code = ?", code).Delete(&ChartOfAccount{}).Error
}

func (r *financeRepository) BulkUpsertCOA(coas []ChartOfAccount) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "type", "parent_code", "updated_at"}),
	}).CreateInBatches(&coas, 100).Error
}

// Cash Transactions
func (r *financeRepository) GetAllCashTransactions() ([]CashTransaction, error) {
	var transactions []CashTransaction
	err := r.db.Preload("Account").Order("date desc").Find(&transactions).Error
	return transactions, err
}

func (r *financeRepository) GetCashTransactionById(id string) (*CashTransaction, error) {
	var transaction CashTransaction
	err := r.db.Preload("Account").Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *financeRepository) CreateCashTransaction(transaction *CashTransaction) error {
	return r.db.Create(transaction).Error
}

func (r *financeRepository) DeleteCashTransaction(id string) error {
	return r.db.Where("id = ?", id).Delete(&CashTransaction{}).Error
}

// General Journal
func (r *financeRepository) GetAllGeneralJournals() ([]GeneralJournal, error) {
	var journals []GeneralJournal
	err := r.db.Preload("Account").Order("date desc").Find(&journals).Error
	return journals, err
}

func (r *financeRepository) CreateGeneralJournal(journal *GeneralJournal) error {
	return r.db.Create(journal).Error
}

func (r *financeRepository) CreateJournalEntries(journals []GeneralJournal) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, j := range journals {
			if err := tx.Create(&j).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *financeRepository) VerifyPayment(req *dto.VerifyPaymentRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get the reservation
		var res map[string]interface{}
		if err := tx.Table("reservations").Where("id = ?", req.ReservationId).Take(&res).Error; err != nil {
			return err
		}

		totalPrice := res["total_price"].(int64)
		var deposit int64 = 0
		if res["deposit"] != nil {
			deposit = res["deposit"].(int64)
		}

		var resCode string
		if res["code"] != nil {
			resCode = res["code"].(string)
		} else {
			resCode = req.ReservationId
		}

		if req.Status == "confirmed" {
			// Update reservation status to confirmed
			if err := tx.Table("reservations").Where("id = ?", req.ReservationId).Update("status", "confirmed").Error; err != nil {
				return err
			}

			// Prepare journals
			var journals []GeneralJournal

			if deposit > 0 {
				journals = append(journals,
					GeneralJournal{
						Date:          time.Now(),
						Description:   "Penerimaan Titipan Deposit - " + resCode,
						AccountCode:   "1-100", // Kas
						Debit:         int(deposit),
						Credit:        0,
						ReservationId: &req.ReservationId,
					},
					GeneralJournal{
						Date:          time.Now(),
						Description:   "Kewajiban Titipan Tamu - " + resCode,
						AccountCode:   "2-100", // Kewajiban / Titipan Tamu
						Debit:         0,
						Credit:        int(deposit),
						ReservationId: &req.ReservationId,
					},
				)
			} else {
				// No deposit, assuming full payment upfront
				journals = append(journals,
					GeneralJournal{
						Date:          time.Now(),
						Description:   "Penerimaan Pembayaran Reservasi - " + resCode,
						AccountCode:   "1-100", // Kas
						Debit:         int(totalPrice),
						Credit:        0,
						ReservationId: &req.ReservationId,
					},
					GeneralJournal{
						Date:          time.Now(),
						Description:   "Pendapatan Sewa Kamar - " + resCode,
						AccountCode:   "4-100", // Pendapatan Sewa
						Debit:         0,
						Credit:        int(totalPrice),
						ReservationId: &req.ReservationId,
					},
				)
			}

			for i := range journals {
				journals[i].Id = uuid.New().String()
				if err := tx.Create(&journals[i]).Error; err != nil {
					return err
				}
			}
		} else if req.Status == "rejected" {
			// Update reservation status to rejected
			if err := tx.Table("reservations").Where("id = ?", req.ReservationId).Update("status", "rejected").Error; err != nil {
				return err
			}
		}

		return nil
	})
}
// Invoices
func (r *financeRepository) GetAllInvoices() ([]Invoice, error) {
	var invoices []Invoice
	err := r.db.Order("date desc").Find(&invoices).Error
	return invoices, err
}

func (r *financeRepository) GetInvoiceById(id string) (*Invoice, error) {
	var invoice Invoice
	err := r.db.Where("id = ?", id).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *financeRepository) CreateInvoice(invoice *Invoice) error {
	return r.db.Create(invoice).Error
}

func (r *financeRepository) UpdateInvoice(invoice *Invoice) error {
	return r.db.Save(invoice).Error
}

func (r *financeRepository) DeleteInvoice(id string) error {
	return r.db.Where("id = ?", id).Delete(&Invoice{}).Error
}
