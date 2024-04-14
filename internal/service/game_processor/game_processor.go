package game_processor

import (
	"context"
	"log/slog"
	"open-api-games/internal/domain"
)

//go:generate mockery --dir . --name Repository --output ./mocks --case=underscore
type Repository interface {
	UserGetByUID(ctx context.Context, uid string) (*domain.User, error)
	SessionGetByUID(ctx context.Context, uid string) (*domain.Session, error)
	BalanceGetByUserUIDAndCurrency(ctx context.Context, userUID, currency string) (*domain.Balance, error)
	BalanceDecrementByUserUIDAndCurrency(ctx context.Context, userUID, currency string, amount int) (*domain.Transaction, error)
	BalanceIncrementByUserUIDAndCurrency(ctx context.Context, userUID, currency string, amount int) (*domain.Transaction, error)
	TransactionGetByUID(ctx context.Context, uid string) (*domain.Transaction, error)
	CurrencyGetByCode(ctx context.Context, code string) (*domain.Currency, error)
}

type Service struct {
	repo   Repository
	logger *slog.Logger
}

func New(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}
