package cqs

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
)

// AppError is a query/command hnd error with context
type AppError struct {
	Name   string      `json:"name"`
	Input  interface{} `json:"input"`
	ErrMsg string      `json:"err_msg"`
}

// QueryHandlerMiddleware is a type for decorating QueryHandlers
type QueryHandlerMiddleware func(h QueryHandler) QueryHandler

// NewQueryHndErrorMiddleware is a middleware constructor to log a contextualized query handler error
func NewQueryHndErrorMiddleware(logger *zerolog.Logger) QueryHandlerMiddleware {
	return func(h QueryHandler) QueryHandler {
		return queryHandlerFunc(func(ctx context.Context, q Query) (any, error) {
			result, err := h.Handle(ctx, q)
			if err != nil {
				logAppErr(logger, AppError{
					Name:   q.Name(),
					Input:  q,
					ErrMsg: err.Error(),
				})
				return nil, err
			}

			return result, nil
		})
	}
}

func logAppErr(logger *zerolog.Logger, appErr AppError) {
	b, err := json.Marshal(&appErr)
	if err != nil {
		logger.Err(err).Msgf("something when wrong when trying to marshal app error from %s: %s", err.Error(), appErr.Name)
		return
	}

	logger.Err(err).Msg(string(b))
}

type CommandHandlerMiddleware func(h CommandHandler) CommandHandler

func NewCommandHndErrorMiddleware(logger *zerolog.Logger) CommandHandlerMiddleware {
	return func(h CommandHandler) CommandHandler {
		return CommandHandlerFunc(func(ctx context.Context, q Command) ([]Event, error) {
			result, err := h.Handle(ctx, q)
			if err != nil {
				logAppErr(logger, AppError{
					Name:   q.Name(),
					Input:  q,
					ErrMsg: err.Error(),
				})
				return nil, err
			}

			return result, nil
		})
	}
}
