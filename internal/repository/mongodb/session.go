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
	sessionTable = "session"

	// errors prefix
	sessionErrorSource = "[repository.mongodb.session]"
)

type sessionDB struct {
	UID     string `bson:"uid"`
	UserUID string `bson:"userUid"`
}

func (mr *Repo) SessionGetByUID(ctx context.Context, uid string) (*domain.Session, error) {
	var result sessionDB
	err := mr.db.Collection(sessionTable).FindOne(ctx, bson.M{"uid": uid}).Decode(&result)
	if err != nil {
		mr.logger.Error("failed to find session", uid, err)
		return nil, domain.NewError(sessionErrorSource).SetCode(domain.ErrNotFound).Add(err)
	}
	return &domain.Session{
		UID:     result.UID,
		UserUID: result.UserUID,
	}, nil
}

func (mr *Repo) SessionCreate(ctx context.Context, sess *domain.Session) error {
	sessionDb := sessionDB{
		UID:     sess.UID,
		UserUID: sess.UserUID,
	}

	_, err := mr.db.Collection(sessionTable).InsertOne(ctx, sessionDb)
	if err != nil {
		mr.logger.Error("failed to create session", "record", sessionDb, "error", err)
		return domain.NewError(sessionErrorSource).SetCode(domain.ErrRepoCreate).Add(err)
	}
	return nil
}

func (mr *Repo) sessionEnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{"uid", -1}}, Options: options.Index().SetUnique(true)},
	}
	_, err := mr.db.Collection(sessionTable).Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return domain.NewError(sessionErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}
	return nil
}
