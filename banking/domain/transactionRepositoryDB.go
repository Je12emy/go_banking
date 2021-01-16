package domain

import (
	"banking/dto"
	"banking/errs"
	"banking/logger"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type TransactionRepositoryDB struct {
	client *sqlx.DB
}

// NewTransaction : Creates a new transaction for a given customer
func (tr TransactionRepositoryDB) NewTransaction(t Transaction) (*Transaction, *errs.AppError) {
	var appErr *errs.AppError

	if strings.ToLower(t.TransactionType) == "withdrawal" {
		isValid, appErr := tr.validateBalance(t.Amount, t.AccoundId, t.TransactionType)

		if !isValid {
			return nil, appErr
		}
	}

	sqlInsert := "INSERT INTO transactions (account_id, transaction_type, amount, transaction_date) VALUES (?, ?, ?, ?);"
	result, err := tr.client.Exec(sqlInsert, t.AccoundId, t.TransactionType, t.Amount, t.TransactionDate)
	if err != nil {
		logger.Error("Error while creating new transaction: " + err.Error())
		return nil, errs.NewUnexpectedError("Error while creating new transaction: " + err.Error())
	}

	id, err := result.LastInsertId()

	if err != nil {
		logger.Error("Error while retrieving the last inserted id: " + err.Error())
		return nil, errs.NewUnexpectedError("Error while retrieving the last inserted id: " + err.Error())
	}

	t.ID = strconv.FormatInt(id, 10)

	// Update the account balance
	a, appErr := tr.getAccountInfo(t.AccoundId)
	if appErr != nil {
		return nil, appErr
	}
	t.Amount = tr.changeAmount(t.Amount, t.TransactionType)
	transac, appErr := tr.updateBalance(*a, t)
	if err != nil {
		return nil, appErr
	}
	return transac, nil
}

func (tr TransactionRepositoryDB) validateBalance(amount float64, id string, transactionType string) (bool, *errs.AppError) {
	accountSelect := `SELECT IF ( ? <= amount, "TRUE", "FALSE") AS hasBalance FROM accounts WHERE account_id =?`
	var isValid dto.HasBalance
	err := tr.client.Get(&isValid, accountSelect, amount, id)

	if err != nil {
		logger.Error("Error while validating account balance " + err.Error())
		return false, errs.NewUnexpectedError("Error while validating account balance " + err.Error())
	}
	if isValid.BalanceIsValid {
		return true, nil
	}
	return false, errs.NewValidationError("Not enough funds for withdrawal")
}

func (tr TransactionRepositoryDB) changeAmount(amount float64, transactionType string) float64 {
	if strings.ToLower(transactionType) == "deposit" {
		return amount
	}
	return amount * -1
}

// getAccountInfo : Returns basic account information
func (tr TransactionRepositoryDB) getAccountInfo(id string) (*Account, *errs.AppError) {
	accountSelect := "SELECT account_id, amount from accounts WHERE account_id =?;"
	var a Account
	err := tr.client.Get(&a, accountSelect, id)

	if err != nil {
		logger.Error("Error while retrieving account: " + err.Error())
		return nil, errs.NewUnexpectedError("Error while retrieving account: " + err.Error())
	}
	return &a, nil
}

// updateBalance : Updates the account's balance
func (tr TransactionRepositoryDB) updateBalance(a Account, t Transaction) (*Transaction, *errs.AppError) {
	sqlUpdate := "UPDATE accounts SET AMOUNT = AMOUNT + ? WHERE account_id = ?"
	_, err := tr.client.Exec(sqlUpdate, t.Amount, a.AccountId)
	if err != nil {
		logger.Error("Error while updating account balance: " + err.Error())
		return nil, errs.NewUnexpectedError("Error while updating account balance: " + err.Error())
	}

	balance, errs := tr.getBalance(a.AccountId)
	if errs != nil {
		return nil, errs
	}
	transac := Transaction{
		ID:      t.ID,
		Balance: *balance,
	}

	return &transac, nil
}

// getBalance : Returns the account's current balance
func (tr TransactionRepositoryDB) getBalance(id string) (*float64, *errs.AppError) {
	sqlSelect := "SELECT amount from accounts WHERE account_id = ?"
	var a Account
	err := tr.client.Get(&a, sqlSelect, id)
	if err != nil {
		logger.Error("Error while retrieving balance: " + err.Error())
		return nil, errs.NewUnexpectedError("Error while retrieving balance: " + err.Error())
	}
	return &a.Amount, nil
}

// NewTransactionRepositoryDB : Creates a new TransactionRepositoryDB
func NewTransactionRepositoryDB(dbClient *sqlx.DB) TransactionRepositoryDB {
	return TransactionRepositoryDB{dbClient}
}
