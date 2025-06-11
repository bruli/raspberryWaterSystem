package telegram

import (
	"context"
	"fmt"

	"github.com/bruli/raspberryWaterSystem/internal/app"
	"github.com/bruli/raspberryWaterSystem/internal/domain/weather"
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
)

type weatherCommand struct{}

func (w weatherCommand) CommandName() CommandName {
	return WeatherCommandName
}

type weatherRunner struct {
	qh cqs.QueryHandler
}

func (w weatherRunner) Run(ctx context.Context, chatID int64, msgs *Messages, _ runnerCommand) error {
	result, err := w.qh.Handle(ctx, app.FindWeatherQuery{})
	if err != nil {
		return fmt.Errorf("failed to find weather: %w", err)
	}
	weath, _ := result.(weather.Weather)
	buildMessage(chatID, msgs, fmt.Sprintf("Current temperature: %v *C", weath.Temperature()))
	buildMessage(chatID, msgs, fmt.Sprintf("Current humidity: %v", weath.Humidity()))
	buildMessage(chatID, msgs, fmt.Sprintf("Is raining:  %v", weath.IsRaining()))
	return nil
}

func newWeatherRunner(qh cqs.QueryHandler) *weatherRunner {
	return &weatherRunner{qh: qh}
}
