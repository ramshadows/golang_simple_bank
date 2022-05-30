package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	// Create a NewStore by passing the db connection
	store := NewStore(testDB)

	// Create two random accounts to transfer from and to
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	// Print initial Balances
	fmt.Printf(">> Before Transfer Balance: Account1: %d, Account2: %d", account1.Balance, account2.Balance)

	// Note: run a concurrent transfer transaction
	n := 5 //run 5 concurrent transactions

	amount := int64(100) // transfer an amount of 100 from account1 to account2

	// declare and create a channel to send to the results and error of the goroutine below
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

	// check results
	existed := make(map[int]bool)

	// Check results by receiving the values from the channels
	// Use a for loop
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		// Receive from the channel
		result := <-results
		// check the result
		require.NotEmpty(t, result)

		// Now check each result
		// The transfer itself
		transfer := result.Transfer
		require.NotEmpty(t, transfer)

		// Next contents of the transfer
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)

		// Then other auto populated fields
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// Now verify that the transfer was successiful
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check from Entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)

		// Then other auto populated fields
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		// Now verify that the entry was successiful
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// Check to Entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)

		// Then other auto populated fields
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// Verify accounts
		fromAccount := result.FromAccount

		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount

		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		//Check account balance
		fmt.Println("Transaction Results:", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)

		// diff1 should be greater than zero
		require.True(t, diff1 > 0)

		// Also should be divisible by the amount transfered
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		// k must be an interger btwn k and n
		// n is the number of executed transactions
		// k must be unique for each transaction up to k = n
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		// the existed map should not contain k
		require.NotContains(t, existed, k)
		// set k to true
		existed[k] = true

	}

	// check the final updated balance from the db
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> Update Balances:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

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

	// check the final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
