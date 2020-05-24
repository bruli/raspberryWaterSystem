package server

import (
	"errors"
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/bruli/raspberryWaterSystem/internal/zone"
)

func checkBadRequestErrors(e error) bool {
	switch {
	case errors.As(e, new(zone.CreateError)):
		return true
	case errors.As(e, new(*zone.InvalidRelay)):
		return true
	case errors.As(e, new(*execution.InvalidCreateData)):
		return true
	case errors.As(e, new(*execution.InvalidCreateExecution)):
		return true
	case errors.As(e, new(*execution.InvalidExecutorData)):
		return true
	}

	return false
}

func checkNotFoundErrors(e error) bool {
	switch {
	case errors.As(e, new(*zone.NotFound)):
		return true
	}
	return false
}
