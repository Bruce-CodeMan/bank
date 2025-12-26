package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"

	"github.com/BruceCompiler/bank/utils"
)

func createRandomUser(t *testing.T) CreateUserRow {
	hashedPassword, err := utils.HashPassword(utils.RandomString(6))
	arg := CreateUserParams{
		PublicID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomString(6),
		Email:          utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByPublicID(context.Background(), user1.PublicID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)

}
