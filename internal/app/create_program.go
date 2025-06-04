package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"github.com/bruli/raspberryWaterSystem/pkg/vo"
)

const (
	CreateDailyProgramCommandName = "createDailyProgram"
	CreateOddProgramCommandName   = "createOddProgram"
	CreateEvenProgramCommandName  = "createEvenProgram"
)

type CreateDailyProgramCommand struct {
	Program *program.Program
}

func (c CreateDailyProgramCommand) Name() string {
	return CreateDailyProgramCommandName
}

type CreateOddProgramCommand struct {
	Program *program.Program
}

func (c CreateOddProgramCommand) Name() string {
	return CreateOddProgramCommandName
}

type CreateEvenProgramCommand struct {
	Program *program.Program
}

func (c CreateEvenProgramCommand) Name() string {
	return CreateEvenProgramCommandName
}

type CreateDailyProgram struct {
	CreateProgram
}

func (c CreateDailyProgram) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(CreateDailyProgramCommand)
	if !ok {
		return nil, cqs.NewInvalidCommandError(CreateDailyProgramCommandName, cmd.Name())
	}
	return nil, c.Create(ctx, co.Program)
}

func NewCreateDailyProgram(programRepo ProgramRepository, zonesRepo ZoneRepository) *CreateDailyProgram {
	return &CreateDailyProgram{
		CreateProgram: CreateProgram{programRepo: programRepo, zonesRepo: zonesRepo},
	}
}

type CreateOddProgram struct {
	CreateProgram
}

func (c CreateOddProgram) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(CreateOddProgramCommand)
	if !ok {
		return nil, cqs.NewInvalidCommandError(CreateOddProgramCommandName, cmd.Name())
	}
	return nil, c.Create(ctx, co.Program)
}

func NewCreateOddProgram(programRepo ProgramRepository, zonesRepo ZoneRepository) *CreateOddProgram {
	return &CreateOddProgram{
		CreateProgram: CreateProgram{programRepo: programRepo, zonesRepo: zonesRepo},
	}
}

type CreateEvenProgram struct {
	CreateProgram
}

func (c CreateEvenProgram) Handle(ctx context.Context, cmd cqs.Command) ([]cqs.Event, error) {
	co, ok := cmd.(CreateEvenProgramCommand)
	if !ok {
		return nil, cqs.NewInvalidCommandError(CreateEvenProgramCommandName, cmd.Name())
	}
	return nil, c.Create(ctx, co.Program)
}

func NewCreateEvenProgram(programRepo ProgramRepository, zonesRepo ZoneRepository) *CreateEvenProgram {
	return &CreateEvenProgram{
		CreateProgram: CreateProgram{programRepo: programRepo, zonesRepo: zonesRepo},
	}
}

type CreateProgramError struct {
	msg string
}

func (c CreateProgramError) Error() string {
	return fmt.Sprintf("failed to create program: %s", c.msg)
}

type CreateProgram struct {
	programRepo ProgramRepository
	zonesRepo   ZoneRepository
}

func (p CreateProgram) Create(ctx context.Context, program *program.Program) error {
	hour := program.Hour()
	_, err := p.programRepo.FindByHour(ctx, &hour)
	switch {
	case err == nil:
		return CreateProgramError{msg: fmt.Sprintf("a program with hour %s, already exists", hour.String())}
	case errors.As(err, &vo.NotFoundError{}):
		if err = p.checkZone(ctx, program.Executions()); err != nil {
			return err
		}
		return p.programRepo.Save(ctx, program)
	default:
		return err
	}
}

func (p CreateProgram) checkZone(ctx context.Context, executions []program.Execution) error {
	for _, e := range executions {
		for _, z := range e.Zones() {
			if _, err := p.zonesRepo.FindByID(ctx, z); err != nil {
				switch {
				case errors.As(err, &vo.NotFoundError{}):
					return CreateProgramError{msg: fmt.Sprintf("a zone with id %s, not found", z)}
				default:
					return err
				}
			}
		}
	}
	return nil
}
