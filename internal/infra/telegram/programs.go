package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/cqs"
	"github.com/bruli/raspberryWaterSystem/internal/domain/program"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type programsCommand struct{}

func (p programsCommand) CommandName() CommandName {
	return ProgramsCommandName
}

type programsRunner struct {
	qh     cqs.QueryHandler
	tracer trace.Tracer
}

func (p programsRunner) Run(ctx context.Context, chatID int64, msgs *Messages, _ runnerCommand) error {
	ctx, span := p.tracer.Start(ctx, "programsRunner.Run")
	defer span.End()
	result, err := p.qh.Handle(ctx, app.FindAllProgramsQuery{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("failed to find programs: %w", err)
	}
	programs, _ := result.(app.AllPrograms)
	p.checkPrograms(programs.Daily, "Daily", chatID, msgs)
	p.checkPrograms(programs.Odd, "Odd", chatID, msgs)
	p.checkPrograms(programs.Even, "Even", chatID, msgs)
	p.checkWeeklyPrograms(programs.Weekly, chatID, msgs)
	p.checkTemperaturePrograms(programs.Temperature, chatID, msgs)
	span.SetStatus(codes.Ok, "programs found")
	return nil
}

func (p programsRunner) checkPrograms(programs []program.Program, prefix string, chatID int64, msgs *Messages) {
	if len(programs) == 0 {
		buildMessage(chatID, msgs, fmt.Sprintf("%s: No programs", prefix))
		return
	}
	for _, daily := range programs {
		var (
			zones   string
			seconds int
		)
		for _, exec := range daily.Executions() {
			zones = strings.Join(exec.Zones(), ", ")
			seconds = exec.Seconds().Int()
			buildMessage(chatID, msgs, fmt.Sprintf("%s: - Hour: %s, Zones: %s, Seconds: %v", prefix, daily.Hour().String(), zones, seconds))
		}
	}
}

func (p programsRunner) checkWeeklyPrograms(weekly []program.Weekly, id int64, msgs *Messages) {
	if len(weekly) == 0 {
		buildMessage(id, msgs, "Weekly: No programs")
		return
	}
	for _, w := range weekly {
		day := w.WeekDay().String()
		for _, pr := range w.Programs() {
			var (
				zones   string
				seconds int
			)
			for _, exe := range pr.Executions() {
				zones = strings.Join(exe.Zones(), ", ")
				seconds = exe.Seconds().Int()
				buildMessage(id, msgs, fmt.Sprintf("Weekly: - Day: %s, Hour: %s, Zones: %s, Seconds: %v", day, pr.Hour().String(), zones, seconds))
			}
		}
	}
}

func (p programsRunner) checkTemperaturePrograms(temperature []program.Temperature, id int64, msgs *Messages) {
	if len(temperature) == 0 {
		buildMessage(id, msgs, "Temperature: No programs")
		return
	}
	for _, t := range temperature {
		temp := t.Temperature()
		for _, pr := range t.Programs() {
			var (
				zones   string
				seconds int
			)
			for _, exe := range pr.Executions() {
				zones = strings.Join(exe.Zones(), ", ")
				seconds = exe.Seconds().Int()
				buildMessage(id, msgs, fmt.Sprintf("Temperature: - Temperature: %v, Hour: %s, Zones: %s, Seconds: %v", temp, pr.Hour().String(), zones, seconds))
			}
		}
	}
}

func newProgramsRunner(qh cqs.QueryHandler, tracer trace.Tracer) *programsRunner {
	return &programsRunner{qh: qh, tracer: tracer}
}
