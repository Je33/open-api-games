package game_processor

import (
	"context"
	"open-api-games/internal/domain"
)

const (
	errorBalanceSource = "[service.game_processor.balance]"
)

// Balance processes get balance amount request
func (s *Service) Balance(ctx context.Context, req *domain.ProcessBalanceReq) (*domain.ProcessBalanceRes, error) {
	cur, err := s.repo.CurrencyGetByCode(ctx, req.Currency)
	if err != nil || cur == nil {
		return nil, domain.NewError(errorBalanceSource).SetCode(domain.ErrUnknownCurrency).Add(err)
	}

	session, err := s.repo.SessionGetByUID(ctx, req.GameSessionUID)
	if err != nil {
		return nil, domain.NewError(errorBalanceSource).SetCode(domain.ErrSessionNotFound).Add(err)
	}

	user, err := s.repo.UserGetByUID(ctx, session.UserUID)
	if err != nil {
		return nil, domain.NewError(errorBalanceSource).SetCode(domain.ErrUserNotFound).Add(err)
	}

	balance, err := s.repo.BalanceGetByUserUIDAndCurrency(ctx, user.UID, req.Currency)
	if err != nil {
		return nil, domain.NewError(errorBalanceSource).SetCode(domain.ErrBalanceNotFound).Add(err)
	}

	return &domain.ProcessBalanceRes{
		UserUID:      user.UID,
		UserNick:     user.Nick,
		Amount:       balance.Amount,
		Currency:     balance.Currency,
		Denomination: balance.Denomination,
		MaxWin:       0,  // TODO: implement MaxWin
		JpKey:        "", // TODO: implement JpKey
	}, nil
}
