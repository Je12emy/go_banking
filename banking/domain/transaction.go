package domain

import (
	"banking/dto"
	"banking/errs"
)

type Transaction struct {
	ID              string
	AccoundId       string
	Amount          float64
	TransactionType string
	TransactionDate string
	Balance         float64
}

type TransactionRepository interface {
	NewTransaction(Transaction) (*Transaction, *errs.AppError)
}

func (t Transaction) ToTransactionResponseDTO() dto.TransactionResponse {
	return dto.TransactionResponse{
		TransactionId: t.ID,
		Balance:       t.Balance,
	}
}
