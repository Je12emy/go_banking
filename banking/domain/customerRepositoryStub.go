package domain

// Stub = Mock Adapter

// CustomerRepositoryStub : Contains a dummy list of Customers
type CustomerRepositoryStub struct {
	customers []Customer
}

// FindAll : Returns the customers slice.
func (s CustomerRepositoryStub) FindAll() ([]Customer, error) {
	return s.customers, nil
}

// NewCustomerRepositoryStub : Helper function for creating a new list of customers and returns the customer repository stub
func NewCustomerRepositoryStub() CustomerRepositoryStub {
	customers := []Customer{
		{ID: "1001", Name: "Jotaro Kujo", City: "Okinawa", Zipcode: "30205", DateofBirth: "1970-01-01", Status: "1"},
		{ID: "1002", Name: "Jonathan Joestar", City: "England", Zipcode: "95457", DateofBirth: "1868-04-04", Status: "0"},
		{ID: "1003", Name: "Joseph Joestar", City: "England", Zipcode: "95825", DateofBirth: "1920-09-27", Status: "1"},
	}
	return CustomerRepositoryStub{customers}
}
