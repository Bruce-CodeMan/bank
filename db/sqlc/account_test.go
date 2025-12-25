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
	user := createRandomUser(t)
	arg := CreateAccountParams{
		PublicID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Balance:       utils.RandomMoney(),
		Currency:      utils.RandomCurrency(),
		PrimaryUserID: user.ID,
	}
	row, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, row)

	require.Equal(t, arg.PublicID, row.PublicID)
	require.Equal(t, arg.Balance, row.Balance)
	require.Equal(t, arg.Currency, row.Currency)

	require.NotZero(t, row.CreatedAt)

	account := Account{
		ID:            row.ID,
		PublicID:      row.PublicID,
		Balance:       row.Balance,
		Currency:      row.Currency,
		CreatedAt:     row.CreatedAt,
		PrimaryUserID: user.ID,
	}

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotZero(t, account2)

	require.Equal(t, account1.ID, account2.ID)
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

	account3, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)

	require.Equal(t, account2.ID, account3.ID)
	require.Equal(t, account2.PublicID, account3.PublicID)
	require.Equal(t, account2.Balance, account2.Balance)
	require.Equal(t, account2.Currency, account3.Currency)
	require.Equal(t, account2.CreatedAt, account3.CreatedAt)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.Zero(t, account2)
}
