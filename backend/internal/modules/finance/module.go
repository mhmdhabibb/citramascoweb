package finance

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

type FinanceModule struct {
	Handler FinanceHandlerInterface
}

func InitModule(db *gorm.DB) *FinanceModule {
	repo := NewFinanceRepository(db)
	service := NewFinanceService(repo)
	handler := NewFinanceHandler(service)

	return &FinanceModule{
		Handler: handler,
	}
}

func (m *FinanceModule) FinanceRoutes(api *gin.RouterGroup) {
	financeGroup := api.Group("/finance")
	{
		// COA
		financeGroup.GET("/coa", m.Handler.GetAllCOA)
		financeGroup.POST("/coa", m.Handler.CreateCOA)
		financeGroup.POST("/coa/import", m.Handler.ImportCOA)
		financeGroup.PUT("/coa/:code", m.Handler.UpdateCOA)
		financeGroup.DELETE("/coa/:code", m.Handler.DeleteCOA)

		// Cash Transactions
		financeGroup.GET("/cash", m.Handler.GetAllCashTransactions)
		financeGroup.POST("/cash", m.Handler.CreateCashTransaction)
		financeGroup.DELETE("/cash/:id", m.Handler.DeleteCashTransaction)

		// General Journal
		financeGroup.GET("/journal", m.Handler.GetAllGeneralJournals)
		financeGroup.POST("/journal", m.Handler.CreateGeneralJournal)

		// Verify Payment
		financeGroup.POST("/verify-payment", m.Handler.VerifyPayment)

		// Invoices
		financeGroup.GET("/invoices", m.Handler.GetAllInvoices)
		financeGroup.POST("/invoices", m.Handler.CreateInvoice)
		financeGroup.POST("/invoices/:id/pay", m.Handler.PayInvoice)
	}
}
