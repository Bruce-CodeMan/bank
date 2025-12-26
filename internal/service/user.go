package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/BruceCompiler/bank/db/sqlc"
	"github.com/BruceCompiler/bank/internal/dto"
	"github.com/BruceCompiler/bank/internal/repository/postgres"
	"github.com/BruceCompiler/bank/utils"
)

// UserService handles the business logic for user operations.
// It uses the Store to interact with the database.
type UserService struct {
	store postgres.Store
}

// NewUserService creates a new UserService with the given store.
func NewUserService(s postgres.Store) *UserService {
	return &UserService{store: s}
}

// CreateUser creates a new user with the given details
//
// Parameters:
//   - ctx: Standard context for request-scoped values and cancellation
//   - req: A CreateUserRequest DTO contains Username, Password, Email, FullName
//
// Returns:
//   - db.CreateUserRow: The created user
//   - error: An error if the creation fails
func (u *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (db.CreateUserRow, error) {

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return db.CreateUserRow{}, err
	}

	return u.store.CreateUser(ctx, db.CreateUserParams{
		PublicID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	})
}
