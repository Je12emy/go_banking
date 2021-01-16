package service

import (
	"banking/domain"
	"banking/dto"
	"banking/errs"
)

// CustomerService : Primary Port
type CustomerService interface {
	GetAllCustomer(status string) ([]dto.CustomerResponse, *errs.AppError)
	GetCustomer(id string) (*dto.CustomerResponse, *errs.AppError)
}

// DefaultCustomerService : Named Default because the could be more implementations
type DefaultCustomerService struct {
	// Has a dependency on the CustomerRepository
	repo domain.CustomerRepository
}

// GetAllCustomer : Helper function which returns a slice of customer dto, status can be passed
func (s DefaultCustomerService) GetAllCustomer(status string) ([]dto.CustomerResponse, *errs.AppError) {
	var query string

	if status == "active" {
		query = "1"
	} else if status == "inactive" {
		query = "0"
	} else {
		query = ""
	}

	var customersResponse []dto.CustomerResponse

	customers, err := s.repo.FindAll(query)
	if err != nil {
		return nil, err
	}

	for _, c := range customers {
		customersResponse = append(customersResponse, c.ToDTO())
	}
	return customersResponse, nil
}

// GetCustomer : Helper function which returns a customer DTO
func (s DefaultCustomerService) GetCustomer(id string) (*dto.CustomerResponse, *errs.AppError) {
	c, err := s.repo.ById(id)
	if err != nil {
		return nil, err
	}
	response := c.ToDTO()
	return &response, nil
}

// NewCustomerService : Return a new DefaultCustomerService which takes a CustomerRepository.
func NewCustomerService(repository domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repository}
}
