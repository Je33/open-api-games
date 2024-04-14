package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"open-api-games/internal/domain"
)

const (
	// table name in DB
	balanceTable = "balance"

	// errors prefix
	balanceErrorSource = "[repository.mongodb.balance]"
)

type balanceDB struct {
	UserUID      string `bson:"userUid"`
	Amount       int    `bson:"amount"`
	Currency     string `bson:"currency"`
	Denomination int    `bson:"denomination"`
}

func (mr *Repo) BalanceGetByUserUIDAndCurrency(ctx context.Context, userUID, currency string) (*domain.Balance, error) {
	var result balanceDB
	err := mr.db.Collection(balanceTable).FindOne(ctx, bson.M{"userUid": userUID, "currency": currency}).Decode(&result)
	if err != nil {
		mr.logger.Error("failed to find balance", userUID, currency, err)
		return nil, domain.NewError(balanceErrorSource).SetCode(domain.ErrNotFound).Add(err)
	}
	return &domain.Balance{
		UserUID:      result.UserUID,
		Amount:       result.Amount,
		Currency:     result.Currency,
		Denomination: result.Denomination,
	}, nil
}

func (mr *Repo) BalanceCreate(ctx context.Context, balance *domain.Balance) error {
	balanceDb := balanceDB{
		UserUID:      balance.UserUID,
		Amount:       balance.Amount,
		Currency:     balance.Currency,
		Denomination: balance.Denomination,
	}

	_, err := mr.db.Collection(balanceTable).InsertOne(ctx, balanceDb)
	if err != nil {
		mr.logger.Error("failed to create balance", "document", balanceDb, "error", err)
		return domain.NewError(balanceErrorSource).SetCode(domain.ErrRepoCreate).Add(err)
	}

	return nil
}

func (mr *Repo) BalanceDecrementByUserUIDAndCurrency(ctx context.Context, userUID, currency string, amount int) (*domain.Transaction, error) {
	// use transaction to avoid race condition
	var transactionDb transactionDB
	err := mr.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		var balanceDb balanceDB
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		err = mr.db.Collection(balanceTable).FindOneAndUpdate(
			sessionContext,
			bson.M{"userUid": userUID, "currency": currency, "amount": bson.M{"$gte": amount}},
			bson.M{"$inc": bson.M{"amount": -amount}},
			opts,
		).Decode(&balanceDb)
		if err != nil {
			if errAbort := sessionContext.AbortTransaction(sessionContext); errAbort != nil {
				return errAbort
			}
			return err
		}

		transactionDb = transactionDB{
			UID:          domain.GenUID(),
			UserUID:      balanceDb.UserUID,
			Amount:       amount,
			Currency:     balanceDb.Currency,
			Denomination: balanceDb.Denomination,
			Type:         domain.TransactionTypeDebit,
		}
		_, err = mr.db.Collection(transactionTable).InsertOne(sessionContext, transactionDb)
		if err != nil {
			if errAbort := sessionContext.AbortTransaction(sessionContext); errAbort != nil {
				return errAbort
			}
			return err
		}
		if err = sessionContext.CommitTransaction(sessionContext); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		mr.logger.Error("failed to decrement balance", "userUid", userUID, "currency", currency, "amount", amount, "error", err)
		return nil, domain.NewError(balanceErrorSource).SetCode(domain.ErrDecrement).Add(err)
	}

	return &domain.Transaction{
		UID:          transactionDb.UID,
		UserUID:      transactionDb.UserUID,
		SessionUID:   transactionDb.SessionUID,
		Amount:       transactionDb.Amount,
		Currency:     transactionDb.Currency,
		Denomination: transactionDb.Denomination,
		Type:         transactionDb.Type,
	}, nil
}

func (mr *Repo) BalanceIncrementByUserUIDAndCurrency(ctx context.Context, userUID, currency string, amount int) (*domain.Transaction, error) {
	// use transaction to avoid race condition
	var transactionDb transactionDB
	err := mr.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		var balanceDb balanceDB
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		err = mr.db.Collection(balanceTable).FindOneAndUpdate(
			sessionContext,
			bson.M{"userUid": userUID, "currency": currency},
			bson.M{"$inc": bson.M{"amount": amount}},
			opts,
		).Decode(&balanceDb)
		if err != nil {
			if errAbort := sessionContext.AbortTransaction(sessionContext); errAbort != nil {
				return errAbort
			}
			return err
		}

		transactionDb = transactionDB{
			UID:          domain.GenUID(),
			UserUID:      balanceDb.UserUID,
			Amount:       amount,
			Currency:     balanceDb.Currency,
			Denomination: balanceDb.Denomination,
			Type:         domain.TransactionTypeCredit,
		}
		_, err = mr.db.Collection(transactionTable).InsertOne(sessionContext, transactionDb)
		if err != nil {
			if errAbort := sessionContext.AbortTransaction(sessionContext); errAbort != nil {
				return errAbort
			}
			return err
		}
		if err = sessionContext.CommitTransaction(sessionContext); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		mr.logger.Error("failed to increment balance", "userUid", userUID, "currency", currency, "amount", amount, "error", err)
		return nil, domain.NewError(balanceErrorSource).SetCode(domain.ErrIncrement).Add(err)
	}

	return &domain.Transaction{
		UID:          transactionDb.UID,
		UserUID:      transactionDb.UserUID,
		SessionUID:   transactionDb.SessionUID,
		Amount:       transactionDb.Amount,
		Currency:     transactionDb.Currency,
		Denomination: transactionDb.Denomination,
		Type:         transactionDb.Type,
	}, nil
}

func (mr *Repo) balanceEnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{"userUid", -1}, {"currency", -1}}, Options: options.Index().SetUnique(true)},
	}
	_, err := mr.db.Collection(balanceTable).Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return domain.NewError(balanceErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
