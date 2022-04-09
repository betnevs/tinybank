package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

var txKey = struct{}{}

type Store interface {
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	Querier
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	TransferID    int64 `json:"transfer_id"`
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	FromEntryID   int64 `json:"from_entry_id"`
	ToEntryID     int64 `json:"to_entry_id"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		res, err := q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
			CreatedAt:     time.Now(),
		})
		if err != nil {
			return err
		}
		insertId, err := res.LastInsertId()
		if err != nil {
			return err
		}
		result.TransferID = insertId

		res, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return err
		}
		insertId, err = res.LastInsertId()
		if err != nil {
			return err
		}
		result.FromEntryID = insertId

		res, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return err
		}
		insertId, err = res.LastInsertId()
		if err != nil {
			return err
		}
		result.ToEntryID = insertId

		if arg.FromAccountID < arg.ToAccountID {
			err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		if err != nil {
			return err
		}
		result.FromAccountID = arg.FromAccountID
		result.ToAccountID = arg.ToAccountID
		return nil
	})
	return result, err
}

func addMoney(ctx context.Context, q *Queries, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64) error {
	err := q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})
	if err != nil {
		return err
	}
	err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	})
	if err != nil {
		return err
	}
	return nil
}
