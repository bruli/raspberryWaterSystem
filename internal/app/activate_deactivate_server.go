package app

import (
	"context"

	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

const ActivateDeactivateServerCmdName = "activateDeactivateServer"

type ActivateDeactivateServerCmd struct {
	Active bool
}

func (a ActivateDeactivateServerCmd) Name() string {
	return ActivateDeactivateServerCmdName
}

type ActivateDeactivateServer struct {
	stRepo StatusRepository
}

func (a ActivateDeactivateServer) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(ActivateDeactivateServerCmd)
	if !ok {
		return nil, cqs.NewInvalidCommandError(ActivateDeactivateServerCmdName, cmd.Name())
	}

	st, err := a.stRepo.Find(ctx)
	if err != nil {
		return nil, err
	}

	switch {
	case co.Active:
		st.Activate()
	default:
		st.Deactivate()
	}
	return nil, a.stRepo.Update(ctx, st)
}

func NewActivateDeactivateServer(stRepo StatusRepository) *ActivateDeactivateServer {
	return &ActivateDeactivateServer{stRepo: stRepo}
}
