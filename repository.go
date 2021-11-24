package gpa

import (
	"context"
	"gorm.io/gorm"
)

type Repository struct {
	*gorm.DB
}

func (r *Repository) PrepareQueryContext(ctx context.Context, dest interface{}, opts []Option) *QueryContext {
	tx, found := GetTxFromContext(ctx)
	if !found {
		tx = r.WithContext(ctx)
	}
	q := QueryContext{Tx: tx, Dest: dest, curr: -1}
	for _, opt := range opts {
		if opt != nil {
			opt(&q)
		}
	}
	return &q
}

func (r *Repository) Find(ctx context.Context, dest interface{}, opts ...Option) error {
	q := r.PrepareQueryContext(ctx, dest, opts)
	return q.Execute(func(tx *gorm.DB) error {
		return tx.Find(dest).Error
	})
}
