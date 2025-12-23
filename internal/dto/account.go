// Package dto contains data transfer object (DTOs) for API requests and responses
package dto

import "github.com/google/uuid"

// CreateAccountRequest represents the request body for creating a new account.
// All fields are required and validated using gin's binding tags.
//	- PublicID is the UUID of the account(external identifier).
//	- Owner is the name of the account owner.
//	- Currency must be either "USD" or "EUR".
type CreateAccountRequest struct {
	PublicID uuid.UUID `json:"public_id" binding:"required"`
	Owner    string    `json:"owner" binding:"required"`
	Currency string    `json:"currency" binding:"required,oneof=USD EUR"`
}

// GetAccountRequest represents the URI parameter used to retrieve an existed account.
//
// In a dual-ID architecture, this struct is used to receive a public-facing UUID
// from the client instead of exposing the internal database ID. It helps to
// abstract and protect internal implementation details, improving security.
type GetAccountRequest struct {
	PublicID string `uri:"public_id" binding:"required"`
}

// ListAccountRequest represents the query paramaters used for paginated account listing.
//
// It includes the page number and the number of items per page.
//	- PageID must be at least 1.
//	- PageSize must be between 5 and 10.
type ListAccountRequest struct {
	PageID   int64 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}
