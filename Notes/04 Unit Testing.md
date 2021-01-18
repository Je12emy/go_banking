# Unit Testing
## Unit Testing: State based test
Unit testing consists of testing the most basic pieces of code in isolation, in our application we can easily test dto validation. Test follow a "AAA" format: Arrange, Act and Assert.

Go provides us with a testing package, though it lacks assert functionality it is enough to write tests with conditionals.

```go
// function name should be as descriptive as possible
func Test_should_return_error_when_transaction_type_is_not_deposit_or_withdrawal(t *testing.T) {
	// Arrange
	request := TransactionRequest{TransactionType: "invalid type"}
	// Act
	err := request.Validate()
	// Assert
	if err.Message != "Transaction type can only be deposit or withdrawal" {
		t.Error("Invalid error message was thrown when validating error message.")
	}
	if err.Code != http.StatusUnprocessableEntity {
		t.Error("Invalid error code was thrown when validating error code.")
	}
}
```

To run this test, we are able to execute it in our editor/IDE or we could move into this directory and run `go test` or `go test -v`

```
> $ go test -v                                                                                                                   
=== RUN   Test_should_return_error_when_transaction_type_is_not_deposit_or_withdrawal
--- PASS: Test_should_return_errorwhen_transaction_type_is_not_deposit_or_withdrawal (0.00s)
```

## Testing Routes: Using Mocks

To test our route handlers, we need to mock the services since they should be tested in issolation. To generate a mock service we will be using Go' [mock tool](https://github.com/golang/mock) which will generate a mock implementation for us.

To generate a moch for our `CustomerService` interface we need to add a decorator, which will be used by to mock tool to create the mock for us.

```go
//go:generate mockgen -destination=../mocks/service/mockCustomerService.go -package=service banking/service CustomerService
type CustomerService interface {
	GetAllCustomer(status string) ([]dto.CustomerResponse, *errs.AppError)
	GetCustomer(id string) (*dto.CustomerResponse, *errs.AppError)
}
```

In our terminal run `go generate ./...` this will create the mock in the specified location.

Let's focus on the "Arrange" part of our test, in here we need to inject into the `customerHandlers` the mock implementation of `customerService`, this mock implementation needs a controller which allows us to test our interactions.

Into this mock implementation we will be codding what should be returned, which is done with the `EXPECT()` function which allows us to access the `GetAllCustomers()` and `RETURN()` the expected result which is a dummy list of customers.

Finally we'll create a new router and a new query to access the resource, this code is a bit large due to the set up, which we will solve later on.

```go
func Test_should_return_customer_with_status_code_200(t *testing.T) {
	// Arrange
	// Create a new controller for the mock service
	ctrl := gomock.NewController(t)

	// Finish allows us to test our interactions
	// defer will call ctrl.Finish() just before the functions returns
	defer ctrl.Finish()
	mockService := service.NewMockCustomerService(ctrl)

	// Insert the getAllCustomers method into the mock implementation
	dummyCustomers := []domain.Customer{
		{ID: "1001", Name: "Jotaro Kujo", City: "Okinawa", Zipcode: "30205", DateofBirth: "1970-01-01", Status: "1"},
		{ID: "1002", Name: "Jonathan Joestar", City: "England", Zipcode: "95457", DateofBirth: "1868-04-04", Status: "0"},
		{ID: "1003", Name: "Joseph Joestar", City: "England", Zipcode: "95825", DateofBirth: "1920-09-27", Status: "1"},
	}
	// Return the dummy list of customers when the GetAllCustomers method is called
	mockService.EXPECT().GetAllCustomer("").Return(dummyCustomers, nil)
	ch := CustomerHandlers{mockService}
	// Set up a new router and route handler with the mock implementation
	router := mux.NewRouter()
	router.HandleFunc("/customers", ch.getAllCustomers)
	// Prepare a query
	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	// Act
	// more code ...
	
	// Assert
	// more code ...
}
```

In the "Act" section, we will be creating a new test response writter which is provided by the `httptest` library and we will be asserting into this writter.

