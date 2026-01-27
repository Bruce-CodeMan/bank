package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/BruceCompiler/bank/db/sqlc"
	"github.com/BruceCompiler/bank/internal/dto"
	"github.com/BruceCompiler/bank/internal/token"
	"github.com/BruceCompiler/bank/utils"
	"github.com/BruceCompiler/bank/worker"
)

// UserService handles the business logic for user operations.
// It uses the Store to interact with the database.
type UserService struct {
	store           db.Store
	tokenMaker      token.Maker
	config          utils.Config
	taskDistributor worker.TaskDistributor
}

// NewUserService creates a new UserService with the given store.
func NewUserService(s db.Store, tokenMaker token.Maker, config utils.Config, taskDistributor worker.TaskDistributor) *UserService {
	return &UserService{
		store:           s,
		tokenMaker:      tokenMaker,
		config:          config,
		taskDistributor: taskDistributor,
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

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			PublicID: pgtype.UUID{
				Bytes: uuid.New(),
				Valid: true,
			},
			Username:       req.Username,
			HashedPassword: hashedPassword,
			FullName:       req.FullName,
			Email:          req.Email,
		},
		AfterCreate: func(createUserRow db.CreateUserRow) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: req.Username,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(5 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			return u.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	txResult, err := u.store.CreateUserTx(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// 23505 是 PostgreSQL 的 unique_violation 错误码
			if pgErr.Code == "23505" {
				// 可以进一步检查是哪个字段冲突
				if pgErr.ConstraintName == "users_username_key" {
					return db.CreateUserRow{}, errors.New("username already exists")
				}
				// 通用的唯一约束冲突
				return db.CreateUserRow{}, errors.New("user already exists")
			}
		}
		return db.CreateUserRow{}, err
	}
	return txResult.CreateUserRow, nil

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
		PublicID: pgtype.UUID{
			Bytes: refreshPayload.ID,
			Valid: true,
		},
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
