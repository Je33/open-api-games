package seed

import (
	"context"
	"log/slog"
	"open-api-games/internal/domain"
	"open-api-games/internal/service/seed/data"
)

const (
	seedErrorSource = "[service.seed]"
)

type Repository interface {
	UserCreate(ctx context.Context, user *domain.User) error
	SessionCreate(ctx context.Context, sess *domain.Session) error
	CurrencyCreate(ctx context.Context, cur *domain.Currency) error
	BalanceCreate(ctx context.Context, balance *domain.Balance) error
}

type Service struct {
	db     Repository
	logger *slog.Logger
}

func New(db Repository, logger *slog.Logger) *Service {
	return &Service{
		db,
		logger,
	}
}

func (s *Service) Seed(ctx context.Context) error {
	currencies := data.Currencies()
	for _, c := range currencies {
		s.db.CurrencyCreate(ctx, &c)
	}

	users := data.Users()
	for _, u := range users {
		s.db.UserCreate(ctx, &u)
	}

	sessions := data.Sessions()
	for _, us := range sessions {
		s.db.SessionCreate(ctx, &us)
	}

	balances := data.Balances()
	for _, bl := range balances {
		s.db.BalanceCreate(ctx, &bl)
	}

	return nil
}
