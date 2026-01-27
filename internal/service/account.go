// Package service contains business logic for the bank application.
// Services orchestrate operations between controllers and repositories
package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/BruceCompiler/bank/db/sqlc"
	"github.com/BruceCompiler/bank/internal/dto"
)

// AccountService handles business logic for account operations.
// It uses the Store to interact with the database.
type AccountService struct {
	store db.Store
}

// NewAccountService creates a new AccountService with the given store.
func NewAccountService(s db.Store) *AccountService {
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
	row, err := s.store.CreateAccount(ctx, db.CreateAccountParams{
		Currency: req.Owner,
		PublicID: pgtype.UUID{
			Bytes: req.PublicID,
			Valid: true,
		},
		Balance: 0,
	})
	if err != nil {
		return db.Account{}, err
	}

	account := db.Account{
		ID:            row.ID,
		PublicID:      row.PublicID,
		Balance:       row.Balance,
		Currency:      row.Currency,
		CreatedAt:     row.CreatedAt,
		PrimaryUserID: 0, // 或者后续补
	}

	return account, nil
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

// ListAccount retrieves a paginated list of bank accounts from the database.
//
// It calculates the SQL offset and limit based on the provided page ID and page size
// and then queries the database for the corresponding account records.
//
// Parameters:
//   - ctx: Standard context for request-scoped values and cancellation.
//   - req: A listAccountRequest DTO containing page_id and page_size.
//
// Returns:
//   - []db.Account: A slice of account records for the requested page.
//   - error: An error if the database query fails or parameter are invalid.
func (s *AccountService) ListAccount(ctx context.Context, req dto.ListAccountRequest) ([]db.Account, error) {
	arg := db.ListAccountsParams{
		Limit:  int32(req.PageSize),
		Offset: int32(req.PageID-1) * int32(req.PageSize),
	}
	return s.store.ListAccounts(ctx, arg)

}
