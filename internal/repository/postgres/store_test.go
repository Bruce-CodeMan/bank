package postgres

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	db "github.com/BruceCompiler/bank/db/sqlc"
	"github.com/BruceCompiler/bank/utils"
)

const (
	dbSource = "postgresql://root:secret@localhost:5432/bank?sslmode=disable"
)

func createRandomAccount(t *testing.T, store *Store) db.Account {
	arg := db.CreateAccountParams{
		Owner:    utils.RandomOwner(),
		PublicID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
	account, err := store.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	return account
}

func TestTransferTX(t *testing.T) {

	pool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect db: ", err)
	}
	defer pool.Close()
	store := NewStore(pool)
	ctx := context.Background()

	account1 := createRandomAccount(t, store)
	account2 := createRandomAccount(t, store)

	fmt.Printf(">> before: account1.Balance: %d, account2.Balance: %d\n",
		account1.Balance, account2.Balance)

	n := 5
	amount := int64(10)

	// 用缓冲，避免 goroutine 因发送阻塞导致测试不稳定
	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)

	for i := 0; i < n; i++ {
		from := account1.ID
		to := account2.ID
		if i%2 == 1 {
			from = account2.ID
			to = account1.ID
		}

		go func(from, to int64) {
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: from,
				ToAccountID:   to,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}(from, to)
	}

	// check each tx result
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotZero(t, result)

		// ---- check transfer ----
		transfer := result.Transfer
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		require.Equal(t, amount, transfer.Amount)

		fromID := transfer.FromAccountID
		toID := transfer.ToAccountID

		// 方向必须是 A->B 或 B->A
		require.True(t,
			(fromID == account1.ID && toID == account2.ID) ||
				(fromID == account2.ID && toID == account1.ID),
		)

		_, err = store.GetTransfer(ctx, transfer.ID)
		require.NoError(t, err)

		// ---- check entries (按本次 transfer 的 from/to 校验) ----
		fromEntry := result.FromEntry
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		require.Equal(t, fromID, fromEntry.AccountID)
		require.Equal(t, transfer.ID, fromEntry.TransferID)
		require.Equal(t, -amount, fromEntry.Amount)

		_, err = store.GetEntry(ctx, fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		require.Equal(t, toID, toEntry.AccountID)
		require.Equal(t, transfer.ID, toEntry.TransferID)
		require.Equal(t, amount, toEntry.Amount)

		_, err = store.GetEntry(ctx, toEntry.ID)
		require.NoError(t, err)

		// ---- check accounts returned (按本次 transfer 的 from/to 校验) ----
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toID, toAccount.ID)
	}

	// ---- check final updated balance ----
	updatedAccount1, err := store.GetAccountById(ctx, account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccountById(ctx, account2.ID)
	require.NoError(t, err)

	fmt.Printf(">> after: account1.Balance: %d, account2.Balance: %d\n",
		updatedAccount1.Balance, updatedAccount2.Balance)

	// 双向净变化
	// 偶数 i: A->B 次数 = ceil(n/2) = (n+1)/2
	// 奇数 i: B->A 次数 = floor(n/2) = n/2
	a2b := (n + 1) / 2
	b2a := n / 2

	expectedA := account1.Balance - int64(a2b)*amount + int64(b2a)*amount
	expectedB := account2.Balance + int64(a2b)*amount - int64(b2a)*amount

	require.Equal(t, expectedA, updatedAccount1.Balance)
	require.Equal(t, expectedB, updatedAccount2.Balance)
}
