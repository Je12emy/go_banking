package dto

import (
	"banking/errs"
	"strings"
)

type TransactionRequest struct {
	AccountID       string  `json:"account_id"`
	TransactionType string  `json:"transaction_type"`
	Amount          float64 `json:amount`
}

// Validate : Validates the request if the amount and type are correct
func (tr TransactionRequest) Validate() *errs.AppError {
	if tr.Amount <= 0 {
		return errs.NewValidationError("Transaction amount must be positive.")
	}

	if strings.ToLower(tr.TransactionType) != "withdrawal" && strings.ToLower(tr.TransactionType) != "deposit" {
		return errs.NewValidationError("Invalid transaction type")
	}
	return nil
}
