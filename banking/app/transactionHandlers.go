package app

import (
	"banking/dto"
	"banking/service"
	"encoding/json"
	"net/http"
)

type TransactionHandlers struct {
	service service.TransactionService
}

func (th TransactionHandlers) newTransaction(w http.ResponseWriter, r *http.Request) {
	var request dto.TransactionRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {

		transaction, err := th.service.CreateTransaction(request)
		if err != nil {
			writeResponse(w, err.Code, err.AsMessage())
		} else {
			writeResponse(w, http.StatusCreated, transaction)
		}
	}
}
