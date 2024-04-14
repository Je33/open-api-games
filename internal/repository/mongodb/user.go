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
	userTable = "user"

	// errors prefix
	userErrorSource = "[repository.mongodb.user]"
)

type userDB struct {
	UID  string `bson:"uid"`
	Nick string `bson:"nick"`
}

func (mr *Repo) UserGetByUID(ctx context.Context, uid string) (*domain.User, error) {
	var result userDB
	err := mr.db.Collection(userTable).FindOne(ctx, bson.M{"uid": uid}).Decode(&result)
	if err != nil {
		mr.logger.Error("failed to find user", uid, err)
		return nil, domain.NewError(userErrorSource).SetCode(domain.ErrNotFound).Add(err)
	}
	return &domain.User{
		UID:  result.UID,
		Nick: result.Nick,
	}, nil
}

func (mr *Repo) UserCreate(ctx context.Context, user *domain.User) error {
	userDb := userDB{
		UID:  user.UID,
		Nick: user.Nick,
	}

	_, err := mr.db.Collection(userTable).InsertOne(ctx, userDb)
	if err != nil {
		mr.logger.Error("failed to create user", "record", userDb, "error", err)
		return domain.NewError(userErrorSource).SetCode(domain.ErrRepoCreate).Add(err)
	}
	return nil
}

func (mr *Repo) userEnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{"uid", -1}}, Options: options.Index().SetUnique(true)},
	}
	_, err := mr.db.Collection(userTable).Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return domain.NewError(userErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
