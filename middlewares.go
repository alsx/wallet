package main

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

// LoggingMiddleware logs http requests.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

// PostPayment wraps request with logger.
func (mw loggingMiddleware) PostPayment(ctx context.Context, p Payment) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostPayment", "id", p.Account, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostPayment(ctx, p)
}

// GetPayments wraps request with logger.
func (mw loggingMiddleware) GetPayments(ctx context.Context) (payments []Payment, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetAddresses", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetPayments(ctx)
}

// GetAccounts wraps request with logger.
func (mw loggingMiddleware) GetAccounts(ctx context.Context) (accounts []Account, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetAddresses", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetAccounts(ctx)
}
