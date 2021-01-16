package service

import (
	"banking/domain"
	"banking/dto"
	"banking/errs"
	"time"
)

type TransactionService interface {
	CreateTransaction(dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppError)
}

type DefaultTransactionService struct {
	repo domain.TransactionRepository
}

// CreateTransaction : Creates a new transaction and returns a transaction DTO
func (s DefaultTransactionService) CreateTransaction(req dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppError) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	t := domain.Transaction{
		AccoundId:       req.AccountID,
		TransactionDate: time.Now().Format("2006-01-02T15:04:05"),
		TransactionType: req.TransactionType,
		Amount:          req.Amount,
	}
	// todo: Validate account balance

	transaction, err := s.repo.NewTransaction(t)
	if err != nil {
		return nil, err
	}
	response := transaction.ToTransactionResponseDTO()
	return &response, nil
}

func NewTransactionService(repo domain.TransactionRepository) DefaultTransactionService {
	return DefaultTransactionService{repo}
}
