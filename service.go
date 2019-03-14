package main

import (
	"context"
	"errors"
	"fmt"
)

// Service is a simple list of interface for user accounts.
type Service interface {
	GetAccounts(ctx context.Context) ([]Account, error)
	PostPayment(ctx context.Context, p Payment) error
	GetPayments(ctx context.Context) ([]Payment, error)
}

// Account represents a single user account.
// ID should be globally unique.
type Account struct {
	ID       string  `json:"id"`
	Balance  float32 `json:"balance"`
	Currency string  `json:"currency"`
}

// Payment represents a user payment.
type Payment struct {
	Account     string  `json:"account"`
	Amount      float32 `json:"amount"`
	ToAccount   string  `json:"to_account,omitempty"`
	FromAccount string  `json:"from_account,omitempty"`
	Direction   string  `json:"direction"` // incoming | outgoing
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type walletService struct {
	model Modeller
}

// NewWalletService creates new wallet service to handle payments between accounts.
func NewWalletService(dsn string) Service {
	return &walletService{
		model: NewModel(dsn),
	}
}

// PostPayment handles payments.
func (s *walletService) PostPayment(ctx context.Context, p Payment) error {
	if p.Amount < 0 {
		return fmt.Errorf("cannot use negative value %.2f as amount", p.Amount)
	}
	if p.FromAccount == "" {
		return fmt.Errorf("from_account cannot be empty %v", p)
	}
	if p.ToAccount == "" {
		return fmt.Errorf("to_account cannot be empty %v", p)
	}
	err := s.model.DoPayment(ctx, p.FromAccount, p.ToAccount, p.Amount)
	return err
}

// GetPayments handles request to list all payments.
func (s *walletService) GetPayments(ctx context.Context) ([]Payment, error) {
	p, err := s.model.SelectPayments(ctx)
	return p, err
}

// GetPayments handles request to list all accounts.
func (s *walletService) GetAccounts(ctx context.Context) ([]Account, error) {
	a, err := s.model.SelectAccounts(ctx)
	return a, err
}
