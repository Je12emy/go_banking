package domain

import (
	"banking/errs"
	"banking/logger"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type AccountRepositoryDB struct {
	client *sqlx.DB
}

// Save : Take a customer and create a new account record
func (d AccountRepositoryDB) Save(a Account) (*Account, *errs.AppError) {
	sqlInsert := "INSERT INTO accounts (customer_id, opening_date, account_type, amount, status) VALUES(?, ?, ?, ?, ?);"
	result, err := d.client.Exec(sqlInsert, a.CustomerId, a.OpeningDate, a.AccountType, a.Amount, a.Status)
	if err != nil {
		logger.Error("Error while creating new account: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database")
	}
	// Returns the id
	id, err := result.LastInsertId()

	if err != nil {
		logger.Error("Error while retrieving last inserted id: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database")
	}
	// Save the account id into the domain object
	a.AccountId = strconv.FormatInt(id, 10)
	return &a, nil
}

// NewAccountRepositoryDB : Returns the account repository
func NewAccountRepositoryDB(dbClient *sqlx.DB) AccountRepositoryDB {
	return AccountRepositoryDB{dbClient}
}
