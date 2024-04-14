package game_processor

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"open-api-games/internal/domain"
	"open-api-games/internal/service/game_processor/mocks"
	"os"
	"testing"
)

func TestBalance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	t.Run("get balance success", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		repoMock.
			On("CurrencyGetByCode", ctx, "USD").
			Return(&domain.Currency{
				Code:         "USD",
				Denomination: 2,
			}, nil)

		repoMock.
			On("SessionGetByUID", ctx, "123").
			Return(&domain.Session{
				UserUID: "123",
				UID:     "123",
			}, nil)

		repoMock.
			On("UserGetByUID", ctx, "123").
			Return(&domain.User{
				UID:  "123",
				Nick: "test",
			}, nil)

		repoMock.
			On("BalanceGetByUserUIDAndCurrency", ctx, "123", "USD").
			Return(&domain.Balance{
				Amount:       100,
				Denomination: 2,
				Currency:     "USD",
			}, nil)

		res, err := service.Balance(ctx, &domain.ProcessBalanceReq{
			GameSessionUID: "123",
			Currency:       "USD",
		})

		assert.NoError(t, err)
		assert.Equal(t, "123", res.UserUID)
		assert.Equal(t, "test", res.UserNick)
		assert.Equal(t, 100, res.Amount)
		assert.Equal(t, "USD", res.Currency)
		assert.Equal(t, 2, res.Denomination)

		repoMock.AssertExpectations(t)
	})

	t.Run("get balance unknown currency", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		repoMock.
			On("CurrencyGetByCode", ctx, "USD").
			Return(nil, domain.NewError(errorBalanceSource).SetCode(domain.ErrNotFound))

		res, err := service.Balance(ctx, &domain.ProcessBalanceReq{
			GameSessionUID: "123",
			Currency:       "USD",
		})

		assert.Equal(t, domain.AsError(err).Code, domain.ErrUnknownCurrency)
		assert.Nil(t, res)

		repoMock.AssertExpectations(t)
	})

	t.Run("get balance session not found", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		repoMock.
			On("CurrencyGetByCode", ctx, "USD").
			Return(&domain.Currency{
				Code:         "USD",
				Denomination: 2,
			}, nil)

		repoMock.
			On("SessionGetByUID", ctx, "123").
			Return(nil, domain.NewError(errorBalanceSource).SetCode(domain.ErrNotFound))

		res, err := service.Balance(ctx, &domain.ProcessBalanceReq{
			GameSessionUID: "123",
			Currency:       "USD",
		})

		assert.Equal(t, domain.AsError(err).Code, domain.ErrSessionNotFound)
		assert.Nil(t, res)

		repoMock.AssertExpectations(t)
	})

	t.Run("get balance user not found", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		repoMock.
			On("CurrencyGetByCode", ctx, "USD").
			Return(&domain.Currency{
				Code:         "USD",
				Denomination: 2,
			}, nil)

		repoMock.
			On("SessionGetByUID", ctx, "123").
			Return(&domain.Session{
				UserUID: "123",
				UID:     "123",
			}, nil)

		repoMock.
			On("UserGetByUID", ctx, "123", mock.Anything).
			Return(nil, domain.NewError(errorBalanceSource).SetCode(domain.ErrNotFound))

		res, err := service.Balance(ctx, &domain.ProcessBalanceReq{
			GameSessionUID: "123",
			Currency:       "USD",
		})

		assert.Equal(t, domain.AsError(err).Code, domain.ErrUserNotFound)
		assert.Nil(t, res)

		repoMock.AssertExpectations(t)
	})

	t.Run("get balance not found", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		repoMock.
			On("CurrencyGetByCode", ctx, "USD").
			Return(&domain.Currency{
				Code:         "USD",
				Denomination: 2,
			}, nil)

		repoMock.
			On("SessionGetByUID", ctx, "123").
			Return(&domain.Session{
				UserUID: "123",
				UID:     "123",
			}, nil)

		repoMock.
			On("UserGetByUID", ctx, "123", mock.Anything).
			Return(&domain.User{
				UID:  "123",
				Nick: "test",
			}, nil)

		repoMock.
			On("BalanceGetByUserUIDAndCurrency", ctx, "123", "USD").
			Return(nil, domain.NewError(errorBalanceSource).SetCode(domain.ErrNotFound))

		res, err := service.Balance(ctx, &domain.ProcessBalanceReq{
			GameSessionUID: "123",
			Currency:       "USD",
		})

		assert.Equal(t, domain.AsError(err).Code, domain.ErrBalanceNotFound)
		assert.Nil(t, res)

		repoMock.AssertExpectations(t)
	})

}
