package dto

import (
	"net/http"
	"testing"
)

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

func Test_should_return_error_when_amount_is_less_than_zeo(t *testing.T) {
	// Arrange
	request := TransactionRequest{TransactionType: DEPOSIT, Amount: -100}
	// Act
	err := request.Validate()
	// Assert
	if err.Message != "Amount cannot be less than zero" {
		t.Error("Invalid error message was thrown when validating transaction amount.")
	}
	if err.Code != http.StatusUnprocessableEntity {
		t.Error("Invalid error code was thrown when validating transaction amount")
	}
}
