// Package service contains business logic for the bank application.
// Services orchestrate operations between controllers and repositories
package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/BruceCompiler/bank/db/sqlc"
	"github.com/BruceCompiler/bank/internal/dto"
	"github.com/BruceCompiler/bank/internal/repository/postgres"
)

// AccountService handles business logic for account operations.
// It uses the Store to interact with the database.
type AccountService struct {
	store *postgres.Store
}

// NewAccountService creates a new AccountService with the given store.
func NewAccountService(s *postgres.Store) *AccountService {
	return &AccountService{store: s}
}

// CreateAccount creates a new bank account with the provided details.
// It initializes the account with a zero balance.
//
// Returns the created account or an error if the operation fails.
func (s *AccountService) CreateAccount(ctx context.Context, req dto.CreateAccountRequest) (db.Account, error) {
	return s.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		PublicID: pgtype.UUID{
			Bytes: req.PublicID,
			Valid: true,
		},
		Balance: 0,
	})
}
