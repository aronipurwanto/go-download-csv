package transaction

import (
	"gorm.io/datatypes"
	"time"
)

type Transaction struct {
	ID                     uint           `gorm:"primaryKey" json:"-"`
	TransactionID          string         `gorm:"uniqueIndex;size:36" json:"transaction_id"`
	NoRef                  string         `gorm:"size:64" json:"no_ref"`
	OrderTypeCode          string         `gorm:"size:32" json:"order_type_code"`
	OrderTypeName          string         `gorm:"size:128" json:"order_type_name"`
	TransactionTypeCode    string         `gorm:"size:32" json:"transaction_type_code"`
	TransactionTypeName    string         `gorm:"size:128" json:"transaction_type_name"`
	TransactionDate        time.Time      `json:"transaction_date"`
	FromAccountNumber      string         `gorm:"size:64" json:"from_account_number"`
	FromAccountName        string         `gorm:"size:128" json:"from_account_name"`
	FromAccountProductName string         `gorm:"size:128" json:"from_account_product_name"`
	ToAccountNumber        string         `gorm:"size:64" json:"to_account_number"`
	ToAccountName          string         `gorm:"size:128" json:"to_account_name"`
	ToAccountProductName   string         `gorm:"size:128" json:"to_account_product_name"`
	Amount                 float64        `json:"amount"`
	Status                 string         `gorm:"size:32" json:"status"`
	Description            string         `gorm:"size:255" json:"description"`
	Method                 string         `gorm:"size:64" json:"method"`
	Currency               string         `gorm:"size:16" json:"currency"`
	Metadata               datatypes.JSON `json:"metadata"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
