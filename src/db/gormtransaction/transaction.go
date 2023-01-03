package gormtransaction

import (
	"context"
	"fmt"
	"log"

	db "github.com/aydinmuzaffer/migration-tool-service/src/db"
	"gorm.io/gorm"
)

var _ db.Transactor = (*gormTransactor)(nil) // compile time proof

type gormTransactor struct {
	db *gorm.DB
}

// NewTransactor creates gormTransactor.
func NewTransactor(db *gorm.DB) db.Transactor {
	return &gormTransactor{db: db}
}

// WithinTransaction runs function within transaction
//
// The transaction commits when function were finished without error.
func (gt *gormTransactor) WithinTransaction(ctx context.Context, tFunc db.TransactionHandleFunc) (err error) {
	// begin transaction
	tx := gt.db.WithContext(ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
		if err != nil {
			if tx.Rollback(); tx.Error != nil {
				fmt.Printf("rollback err: %v", tx.Error)
			}
		}
	}()

	// run callback
	err = tFunc(injectTx(ctx, tx))

	if err != nil {
		return err //nolint
	}

	// if no error, commit
	if errCommit := tx.Commit().Error; errCommit != nil {
		log.Printf("commit transaction: %v", errCommit)
	}
	return nil
}

type txKey struct{}

// injectTx injects transaction to context.
func injectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx extracts transaction from context.
func extractTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return nil
}
