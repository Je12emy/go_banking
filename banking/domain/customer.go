package domain

import (
	"banking/dto"
	"banking/errs"
)

// Customer : Bussiness Object
type Customer struct {
	ID          string `db:"customer_id"`
	Name        string
	City        string
	Zipcode     string
	DateofBirth string `db:"date_of_birth"`
	Status      string
}

// ToDTO : Returns a customer domain object as a DTO customer object
func (c Customer) ToDTO() dto.CustomerResponse {
	return dto.CustomerResponse{
		ID:          c.ID,
		Name:        c.Name,
		City:        c.City,
		Zipcode:     c.Zipcode,
		DateofBirth: c.DateofBirth,
		Status:      c.AsStatusText(),
	}
}

// AsStatusText : Return the number for status as a new string
func (c Customer) AsStatusText() string {
	statusText := "Active"
	if c.Status == "0" {
		statusText = "Inactive"
	}
	return statusText
}

// CustomerRepository : Secondary Port, boundary of the domain
type CustomerRepository interface {
	FindAll(status string) ([]Customer, *errs.AppError)
	ById(string) (*Customer, *errs.AppError)
}
