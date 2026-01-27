package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(createUserRow CreateUserRow) error
}

type CreateUserTxResult struct {
	CreateUserRow CreateUserRow
}

func (store *PGStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.CreateUserRow, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}
		return arg.AfterCreate(result.CreateUserRow)

	})

	return result, err
}
