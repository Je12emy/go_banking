package app

import (
	"banking/dto"
	"banking/errs"
	"banking/mocks/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

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

func Test_should_return_customer_with_status_code_200(t *testing.T) {
	// Arrange
	teardown := setup(t)
	defer teardown()

	// Insert the getAllCustomers method into the mock implementation
	dummyCustomers := []dto.CustomerResponse{
		{ID: "1001", Name: "Jotaro Kujo", City: "Okinawa", Zipcode: "30205", DateofBirth: "1970-01-01", Status: "1"},
		{ID: "1002", Name: "Jonathan Joestar", City: "England", Zipcode: "95457", DateofBirth: "1868-04-04", Status: "0"},
		{ID: "1003", Name: "Joseph Joestar", City: "England", Zipcode: "95825", DateofBirth: "1920-09-27", Status: "1"},
	}
	// Return the dummy list of customers when the GetAllCustomers method is called
	mockService.EXPECT().GetAllCustomer("").Return(dummyCustomers, nil)
	// Prepare a query
	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	// Assert
	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing the status code.")
	}
}

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
