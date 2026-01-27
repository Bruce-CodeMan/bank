package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all function to execute db quries and transaction.
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
}

// Store provides all functions to execute Postgres DB queries and transaction
type PGStore struct {
	pool *pgxpool.Pool
	*Queries
}

// NewStore creates a new Store
func NewPGStore(pool *pgxpool.Pool) *PGStore {
	return &PGStore{
		pool:    pool,
		Queries: New(pool), // New接收的是DBTX, 而pool实现了DBTX的 Exec/Query/QueryRow
	}
}

func NewStore(pool *pgxpool.Pool) Store {
	return NewPGStore(pool)
}

// execTx executes a function within a database transaction
func (s *PGStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := s.Queries.WithTx(tx)

	if err := fn(q); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult不是数据库表, 只是Go为了"把一次转账事务的结果打包返回"定义的结构体
// 这些字段的类型(Transfer/Entry/Account)一般都是sqlc根据你的表生成的Go struct(对应表的一行记录)
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, and account entries, and update account's balance with a single database transaction
func (s *PGStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		// 1) transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		// 2) Two entries
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID:  arg.FromAccountID,
			TransferID: result.Transfer.ID,
			Amount:     -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID:  arg.ToAccountID,
			TransferID: result.Transfer.ID,
			Amount:     arg.Amount,
		})
		if err != nil {
			return err
		}

		fromID := arg.FromAccountID
		toId := arg.ToAccountID

		if fromID < toId {
			result.FromAccount, err = s.addMoney(ctx, q, fromID, -arg.Amount)
			if err != nil {
				return err
			}
			result.ToAccount, err = s.addMoney(ctx, q, toId, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			// 注意这里先更新较小的
			result.ToAccount, err = s.addMoney(ctx, q, toId, arg.Amount)
			if err != nil {
				return err
			}
			result.FromAccount, err = s.addMoney(ctx, q, fromID, -arg.Amount)
			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func (s *PGStore) addMoney(ctx context.Context, q *Queries, id int64, amount int64) (Account, error) {
	return q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     id,
		Amount: amount,
	})
}
