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

func TestDebit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	t.Run("debit by session success", func(t *testing.T) {
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
			On("BalanceDecrementByUserUIDAndCurrency", ctx, "123", "USD", 100).
			Return(&domain.Transaction{
				UID:          "123",
				Amount:       100,
				Currency:     "USD",
				Denomination: 2,
				Type:         domain.TransactionTypeDebit,
			}, nil)

		res, err := service.Debit(ctx, &domain.ProcessDebitCreditRollbackReq{
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

	t.Run("debit by user success", func(t *testing.T) {
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
			On("BalanceDecrementByUserUIDAndCurrency", ctx, "123", "USD", 100).
			Return(&domain.Transaction{
				UID:          "123",
				Amount:       100,
				Currency:     "USD",
				Denomination: 2,
				Type:         domain.TransactionTypeDebit,
			}, nil)

		res, err := service.Debit(ctx, &domain.ProcessDebitCreditRollbackReq{
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

	t.Run("debit by session not found", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		repoMock.
			On("SessionGetByUID", ctx, "123").
			Return(nil, domain.NewError(errorDebitSource).SetCode(domain.ErrNotFound))

		res, err := service.Debit(ctx, &domain.ProcessDebitCreditRollbackReq{
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
			Return(nil, domain.NewError(errorDebitSource).SetCode(domain.ErrNotFound))

		res, err := service.Debit(ctx, &domain.ProcessDebitCreditRollbackReq{
			UserUID:  "123",
			Currency: "USD",
			Amount:   100,
		})

		assert.Equal(t, domain.AsError(err).Code, domain.ErrUserNotFound)
		assert.Nil(t, res)

		repoMock.AssertExpectations(t)
	})

	t.Run("debit unknown currency", func(t *testing.T) {
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
			Return(nil, domain.NewError(errorDebitSource).SetCode(domain.ErrNotFound))

		res, err := service.Debit(ctx, &domain.ProcessDebitCreditRollbackReq{
			UserUID:  "123",
			Currency: "USD",
			Amount:   100,
		})

		assert.Equal(t, domain.AsError(err).Code, domain.ErrUnknownCurrency)
		assert.Nil(t, res)

		repoMock.AssertExpectations(t)
	})

	t.Run("debit insufficient funds", func(t *testing.T) {
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
			On("BalanceDecrementByUserUIDAndCurrency", ctx, "123", "USD", 100).
			Return(nil, domain.NewError(errorDebitSource).SetCode(domain.ErrDecrement))

		res, err := service.Debit(ctx, &domain.ProcessDebitCreditRollbackReq{
			UserUID:  "123",
			Currency: "USD",
			Amount:   100,
		})

		assert.Equal(t, domain.AsError(err).Code, domain.ErrDecrement)
		assert.Nil(t, res)

		repoMock.AssertExpectations(t)
	})
}
