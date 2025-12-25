// Package service contains business logic for the transfer
// Services orchestrate operations between controllers and repositories
package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/BruceCompiler/bank/internal/dto"
	"github.com/BruceCompiler/bank/internal/repository/postgres"
)

// TransferService handles business logic for transfer operations.
// It uses the store to interact with the database.
type TransferService struct {
	store postgres.Store
}

// NewTransferService creates a new TransferService with the given store
func NewTransferService(s postgres.Store) *TransferService {
	return &TransferService{store: s}
}

// CreateTransfer creates a new transfer with the given details
//
// Parameters:
//   - ctx: Standard context for request-scoped values and cancellation
//   - req: A CreateTransferRequest DTO contains FromAccountPublicID, ToAccountPublicID, Amount, Currency
//
// Returns
//   - postgres.TransferResult: The created transfer
//   - error: An error if the creation fails
func (ts *TransferService) CreateTransfer(ctx context.Context, req dto.CreateTransferRequest) (postgres.TransferTxResult, error) {
	arg := postgres.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	err := ts.validAccount(ctx, req.FromAccountID, req.Currency)
	if err != nil {
		return postgres.TransferTxResult{}, err
	}
	err = ts.validAccount(ctx, req.ToAccountID, req.Currency)
	if err != nil {
		return postgres.TransferTxResult{}, err
	}

	return ts.store.TransferTx(ctx, arg)
}

func (ts *TransferService) validAccount(ctx context.Context, accountID int64, currency string) error {
	account, err := ts.store.GetAccountById(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}

	if account.Currency != currency {
		return errors.New("account currency mismatch")
	}
	return nil
}
