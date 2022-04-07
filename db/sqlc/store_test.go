package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// check result
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)
		// check transfer
		transferID := result.TransferID
		require.NotZero(t, transferID)

		transfer, err := store.GetTransfer(context.Background(), transferID)
		require.NoError(t, err)
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		require.Equal(t, transferID, transfer.ID)

		fromEntryID := result.FromEntryID
		require.NotZero(t, fromEntryID)

		fromEntry, err := store.GetEntry(context.Background(), fromEntryID)
		require.NoError(t, err)
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntryID, fromEntry.ID)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.CreatedAt)

		toEntryID := result.ToEntryID
		require.NotZero(t, toEntryID)

		toEntry, err := store.GetEntry(context.Background(), toEntryID)
		require.NoError(t, err)
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntryID, toEntry.ID)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.CreatedAt)

	}

}
