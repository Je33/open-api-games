package game_processor

import (
	"context"
	"open-api-games/internal/domain"
)

const (
	errorCreditSource = "[service.game_processor.credit]"
)

func (s *Service) Credit(ctx context.Context, req *domain.ProcessDebitCreditRollbackReq) (*domain.ProcessDebitCreditRollbackRes, error) {
	userUid := req.UserUID
	if req.GameSessionUID != "" {
		session, err := s.repo.SessionGetByUID(ctx, req.GameSessionUID)
		if err != nil {
			return nil, domain.NewError(errorBalanceSource).SetCode(domain.ErrSessionNotFound).Add(err)
		}
		userUid = session.UserUID
	}

	user, err := s.repo.UserGetByUID(ctx, userUid)
	if err != nil {
		return nil, domain.NewError(errorCreditSource).SetCode(domain.ErrUserNotFound).Add(err)
	}

	cur, err := s.repo.CurrencyGetByCode(ctx, req.Currency)
	if err != nil || cur == nil {
		return nil, domain.NewError(errorCreditSource).SetCode(domain.ErrUnknownCurrency).Add(err)
	}

	txn, err := s.repo.BalanceIncrementByUserUIDAndCurrency(ctx, userUid, req.Currency, req.Amount)
	if err != nil {
		return nil, domain.NewError(errorCreditSource).SetCode(domain.ErrIncrement).Add(err)
	}

	return &domain.ProcessDebitCreditRollbackRes{
		TransactionUID: txn.UID,
		UserNick:       user.Nick,
		Amount:         txn.Amount,
		Currency:       txn.Currency,
		Denomination:   txn.Denomination,
		MaxWin:         0, // TODO: implement MaxWin
	}, nil
}
