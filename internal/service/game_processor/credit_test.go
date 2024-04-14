package game_processor

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"open-api-games/internal/domain"
	"open-api-games/internal/service/game_processor/mocks"
	"os"
	"testing"
)

func TestCredit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	t.Run("credit by session success", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

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
			On("CurrencyGetByCode", ctx, "USD").
			Return(&domain.Currency{
				Code:         "USD",
				Denomination: 2,
			}, nil)

		repoMock.
			On("BalanceIncrementByUserUIDAndCurrency", ctx, "123", "USD", 100).
			Return(&domain.Transaction{
				UID:          "123",
				Amount:       100,
				Currency:     "USD",
				Denomination: 2,
				Type:         domain.TransactionTypeCredit,
			}, nil)

		res, err := service.Credit(ctx, &domain.ProcessDebitCreditRollbackReq{
			GameSessionUID: "123",
			Currency:       "USD",
			Amount:         100,
		})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "123", res.TransactionUID)
		assert.Equal(t, "test", res.UserNick)
		assert.Equal(t, "USD", res.Currency)
		assert.Equal(t, 100, res.Amount)
		assert.Equal(t, 2, res.Denomination)

		repoMock.AssertExpectations(t)
	})

	t.Run("credit by user success", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		repoMock.
			On("UserGetByUID", ctx, "123").
			Return(&domain.User{
				UID:  "123",
				Nick: "test",
			}, nil)

		repoMock.
			On("CurrencyGetByCode", ctx, "USD").
			Return(&domain.Currency{
				Code:         "USD",
				Denomination: 2,
			}, nil)

		repoMock.
			On("BalanceIncrementByUserUIDAndCurrency", ctx, "123", "USD", 100).
			Return(&domain.Transaction{
				UID:          "123",
				Amount:       100,
				Currency:     "USD",
				Denomination: 2,
				Type:         domain.TransactionTypeCredit,
			}, nil)

		res, err := service.Credit(ctx, &domain.ProcessDebitCreditRollbackReq{
			UserUID:  "123",
			Currency: "USD",
			Amount:   100,
		})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "123", res.TransactionUID)
		assert.Equal(t, "test", res.UserNick)
		assert.Equal(t, "USD", res.Currency)
		assert.Equal(t, 100, res.Amount)
		assert.Equal(t, 2, res.Denomination)

		repoMock.AssertExpectations(t)
	})

	t.Run("credit by session not found", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		repoMock.
			On("SessionGetByUID", ctx, "123").
			Return(nil, domain.NewError(errorCreditSource).SetCode(domain.ErrNotFound))

		res, err := service.Credit(ctx, &domain.ProcessDebitCreditRollbackReq{
			GameSessionUID: "123",
			Currency:       "USD",
			Amount:         100,
		})

		assert.Equal(t, domain.AsError(err).Code, domain.ErrSessionNotFound)
		assert.Nil(t, res)

		repoMock.AssertExpectations(t)
	})

	t.Run("credit by user not found", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		repoMock.
			On("UserGetByUID", ctx, "123").
			Return(nil, domain.NewError(errorCreditSource).SetCode(domain.ErrNotFound))

		res, err := service.Credit(ctx, &domain.ProcessDebitCreditRollbackReq{
			UserUID:  "123",
			Currency: "USD",
			Amount:   100,
		})

		assert.Equal(t, domain.AsError(err).Code, domain.ErrUserNotFound)
		assert.Nil(t, res)

		repoMock.AssertExpectations(t)
	})

	t.Run("credit unknown currency", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		repoMock.
			On("UserGetByUID", ctx, "123").
			Return(&domain.User{
				UID:  "123",
				Nick: "test",
			}, nil)

		repoMock.
			On("CurrencyGetByCode", ctx, "USD").
			Return(nil, domain.NewError(errorCreditSource).SetCode(domain.ErrNotFound))

		res, err := service.Credit(ctx, &domain.ProcessDebitCreditRollbackReq{
			UserUID:  "123",
			Currency: "USD",
			Amount:   100,
		})

		assert.Equal(t, domain.AsError(err).Code, domain.ErrUnknownCurrency)
		assert.Nil(t, res)

		repoMock.AssertExpectations(t)
	})

	t.Run("credit error increment", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		repoMock.
			On("UserGetByUID", ctx, "123").
			Return(&domain.User{
				UID:  "123",
				Nick: "test",
			}, nil)

		repoMock.
			On("CurrencyGetByCode", ctx, "USD").
			Return(&domain.Currency{
				Code:         "USD",
				Denomination: 2,
			}, nil)

		repoMock.
			On("BalanceIncrementByUserUIDAndCurrency", ctx, "123", "USD", 100).
			Return(nil, domain.NewError(errorCreditSource).SetCode(domain.ErrIncrement))

		res, err := service.Credit(ctx, &domain.ProcessDebitCreditRollbackReq{
			UserUID:  "123",
			Currency: "USD",
			Amount:   100,
		})

		assert.Equal(t, domain.AsError(err).Code, domain.ErrIncrement)
		assert.Nil(t, res)

		repoMock.AssertExpectations(t)
	})
}
