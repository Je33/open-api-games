package game_processor

import (
	"context"
	"open-api-games/internal/domain"
)

const (
	errorMetadataSource = "[service.game_processor.metadata]"
)

func (s *Service) MetaData(ctx context.Context, req *domain.ProcessMetaDataReq) (*domain.ProcessMetaDataRes, error) {
	session, err := s.repo.SessionGetByUID(ctx, req.GameSessionUID)
	if err != nil {
		return nil, domain.NewError(errorMetadataSource).SetCode(domain.ErrSessionNotFound).Add(err)
	}

	// TODO: implement metadata api

	return &domain.ProcessMetaDataRes{
		Api:  "",
		Data: session.UID,
	}, nil
}
