package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new store object and returns it
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db), // *sql.DB created by calling New func created by sqlc
	}
}

// execTx Executes a function within the generated db transactions - fn is a callback function
func (store *SQLStore) execTx(cxt context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(cxt, nil)

	if err != nil {
		return err

	}

	// Otherwise call New() with the created transaction
	q := New(tx)

	// Callback the callback function with the new transaction
	// Should return an error
	err = fn(q)

	if err != nil {
		// Generate a rollback error
		rbErr := tx.Rollback()

		if rbErr != nil {
			// Combine both errors and return them
			return fmt.Errorf("transaction Error: %v, rollback Error: %v", err, rbErr)
		}
		// Otherwise return the original transaction error
		return err

	}
	// If we have reached here - commit the transaction
	return tx.Commit()

}

// TransferTxParams contains all the input params of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult contains all the results of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"` // Created transfer record
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"` // the entry for the account money is debited from
	ToEntry     Entry    `json:"to_entry"`   // the entry for the account money is credited to

}

// TransferTx performs money transfer from one account to another
// It creates a transfer record, adds account entries, and updates accounts bal within a
// single db transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		// Creates the from entry records
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // negative arg since money is moving out
		})

		if err != nil {
			return err
		}

		// Creates the to entry records
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return err
	})

	return result, err

}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
