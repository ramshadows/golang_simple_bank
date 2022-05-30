package db

import (
	"context"
	"database/sql"
	"simple_bank/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	// call the require function passing the t -> testing.T object and the error
	require.NoError(t, err)
	// Check that the returned account is not empty
	require.NotEmpty(t, account)
	// Check that the returned owner is equal to the account owner
	require.Equal(t, arg.Owner, account.Owner)
	// Check that the returned balance is equal to the account balance
	require.Equal(t, arg.Balance, account.Balance)
	// Check that the returned currency is equal to the account owner
	require.Equal(t, arg.Currency, account.Currency)

	// Check if the account id is autmatically populated by postgres
	require.NotZero(t, account.ID)

	// Check if the account create at is autmatically populated by postgres
	require.NotZero(t, account.CreatedAt)

	return account

}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)

}

func TestGetAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)

	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	// Verifies that the account is really deleted
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	// The call should return an error
	require.Error(t, err)

	// Also should return an error equal to sql error no row found
	require.EqualError(t, err, sql.ErrNoRows.Error())

	// Also check the account2 object should be empty
	require.Empty(t, account2)

}

func TestListAccount(t *testing.T) {
	var lastAccount Account
	// Use a loop to create 10 random accounts
	for i := 0; i < 10; i++ {
		lastAccount = CreateRandomAccount(t)

	}

	// since we expect at least 10 random accounts
	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5, // skip 5
		Offset: 0, // return next 0
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	// Should return no error
	require.NoError(t, err)

	// Should return at least 5 accounts
	require.NotEmpty(t, accounts)

	// Iterate thro the list of accounts and require each of them to be not empty
	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)

	}
}
