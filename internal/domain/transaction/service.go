package transaction

import (
	"context"
	"errors"
)

type Service interface {
	Create(ctx context.Context, in CreateRequest) (Response, error)
	Get(ctx context.Context, txID string) (Response, error)
	List(ctx context.Context, page, size int) ([]Response, int, int64, error)
	Update(ctx context.Context, txID string, in UpdateRequest) (Response, error)
	Delete(ctx context.Context, txID string) error
}
type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Create(ctx context.Context, in CreateRequest) (Response, error) {
	if err := ValidateCreate(in); err != nil {
		return Response{}, err
	}
	entity := &Transaction{
		TransactionID:          in.TransactionID,
		NoRef:                  in.NoRef,
		OrderTypeCode:          in.OrderTypeCode,
		OrderTypeName:          in.OrderTypeName,
		TransactionTypeCode:    in.TransactionTypeCode,
		TransactionTypeName:    in.TransactionTypeName,
		TransactionDate:        in.TransactionDate,
		FromAccountNumber:      in.FromAccountNumber,
		FromAccountName:        in.FromAccountName,
		FromAccountProductName: in.FromAccountProductName,
		ToAccountNumber:        in.ToAccountNumber,
		ToAccountName:          in.ToAccountName,
		ToAccountProductName:   in.ToAccountProductName,
		Amount:                 in.Amount,
		Status:                 in.Status,
		Description:            in.Description,
		Method:                 in.Method,
		Currency:               in.Currency,
		Metadata:               in.Metadata,
	}
	if err := s.repo.Create(ctx, entity); err != nil {
		return Response{}, err
	}
	return ToResponse(entity), nil
}

func (s *service) Get(ctx context.Context, txID string) (Response, error) {
	found, err := s.repo.GetByTxID(ctx, txID)
	if err != nil {
		return Response{}, err
	}
	if found == nil {
		return Response{}, errors.New("not_found")
	}
	return ToResponse(found), nil
}

func (s *service) List(ctx context.Context, page, size int) ([]Response, int, int64, error) {
	items, total, err := s.repo.List(ctx, page, size)
	if err != nil {
		return nil, 0, 0, err
	}
	out := make([]Response, 0, len(items))
	for i := range items {
		out = append(out, ToResponse(&items[i]))
	}
	return out, page, total, nil
}

func (s *service) Update(ctx context.Context, txID string, in UpdateRequest) (Response, error) {
	if err := in.Validate(); err != nil {
		return Response{}, err
	}
	found, err := s.repo.GetByTxID(ctx, txID)
	if err != nil {
		return Response{}, err
	}
	if found == nil {
		return Response{}, errors.New("not_found")
	}
	// patch
	if in.NoRef != nil {
		found.NoRef = *in.NoRef
	}
	if in.OrderTypeCode != nil {
		found.OrderTypeCode = *in.OrderTypeCode
	}
	if in.OrderTypeName != nil {
		found.OrderTypeName = *in.OrderTypeName
	}
	if in.TransactionTypeCode != nil {
		found.TransactionTypeCode = *in.TransactionTypeCode
	}
	if in.TransactionTypeName != nil {
		found.TransactionTypeName = *in.TransactionTypeName
	}
	if in.TransactionDate != nil {
		found.TransactionDate = *in.TransactionDate
	}
	if in.FromAccountNumber != nil {
		found.FromAccountNumber = *in.FromAccountNumber
	}
	if in.FromAccountName != nil {
		found.FromAccountName = *in.FromAccountName
	}
	if in.FromAccountProductName != nil {
		found.FromAccountProductName = *in.FromAccountProductName
	}
	if in.ToAccountNumber != nil {
		found.ToAccountNumber = *in.ToAccountNumber
	}
	if in.ToAccountName != nil {
		found.ToAccountName = *in.ToAccountName
	}
	if in.ToAccountProductName != nil {
		found.ToAccountProductName = *in.ToAccountProductName
	}
	if in.Amount != nil {
		found.Amount = *in.Amount
	}
	if in.Status != nil {
		found.Status = *in.Status
	}
	if in.Description != nil {
		found.Description = *in.Description
	}
	if in.Method != nil {
		found.Method = *in.Method
	}
	if in.Currency != nil {
		found.Currency = *in.Currency
	}
	if in.Metadata != nil {
		found.Metadata = *in.Metadata
	}

	if err := s.repo.Update(ctx, found); err != nil {
		return Response{}, err
	}
	return ToResponse(found), nil
}

func (s *service) Delete(ctx context.Context, txID string) error {
	return s.repo.DeleteByTxID(ctx, txID)
}
