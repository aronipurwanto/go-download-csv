package transaction

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/datatypes"
	"time"
)

type CreateRequest struct {
	TransactionID          string         `json:"transaction_id"`
	NoRef                  string         `json:"no_ref"`
	OrderTypeCode          string         `json:"order_type_code" validate:"required"`
	OrderTypeName          string         `json:"order_type_name" validate:"required"`
	TransactionTypeCode    string         `json:"transaction_type_code" validate:"required"`
	TransactionTypeName    string         `json:"transaction_type_name" validate:"required"`
	TransactionDate        time.Time      `json:"transaction_date" validate:"required"`
	FromAccountNumber      string         `json:"from_account_number" validate:"required"`
	FromAccountName        string         `json:"from_account_name" validate:"required"`
	FromAccountProductName string         `json:"from_account_product_name" validate:"required"`
	ToAccountNumber        string         `json:"to_account_number" validate:"required"`
	ToAccountName          string         `json:"to_account_name" validate:"required"`
	ToAccountProductName   string         `json:"to_account_product_name" validate:"required"`
	Amount                 float64        `json:"amount" validate:"required,gt=0"`
	Status                 string         `json:"status" validate:"required,oneof=PENDING SUCCESS FAILED"`
	Description            string         `json:"description" validate:"omitempty,max=255"`
	Method                 string         `json:"method" validate:"required"`
	Currency               string         `json:"currency" validate:"required,len=3"`
	Metadata               datatypes.JSON `json:"metadata" validate:"omitempty"`
}

type UpdateRequest struct {
	NoRef                  *string         `json:"no_ref"`
	OrderTypeCode          *string         `json:"order_type_code"`
	OrderTypeName          *string         `json:"order_type_name"`
	TransactionTypeCode    *string         `json:"transaction_type_code"`
	TransactionTypeName    *string         `json:"transaction_type_name"`
	TransactionDate        *time.Time      `json:"transaction_date"`
	FromAccountNumber      *string         `json:"from_account_number"`
	FromAccountName        *string         `json:"from_account_name"`
	FromAccountProductName *string         `json:"from_account_product_name"`
	ToAccountNumber        *string         `json:"to_account_number"`
	ToAccountName          *string         `json:"to_account_name"`
	ToAccountProductName   *string         `json:"to_account_product_name"`
	Amount                 *float64        `json:"amount"`
	Status                 *string         `json:"status"`
	Description            *string         `json:"description"`
	Method                 *string         `json:"method"`
	Currency               *string         `json:"currency"`
	Metadata               *datatypes.JSON `json:"metadata"`
}

// Response DTO (what we expose)

type Response struct {
	TransactionID          string         `json:"transaction_id"`
	NoRef                  string         `json:"no_ref"`
	OrderTypeCode          string         `json:"order_type_code"`
	OrderTypeName          string         `json:"order_type_name"`
	TransactionTypeCode    string         `json:"transaction_type_code"`
	TransactionTypeName    string         `json:"transaction_type_name"`
	TransactionDate        time.Time      `json:"transaction_date"`
	FromAccountNumber      string         `json:"from_account_number"`
	FromAccountName        string         `json:"from_account_name"`
	FromAccountProductName string         `json:"from_account_product_name"`
	ToAccountNumber        string         `json:"to_account_number"`
	ToAccountName          string         `json:"to_account_name"`
	ToAccountProductName   string         `json:"to_account_product_name"`
	Amount                 float64        `json:"amount"`
	Status                 string         `json:"status"`
	Description            string         `json:"description"`
	Method                 string         `json:"method"`
	Currency               string         `json:"currency"`
	Metadata               datatypes.JSON `json:"metadata"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
}

func ToResponse(e *Transaction) Response {
	return Response{
		TransactionID:          e.TransactionID,
		NoRef:                  e.NoRef,
		OrderTypeCode:          e.OrderTypeCode,
		OrderTypeName:          e.OrderTypeName,
		TransactionTypeCode:    e.TransactionTypeCode,
		TransactionTypeName:    e.TransactionTypeName,
		TransactionDate:        e.TransactionDate,
		FromAccountNumber:      e.FromAccountNumber,
		FromAccountName:        e.FromAccountName,
		FromAccountProductName: e.FromAccountProductName,
		ToAccountNumber:        e.ToAccountNumber,
		ToAccountName:          e.ToAccountName,
		ToAccountProductName:   e.ToAccountProductName,
		Amount:                 e.Amount,
		Status:                 e.Status,
		Description:            e.Description,
		Method:                 e.Method,
		Currency:               e.Currency,
		Metadata:               e.Metadata,
		CreatedAt:              e.CreatedAt,
		UpdatedAt:              e.UpdatedAt,
	}
}

var validate = validator.New()

func ValidateCreate(r CreateRequest) error {
	return validate.Struct(r)
}

func (u UpdateRequest) Validate() error {
	return validate.Struct(u) // fields are optional
}
