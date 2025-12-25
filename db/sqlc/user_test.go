package db

import (
	"context"
	"testing"

	"github.com/BruceCompiler/bank/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		PublicID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Username:       utils.RandomOwner(),
		HashedPassword: utils.RandomString(6),
		FullName:       utils.RandomString(6),
		Email:          utils.RandomEmail(),
	}

	row, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, row)

	user := User{
		ID:                row.ID,
		PublicID:          row.PublicID,
		Username:          row.Username,
		HashedPassword:    row.HashedPassword,
		FullName:          row.FullName,
		Email:             row.Email,
		PasswordChangedAt: row.PasswordChangedAt,
		CreatedAt:         row.CreatedAt,
	}

	return user
}
