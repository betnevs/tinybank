package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

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

		// TODO update account's balance

		return nil
	})

	return result, err

}
