package repository

import (
	"context"
	"log/slog"
	"open-api-games/internal/service/game_processor"

	"open-api-games/internal/repository/mongodb"
)

var (
	// test services interfaces
	_ game_processor.Repository = (*mongodb.Repo)(nil)
)

func NewRepo(ctx context.Context, logger *slog.Logger) (*mongodb.Repo, error) {
	// connect and initialize db
	dbr, err := mongodb.Connect(ctx, logger)
	if err != nil {
		return nil, err
	}

	return dbr, nil
}
