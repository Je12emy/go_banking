package app

import (
	"banking/dto"
	"banking/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type AccountHandler struct {
	service service.AccountService
}

func (h AccountHandler) newAccount(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	customerId := vars["customer_id"]

	var request dto.NewAccountRequest
	// Decode the request into the "request" variable
	err := json.NewDecoder(r.Body).Decode(&request)
	// Check for errors in the decoding, BAD REQUEST
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		request.CustomerId = customerId
		// Create a new account by accesing the service
		account, appError := h.service.NewAccount(request)
		if appError != nil {
			// Write into the response if a error is found
			writeResponse(w, appError.Code, appError.Message)
		} else {
			// Write into the response the dto
			writeResponse(w, http.StatusCreated, account)
		}

	}
}

// /customers/2000/accounts/90720
func (h AccountHandler) MakeTransaction(w http.ResponseWriter, r *http.Request) {
	// get the account_id and customer_id from the URL
	vars := mux.Vars(r)
	accountId := vars["account_id"]
	customerId := vars["customer_id"]

	// decode incoming request
	var request dto.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {

		//build the request object
		request.AccountId = accountId
		request.CustomerId = customerId

		// make transaction
		account, appError := h.service.MakeTransaction(request)

		if appError != nil {
			writeResponse(w, appError.Code, appError.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, account)
		}
	}

}
