// Package dto contains data transfer object (DTO) for transfer API request and response
package dto

// TransferRequest represents the request body for creating a transfer info
type CreateTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required"`
	ToAccountID   int64  `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,oneof=USD EUR"`
}
