package domain

import (
	"banking/dto"
	"banking/errs"
)

type Account struct {
	AccountId   string  `db:"account_id"`
	CustomerId  string  `db:"customer_id"`
	OpeningDate string  `db:"opening_date"`
	AccountType string  `db:"account_type"`
	Amount      float64 `db:"amount"`
	Status      string  `db:"status"`
}

type AccountRepository interface {
	Save(Account) (*Account, *errs.AppError)
}

// ToNewAccountResponseDTO : Transforms a account object into a dto.account response
func (a Account) ToNewAccountResponseDTO() dto.NewAccountResponse {
	return dto.NewAccountResponse{AccountId: a.AccountId}
}
