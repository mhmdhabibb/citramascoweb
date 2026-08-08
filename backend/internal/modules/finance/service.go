package finance

import (
	"citramascoweb-backend/internal/dto"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type FinanceServiceInterface interface {
	// COA
	GetAllCOA() ([]ChartOfAccount, error)
	GetCOAByCode(code string) (*ChartOfAccount, error)
	CreateCOA(req *dto.CreateCOARequest) (*ChartOfAccount, error)
	UpdateCOA(code string, req *dto.UpdateCOARequest) error
	DeleteCOA(code string) error
	ImportCOA(file multipart.File) error

	// Cash Transactions
	GetAllCashTransactions() ([]CashTransaction, error)
	GetCashTransactionById(id string) (*CashTransaction, error)
	CreateCashTransaction(req *dto.CreateCashTransactionRequest) (*CashTransaction, error)
	DeleteCashTransaction(id string) error

	// General Journal
	GetAllGeneralJournals() ([]GeneralJournal, error)
	CreateGeneralJournal(req *dto.CreateGeneralJournalRequest) (*GeneralJournal, error)
	CreateJournalEntries(entries []GeneralJournal) error

	// Verify Payment
	VerifyPayment(req *dto.VerifyPaymentRequest) error

	// Invoices
	GetAllInvoices() ([]Invoice, error)
	CreateInvoice(req *dto.CreateInvoiceRequest) (*Invoice, error)
	PayInvoice(id string, req *dto.PayInvoiceRequest) error
}

type financeService struct {
	repo FinanceRepositoryInterface
}

func NewFinanceService(repo FinanceRepositoryInterface) FinanceServiceInterface {
	return &financeService{repo: repo}
}

// COA
func (s *financeService) GetAllCOA() ([]ChartOfAccount, error) {
	return s.repo.GetAllCOA()
}

func (s *financeService) GetCOAByCode(code string) (*ChartOfAccount, error) {
	return s.repo.GetCOAByCode(code)
}

func (s *financeService) CreateCOA(req *dto.CreateCOARequest) (*ChartOfAccount, error) {
	coa := &ChartOfAccount{
		Code:      req.Code,
		Name:      req.Name,
		Type:      req.Type,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.repo.CreateCOA(coa)
	if err != nil {
		return nil, err
	}
	return coa, nil
}

func (s *financeService) UpdateCOA(code string, req *dto.UpdateCOARequest) error {
	coa := &ChartOfAccount{
		UpdatedAt: time.Now(),
	}
	if req.Name != "" {
		coa.Name = req.Name
	}
	if req.Type != "" {
		coa.Type = req.Type
	}
	return s.repo.UpdateCOA(coa, code)
}

func (s *financeService) DeleteCOA(code string) error {
	return s.repo.DeleteCOA(code)
}

func (s *financeService) ImportCOA(file multipart.File) error {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return err
	}
	defer f.Close()

	// Asumsikan sheet pertama
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		sheetName = "Sheet1"
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return err
	}

	var coas []ChartOfAccount
	// Loop rows (skip header if index 0 has headers)
	for i, row := range rows {
		// skip header
		if i == 0 {
			continue
		}

		if len(row) < 3 {
			continue // skip empty or incomplete row
		}

		code := row[0]
		name := row[1]
		accType := row[2]
		
		if code == "" {
			continue
		}

		var parentCode *string
		if len(row) >= 4 && row[3] != "" {
			pc := row[3]
			parentCode = &pc
		}

		coas = append(coas, ChartOfAccount{
			Code:       code,
			Name:       name,
			Type:       accType,
			ParentCode: parentCode,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})
	}

	if len(coas) > 0 {
		return s.repo.BulkUpsertCOA(coas)
	}

	return nil
}

// Cash Transactions
func (s *financeService) GetAllCashTransactions() ([]CashTransaction, error) {
	return s.repo.GetAllCashTransactions()
}

func (s *financeService) GetCashTransactionById(id string) (*CashTransaction, error) {
	return s.repo.GetCashTransactionById(id)
}

