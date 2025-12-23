// Package service contains business logic for the bank application.
// Services orchestrate operations between controllers and repositories
package service

import (
	"context"

	"github.com/google/uuid"
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
//
// It stores the account in the database with a zero balance
// and returns the created account object
//
// Parameters:
//   - ctx: Standard context for request-scoped values and cancellation.
//   - req: A createAccountRequest DTO containing Owner, Currency and PublicID
//
// Returns:
//   - db.Account: The created account object.
//   - error: An error if the creation fails.
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

// GetAccount retrieves a bank account from the database using the provided UUID.
//
// It parses the UUID from the request, validates it and then queries the database.
// If the UUID is invalid or the database operation fails, it returns an error.
//
// Parameters:
//   - ctx: Standard context for request-scoped values and cancellation.
//   - req: A GetAccountRequest DTO containing the PublicID(UUID) of the account.
//
// Returns:
//   - db.Account: The account record if found.
//   - error: An error if the UUID is invalid or the account cannot be retrieved
func (s *AccountService) GetAccountByPublicID(ctx context.Context, req dto.GetAccountRequest) (db.Account, error) {

	uuidParsed, err := uuid.Parse(req.PublicID)
	if err != nil {
		return db.Account{}, err
	}
	return s.store.GetAccountByUUID(ctx, pgtype.UUID{
		Bytes: uuidParsed,
		Valid: true,
	})
}
