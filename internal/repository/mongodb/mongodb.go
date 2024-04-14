package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"log/slog"
	"open-api-games/internal/config"
	"open-api-games/internal/domain"
	"time"
)

const (
	mongodbErrorSource = "[repository.mongodb]"
)

var (
	heartbeatInterval = 5 * time.Second
	reconnectInterval = 5 * time.Second
)

type Repo struct {
	db     *mongo.Database
	logger *slog.Logger
}

func Connect(ctx context.Context, logger *slog.Logger) (*Repo, error) {
	cfg := config.Get(logger)

	if cfg.MongoURI == "" {
		logger.Error("missing mongo uri", "config", cfg)
		return nil, domain.NewError(mongodbErrorSource).SetCode(domain.ErrConfig)
	}

	cs, err := connstring.ParseAndValidate(cfg.MongoURI)
	if err != nil {
		logger.Error("failed to parse mongo uri", "uri", cfg.MongoURI, "error", err)
		return nil, domain.NewError(mongodbErrorSource).SetCode(domain.ErrConfig).Add(err)
	}
	if cs.Database == "" {
		logger.Error("missing database name", "connstring", cs)
		return nil, domain.NewError(mongodbErrorSource).SetCode(domain.ErrConfig)
	}

	repo := &Repo{
		logger: logger,
	}

	// Debug
	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			logger.Debug("database command", "command", evt.Command.String())
		},
	}

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(cfg.MongoURI).
		SetServerAPIOptions(serverAPI).
		SetMonitor(cmdMonitor)

	client, err := repo.connect(ctx, opts, logger)
	if err != nil {
		return nil, domain.NewError(mongodbErrorSource).SetCode(domain.ErrConnect).Add(err)
	}

	repo.db = client.Database(cs.Database)

	// Send a ping to confirm a successful connection
	if err = repo.ping(ctx); err != nil {
		return nil, domain.NewError(mongodbErrorSource).SetCode(domain.ErrConnect).Add(err)
	}

	go repo.KeepAlive(ctx, logger)

	return repo, nil
}

func (mr *Repo) Close(ctx context.Context) error {
	return mr.db.Client().Disconnect(ctx)
}

func (mr *Repo) KeepAlive(ctx context.Context, logger *slog.Logger) {
	ticker := time.NewTicker(heartbeatInterval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
		case <-ticker.C:
			err := mr.ping(ctx)
			if err != nil {
				ticker.Stop()
				mr.reconnect(ctx, logger)
			}
		}
	}
}

func (mr *Repo) DropTest(ctx context.Context, logger *slog.Logger) error {
	cfg := config.Get(logger)
	if cfg.Env != "dev" {
		logger.Info("not in test environment", "env", cfg.Env)
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrConfig)
	}
	return mr.db.Drop(ctx)
}

func (mr *Repo) ping(ctx context.Context) error {
	// Send a ping to check connection
	var result bson.M
	if err := mr.db.Client().Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		mr.logger.Error("failed to ping mongo", "error", err)
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrConnect).Add(err)
	}
	return nil
}

func (mr *Repo) connect(ctx context.Context, opts *options.ClientOptions, logger *slog.Logger) (*mongo.Client, error) {
	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		logger.Error("failed to connect to mongo", "error", err)
		return nil, domain.NewError(mongodbErrorSource).SetCode(domain.ErrConnect).Add(err)
	}

	return client, nil
}

func (mr *Repo) reconnect(ctx context.Context, logger *slog.Logger) {
	attempts := 0
	for {
		attempts++
		repo, err := Connect(ctx, logger)
		if err != nil {
			logger.Error("failed to reconnect to mongo", "attempts", attempts, "error", err)
			time.Sleep(reconnectInterval)
		} else {
			mr.db = repo.db
			break
		}
	}
}

func (mr *Repo) EnsureIndexes(ctx context.Context) error {
	err := mr.userEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	err = mr.balanceEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	err = mr.currencyEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	err = mr.sessionEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	err = mr.transactionEnsureIndexes(ctx)
	if err != nil {
		return domain.NewError(mongodbErrorSource).SetCode(domain.ErrRepoInit).Add(err)
	}

	// TODO: Add indexes

	return nil
}
