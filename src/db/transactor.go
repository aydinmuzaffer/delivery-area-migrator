package db

import "context"

// TransactionHandleFunc is pattern of function.
type TransactionHandleFunc func(ctx context.Context) error

// Transactor runs logic inside a single database transaction.
type Transactor interface {
	// WithinTransaction runs a function within a database transaction.
	//
	// Transaction is propagated in the context,
	// so it is important to propagate it to underlying repositories.
	// Function commits if error is nil, and rollbacks if not.
	// It returns the same error.
	WithinTransaction(context.Context, TransactionHandleFunc) error
}
