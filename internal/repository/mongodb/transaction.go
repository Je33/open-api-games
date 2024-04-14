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
	transactionTable = "transaction"

	// errors prefix
	transactionErrorSource = "[repository.mongodb.transaction]"
)

type transactionDB struct {
	UID          string                 `bson:"uid"`
	UserUID      string                 `bson:"userUid"`
	SessionUID   string                 `bson:"sessionUid"`
	Amount       int                    `bson:"amount"`
	Currency     string                 `bson:"currency"`
	Denomination int                    `bson:"denomination"`
	Type         domain.TransactionType `bson:"type"`
}

func (mr *Repo) TransactionGetByUID(ctx context.Context, uid string) (*domain.Transaction, error) {
	var result transactionDB
	err := mr.db.Collection(transactionTable).FindOne(ctx, bson.M{"uid": uid}).Decode(&result)
	if err != nil {
		mr.logger.Error("failed to find transaction", uid, err)
		return nil, domain.NewError(transactionErrorSource).SetCode(domain.ErrNotFound).Add(err)
	}
	return &domain.Transaction{
		UID:          result.UID,
		UserUID:      result.UserUID,
		SessionUID:   result.SessionUID,
		Amount:       result.Amount,
		Currency:     result.Currency,
		Denomination: result.Denomination,
		Type:         result.Type,
	}, nil
}

func (mr *Repo) TransactionCreate(ctx context.Context, transaction *domain.Transaction) error {
	transactionDb := transactionDB{
		UID:          transaction.UID,
		UserUID:      transaction.UserUID,
		SessionUID:   transaction.SessionUID,
		Amount:       transaction.Amount,
		Currency:     transaction.Currency,
		Denomination: transaction.Denomination,
		Type:         transaction.Type,
	}

	_, err := mr.db.Collection(transactionTable).InsertOne(ctx, transactionDb)
	if err != nil {
		mr.logger.Error("failed to create transaction", "record", transactionDb, "error", err)
		return domain.NewError(transactionErrorSource).SetCode(domain.ErrRepoCreate).Add(err)
	}
	return nil
}

func (mr *Repo) transactionEnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{"uid", -1}}, Options: options.Index().SetUnique(true)},
	}
	_, err := mr.db.Collection(transactionTable).Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return domain.NewError(transactionErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
