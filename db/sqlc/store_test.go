package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before: ", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
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

		// check balance
		fromAccountID := result.FromAccountID
		require.NotZero(t, fromAccountID)
		require.Equal(t, account1.ID, fromAccountID)

		toAccountID := result.ToAccountID
		require.NotZero(t, toAccountID)
		require.Equal(t, account2.ID, toAccountID)
	}

	// check final balance
	updateAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">> after: ", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)
}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> before: ", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error)
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updateAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">> after: ", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}
