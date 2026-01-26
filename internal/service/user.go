package service

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/BruceCompiler/bank/db/sqlc"
	"github.com/BruceCompiler/bank/internal/dto"
	"github.com/BruceCompiler/bank/internal/repository/postgres"
	"github.com/BruceCompiler/bank/internal/token"
	"github.com/BruceCompiler/bank/utils"
)

// UserService handles the business logic for user operations.
// It uses the Store to interact with the database.
type UserService struct {
	store      postgres.Store
	tokenMaker token.Maker
	config     utils.Config
}

// NewUserService creates a new UserService with the given store.
func NewUserService(s postgres.Store, tokenMaker token.Maker, config utils.Config) *UserService {
	return &UserService{
		store:      s,
		tokenMaker: tokenMaker,
		config:     config,
	}
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

// Login
//
// Parameters:
//   - ctx: Standard context for request-scoped values and cancellation
//   - req: A LoginUserRequest DTO caontains Username and Password
func (u *UserService) Login(ctx context.Context, req dto.LoginUserRequest) (dto.LoginUserResponse, error) {
	user, err := u.store.GetUserByName(ctx, req.Username)
	if err != nil {
		return dto.LoginUserResponse{}, err
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return dto.LoginUserResponse{}, err
	}

	accessToken, accessPayload, err := u.tokenMaker.CreateToken(
		user.Username,
		u.config.AccessTokenDuration,
	)
	if err != nil {
		return dto.LoginUserResponse{}, err
	}

	refreshToken, refreshPayload, err := u.tokenMaker.CreateToken(
		user.Username,
		u.config.RefreshTokenDuration,
	)
	if err != nil {
		return dto.LoginUserResponse{}, err
	}

	_, err = u.store.CreateSession(ctx, db.CreateSessionParams{
		PublicID:     user.PublicID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt: pgtype.Timestamptz{
			Time:  refreshPayload.ExpiresAt.Time,
			Valid: true,
		},
	})
	if err != nil {
		return dto.LoginUserResponse{}, err
	}

	rsp := dto.LoginUserResponse{
		AccessToken: accessToken,
		User: dto.CreateUserResponse{
			ID:       strconv.FormatInt(user.ID, 10),
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
		},
		AccessTokenExpiresAt:  accessPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt.Time,
	}
	return rsp, nil
}
