// Package dto contains data transfer object (DTOs) for API requests and responses
package dto

import "github.com/google/uuid"

// CreateAccountRequest represents the request body for creating a new account.
// All fields are required and validated using gin's binding tags.
type CreateAccountRequest struct {
	PublicID uuid.UUID `json:"public_id" binding:"required"`
	Owner    string    `json:"owner" binding:"required"`
	Currency string    `json:"currency" binding:"required,oneof=USD EUR"`
}
