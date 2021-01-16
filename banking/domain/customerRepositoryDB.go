package domain

import (
	"banking/errs"
	"banking/logger"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// CustomerRepositoryDb : Repository for querying the database
type CustomerRepositoryDb struct {
	client *sqlx.DB
}

// FindAll : Queries the database for the list of customers.
func (d CustomerRepositoryDb) FindAll(status string) ([]Customer, *errs.AppError) {
	var err error
	var findAllSQL string
	// Create a Customers slice
	customers := make([]Customer, 0)
	if status != "" {
		findAllSQL = "SELECT customer_id, name, city, zipcode, date_of_birth, status FROM customers WHERE status=?;"
		err = d.client.Select(&customers, findAllSQL, status)

	} else {

		findAllSQL = "SELECT customer_id, name, city, zipcode, date_of_birth, status FROM customers;"
		err = d.client.Select(&customers, findAllSQL)
	}
	if err != nil {
		logger.Error("Error while querying customer table " + err.Error())
		return nil, errs.NewUnexpectedError("Error while querying customer table " + err.Error())
	}
	return customers, nil
}

// ById : Returns a single customer by his id.
func (d CustomerRepositoryDb) ById(id string) (*Customer, *errs.AppError) {
	// Create the query
	customerSQL := "SELECT customer_id, name, city, zipcode, date_of_birth, status FROM customers WHERE customer_id = ?;"
	var c Customer
	err := d.client.Get(&c, customerSQL, id)
	if err != nil {
		// If no customer is found at all
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Customer not found")
		} else {
			logger.Error("Error while scanning customer " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}
	return &c, nil
}

// NewCustomerRepositoryDb : Creates a sql client and returns the CustomerRepositoryDB
func NewCustomerRepositoryDb(dbClient *sqlx.DB) CustomerRepositoryDb {
	return CustomerRepositoryDb{dbClient}
}
