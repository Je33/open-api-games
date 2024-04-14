package game_processor

import (
	"context"
	"open-api-games/internal/domain"
)

const (
	errorRollbackSource = "[service.game_processor.rollback]"
)

func (s *Service) Rollback(ctx context.Context, req *domain.ProcessDebitCreditRollbackReq) (*domain.ProcessDebitCreditRollbackRes, error) {
	if req.TransactionUID == "" {
		return nil, domain.NewError(errorRollbackSource).SetCode(domain.ErrEmptyTransactionUID)
	}

	txn, err := s.repo.TransactionGetByUID(ctx, req.TransactionUID)
	if err != nil {
		return nil, domain.NewError(errorCreditSource).SetCode(domain.ErrTransactionNotFound).Add(err)
	}

	userUid := txn.UserUID

	user, err := s.repo.UserGetByUID(ctx, userUid)
	if err != nil {
		return nil, domain.NewError(errorCreditSource).SetCode(domain.ErrUserNotFound).Add(err)
	}

	// TODO: implement custom rollback logic to save already rollbacked transaction
	// just only create the opposite transaction for now

	var txnRollback *domain.Transaction
	switch txn.Type {
	case domain.TransactionTypeCredit:
		txnRollback, err = s.repo.BalanceDecrementByUserUIDAndCurrency(ctx, userUid, txn.Currency, txn.Amount)
	case domain.TransactionTypeDebit:
		txnRollback, err = s.repo.BalanceIncrementByUserUIDAndCurrency(ctx, userUid, txn.Currency, txn.Amount)
	default:
		return nil, domain.NewError(errorRollbackSource).SetCode(domain.ErrInvalidTransactionType)
	}
	if err != nil {
		return nil, domain.NewError(errorRollbackSource).SetCode(domain.ErrRollback).Add(err)
	}

	return &domain.ProcessDebitCreditRollbackRes{
		TransactionUID: txnRollback.UID,
		UserNick:       user.Nick,
		Amount:         txnRollback.Amount,
		Currency:       txnRollback.Currency,
		Denomination:   txnRollback.Denomination,
		MaxWin:         0, // TODO: implement MaxWin
	}, nil
}
