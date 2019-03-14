package main

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all of the endpoints that compose a wallet service.
type Endpoints struct {
	PostPaymentEndpoint endpoint.Endpoint
	GetPaymentsEndpoint endpoint.Endpoint
	GetAccountsEndpoint endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service.
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostPaymentEndpoint: MakePostPaymentEndpoint(s),
		GetPaymentsEndpoint: MakeGetPaymentsEndpoint(s),
		GetAccountsEndpoint: MakeGetAccountsEndpoint(s),
	}
}

// MakePostProfileEndpoint returns an endpoint via the passed service.
func MakePostPaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postPaymentRequest)
		e := s.PostPayment(ctx, req.Payment)
		return postPaymentResponse{Err: e}, nil
	}
}

// MakeGetAddressesEndpoint returns an endpoint via the passed service.
func MakeGetPaymentsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		a, e := s.GetPayments(ctx)
		return getPaymentsResponse{Payments: a, Err: e}, nil
	}
}

// MakeGetAddressesEndpoint returns an endpoint via the passed service.
func MakeGetAccountsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		a, e := s.GetAccounts(ctx)
		return getAccountsResponse{Accounts: a, Err: e}, nil
	}
}

type postPaymentRequest struct {
	Payment Payment
}

type postPaymentResponse struct {
	Err error `json:"err,omitempty"`
}

func (r postPaymentResponse) error() error { return r.Err }

type getPaymentsRequest struct {
	ID string
}

type getPaymentsResponse struct {
	Payments []Payment `json:"payments,omitempty"`
	Err      error     `json:"err,omitempty"`
}

func (r getPaymentsResponse) error() error { return r.Err }

type getAccountsRequest struct {
	ID string
}

type getAccountsResponse struct {
	Accounts []Account `json:"accounts,omitempty"`
	Err      error     `json:"err,omitempty"`
}

func (r getAccountsResponse) error() error { return r.Err }
