package main

// The profilesvc is just over HTTP, so we just have a single transport.go.

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// var (
// 	// ErrBadRouting is returned when an expected path variable is missing.
// 	// It always indicates programmer error.
// 	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
// )

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST    /v1/payments/            adds another payments
	// GET     /v1/payments/            retrieve all payments
	// GET     /v1/accounts/            retrieve all accounts

	r.Methods("POST").Path("/v1/payments/").Handler(httptransport.NewServer(
		e.PostPaymentEndpoint,
		decodePostPaymentRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/v1/payments/").Handler(httptransport.NewServer(
		e.GetPaymentsEndpoint,
		decodeGetPaymentsRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/v1/accounts/").Handler(httptransport.NewServer(
		e.GetAccountsEndpoint,
		decodeGetAccountsRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodePostPaymentRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postPaymentRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Payment); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetPaymentsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return getPaymentsRequest{}, nil
}

func decodeGetAccountsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return getAccountsRequest{}, nil
}

// errorer is implemented by all concrete response types that may contain errors.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
