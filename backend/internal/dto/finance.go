package dto

type CreateCOARequest struct {
	Code string `json:"code" form:"code" binding:"required"`
	Name string `json:"name" form:"name" binding:"required"`
	Type string `json:"type" form:"type" binding:"required"`
}

type UpdateCOARequest struct {
	Name string `json:"name" form:"name"`
	Type string `json:"type" form:"type"`
}

type CreateCashTransactionRequest struct {
	Date        string `json:"date" form:"date" binding:"required"`
	Description string `json:"description" form:"description" binding:"required"`
	Type        string `json:"type" form:"type" binding:"required"`
	Amount      int    `json:"amount" form:"amount" binding:"required"`
	AccountId   string `json:"account_id" form:"account_id" binding:"required"`
}

type CreateGeneralJournalRequest struct {
	Date          string  `json:"date" form:"date" binding:"required"`
	Description   string  `json:"description" form:"description" binding:"required"`
	AccountCode   string  `json:"account_code" form:"account_code" binding:"required"`
	Debit         int     `json:"debit" form:"debit"`
	Credit        int     `json:"credit" form:"credit"`
	ReservationId *string `json:"reservation_id" form:"reservation_id"`
}

type VerifyPaymentRequest struct {
	ReservationId string `json:"reservation_id" binding:"required"`
	Status        string `json:"status" binding:"required"` // "confirmed" or "rejected"
}

type CreateInvoiceRequest struct {
	Type        string `json:"type" binding:"required"` // AP or AR
	PartnerName string `json:"partner_name" binding:"required"`
	Date        string `json:"date" binding:"required"`
	DueDate     string `json:"due_date" binding:"required"`
	Amount      int    `json:"amount" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type PayInvoiceRequest struct {
	AccountId string `json:"account_id" binding:"required"` // Account used for payment (e.g. Kas)
}
