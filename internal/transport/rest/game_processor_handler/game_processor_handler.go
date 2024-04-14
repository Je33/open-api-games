package game_processor_handler

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"log/slog"
	"open-api-games/internal/config"
	"open-api-games/internal/domain"
	"open-api-games/internal/transport/rest/model"
)

const (
	errorSource = "[transport.rest.game_processor_handler]"
)

type GameProcessorService interface {
	Balance(ctx context.Context, req *domain.ProcessBalanceReq) (*domain.ProcessBalanceRes, error)
	Debit(ctx context.Context, req *domain.ProcessDebitCreditRollbackReq) (*domain.ProcessDebitCreditRollbackRes, error)
	Credit(ctx context.Context, req *domain.ProcessDebitCreditRollbackReq) (*domain.ProcessDebitCreditRollbackRes, error)
	Rollback(ctx context.Context, req *domain.ProcessDebitCreditRollbackReq) (*domain.ProcessDebitCreditRollbackRes, error)
	MetaData(ctx context.Context, req *domain.ProcessMetaDataReq) (*domain.ProcessMetaDataRes, error)
}

type Handler struct {
	gameProcessorService GameProcessorService
	logger               *slog.Logger
}

func New(gameProcessorService GameProcessorService, logger *slog.Logger) *Handler {
	return &Handler{
		gameProcessorService: gameProcessorService,
		logger:               logger,
	}
}

func (h *Handler) Process(c echo.Context) error {
	ctx := c.Request().Context()

	b, err := io.ReadAll(c.Request().Body)
	h.logger.Debug("request processing", "path", c.Request().URL.Path, "body", string(b), "error", err)

	apiCommand := &model.ProcessCommand{}
	err = json.Unmarshal(b, apiCommand)
	if err != nil {
		h.logger.Error("error parsing request api command", "error", err)
		return c.JSON(400, makeError[*model.ProcessBalanceRes](apiCommand.Api, domain.AsError(err).Code))
	}
	if !apiCommand.Api.IsValid() {
		h.logger.Error("invalid request api command", "api", apiCommand.Api)
		return c.JSON(400, makeError[*model.ProcessBalanceRes](apiCommand.Api, domain.AsError(err).Code))
	}

	switch apiCommand.Api {
	case model.ProcessApiCommandBalance:
		req := &model.ProcessReq[*model.ProcessBalanceReq]{}

		err = json.Unmarshal(b, req)
		if err != nil {
			h.logger.Error("error parsing request", "error", err)
			return c.JSON(400, makeError[*model.ProcessBalanceRes](apiCommand.Api, domain.AsError(err).Code))
		}

		resp, err := h.gameProcessorService.Balance(ctx, h.balanceFromTransport(req.Data))
		if err != nil {
			h.logger.Error("error processing request", "error", err)
			return c.JSON(500, makeError[*model.ProcessBalanceRes](apiCommand.Api, domain.AsError(err).Code))
		}

		return c.JSON(200, model.ProcessRes[*model.ProcessBalanceRes]{
			Api: model.ProcessApiCommandBalance,
			Data: &model.ProcessBalanceRes{
				UserUID:      resp.UserUID,
				UserNick:     resp.UserNick,
				Amount:       resp.Amount,
				Currency:     resp.Currency,
				Denomination: resp.Denomination,
				MaxWin:       resp.MaxWin,
				JpKey:        resp.JpKey,
			},
			IsSuccess: true,
			Error:     "",
			ErrorMsg:  domain.ErrNone,
		})

	case model.ProcessApiCommandDebit, model.ProcessApiCommandCredit, model.ProcessApiCommandRollback:
		req := &model.ProcessReq[*model.ProcessDebitCreditRollbackReq]{}

		err = json.Unmarshal(b, req)
		if err != nil {
			h.logger.Error("error parsing request", err)
			return c.JSON(400, makeError[*model.ProcessDebitCreditRollbackRes](apiCommand.Api, domain.AsError(err).Code))
		}

		var resp *domain.ProcessDebitCreditRollbackRes

		switch apiCommand.Api {
		case model.ProcessApiCommandDebit:
			resp, err = h.gameProcessorService.Debit(ctx, h.debitCreditRollbackFromTransport(req.Data))
		case model.ProcessApiCommandCredit:
			resp, err = h.gameProcessorService.Credit(ctx, h.debitCreditRollbackFromTransport(req.Data))
		case model.ProcessApiCommandRollback:
			resp, err = h.gameProcessorService.Rollback(ctx, h.debitCreditRollbackFromTransport(req.Data))
		}
		if err != nil {
			h.logger.Error("error processing request", err)
			return c.JSON(500, makeError[*model.ProcessDebitCreditRollbackRes](apiCommand.Api, domain.AsError(err).Code))
		}

		return c.JSON(200, model.ProcessRes[*model.ProcessDebitCreditRollbackRes]{
			Api: model.ProcessApiCommandBalance,
			Data: &model.ProcessDebitCreditRollbackRes{
				TransactionUID: resp.TransactionUID,
				UserNick:       resp.UserNick,
				Amount:         resp.Amount,
				Currency:       resp.Currency,
				Denomination:   resp.Denomination,
				MaxWin:         resp.MaxWin,
			},
			IsSuccess: true,
			Error:     "",
			ErrorMsg:  domain.ErrNone,
		})

	case model.ProcessApiCommandMetaData:
		req := &model.ProcessReq[*model.ProcessMetaDataReq]{}

		err = json.Unmarshal(b, req)
		if err != nil {
			h.logger.Error("error parsing request", err)
			return c.JSON(400, makeError[*model.ProcessMetaDataRes](apiCommand.Api, domain.AsError(err).Code))
		}

		resp, err := h.gameProcessorService.MetaData(ctx, h.metaDataFromTransport(req.Data))
		if err != nil {
			h.logger.Error("error processing request", err)
			return c.JSON(500, makeError[*model.ProcessMetaDataRes](apiCommand.Api, domain.AsError(err).Code))
		}

		return c.JSON(200, model.ProcessRes[*model.ProcessMetaDataRes]{
			Api: model.ProcessApiCommandMetaData,
			Data: &model.ProcessMetaDataRes{
				Api:  req.Data.Api,
				Data: resp.Data,
			},
			IsSuccess: true,
			Error:     "",
			ErrorMsg:  domain.ErrNone,
		})

	default:
		h.logger.Error("invalid request api command", "api", apiCommand.Api)
		return c.JSON(400, makeError[*model.ProcessMetaDataRes](apiCommand.Api, domain.ErrInvalidApiCommand))
	}
}

