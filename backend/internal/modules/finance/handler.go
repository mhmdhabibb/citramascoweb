package finance

import (
	"citramascoweb-backend/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FinanceHandlerInterface interface {
	// COA
	GetAllCOA(c *gin.Context)
	CreateCOA(c *gin.Context)
	UpdateCOA(c *gin.Context)
	DeleteCOA(c *gin.Context)
	ImportCOA(c *gin.Context)

	// Cash Transactions
	GetAllCashTransactions(c *gin.Context)
	CreateCashTransaction(c *gin.Context)
	DeleteCashTransaction(c *gin.Context)

	// General Journal
	GetAllGeneralJournals(c *gin.Context)
	CreateGeneralJournal(c *gin.Context)

	// Verify Payment
	VerifyPayment(c *gin.Context)

	// Invoices
	GetAllInvoices(c *gin.Context)
	CreateInvoice(c *gin.Context)
	PayInvoice(c *gin.Context)
}

type financeHandler struct {
	service FinanceServiceInterface
}

func NewFinanceHandler(service FinanceServiceInterface) FinanceHandlerInterface {
	return &financeHandler{service: service}
}

// COA
func (h *financeHandler) GetAllCOA(c *gin.Context) {
	coas, err := h.service.GetAllCOA()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "successfully retrieved coa list", "data": coas})
}

func (h *financeHandler) CreateCOA(c *gin.Context) {
	var req dto.CreateCOARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	coa, err := h.service.CreateCOA(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "successfully created coa", "data": coa})
}

func (h *financeHandler) UpdateCOA(c *gin.Context) {
	code := c.Param("code")
	var req dto.UpdateCOARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	err := h.service.UpdateCOA(code, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "successfully updated coa"})
}

func (h *financeHandler) DeleteCOA(c *gin.Context) {
	code := c.Param("code")
	err := h.service.DeleteCOA(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "successfully deleted coa"})
}

func (h *financeHandler) ImportCOA(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to get file: " + err.Error()})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to open file: " + err.Error()})
		return
	}
	defer file.Close()

	err = h.service.ImportCOA(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to import COA: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "successfully imported coa"})
}

// Cash Transactions
func (h *financeHandler) GetAllCashTransactions(c *gin.Context) {
	transactions, err := h.service.GetAllCashTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "successfully retrieved cash transactions", "data": transactions})
}

func (h *financeHandler) CreateCashTransaction(c *gin.Context) {
	var req dto.CreateCashTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	transaction, err := h.service.CreateCashTransaction(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "successfully created cash transaction", "data": transaction})
}

func (h *financeHandler) DeleteCashTransaction(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteCashTransaction(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "successfully deleted cash transaction"})
}

// General Journal
func (h *financeHandler) GetAllGeneralJournals(c *gin.Context) {
	journals, err := h.service.GetAllGeneralJournals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "successfully retrieved general journals", "data": journals})
}

func (h *financeHandler) CreateGeneralJournal(c *gin.Context) {
	var req dto.CreateGeneralJournalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	journal, err := h.service.CreateGeneralJournal(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "successfully created general journal", "data": journal})
}

func (h *financeHandler) VerifyPayment(c *gin.Context) {
	var req dto.VerifyPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	err := h.service.VerifyPayment(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "successfully verified payment"})
}
// Invoices
func (h *financeHandler) GetAllInvoices(c *gin.Context) {
	invoices, err := h.service.GetAllInvoices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": invoices})
}

func (h *financeHandler) CreateInvoice(c *gin.Context) {
	var req dto.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	invoice, err := h.service.CreateInvoice(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": invoice})
}

func (h *financeHandler) PayInvoice(c *gin.Context) {
	id := c.Param("id")
	var req dto.PayInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	
	err := h.service.PayInvoice(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Invoice paid successfully"})
}
