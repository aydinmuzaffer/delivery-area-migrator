package gormtransaction

import (
	"context"

	"gorm.io/gorm"
)

var _ DBWrapper = (*gormDBWrapper)(nil) // compile time proof

// DBWrapper wraps db connection.
type DBWrapper interface {
	GetDB(ctx context.Context) *gorm.DB
}

type gormDBWrapper struct {
	db *gorm.DB
}

// NewGormDBWrapper creates DBWrapper.
func NewGormDBWrapper(db *gorm.DB) DBWrapper {
	return &gormDBWrapper{db: db}
}

func (dw *gormDBWrapper) GetDB(ctx context.Context) *gorm.DB {
	if txDB := extractTx(ctx); txDB != nil {
		return txDB
	}

	return dw.db.WithContext(ctx)
}