```go
func Test_should_return_customer_with_status_code_200(t *testing.T) {
	// Arrange
	// more code ...

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	
	// Assert
	// more code ...
}
```

We can easily assert into this `recorder`, let's assert it's status code response.

```go
func Test_should_return_customer_with_status_code_200(t *testing.T) {
	// Arrange
	// more code ...

	// Act
	// more code ...
	
	// Assert
	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing the status code.")
	}
}
```

If we test for an error, the should should be pretty similar, instead of returning a list of customers we only need to return nil and the new error.

```go
func Test_shold_return_status_code_500_with_error_message(t *testing.T) {
	// Arrange

	// Create a new controller for the mock service
	ctrl := gomock.NewController(t)

	// Finish allows us to test our interactions
	// defer will call ctrl.Finish() just before the functions returns
	defer ctrl.Finish()
	mockService := service.NewMockCustomerService(ctrl)

	// Return the dummy list of customers when the GetAllCustomers method is called
	mockService.EXPECT().GetAllCustomer("").Return(nil, errs.NewUnexpectedError("test database error"))
	ch := CustomerHandlers{mockService}
	// Set up a new router and route handler with the mock implementation
	router := mux.NewRouter()
	router.HandleFunc("/customers", ch.getAllCustomers)
	// Prepare a query
	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing the status code.")
	}
}
```

Before we finish let's refractor the setup for both test, and make it reusable between them.

```go
var ch CustomerHandlers
var router *mux.Router
var mockService *service.MockCustomerService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockService = service.NewMockCustomerService(ctrl)
	ch = CustomerHandlers{mockService}
	router = mux.NewRouter()
	router.HandleFunc("/customers", ch.getAllCustomers)

	return func() {
		router = nil
		defer ctrl.Finish()
	}
}
```

This will allow our test's code to be much smaller.

```go
func Test_shold_return_status_code_500_with_error_message(t *testing.T) {
	// Arrange
	teardown := setup(t)
	defer teardown()

	// Return the dummy list of customers when the GetAllCustomers method is called
	mockService.EXPECT().GetAllCustomer("").Return(nil, errs.NewUnexpectedError("test database error"))
	// Prepare a query
	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing the status code.")
	}
}
```

## Unit Testing: Using Mocks

The same logic is used to test our services, in this case we'll need to mock the domain

```go
//go:generate mockgen -destination=../mocks/domain/mockAccountRepository.go -package=domain banking/domain AccountRepository
type AccountRepository interface {
	Save(Account) (*Account, *errs.AppError)
	SaveTransaction(transaction Transaction) (*Transaction, *errs.AppError)
	FindBy(accountId string) (*Account, *errs.AppError)
}
```

In our test, we will be injecting into the `Save()` method which needs a account domain object. since we are using the `banking/mocks/domain` and `banking/domain` we'll need to create a "alias" for the real domain in order to create the needed domain object.

```go
import (
	realdomain "banking/domain"
	"banking/mocks/domain"
)

var mockRepo *domain.MockAccountRepository
var ctrl gomock.Controller
var service AccountService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockRepo = domain.NewMockAccountRepository(ctrl)
	service = NewAccountService(mockRepo)
	return func() {
		service = nil
		defer ctrl.Finish()
	}
}
```

The test will asert upon the injected `save` function.

```go
func Test_should_return_an_error_from_the_server_side_if_the_new_account_cannot_be_saved(t *testing.T) {
	// Arrange
	teardown := setup(t)
	defer teardown()

	req := dto.NewAccountRequest{
		CustomerId:  "100",
		AccountType: "savings",
		Amount:      6000,
	}
	account := realdomain.Account{
		CustomerId:  req.CustomerId,
		OpeningDate: time.Now().Format("2006-01-02T15:04:05"),
		AccountType: req.AccountType,
		Amount:      req.Amount,
		Status:      "1",
	}
	mockRepo.EXPECT().Save(account).Return(nil, errs.NewUnexpectedError("Unexpected database error"))
	// Act
	_, appError := service.NewAccount(req)

	// Assert
	if appError == nil {
		t.Error("Test failed while validating error for new account")
	}
}
```