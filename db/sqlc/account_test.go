package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"

	"github.com/BruceCompiler/bank/utils"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner: utils.RandomOwner(),
		PublicID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.PublicID, account.PublicID)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotZero(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.PublicID, account2.PublicID)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)

	account3, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	require.Equal(t, account2.ID, account3.ID)
	require.Equal(t, account2.PublicID, account3.PublicID)
	require.Equal(t, account2.Owner, account3.Owner)
	require.Equal(t, account2.Balance, account2.Balance)
	require.Equal(t, account2.Currency, account3.Currency)
	require.Equal(t, account2.CreatedAt, account3.CreatedAt)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.Zero(t, account2)
}
