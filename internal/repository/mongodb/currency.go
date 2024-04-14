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
	currencyTable = "currency"

	// errors prefix
	currencyErrorSource = "[repository.mongodb.currency]"
)

type currencyDB struct {
	Code         string `bson:"code"`
	Denomination int    `bson:"denomination"`
}

func (mr *Repo) CurrencyGetByCode(ctx context.Context, code string) (*domain.Currency, error) {
	var result currencyDB
	err := mr.db.Collection(currencyTable).FindOne(ctx, bson.M{"code": code}).Decode(&result)
	if err != nil {
		mr.logger.Error("failed to find currency", "code", code, "error", err)
		return nil, domain.NewError(currencyErrorSource).SetCode(domain.ErrNotFound).Add(err)
	}
	return &domain.Currency{
		Code:         result.Code,
		Denomination: result.Denomination,
	}, nil
}

func (mr *Repo) CurrencyCreate(ctx context.Context, cur *domain.Currency) error {
	currencyDb := currencyDB{
		Code:         cur.Code,
		Denomination: cur.Denomination,
	}

	_, err := mr.db.Collection(currencyTable).InsertOne(ctx, currencyDb)
	if err != nil {
		mr.logger.Error("failed to create session", "document", currencyDb, "error", err)
		return domain.NewError(currencyErrorSource).SetCode(domain.ErrRepoCreate).Add(err)
	}
	return nil
}

func (mr *Repo) currencyEnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{"code", -1}}, Options: options.Index().SetUnique(true)},
	}
	_, err := mr.db.Collection(currencyTable).Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return domain.NewError(currencyErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
