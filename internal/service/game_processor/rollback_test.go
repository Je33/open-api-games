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

func TestRollback(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	t.Run("rollback debit success", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		repoMock.On("TransactionGetByUID", ctx, "123").Return(&domain.Transaction{
			UID:          "123",
			UserUID:      "123",
			Amount:       100,
			Currency:     "USD",
			Denomination: 2,
			Type:         domain.TransactionTypeDebit,
		}, nil)

		repoMock.
			On("UserGetByUID", ctx, "123").
			Return(&domain.User{
				UID:  "123",
				Nick: "test",
			}, nil)

		repoMock.
			On("BalanceIncrementByUserUIDAndCurrency", ctx, "123", "USD", 100).
			Return(&domain.Transaction{
				UID:          "1234",
				Amount:       100,
				Currency:     "USD",
				Denomination: 2,
				Type:         domain.TransactionTypeCredit,
			}, nil)

		res, err := service.Rollback(ctx, &domain.ProcessDebitCreditRollbackReq{
			TransactionUID: "123",
		})

		assert.NoError(t, err)
		assert.Equal(t, "1234", res.TransactionUID)
		assert.Equal(t, "test", res.UserNick)
		assert.Equal(t, 100, res.Amount)
		assert.Equal(t, "USD", res.Currency)
		assert.Equal(t, 2, res.Denomination)

		repoMock.AssertExpectations(t)
	})

	t.Run("rollback credit success", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		repoMock.On("TransactionGetByUID", ctx, "123").Return(&domain.Transaction{
			UID:          "123",
			UserUID:      "123",
			Amount:       100,
			Currency:     "USD",
			Denomination: 2,
			Type:         domain.TransactionTypeCredit,
		}, nil)

		repoMock.
			On("UserGetByUID", ctx, "123").
			Return(&domain.User{
				UID:  "123",
				Nick: "test",
			}, nil)

		repoMock.
			On("BalanceDecrementByUserUIDAndCurrency", ctx, "123", "USD", 100).
			Return(&domain.Transaction{
				UID:          "1234",
				Amount:       100,
				Currency:     "USD",
				Denomination: 2,
				Type:         domain.TransactionTypeDebit,
			}, nil)

		res, err := service.Rollback(ctx, &domain.ProcessDebitCreditRollbackReq{
			TransactionUID: "123",
		})

		assert.NoError(t, err)
		assert.Equal(t, "1234", res.TransactionUID)
		assert.Equal(t, "test", res.UserNick)
		assert.Equal(t, 100, res.Amount)
		assert.Equal(t, "USD", res.Currency)
		assert.Equal(t, 2, res.Denomination)

		repoMock.AssertExpectations(t)
	})

	t.Run("rollback transaction uid empty", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		res, err := service.Rollback(ctx, &domain.ProcessDebitCreditRollbackReq{
			TransactionUID: "",
		})

		assert.Equal(t, domain.AsError(err).Code, domain.ErrEmptyTransactionUID)
		assert.Nil(t, res)

		repoMock.AssertExpectations(t)
	})

	t.Run("rollback transaction not found", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		repoMock.
			On("TransactionGetByUID", ctx, "321").
			Return(nil, domain.NewError(errorRollbackSource).SetCode(domain.ErrNotFound))

		res, err := service.Rollback(ctx, &domain.ProcessDebitCreditRollbackReq{
			TransactionUID: "321",
		})

		assert.Equal(t, domain.AsError(err).Code, domain.ErrTransactionNotFound)
		assert.Nil(t, res)

		repoMock.AssertExpectations(t)
	})

	t.Run("rollback transaction user not found", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		repoMock.On("TransactionGetByUID", ctx, "123").Return(&domain.Transaction{
			UID:          "123",
			UserUID:      "123",
			Amount:       100,
			Currency:     "USD",
			Denomination: 2,
			Type:         domain.TransactionTypeDebit,
		}, nil)

		repoMock.
			On("UserGetByUID", ctx, "123").
			Return(nil, domain.NewError(errorRollbackSource).SetCode(domain.ErrNotFound))

		res, err := service.Rollback(ctx, &domain.ProcessDebitCreditRollbackReq{
			TransactionUID: "123",
		})

		assert.Equal(t, domain.AsError(err).Code, domain.ErrUserNotFound)
		assert.Nil(t, res)

		repoMock.AssertExpectations(t)
	})

	t.Run("rollback transaction type invalid", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		repoMock.On("TransactionGetByUID", ctx, "123").Return(&domain.Transaction{
			UID:          "123",
			UserUID:      "123",
			Amount:       100,
			Currency:     "USD",
			Denomination: 2,
			Type:         domain.TransactionTypeRollback,
		}, nil)

		repoMock.
			On("UserGetByUID", ctx, "123").
			Return(&domain.User{
				UID:  "123",
				Nick: "test",
			}, nil)

		res, err := service.Rollback(ctx, &domain.ProcessDebitCreditRollbackReq{
			TransactionUID: "123",
		})

		assert.Equal(t, domain.AsError(err).Code, domain.ErrInvalidTransactionType)
		assert.Nil(t, res)

		repoMock.AssertExpectations(t)
	})

	t.Run("rollback transaction type invalid", func(t *testing.T) {
		repoMock := &mocks.Repository{}
		service := New(repoMock, logger)

		repoMock.On("TransactionGetByUID", ctx, "123").Return(&domain.Transaction{
			UID:          "123",
			UserUID:      "123",
			Amount:       100,
			Currency:     "USD",
			Denomination: 2,
			Type:         domain.TransactionTypeDebit,
		}, nil)

		repoMock.
			On("UserGetByUID", ctx, "123").
			Return(&domain.User{
				UID:  "123",
				Nick: "test",
			}, nil)

		repoMock.
			On("BalanceIncrementByUserUIDAndCurrency", ctx, "123", "USD", 100).
			Return(nil, domain.NewError(errorRollbackSource).SetCode(domain.ErrIncrement))

		res, err := service.Rollback(ctx, &domain.ProcessDebitCreditRollbackReq{
			TransactionUID: "123",
		})

		assert.Equal(t, domain.AsError(err).Code, domain.ErrRollback)
		assert.Nil(t, res)

		repoMock.AssertExpectations(t)
	})
}
