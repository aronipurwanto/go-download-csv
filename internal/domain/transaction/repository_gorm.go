package transaction

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type gormRepository struct{ db *gorm.DB }

func NewGormRepository(db *gorm.DB) Repository { return &gormRepository{db: db} }

func (r *gormRepository) Create(ctx context.Context, t *Transaction) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *gormRepository) GetByTxID(ctx context.Context, txID string) (*Transaction, error) {
	var out Transaction
	err := r.db.WithContext(ctx).Where("transaction_id = ?", txID).First(&out).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &out, err
}

func (r *gormRepository) List(ctx context.Context, page, size int) ([]Transaction, int64, error) {
	var (
		items []Transaction
		total int64
	)
	db := r.db.WithContext(ctx).Model(&Transaction{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *gormRepository) Update(ctx context.Context, t *Transaction) error {
	return r.db.WithContext(ctx).Where("transaction_id = ?", t.TransactionID).Updates(t).Error
}

func (r *gormRepository) DeleteByTxID(ctx context.Context, txID string) error {
	return r.db.WithContext(ctx).Where("transaction_id = ?", txID).Delete(&Transaction{}).Error
}

func (r *gormRepository) ListByDateRange(ctx context.Context, from, to time.Time) ([]Transaction, error) {
	var items []Transaction
	db := r.db.WithContext(ctx).Model(&Transaction{})
	if !from.IsZero() {
		db = db.Where("transaction_date >= ?", from)
	}
	if !to.IsZero() {
		db = db.Where("transaction_date <= ?", to)
	}
	if err := db.Order("transaction_date ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
