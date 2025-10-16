package transaction

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, t *Transaction) error
	GetByTxID(ctx context.Context, txID string) (*Transaction, error)
	List(ctx context.Context, page, size int) ([]Transaction, int64, error)
	Update(ctx context.Context, t *Transaction) error
	DeleteByTxID(ctx context.Context, txID string) error // soft delete

	// Export helpers
	ListByDateRange(ctx context.Context, from, to time.Time) ([]Transaction, error)
}
