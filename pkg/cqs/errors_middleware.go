package cqs

import (
	"context"
	"encoding/json"
	"log/slog"
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
func NewQueryHndErrorMiddleware(logger *slog.Logger) QueryHandlerMiddleware {
	return func(h QueryHandler) QueryHandler {
		return queryHandlerFunc(func(ctx context.Context, q Query) (any, error) {
			result, err := h.Handle(ctx, q)
			if err != nil {
				logAppErr(ctx, logger, AppError{
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

func logAppErr(ctx context.Context, logger *slog.Logger, appErr AppError) {
	b, err := json.Marshal(&appErr)
	if err != nil {
		logger.ErrorContext(
			ctx,
			"failed marshaling app error",
			slog.String("error", err.Error()),
			slog.String("app_error", appErr.Name),
		)
		return
	}
	logger.ErrorContext(ctx, string(b), slog.String("app_error", appErr.Name))
}

type CommandHandlerMiddleware func(h CommandHandler) CommandHandler

func NewCommandHndErrorMiddleware(logger *slog.Logger) CommandHandlerMiddleware {
	return func(h CommandHandler) CommandHandler {
		return CommandHandlerFunc(func(ctx context.Context, q Command) ([]Event, error) {
			result, err := h.Handle(ctx, q)
			if err != nil {
				logAppErr(ctx, logger, AppError{
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