func (s *financeService) CreateCashTransaction(req *dto.CreateCashTransactionRequest) (*CashTransaction, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		date = time.Now() // Fallback
	}

	transaction := &CashTransaction{
		Id:          uuid.New().String(),
		Date:        date,
		Description: req.Description,
		Type:        req.Type,
		Amount:      req.Amount,
		AccountId:   req.AccountId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.repo.CreateCashTransaction(transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *financeService) DeleteCashTransaction(id string) error {
	return s.repo.DeleteCashTransaction(id)
}

// General Journal
func (s *financeService) GetAllGeneralJournals() ([]GeneralJournal, error) {
	return s.repo.GetAllGeneralJournals()
}

func (s *financeService) CreateGeneralJournal(req *dto.CreateGeneralJournalRequest) (*GeneralJournal, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		date = time.Now()
	}

	journal := &GeneralJournal{
		Id:            uuid.New().String(),
		Date:          date,
		Description:   req.Description,
		AccountCode:   req.AccountCode,
		Debit:         req.Debit,
		Credit:        req.Credit,
		ReservationId: req.ReservationId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = s.repo.CreateGeneralJournal(journal)
	if err != nil {
		return nil, err
	}
	return journal, nil
}

func (s *financeService) CreateJournalEntries(entries []GeneralJournal) error {
	// Add IDs and timestamps if not present
	for i := range entries {
		if entries[i].Id == "" {
			entries[i].Id = uuid.New().String()
		}
		if entries[i].CreatedAt.IsZero() {
			entries[i].CreatedAt = time.Now()
			entries[i].UpdatedAt = time.Now()
		}
	}
	return s.repo.CreateJournalEntries(entries)
}

func (s *financeService) VerifyPayment(req *dto.VerifyPaymentRequest) error {
	// Let's rely on gorm DB for reservation to avoid circular dependencies
	// Wait, we need the DB connection here, or pass through repository.
	// We'll update the repository to have VerifyPayment(req *dto.VerifyPaymentRequest) error
	// so the service will just call it.
	return s.repo.VerifyPayment(req)
}
// Invoices
func (s *financeService) GetAllInvoices() ([]Invoice, error) {
	return s.repo.GetAllInvoices()
}

func (s *financeService) CreateInvoice(req *dto.CreateInvoiceRequest) (*Invoice, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		date = time.Now()
	}
	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		dueDate = time.Now()
	}

	invoice := &Invoice{
		Id:            uuid.New().String(),
		InvoiceNumber: "INV-" + uuid.New().String()[:6],
		Type:          req.Type,
		PartnerName:   req.PartnerName,
		Date:          date,
		DueDate:       dueDate,
		Amount:        req.Amount,
		Status:        "Unpaid",
		Description:   req.Description,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = s.repo.CreateInvoice(invoice)
	if err != nil {
		return nil, err
	}
	return invoice, nil
}

func (s *financeService) PayInvoice(id string, req *dto.PayInvoiceRequest) error {
	invoice, err := s.repo.GetInvoiceById(id)
	if err != nil {
		return err
	}

	if invoice.Status == "Paid" {
		return nil
	}

	// Buat Jurnal
	var journals []GeneralJournal
	
	if invoice.Type == "AP" {
		// Hutang: Debit Hutang (Misal: 2-200), Kredit Kas
		journals = append(journals, 
			GeneralJournal{
				Date: time.Now(),
				Description: "Pelunasan Hutang - " + invoice.InvoiceNumber,
				AccountCode: "2-200", // Contoh kode Hutang Usaha, pastikan ada di COA
				Debit: invoice.Amount,
				Credit: 0,
			},
			GeneralJournal{
				Date: time.Now(),
				Description: "Pembayaran Hutang - " + invoice.InvoiceNumber,
				AccountCode: req.AccountId,
				Debit: 0,
				Credit: invoice.Amount,
			},
		)
	} else {
		// Piutang: Debit Kas, Kredit Piutang (Misal: 1-200)
		journals = append(journals, 
			GeneralJournal{
				Date: time.Now(),
				Description: "Penerimaan Piutang - " + invoice.InvoiceNumber,
				AccountCode: req.AccountId,
				Debit: invoice.Amount,
				Credit: 0,
			},
			GeneralJournal{
				Date: time.Now(),
				Description: "Pelunasan Piutang - " + invoice.InvoiceNumber,
				AccountCode: "1-200", // Contoh kode Piutang Usaha
				Debit: 0,
				Credit: invoice.Amount,
			},
		)
	}

	err = s.repo.CreateJournalEntries(journals)
	if err != nil {
		return err
	}

	invoice.Status = "Paid"
	invoice.UpdatedAt = time.Now()
	return s.repo.UpdateInvoice(invoice)
}