func (h *Handler) CheckSign(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := config.Get(h.logger)

		bytesToHash, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(400, makeError[*model.ProcessMetaDataRes](model.ProcessApiCommandMetaData, domain.ErrReadBody))
		}

		c.Request().Body = io.NopCloser(bytes.NewBuffer(bytesToHash))

		apiCommand := &model.ProcessCommand{}
		err = json.Unmarshal(bytesToHash, apiCommand)
		if err != nil {
			return c.JSON(400, makeError[*model.ProcessMetaDataRes](model.ProcessApiCommandMetaData, domain.ErrReadBody))
		}

		headerSign := c.Request().Header.Get("Sign")
		if len(headerSign) == 0 {
			return c.JSON(400, makeError[*model.ProcessMetaDataRes](apiCommand.Api, domain.ErrSignEmpty))
		}

		stringToHash := string(bytesToHash) + cfg.ApiKey
		hash := md5.Sum([]byte(stringToHash))

		if hex.EncodeToString(hash[:]) != headerSign {
			return c.JSON(400, makeError[*model.ProcessMetaDataRes](apiCommand.Api, domain.ErrSignInvalid))
		}

		return next(c)
	}
}

func makeError[T model.ProcessApiResData](api model.ProcessApiCommand, error string) model.ProcessRes[T] {
	return model.ProcessRes[T]{
		Api:       api,
		Data:      nil,
		IsSuccess: false,
		Error:     error,
		ErrorMsg:  error,
	}
}

func (h *Handler) balanceFromTransport(req *model.ProcessBalanceReq) *domain.ProcessBalanceReq {
	if req == nil {
		return nil
	}
	return &domain.ProcessBalanceReq{
		GameSessionUID: req.GameSessionUID,
		Currency:       req.Currency,
	}
}

func (h *Handler) debitCreditRollbackFromTransport(req *model.ProcessDebitCreditRollbackReq) *domain.ProcessDebitCreditRollbackReq {
	if req == nil {
		return nil
	}
	return &domain.ProcessDebitCreditRollbackReq{
		TransactionUID: req.TransactionUID,
		GameSessionUID: req.GameSessionUID,
		UserUID:        req.UserUID,
		UserNick:       req.UserNick,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Denomination:   req.Denomination,
		MaxWin:         req.MaxWin,
		JpKey:          req.JpKey,
		SpinMeta:       req.SpinMeta,
		BetMeta:        req.BetMeta,
	}
}

func (h *Handler) metaDataFromTransport(req *model.ProcessMetaDataReq) *domain.ProcessMetaDataReq {
	if req == nil {
		return nil
	}
	return &domain.ProcessMetaDataReq{
		UserUID:        req.UserUID,
		GameSessionUID: req.GameSessionUID,
		Currency:       req.Currency,
		Api:            domain.ProcessApiDataApi(req.Api),
		Data: domain.ProcessApiDataData{
			BetUID: req.Data.BetId,
		},
	}
}
