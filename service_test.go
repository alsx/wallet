package main

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
)

// ModelMock fakes db models.
type ModelMock struct {
	db *sql.DB
}

// DoPayment fakes update/insert into db.
func (m ModelMock) DoPayment(ctx context.Context, fromAccount, toAccount string, amount float32) error {
	return nil
}

// SelectPayments fakes selecting payments.
func (m ModelMock) SelectPayments(ctx context.Context) ([]Payment, error) {
	return []Payment{{FromAccount: "bob123", Amount: 0.05, ToAccount: "alice456"}}, nil
}

// SelectAccounts fakes selecting accounts.
func (m ModelMock) SelectAccounts(ctx context.Context) ([]Account, error) {
	return []Account{{ID: "alice456", Balance: 0.01}}, nil
}

// Test_walletService_PostPayment tests do payments.
func Test_walletService_PostPayment(t *testing.T) {
	type fields struct {
		model ModelMock
	}
	type args struct {
		ctx context.Context
		p   Payment
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"check valid", fields{model: ModelMock{}}, args{p: Payment{FromAccount: "bob123", Amount: 0.05, ToAccount: "alice456"}}, false},
		{"check valid", fields{model: ModelMock{}}, args{p: Payment{FromAccount: "", Amount: 0.05, ToAccount: "alice456"}}, true},
		{"check valid", fields{model: ModelMock{}}, args{p: Payment{FromAccount: "John789", Amount: -5, ToAccount: "alice456"}}, true},
		{"check valid", fields{model: ModelMock{}}, args{p: Payment{FromAccount: "John789", Amount: 50, ToAccount: ""}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &walletService{
				model: tt.fields.model,
			}
			if err := s.PostPayment(tt.args.ctx, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("walletService.PostPayment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test_walletService_GetPayments tests getting all payments.
func Test_walletService_GetPayments(t *testing.T) {
	type fields struct {
		model ModelMock
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Payment
		wantErr bool
	}{
		{"check valid", fields{model: ModelMock{}}, args{}, []Payment{{FromAccount: "bob123", Amount: 0.05, ToAccount: "alice456"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &walletService{
				model: tt.fields.model,
			}
			got, err := s.GetPayments(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("walletService.GetPayments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("walletService.GetPayments() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_walletService_GetAccounts tests getting all accounts.
func Test_walletService_GetAccounts(t *testing.T) {
	type fields struct {
		model ModelMock
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Account
		wantErr bool
	}{
		{"check valid", fields{model: ModelMock{}}, args{}, []Account{{ID: "alice456", Balance: 0.01}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &walletService{
				model: tt.fields.model,
			}
			got, err := s.GetAccounts(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("walletService.GetAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("walletService.GetAccounts() = %v, want %v", got, tt.want)
			}
		})
	}
}
