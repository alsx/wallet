package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Modeller defines method to works with payments and accounts.
type Modeller interface {
	DoPayment(ctx context.Context, fromAccount, toAccount string, amount float32) error
	SelectPayments(ctx context.Context) ([]Payment, error)
	SelectAccounts(ctx context.Context) ([]Account, error)
}

// Model is a key sttucture to tie methods and share db connection.
type Model struct {
	db *sql.DB
}

// NewModel creates new Model with db connection.
func NewModel(dsn string) Model {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
	return Model{db}
}

// DoPayment transfers money from one account to another, logging this in payment table.
func (m Model) DoPayment(ctx context.Context, fromAccount, toAccount string, amount float32) error {
	conn, err := m.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("cannot create db conn:%v", err)
	}
	defer conn.Close()
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("cannot begin transaction:%v", err)
	}
	_, execErr := tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2;", amount, fromAccount)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("update failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
		}
		return fmt.Errorf("update failed: %v", execErr)
	}
	_, execErr = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2;", amount, toAccount)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("update failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
		}
		return fmt.Errorf("update failed: %v", execErr)
	}
	_, execErr = tx.ExecContext(ctx, "INSERT INTO payments VALUES ($1, $2, $3);", fromAccount, amount, toAccount)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("insert failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
		}
		return fmt.Errorf("insert failed: %v", execErr)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("cannot commit:%v", err)
	}
	return nil
}

// SelectPayments selects all payments.
// Returns a double entry bookkeeping.
func (m Model) SelectPayments(ctx context.Context) ([]Payment, error) {
	conn, err := m.db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close() // Return the connection to the pool.

	payments := []Payment{}
	// TODO: good to have a paging.
	rows, err := conn.QueryContext(ctx, "SELECT from_account, amount::money::numeric::float8, to_account FROM payments")
	if err != nil {
		return payments, fmt.Errorf("select payments error:%v", err)
	}
	defer rows.Close()
	for rows.Next() {
		payment := Payment{}
		if err := rows.Scan(&payment.FromAccount, &payment.Amount, &payment.ToAccount); err != nil {
			log.Fatal(err)
		}
		outgoingPayment := Payment{
			Direction:   "outgoing",
			Account:     payment.FromAccount,
			ToAccount:   payment.ToAccount,
			Amount:      payment.Amount,
			FromAccount: "",
		}
		incomingPayment := Payment{
			Direction:   "incoming",
			Account:     payment.ToAccount,
			FromAccount: payment.FromAccount,
			Amount:      payment.Amount,
			ToAccount:   "",
		}
		payments = append(payments, outgoingPayment, incomingPayment)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return payments, nil
}

// SelectAccounts select a list of all accounts.
func (m Model) SelectAccounts(ctx context.Context) ([]Account, error) {
	conn, err := m.db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close() // Return the connection to the pool.
	accounts := []Account{}
	// TODO: good to have a paging.
	rows, err := conn.QueryContext(ctx, "SELECT id, balance::money::numeric::float8 FROM accounts")
	if err != nil {
		return accounts, fmt.Errorf("select accounts error:%v", err)
	}
	defer rows.Close()
	for rows.Next() {
		account := Account{}
		if err := rows.Scan(&account.ID, &account.Balance); err != nil {
			log.Fatal(err)
		}
		accounts = append(accounts, account)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return accounts, nil
}
