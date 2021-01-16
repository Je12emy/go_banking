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
